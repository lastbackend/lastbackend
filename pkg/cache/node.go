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
	"github.com/lastbackend/lastbackend/pkg/log"
	"sync"
)

type NodeCache struct {
	lock    sync.RWMutex
	storage map[string]*types.Node
}

func (ec *NodeCache) Get(hostname string) *types.Node {
	log.V(logLevel).Debugf("Cache: NodeCache: get node by hostname: %s", hostname)
	if d, ok := ec.storage[hostname]; !ok {
		return nil
	} else {
		return d
	}
}

func (ec *NodeCache) List() map[string]*types.Node {
	log.V(logLevel).Debugf("Cache: NodeCache: get nodes list")
	return ec.storage
}

func (ec *NodeCache) Set(hostname string, node *types.Node) error {
	log.V(logLevel).Debugf("Cache: NodeCache: set node by hostname: %s", hostname)
	ec.lock.Lock()
	ec.storage[hostname] = node
	ec.lock.Unlock()
	return nil
}

func (ec *NodeCache) Del(hostname string) error {
	log.V(logLevel).Debugf("Cache: NodeCache: del node: %s", hostname)
	ec.lock.Lock()
	if _, ok := ec.storage[hostname]; ok {
		delete(ec.storage, hostname)
	}
	ec.lock.Unlock()
	return nil
}

func NewNodeCache() *NodeCache {
	log.V(logLevel).Debug("Cache: NodeCache: initialization storage")
	return &NodeCache{
		storage: make(map[string]*types.Node),
	}
}
