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
//
//import (
//	"sync"
//	"time"
//)
//
//type EndpointCache struct {
//	mutex sync.RWMutex
//	ttl   time.Duration
//	items map[string]*Item
//}
//
//func (ec *EndpointCache) Set(key string, data []string) {
//	ec.mutex.Lock()
//	defer ec.mutex.Unlock()
//
//	item := &Item{data: data}
//	item.setExpireTime(ec.ttl)
//	ec.items[key] = item
//}
//
//func (ec *EndpointCache) Get(key string) (data []string) {
//	ec.mutex.Lock()
//	defer ec.mutex.Unlock()
//
//	item, ok := ec.items[key]
//	if !ok || item.expired() {
//		data = nil
//		return
//	}
//
//	item.setExpireTime(ec.ttl)
//	data = item.data
//
//	return
//}
//
//func (ec *EndpointCache) Del(key string) {
//	ec.mutex.Lock()
//	defer ec.mutex.Unlock()
//	delete(ec.items, key)
//	return
//}
//
//func (ec *EndpointCache) Clear() {
//	ec.items = make(map[string]*Item, 0)
//}
//
//func (ec *EndpointCache) Count() int {
//	ec.mutex.RLock()
//	defer ec.mutex.RUnlock()
//	return len(ec.items)
//}
//
//func (ec *EndpointCache) cleanup() {
//	ec.mutex.Lock()
//	defer ec.mutex.Unlock()
//
//	for key, item := range ec.items {
//		if item.expired() {
//			delete(ec.items, key)
//		}
//	}
//}
//
//func (ec *EndpointCache) cleanerTimer() {
//	duration := ec.ttl
//	if duration < time.Second {
//		duration = time.Second
//	}
//	ticker := time.Tick(duration)
//	go (func() {
//		for {
//			select {
//			case <-ticker:
//				ec.cleanup()
//			}
//		}
//	})()
//}
//
//func NewEndpointCache(duration time.Duration) *EndpointCache {
//	cache := &EndpointCache{
//		ttl:   duration,
//		items: make(map[string]*Item, 0),
//	}
//	cache.cleanerTimer()
//	return cache
//}
