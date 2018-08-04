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

package state

import (
	"context"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/controller/state/cluster"
	"github.com/lastbackend/lastbackend/pkg/controller/state/service"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/log"
	)

type State struct {
	Cluster *cluster.ClusterState
	Service map[string]*service.ServiceState
}

func (s *State) Restore() {

	println()
	println()
	log.Info("start cluster restore")
	s.Cluster.Restore()
	log.Info("finish cluster restore\n\n")


	log.Info("start services restore")
	nm := distribution.NewNamespaceModel(context.Background(), envs.Get().GetStorage())
	sm := distribution.NewServiceModel(context.Background(), envs.Get().GetStorage())
	ns, err := nm.List()
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	for _, n := range ns.Items {
		log.Debugf("\n\nrestore service in namespace: %s", n.SelfLink())
		ss, err := sm.List(n.SelfLink())
		if err != nil {
			log.Errorf("%s", err.Error())
			return
		}

		for _, svc := range ss.Items {

			log.Debugf("restore service state: %s \n", svc.SelfLink())
			if _, ok := s.Service[svc.SelfLink()]; !ok {
				s.Service[svc.SelfLink()] = service.NewServiceState(svc)
			}

			s.Service[svc.SelfLink()].Restore()
		}

	}
	log.Info("finish services restore\n\n")
}

func NewState() *State {
	var state = new(State)
	state.Cluster = cluster.NewClusterState()
	state.Service = make(map[string]*service.ServiceState)
	return state
}
