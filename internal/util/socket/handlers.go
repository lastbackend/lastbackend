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
	"context"
	"sync"
)

type Emitter struct {
	mu             sync.Mutex
	handlers       map[string][]HandleFunc
	defaultHandler HandleFunc
}

type HandleFunc func(ctx context.Context, event string, sock *Socket, pools map[string]*Pool)

func newEmitter() *Emitter {
	e := new(Emitter)
	e.handlers = make(map[string][]HandleFunc, 0)
	e.defaultHandler = func(ctx context.Context, event string, sock *Socket, pools map[string]*Pool) {}
	return e
}

func (e *Emitter) SetDefaultHandler(handler HandleFunc) {
	e.defaultHandler = handler
}

func (e *Emitter) AddHandler(name string, handlers ...HandleFunc) {
	if len(handlers) == 0 {
		return
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.handlers == nil {
		e.handlers = make(map[string][]HandleFunc, 0)
	}

	h := e.handlers[name]

	if h == nil {
		h = make([]HandleFunc, 0)
	}

	e.handlers[name] = append(h, handlers...)
}

func (e *Emitter) Remove(name string) bool {
	if e.handlers == nil {
		return false
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if h := e.handlers[name]; h != nil {
		delete(e.handlers, name)
		return true
	}
	return false
}

func (e *Emitter) Clear() {
	e.handlers = make(map[string][]HandleFunc, 0)
}

func (e Emitter) Call(ctx context.Context, name string, sock *Socket, pools map[string]*Pool) {
	if e.handlers == nil {
		return
	}

	if h := e.handlers[name]; h != nil && len(h) > 0 {
		for i := range h {
			l := h[i]
			if l != nil {
				l(ctx, name, sock, pools)
			}
		}
	} else {
		e.defaultHandler(ctx, name, sock, pools)
	}
}
