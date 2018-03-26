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

package runtime

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime/deployment"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime/service"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/system"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime/pod"
)

// watch service state and specs
// generate pods by specs

// watch service builds
// generate build spec after build creation

// watch service build state
// update pods after build passed state

type Runtime struct {
	process *system.Process

	sc      *service.Controller
	dc      *deployment.Controller
	pc      *pod.Controller

	active bool
}

func NewRuntime(ctx context.Context) *Runtime {
	r := new(Runtime)
	r.process = new(system.Process)
	r.process.Register(ctx, types.KindController, envs.Get().GetStorage())

	r.sc = service.NewServiceController(ctx)
	r.dc = deployment.NewDeploymentController(ctx)
	r.pc = pod.NewPodController(ctx)

	go r.sc.WatchSpec()
	go r.dc.WatchSpec()
	go r.dc.WatchStatus()
	go r.pc.WatchStatus()

	return r
}

func (r *Runtime) Loop() {

	var (
		lead = make(chan bool)
	)

	log.Debug("Controller: Runtime: Loop")

	go func() {
		for {
			select {
			case l := <-lead:
				{
					if l {

						if r.active {
							log.Debug("Runtime: is already marked as lead -> skip")
							continue
						}

						log.Debug("Runtime: Mark as lead")

						r.active = true
						r.sc.Resume()
						r.dc.Resume()
						r.pc.Resume()

					} else {

						if !r.active {
							log.Debug("Runtime: is already marked as slave -> skip")
							continue
						}

						log.Debug("Runtime: Mark as slave")

						r.active = false
						r.sc.Pause()
						r.dc.Pause()
						r.pc.Pause()
					}
				}
			}
		}
	}()

	if err := r.process.WaitElected(lead); err != nil {
		log.Errorf("Runtime: Elect Wait error: %s", err.Error())
	}
}
