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

package proxy

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	protoio "github.com/gogo/protobuf/io"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
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

	msg := new(models.ProxyMessage)
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
			var msg models.ProxyMessage

			err := dec.ReadMsg(&msg)
			if err != nil {
				if err == io.EOF {
					c.done <- true
					return
				} else {
					dec = protoio.NewUint32DelimitedReader(c, binary.BigEndian, 1e6)
					continue
				}
			}

			switch msg.Type {
			case KindMSG:
				if handler != nil {
					if err := handler(msg); err != nil {
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
			fmt.Println(e)
			return
		case <-c.done:
			return
		}
	}

}

func (c *Conn) Ping() error {

	msg := new(models.ProxyMessage)
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

	msg := new(models.ProxyMessage)
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
