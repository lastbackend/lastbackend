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

const logVolumePrefix = "state:volume:> "

type VolumesState struct {
	lock     sync.RWMutex
	volumes  map[string]types.VolumeStatus
	local    map[string]bool
	watchers map[chan string]bool
	claims   map[string]types.VolumeClaim
}

func (s *VolumesState) dispatch(pod string) {
	for w := range s.watchers {
		w <- pod
	}
}

func (s *VolumesState) Watch(watcher chan string, done chan bool) {
	s.watchers[watcher] = true
	defer delete(s.watchers, watcher)
	<-done
}

func (s *VolumesState) GetVolumes() map[string]types.VolumeStatus {
	log.V(logLevel).Debugf("%s get volumes", logVolumePrefix)
	return s.volumes
}

func (s *VolumesState) SetVolumes(key string, volumes []*types.VolumeStatus) {
	log.V(logLevel).Debugf("%s set volumes: %#v", logVolumePrefix, volumes)
	for _, vol := range volumes {
		s.volumes[key] = *vol
	}
}

func (s *VolumesState) GetVolume(key string) *types.VolumeStatus {
	log.V(logLevel).Debugf("%s get volume: %s", logVolumePrefix, key)
	s.lock.Lock()
	defer s.lock.Unlock()
	v, ok := s.volumes[key]
	if !ok {
		return nil
	}
	return &v
}

func (s *VolumesState) AddVolume(key string, v *types.VolumeStatus) {
	log.V(logLevel).Debugf("%s add volume: %s > %s", logVolumePrefix, key, v.State)
	s.SetVolume(key, v)
}

func (s *VolumesState) SetVolume(key string, v *types.VolumeStatus) {
	log.V(logLevel).Debugf("%s set volume: %s > %s", logVolumePrefix, key, v.State)
	s.lock.Lock()
	s.volumes[key] = *v
	s.lock.Unlock()
	s.dispatch(key)
}

func (s *VolumesState) DelVolume(key string) {
	log.V(logLevel).Debugf("%s del volume: %#v", logVolumePrefix, key)
	s.lock.Lock()
	if _, ok := s.volumes[key]; ok {
		delete(s.volumes, key)
	}
	s.lock.Unlock()
	s.dispatch(key)
}

func (s *VolumesState) GetClaim(key string) *types.VolumeClaim {
	log.V(logLevel).Debugf("%s get claim: %s", logVolumePrefix, key)
	v, ok := s.claims[key]
	if !ok {
		return nil
	}
	return &v
}

func (s *VolumesState) AddClaim(key string, vc *types.VolumeClaim) {
	log.V(logLevel).Debugf("%s add claim: %s", logVolumePrefix, key)
	s.SetClaim(key, vc)
}

func (s *VolumesState) SetClaim(key string, vc *types.VolumeClaim) {
	log.V(logLevel).Debugf("%s set claim: %s", logVolumePrefix, key)
	s.lock.Lock()
	s.claims[key] = *vc
	s.lock.Unlock()
}

func (s *VolumesState) DelClaim(key string) {
	log.V(logLevel).Debugf("%s del claim: %#v", logVolumePrefix, key)
	s.lock.Lock()
	if _, ok := s.claims[key]; ok {
		delete(s.claims, key)
	}
	s.lock.Unlock()
}

func (s *VolumesState) SetLocal(key string) {
	log.V(logLevel).Debugf("%s set volume: %s as local", logVolumePrefix, key)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.local[key] = true
}

func (s *VolumesState) DelLocal(key string) {
	log.V(logLevel).Debugf("%s del volume: %s from local", logVolumePrefix, key)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.local[key] = true
}

func (s *VolumesState) IsLocal(key string) bool {
	log.V(logLevel).Debugf("%s check volume: %s is local", logVolumePrefix, key)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.local[key]; ok {
		return true
	}

	return false
}
