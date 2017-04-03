package pod

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"sync"
)

type PodManager struct {
	lock sync.RWMutex
	pods map[string]*types.Pod
}

func (pm *PodManager) GetPods() types.PodList {
	return nil
}

func (pm *PodManager) SetPods(pods types.PodList) {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pm.pods = make(map[string]*types.Pod)

	for _, pod := range pods {
		pm.SetPod(pod)
	}

}

func (pm *PodManager) GetPod(uuid string) *types.Pod {
	return nil
}

func (pm *PodManager) AddPod(pod *types.Pod) {

}

func (pm *PodManager) SetPod(pod *types.Pod) {
	pm.pods[pod.Meta.ID] = pod
}

func (pm *PodManager) DelPod(pod *types.Pod) {

}

func (pm *PodManager) SyncPod(pod *types.Pod) {
	p := pm.pods[pod.Meta.ID]

	if p == nil {
		pm.SetPod(pod)
		return
	}

	ohash, _ := json.Marshal(p.Spec)
	nhash, _ := json.Marshal(pod.Spec)

	if string(ohash) == string(nhash) {
		return
	}

}

func NewPodManager() Manager {
	pm := &PodManager{}
	pm.SetPods(nil)

	return pm
}
