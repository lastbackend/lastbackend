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

package pod

import (
	"github.com/lastbackend/lastbackend/pkg/cache"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/scheduler/context"
)

type PodController struct {
	context *context.Context
	pods    chan *types.Pod

	pending *cache.PodCache

	active bool
}

func (pc *PodController) Watch(node chan *types.Node) {
	var (
		log = pc.context.GetLogger()
		stg = pc.context.GetStorage()
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

					// If pod state set to provision then need run provision action
					if p.State.Provision {
						log.Debugf("PodController: pod needs to be allocated to node: %s", p.Meta.Name)
						if err := Provision(p); err != nil {
							if err.Error() != errors.NodeNotFound {
								pc.pending.AddPod(p)
							} else {
								log.Errorf("Error: PodController: pod provision: %s", err.Error())
							}
							continue
						}

						pc.pending.DelPod(p)
						continue
					}

					// If pod state not set in provision and status
					// destroyed then need remove pod from node
					if p.State.State == types.StateDestroy {
						if err := Remove(p); err != nil {
							log.Errorf("Error: PodController: remove pod from node: %s", err.Error())
						}
						continue
					}

					// If pod state not set in provision and status
					// not destroyed then need update pod for node
					if err := Update(p); err != nil {
						log.Errorf("Error: PodController: update pod to node: %s", err.Error())
					}
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

	stg.Pod().Watch(pc.context.Background(), pc.pods)
}

func (pc *PodController) Pause() {
	pc.active = false
}

func (pc *PodController) Resume() {

	var (
		log = pc.context.GetLogger()
		stg = pc.context.GetStorage()
	)

	pc.active = true

	log.Debug("PodController: start check pods state")
	nss, err := stg.Namespace().List(pc.context.Background())
	if err != nil {
		log.Errorf("PodController: Get namespaces list err: %s", err.Error())
	}

	for _, ns := range nss {
		log.Debugf("Get pods in namespace: %s", ns.Meta.Name)
		pods, err := stg.Pod().ListByNamespace(pc.context.Background(), ns.Meta.Name)
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

func NewPodController(ctx *context.Context) *PodController {
	sc := new(PodController)
	sc.context = ctx
	sc.active = false
	sc.pods = make(chan *types.Pod)
	sc.pending = cache.NewPodCache()
	return sc
}
