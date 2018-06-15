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
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"sync"
	"strings"
)

type CacheNodeSpec struct {
	lock sync.RWMutex
	spec map[string]*types.NodeSpec
}

type NetworkSpecWatcher func(ctx context.Context, event chan *types.Event) error

type PodSpecWatcher func(ctx context.Context, event chan *types.Event) error

type VolumeSpecWatcher func(ctx context.Context, event chan *types.Event) error

type NodeStatusWatcher func(ctx context.Context, event chan *types.Event) error

type EndpointSpecWatcher func(ctx context.Context, event chan *types.Event) error

func (c *CacheNodeSpec) SetPodSpec(node, pod string, s types.PodSpec) {
	log.Info("api:cache:setpodspec:> %s, %s, %#v", node, pod, s)
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.spec[node]; !ok {
		c.spec[node] = new(types.NodeSpec)
	}

	if c.spec[node].Pods == nil {
		sp := c.spec[node]
		sp.Pods = make(map[string]types.PodSpec, 0)
	}

	c.spec[node].Pods[pod] = s
}

func (c *CacheNodeSpec) DelPodSpec(node, pod string) {
	log.Info("api:cache:delpodspec:> %s, %s", node, pod)
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.spec[node].Pods, pod)
}

func (c *CacheNodeSpec) SetVolumeSpec(node, volume string, s types.VolumeSpec) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.spec[node]; !ok {
		c.spec[node] = new(types.NodeSpec)
	}

	if c.spec[node].Volumes == nil {
		sp := c.spec[node]
		sp.Volumes = make(map[string]types.VolumeSpec, 0)
	}

	c.spec[node].Volumes[volume] = s
}

func (c *CacheNodeSpec) DelVolumeSpec(node, volume string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.spec[node].Volumes, volume)
}

func (c *CacheNodeSpec) SetNetworkSpec(node string, s types.NetworkSpec) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for n := range c.spec {
		if c.spec[node].Network == nil {
			sp := c.spec[node]
			sp.Network = make(map[string]types.NetworkSpec, 0)
		}

		c.spec[n].Network[node] = s
	}
}

func (c *CacheNodeSpec) DelNetworkSpec(node string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for n := range c.spec {
		delete(c.spec[n].Network, node)
	}

}

func (c *CacheNodeSpec) SetEndpointSpec(endpoint string, s types.EndpointSpec) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, n := range c.spec {
		if n.Endpoints == nil {
			n.Endpoints = make(map[string]types.EndpointSpec, 0)
		}
		n.Endpoints[endpoint] = s
	}
}

func (c *CacheNodeSpec) DelEndpointSpec(endpoint string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, n := range c.spec {
		delete(n.Endpoints, endpoint)
	}
}

func (c *CacheNodeSpec) Get(node string) *types.NodeSpec {
	c.lock.Lock()
	defer c.lock.Unlock()
	if s, ok := c.spec[node]; !ok {
		return nil
	} else {
		return s
	}
}

func (c *CacheNodeSpec) Flush(node string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.spec[node] = new(types.NodeSpec)
}

func (c *CacheNodeSpec) Clear(node string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.spec, node)
}

func (c *CacheNodeSpec) CachePods(ps PodSpecWatcher) error {
	evs := make(chan *types.Event)

	go func() {
		for {
			select {
			case e := <-evs:
				{

					if e.Data == nil {
						continue
					}

					spec := e.Data.(types.PodSpec)
					parse := strings.Split(e.Name, ":")
					node := parse[0]
					pod := parse[1]

					switch e.Action {
					case types.EventActionCreate:
						fallthrough
					case types.EventActionUpdate:
						c.SetPodSpec(node, pod, spec)
					case types.EventActionDelete:
						c.DelPodSpec(node, pod)
					}

				}
			}
		}
	}()

	return ps(context.Background(), evs)
}

func (c *CacheNodeSpec) CacheVolumes(vs VolumeSpecWatcher) error {
	evs := make(chan *types.Event)

	go func() {
		for {
			select {
			case e := <-evs:
				{

					if e.Data == nil {
						continue
					}

					spec := e.Data.(types.VolumeSpec)
					parse := strings.Split(e.Name, ":")
					node := parse[0]
					volume := parse[1]

					switch e.Action {
					case types.EventActionCreate:
						fallthrough
					case types.EventActionUpdate:
						c.SetVolumeSpec(node, volume, spec)
					case types.EventActionDelete:
						c.DelVolumeSpec(node, volume)
					}

				}
			}
		}
	}()

	return vs(context.Background(), evs)
}

func (c *CacheNodeSpec) CacheNetwork(ns NetworkSpecWatcher) error {
	evs := make(chan *types.Event)
	go func() {
		for {
			select {
			case e := <-evs:
				{

					if e.Data == nil {
						continue
					}

					spec := e.Data.(types.NetworkSpec)
					node := e.Name

					switch e.Action {
					case types.EventActionCreate:
						fallthrough
					case types.EventActionUpdate:
						c.SetNetworkSpec(node, spec)
					case types.EventActionDelete:
						c.DelNetworkSpec(node)
					}

				}
			}
		}
	}()

	return ns(context.Background(), evs)
}

func (c *CacheNodeSpec) CacheEndpoints(es EndpointSpecWatcher) error {

	//evs := make(chan *types.EndpointSpecEvent)
	evs := make(chan *types.Event)

	go func() {
		for {
			select {
			case e := <-evs:
				{

					if e.Data == nil {
						continue
					}

					spec := e.Data.(types.EndpointSpec)

					switch e.Action {
					case types.EventActionCreate:
						fallthrough
					case types.EventActionUpdate:
						c.SetEndpointSpec(spec.IP, spec)
					case types.EventActionDelete:
						c.DelEndpointSpec(spec.IP)
					}

				}
			}
		}
	}()

	return es(context.Background(), evs)
}

func (c *CacheNodeSpec) Del(dw NodeStatusWatcher) error {
	evs := make(chan *types.Event)
	go func() {
		for {
			select {
			case e := <-evs:

				if e.Data == nil {
					continue
				}

				online := e.Data.(bool)
				node := e.Name

				if !online {
					delete(c.spec, node)
				}
			}
		}
	}()

	return dw(context.Background(), evs)
}

func NewCacheNodeSpec() *CacheNodeSpec {
	c := new(CacheNodeSpec)
	c.spec = make(map[string]*types.NodeSpec, 0)
	return c
}
