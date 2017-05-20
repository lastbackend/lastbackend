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
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/util/system"
	"sync"
)

type PodManager struct {
	lock    sync.RWMutex
	workers map[string]*Worker
}

func (pm *PodManager) GetPodList() []*types.Pod {
	pods := context.Get().GetCache().Pods().GetPods()
	list := []*types.Pod{}

	for _, pod := range pods {
		list = append(list, pod)
	}

	return list
}

func (pm *PodManager) GetPods() map[string]*types.Pod {
	return context.Get().GetCache().Pods().GetPods()
}

func (pm *PodManager) SyncPod(pod types.PodNodeSpec) {

	log := context.Get().GetLogger()
	log.Debugf("Pod %s sync", pod.Meta.Name)

	p := context.Get().GetCache().Pods().GetPod(pod.Meta.Name)
	if p == nil {
		log.Debugf("Pod %s not found, create new one", pod.Meta.Name)
		p := types.NewPod()
		p.Meta = pod.Meta
		p.Meta.Hostname, _ = system.GetHostname()
		context.Get().GetCache().Pods().SetPod(p)
		pm.sync(pod.Meta, pod.State, pod.Spec, p)
		return
	}

	if p.State.Provision {
		log.Debugf("Pod %s is not in %s state > skip sync", p.Meta.Name, types.StateReady)
		return
	}

	log.Debugf("Pod %s found", pod.Meta.Name)

	if len(pod.Spec.Containers) != len(p.Containers) {
		log.Debugf("Pod %s containers len different from spec count %d(%d)", pod.Meta.Name, len(p.Containers),
			len(p.Containers))

		pm.sync(pod.Meta, pod.State, pod.Spec, p)
		return
	}

	if (p.Spec.ID == pod.Spec.ID) && (p.Spec.State == pod.Spec.State) {
		log.Debugf("Pod %s in correct state", pod.Meta.Name)
		return
	}

	if p.Spec.ID != pod.Spec.ID {
		log.Debugf("Pod %s need to spec update: %s (%s) ", pod.Meta.Name, pod.Spec.ID, p.Spec.ID)
	}

	if p.Spec.State != pod.Spec.State {
		log.Debugf("Pod %s need to change state to: %s (%s) ", pod.Meta.Name, pod.Spec.State, p.Spec.State)
	}

	pm.sync(pod.Meta, pod.State, pod.Spec, p)
}

func (pm *PodManager) sync(meta types.PodMeta, state types.PodState, spec types.PodSpec, pod *types.Pod) {
	// Create new worker to sync pod
	// Check if pod worker exists
	log := context.Get().GetLogger()
	log.Debugf("Pod %s sync start", pod.Meta.Name)
	w, ok := pm.workers[pod.Meta.Name]
	if !ok {
		log.Debugf("Pod %s sync create new worker", pod.Meta.Name)
		w = NewWorker()
		w.pod = pod.Meta.Name
		pm.lock.Lock()
		pm.workers[pod.Meta.Name] = w
		pm.lock.Unlock()

		// Start worker watcher
		go func() {
			<-w.done
			log.Debugf("Pod %s worker deletion", pod.Meta.Name)
			pm.lock.Lock()
			delete(pm.workers, pod.Meta.Name)
			pm.lock.Unlock()
		}()
	}

	log.Debugf("Pod %s sync proceed", pod.Meta.Name)
	w.Provision(meta, state, spec, pod)
}

func NewPodManager() (*PodManager, error) {

	log := context.Get().GetLogger()

	log.Debug("Create new pod manager")

	crii := context.Get().GetCri()

	pm := &PodManager{}
	pm.workers = make(map[string]*Worker)

	log.Debug("Restore new pod manager state")

	pods, err := crii.PodList(context.Get())
	if err != nil {
		return pm, err
	}

	log.Debugf("Runtime: new pods manager: restore state: %d pods found", len(pods))
	s := context.Get().GetCache().Pods()
	s.SetPods(pods)

	return pm, nil
}
