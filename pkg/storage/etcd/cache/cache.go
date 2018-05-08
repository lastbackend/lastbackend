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
	"sync"
	"time"
)

type Cache struct {
	mutex sync.RWMutex
	ttl   time.Duration
	items map[string]*Item
}

func (cache *Cache) Set(key string, data interface{}) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	item := &Item{data: data}
	item.setExpireTime(cache.ttl)
	cache.items[key] = item
}

func (cache *Cache) Get(key string) (data interface{}) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	item, ok := cache.items[key]
	if !ok || item.expired() {
		data = nil
		return
	}

	item.setExpireTime(cache.ttl)
	data = item.data

	return
}

func (cache *Cache) Clear() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.items = make(map[string]*Item, 0)
}

func (cache *Cache) Count() int {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	return len(cache.items)
}

func (cache *Cache) cleanup() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	for key, item := range cache.items {
		if item.expired() {
			delete(cache.items, key)
		}
	}
}

func (cache *Cache) cleanerTimer() {
	duration := cache.ttl
	if duration < time.Second {
		duration = time.Second
	}
	ticker := time.Tick(duration)
	go (func() {
		for {
			select {
			case <-ticker:
				cache.cleanup()
			}
		}
	})()
}

func NewCache(duration time.Duration) *Cache {
	cache := &Cache{
		ttl:   duration,
		items: make(map[string]*Item, 0),
	}
	cache.cleanerTimer()
	return cache
}
