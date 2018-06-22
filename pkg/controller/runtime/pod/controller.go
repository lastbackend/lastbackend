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

package pod

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime/cache"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd"
	"fmt"
)

type Controller struct {
	podChan chan *types.Pod
	storage storage.Storage
	cache   *cache.Cache
	service *types.Service
	active  bool
}

//// Watch pod spec changes
//func (ctrl *Controller) WatchStatus() {
//
//	var (
//		stg   = envs.Get().GetStorage()
//		event = make(chan *stgtypes.WatcherEvent)
//	)
//
//	log.Debug("controller:pod:controller: start watch pod spec")
//	go func() {
//		for {
//			select {
//			case s := <-ctrl.status:
//				{
//					if !ctrl.active {
//						log.Debug("controller:pod:controller: skip management course it is in slave mode")
//						continue
//					}
//
//					if s == nil {
//						log.Debug("controller:pod:controller: skip because service is nil")
//						continue
//					}
//
//					log.Debugf("controller:pod:controller: Service needs to be provisioned: %s:%s", s.Meta.Namespace, s.Meta.Name)
//					if err := HandleStatus(s); err != nil {
//						log.Errorf("controller:pod:controller: service provision: %s err: %s", s.Meta.Name, err.Error())
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
//				pod := new(types.Pod)
//
//				if err := json.Unmarshal(e.Data.([]byte), &pod); err != nil {
//					log.Errorf("controller:pod:controller: parse json err: %v", err)
//					continue
//				}
//
//				ctrl.status <- pod
//			}
//		}
//	}()
//
//	stg.Watch(context.Background(), storage.PodKind, event)
//}

// Pause pod controller because not lead
func (ctrl *Controller) Pause() {
	ctrl.active = false
}

// Resume pod controller management
func (ctrl *Controller) Resume() {

	var (
		stg = envs.Get().GetStorage()
	)

	ctrl.active = true

	log.Debug("controller:pod:controller:resume start check pod states")

	nss := make(map[string]*types.Namespace, 0)

	err := stg.Map(context.Background(), storage.NamespaceKind, "", &nss)
	if err != nil {
		log.Errorf("controller:pod:controller:resume get namespaces list err: %s", err.Error())
	}

	for _, ns := range nss {

		pl := make(map[string]*types.Pod, 0)

		err := stg.Map(context.Background(), storage.PodKind, etcd.BuildPodQuery(ns.Meta.Name, "", ""), &pl)
		if err != nil {
			log.Errorf("controller:pod:controller:resume get pod list err: %s", err.Error())
		}

		for _, p := range pl {
			ctrl.podChan <- p
		}
	}
}

func (ctrl *Controller) Observe(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case p := <-ctrl.podChan:

				ctrl.cache.Pods.Set(p.Meta.SelfLink, p)

				// todo: run handlers
			}
		}
	}()

	event := make(chan cache.PodEvent)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-event:
				fmt.Println("pod event", e)
				// todo: run handlers
			}
		}
	}()

	ctrl.cache.Pods.Subscribe(event)
}

func (ctrl *Controller) UpdatePod(ctx context.Context, pod *types.Pod) {
	ctrl.podChan <- pod
}

// NewDeploymentController return new controller instance
func NewPodController(stg storage.Storage, c *cache.Cache, s *types.Service) *Controller {
	ctrl := new(Controller)
	ctrl.active = false
	ctrl.storage = stg
	ctrl.cache = c
	ctrl.service = s
	ctrl.podChan = make(chan *types.Pod)
	return ctrl
}
