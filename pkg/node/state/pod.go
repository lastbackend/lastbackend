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

func (s *PodState) GetPods() map[string]types.Pod {
	log.V(logLevel).Debug("Cache: PodCache: get pods")
	return s.pods
}

func (s *PodState) SetPods(pods []*types.Pod) {
	log.V(logLevel).Debugf("Cache: PodCache: set pods: %#v", pods)
	for _, pod := range pods {
		s.pods[pod.Meta.Name] = *pod
		s.stats.pods++
	}
}

func (s *PodState) GetPod(id string) *types.Pod {
	log.V(logLevel).Debugf("Cache: PodCache: get pod: %s", id)
	s.lock.Lock()
	defer s.lock.Unlock()
	pod, ok := s.pods[id]
	if !ok {
		return nil
	}
	return &pod
}

func (s *PodState) AddPod(pod *types.Pod) {
	log.V(logLevel).Debugf("Cache: PodCache: add pod: %#v", pod)
	s.SetPod(pod)
}

func (s *PodState) SetPod(pod *types.Pod) {
	log.V(logLevel).Debugf("Cache: PodCache: set pod: %#v", pod)
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.pods[pod.Meta.Name]; ok {
		delete(s.pods, pod.Meta.Name)
		s.stats.pods--
	}

	s.pods[pod.Meta.Name] = *pod
	s.stats.pods++
}

func (s *PodState) DelPod(pod *types.Pod) {
	log.V(logLevel).Debugf("Cache: PodCache: del pod: %#v", pod)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.pods[pod.Meta.Name]; ok {
		delete(s.pods, pod.Meta.Name)
		s.stats.pods--
	}
}

func (s *PodState) GetContainer(id string) *types.Container {
	log.V(logLevel).Debugf("Cache: PodCache: get container: %s", id)
	c, ok := s.containers[id]
	if !ok {
		return nil
	}
	return &c
}

func (s *PodState) AddContainer(c *types.Container) {
	log.V(logLevel).Debugf("Cache: PodCache: add container: %#v", c)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.containers[c.ID]; !ok {
		s.stats.containers++
	}
	s.containers[c.ID] = *c

}

func (s *PodState) SetContainer(c *types.Container) {
	log.V(logLevel).Debugf("Cache: PodCache: set container: %#v", c)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.containers[c.ID]; !ok {
		s.stats.containers++
	}
	s.containers[c.ID] = *c
}

func (s *PodState) DelContainer(c *types.Container) {
	log.V(logLevel).Debugf("Cache: PodCache: del container: %s", c.ID)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.containers[c.ID]; ok {
		delete(s.containers, c.ID)
		s.stats.containers--
	}
}
