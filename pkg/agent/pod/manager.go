package pod

import (
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"sync"
)

type PodManager struct {
	lock    sync.RWMutex
	workers map[types.PodID]*Worker
}

func (pm *PodManager) SyncPod(pod *types.Pod) {
	s := context.Get().GetStorage().Pods()
	pods := s.GetPods()

	p := pods[pod.Meta.ID]

	if p == nil {
		s.SetPod(pod)
		pm.sync(pod.Policy, pod.Spec, pod)
		return
	}

	if p.Spec.NotEqual(pod.Spec) {
		pm.sync(pod.Policy, pod.Spec, p)
	}
}

func (pm *PodManager) sync(policy types.PodPolicy, spec types.PodSpec, p *types.Pod) {
	// Create new worker to sync pod
	// Check if pod worker exists
	w := pm.workers[p.ID()]
	if w == nil {
		w = NewWorker()

		// Start worker watcher
		go func() {
			<-w.done
			pm.lock.Lock()
			delete(pm.workers, p.ID())
			pm.lock.Unlock()
		}()

	}

	w.Proceed(policy, spec, p)
}

func NewPodManager() (*PodManager, error) {
	log := context.Get().GetLogger()
	log.Debug("Create new pod manager")

	//	s := context.Get().GetStorage().Pods()

	var (
	//	err error
	)

	crii := context.Get().GetCri()

	pm := &PodManager{}

	pm.workers = make(map[types.PodID]*Worker)

	log.Debug("Restore new pod manager state")

	crii.PodList()
	return pm, nil
}
