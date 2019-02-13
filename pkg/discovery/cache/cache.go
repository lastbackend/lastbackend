//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"time"

	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/spf13/viper"
)

const (
	logLevel          = 7
	defaultExpireTime = 24 // 24 hours
)

type Cache struct {
	endpoints *EndpointCache
}

func New() *Cache {
	log.V(logLevel).Debug("Cache: initialization cache storage")

	var duration = viper.GetDuration("discovery.cache.duration")
	if duration == 0 {
		duration = defaultExpireTime
	}

	return &Cache{
		endpoints: NewEndpointCache(duration * time.Minute),
	}
}

// Return endpoint storage
func (s *Cache) Endpoint() *EndpointCache {
	return s.endpoints
}
