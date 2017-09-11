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
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/discovery/context"
	"github.com/lastbackend/lastbackend/pkg/discovery/endpoint"
	"github.com/lastbackend/lastbackend/pkg/discovery/pod"
	"github.com/lastbackend/lastbackend/pkg/discovery/service"
	"github.com/lastbackend/lastbackend/pkg/system"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const debugLevel = 2

// watch service state and specs
// generate pods by specs

// watch service builds
// generate build spec after build creation

// watch service build state
// update pods after build passed state

type Runtime struct {
	context *context.Context
	process *system.Process
	ec      *endpoint.EndpointController
	sc      *service.ServiceController
	pc      *pod.PodController

	active bool
}

func NewRuntime(ctx *context.Context) *Runtime {
	r := new(Runtime)
	r.context = ctx
	r.process = new(system.Process)
	r.process.Register(ctx, types.KindDiscovery)

	r.ec = endpoint.NewEndpointController(ctx)

	s := make(chan *types.Service)
	r.sc = service.NewServiceController(ctx)
	r.pc = pod.NewPodController(ctx)

	go r.ec.Watch()
	go r.sc.Watch(s)
	go r.pc.Watch(s)

	return r
}

func (r *Runtime) Loop() {

	var (
		lead = make(chan bool)
	)

	log.V(debugLevel).Debug("Runtime: Loop")

	go func() {
		for {
			select {
			case l := <-lead:
				{
					if l {
						if r.active {
							log.V(debugLevel).Debug("Runtime: is already marked as lead -> skip")
							continue
						}
						r.active = true
						log.V(debugLevel).Debug("Runtime: Mark as lead")
						r.ec.Resume()
						r.sc.Resume()
						r.pc.Resume()

					} else {
						if !r.active {
							log.V(debugLevel).Debug("Runtime: is already marked as slave -> skip")
							continue
						}
						log.V(debugLevel).Debug("Runtime: Mark as slave")
						r.active = false
						r.ec.Pause()
						r.sc.Pause()
						r.pc.Pause()
					}
				}
			}
		}
	}()

	if err := r.process.WaitElected(lead); err != nil {
		log.V(debugLevel).Errorf("Runtime: Elect Wait err: %s", err.Error())
	}
}
