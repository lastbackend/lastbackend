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

package runtime

import (
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"sync"
)

type PodManager struct {
	lock    sync.RWMutex
	workers map[string]*Worker
}

func (pm *PodManager) GetPodList() ([]*types.Pod) {
	pods := context.Get().GetStorage().Pods().GetPods()
	list := []*types.Pod{}

	for _, pod := range pods {
		list = append(list, pod)
	}

	return list
}

func (pm *PodManager) GetPods() (map[string]*types.Pod) {
	return context.Get().GetStorage().Pods().GetPods()
}

func (pm *PodManager) SyncPod(pod *types.PodNodeSpec) {
	log := context.Get().GetLogger()
	log.Debugf("Pod %s sync", pod.Meta.ID)

	p := context.Get().GetStorage().Pods().GetPod(pod.Meta.ID)

	if p == nil {
		log.Debugf("Pod %s not found, create new one", pod.Meta.ID)
		p := types.NewPod()
		p.Meta.ID = pod.Meta.ID
		context.Get().GetStorage().Pods().SetPod(p)
		pm.sync(pod.Meta, pod.Spec, p)
		return
	}

	log.Debugf("Pod %s found", pod.Meta.ID)
	if (p.Spec.ID == pod.Spec.ID) && p.Meta.State.State == pod.Meta.State.State {
		log.Debugf("Pod %s in correct state", pod.Meta.ID)
		return
	}
	pm.sync(pod.Meta, pod.Spec, p)
}

func (pm *PodManager) sync(meta types.PodMeta, spec types.PodSpec, pod *types.Pod) {
	// Create new worker to sync pod
	// Check if pod worker exists
	log := context.Get().GetLogger()
	log.Debugf("Pod %s sync start", pod.Meta.ID)
	w := pm.workers[pod.Meta.ID]
	if w == nil {
		log.Debugf("Pod %s sync create new worker", pod.Meta.ID)
		w = NewWorker()

		// Start worker watcher
		go func() {
			<-w.done
			pm.lock.Lock()
			delete(pm.workers, pod.Meta.ID)
			pm.lock.Unlock()
		}()
	}
	log.Debugf("Pod %s sync proceed", pod.Meta.ID)
	w.Proceed(meta, spec, pod)
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

	pods, err := crii.PodList()
	if err != nil {
		return pm, err
	}

	s := context.Get().GetStorage().Pods()
	s.SetPods(pods)

	return pm, nil
}
