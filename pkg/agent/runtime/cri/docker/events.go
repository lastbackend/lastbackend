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
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
)

func (r *Runtime) Subscribe() chan types.ContainerEvent {

	var (
		container *types.Container
	)

	log := context.Get().GetLogger()
	s   := context.Get().GetStorage().Pods()
	log.Debug("Create new event listener subscribe")

	var ch = make(chan types.ContainerEvent)
	go func() {
		es, errors := r.client.Events(context.Background(), docker.EventsOptions{})
		for {
			select {
			case e := <-es:

				log.Debugf("Event type: %s action: %s", e.Type, e.Action)
				if e.Type != events.ContainerEventType {
					continue
				}

				if (e.Action == "created") || (e.Action == "kill") {
					continue
				}

				container = s.GetContainer(e.ID)
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
