package manager

import (
	"github.com/golang/glog"
	"github.com/lastbackend/lastbackend/pkg/api/types"
)


type PodManager struct {
	update chan types.PodList
	close  chan bool
}

func NewPodManager() *PodManager {
	glog.V(4).Info("Create new pod Manager")
	var pm = new(PodManager)

	pm.update = make(chan types.PodList)
	pm.close  = make(chan bool)

	return pm
}

func ReleasePodManager (pm *PodManager) error {
	glog.V(4).Info("Release pod manager")
	close (pm.update)
	close (pm.close)
	return nil
}


func (pm *PodManager) watch () error {
	glog.V(4).Info("Start new pod watcher")

	for {
		select {
		case _= <- pm.close: return ReleasePodManager(pm)
		case pods := <- pm.update: {
			for _, pod := range pods {
				pm.sync(&pod)
			}
		}
		}
	}

	return nil
}

func (pm *PodManager) sync (p *types.Pod) {
	glog.V(4).Infof("Pod update: %s", p)

}