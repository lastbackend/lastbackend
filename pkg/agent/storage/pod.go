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
	lock sync.RWMutex
	pods map[string]*types.Pod
}

func (ps *PodStorage) GetPods() map[string]*types.Pod {
	return ps.pods
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
}

func (ps *PodStorage) SetPod(pod *types.Pod) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.pods[pod.Meta.ID] = pod
}

func (ps *PodStorage) DetPod(pod *types.Pod) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	delete(ps.pods, pod.Meta.ID)
}

func (ps *PodStorage) SetPods(pods []*types.Pod) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	for _, pod := range pods {
		ps.pods[pod.Meta.ID] = pod
	}
}

func NewPodStorage() *PodStorage {
	pods := make(map[string]*types.Pod)
	return &PodStorage{
		pods: pods,
	}
}
