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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logLevel = 3

type State struct {
	Cluster *cluster.ClusterState
	Service map[string]*service.ServiceState
}

func (s *State) Loop() {

	log.Info("start cluster restore")
	s.Cluster.Loop()
	log.Info("finish cluster restore\n\n")

	log.Info("start services restore")
	nm := distribution.NewNamespaceModel(context.Background(), envs.Get().GetStorage())
	sm := distribution.NewServiceModel(context.Background(), envs.Get().GetStorage())
	vm := distribution.NewVolumeModel(context.Background(), envs.Get().GetStorage())
	dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
	pm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	cm := distribution.NewConfigModel(context.Background(), envs.Get().GetStorage())
	sc := distribution.NewSecretModel(context.Background(), envs.Get().GetStorage())

	dr, err := dm.Runtime()
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	sr, err := sm.Runtime()
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	vr, err := vm.Runtime()
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	pr, err := pm.Runtime()
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	cr, err := cm.Runtime()
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	scr, err := sc.Runtime()
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	ns, err := nm.List()
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	for _, n := range ns.Items {
		log.V(logLevel).Debugf("\n\nrestore service in namespace: %s", n.SelfLink())
		ss, err := sm.List(n.SelfLink())
		if err != nil {
			log.Errorf("%s", err.Error())
			return
		}

		for _, svc := range ss.Items {

			log.V(logLevel).Debugf("restore service state: %s \n", svc.SelfLink())
			if _, ok := s.Service[svc.SelfLink()]; !ok {
				s.Service[svc.SelfLink()] = service.NewServiceState(s.Cluster, svc)
			}

			s.Service[svc.SelfLink()].Restore()
		}

		vl, err := vm.ListByNamespace(n.SelfLink())
		if err != nil {
			log.Errorf("%s", err.Error())
			return
		}

		for _, v := range vl.Items {

			log.V(logLevel).Debugf("restore volume state: %s \n", v.SelfLink())
			s.Cluster.SetVolume(v)
		}

	}

	go s.watchPods(context.Background(), &pr.System.Revision)
	go s.watchDeployments(context.Background(), &dr.System.Revision)
	go s.watchServices(context.Background(), &sr.System.Revision)
	go s.watchVolumes(context.Background(), &vr.System.Revision)
	go s.watchSecrets(context.Background(), &scr.System.Revision)
	go s.watchConfigs(context.Background(), &cr.System.Revision)

	log.Info("finish services restore\n\n")
}

func (s *State) watchServices(ctx context.Context, rev *int64) {

	var (
		svc = make(chan types.ServiceEvent)
	)

	sm := distribution.NewServiceModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-svc:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					_, ok := s.Service[w.Data.SelfLink()]
					if ok {
						delete(s.Service, w.Data.SelfLink())
					}
					continue
				}

				_, ok := s.Service[w.Data.SelfLink()]
				if !ok {
					s.Service[w.Data.SelfLink()] = service.NewServiceState(s.Cluster, w.Data)
				}

				s.Service[w.Data.SelfLink()].SetService(w.Data)
			}
		}
	}()

	sm.Watch(svc, rev)
}

func (s *State) watchDeployments(ctx context.Context, rev *int64) {

	// Watch pods change
	var (
		d = make(chan types.DeploymentEvent)
	)

	dm := distribution.NewDeploymentModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-d:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					_, ok := s.Service[w.Data.ServiceLink()]
					if ok {
						s.Service[w.Data.ServiceLink()].DelDeployment(w.Data)
					}
					continue
				}

				_, ok := s.Service[w.Data.ServiceLink()]
				if !ok {
					break
				}

				s.Service[w.Data.ServiceLink()].SetDeployment(w.Data)
			}
		}
	}()

	dm.Watch(d, rev)
}

func (s *State) watchPods(ctx context.Context, rev *int64) {

	// Watch pods change
	var (
		p = make(chan types.PodEvent)
	)

	pm := distribution.NewPodModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-p:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					_, ok := s.Service[w.Data.ServiceLink()]
					if ok {
						s.Service[w.Data.ServiceLink()].DelPod(w.Data)
					}
					continue
				}

				_, ok := s.Service[w.Data.ServiceLink()]
				if !ok {
					break
				}

				s.Service[w.Data.ServiceLink()].SetPod(w.Data)
			}
		}
	}()

	pm.Watch(p, rev)
}

func (s *State) watchVolumes(ctx context.Context, rev *int64) {
	var (
		vl = make(chan types.VolumeEvent)
	)

	vm := distribution.NewVolumeModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-vl:

				if w.Data == nil {
					continue
				}

				if w.IsActionRemove() {
					s.Cluster.DelVolume(w.Data)
					continue
				}

				s.Cluster.SetVolume(w.Data)
				for _, ss := range s.Service {

					if ss.Namespace() != w.Data.Meta.Namespace {
						continue
					}

					ss.CheckDeps(types.DeploymentStatusDependency{
						Name: w.Data.Meta.Name,
						Type: types.KindVolume,
						Status: w.Data.Status.State,
					})
				}
			}
		}
	}()

	vm.Watch(vl, rev)
}

func (s *State) watchSecrets(ctx context.Context, rev *int64) {

	var (
		vl = make(chan types.SecretEvent)
	)

	sm := distribution.NewSecretModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-vl:

				if w.Data == nil {
					continue
				}

				var dep = types.DeploymentStatusDependency{
					Name: w.Data.Meta.Name,
					Type: types.KindSecret,
					Status: types.StateReady,
				}

				if w.IsActionRemove() {
					dep.Status = types.StateNotReady
				}


				for _, ss := range s.Service {
					if ss.Namespace() != w.Data.Meta.Namespace {
						continue
					}
					ss.CheckDeps(dep)
				}
			}
		}
	}()

	sm.Watch(vl, rev)
}

func (s *State) watchConfigs(ctx context.Context, rev *int64) {

	var (
		vl = make(chan types.ConfigEvent)
	)

	sm := distribution.NewConfigModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-vl:

				if w.Data == nil {
					continue
				}

				var dep = types.DeploymentStatusDependency{
					Name: w.Data.Meta.Name,
					Type: types.KindConfig,
					Status: types.StateReady,
				}

				if w.IsActionRemove() {
					dep.Status = types.StateNotReady
				}


				for _, ss := range s.Service {
					if ss.Namespace() != w.Data.Meta.Namespace {
						continue
					}
					ss.CheckDeps(dep)
				}
			}
		}
	}()

	sm.Watch(vl, rev)
}

func NewState() *State {
	var state = new(State)
	state.Cluster = cluster.NewClusterState()
	state.Service = make(map[string]*service.ServiceState)
	return state
}
