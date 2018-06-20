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
)

const logPrefix = "controller:service"

type Controller struct {
	spec   chan *types.Service
	status chan *types.Service
	active bool
}

// Watch services spec changes
func (sc *Controller) WatchSpec() {

	var (
		stg   = envs.Get().GetStorage()
		event = make(chan *stgtypes.WatcherEvent)
	)

	log.Debug("controller:service:controller: start watch service spec")
	go func() {
		for {
			select {
			case s := <-sc.spec:
				{
					if !sc.active {
						log.Debug("controller:service:controller: skip management course it is in slave mode")
						continue
					}

					if s == nil {
						log.Debug("controller:service:controller: skip because service is nil")
						continue
					}

					log.Debugf("controller:service:controller: Service needs to be provisioned: %s:%s", s.Meta.Namespace, s.Meta.Name)
					if err := Provision(s); err != nil {
						log.Errorf("controller:service:controller: service provision: %s err: %s", s.Meta.Name, err.Error())
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
					log.Errorf("controller:service:controller: parse json err: %s", err.Error())
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
						log.Errorf("%s:watch_status> service provision: %s err: %s", logPrefix, s.SelfLink(), err.Error())
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
					log.Errorf("controller:service:controller: parse json err: %s", err.Error())
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
		log.Errorf("%s:resume> get namespaces list err: %s", logPrefix, err.Error())
	}

	for _, ns := range nss {

		svcs := make(map[string]*types.Service)

		err := stg.Map(context.Background(), storage.ServiceKind, etcd.BuildServiceQuery(ns.Meta.Name), &svcs)
		if err != nil {
			log.Errorf("%s:resume> get services list err: %s", logPrefix, err.Error())
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

// NewServiceController return new controller instance
func NewServiceController(_ context.Context) *Controller {
	sc := new(Controller)
	sc.active = false
	sc.spec = make(chan *types.Service)
	sc.status = make(chan *types.Service)
	return sc
}
