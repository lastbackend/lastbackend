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

package listener

import (
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"github.com/lastbackend/lastbackend/pkg/api/node/views/v1"
	"github.com/lastbackend/lastbackend/pkg/util/http"
)

type EventListener struct {
	events chan *types.Event
	spec   chan *types.NodeSpec
	http   *http.RawReq
}

func (el *EventListener) Loop() {
	for {
		select {
		case e := <-el.events:
			spec, err := el.request(e)
			if err != nil {
				// TODO: try to send after a small timeout
			}
			if spec == nil {
				continue
			}

			el.spec <- spec
		}
	}
}

func (el *EventListener) request (event *types.Event) (*types.NodeSpec, error) {
	var (
		er       = new(errors.Http)
		endpoint = "/node/event"
		spec     = v1.Spec{}
	)

	_, _, err := el.http.
		PUT(endpoint).
		AddHeader("Content-Type", "application/json").
		BodyJSON(event).
		Request(&spec, er)
	if err != nil {
		return nil, err
	}

	if er.Code == 401 {
		return nil, nil
	}

	if er.Code != 0 {
		return nil, errors.New(er.Message)
	}

	return v1.FromNodeSpec(spec), nil
}

func (el *EventListener) Send(event *types.Event)  {
	el.events <- event
}

func New(http *http.RawReq, spec chan *types.NodeSpec) *EventListener {
	el := new(EventListener)
	el.http   = http
	el.spec   = spec
	el.events = make(chan *types.Event)

	go el.Loop()

	return el
}