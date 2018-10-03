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

package state

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"sync"
)

const logRoutePrefix = "state:routes:>"

type RouteState struct {
	lock      sync.RWMutex
	hash      string
	routes map[string]*types.RouteManifest
}

func (es *RouteState) GetHash() string {
	return es.hash
}

func (es *RouteState) SetHash(hash string) {
	es.hash = hash
}

func (es *RouteState) GetRoutes() map[string]*types.RouteManifest {
	log.V(logLevel).Debugf("%s get routes", logRoutePrefix)
	return es.routes
}

func (es *RouteState) SetRoutes(routes map[string]*types.RouteManifest) {
	es.lock.Lock()
	defer es.lock.Unlock()

	for key, route := range routes {
		es.routes[key] = route
	}
}

func (es *RouteState) GetRoute(key string) *types.RouteManifest {
	log.V(logLevel).Debugf("%s: get route: %s", logRoutePrefix, key)
	es.lock.Lock()
	defer es.lock.Unlock()

	ep, ok := es.routes[key]
	if !ok {
		return nil
	}

	return ep
}

func (es *RouteState) AddRoute(key string, route *types.RouteManifest) {
	log.V(logLevel).Debugf("%s: add route: %s", logRoutePrefix, key)
	es.lock.Lock()
	defer es.lock.Unlock()
	es.routes[key] = route
}

func (es *RouteState) SetRoute(key string, route *types.RouteManifest) {
	es.lock.Lock()
	defer es.lock.Unlock()
	log.V(logLevel).Debugf("%s: set route: %s", logRoutePrefix, key)
	es.routes[key] = route
}

func (es *RouteState) DelRoute(key string) {
	es.lock.Lock()
	defer es.lock.Unlock()
	log.V(logLevel).Debugf("%s: del route: %s", logRoutePrefix, key)
	delete(es.routes, key)
}
