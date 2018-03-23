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

package docker

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/node/state"

	d "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"os"
	"time"
)

func (r *Runtime) Subscribe(ctx context.Context, state *state.PodState, p chan string) {

	log.Debug("Create new event listener subscribe")

	go func() {

		if _, err := r.client.Ping(ctx); err != nil {
			log.Errorf("Can not ping docker client")
			return
		}

		es, errors := r.client.Events(ctx, d.EventsOptions{})
		for {
			select {
			case e := <-es:

				if e.Type != events.ContainerEventType {
					continue
				}

				log.Debugf("Event type: %s action: %s", e.Type, e.Action)

				container := state.GetContainer(e.ID)
				if container == nil {
					log.Debugf("Container not found")
					continue
				}

				log.Debugf("Container %s update in pod %s", container.ID, container.Pod)

				if e.Action == types.EventStateDestroy {
					state.DelContainer(container)
					p <- container.Pod
					break
				}

				c, err := r.ContainerInspect(ctx, e.ID)
				container.Pod = c.Pod

				switch c.State {
				case types.StateCreated:
					container.State = types.PodContainerState{
						Created: types.PodContainerStateCreated{
							Created: time.Now().UTC(),
						},
					}
				case types.StateStarted:
					if container.State.Started.Started {
						continue
					}
					container.State = types.PodContainerState{
						Started: types.PodContainerStateStarted{
							Started: true,
							Timestamp: time.Now().UTC(),
						},
					}
					container.State.Stopped.Stopped = false
				case types.StateStopped:
					if container.State.Stopped.Stopped {
						continue
					}
					container.State.Stopped.Stopped = true
					container.State.Stopped.Exit = types.PodContainerStateExit{
						Code:      c.ExitCode,
						Timestamp: time.Now().UTC(),
					}
					container.State.Started.Started = false
				case types.StateError:
					if container.State.Error.Error {
						continue
					}
					container.State.Error.Error = true
					container.State.Error.Message = c.Status
					container.State.Error.Exit = types.PodContainerStateExit{
						Code:      c.ExitCode,
						Timestamp: time.Now().UTC(),
					}
					container.State.Started.Started = false
					container.State.Stopped.Stopped = false
					container.State.Stopped.Exit = types.PodContainerStateExit{
						Code:      c.ExitCode,
						Timestamp: time.Now().UTC(),
					}
					container.State.Started.Started = false
				}

				if err != nil {
					log.Errorf("Container: can-not inspect")
					break
				}
				if c == nil {
					log.Errorf("Container: container not found")
					break
				}

				state.SetContainer(container)
				p <- container.Pod

				break

			case err := <-errors:
				log.Errorf("Event listening error: %s", err)
				os.Exit(0)
			}
		}
	}()
}
