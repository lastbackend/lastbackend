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

package docker

import (
	docker "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/lastbackend/lastbackend/pkg/cache"
	"github.com/lastbackend/lastbackend/pkg/common/context"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

func (r *Runtime) Subscribe(ctx context.IContext, stg *cache.PodCache) chan types.ContainerEvent {

	var (
		container *types.Container
	)

	log.Debug("Create new event listener subscribe")

	var ch = make(chan types.ContainerEvent)
	go func() {
		es, errors := r.client.Events(ctx.Background(), docker.EventsOptions{})
		for {
			select {
			case e := <-es:

				log.Debugf("Event type: %s action: %s", e.Type, e.Action)
				if e.Type != events.ContainerEventType {
					continue
				}

				if (e.Action == types.EventStateCreated) || (e.Action == types.EventStateKill) {
					continue
				}

				container = stg.GetContainer(e.ID)
				if container == nil {
					log.Debugf("Container not found")
					continue
				}

				log.Debugf("Contaniner %s update in pod %s", container.ID, container.Pod)

				ch <- types.ContainerEvent{
					Event:     e.Action,
					Container: container,
				}

				break

			case err := <-errors:
				log.Errorf("Event listening error: %s", err.Error())
			}
		}
	}()

	return ch
}
