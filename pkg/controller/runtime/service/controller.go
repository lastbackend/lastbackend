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

package service

import (
	"context"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"

	stgtypes "github.com/lastbackend/lastbackend/pkg/storage/etcd/types"
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime/cache"
	"fmt"
)

const logPrefix = "controller:service:controller"

type Controller struct {
	spec   chan *types.Service
	status chan *types.Service

	storage storage.Storage
	active  bool
}

// Watch services spec changes
func (sc *Controller) WatchSpec() {

	var (
		stg   = envs.Get().GetStorage()
		event = make(chan *stgtypes.WatcherEvent)
	)

	log.Debug("%s:watch_spec:> start watch service spec", logPrefix)
	go func() {
		for {
			select {
			case s := <-sc.spec:
				{
					if !sc.active {
						log.Debug("%s:watch_spec:> skip management course it is in slave mode", logPrefix)
						continue
					}

					if s == nil {
						log.Debug("%s:watch_spec:> skip because service is nil", logPrefix)
						continue
					}

					log.Debugf("%s:watch_spec:> service needs to be provisioned: %s:%s", logPrefix, s.Meta.Namespace, s.Meta.Name)
					if err := Provision(s); err != nil {
						log.Errorf("%s:watch_spec:> service provision: %s err: %v", logPrefix, s.Meta.Name, err)
						continue
					}
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case e := <-event:
				if e.Data == nil {
					continue
				}

				service := new(types.Service)

				if err := json.Unmarshal(e.Data.([]byte), &service); err != nil {
					log.Errorf("%s:watch_spec:> parse json err: %v", logPrefix, err)
					continue
				}

				sc.spec <- service
			}
		}
	}()

	stg.Watch(context.Background(), storage.ServiceKind, event)
}

// Watch services spec changes
func (sc *Controller) WatchStatus() {

	var (
		stg   = envs.Get().GetStorage()
		event = make(chan *stgtypes.WatcherEvent)
	)

	log.Debugf("%s:watch_status> start watch service status", logPrefix)
	go func() {
		for {
			select {
			case s := <-sc.status:
				{
					if !sc.active {
						log.Debugf("%s:watch_status> skip management course it is in slave mode", logPrefix)
						continue
					}

					if s == nil {
						log.Debugf("%s:watch_status> skip because service is nil", logPrefix)
						continue
					}

					log.Debugf("%s:watch_status> Service needs to be provisioned: %s", logPrefix, s.SelfLink())
					if err := HandleStatus(s); err != nil {
						log.Errorf("%s:watch_status> service provision: %s err: %v", logPrefix, s.SelfLink(), err)
					}
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case e := <-event:
				if e.Data == nil {
					continue
				}

				service := new(types.Service)

				if err := json.Unmarshal(e.Data.([]byte), &service); err != nil {
					log.Errorf("%s:watch_status:> parse json err: %v", logPrefix, err)
					continue
				}

				sc.status <- service
			}
		}
	}()

	stg.Watch(context.Background(), storage.ServiceKind, event)
}

// Pause service controller because not lead
func (sc *Controller) Pause() {
	sc.active = false
}

// Resume service controller management
func (sc *Controller) Resume() {

	var (
		stg = envs.Get().GetStorage()
	)

	sc.active = true

	log.Debugf("%s:resume> start check services states", logPrefix)

	nss := make(map[string]*types.Namespace)

	err := stg.Map(context.Background(), storage.NamespaceKind, "", &nss)
	if err != nil {
		log.Errorf("%s:resume> get namespaces list err: %v", logPrefix, err)
	}

	for _, ns := range nss {

		svcs := make(map[string]*types.Service)

		err := stg.Map(context.Background(), storage.ServiceKind, etcd.BuildServiceQuery(ns.Meta.Name), &svcs)
		if err != nil {
			log.Errorf("%s:resume> get services list err: %v", logPrefix, err)
		}

		for _, svc := range svcs {
			sc.spec <- svc
		}

		for _, svc := range svcs {
			log.Debugf("%s:resume> check service [%s] status", logPrefix, svc.SelfLink())
			sc.status <- svc
		}
	}
}

func (sc *Controller) Observe(ctx context.Context, cache *cache.Cache) {

	done := make(chan bool)
	event := make(chan *stgtypes.WatcherEvent)

	go func() {
		for {
			select {
			case <-ctx.Done():
				done <- true
				return
			case e := <-event:
				if e.Data == nil {
					continue
				}

				res := types.ServiceEvent{}
				res.Action = e.Action
				res.Name = e.Name

				service := new(types.Service)

				if err := json.Unmarshal(e.Data.([]byte), *service); err != nil {
					log.Errorf("%s:> parse data err: %v", logPrefix, err)
					continue
				}

				res.Data = service

				// TODO: service status handlers
			}
		}
	}()

	if err := sc.storage.Watch(ctx, storage.ServiceKind, event); err != nil {
		return
	}

	emit := cache.Deployments.Subscribe()

	for {
		select {
		case <-ctx.Done():
			done <- true
			return
		case e := <-emit:
			fmt.Println("change deployment", e)
		}
	}
}

// NewServiceController return new controller instance
func NewServiceController(_ context.Context) *Controller {
	sc := new(Controller)
	sc.active = false
	sc.spec = make(chan *types.Service, 0)
	sc.status = make(chan *types.Service, 0)
	return sc
}
