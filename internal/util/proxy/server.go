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
	protoio "github.com/gogo/protobuf/io"
	"net"
	"time"
)

type Server struct {
	Addr        string
	IdleTimeout time.Duration
	inShutdown  bool
	conns       map[Conn]bool
	listener    *net.Listener
}

func (srv Server) Listen(handler Handler) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":2963"
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {

		c, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		conn := Conn{
			Conn:        c,
			IdleTimeout: srv.IdleTimeout,
			done:        make(chan bool),
			error:       make(chan string),
			writer:      protoio.NewUint32DelimitedWriter(c, binary.BigEndian),
		}

		go conn.Handle(handler)

		srv.conns[conn] = true
	}

}
