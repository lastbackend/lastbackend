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

package observe

import (
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime/cache"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime/deployment"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime/pod"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime/service"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"golang.org/x/net/context"
	"sync"
)

const (
	logPrefix = "controller:runtime:observe"
)

type Controllers struct {
	pc *pod.Controller
	dc *deployment.Controller
	sc *service.Controller
}

type Observer struct {
	stg storage.Storage

	observers map[string]*ServiceObserver
}

func New(ctx context.Context, stg storage.Storage) *Observer {
	o := new(Observer)
	o.stg = stg
	o.observers = make(map[string]*ServiceObserver, 0)

	go o.watchServices(ctx)
	go o.watchDeployments(ctx)
	go o.watchPods(ctx)

	sl := make(map[string]*types.Service, 0)

	if err := o.stg.Map(ctx, storage.ServiceKind, types.EmptyString, &sl); err != nil {
		log.Errorf("$s:> ger services list err: %v", logPrefix, err)
		panic(err)
	}

	for _, s := range sl {
		if _, ok := o.observers[s.Meta.SelfLink]; !ok {
			o.newServiceObserver(ctx, s)
		}
	}

	return o
}

func (o Observer) newServiceObserver(ctx context.Context, service *types.Service) {
	var mutex = &sync.Mutex{}

	mutex.Lock()
	defer mutex.Unlock()

	so := NewServiceObserver(ctx, service)
	so.Run()
	o.observers[service.Meta.SelfLink] = so
}

func (o Observer) watchServices(ctx context.Context) {

	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-watcher:

				if w.Data == nil {
					continue
				}

				srv := w.Data.(*types.Service)

				if observer, ok := o.observers[srv.Meta.SelfLink]; !ok {
					o.newServiceObserver(ctx, srv)
				} else {
					observer.ctrl.sc.UpdateService(ctx, srv)
				}

			}
		}
	}()

	go envs.Get().GetStorage().Watch(ctx, storage.ServiceKind, watcher)
}

func (o Observer) watchDeployments(ctx context.Context) {
	// Watch deployments change

	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-watcher:

				if w.Data == nil {
					continue
				}

				dp := w.Data.(*types.Deployment)

				s := types.Service{}.CreateSelfLink(dp.Meta.Namespace, dp.Meta.Service)

				if observer, ok := o.observers[s]; !ok {

					srv := new(types.Service)
					if err := envs.Get().GetStorage().Get(ctx, storage.ServiceKind, s, &srv); err != nil {
						log.Errorf("%s:> get service err: %v", err)
						continue
					}

					o.newServiceObserver(ctx, srv)
				} else {
					observer.ctrl.dc.UpdateDeployment(ctx, dp)
				}

			}
		}
	}()

	go envs.Get().GetStorage().Watch(ctx, storage.ServiceKind, watcher)
}

func (o Observer) watchPods(ctx context.Context) {
	// Watch pods change

	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-watcher:

				if w.Data == nil {
					continue
				}

				pd := w.Data.(*types.Pod)

				s := types.Service{}.CreateSelfLink(pd.Meta.Namespace, pd.Meta.Service)

				if observer, ok := o.observers[s]; !ok {

					srv := new(types.Service)
					if err := envs.Get().GetStorage().Get(ctx, storage.ServiceKind, s, &srv); err != nil {
						log.Errorf("%s:> get service err: %v", err)
						continue
					}

					o.newServiceObserver(ctx, srv)
				} else {
					observer.ctrl.pc.UpdatePod(ctx, pd)
				}

			}
		}
	}()

	go envs.Get().GetStorage().Watch(ctx, storage.ServiceKind, watcher)
}

type ServiceObserver struct {
	ctx  context.Context
	ctrl Controllers
}

func NewServiceObserver(ctx context.Context, s *types.Service) *ServiceObserver {
	o := new(ServiceObserver)

	o.ctx = ctx
	c := cache.New()

	o.ctrl.sc = service.NewServiceController(envs.Get().GetStorage(), c, s)
	o.ctrl.dc = deployment.NewDeploymentController(envs.Get().GetStorage(), c, s)
	o.ctrl.pc = pod.NewPodController(envs.Get().GetStorage(), c, s)

	return o
}

func (o ServiceObserver) Run() {
	go o.ctrl.pc.Observe(o.ctx)
	go o.ctrl.dc.Observe(o.ctx)
	go o.ctrl.sc.Observe(o.ctx)
}

func (so ServiceObserver) Stop() {
	// todo: close all observers
}
