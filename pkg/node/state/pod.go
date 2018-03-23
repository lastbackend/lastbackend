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

func (s *PodState) GetPodsCount() int {
	log.V(logLevel).Debugf("Cache: PodCache: get pods count: %d", s.stats.pods)
	return s.stats.pods
}

func (s *PodState) GetContainersCount() int {
	log.V(logLevel).Debugf("Cache: PodCache: get containers count: %d", s.stats.containers)
	return s.stats.containers
}

func (s *PodState) GetPods() map[string]types.PodStatus {
	log.V(logLevel).Debug("Cache: PodCache: get pods")
	return s.pods
}

func (s *PodState) SetPods(pods map[string]*types.PodStatus) {
	log.V(logLevel).Debugf("Cache: PodCache: set pods: %#v", pods)
	for key, pod := range pods {
		s.pods[key] = *pod
		s.stats.pods++
	}
}

func (s *PodState) GetPod(key string) *types.PodStatus {
	log.V(logLevel).Debugf("Cache: PodCache: get pod: %s", key)
	s.lock.Lock()
	defer s.lock.Unlock()
	pod, ok := s.pods[key]
	if !ok {
		return nil
	}
	return &pod
}

func (s *PodState) AddPod(key string, pod *types.PodStatus) {
	log.V(logLevel).Debugf("Cache: PodCache: add pod: %#v", pod)
	s.SetPod(key, pod)
}

func (s *PodState) SetPod(key string, pod *types.PodStatus) {
	log.V(logLevel).Debugf("Cache: PodCache: set pod: %#v", pod)
	s.lock.Lock()
	if _, ok := s.pods[key]; ok {
		delete(s.pods, key)
		s.stats.pods--
	}

	s.pods[key] = *pod
	s.stats.pods++

	s.lock.Unlock()
	for _, c := range pod.Containers {
		s.SetContainer(c)
	}
}

func (s *PodState) DelPod(key string) {
	log.V(logLevel).Debugf("Cache: PodCache: del pod: %s", key)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.pods[key]; ok {
		delete(s.pods, key)
		s.stats.pods--
	}
}

func (s *PodState) GetContainer(id string) *types.PodContainer {
	log.V(logLevel).Debugf("Cache: PodCache: get container: %s", id)
	c, ok := s.containers[id]
	if !ok {
		return nil
	}
	return &c
}

func (s *PodState) AddContainer(c *types.PodContainer) {
	log.V(logLevel).Debugf("Cache: PodCache: add container: %#v", c)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.containers[c.ID]; !ok {
		s.stats.containers++
	}
	s.containers[c.ID] = *c

}

func (s *PodState) SetContainer(c *types.PodContainer) {
	log.V(logLevel).Debugf("Cache: PodCache: set container: %#v", c)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.containers[c.ID]; !ok {
		s.stats.containers++
	}
	s.containers[c.ID] = *c
}

func (s *PodState) DelContainer(c *types.PodContainer) {
	log.V(logLevel).Debugf("Cache: PodCache: del container: %s", c.ID)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.containers[c.ID]; ok {
		delete(s.containers, c.ID)
		s.stats.containers--
	}
}
