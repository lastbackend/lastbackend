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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/discovery/context"
	"strings"
)

type ServiceController struct {
	context  *context.Context
	services chan *types.Service

	active bool
}

func (sc *ServiceController) Watch(services chan *types.Service) {
	var (
		log = sc.context.GetLogger()
		stg = sc.context.GetStorage()
	)

	log.V(logLevel).Debug("ServiceController: start watch")

	go func() {
		for {
			select {
			case s := <-sc.services:
				{

					if !sc.active {
						log.V(logLevel).Debug("ServiceController: skip management cause it is in slave mode")
						continue
					}

					if s == nil {
						continue
					}

					endpoint := fmt.Sprintf("%s-%s.%s", s.Meta.Name, s.Meta.Namespace, *context.Get().GetConfig().SystemDomain)
					endpoint = strings.Replace(endpoint, ":", "-", -1)

					if s.State.State == types.StateDestroyed {
						if err := stg.Endpoint().Remove(context.Get().Background(), endpoint); err != nil {
							log.V(logLevel).Errorf("ServiceController: remove service endpoint error %s", err.Error())
						}
						continue
					}

					hosts := make(map[string]string)
					ips := []string{}
					for _, pod := range s.Pods {
						if _, ok := hosts[pod.Node.ID]; ok || pod.Spec.State == types.StateDestroyed {
							continue
						}

						node, err := stg.Node().Get(context.Get().Background(), pod.Node.ID)
						if err != nil {
							log.V(logLevel).Errorf("ServiceController: get node error %s", err.Error())
							break
						}

						if node == nil {
							log.V(logLevel).Errorf("ServiceController: node not found")
							break
						}

						hosts[pod.Node.ID] = node.Meta.IP
						ips = append(ips, node.Meta.IP)
					}

					if err := stg.Endpoint().Upsert(context.Get().Background(), endpoint, ips); err != nil {
						log.V(logLevel).Errorf("ServiceController: upsert service endpoint error %s", err.Error())
						continue
					}
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case s := <-services:
				{
					sc.services <- s
				}
			}
		}
	}()

	stg.Service().Watch(sc.context.Background(), sc.services)
}

func (sc *ServiceController) Pause() {
	sc.context.GetLogger().Debugf("ServiceController: pause")
	sc.active = false
}

func (sc *ServiceController) Resume() {
	sc.context.GetLogger().Debugf("ServiceController: pause")
	sc.active = true
}

func NewServiceController(ctx *context.Context) *ServiceController {
	sc := new(ServiceController)
	sc.context = ctx
	sc.active = false
	sc.services = make(chan *types.Service)
	return sc
}
