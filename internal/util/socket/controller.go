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

package socket

import (
	"github.com/lastbackend/lastbackend/tools/log"
	"sync"
)

// Controller - controller of pools of connections
type Controller struct {
	sync.Mutex

	pools map[string]*Pool
	conns map[*Socket]*Pool

	join  chan *Pool
	leave chan *Pool

	Leave   chan *Socket
	Message chan *Message

	Events *Emitter
}

// Listen controller channels to manage pools
func (c *Controller) Listen() {
	log.V(logLevel).Debugf("%s:controller:listen:> listen controller channels to manage pools", logPrefix)

	go func() {
		for {
			select {
			case sock := <-c.Leave:

				if _, ok := c.conns[sock]; !ok {
					continue
				}

				c.Lock()
				c.conns[sock].leave <- sock
				c.Unlock()

				c.Lock()
				delete(c.pools, c.conns[sock].ID)
				c.Unlock()

				c.Lock()
				delete(c.conns, sock)
				c.Unlock()

			case e := <-c.Message:

				if _, ok := c.conns[e.Socket]; ok {
					c.Lock()
					c.conns[e.Socket].leave <- e.Socket
					c.Unlock()
				}

				ev := string(e.Data)
				c.Events.Call(e.Socket.Context(), ev, e.Socket, c.pools)

				if p, ok := c.pools[ev]; ok {
					c.Lock()
					c.conns[e.Socket] = p
					c.Unlock()
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case p := <-c.join:

				c.Lock()
				c.pools[p.ID] = p
				c.Unlock()

			case p := <-c.leave:

				c.Lock()
				delete(c.pools, p.ID)
				if err := c.Clean(p.ID); err != nil {
					log.Errorf("%s:controller:listen:> connection pool clean err: %v", logPrefix, err)
				}
				c.Unlock()
			}
		}
	}()
}

// Broadcast message to all pools
// TODO: need optimization
func (c *Controller) Broadcast(event, op, entity string, data []byte) error {
	for _, p := range c.pools {
		p.Broadcast(event, op, entity, data)
	}
	return nil
}

// Get returns named broker by id
func (c *Controller) Get(id string) *Pool {
	log.V(logLevel).Debugf("%s:controller:get:> get returns named broker by id %s", logPrefix, id)
	return c.pools[id]
}

// Add create and return new connections broker
func (c *Controller) Add(id string, sock *Socket) *Pool {
	log.V(logLevel).Debugf("%s:controller:add:> create new connections broker %s", logPrefix, id)

	var p = new(Pool)
	p.ID = id

	p.conns = make(map[*Socket]bool)

	p.join = make(chan *Socket)
	p.leave = make(chan *Socket)
	p.close = c.leave

	p.ignore = make(chan []byte)
	p.broadcast = make(chan []byte)
	p.Listen()

	p.manage(sock)

	c.Lock()
	c.conns[sock] = p
	c.Unlock()

	c.join <- p

	return p
}

// Attach connection to pool
func (c *Controller) Attach(pool *Pool, sock *Socket) {
	log.V(logLevel).Debugf("%s:controller:attach:> attach connection to pool", logPrefix)
	pool.manage(sock)

	c.Lock()
	c.conns[sock] = pool
	c.Unlock()
}

func (c *Controller) Clean(id string) error {
	log.V(logLevel).Debugf("%s:controller:clean:> remove broker session %s", logPrefix, id)
	return nil
}

// Ping all pools and internal connections
func (c *Controller) Ping() {
	log.V(7).Debugf("%s:ping:> ping all pools and internal connections", logPrefix)
	for _, p := range c.pools {
		p.Ping()
	}
}

func New() *Controller {
	var c = new(Controller)

	c.pools = make(map[string]*Pool, 0)
	c.conns = make(map[*Socket]*Pool, 0)
	c.join = make(chan *Pool)
	c.leave = make(chan *Pool)

	c.Message = make(chan *Message)
	c.Leave = make(chan *Socket)

	c.Events = newEmitter()

	c.Listen()

	return c
}
