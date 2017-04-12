package docker

import (
	docker "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
)

func (r *Runtime) Subscribe() chan types.ContainerEvent {
	log := context.Get().GetLogger()
	log.Debug("Create new event listener subscribe")

	var ch = make(chan types.ContainerEvent)
	go func() {
		es, errors := r.client.Events(context.Background(), docker.EventsOptions{})
		for {
			select {
			case e := <-es:

				log.Debugf("Event type: %s", e.Type)
				if e.Type != events.ContainerEventType {
					continue
				}

				log.Debugf("Proceed container update: %s", e.ID)
				container, pod, err := r.ContainerInspect(e.ID)
				if err != nil || container == nil {
					continue
				}

				log.Debugf("Contaniner %s update in pod %s", container.ID, pod)

				ch <- types.ContainerEvent{
					Pod:       pod,
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
