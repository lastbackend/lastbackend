//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/controller/context"
)

type ServiceController struct {
	context  *context.Context
	services chan *types.Service
	active   bool
}

func (sc *ServiceController) Watch() {
	var (
		log = sc.context.GetLogger()
		stg = sc.context.GetStorage()
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

	stg.Service().SpecWatch(sc.context.Background(), sc.services)
}

func (sc *ServiceController) Pause() {
	sc.active = false
}

func (sc *ServiceController) Resume() {

	var (
		log = sc.context.GetLogger()
		stg = sc.context.GetStorage()
	)

	sc.active = true

	log.Debug("Service: start check services states")
	nss, err := stg.Namespace().List(sc.context.Background())
	if err != nil {
		log.Errorf("Service: Get namespaces list err: %s", err.Error())
	}

	for _, ns := range nss {
		svcs, err := stg.Service().ListByNamespace(sc.context.Background(), ns.Meta.Name)
		if err != nil {
			log.Errorf("Service: Get services list err: %s", err.Error())
		}

		for _, svc := range svcs {
			svc, err := stg.Service().GetByName(sc.context.Background(), svc.Meta.Namespace, svc.Meta.Name)
			if err != nil {
				log.Errorf("Service: Get service err: %s", err.Error())
			}
			sc.services <- svc
		}
	}
}

func NewServiceController(ctx *context.Context) *ServiceController {
	sc := new(ServiceController)
	sc.context = ctx
	sc.active = false
	sc.services = make(chan *types.Service)
	return sc
}
