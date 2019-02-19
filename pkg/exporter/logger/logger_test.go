//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/util/proxy"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {

	viper.Set("exporter.dir", "/var/log/lastbackend")

	t.Log("start logger")
	logger, err := NewLogger()
	if !assert.NoError(t, err, "can not create logger") {
		return
	}

	go func() {
		fmt.Println("start logger listen")
		if err := logger.Listen(); err != nil {
			assert.NoError(t, err, "logger listen error")
			return
		}
		fmt.Println("stop logger listen")
	}()

	<-time.NewTimer(time.Second).C
	fmt.Println("start send messages")

	cl := proxy.NewClient("test", types.EmptyString, nil)
	if !assert.NotNil(t, cl, "client can not be nil") {
		return
	}

	var (
		count  = 10
		done   = make(chan bool)
		tail   = make(chan bool)
		sl     = types.NewTaskSelfLink("ns", "job", "task").String()
		psl, _ = types.NewPodSelfLink(types.KindTask, sl, "pod")
	)

	go func() {

		var stream *File

		for {
			stream = logger.storage.Collection[types.KindTask][sl]
			if stream != nil {
				break
			}
		}

		var (
			data = []byte{}
			buf  = bytes.NewBuffer(data)
		)

		err = stream.Tail(count, false, buf)
		if err != nil {
			t.Error(err.Error())
			return
		}

		var l = 0
		for {

			if l >= count {
				tail <- true
				break
			}

			l++
			var b = []byte{}
			_, err = buf.Read(b)
			if err != nil {
				t.Error(err.Error())
				return
			}
			fmt.Println("tail:>", string(b))
		}
	}()

	go func() {
		var i = 0

		for {

			if i >= count {
				fmt.Println("stop sending messages")
				done <- true
				return
			}

			fmt.Printf("send message: %d\n", i)

			log := types.LogMessage{
				Selflink: psl.String(),
				Data:     fmt.Sprintf("log: %d", i),
			}

			i++

			b, err := json.Marshal(log)
			if !assert.NoError(t, err, "can not marshal log message") {
				return
			}

			if err := cl.Send(b); err != nil {
				assert.NoError(t, err, "logger listen error")
				return
			}

			<-time.NewTimer(time.Second).C
		}
	}()

	<-done

	_, ok := logger.storage.Collection[types.KindTask][sl]
	if !ok {
		t.Error("task storage should exists")
		return
	}

	lines, err := logger.storage.Collection[types.KindTask][sl].ReadLines(0, false)
	if !assert.NoError(t, err, "read lines err") {
		return
	}

	if !assert.Equal(t, count, len(lines), "should read all lines") {
		return
	}

	for _, l := range lines {
		fmt.Println(l)
	}

	<-tail

}
