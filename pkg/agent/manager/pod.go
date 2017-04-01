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
	"fmt"
	_types "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
)

type PodManager struct {
	client *client.Client

	sync  chan types.PodList
	close chan bool
}

func NewPodManager(c *client.Client) *PodManager {
	ctx := context.Get()
	ctx.Log.Debug("Create new pod Manager")

	var pm = &PodManager{client: c}
	pm.sync = make(chan types.PodList)
	pm.close = make(chan bool)

	return pm
}

func ReleasePodManager(pm *PodManager) error {
	ctx := context.Get()
	ctx.Log.Debug("Release pod manager")
	close(pm.sync)
	close(pm.close)
	return nil
}

func (pm *PodManager) Run() {
	ctx := context.Get()
	ctx.Log.Debug("Restore pod manager state")

	containers, err := pm.client.ContainerList(context.Background(), _types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}

func (pm *PodManager) watch() error {
	ctx := context.Get()
	ctx.Log.Debug("Start pod watcher")

	for {
		select {
		case _ = <-pm.close:
			return ReleasePodManager(pm)
		case pods := <-pm.sync:
			{
				for _, pod := range pods {
					pm.patch(&pod)
				}
			}
		}
	}

	return nil
}

func (pm *PodManager) patch(p *types.Pod) {
	ctx := context.Get()
	ctx.Log.Debug("pod sync")

}
