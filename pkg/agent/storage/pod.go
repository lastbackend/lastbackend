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

package storage

import (
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"sync"
)

type PodStorage struct {
	lock       sync.RWMutex
	stats      PodStorageStats
	containers map[string]*types.Container
	pods       map[string]*types.Pod
}

type PodStorageStats struct {
	pods       int
	containers int
}

func (ps *PodStorage) GetPodsCount() int {
	return ps.stats.pods
}

func (ps *PodStorage) GetContainersCount() int {
	return ps.stats.containers
}

func (ps *PodStorage) GetPods() map[string]*types.Pod {
	return ps.pods
}

func (ps *PodStorage) GetContainer(id string) *types.Container {
	c, ok := ps.containers[id]
	if !ok {
		return nil
	}
	return c
}

func (ps *PodStorage) AddContainer(c *types.Container) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	ps.containers[c.ID] = c
}

func (ps *PodStorage) DelContainer(id string) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	delete(ps.containers, id)
}

func (ps *PodStorage) GetPod(id string) *types.Pod {
	pod, ok := ps.pods[id]
	if !ok {
		return nil
	}
	return pod
}

func (ps *PodStorage) AddPod(pod *types.Pod) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.pods[pod.Meta.ID] = pod
	ps.stats.pods++
	ps.stats.containers += len(pod.Containers)

	for _, c := range pod.Containers {
		ps.containers[c.ID] = c
	}
}

func (ps *PodStorage) SetPod(pod *types.Pod) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if p, ok := ps.pods[pod.Meta.ID]; ok {
		ps.stats.containers--
		for _, c := range p.Containers {
			delete(ps.containers, c.ID)
		}
	}

	ps.pods[pod.Meta.ID] = pod
	for _, c := range pod.Containers {
		ps.containers[c.ID] = c
	}
}

func (ps *PodStorage) DelPod(pod *types.Pod) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if p, ok := ps.pods[pod.Meta.ID]; ok {
		ps.stats.containers--
		for _, c := range p.Containers {
			delete(ps.containers, c.ID)
		}
	}

	delete(ps.pods, pod.Meta.ID)
	ps.stats.pods--
}

func (ps *PodStorage) SetPods(pods []*types.Pod) {
	for _, pod := range pods {
		ps.AddPod(pod)
	}
}

func NewPodStorage() *PodStorage {
	pods := make(map[string]*types.Pod)
	containers := make(map[string]*types.Container)
	return &PodStorage{
		stats:      PodStorageStats{},
		containers: containers,
		pods:       pods,
	}
}
