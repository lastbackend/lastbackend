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

type PodCache struct {
	data    map[string]pod
	IsReady bool

	ch    chan PodEvent
	ready chan bool
}

type pod struct {
	event string
	data  types.Pod
}

type PodEvent struct {
	Event string
	Data  types.Pod
}

func NewPodCache() *PodCache {
	cache := new(PodCache)
	cache.data = make(map[string]pod, 0)
	cache.ch = make(chan PodEvent)
	cache.ready = make(chan bool)
	return cache
}

func (cache PodCache) Set(id string, obj *types.Pod) {

	item := pod{data: *obj}

	if _, ok := cache.data[id]; ok {
		item.event = EventUpdate
	} else {
		item.event = EventCreate
	}

	cache.data[id] = item

	cache.ch <- PodEvent{
		Event: item.event,
		Data:  item.data,
	}

	return
}

func (cache PodCache) Get(id string) *types.Pod {

	if _, ok := cache.data[id]; !ok {
		return nil
	}

	obj := cache.data[id].data

	return &obj
}

func (cache PodCache) Remove(id string) error {
	return nil
}

func (cache PodCache) Subscribe(ch chan PodEvent) {
	for event := range cache.ch {
		cache.ch <- event
	}
}

func (cache PodCache) SetReady() {
	cache.IsReady = true
	cache.ready <- true
}

func (cache PodCache) Ready() <-chan bool {
	return cache.ready
}
