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
		spec:   make(chan *types.NodeSpec),
	}
}

type Runtime struct {
	pods   chan *types.Pod
	events chan *types.Event

	pManager  *PodManager
	eListener *EventListener

	spec chan *types.NodeSpec
}

func Get() *Runtime {
	return &runtime
}

func (r *Runtime) GetSpecChan() chan *types.NodeSpec {
	return r.spec
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

func (r *Runtime) Recovery(pods map[string]types.PodNodeSpec) {

	log := context.Get().GetLogger()
	ps := context.Get().GetStorage().Pods().GetPods()

	for _, pod := range ps {
		if _, ok := pods[pod.Meta.ID]; !ok {
			log.Debugf("Mark pod %s for removable", pod.Meta.ID)
			pods[pod.Meta.ID] = types.PodNodeSpec{
				Meta: pod.Meta,
				State: types.PodState{
					State: types.StateDestroy,
				},
				Spec: types.PodSpec{
					State: types.StateDestroy,
				},
			}
		}
	}

	r.Sync(pods)
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

	var recovery = false

	go r.HeartBeat()
	go r.Events()

	go func() {
		for {
			select {
			case spec := <-r.spec:
				log.Debug("Runtime: Loop: new spec received")
				if !recovery {
					recovery = true
					log.Debug("Runtime: Loop: recovery state")
					r.Recovery(spec.Pods)
				} else {
					log.Debug("Runtime: Loop: sync pods")
					r.Sync(spec.Pods)
				}
			}

		}
	}()

	events.NewInitialEvent(r.pManager.GetPodList())
}

func (r *Runtime) HeartBeat() {
	ticker := time.NewTicker(time.Second * 10)
	for _ = range ticker.C {
		events.NewTickerEvent()
	}
}

func (r *Runtime) Events() {

	log := context.Get().GetLogger()
	log.Debug("Runtime: Events listener")
	ev, host := r.eListener.Subscribe()

	for {
		select {
		case e := <-ev:
			log.Debugf("Runtime: Loop: send pod update event: %s", e.Event)
			events.SendPodState(e.Pod)
		case host := <-host:
			log.Debugf("Runtime: Loop: send host update event: %s", host.Event)
		}
	}
}
