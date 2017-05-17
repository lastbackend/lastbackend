//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2016] Last.Backend LLC
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

type Cache struct {
	sync.Mutex
	storage map[string][]net.IP
}

func New() *Cache {
	var c = new(Cache)
	c.storage = make(map[string][]net.IP)
	return c
}

func (c *Cache) Insert(domain string, ips []net.IP) error {
	c.storage[domain] = ips
	return nil
}

func (c *Cache) Remove(domain string) error {
	if _, ok := c.storage[domain]; ok {
		delete(c.storage, domain)
	}
	return nil
}

func (c *Cache) IPList(domain string) []net.IP {
	d, ok := c.storage[domain]
	if !ok || len(d) == 0 {
		return nil
	}
	if len(d) > 1 {
		d = append(d[1:len(d)], d[0:1]...)
	}
	return d
}
