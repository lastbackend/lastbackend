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

package runtime

import (
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/controller/state"
	"github.com/lastbackend/lastbackend/pkg/controller/state/service"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"golang.org/x/net/context"
)

const (
	logPrefix = "controller:runtime:observer"
)

type Observer struct {
	stg   storage.Storage
	state *state.State
}

func NewObserver(ctx context.Context) *Observer {

	o := new(Observer)
	o.stg = envs.Get().GetStorage()

	o.state = state.NewState()

	go o.watchServices(ctx)
	go o.watchNodes(ctx)
	go o.watchPods(ctx)

	o.state.Restore()

	return o
}

func (o *Observer) watchServices(ctx context.Context) {

	var (
		svc = make(chan types.ServiceEvent)
	)

	sm := distribution.NewServiceModel(ctx, o.stg)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-svc:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					_, ok := o.state.Service[w.Data.SelfLink()]
					if ok {
						delete(o.state.Service, w.Data.SelfLink())
					}
					continue
				}

				_, ok := o.state.Service[w.Data.SelfLink()]
				if !ok {
					o.state.Service[w.Data.SelfLink()] = service.NewServiceState(w.Data)
				}

				o.state.Service[w.Data.SelfLink()].SetService(w.Data)
			}
		}
	}()

	sm.Watch(svc)
}

func (o *Observer) watchPods(ctx context.Context) {

	// Watch pods change
	var (
		p = make(chan types.PodEvent)
	)

	pm := distribution.NewPodModel(ctx, o.stg)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-p:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					_, ok := o.state.Service[w.Data.ServiceLink()]
					if ok {
						o.state.Service[w.Data.ServiceLink()].DelPod(w.Data)
					}
					continue
				}

				log.Info("send pod to cluster state")
				o.state.Cluster.SetPod(w.Data)
				log.Info("pod updated in cluster state")

				_, ok := o.state.Service[w.Data.ServiceLink()]
				if !ok {
					log.Info("service state not found: skip")
					break
				}

				log.Info("send pod to service state")
				o.state.Service[w.Data.ServiceLink()].SetPod(w.Data)
				log.Info("pod updated in service state")
			}
		}
	}()

	pm.Watch(p)
}

func (o *Observer) watchNodes(ctx context.Context) {
	var (
		p = make(chan types.NodeEvent)
	)

	nm := distribution.NewNodeModel(ctx, o.stg)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-p:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					o.state.Cluster.DelNode(w.Data)
					continue
				}

				o.state.Cluster.SetNode(w.Data)
			}
		}
	}()

	nm.Watch(p)
}
