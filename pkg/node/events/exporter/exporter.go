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
	"context"
	"sync"
	"time"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
)

type Exporter struct {
	lock   sync.RWMutex
	dispatcher Dispatcher

	resources types.NodeStatus
	pods   map[string]*types.PodStatus
}

type Dispatcher func(options *request.NodeStatusOptions) error

func (e *Exporter) Loop() {

	go func(ctx context.Context) {
		for {
			select {}
		}
	}(context.Background())

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second * 3)

		for range ticker.C {
			opts := new(request.NodeStatusOptions)
			opts.Pods = make(map[string]*request.NodePodStatusOptions)
			opts.Resources.Capacity = e.resources.Capacity
			opts.Resources.Allocated = e.resources.Allocated

			e.lock.Lock()
			var  i = 0
			for p, status := range e.pods {
				i++
				if i > 10 {
					break
				}
				opts.Pods[p] = getPodOptions(status)
			}

			for p := range opts.Pods {
				delete(e.pods, p)
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

func (e *Exporter) Resources(res types.NodeStatus) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.resources.Allocated = res.Allocated
	e.resources.Capacity  = res.Capacity
}

func (e *Exporter) PodStatus(pod string, status *types.PodStatus) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.pods[pod] = status
}

func (e *Exporter) SetDispatcher(f Dispatcher) {
	e.dispatcher = f
}

func NewExporter() *Exporter {
	d := new(Exporter)

	d.pods = make(map[string]*types.PodStatus)

	return d
}


func getPodOptions(p *types.PodStatus) *request.NodePodStatusOptions {
	opts := v1.Request().Node().NodePodStatusOptions()
	opts.State = p.Stage
	opts.Message = p.Message
	opts.Containers = p.Containers
	opts.Network = p.Network
	opts.Steps = p.Steps
	return opts
}