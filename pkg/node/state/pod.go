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

package state

import (
	"errors"
	"sync"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logPodPrefix = "node:state:pods:>"

type PodState struct {
	lock       sync.RWMutex
	stats      PodStateStats
	local      map[string]bool
	containers map[string]*types.PodContainer
	pods       map[string]*types.PodStatus
	watchers   map[chan string]bool
}

type PodStateStats struct {
	pods       int
	containers int
}

func (s *PodState) dispatch(pod string) {
	for w := range s.watchers {
		w <- pod
	}
}

func (s *PodState) Watch(watcher chan string, done chan bool) {
	s.watchers[watcher] = true
	defer delete(s.watchers, watcher)
	<-done
}

func (s *PodState) GetPodsCount() int {
	log.V(logLevel).Debugf("%s: get pods count: %d", logPodPrefix, s.stats.pods)
	return s.stats.pods
}

func (s *PodState) GetContainersCount() int {
	log.V(logLevel).Debugf("%s: get containers count: %d", logPodPrefix, s.stats.containers)
	return s.stats.containers
}

func (s *PodState) GetPods() map[string]*types.PodStatus {
	log.V(logLevel).Debugf("%s: get pods", logPodPrefix)
	return s.pods
}

func (s *PodState) SetPods(pods map[string]*types.PodStatus) {
	log.V(logLevel).Debugf("%s: set pods: %d", logPodPrefix, len(pods))
	for key, pod := range pods {
		state(pod)
		s.pods[key] = pod
		s.stats.pods++
	}
}

func (s *PodState) GetPod(key string) *types.PodStatus {
	log.V(logLevel).Debugf("%s: get pod: %s", logPodPrefix, key)
	s.lock.Lock()
	defer s.lock.Unlock()
	pod, ok := s.pods[key]
	if !ok {
		return nil
	}
	return pod
}

func (s *PodState) AddPod(key string, pod *types.PodStatus) {
	log.V(logLevel).Debugf("%s: add pod: %s: %s ", logPodPrefix, key, pod.Status)
	s.SetPod(key, pod)
}

func (s *PodState) SetLocal(key string) {
	log.V(logLevel).Debugf("%s: set pod: %s as local", logPodPrefix, key)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.local[key] = true
}

func (s *PodState) IsLocal(key string) bool {
	log.V(logLevel).Debugf("%s: check pod: %s is local", logPodPrefix, key)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.local[key]; ok {
		return true
	}

	return false
}

func (s *PodState) SetPod(key string, pod *types.PodStatus) {
	log.V(logLevel).Debugf("%s: set pod %s: %s", logPodPrefix, key, pod.Status)

	s.lock.Lock()
	if _, ok := s.pods[key]; ok {
		delete(s.pods, key)
		s.stats.pods--
	}

	s.pods[key] = pod
	s.stats.pods++

	s.lock.Unlock()
	for _, c := range pod.Containers {
		s.SetContainer(c)
	}

	s.lock.Lock()
	state(pod)
	s.lock.Unlock()
	s.dispatch(key)
}

func (s *PodState) DelPod(key string) {
	log.V(logLevel).Debugf("%s: del pod: %s", logPodPrefix, key)
	s.lock.Lock()
	if _, ok := s.pods[key]; ok {
		delete(s.pods, key)
		s.stats.pods--
	}
	s.lock.Unlock()
	s.dispatch(key)
}

func (s *PodState) GetContainer(id string) *types.PodContainer {
	log.V(logLevel).Debugf("%s: get container: %s", logPodPrefix, id)
	c, ok := s.containers[id]
	if !ok {
		return nil
	}
	return c
}

func (s *PodState) AddContainer(c *types.PodContainer) {
	log.V(logLevel).Debugf("%s: add container: %s", logPodPrefix, c.ID)
	s.lock.Lock()
	if _, ok := s.containers[c.ID]; !ok {
		s.stats.containers++
	}
	s.containers[c.ID] = c

	s.lock.Unlock()
}

func (s *PodState) SetContainer(c *types.PodContainer) {
	log.V(logLevel).Debugf("%s: set container: %s", logPodPrefix, c.ID)
	s.lock.Lock()

	if _, ok := s.containers[c.ID]; !ok {
		s.stats.containers++
	}
	s.containers[c.ID] = c

	s.lock.Unlock()
}

func (s *PodState) DelContainer(c *types.PodContainer) {
	log.V(logLevel).Debugf("%s: del container: %s", logPodPrefix, c.ID)
	s.lock.Lock()
	if _, ok := s.containers[c.ID]; ok {
		delete(s.containers, c.ID)
		s.stats.containers--
	}
	s.lock.Unlock()
	pod := s.GetPod(c.Pod)
	if pod == nil {
		return
	}

	s.lock.Lock()
	delete(pod.Containers, c.ID)
	state(pod)
	s.lock.Unlock()
}

func state(s *types.PodStatus) {

	var sts = make(map[string]int)
	var ems string

	switch s.State {
	case types.StateExited:
		return
	case types.StateDestroyed:
		return
	case types.StateError:
		return
	case types.StateProvision:
		return
	case types.StateCreated:
		return
	case types.StatusPull:
		return
	}

	if len(s.Containers) == 0 {
		s.State = types.StateDegradation
		return
	}

	for _, cn := range s.Containers {

		switch true {
		case cn.State.Error.Error:
			sts[types.StateError] += 1
			ems = cn.State.Error.Message
			break
		case cn.State.Stopped.Stopped:
			sts[types.StatusStopped] += 1
			break
		case cn.State.Started.Started:
			sts[types.StateStarted] += 1
			break
		}
	}

	switch true {
	case len(s.Containers) == sts[types.StateError]:
		s.SetError(errors.New(ems))
		break
	case len(s.Containers) == sts[types.StateStarted]:
		s.SetRunning()
		break
	case len(s.Containers) == sts[types.StatusStopped]:
		s.SetStopped()
		break
	}
}
