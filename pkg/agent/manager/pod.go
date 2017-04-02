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

package manager

import (
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
)

type PodManager struct {
	update chan types.PodList
	close  chan bool
}

func NewPodManager() *PodManager {
	ctx := context.Get()
	ctx.Log.Info("Create new pod Manager")
	var pm = new(PodManager)

	pm.update = make(chan types.PodList)
	pm.close = make(chan bool)

	return pm
}

func ReleasePodManager(pm *PodManager) error {
	ctx := context.Get()
	ctx.Log.Info("Release pod manager")
	close(pm.update)
	close(pm.close)
	return nil
}

func (pm *PodManager) watch() error {
	ctx := context.Get()
	ctx.Log.Info("Start pod watcher")

	for {
		select {
		case _ = <-pm.close:
			return ReleasePodManager(pm)
		case pods := <-pm.update:
			{
				for _, pod := range pods {
					pm.sync(&pod)
				}
			}
		}
	}

	return nil
}

func (pm *PodManager) sync(p *types.Pod) {
	ctx := context.Get()
	ctx.Log.Infof("pod sync")

}
