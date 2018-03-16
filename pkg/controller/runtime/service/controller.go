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
)

type ServiceController struct {
	services chan *types.Service
	active   bool
}

func (sc *ServiceController) Watch() {
	var (
		stg = envs.Get().GetStorage()
	)

	log.Debug("ServiceController: start watch")
	go func() {
		for {
			select {
			case s := <-sc.services:
				{
					if !sc.active {
						log.Debug("ServiceController: skip management course it is in slave mode")
						continue
					}

					if s == nil {
						log.Debug("ServiceController: skip because service is nil")
						continue
					}

					log.Debugf("Service needs to be provisioned: %s:%s", s.Meta.Namespace, s.Meta.Name)
					if err := Provision(s); err != nil {
						log.Errorf("Error: ServiceController: Service provision: %s", err.Error())
					}
				}
			}
		}
	}()

	stg.Service().SpecWatch(context.Background(), sc.services)
}

func (sc *ServiceController) Pause() {
	sc.active = false
}

func (sc *ServiceController) Resume() {

	var (
		stg = envs.Get().GetStorage()
	)

	sc.active = true

	log.Debug("Service: start check services states")
	nss, err := stg.Namespace().List(context.Background())
	if err != nil {
		log.Errorf("Service: Get namespaces list err: %s", err.Error())
	}

	for _, ns := range nss {
		svcs, err := stg.Service().ListByNamespace(context.Background(), ns.Meta.Name)
		if err != nil {
			log.Errorf("Service: Get services list err: %s", err.Error())
		}

		for _, svc := range svcs {
			svc, err := stg.Service().Get(context.Background(), svc.Meta.Namespace, svc.Meta.Name)
			if err != nil {
				log.Errorf("Service: Get service err: %s", err.Error())
			}
			sc.services <- svc
		}
	}
}

func NewServiceController(_ context.Context) *ServiceController {
	sc := new(ServiceController)
	sc.active = false
	sc.services = make(chan *types.Service)
	return sc
}
