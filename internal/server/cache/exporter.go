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

const logCacheExporter = "api:cache:exporter"

type CacheExporterManifest struct {
	lock      sync.RWMutex
	manifests map[string]*models.ExporterManifest
}

func (c *CacheExporterManifest) SetSubnetManifest(cidr string, s *models.SubnetManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()
}

func (c *CacheExporterManifest) Get(exporter string) *models.ExporterManifest {
	c.lock.Lock()
	defer c.lock.Unlock()
	if s, ok := c.manifests[exporter]; !ok {
		return nil
	} else {
		return s
	}
}

func (c *CacheExporterManifest) Flush(exporter string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.manifests[exporter] = new(models.ExporterManifest)
}

func (c *CacheExporterManifest) Clear(exporter string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.manifests, exporter)
}

func NewCacheExporterManifest() *CacheExporterManifest {
	c := new(CacheExporterManifest)
	c.manifests = make(map[string]*models.ExporterManifest, 0)
	return c
}
