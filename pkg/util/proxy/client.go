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

package proxy

import (
	"encoding/binary"
	protoio "github.com/gogo/protobuf/io"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"io"
	"net"
	"time"
)

const defaultAddr = "127.0.0.1:2963"

type Client struct {
	io.Writer
	name   string
	conn   net.Conn
	done   chan bool
	writer protoio.WriteCloser
	active bool
}

func (p *Client) Connect(addr string, cl Handler) error {

	if addr == types.EmptyString {
		addr = defaultAddr
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	defer func() { _ = conn.Close() }()

	err = conn.(*net.TCPConn).SetKeepAlive(true)
	if err != nil {
		return err
	}

	p.conn = conn
	p.writer = protoio.NewUint32DelimitedWriter(p.conn, binary.BigEndian)
	p.active = true

	defer func() {
		log.Debugf("closing connection to %v", conn.RemoteAddr())
		_ = conn.Close()
	}()

	dec := protoio.NewUint32DelimitedReader(conn, binary.BigEndian, 1e6)
	defer func() { _ = dec.Close() }()

	go func() {
		for {
			var msg types.ProxyMessage

			err := dec.ReadMsg(&msg)
			if err != nil {
				if err == io.EOF {
					log.Debug("shutting down logger goroutine due to file EOF")
					p.active = false
					p.done <- true
					return
				} else {
					log.Warn("client: error reading message")
					dec = protoio.NewUint32DelimitedReader(conn, binary.BigEndian, 1e6)
					return
				}
			}

			switch msg.Type {
			case KindPing:
				if err := p.Pong(); err != nil {
					return
				}
			case KindPong:
			case KindMSG:
				if cl != nil {
					if err := cl(msg); err != nil {
						log.Debug("msg handle err")
						p.done <- true
					}
				}
			}

			msg.Reset()
		}

	}()

	<-p.done
	p.active = false

	return nil
}

func (p *Client) Write(msg []byte) error {

	if !p.active {
		return nil
	}

	return p.Send(msg)
}

func (p *Client) Send(data []byte) error {

	if !p.active {
		return nil
	}

	msg := new(types.ProxyMessage)
	msg.Type = KindMSG
	msg.Partial = false
	msg.Source = p.name
	msg.Line = data
	msg.TimeNano = time.Now().Unix()

	if err := p.writer.WriteMsg(msg); err != nil {
		return err
	}

	return nil
}

func (p *Client) Ping() error {

	msg := new(types.ProxyMessage)
	msg.Type = KindPing
	msg.Partial = false
	msg.Source = p.name
	msg.TimeNano = time.Now().Unix()

	if err := p.writer.WriteMsg(msg); err != nil {
		return err
	}

	return nil
}

func (p *Client) Pong() error {

	msg := new(types.ProxyMessage)
	msg.Type = KindPong
	msg.Partial = false
	msg.Source = p.name
	msg.Line = []byte{}
	msg.TimeNano = time.Now().Unix()

	if err := p.writer.WriteMsg(msg); err != nil {
		return err
	}

	return nil
}

func (p *Client) updateDeadline() error {
	idleDeadline := time.Now().Add(DeadlineWrite)
	return p.conn.SetDeadline(idleDeadline)
}

func NewClient(name, addr string, handler Handler) *Client {
	p := new(Client)
	p.name = name
	p.done = make(chan bool)
	go func() {
		for {
			_ = p.Connect(addr, handler)
			<-time.NewTimer(time.Second).C
		}
	}()
	return p
}
