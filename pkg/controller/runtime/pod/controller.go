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
	"encoding/json"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd"
	"github.com/lastbackend/lastbackend/pkg/controller/runtime/cache"

	stgtypes "github.com/lastbackend/lastbackend/pkg/storage/types"
)

type Controller struct {
	status chan *types.Pod
	active bool
}

// Watch pod spec changes
func (pc *Controller) WatchStatus() {

	var (
		stg   = envs.Get().GetStorage()
		event = make(chan *stgtypes.WatcherEvent)
	)

	log.Debug("controller:pod:controller: start watch pod spec")
	go func() {
		for {
			select {
			case s := <-pc.status:
				{
					if !pc.active {
						log.Debug("controller:pod:controller: skip management course it is in slave mode")
						continue
					}

					if s == nil {
						log.Debug("controller:pod:controller: skip because service is nil")
						continue
					}

					log.Debugf("controller:pod:controller: Service needs to be provisioned: %s:%s", s.Meta.Namespace, s.Meta.Name)
					if err := HandleStatus(s); err != nil {
						log.Errorf("controller:pod:controller: service provision: %s err: %s", s.Meta.Name, err.Error())
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

				pod := new(types.Pod)

				if err := json.Unmarshal(e.Data.([]byte), &pod); err != nil {
					log.Errorf("controller:pod:controller: parse json err: %v", err)
					continue
				}

				pc.status <- pod
			}
		}
	}()

	stg.Watch(context.Background(), storage.PodKind, event)
}

// Pause pod controller because not lead
func (pc *Controller) Pause() {
	pc.active = false
}

// Resume pod controller management
func (pc *Controller) Resume() {

	var (
		stg = envs.Get().GetStorage()
	)

	pc.active = true

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
			pc.status <- p
		}
	}
}

func (pc *Controller) Observe(ctx context.Context, cache *cache.Cache) {
	// TODO: watch etcd: pod collection
}

// NewDeploymentController return new controller instance
func NewPodController(_ context.Context) *Controller {
	sc := new(Controller)
	sc.active = false
	sc.status = make(chan *types.Pod)
	return sc
}
