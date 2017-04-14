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
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
)

type EventListener struct {
	pods chan *types.PodEvent
	host chan *types.HostEvent
}

func (el *EventListener) Subscribe() (chan *types.PodEvent, chan *types.HostEvent) {
	log := context.Get().GetLogger()
	log.Debug("Runtime: EventListener: Subscribe")

	return el.pods, el.host
}

func (el *EventListener) Listen() {
	log := context.Get().GetLogger()
	log.Debug("Runtime: EventListener: Listen")

	pods := context.Get().GetStorage().Pods()

	crii := context.Get().GetCri()

	events := crii.Subscribe()
	go func() {
		for {
			select {
			case event := <-events:
				{
					log.Debugf("Runtime: New event receive: %s", event.Event)
					if (event.Event != "start") && (event.Event != "stop") {
						continue
					}

					log.Debugf("Runtime: New event %s type proceed", event.Event)
					pod := pods.GetPod(event.Pod)
					if pod == nil {
						log.Debugf("Runtime: Pod %s not found", event.Pod)
						continue
					}
					log.Debugf("Runtime: Pod %s found > update container", event.Pod)
					pod.SetContainer(event.Container)
					pod.UpdateState()

					el.pods <- &types.PodEvent{
						Meta:       pod.Meta,
						Containers: pod.Containers,
					}

					jp, _ := json.Marshal(pod)
					log.Debug(string(jp))
				}
			}
		}
	}()
}

func NewEventListener() (*EventListener, error) {

	log := context.Get().GetLogger()
	log.Debug("Create new event listener")
	el := &EventListener{
		pods: make(chan *types.PodEvent),
		host: make(chan *types.HostEvent),
	}

	el.Listen()
	return el, nil
}
