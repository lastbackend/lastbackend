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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"sync"
)

type ConfigState struct {
	lock    sync.RWMutex
	configs map[string]*types.ConfigManifest
}

func (s *ConfigState) GetConfigs() map[string]*types.ConfigManifest {
	log.V(logLevel).Debug("Cache: ConfigCache: get pods")
	return s.configs
}

func (s *ConfigState) SetConfigs(configs map[string]*types.ConfigManifest) {
	log.V(logLevel).Debugf("Cache: ConfigCache: set configs: %#v", configs)
	for h, config := range configs {
		s.configs[h] = config
	}
}

func (s *ConfigState) GetConfig(name string) *types.ConfigManifest {
	log.V(logLevel).Debugf("Cache: ConfigCache: get config: %s", name)
	s.lock.Lock()
	defer s.lock.Unlock()
	cfg, ok := s.configs[name]
	if !ok {
		return nil
	}
	return cfg
}

func (s *ConfigState) AddConfig(name string, config *types.ConfigManifest) {
	log.V(logLevel).Debugf("Cache: ConfigCache: add config: %#v", config)
	s.SetConfig(name, config)
}

func (s *ConfigState) SetConfig(name string, config *types.ConfigManifest) {
	log.V(logLevel).Debugf("Cache: ConfigCache: set config: %#v", config)
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.configs[name]; ok {
		delete(s.configs, name)
	}

	s.configs[name] = config
}

func (s *ConfigState) DelConfig(name string) {
	log.V(logLevel).Debugf("Cache: ConfigCache: del config: %s", name)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.configs[name]; ok {
		delete(s.configs, name)
	}
}
