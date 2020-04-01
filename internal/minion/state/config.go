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

package state

import (
	"sync"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/tools/log"
)

const logConfigPrefix = "state:config:>"

type ConfigState struct {
	lock    sync.RWMutex
	configs map[string]*models.ConfigManifest
}

func (s *ConfigState) GetConfigs() map[string]*models.ConfigManifest {
	log.V(logLevel).Debugf("%s get pods", logConfigPrefix)
	return s.configs
}

func (s *ConfigState) SetConfigs(configs map[string]*models.ConfigManifest) {
	log.V(logLevel).Debugf("%s set configs: %d", logConfigPrefix, len(configs))
	for h, config := range configs {
		s.configs[h] = config
	}
}

func (s *ConfigState) GetConfig(name string) *models.ConfigManifest {
	log.V(logLevel).Debugf("%s get config: %s", logConfigPrefix, name)
	s.lock.Lock()
	defer s.lock.Unlock()
	cfg, ok := s.configs[name]
	if !ok {
		return nil
	}
	return cfg
}

func (s *ConfigState) AddConfig(name string, config *models.ConfigManifest) {
	log.V(logLevel).Debugf("%s add config: %s", logConfigPrefix, name)
	s.SetConfig(name, config)
}

func (s *ConfigState) SetConfig(name string, config *models.ConfigManifest) {
	log.V(logLevel).Debugf("%s set config: %s", logConfigPrefix, name)
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.configs[name]; ok {
		delete(s.configs, name)
	}

	s.configs[name] = config
}

func (s *ConfigState) DelConfig(name string) {
	log.V(logLevel).Debugf("%s del config: %s", logConfigPrefix, name)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.configs[name]; ok {
		delete(s.configs, name)
	}
}
