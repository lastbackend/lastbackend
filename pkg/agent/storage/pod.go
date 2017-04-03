package storage

import (
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"sync"
)

type PodStorage struct {
	lock sync.RWMutex
	pods map[types.PodID]*types.Pod
}

func (ps *PodStorage) GetPods() map[types.PodID]*types.Pod {
	return ps.pods
}

func (ps *PodStorage) GetPod(id types.PodID) *types.Pod {
	pod, ok := ps.pods[id]
	if !ok {
		return nil
	}
	return pod
}

func (ps *PodStorage) AddPod(pod *types.Pod) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.pods[pod.ID()] = pod
}

func (ps *PodStorage) SetPod(pod *types.Pod) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.pods[pod.ID()] = pod
}

func (ps *PodStorage) DetPod(pod *types.Pod) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	delete(ps.pods, pod.ID())
}

func (ps *PodStorage) SetPods(pods []*types.Pod) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	for _, pod := range pods {
		ps.pods[pod.ID()] = pod
	}
}

func NewPodStorage() *PodStorage {
	pods := make(map[types.PodID]*types.Pod)
	return &PodStorage{
		pods: pods,
	}
}
