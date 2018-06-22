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
	"github.com/lastbackend/lastbackend/pkg/controller/runtime/cache"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"fmt"
)

const logPrefix = "controller:service:controller"

type Controller struct {
	serviceChan chan *types.Service

	cache   *cache.Cache
	storage storage.Storage

	service *types.Service

	active bool
}

//// Watch services spec changes
//func (ctrl *Controller) WatchSpec() {
//
//	var (
//		stg   = envs.Get().GetStorage()
//		event = make(chan *stgtypes.WatcherEvent)
//	)
//
//	log.Debug("%s:watch_spec:> start watch service spec", logPrefix)
//	go func() {
//		for {
//			select {
//			case s := <-ctrl.spec:
//				{
//					if !ctrl.active {
//						log.Debug("%s:watch_spec:> skip management course it is in slave mode", logPrefix)
//						continue
//					}
//
//					if s == nil {
//						log.Debug("%s:watch_spec:> skip because service is nil", logPrefix)
//						continue
//					}
//
//					log.Debugf("%s:watch_spec:> service needs to be provisioned: %s:%s", logPrefix, s.Meta.Namespace, s.Meta.Name)
//					if err := Provision(s); err != nil {
//						log.Errorf("%s:watch_spec:> service provision: %s err: %v", logPrefix, s.Meta.Name, err)
//						continue
//					}
//				}
//			}
//		}
//	}()
//
//	go func() {
//		for {
//			select {
//			case e := <-event:
//				if e.Data == nil {
//					continue
//				}
//
//				service := new(types.Service)
//
//				if err := json.Unmarshal(e.Data.([]byte), &service); err != nil {
//					log.Errorf("%s:watch_spec:> parse json err: %v", logPrefix, err)
//					continue
//				}
//
//				ctrl.spec <- service
//			}
//		}
//	}()
//
//	stg.Watch(context.Background(), storage.ServiceKind, event)
//}
//
//// Watch services spec changes
//func (ctrl *Controller) WatchStatus() {
//
//	var (
//		stg   = envs.Get().GetStorage()
//		event = make(chan *stgtypes.WatcherEvent)
//	)
//
//	log.Debugf("%s:watch_status> start watch service status", logPrefix)
//	go func() {
//		for {
//			select {
//			case s := <-ctrl.status:
//				{
//					if !ctrl.active {
//						log.Debugf("%s:watch_status> skip management course it is in slave mode", logPrefix)
//						continue
//					}
//
//					if s == nil {
//						log.Debugf("%s:watch_status> skip because service is nil", logPrefix)
//						continue
//					}
//
//					log.Debugf("%s:watch_status> Service needs to be provisioned: %s", logPrefix, s.SelfLink())
//					if err := HandleStatus(s); err != nil {
//						log.Errorf("%s:watch_status> service provision: %s err: %v", logPrefix, s.SelfLink(), err)
//					}
//				}
//			}
//		}
//	}()
//
//	go func() {
//		for {
//			select {
//			case e := <-event:
//				if e.Data == nil {
//					continue
//				}
//
//				service := new(types.Service)
//
//				if err := json.Unmarshal(e.Data.([]byte), &service); err != nil {
//					log.Errorf("%s:watch_status:> parse json err: %v", logPrefix, err)
//					continue
//				}
//
//				ctrl.status <- service
//			}
//		}
//	}()
//
//	stg.Watch(context.Background(), storage.ServiceKind, event)
//}

// Pause service controller because not lead
func (ctrl *Controller) Pause() {
	ctrl.active = false
}

// Resume service controller management
func (ctrl *Controller) Resume() {

	ctrl.active = true

	log.Debugf("%s:resume> start check services states", logPrefix)

	err := ctrl.storage.Get(context.Background(), storage.ServiceKind, ctrl.service.Meta.SelfLink, &ctrl.service)
	if err != nil {
		log.Errorf("%s:resume> get services list err: %v", logPrefix, err)
	}

}

func (ctrl *Controller) Observe(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case s := <-ctrl.serviceChan:

				ctrl.cache.Services.Set(s.Meta.SelfLink, s)

				// todo: run handlers
			}
		}
	}()

	event := make(chan cache.ServiceEvent)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-event:
				fmt.Println("service event", e)
				// todo: run handlers
			}
		}
	}()

	ctrl.cache.Services.Subscribe(event)
}

func (ctrl *Controller) UpdateService(ctx context.Context, service *types.Service) {
	ctrl.serviceChan <- service
}

// NewServiceController return new controller instance
func NewServiceController(stg storage.Storage, c *cache.Cache, service *types.Service) *Controller {
	ctrl := new(Controller)
	ctrl.active = false
	ctrl.storage = stg
	ctrl.service = service
	ctrl.cache = c
	ctrl.serviceChan = make(chan *types.Service, 0)
	return ctrl
}
