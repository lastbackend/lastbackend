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

package cache

import (
	"sync"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

const logCacheIngress = "api:cache:ingress"

type CacheIngressManifest struct {
	lock      sync.RWMutex
	ingress   map[string]*models.Ingress
	discovery map[string]*models.Discovery
	routes    map[string]*models.RouteManifest
	manifests map[string]*models.IngressManifest
}

func (c *CacheIngressManifest) SetSubnetManifest(cidr string, s *models.SubnetManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for n := range c.manifests {

		if _, ok := c.manifests[n].Network[cidr]; !ok {
			c.manifests[n].Network = make(map[string]*models.SubnetManifest)
		}

		c.manifests[n].Network[cidr] = s
	}
}

func (c *CacheIngressManifest) SetRouteManifest(ingress, name string, s *models.RouteManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if s.State == models.StateDestroyed {
		delete(c.routes, name)
	} else {
		c.routes[name] = s
	}

	if _, ok := c.manifests[ingress]; ok {
		if _, ok := c.manifests[ingress].Routes[name]; !ok {
			c.manifests[ingress].Routes = make(map[string]*models.RouteManifest, 0)
		}
		c.manifests[ingress].Routes[name] = s
	}

}

func (c *CacheIngressManifest) DelRouteManifest(ingress, name string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.routes, name)
	if _, ok := c.manifests[ingress]; ok {
		delete(c.manifests[ingress].Routes, name)
	}
}

func (c *CacheIngressManifest) SetEndpointManifest(addr string, s *models.EndpointManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()


	for _, n := range c.manifests {
		if n.Endpoints == nil {
			n.Endpoints = make(map[string]*models.EndpointManifest, 0)
		}
		n.Endpoints[addr] = s
	}
}

func (c *CacheIngressManifest) SetIngress(ingress *models.Ingress) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.ingress[ingress.SelfLink().String()] = ingress
}

func (c *CacheIngressManifest) DelIngress(selflink string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.ingress, selflink)
	delete(c.manifests, selflink)
}

func (c *CacheIngressManifest) SetDiscovery(discovery *models.Discovery) {
	c.lock.Lock()
	defer c.lock.Unlock()

	dvc, ok := c.discovery[discovery.SelfLink().String()]

	if !ok {
		c.discovery[discovery.SelfLink().String()] = discovery
		c.SetResolvers()
		return
	}

	var update = false
	switch true {
	case dvc.Status.IP != discovery.Status.IP:
		update = true
		break
	case dvc.Status.Port != discovery.Status.Port:
		update = true
		break
	case dvc.Status.Ready != discovery.Status.Ready:
		update = true
		break
	}
	if update {
		c.discovery[discovery.SelfLink().String()] = discovery
		c.SetResolvers()
	}
	return
}

func (c *CacheIngressManifest) DelDiscovery(selflink string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.discovery, selflink)

	resolvers := make(map[string]*models.ResolverManifest, 0)

	for _, d := range c.discovery {
		if d.Status.Ready {
			resolvers[d.Status.IP] = &models.ResolverManifest{
				IP:   d.Status.IP,
				Port: d.Status.Port,
			}
		}
	}

	for _, n := range c.manifests {
		n.Resolvers = resolvers
	}
}

func (c *CacheIngressManifest) SetResolvers() {
	resolvers := make(map[string]*models.ResolverManifest, 0)

	for _, d := range c.discovery {
		if d.Status.Ready {
			resolvers[d.Status.IP] = &models.ResolverManifest{
				IP:   d.Status.IP,
				Port: d.Status.Port,
			}
		}
	}

	for _, n := range c.manifests {
		n.Resolvers = resolvers
	}
}

func (c *CacheIngressManifest) GetResolvers() map[string]*models.ResolverManifest {

	resolvers := make(map[string]*models.ResolverManifest, 0)

	for _, d := range c.discovery {
		if d.Status.Ready {
			resolvers[d.Status.IP] = &models.ResolverManifest{
				IP:   d.Status.IP,
				Port: d.Status.Port,
			}
		}
	}

	return resolvers
}

func (c *CacheIngressManifest) Get(ingress string) *models.IngressManifest {
	c.lock.Lock()
	defer c.lock.Unlock()
	if s, ok := c.manifests[ingress]; !ok {
		return nil
	} else {
		return s
	}
}

func (c *CacheIngressManifest) GetRoutes(ingress string) map[string]*models.RouteManifest {
	return c.routes
}

func (c *CacheIngressManifest) Flush(ingress string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.manifests[ingress] = new(models.IngressManifest)
}

func (c *CacheIngressManifest) Clear(ingress string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.manifests, ingress)
}

func NewCacheIngressManifest() *CacheIngressManifest {
	c := new(CacheIngressManifest)
	c.manifests = make(map[string]*models.IngressManifest, 0)
	c.discovery = make(map[string]*models.Discovery, 0)
	c.routes = make(map[string]*models.RouteManifest, 0)
	c.ingress = make(map[string]*models.Ingress)
	return c
}
