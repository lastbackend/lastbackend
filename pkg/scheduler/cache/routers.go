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
)

// RouterCache is used for storing routers states in to scheduler to quick access
type RouterCache struct {
	lock sync.RWMutex
	routers map[string]*types.NodeMeta
}

// Add node info to cache
func (nc *RouterCache) Add(node *types.Node) {
	nc.lock.Lock()
	defer nc.lock.Unlock()

	if node.Roles.Router.Enabled {
		nc.routers[node.SelfLink()] = &node.Meta
	}
}

// Del node info from cache
func (nc *RouterCache) Del(node *types.Node) {
	nc.lock.Lock()
	defer nc.lock.Unlock()

	if node.Roles.Router.Enabled {
		delete(nc.routers, node.SelfLink())
	}
}

// Set node info in cache
func (nc *RouterCache) Set(node *types.Node) {
	nc.lock.Lock()
	defer nc.lock.Unlock()

	if _, ok := nc.routers[node.SelfLink()]; ok && !node.Roles.Router.Enabled {
		delete(nc.routers, node.SelfLink())
	}

	if node.Roles.Router.Enabled {
		nc.routers[node.SelfLink()] = &node.Meta
	}
}

// List routers from cache
func (nc *RouterCache) List() map[string]*types.Node {
	return nc.routers
}

// NewNodeCache returns new node cache
func NewRouterCache() *RouterCache {
	nc := new(RouterCache)
	return nc
}
