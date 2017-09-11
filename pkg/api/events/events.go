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
	"github.com/lastbackend/lastbackend/pkg/api/context"
	"github.com/lastbackend/lastbackend/pkg/api/app"
	"github.com/lastbackend/lastbackend/pkg/api/service/views/v1"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const (
	logLevel    = 1
	defaultRoom = "lastbackend"
)

type Events struct{}

func (e *Events) Listen() error {

	var (
		hub = context.Get().GetWssHub()
		ctx = context.Get().Background()

		ns      = app.New(ctx)
		service = make(chan *types.Service)
	)

	log.V(logLevel).Debug("Events: start events listener")

	go func() {
		for {
			select {
			case s := <-service:
				{
					if s == nil {
						continue
					}

					log.V(logLevel).Debugf("Events: service %s changed", s.Meta.Name)

					if obj, err := v1.New(s).ToJson(); err == nil {
						if room := hub.GetRoom(defaultRoom); room != nil {
							log.V(logLevel).Debug("Events: room founded, try to broadcast")
							room.Broadcast <- obj
						}

					}
					log.V(logLevel).Debug("Events: send update finished")
				}
			}
		}
	}()

	if err := ns.WatchService(service); err != nil {
		log.V(logLevel).Errorf("Events: watch services in app err: %s", err.Error())
		return err
	}

	return nil
}

func NewEventListener() *Events {
	listener := new(Events)
	return listener
}
