//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package pod

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/cache"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/scheduler/envs"
)

type Controller struct {
	context context.Context
	spec    chan *types.Pod
	status  chan *types.Pod
	pending *cache.PodCache

	active bool
}

func (pc *Controller) WatchSpec(node chan *types.Node) {
	var (
		stg   = envs.Get().GetStorage()
		event = make(chan *types.Event)
	)

	log.Debug("PodController: start watch")
	go func() {
		for {
			select {
			case p := <-pc.spec:
				{
					if !pc.active {
						log.Debug("PodController: skip management cause it is in slave mode")
						pc.pending.DelPod(p)
						continue
					}

					log.Debugf("PodController: provision for pod: %s", p.Meta.Name)
					if err := Provision(p); err != nil {
						if err.Error() != errors.NodeNotFound {
							pc.pending.AddPod(p)
						} else {
							log.Errorf("PodController: pod provision err: %s", err.Error())
						}
						continue
					}

					pc.pending.DelPod(p)
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case _ = <-node:
				{
					for _, p := range pc.pending.GetPods() {
						pc.spec <- p
					}
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case e := <-event:
				if e.Data == nil {
					continue
				}

				pc.spec <- e.Data.(*types.Pod)
			}
		}
	}()

	stg.Pod().WatchSpec(context.Background(), event)
}

func (pc *Controller) WatchStatus(node chan *types.Node) {
	var (
		stg   = envs.Get().GetStorage()
		event = make(chan *types.Event)
	)

	log.Debug("PodController: start watch")
	go func() {
		for {
			select {
			case p := <-pc.status:
				{
					if !pc.active {
						log.Debug("PodController: skip management cause it is in slave mode")
						pc.pending.DelPod(p)
						continue
					}

					// If pod state not set to provision then need skip
					if p.Status.Stage == types.StateProvision {
						continue
					}

					log.Debugf("PodController: handle status for pod: %s", p.Meta.Name)
					if err := HandleStatus(p); err != nil {
						if err.Error() != errors.NodeNotFound {
							pc.pending.AddPod(p)
						} else {
							log.Errorf("PodController: pod provision err: %s", err.Error())
						}
						continue
					}

					pc.pending.DelPod(p)
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case _ = <-node:
				{
					for _, p := range pc.pending.GetPods() {
						pc.status <- p
					}
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case e := <-event:
				if e.Data == nil {
					continue
				}

				pc.status <- e.Data.(*types.Pod)
			}
		}
	}()

	stg.Pod().WatchStatus(context.Background(), event)
}

func (pc *Controller) Pause() {
	pc.active = false
}

func (pc *Controller) Resume() {

	var (
		stg = envs.Get().GetStorage()
		msg = "scheduler:controller:pod:resume"
	)

	pc.active = true

	log.Debug("PodController: start check pods state")
	namespaces, err := stg.Namespace().List(context.Background())
	if err != nil {
		log.Errorf("PodController: Get apps list err: %s", err.Error())
	}

	for _, ns := range namespaces {
		log.Debugf("PodController: Get pods in namespace: %s", ns.Meta.Name)
		pods, err := stg.Pod().ListByNamespace(context.Background(), ns.Meta.Name)
		if err != nil {
			log.Errorf("PodController: Get pods list err: %s", err.Error())
		}

		for _, p := range pods {

			log.Debugf("%s: restore pod: %s> status:[%s], state:[%s]", msg, p.SelfLink(), p.Status.Stage, p.Spec.State)
			if p.Status.Stage == types.StateProvision || p.Spec.State.Destroy {
				log.Debugf("%s: provision pod: %s", msg, p.SelfLink())
				pc.spec <- p
			}
		}

		for _, p := range pods {
			if p.Status.Stage != types.StateProvision {
				pc.status <- p
			}
		}
	}
}

func NewPodController(ctx context.Context) *Controller {
	sc := new(Controller)
	sc.context = ctx
	sc.active = false
	sc.spec = make(chan *types.Pod)
	sc.status = make(chan *types.Pod)
	sc.pending = cache.NewPodCache()
	return sc
}
