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
	"github.com/lastbackend/lastbackend/pkg/controller/context"
	"github.com/lastbackend/lastbackend/pkg/common/types"
)

type ServiceController struct {
	context context.Context
	active bool
}


func (sc *ServiceController) Watch () {
	var (
		log = sc.context.GetLogger()
		stg = sc.context.GetStorage()
		svc = make (chan *types.Service)
	)

	log.Debug("Controller:ServiceController: start watch")
	go func(){
		for {
			select {
			case s := <- svc : {
				if !sc.active {
					continue
				}


			}
			}
		}
	}()
	stg.Service().SpecWatch(sc.context.Background(), svc)
}

func (sc *ServiceController) Pause () {

}

func (sc *ServiceController) Resume () {

}


func NewServiceController (ctx context.Context) *ServiceController {
	sc := new(ServiceController)
	sc.active = false

	return sc
}
