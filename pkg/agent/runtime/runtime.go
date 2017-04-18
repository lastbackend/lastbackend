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
	"github.com/lastbackend/lastbackend/pkg/agent/events"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"time"
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

func (r *Runtime) Sync(pods map[string]types.PodNodeSpec) {
	log := context.Get().GetLogger()
	log.Debug("Runtime: start sync")
	for _, pod := range pods {
		r.pManager.SyncPod(pod)
	}
}

func (r *Runtime) Loop() {

	log := context.Get().GetLogger()
	log.Debug("Runtime: start Loop")

	spec, err := events.New().Send(events.NewInitialEvent(GetNodeMeta(), r.pManager.GetPodList()))
	if err != nil {
		log.Errorf("Send initial event error %s", err.Error())
	}

	pods, host := r.eListener.Subscribe()

	go func() {
		log := context.Get().GetLogger()
		log.Debug("Runtime: Loop")
		ticker := time.NewTicker(time.Second * 10)

		go func() {
			for _ = range ticker.C {
				spec, err := events.New().Send(events.NewTickerEvent(GetNodeMeta()))

				if err != nil {
					log.Errorf("Runtime: send event error: %s", err.Error())
					continue
				}

				r.Sync(spec.Pods)
			}
		}()

		for {
			select {
			case pod := <-pods:
				log.Debugf("Runtime: Loop: send pod update event: %s", pod.Event)
				ps := []*types.Pod{}

				spec, err := events.New().Send(events.NewEvent(GetNodeMeta(), append(ps, &types.Pod{
					Meta:       pod.Meta,
					State:      pod.State,
					Containers: pod.Containers,
				})))

				if err != nil {
					log.Errorf("Runtime: send event error: %s", err.Error())
					continue
				}

				log.Debugf("pod contaienrs length: %d", len(pod.Containers))
				r.Sync(spec.Pods)

			case host := <-host:
				log.Debugf("Runtime: Loop: send host update event: %s", host.Event)
			}
		}
	}()

	r.Sync(spec.Pods)
}
