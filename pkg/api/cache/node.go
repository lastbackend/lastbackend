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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"sync"
	"context"
	"github.com/lastbackend/lastbackend/pkg/log"
)

type CacheNodeSpec struct {
	lock sync.RWMutex
	spec map[string]*types.NodeSpec
}

type NetworkSpecWatcher func(ctx context.Context, event chan *types.NetworkSpecEvent) error

type PodSpecWatcher func(ctx context.Context, event chan *types.PodSpecEvent) error

type VolumeSpecWatcher func(ctx context.Context, event chan *types.VolumeSpecEvent) error

type NodeStatusWatcher func(ctx context.Context, event chan *types.NodeStatusEvent) error


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
	evs := make(chan *types.PodSpecEvent)
	go func() {
		for {
			select {
			case e := <-evs:
				{
					if e.Event == "create" || e.Event == "update" {
						c.SetPodSpec(e.Node, e.Name, e.Spec)
						continue
					}

					if e.Event == "delete" {
						c.DelPodSpec(e.Node, e.Name)
						continue
					}
				}
			}
		}
	}()

	return ps(context.Background(), evs)
}

func (c *CacheNodeSpec) CacheVolumes(vs VolumeSpecWatcher) error {
	evs := make(chan *types.VolumeSpecEvent)
	go func() {
		for {
			select {
			case e := <-evs:
				{
					if e.Event == "create" || e.Event == "update" {
						c.SetVolumeSpec(e.Node, e.Name, e.Spec)
						continue
					}

					if e.Event == "delete" {
						c.DelVolumeSpec(e.Node, e.Name)
						continue
					}
				}
			}
		}
	}()

	return vs(context.Background(), evs)
}

func (c *CacheNodeSpec) CacheNetwork(ns NetworkSpecWatcher) error {
	evs := make(chan *types.NetworkSpecEvent)
	go func() {
		for {
			select {
			case e := <-evs:
				{
					if e.Event == "create" || e.Event == "update" {
						c.SetNetworkSpec(e.Node, e.Spec)
						continue
					}

					if e.Event == "delete" {
						c.DelNetworkSpec(e.Node)
						continue
					}
				}
			}
		}
	}()

	return ns(context.Background(), evs)
}


func (c *CacheNodeSpec) Del(dw NodeStatusWatcher) error {
	evs := make(chan *types.NodeStatusEvent)
	go func() {
		for {
			select {
			case e := <-evs:
				if !e.Online {
					delete(c.spec, e.Node)
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
