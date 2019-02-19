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
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/spf13/viper"

	"github.com/lastbackend/lastbackend/pkg/util/proxy"
	"net/http"
)

const (
	logPrefix = "exporter:logger"
	logLevel  = 3
)

type Logger struct {
	server      *proxy.Server
	connections map[string]map[http.ResponseWriter]bool
	storage     *Storage
}

func (l *Logger) Listen() error {
	return l.server.Listen(l.Handle)
}

func (l *Logger) Handle(msg types.ProxyMessage) error {

	m := types.LogMessage{}
	if err := json.Unmarshal(msg.Line, &m); err != nil {
		_ = fmt.Errorf("%s:>unmarshal json: %s", logPrefix, err.Error())
		return nil
	}

	fmt.Println(m.Selflink, " ", m.Data)

	pod := types.PodSelfLink{}
	if err := pod.Parse(m.Selflink); err != nil {
		return nil
	}

	k, parent := pod.Parent()

	var stream *File

	switch k {
	case types.KindDeployment:
		kind, p := parent.Parent()
		stream = l.storage.GetStream(kind, p.String())
	case types.KindTask:
		stream = l.storage.GetStream(k, parent.String())
	}

	if stream == nil {
		log.Errorf("stream is nil")
		return nil
	}

	if err := stream.Write(string(msg.Line)); err != nil {
		return nil
	}

	return nil
}

func (l *Logger) Stream(ctx context.Context, kind, selflink string, writer http.ResponseWriter) error {

	var (
		data = []byte{}
		buf  = bytes.NewBuffer(data)
	)

	err := l.storage.GetStream(kind, selflink).Tail(0, false, buf)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(data)
	r := bufio.NewReader(reader)

	for {

		select {
		case <-ctx.Done():
			return nil
		default:
			var b = []byte{}

			_, err = r.Read(b)
			if err != nil {
				return err
			}

			if _, err := writer.Write(b); err != nil {
				return err
			}

			if f, ok := writer.(http.Flusher); ok {
				f.Flush()
			}
		}
	}

}

func NewLogger() (*Logger, error) {

	var (
		err  error
		root = viper.GetString("exporter.dir")
	)

	l := new(Logger)
	l.connections = make(map[string]map[http.ResponseWriter]bool, 0)
	l.server, err = proxy.NewServer(proxy.DefaultServer)
	l.storage = NewStorage(root)
	if err != nil {
		return nil, err
	}

	return l, nil
}
