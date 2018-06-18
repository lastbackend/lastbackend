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
)

const logCacheNode = "api:cache:node"

type CacheNodeManifest struct {
	lock      sync.RWMutex
	manifests map[string]*types.NodeManifest
}

type NetworkManifestWatcher func(ctx context.Context, event chan *types.NetworkManifestEvent) error

type PodManifestWatcher func(ctx context.Context, event chan *types.PodManifestEvent) error

type VolumeManifestWatcher func(ctx context.Context, event chan *types.VolumeManifestEvent) error

type EndpointManifestWatcher func(ctx context.Context, event chan *types.EndpointManifestEvent) error

func (c *CacheNodeManifest) checkNode(node string) {
	if _, ok := c.manifests[node]; !ok {
		c.manifests[node] = new(types.NodeManifest)
	}
}

func (c *CacheNodeManifest) SetPodManifest(node, pod string, s types.PodManifest) {
	log.Infof("%s:SetPodManifest:> %s, %s, %#v", logCacheNode, node, pod, s)
	c.lock.Lock()
	defer c.lock.Unlock()

	c.checkNode(node)

	if c.manifests[node].Pods == nil {
		sp := c.manifests[node]
		sp.Pods = make(map[string]types.PodManifest, 0)
	}

	c.manifests[node].Pods[pod] = s
}

func (c *CacheNodeManifest) DelPodManifest(node, pod string) {
	log.Infof("%s:DelPodManifest:> %s, %s", node, pod)
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.manifests[node]; !ok {
		return
	}

	delete(c.manifests[node].Pods, pod)
}

func (c *CacheNodeManifest) SetVolumeManifest(node, volume string, s types.VolumeManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.checkNode(node)

	if c.manifests[node].Volumes == nil {
		sp := c.manifests[node]
		sp.Volumes = make(map[string]types.VolumeManifest, 0)
	}

	c.manifests[node].Volumes[volume] = s
}

func (c *CacheNodeManifest) DelVolumeManifest(node, volume string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.manifests[node]; !ok {
		return
	}

	delete(c.manifests[node].Volumes, volume)
}

func (c *CacheNodeManifest) SetNetworkManifest(cidr string, s types.NetworkManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for n := range c.manifests {
		c.manifests[n].Network[cidr] = s
	}
}

func (c *CacheNodeManifest) SetEndpointManifest(addr string, s types.EndpointManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, n := range c.manifests {
		if n.Endpoints == nil {
			n.Endpoints = make(map[string]types.EndpointManifest, 0)
		}
		n.Endpoints[addr] = s
	}
}

func (c *CacheNodeManifest) Get(node string) *types.NodeManifest {
	c.lock.Lock()
	defer c.lock.Unlock()
	if s, ok := c.manifests[node]; !ok {
		return nil
	} else {
		return s
	}
}

func (c *CacheNodeManifest) Flush(node string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.manifests[node] = new(types.NodeManifest)
}

func (c *CacheNodeManifest) Clear(node string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.manifests, node)
}

func (c *CacheNodeManifest) Status(isw IngressStatusWatcher) error {
	evs := make(chan *types.NodeStatusEvent)
	go func() {
		for {
			select {
			case e := <-evs:
				if !e.Ready {
					delete(c.manifests, e.Node)
				}
			}
		}
	}()

	return isw(context.Background(), evs)
}

func (c *CacheNodeManifest) CachePods(ps PodManifestWatcher) error {
	evs := make(chan *types.PodManifestEvent)
	go func() {
		for {
			select {
			case e := <-evs:
				{
					if e.Event == "create" || e.Event == "update" {
						c.SetPodManifest(e.Node, e.Name, e.Manifest)
						continue
					}

					if e.Event == "delete" {
						c.DelPodManifest(e.Node, e.Name)
						continue
					}
				}
			}
		}
	}()

	return ps(context.Background(), evs)
}

func (c *CacheNodeManifest) CacheVolumes(vs VolumeManifestWatcher) error {
	evs := make(chan *types.VolumeManifestEvent)
	go func() {
		for {
			select {
			case e := <-evs:
				{
					if e.Event == "create" || e.Event == "update" {
						c.SetVolumeManifest(e.Node, e.Name, e.Manifest)
						continue
					}

					if e.Event == "delete" {
						c.DelVolumeManifest(e.Node, e.Name)
						continue
					}
				}
			}
		}
	}()

	return vs(context.Background(), evs)
}

func (c *CacheNodeManifest) CacheNetwork(ns NetworkManifestWatcher) error {
	evs := make(chan *types.NetworkManifestEvent)
	go func() {
		for {
			select {
			case e := <-evs:
				{
					if e.Event == "create" || e.Event == "update" {
						c.SetNetworkManifest(e.Manifest.CIDR, e.Manifest)
						continue
					}

					if e.Event == "delete" {
						e.Manifest.State = types.StateDestroy
						c.SetNetworkManifest(e.Manifest.CIDR, e.Manifest)
						continue
					}
				}
			}
		}
	}()

	return ns(context.Background(), evs)
}

func (c *CacheNodeManifest) CacheEndpoints(es EndpointManifestWatcher) error {
	evs := make(chan *types.EndpointManifestEvent)
	go func() {
		for {
			select {
			case e := <-evs:
				{
					if e.Event == "create" || e.Event == "update" {
						c.SetEndpointManifest(e.Manifest.IP, e.Manifest)
						continue
					}

					if e.Event == "delete" {
						e.Manifest.State = types.StateDestroy
						c.SetEndpointManifest(e.Manifest.IP, e.Manifest)
						continue
					}
				}
			}
		}
	}()

	return es(context.Background(), evs)
}

func NewCacheNodeManifest() *CacheNodeManifest {
	c := new(CacheNodeManifest)
	c.manifests = make(map[string]*types.NodeManifest, 0)
	return c
}
