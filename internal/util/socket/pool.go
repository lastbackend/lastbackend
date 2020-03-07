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

	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logLevel  = 3
	logPrefix = "utils:socket"
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
	log.V(logLevel).Debugf("%s:pool:listen:> listen broker channels to manage connections and broadcast data", logPrefix)

	go func() {
		for {
			select {
			case m := <-p.broadcast:
				log.V(logLevel).Debugf("%s:pool:listen:> broker %s broadcast: %s", logPrefix, p.ID, string(m))

				for c := range p.conns {
					c.Write(m)
				}

			case c := <-p.join:
				log.V(logLevel).Debugf("%s:pool:listen:> join connection to broker: %s", logPrefix, p.ID)
				p.conns[c] = true

			case c := <-p.leave:
				log.V(logLevel).Debugf("%s:pool:listen:> leave connection from broker: %s", logPrefix, p.ID)

				delete(p.conns, c)

				if len(p.conns) == 0 {

					close(p.broadcast)
					close(p.join)
					close(p.leave)

					log.V(logLevel).Debugf("%s:pool:listen:> broker closed successful", logPrefix)

					p.close <- p
					return
				}

			case m := <-p.ignore:
				log.V(logLevel).Debugf("%s:pool:listen:> incoming data processing disabled: %s", logPrefix, string(m))
			}

		}
	}()

}

// Broadcast message to connections
func (p Pool) Broadcast(event, op, entity string, msg []byte) {
	log.V(logLevel).Debugf("%s:pool:broadcast:> broadcast message to connections", logPrefix)
	p.broadcast <- []byte(fmt.Sprintf("{\"event\":\"%s\", \"operation\":\"%s\", \"entity\":\"%s\", \"payload\":%s}", event, op, entity, string(msg)))
}

// manage connection and attach it to broker
func (p Pool) Leave(s *Socket) {
	log.V(logLevel).Debugf("%s:pool:manage:> drop connection from broker", logPrefix)
	p.leave <- s
}

// Ping connection to stay it online
func (p Pool) Ping() {
	log.V(7).Debugf("%s:pool:ping:> ping connection to stay it online", logPrefix)
	for c := range p.conns {
		c.Ping()
	}
}

// manage connection and attach it to pool
func (p Pool) manage(sock *Socket) {
	p.join <- sock
}
