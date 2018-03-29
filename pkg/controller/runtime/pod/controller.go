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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

type Controller struct {
	status chan *types.Pod
	active bool
}

// Watch pod spec changes
func (pc *Controller) WatchStatus() {

	var (
		stg = envs.Get().GetStorage()
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

	stg.Pod().WatchStatus(context.Background(), pc.status)
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
	nss, err := stg.Namespace().List(context.Background())
	if err != nil {
		log.Errorf("controller:pod:controller:resume get namespaces list err: %s", err.Error())
	}

	for _, ns := range nss {
		pl, err := stg.Pod().ListByNamespace(context.Background(), ns.Meta.Name)
		if err != nil {
			log.Errorf("controller:pod:controller:resume get pod list err: %s", err.Error())
		}

		for _, p := range pl {
			pc.status <- p
		}
	}
}

// NewDeploymentController return new controller instance
func NewPodController(_ context.Context) *Controller {
	sc := new(Controller)
	sc.active = false
	sc.status = make(chan *types.Pod)
	return sc
}
