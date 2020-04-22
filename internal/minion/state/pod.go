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
	"context"
	"errors"
	"github.com/lastbackend/lastbackend/tools/logger"
	"sync"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

const logPodPrefix = "node:state:pods:>"

type PodState struct {
	lock       sync.RWMutex
	stats      PodStateStats
	local      map[string]bool
	containers map[string]*models.PodContainer
	pods       map[string]*models.PodStatus
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
	log := logger.WithContext(context.Background())
	log.Debugf("%s: get pods count: %d", logPodPrefix, s.stats.pods)
	return s.stats.pods
}

func (s *PodState) GetContainersCount() int {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: get containers count: %d", logPodPrefix, s.stats.containers)
	return s.stats.containers
}

func (s *PodState) GetPods() map[string]*models.PodStatus {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: get pods", logPodPrefix)
	return s.pods
}

func (s *PodState) SetPods(pods map[string]*models.PodStatus) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: set pods: %d", logPodPrefix, len(pods))
	for key, pod := range pods {
		state(pod)
		s.pods[key] = pod
		s.stats.pods++
	}
}

func (s *PodState) GetPod(key string) *models.PodStatus {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: get pod: %s", logPodPrefix, key)
	s.lock.Lock()
	defer s.lock.Unlock()
	pod, ok := s.pods[key]
	if !ok {
		return nil
	}
	return pod
}

func (s *PodState) AddPod(key string, pod *models.PodStatus) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: add pod: %s: %s ", logPodPrefix, key, pod.Status)
	s.SetPod(key, pod)
}

func (s *PodState) SetLocal(key string) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: set pod: %s as local", logPodPrefix, key)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.local[key] = true
}

func (s *PodState) IsLocal(key string) bool {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: check pod: %s is local", logPodPrefix, key)
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.local[key]; ok {
		return true
	}

	return false
}

func (s *PodState) SetPod(key string, pod *models.PodStatus) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: set pod %s: %s", logPodPrefix, key, pod.Status)

	s.lock.Lock()
	if _, ok := s.pods[key]; ok {
		delete(s.pods, key)
		s.stats.pods--
	}

	s.pods[key] = pod
	s.stats.pods++

	s.lock.Unlock()
	for _, c := range pod.Runtime.Services {
		s.SetContainer(c)
	}

	s.lock.Lock()
	state(pod)
	s.lock.Unlock()
	s.dispatch(key)
}

func (s *PodState) DelPod(key string) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: del pod: %s", logPodPrefix, key)
	s.lock.Lock()
	if _, ok := s.pods[key]; ok {
		delete(s.pods, key)
		s.stats.pods--
	}
	s.lock.Unlock()
	s.dispatch(key)
}

func (s *PodState) GetContainer(id string) *models.PodContainer {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: get container: %s", logPodPrefix, id)
	c, ok := s.containers[id]
	if !ok {
		return nil
	}
	return c
}

func (s *PodState) AddContainer(c *models.PodContainer) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: add container: %s", logPodPrefix, c.ID)
	s.lock.Lock()
	if _, ok := s.containers[c.ID]; !ok {
		s.stats.containers++
	}
	s.containers[c.ID] = c

	s.lock.Unlock()
}

func (s *PodState) SetContainer(c *models.PodContainer) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: set container: %s", logPodPrefix, c.ID)
	s.lock.Lock()

	if _, ok := s.containers[c.ID]; !ok {
		s.stats.containers++
	}
	s.containers[c.ID] = c

	s.lock.Unlock()
}

func (s *PodState) DelContainer(c *models.PodContainer) {
	log := logger.WithContext(context.Background())
	log.Debugf("%s: del container: %s", logPodPrefix, c.ID)
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
	delete(pod.Runtime.Services, c.ID)
	state(pod)
	s.lock.Unlock()
}

func state(s *models.PodStatus) {

	var sts = make(map[string]int)
	var ems string

	switch s.State {
	case models.StateExited:
		return
	case models.StateDestroyed:
		return
	case models.StateError:
		return
	case models.StateProvision:
		return
	case models.StateCreated:
		return
	case models.StatusPull:
		return
	}

	if len(s.Runtime.Services) == 0 {
		s.State = models.StateDegradation
		return
	}

	for _, cn := range s.Runtime.Services {

		switch true {
		case cn.State.Error.Error:
			sts[models.StateError] += 1
			ems = cn.State.Error.Message
			break
		case cn.State.Stopped.Stopped:
			sts[models.StatusStopped] += 1
			break
		case cn.State.Started.Started:
			sts[models.StateStarted] += 1
			break
		}
	}

	switch true {
	case len(s.Runtime.Services) == sts[models.StateError]:
		s.SetError(errors.New(ems))
		break
	case len(s.Runtime.Services) == sts[models.StateStarted]:
		s.SetRunning()
		break
	case len(s.Runtime.Services) == sts[models.StatusStopped]:
		s.SetStopped()
		break
	}
}
