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

type ServiceCache struct {
	item    *service
	IsReady bool
	ch      chan ServiceEvent
	ready   chan bool
}

type service struct {
	event string
	data  types.Service
}

type ServiceEvent struct {
	Event string
	Data  types.Service
}

func NewServiceCache() *ServiceCache {
	cache := new(ServiceCache)
	cache.ch = make(chan ServiceEvent)
	cache.ready = make(chan bool)
	return cache
}

func (cache ServiceCache) Set(id string, obj *types.Service) {

	item := service{data: *obj}

	if cache.item == nil {
		item.event = EventCreate
	} else {
		item.event = EventUpdate
	}

	cache.item = &item

	cache.ch <- ServiceEvent{
		Event: item.event,
		Data:  item.data,
	}

	return
}

func (cache ServiceCache) Get(id string) *types.Service {

	if cache.item == nil {
		return nil
	}

	obj := cache.item.data

	return &obj
}

func (cache ServiceCache) Remove(id string) {

	cache.ch <- ServiceEvent{
		Event: EventRemove,
		Data:  cache.item.data,
	}

	cache.item = nil

}

func (cache ServiceCache) Subscribe(ch chan ServiceEvent) {
	for event := range cache.ch {
		cache.ch <- event
	}
}

func (cache ServiceCache) SetReady() {
	cache.ready <- true
}

func (cache ServiceCache) Ready() <-chan bool {
	cache.IsReady = true
	return cache.ready
}
