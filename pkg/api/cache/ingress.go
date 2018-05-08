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

package cache

import (
	"sync"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"context"
	"github.com/lastbackend/lastbackend/pkg/log"
)

type CacheIngressSpec struct {
	lock sync.RWMutex
	spec map[string]*types.IngressSpec
}

type IngressStatusWatcher func(ctx context.Context, event chan *types.IngressStatusEvent) error

type RouteSpecWatcher func(ctx context.Context, event chan *types.RouteSpecEvent) error

func (c *CacheIngressSpec) SetRouteSpec(route string, s types.RouteSpec) {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Debugf("add route spec: %s", route)
	for i := range c.spec {
		if _, ok := c.spec[i]; !ok {
			c.spec[i] = new(types.IngressSpec)
		}

		if c.spec[i].Routes == nil {
			sp := c.spec[i]
			sp.Routes = make(map[string]types.RouteSpec, 0)
		}

		c.spec[i].Routes[route] = s
	}
}

func (c *CacheIngressSpec) DelRouteSpec(route string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.Debugf("del route spec: %s", route)
	for i := range c.spec {
		delete(c.spec[i].Routes, route)
	}
}

func (c *CacheIngressSpec) CacheRoutes(rs RouteSpecWatcher) error {
	evs := make(chan *types.RouteSpecEvent)
	go func() {
		for {
			select {
			case e := <-evs:
				{
					if e.Event == "create" || e.Event == "update" {
						c.SetRouteSpec(e.Name, e.Spec)
						continue
					}

					if e.Event == "delete" {
						c.DelRouteSpec(e.Name)
						continue
					}
				}
			}
		}
	}()

	return rs(context.Background(), evs)
}

func (c *CacheIngressSpec) Get(ingress string) *types.IngressSpec {
	c.lock.Lock()
	defer c.lock.Unlock()
	if s, ok := c.spec[ingress]; !ok {
		return nil
	} else {
		return s
	}
}

func (c *CacheIngressSpec) Flush(ingress string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.spec[ingress] = new(types.IngressSpec)
	c.spec[ingress].Routes = make(map[string]types.RouteSpec, 0)
}

func (c *CacheIngressSpec) Clear(ingress string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.spec, ingress)
}

func (c *CacheIngressSpec) Status(isw IngressStatusWatcher) error {
	evs := make(chan *types.IngressStatusEvent)
	go func() {
		for {
			select {
			case e := <-evs:
				if !e.Status.Ready {
					delete(c.spec, e.Name)
				}
			}
		}
	}()

	return isw(context.Background(), evs)
}

func NewCacheIngressSpec() *CacheIngressSpec {
	c := new(CacheIngressSpec)
	c.spec = make(map[string]*types.IngressSpec, 0)
	return c
}
