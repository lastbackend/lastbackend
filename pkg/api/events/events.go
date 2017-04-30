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

package events

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/api/context"
	"github.com/lastbackend/lastbackend/pkg/api/namespace"
	"github.com/lastbackend/lastbackend/pkg/api/service/views/v1"
)

type Events struct {
}

func (e *Events) Listen() {
	var (
		log = context.Get().GetLogger()
		hub = context.Get().GetWssHub()
	)

	ns := namespace.New(context.Get().Background())
	service := make(chan *types.Service)

	go func() {
		for {
			select {
			case s := <-service:
				{
					if s == nil {
						continue
					}

					log.Debugf("%s changed", s.Meta.Name)
					if obj, err := v1.New(s).ToJson(); err == nil {
						if room := hub.GetRoom("lastbackend"); room != nil {
							log.Debug("Room founded, try to broadcast")
							fmt.Println(string(obj))
							room.Broadcast <- obj
						}

					}
					log.Debug("Send update finished")
				}
			}
		}
	}()

	ns.Watch(service)
}

func NewEventListener() *Events {
	listener := new(Events)
	return listener
}
