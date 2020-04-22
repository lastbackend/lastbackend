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

package state

import (
	"sync"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/tools/log"
)

const logRoutePrefix = "state:routes:>"

type RouteState struct {
	lock   sync.RWMutex
	hash   string
	routes map[string]struct {
		status   *models.RouteStatus
		manifest *models.RouteManifest
	}
	watchers map[chan string]bool
}

func (rs *RouteState) dispatch(route string) {
	for w := range rs.watchers {
		w <- route
	}
}

func (rs *RouteState) Watch(watcher chan string, done chan bool) {
	rs.watchers[watcher] = true
	defer delete(rs.watchers, watcher)
	<-done
}

func (rs *RouteState) GetHash() string {
	return rs.hash
}

func (rs *RouteState) SetHash(hash string) {
	rs.hash = hash
}

func (rs *RouteState) GetRouteManifests() map[string]*models.RouteManifest {
	log.Debugf("%s get route manifests", logRoutePrefix)

	var manifests = make(map[string]*models.RouteManifest, 0)
	for k, route := range rs.routes {
		if route.manifest != nil {
			manifests[k] = route.manifest
		}
	}

	return manifests
}

func (rs *RouteState) GetRouteManifest(key string) *models.RouteManifest {
	log.Debugf("%s: get route manifest: %s", logRoutePrefix, key)
	rs.lock.Lock()
	defer rs.lock.Unlock()

	ep, ok := rs.routes[key]
	if !ok {
		return nil
	}

	return ep.manifest
}

func (rs *RouteState) AddRouteManifest(key string, route *models.RouteManifest) {
	log.Debugf("%s: add route manifest: %s", logRoutePrefix, key)
	rs.lock.Lock()
	rt, ok := rs.routes[key]
	if !ok {
		rs.routes[key] = struct {
			status   *models.RouteStatus
			manifest *models.RouteManifest
		}{status: nil, manifest: route}
	} else {
		rt.manifest = route
		rs.routes[key] = rt
	}

	rs.lock.Unlock()
}

func (rs *RouteState) SetRouteManifest(key string, route *models.RouteManifest) {
	rs.lock.Lock()
	log.Debugf("%s: set route manifest: %s", logRoutePrefix, key)
	rt, ok := rs.routes[key]
	if !ok {
		rs.routes[key] = struct {
			status   *models.RouteStatus
			manifest *models.RouteManifest
		}{status: nil, manifest: route}
	} else {
		rt.manifest = route
		rs.routes[key] = rt
	}

	rs.lock.Unlock()
}

func (rs *RouteState) DelRouteManifests(key string) {
	rs.lock.Lock()
	log.Debugf("%s: del route manifest: %s", logRoutePrefix, key)
	rt, ok := rs.routes[key]
	if ok {
		rt.manifest = nil
		rs.routes[key] = rt
	}
	rs.lock.Unlock()
}

func (rs *RouteState) GetRouteStatuses() map[string]*models.RouteStatus {
	log.Debugf("%s get route statuses", logRoutePrefix)

	var statuses = make(map[string]*models.RouteStatus, 0)
	for k, route := range rs.routes {
		statuses[k] = route.status
	}

	return statuses
}

func (rs *RouteState) GetRouteStatus(key string) *models.RouteStatus {
	log.Debugf("%s: get route status: %s", logRoutePrefix, key)
	rs.lock.Lock()
	defer rs.lock.Unlock()

	ep, ok := rs.routes[key]
	if !ok {
		return nil
	}

	return ep.status
}

func (rs *RouteState) AddRouteStatus(key string, status *models.RouteStatus) {
	log.Debugf("%s: add route status: %s", logRoutePrefix, key)
	rs.lock.Lock()
	rt, ok := rs.routes[key]
	if !ok {
		rs.routes[key] = struct {
			status   *models.RouteStatus
			manifest *models.RouteManifest
		}{status: status, manifest: nil}
	} else {
		rt.status = status
	}
	rs.routes[key] = rt
	rs.lock.Unlock()
	rs.dispatch(key)
}

func (rs *RouteState) SetRouteStatus(key string, status *models.RouteStatus) {
	rs.lock.Lock()
	log.Debugf("%s: set route status: %s", logRoutePrefix, key)
	rt, ok := rs.routes[key]
	if !ok {
		rs.routes[key] = struct {
			status   *models.RouteStatus
			manifest *models.RouteManifest
		}{status: status, manifest: nil}
	} else {
		rt.status = status
	}
	rs.routes[key] = rt
	rs.lock.Unlock()
	rs.dispatch(key)
}

func (rs *RouteState) DelRoute(key string) {
	rs.lock.Lock()
	log.Debugf("%s: del route: %s", logRoutePrefix, key)
	delete(rs.routes, key)
	rs.lock.Unlock()
	rs.dispatch(key)
}
