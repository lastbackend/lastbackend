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
	"github.com/lastbackend/lastbackend/pkg/logger"
	"net"
	"sync"
)

type EndpointCache struct {
	log     logger.ILogger
	lock    sync.RWMutex
	storage map[string][]net.IP
}

func (ec *EndpointCache) Get(domain string) []net.IP {
	ec.log.V(logLevel).Debugf("Cache: EndpointCache: get ips for domain: %s", domain)

	d, ok := ec.storage[domain]
	if !ok || len(d) == 0 {
		return nil
	}
	return d
}

func (ec *EndpointCache) Set(domain string, ips []net.IP) error {
	ec.log.V(logLevel).Debugf("Cache: EndpointCache: set ips for domain: %s", domain)

	ec.lock.Lock()
	ec.storage[domain] = ips
	ec.lock.Unlock()
	return nil
}

func (ec *EndpointCache) Del(domain string) error {
	ec.log.V(logLevel).Debugf("Cache: EndpointCache: del domain: %s", domain)

	ec.lock.Lock()
	if _, ok := ec.storage[domain]; ok {
		delete(ec.storage, domain)
	}
	ec.lock.Unlock()
	return nil
}

func NewEndpointCache(log logger.ILogger) *EndpointCache {
	log.V(logLevel).Debug("Cache: EndpointCache: initialization storage")
	return &EndpointCache{
		log:     log,
		storage: make(map[string][]net.IP),
	}
}
