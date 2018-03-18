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
	pods    chan *types.Pod

	pending *cache.PodCache

	active bool
}

func (pc *Controller) Watch(node chan *types.Node) {
	var (
		stg = envs.Get().GetStorage()
	)

	log.Debug("PodController: start watch")
	go func() {
		for {
			select {
			case p := <-pc.pods:
				{
					if !pc.active {
						log.Debug("PodController: skip management cause it is in slave mode")
						pc.pending.DelPod(p)
						continue
					}

					// If pod state not set to provision then need skip
					if !p.State.Provision {
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
						pc.pods <- p
					}
				}
			}
		}
	}()

	stg.Pod().Watch(context.Background(), pc.pods)
}

func (pc *Controller) Pause() {
	pc.active = false
}

func (pc *Controller) Resume() {

	var (
		stg = envs.Get().GetStorage()
	)

	pc.active = true

	log.Debug("PodController: start check pods state")
	namespaces, err := stg.Namespace().List(context.Background())
	if err != nil {
		log.Errorf("PodController: Get apps list err: %s", err.Error())
	}

	for _, ns := range namespaces {
		log.Debugf("PodController: Get pods in app: %s", ns.Meta.Name)
		pods, err := stg.Pod().ListByNamespace(context.Background(), ns.Meta.Name)
		if err != nil {
			log.Errorf("PodController: Get pods list err: %s", err.Error())
		}

		for _, p := range pods {
			if p.State.Provision == true {
				pc.pods <- p
			}
		}
	}
}

func NewPodController(ctx context.Context) *Controller {
	sc := new(Controller)
	sc.context = ctx
	sc.active = false
	sc.pods = make(chan *types.Pod)
	sc.pending = cache.NewPodCache()
	return sc
}
