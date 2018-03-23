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
)

type Controller struct {
	status chan *types.Deployment
	spec chan *types.Deployment
	active     bool
}

// Watch deployment spec changes
func (dc *Controller) WatchSpec() {

	var (
		stg = envs.Get().GetStorage()
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

	stg.Deployment().WatchSpec(context.Background(), dc.spec)
}

// Watch deployment spec changes
func (dc *Controller) WatchStatus() {

	var (
		stg = envs.Get().GetStorage()
	)

	log.Debug("controller:deployment:controller: start watch deployment status")
	go func() {
		for {
			select {
			case s := <-dc.status:
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
					if err := HandleStatus(s); err != nil {
						log.Errorf("controller:deployment:controller: service provision: %s err: %s", s.Meta.Name, err.Error())
					}
				}
			}
		}
	}()

	stg.Deployment().WatchStatus(context.Background(), dc.status)
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

	log.Debug("controller:deployment:controller:resume start check deployment states")
	nss, err := stg.Namespace().List(context.Background())
	if err != nil {
		log.Errorf("controller:deployment:controller:resume get namespaces list err: %s", err.Error())
	}

	for _, ns := range nss {
		dl, err := stg.Deployment().ListByNamespace(context.Background(), ns.Meta.Name)
		if err != nil {
			log.Errorf("controller:deployment:controller:resume get deployment list err: %s", err.Error())
		}

		for _, d := range dl {
			d, err := stg.Deployment().Get(context.Background(), d.Meta.Namespace, d.Meta.Service, d.Meta.Name)
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
	sc.spec   = make(chan *types.Deployment)
	return sc
}
