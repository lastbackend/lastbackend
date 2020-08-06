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

package socket

import (
	"fmt"
	"sync"
)

// Pool contains a connections used by the same id
type Pool struct {
	sync.Mutex

	ID string

	conns map[*Socket]bool

	join  chan *Socket
	leave chan *Socket
	close chan *Pool

	ignore chan []byte

	broadcast chan []byte
}

// Listen broker channels to manage connections and broadcast data
func (p *Pool) Listen() {

	go func() {
		for {
			select {
			case m := <-p.broadcast:

				for c := range p.conns {
					c.Write(m)
				}

			case c := <-p.join:
				p.conns[c] = true

			case c := <-p.leave:

				delete(p.conns, c)

				if len(p.conns) == 0 {

					close(p.broadcast)
					close(p.join)
					close(p.leave)

					p.close <- p
					return
				}

			case m := <-p.ignore:
				fmt.Println("incoming data processing disabled ", string(m))
			}

		}
	}()

}

// Broadcast message to connections
func (p Pool) Broadcast(event, op, entity string, msg []byte) {
	p.broadcast <- []byte(fmt.Sprintf("{\"event\":\"%s\", \"operation\":\"%s\", \"entity\":\"%s\", \"payload\":%s}", event, op, entity, string(msg)))
}

// manage connection and attach it to broker
func (p Pool) Leave(s *Socket) {
	p.leave <- s
}

// Ping connection to stay it online
func (p Pool) Ping() {
	for c := range p.conns {
		c.Ping()
	}
}

// manage connection and attach it to pool
func (p Pool) manage(sock *Socket) {
	p.join <- sock
}
