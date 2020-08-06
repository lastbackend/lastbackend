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
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/tools/logger"
	"sync"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

const logCacheNode = "api:cache:node"

type CacheNodeManifest struct {
	lock      sync.RWMutex
	nodes     map[string]*models.Node
	ingress   map[string]*models.Ingress
	exporter  map[string]*models.Exporter
	discovery map[string]*models.Discovery
	configs   map[string]*models.ConfigManifest
	manifests map[string]*models.NodeManifest
}

func (c *CacheNodeManifest) checkNode(node string) {
	if _, ok := c.manifests[node]; !ok {
		c.manifests[node] = new(models.NodeManifest)
	}
}

func (c *CacheNodeManifest) SetPodManifest(node, pod string, s *models.PodManifest) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:PodManifestSet:> %s, %s, %#v", logCacheNode, node, pod, s)
	c.lock.Lock()
	defer c.lock.Unlock()

	c.checkNode(node)

	if c.manifests[node].Pods == nil {
		sp := c.manifests[node]
		sp.Pods = make(map[string]*models.PodManifest, 0)
	}

	c.manifests[node].Pods[pod] = s
}

func (c *CacheNodeManifest) DelPodManifest(node, pod string) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:PodManifestDel:> %s, %s", logCacheNode, node, pod)
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.manifests[node]; !ok {
		return
	}

	delete(c.manifests[node].Pods, pod)
}

func (c *CacheNodeManifest) SetVolumeManifest(node, volume string, s *models.VolumeManifest) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:SetVolumeManifest:> %s, %s", logCacheNode, node, volume)

	c.lock.Lock()
	defer c.lock.Unlock()

	c.checkNode(node)

	if c.manifests[node].Volumes == nil {
		sp := c.manifests[node]
		sp.Volumes = make(map[string]*models.VolumeManifest, 0)
	}

	c.manifests[node].Volumes[volume] = s
}

func (c *CacheNodeManifest) DelVolumeManifest(node, volume string) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:DelVolumeManifest:> %s, %s", logCacheNode, node, volume)

	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.manifests[node]; !ok {
		return
	}

	delete(c.manifests[node].Volumes, volume)
}

func (c *CacheNodeManifest) SetSubnetManifest(cidr string, s *models.SubnetManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for n := range c.manifests {

		if _, ok := c.manifests[n].Network[cidr]; !ok {
			c.manifests[n].Network = make(map[string]*models.SubnetManifest)
		}

		c.manifests[n].Network[cidr] = s
	}
}

func (c *CacheNodeManifest) SetSecretManifest(name string, s *models.SecretManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for n := range c.manifests {

		if _, ok := c.manifests[n].Secrets[name]; !ok {
			c.manifests[n].Secrets = make(map[string]*models.SecretManifest)
		}

		c.manifests[n].Secrets[name] = s
	}
}

func (c *CacheNodeManifest) SetConfigManifest(name string, s *models.ConfigManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.configs[name] = s
	for n := range c.manifests {
		if _, ok := c.manifests[n].Configs[name]; !ok {
			c.manifests[n].Configs = make(map[string]*models.ConfigManifest)
		}

		c.manifests[n].Configs[name] = s
	}
}

func (c *CacheNodeManifest) SetEndpointManifest(addr string, s *models.EndpointManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Debugf("%s set endpoint manifest: %s > %s", logCacheNode, addr, s.IP)

	for _, n := range c.manifests {
		if n.Endpoints == nil {
			n.Endpoints = make(map[string]*models.EndpointManifest, 0)
		}
		n.Endpoints[addr] = s
	}
}

func (c *CacheNodeManifest) SetIngress(ingress *models.Ingress) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.ingress[ingress.SelfLink().String()] = ingress
}

func (c *CacheNodeManifest) DelIngress(selflink string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.ingress, selflink)
}

func (c *CacheNodeManifest) SetDiscovery(discovery *models.Discovery) {
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

func (c *CacheNodeManifest) DelDiscovery(selflink string) {
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

func (c *CacheNodeManifest) SetExporter(exporter *models.Exporter) {
	c.lock.Lock()
	defer c.lock.Unlock()

	dvc, ok := c.exporter[exporter.SelfLink().String()]

	if !ok {
		c.exporter[exporter.SelfLink().String()] = exporter
		c.SetExporterEndpoint()
		return
	}

	var update = false
	switch true {
	case dvc.Status.Listener.IP != exporter.Status.Listener.IP:
		update = true
		break
	case dvc.Status.Listener.Port != exporter.Status.Listener.Port:
		update = true
		break
	case dvc.Status.Ready != exporter.Status.Ready:
		update = true
		break
	}
	if update {
		c.exporter[exporter.SelfLink().String()] = exporter
		c.SetExporterEndpoint()
	}
	return
}

func (c *CacheNodeManifest) DelExporter(selflink string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.exporter, selflink)

	for _, n := range c.manifests {
		n.Exporter = nil
	}

	for _, d := range c.exporter {
		if d.Status.Ready {

			exporter := &models.ExporterManifest{
				Endpoint: fmt.Sprintf("%s:%d", d.Status.Listener.IP, d.Status.Listener.Port),
			}

			for _, n := range c.manifests {
				n.Exporter = exporter
			}

			break
		}
	}
}

func (c *CacheNodeManifest) SetExporterEndpoint() {

	for _, n := range c.manifests {
		n.Exporter = nil
	}

	for _, d := range c.exporter {
		if d.Status.Ready {

			exporter := &models.ExporterManifest{
				Endpoint: fmt.Sprintf("%s:%d", d.Status.Listener.IP, d.Status.Listener.Port),
			}

			for _, n := range c.manifests {
				n.Exporter = exporter
			}

			break
		}
	}
}

func (c *CacheNodeManifest) GetExporterEndpoint() *models.ExporterManifest {

	c.lock.Lock()
	defer c.lock.Unlock()

	exporter := new(models.ExporterManifest)

	for _, d := range c.exporter {
		if d.Status.Ready {
			exporter.Endpoint = fmt.Sprintf("%s:%d", d.Status.Listener.IP, d.Status.Listener.Port)
		}
	}

	return exporter
}

func (c *CacheNodeManifest) SetResolvers() {
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

func (c *CacheNodeManifest) GetResolvers() map[string]*models.ResolverManifest {

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

func (c *CacheNodeManifest) GetConfigs() map[string]*models.ConfigManifest {
	return c.configs
}

func (c *CacheNodeManifest) SetNode(node *models.Node) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.nodes[node.SelfLink().String()] = node
}

func (c *CacheNodeManifest) DelNode(node *models.Node) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.nodes, node.SelfLink().String())
	delete(c.manifests, node.SelfLink().String())
}

func (c *CacheNodeManifest) Get(node string) *models.NodeManifest {
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
	c.manifests[node] = new(models.NodeManifest)
}

func (c *CacheNodeManifest) Clear(node string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.manifests, node)
}

func NewCacheNodeManifest() *CacheNodeManifest {
	c := new(CacheNodeManifest)
	c.exporter = make(map[string]*models.Exporter, 0)
	c.manifests = make(map[string]*models.NodeManifest, 0)
	c.ingress = make(map[string]*models.Ingress, 0)
	c.discovery = make(map[string]*models.Discovery, 0)
	c.configs = make(map[string]*models.ConfigManifest, 0)
	return c
}
