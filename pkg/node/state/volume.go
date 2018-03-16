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
)

func (s *VolumesState) GetVolumes() map[string]types.Volume {
	log.V(logLevel).Debug("Cache: VolumeCache: get pods")
	return s.volumes
}

func (s *VolumesState) SetVolumes(volumes []*types.Volume) {
	log.V(logLevel).Debugf("Cache: VolumeCache: set volumes: %#v", volumes)
	for _, secret := range volumes {
		s.volumes[secret.Meta.Name] = *secret
	}
}

func (s *VolumesState) GetVolume(hash string) *types.Volume {
	log.V(logLevel).Debugf("Cache: VolumeCache: get secret: %s", hash)
	s.lock.Lock()
	defer s.lock.Unlock()
	pod, ok := s.volumes[hash]
	if !ok {
		return nil
	}
	return &pod
}

func (s *VolumesState) AddVolume(secret *types.Volume) {
	log.V(logLevel).Debugf("Cache: VolumeCache: add secret: %#v", secret)
	s.SetVolume(secret)
}

func (s *VolumesState) SetVolume(secret *types.Volume) {
	log.V(logLevel).Debugf("Cache: VolumeCache: set secret: %#v", secret)
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.volumes[secret.Meta.Name]; ok {
		delete(s.volumes, secret.Meta.Name)
	}

	s.volumes[secret.Meta.Name] = *secret
}

func (s *VolumesState) DelVolume(secret *types.Volume) {
	log.V(logLevel).Debugf("Cache: VolumeCache: del secret: %#v", secret)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.volumes[secret.Meta.Name]; ok {
		delete(s.volumes, secret.Meta.Name)
	}
}
