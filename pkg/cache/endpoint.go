//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
	"net"
	"sync"
)

type EndpointCache struct {
	lock    sync.RWMutex
	storage map[string][]net.IP
}

func (ec *EndpointCache) Get(domain string) []net.IP {
	d, ok := ec.storage[domain]
	if !ok || len(d) == 0 {
		return nil
	}
	if len(d) > 1 {
		d = append(d[1:len(d)], d[0:1]...)
	}
	return d
}

func (ec *EndpointCache) Set(domain string, ips []net.IP) error {
	ec.lock.Lock()
	ec.storage[domain] = ips
	ec.lock.Unlock()
	return nil
}

func (ec *EndpointCache) Del(domain string) error {
	ec.lock.Lock()
	if _, ok := ec.storage[domain]; ok {
		delete(ec.storage, domain)
	}
	ec.lock.Unlock()
	return nil
}

func NewEndpointCache() *EndpointCache {
	return &EndpointCache{
		storage: make(map[string][]net.IP),
	}
}
