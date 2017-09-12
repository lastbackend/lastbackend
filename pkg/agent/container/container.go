//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package container

import (
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"io"
	"net/http"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const buffer_size = 1024

func Logs(cid string, follow bool, stream io.Writer, done chan bool) error {

	var (
		stdout   = bool(true)
		stderr   = bool(true)
		cri      = context.Get().GetCri()
		buffer   = make([]byte, buffer_size)
		doneChan = make(chan bool, 1)
	)

	req, err := cri.ContainerLogs(context.Get(), cid, stdout, stderr, follow)
	if err != nil {
		log.Errorf("Error get logs stream %s", err)
		return err
	}
	defer req.Close()

	go func() {
		for {
			select {
			case <-doneChan:
				req.Close()
				return
			default:
				n, err := req.Read(buffer)
				if err != nil {
					log.Errorf("Error read bytes from stream %s", err)
					req.Close()
					return
				}

				_, err = func(p []byte) (n int, err error) {
					n, err = stream.Write(p)
					if err != nil {
						log.Errorf("Error write bytes to stream %s", err)
						return n, err
					}
					if f, ok := stream.(http.Flusher); ok {
						f.Flush()
					}
					return n, nil
				}(buffer[0:n])
				if err != nil {
					log.Errorf("Error written to stream %s", err)
					return
				}

				for i := 0; i < n; i++ {
					buffer[i] = 0
				}
			}
		}
	}()

	<-done

	close(doneChan)

	return nil
}
