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

type DeploymentController struct {
	deployment chan *types.Deployment
	active     bool
}

func (dc *DeploymentController) Watch() {
	var (
		stg = envs.Get().GetStorage()
	)

	log.Debug("DeploymentController: start watch")
	go func() {
		for {
			select {
			case s := <-dc.deployment:
				{
					if !dc.active {
						log.Debug("DeploymentController: skip management course it is in slave mode")
						continue
					}

					if s == nil {
						log.Debug("DeploymentController: skip because service is nil")
						continue
					}

					log.Debugf("Service needs to be provisioned: %s:%s", s.Meta.Namespace, s.Meta.Name)
					if err := Provision(s); err != nil {
						log.Errorf("Error: DeploymentController: Service provision: %s", err.Error())
					}
				}
			}
		}
	}()

	stg.Deployment().SpecWatch(context.Background(), dc.deployment)
}

func (dc *DeploymentController) Pause() {
	dc.active = false
}

func (dc *DeploymentController) Resume() {

	var (
		stg = envs.Get().GetStorage()
	)

	dc.active = true

	log.Debug("Service: start check services states")
	nss, err := stg.Namespace().List(context.Background())
	if err != nil {
		log.Errorf("Service: Get namespaces list err: %s", err.Error())
	}

	for _, ns := range nss {
		dl, err := stg.Deployment().ListByNamespace(context.Background(), ns.Meta.Name)
		if err != nil {
			log.Errorf("Service: Get services list err: %s", err.Error())
		}

		for _, d := range dl {
			d, err := stg.Deployment().Get(context.Background(), d.Meta.Namespace, d.Meta.Name)
			if err != nil {
				log.Errorf("Service: Get service err: %s", err.Error())
			}
			dc.deployment <- d
		}
	}
}

func NewDeploymentController(_ context.Context) *DeploymentController {
	sc := new(DeploymentController)
	sc.active = false
	sc.deployment = make(chan *types.Deployment)
	return sc
}
