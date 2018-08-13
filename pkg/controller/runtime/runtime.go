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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/system"
)

type Runtime struct {
	ctx      context.Context
	process  *system.Process
	observer *Observer
	active   bool
}

func NewRuntime(ctx context.Context) *Runtime {
	r := new(Runtime)

	r.ctx = ctx
	r.process = new(system.Process)
	_, err := r.process.Register(ctx, types.KindController, envs.Get().GetStorage())
	if err != nil {
		log.Error(err)
		return nil
	}

	r.observer = NewObserver(ctx)

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
			case <-r.ctx.Done():
				return
			case l := <-lead:
				{

					if l {

						if r.active {
							log.Debug("Runtime: is already marked as lead -> skip")
							continue
						}

						log.Debug("Runtime: Mark as lead")

						r.active = true

					} else {

						if !r.active {
							log.Debug("Runtime: is already marked as slave -> skip")
							continue
						}

						log.Debug("Runtime: Mark as slave")

						r.active = false
					}

				}
			}
		}
	}()

	if err := r.process.WaitElected(r.ctx, lead); err != nil {
		log.Errorf("Runtime: Elect Wait error: %s", err.Error())
	}

}
