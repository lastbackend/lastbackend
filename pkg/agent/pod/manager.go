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

package pod

import (
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/satori/go.uuid"
	"sync"
)

type PodManager struct {
	lock    sync.RWMutex
	workers map[string]*Worker
}

func (pm *PodManager) SyncPod(pod *types.Pod) {
	s := context.Get().GetStorage().Pods()
	pods := s.GetPods()

	p, ok := pods[pod.Meta.ID]

	if !ok {
		p := types.NewPod()
		p.Meta = pod.Meta
		p.Meta.Spec = uuid.NewV4().String()
		s.SetPod(p)
		pm.sync(pod.State, pod.Meta, pod.Spec, pod)
		return
	}

	if p.Meta.Spec != (pod.Meta.Spec) || p.State.State != pod.State.State {
		pm.sync(pod.State, pod.Meta, pod.Spec, p)
	}
}

// meta - new pod meta
// spec - new pod spec
// pod - current pod information
func (pm *PodManager) sync(state types.PodState, meta types.PodMeta, spec types.PodSpec, pod *types.Pod) {
	// Create new worker to sync pod
	// Check if pod worker exists
	w := pm.workers[pod.Meta.ID]
	if w == nil {
		w = NewWorker()

		// Start worker watcher
		go func() {
			<-w.done
			pm.lock.Lock()
			delete(pm.workers, pod.Meta.ID)
			pm.lock.Unlock()
		}()
	}

	w.Proceed(state, meta, spec, pod)
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

	pm.workers = make(map[string]*Worker)

	log.Debug("Restore new pod manager state")

	crii.PodList()

	return pm, nil
}
