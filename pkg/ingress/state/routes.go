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

package state

import (
	"github.com/lastbackend/lastbackend/pkg/log"
	"sync"
)

const stringEmpty = ""

type RoutesState struct {
	lock    sync.RWMutex
	configs map[string]string
}

func (s *RoutesState) GetRoutes() map[string]string {
	log.V(logLevel).Debug("Cache: List: get routes")
	return s.configs
}

// Set config hash by file name
func (s *RoutesState) Set(name, hash string) {
	log.V(logLevel).Debugf("Cache: RouterCache: add cancel func config: %s", name)
	s.configs[name] = hash
}

// Get config hash by file name
func (s *RoutesState) Get(name string) string {
	log.V(logLevel).Debugf("Cache: RouterCache: get cancel func config: %s", name)

	if _, ok := s.configs[name]; ok {
		return s.configs[name]
	}
	return stringEmpty
}

// Remove config hash by file name
func (s *RoutesState) Del(name string) {
	log.V(logLevel).Debugf("Cache: RouterCache: del cancel func config: %s", name)
	delete(s.configs, name)
}
