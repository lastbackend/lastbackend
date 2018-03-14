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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/scheduler/runtime/node"
	"github.com/lastbackend/lastbackend/pkg/scheduler/runtime/pod"
	"github.com/lastbackend/lastbackend/pkg/system"
	"github.com/lastbackend/lastbackend/pkg/log"
	"context"
	"github.com/lastbackend/lastbackend/pkg/scheduler/envs"
)

// watch service state and specs
// generate pods by specs

// watch service builds
// generate build spec after build creation

// watch service build state
// update pods after build passed state

type Runtime struct {
	process *system.Process

	pc *pod.PodController
	nc *node.NodeController

	active bool
}

func NewRuntime(ctx context.Context) *Runtime {
	r := new(Runtime)
	r.process = new(system.Process)
	r.process.Register(ctx, types.KindScheduler, envs.Get().GetStorage())

	r.pc = pod.NewPodController(ctx)
	r.nc = node.NewNodeController(ctx)

	n := make(chan *types.Node)
	go r.pc.Watch(n)
	go r.nc.Watch(n)

	return r
}

func (r *Runtime) Loop() {

	var (
		lead = make(chan bool)
	)

	log.Debug("Runtime: Loop")

	go func() {
		for {
			select {
			case l := <-lead:
				{
					if l {
						if r.active {
							log.Debug(" Runtime: is already marked as lead -> skip")
							continue
						}
						r.active = true
						log.Debug(" Runtime: Mark as lead")
						r.pc.Resume()
						r.nc.Resume()

					} else {
						if !r.active {
							log.Debug(" Runtime: is already marked as slave -> skip")
							continue
						}
						log.Debug(" Runtime: Mark as slave")
						r.active = false
						r.pc.Pause()
						r.nc.Pause()
					}
				}
			}
		}
	}()

	if err := r.process.WaitElected(lead); err != nil {
		log.Errorf("Runtime: Elect Wait error: %s", err.Error())
	}
}
