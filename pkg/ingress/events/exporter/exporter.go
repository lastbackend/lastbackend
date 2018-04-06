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

package exporter

import (
	"time"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"sync"
	"context"
	"github.com/lastbackend/lastbackend/pkg/log"
)

type Exporter struct {
	lock   sync.RWMutex
	dispatcher Dispatcher

	routes   map[string]*types.RouteStatus
}

type Dispatcher func(options *request.IngressStatusOptions) error

func (e *Exporter) Loop() {

	go func(ctx context.Context) {
		for {
			select {}
		}
	}(context.Background())

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second * 10)

		for range ticker.C {
			opts := new(request.IngressStatusOptions)
			opts.Ready = true
			opts.Routes = make(map[string]*request.IngressRouteStatusOptions)

			e.lock.Lock()
			var  i = 0
			for r, status := range e.routes {
				i++
				if i > 10 {
					break
				}
				opts.Routes[r] = getRouteOptions(status)
			}

			for r := range opts.Routes {
				delete(e.routes, r)
			}

			e.lock.Unlock()
			if e.dispatcher == nil {
				continue
			}

			if err := e.dispatcher(opts); err != nil {
				log.Errorf("node:exporter:dispatch err: %", err.Error())
			}
		}

	}(context.Background())

}

func (e *Exporter) PodStatus(routes string, status *types.RouteStatus) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.routes[routes] = status
}

func (e *Exporter) SetDispatcher(f Dispatcher) {
	e.dispatcher = f
}

func NewExporter() *Exporter {
	d := new(Exporter)

	d.routes = make(map[string]*types.RouteStatus)

	return d
}


func getRouteOptions(r *types.RouteStatus) *request.IngressRouteStatusOptions {
	opts := v1.Request().Ingress().IngressRouteStatusOptions()
	opts.State = r.State
	opts.Message = r.Message
	return opts
}