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

package runtime

import (
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
)

var runtime Runtime

func init() {
	runtime = Runtime{
		pods:   make(chan *types.Pod),
		events: make(chan *types.Event),
	}
}

type Runtime struct {
	pods   chan *types.Pod
	events chan *types.Event

	pManager  *PodManager
	eListener *EventListener
}

func Get() *Runtime {
	return &runtime
}

func (r *Runtime) StartPodManager() error {
	var err error
	if r.pManager, err = NewPodManager(); err != nil {
		return err
	}
	return nil
}

func (r *Runtime) StartEventListener() error {
	var err error
	if r.eListener, err = NewEventListener(); err != nil {
		return err
	}

	return nil
}

func (r *Runtime) Sync(pods []*types.Pod) {
	for _, pod := range pods {
		r.pManager.SyncPod(pod)
	}
}

func (r *Runtime) Loop() {

	pods, host := r.eListener.Subscribe()
	go func() {
		log := context.Get().GetLogger()
		log.Debug("Runtime: Loop")

		for {
			select {
			case pod := <-pods:
				log.Debugf("Runtime: Loop: send pod update event: %s", pod.Event)

			case host := <-host:
				log.Debugf("Runtime: Loop: send host update event: %s", host.Event)
			}
		}
	}()
}
