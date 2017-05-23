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

package cache

import (
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"sync"
)

type PodCache struct {
	lock       sync.RWMutex
	stats      PodCacheStats
	containers map[string]*types.Container
	pods       map[string]*types.Pod
}

type PodCacheStats struct {
	pods       int
	containers int
}

func (ps *PodCache) GetPodsCount() int {
	return ps.stats.pods
}

func (ps *PodCache) GetContainersCount() int {
	return ps.stats.containers
}

func (ps *PodCache) GetPods() map[string]*types.Pod {
	return ps.pods
}

func (ps *PodCache) GetContainer(id string) *types.Container {
	c, ok := ps.containers[id]
	if !ok {
		return nil
	}
	return c
}

func (ps *PodCache) AddContainer(c *types.Container) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	ps.addContainer(c)
}

func (ps *PodCache) DelContainer(id string) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	ps.delContainer(id)
}

func (ps *PodCache) GetPod(id string) *types.Pod {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	pod, ok := ps.pods[id]
	if !ok {
		return nil
	}
	return pod
}

func (ps *PodCache) AddPod(pod *types.Pod) {
	ps.SetPod(pod)
}

func (ps *PodCache) SetPod(pod *types.Pod) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if p, ok := ps.pods[pod.Meta.Name]; ok {
		ps.delPod(p)
	}
	ps.addPod(pod)
}

func (ps *PodCache) DelPod(pod *types.Pod) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	if p, ok := ps.pods[pod.Meta.Name]; ok {
		ps.delPod(p)
	}
}

func (ps *PodCache) SetPods(pods []*types.Pod) {
	for _, pod := range pods {
		ps.AddPod(pod)
	}
}

func (ps *PodCache) addPod(pod *types.Pod) {
	ps.pods[pod.Meta.Name] = pod
	ps.stats.pods++
	for _, c := range pod.Containers {
		ps.addContainer(c)
	}
}

func (ps *PodCache) delPod(pod *types.Pod) {
	for _, c := range pod.Containers {
		ps.delContainer(c.ID)
	}
	delete(ps.pods, pod.Meta.Name)
	ps.stats.pods--
}

func (ps *PodCache) addContainer(c *types.Container) {
	if _, ok := ps.containers[c.ID]; !ok {
		ps.stats.containers++
	}
	ps.containers[c.ID] = c
}

func (ps *PodCache) delContainer(id string) {
	if _, ok := ps.containers[id]; ok {
		delete(ps.containers, id)
		ps.stats.containers--
	}
}

func NewPodCache() *PodCache {
	return &PodCache{
		stats:      PodCacheStats{},
		containers: make(map[string]*types.Container),
		pods:       make(map[string]*types.Pod),
	}
}
