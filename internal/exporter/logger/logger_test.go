//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/lastbackend/lastbackend/internal/pkg/models"
//	"github.com/lastbackend/lastbackend/internal/util/proxy"
//	"github.com/stretchr/testify/assert"
//	"testing"
//	"time"
//)
//
//func TestNewLogger(t *testing.T) {
//
//	t.Log("start logger")
//
//	opts := new(LoggerOpts)
//	opts.Workdir = "/tmp/log/lastbackend"
//	opts.Host = "127.0.0.1"
//	opts.Port = 2963
//
//	l, err := New(opts)
//	if !assert.NoError(t, err, "can not create logger") {
//		return
//	}
//
//	go func() {
//		if err := l.Listen(); err != nil {
//			assert.NoError(t, err, "logger listen error")
//			return
//		}
//	}()
//
//	<-time.NewTimer(time.Second).C
//
//	cl := proxy.NewClient("test", models.EmptyString, nil)
//	if !assert.NotNil(t, cl, "client can not be nil") {
//		return
//	}
//
//	var (
//		ctx, cf = context.WithCancel(context.Background())
//		count   = 500
//		lines   = 100
//		total   = 0
//		done    = make(chan bool)
//		sl      = models.NewTaskSelfLink("ns", "job", "task").String()
//		psl, _  = models.NewPodSelfLink(models.KindTask, sl, "pod")
//	)
//
//	<-time.NewTimer(time.Millisecond * 100).C
//
//	var i = 0
//	for {
//
//		if i >= lines {
//			break
//		}
//
//		i++
//		log := models.LogMessage{
//			Selflink: psl.String(),
//			Data:     fmt.Sprintf("stored log: %d", i),
//		}
//
//		b, err := json.Marshal(log)
//		if !assert.NoError(t, err, "can not marshal log message") {
//			return
//		}
//
//		if err := cl.Send(b); err != nil {
//			assert.NoError(t, err, "logger send message error")
//			return
//		}
//
//		<-time.NewTimer(time.Millisecond * 10).C
//	}
//
//	go func() {
//
//		var stream *File
//
//		for {
//			stream, err = l.storage.GetStream(models.KindTask, sl, false)
//			if err != nil {
//				t.Error(err.Error())
//				break
//			}
//			if stream != nil {
//				break
//			}
//		}
//
//		var (
//			lch = make(chan string)
//		)
//
//		go func() {
//			err = stream.Read(ctx, lines, true, lch)
//			if err != nil {
//				t.Error(err.Error())
//				return
//			}
//		}()
//
//		for range lch {
//			total++
//
//			if total == count+lines {
//				done <- true
//				break
//			}
//
//		}
//	}()
//
//	go func() {
//		var i = 0
//		for {
//
//			if i >= count {
//				break
//			}
//
//			i++
//			log := models.LogMessage{
//				Selflink: psl.String(),
//				Data:     fmt.Sprintf("realtime log: %d", i),
//			}
//
//			b, err := json.Marshal(log)
//			if !assert.NoError(t, err, "can not marshal log message") {
//				return
//			}
//
//			if err := cl.Send(b); err != nil {
//				assert.NoError(t, err, "logger listen error")
//				return
//			}
//
//			<-time.NewTimer(time.Millisecond * 10).C
//		}
//	}()
//
//	for {
//		select {
//		case <-done:
//			cf()
//			if !assert.Equal(t, count+lines, total, "should read all lines") {
//				return
//			}
//
//			return
//		case <-time.NewTimer(time.Second * 30).C:
//			cf()
//			t.Errorf("messages not recevied")
//			return
//		}
//	}
//
//}
