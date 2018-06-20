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

package deployment

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd"

	stgtypes "github.com/lastbackend/lastbackend/pkg/storage/etcd/types"
	"encoding/json"
)

const logPrefix = "controller:deployment"

type Controller struct {
	status chan *types.Deployment
	spec   chan *types.Deployment
	active bool
}

// Watch deployment spec changes
func (dc *Controller) WatchSpec() {

	var (
		stg   = envs.Get().GetStorage()
		event = make(chan *stgtypes.WatcherEvent)
	)

	log.Debug("controller:deployment:controller: start watch deployment spec")
	go func() {
		for {
			select {
			case s := <-dc.spec:
				{
					if !dc.active {
						log.Debug("controller:deployment:controller: skip management couse it is in slave mode")
						continue
					}

					if s == nil {
						log.Debug("controller:deployment:controller: skip because service is nil")
						continue
					}

					log.Debugf("controller:deployment:controller: Service needs to be provisioned: %s:%s", s.Meta.Namespace, s.Meta.Name)
					if err := Provision(s); err != nil {
						log.Errorf("controller:deployment:controller: service provision: %s err: %s", s.Meta.Name, err.Error())
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

				deployment := new(types.Deployment)

				if err := json.Unmarshal(e.Data.([]byte), &deployment); err != nil {
					log.Errorf("controller:deployment:controller: parse json err: %s", err.Error())
					continue
				}

				dc.spec <- deployment
			}
		}
	}()

	stg.Watch(context.Background(), storage.DeploymentKind, event)
}

// Watch deployment spec changes
func (dc *Controller) WatchStatus() {

	var (
		stg   = envs.Get().GetStorage()
		event = make(chan *stgtypes.WatcherEvent)
	)

	log.Debugf("%s:status> start watch deployment status", logPrefix)
	go func() {
		for {
			select {
			case s := <-dc.status:
				{
					if !dc.active {
						log.Debug("%s:status> skip management couse it is in slave mode", logPrefix)
						continue
					}

					if s == nil {
						log.Debug("%s:status> skip because service is nil", logPrefix)
						continue
					}

					log.Debugf("%s:status> Service needs to be provisioned: %s", logPrefix, s.SelfLink())
					if err := HandleStatus(s); err != nil {
						log.Errorf("%s:status> service provision: %s err: %s", logPrefix, s.SelfLink(), err.Error())
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

				deployment := new(types.Deployment)

				if err := json.Unmarshal(e.Data.([]byte), &deployment); err != nil {
					log.Errorf("%s:status parse json err: %s", logPrefix, err.Error())
					continue
				}

				dc.status <- deployment
			}
		}
	}()

	stg.Watch(context.Background(), storage.DeploymentKind, event)
}

// Pause deployment controller because not lead
func (dc *Controller) Pause() {
	dc.active = false
}

// Resume deployment controller management
func (dc *Controller) Resume() {

	var (
		stg = envs.Get().GetStorage()
	)

	dc.active = true

	nss := make(map[string]*types.Namespace)

	log.Debug("controller:deployment:controller:resume start check deployment states")
	err := stg.Map(context.Background(), storage.NamespaceKind, "", &nss)
	if err != nil {
		log.Errorf("controller:deployment:controller:resume get namespaces list err: %s", err.Error())
	}

	for _, ns := range nss {

		dl := make(map[string]*types.Deployment)

		err := stg.Map(context.Background(), storage.DeploymentKind, etcd.BuildDeploymentQuery(ns.Meta.Name, ""), ns.Meta.Name)
		if err != nil {
			log.Errorf("controller:deployment:controller:resume get deployment list err: %s", err.Error())
		}

		for _, d := range dl {

			dp := new(types.Deployment)

			err := stg.Get(context.Background(), storage.DeploymentKind, dp.Meta.SelfLink, &dp)
			if err != nil {
				log.Errorf("controller:deployment:controller:resume get deployment err: %s", err.Error())
			}

			dc.spec <- d
			dc.status <- d
		}
	}
}

// NewDeploymentController return new controller instance
func NewDeploymentController(_ context.Context) *Controller {
	sc := new(Controller)
	sc.active = false
	sc.status = make(chan *types.Deployment)
	sc.spec = make(chan *types.Deployment)
	return sc
}
