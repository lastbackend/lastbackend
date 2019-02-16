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
	"time"
)

const (
	KindPing string = "ping"
	KindPong string = "pong"
	KindMSG  string = "msg"

	DeadlineWrite = 10 * time.Second
	DeadlineRead  = 5 * time.Second

	DefaultServer = ":2963"
)

func NewServer(addr string) (*Server, error) {
	s := new(Server)
	s.Addr = addr
	s.IdleTimeout, _ = time.ParseDuration("30s")
	s.conns = make(map[Conn]bool, 0)
	return s, nil
}
