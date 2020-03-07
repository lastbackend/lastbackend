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
	"github.com/lastbackend/lastbackend/tools/log"
	"sync"

	"github.com/lastbackend/lastbackend/internal/pkg/types"
)

const logCacheExporter = "api:cache:exporter"

type CacheExporterManifest struct {
	lock      sync.RWMutex
	manifests map[string]*types.ExporterManifest
}

func (c *CacheExporterManifest) SetSubnetManifest(cidr string, s *types.SubnetManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()
}

func (c *CacheExporterManifest) Get(exporter string) *types.ExporterManifest {
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
	c.manifests[exporter] = new(types.ExporterManifest)
}

func (c *CacheExporterManifest) Clear(exporter string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.Debugf("clear cache for exporter: %s", exporter)
	delete(c.manifests, exporter)
}

func NewCacheExporterManifest() *CacheExporterManifest {
	c := new(CacheExporterManifest)
	c.manifests = make(map[string]*types.ExporterManifest, 0)
	return c
}
