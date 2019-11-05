//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package proxy

import (
	"encoding/binary"
	protoio "github.com/gogo/protobuf/io"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/tools/log"
	"io"
	"net"
	"time"
)

type Conn struct {
	net.Conn
	name        string
	IdleTimeout time.Duration
	done        chan bool
	error       chan string
	writer      protoio.WriteCloser
}

func (c *Conn) Write(p []byte) (int, error) {
	return c.Conn.Write(p)
}

func (c *Conn) Send(data []byte) error {

	msg := new(types.ProxyMessage)
	msg.Type = KindMSG
	msg.Partial = false
	msg.Source = c.name
	msg.Line = data
	msg.TimeNano = time.Now().Unix()

	if err := c.writer.WriteMsg(msg); err != nil {
		return err
	}

	return nil
}

func (c *Conn) Handle(handler Handler) {

	dec := protoio.NewUint32DelimitedReader(c, binary.BigEndian, 1e6)
	defer func() {
		_ = dec.Close()
		_ = c.Close()
		close(c.error)
		close(c.done)
	}()

	go func() {
		for {
			var msg types.ProxyMessage

			err := dec.ReadMsg(&msg)
			if err != nil {
				if err == io.EOF {
					log.Debug("shutting down logger goroutine due to file EOF")
					c.done <- true
					return
				} else {
					log.Warn("consume: error reading encoded message, trying to continue")
					dec = protoio.NewUint32DelimitedReader(c, binary.BigEndian, 1e6)
					continue
				}
			}

			switch msg.Type {
			case KindMSG:
				if handler != nil {
					if err := handler(msg); err != nil {
						log.Debug("msg handle err")
						c.done <- true
					}
				}
			}

			msg.Reset()
		}

	}()

	for {
		select {
		case e := <-c.error:
			log.Errorf(e)
			return
		case <-c.done:
			return
		}
	}

}

func (c *Conn) Ping() error {

	msg := new(types.ProxyMessage)
	msg.Type = KindPing
	msg.Partial = false
	msg.Source = c.name
	msg.Line = []byte{}
	msg.TimeNano = time.Now().Unix()

	if err := c.writer.WriteMsg(msg); err != nil {
		return err
	}

	msg.Reset()

	return nil
}

func (c *Conn) Pong() error {

	msg := new(types.ProxyMessage)
	msg.Type = KindPong
	msg.Partial = false
	msg.Source = c.name
	msg.TimeNano = time.Now().Unix()

	if err := c.writer.WriteMsg(msg); err != nil {
		return err
	}

	msg.Reset()

	return nil
}
