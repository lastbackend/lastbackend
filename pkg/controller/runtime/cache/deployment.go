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
)

type DeploymentCache struct {
	data    map[string]deployment
	IsReady bool
	ch      chan DeploymentEvent
	ready   chan bool
}

type deployment struct {
	event string
	data  types.Deployment
}

type DeploymentEvent struct {
	Event string
	Data  types.Deployment
}

func NewDeploymentCache() *DeploymentCache {
	cache := new(DeploymentCache)
	cache.data = make(map[string]deployment, 0)
	cache.ch = make(chan DeploymentEvent)
	cache.ready = make(chan bool)
	return cache
}

func (cache DeploymentCache) Set(id string, obj *types.Deployment) {

	item := deployment{data: *obj}

	if _, ok := cache.data[id]; ok {
		item.event = EventUpdate
	} else {
		item.event = EventCreate
	}

	cache.data[id] = item

	cache.ch <- DeploymentEvent{
		Event: item.event,
		Data:  item.data,
	}

	return
}

func (cache DeploymentCache) Get(id string) *types.Deployment {

	if _, ok := cache.data[id]; !ok {
		return nil
	}

	obj := cache.data[id].data

	return &obj
}

func (cache DeploymentCache) Remove(id string) error {
	return nil
}

func (cache DeploymentCache) Subscribe(ch chan DeploymentEvent) {
	for event := range cache.ch {
		cache.ch <- event
	}
}

func (cache DeploymentCache) SetReady() {
	cache.ready <- true
}

func (cache DeploymentCache) Ready() <-chan bool {
	cache.IsReady = true
	return cache.ready
}
