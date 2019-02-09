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
	"github.com/lastbackend/lastbackend/pkg/controller/state/job"

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
	Job     map[string]*job.JobState
}

func (s *State) Loop() {

	log.Info("start cluster restore")
	s.Cluster.Loop()
	log.Info("finish cluster restore\n\n")

	log.Info("start namespace restore")
	nm := distribution.NewNamespaceModel(context.Background(), envs.Get().GetStorage())
	sm := distribution.NewServiceModel(context.Background(), envs.Get().GetStorage())
	jm := distribution.NewJobModel(context.Background(), envs.Get().GetStorage())
	vm := distribution.NewVolumeModel(context.Background(), envs.Get().GetStorage())
	dm := distribution.NewDeploymentModel(context.Background(), envs.Get().GetStorage())
	tm := distribution.NewTaskModel(context.Background(), envs.Get().GetStorage())
	pm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	cm := distribution.NewConfigModel(context.Background(), envs.Get().GetStorage())
	sc := distribution.NewSecretModel(context.Background(), envs.Get().GetStorage())

	dr, err := dm.Runtime()
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	tr, err := tm.Runtime()
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	sr, err := sm.Runtime()
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	jr, err := jm.Runtime()
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
		log.V(logLevel).Debugf("\n\nrestore namespace: %s", n.SelfLink())
		ss, err := sm.List(n.SelfLink().String())
		if err != nil {
			log.Errorf("%s", err.Error())
			return
		}

		for _, svc := range ss.Items {

			log.V(logLevel).Debugf("restore service state: %s \n", svc.SelfLink())
			if _, ok := s.Service[svc.SelfLink().String()]; !ok {
				s.Service[svc.SelfLink().String()] = service.NewServiceState(s.Cluster, svc)
			}

			s.Service[svc.SelfLink().String()].Restore()
		}

		js, err := jm.ListByNamespace(n.SelfLink().String())
		if err != nil {
			log.Errorf("%s", err.Error())
			return
		}

		for _, jb := range js.Items {

			log.V(logLevel).Debugf("restore jobs state: %s \n", jb.SelfLink())
			if _, ok := s.Job[jb.SelfLink().String()]; !ok {
				s.Job[jb.SelfLink().String()] = job.NewJobState(s.Cluster, jb)
			}

			s.Job[jb.SelfLink().String()].Restore()
		}

		vl, err := vm.ListByNamespace(n.SelfLink().String())
		if err != nil {
			log.Errorf("%s", err.Error())
			return
		}

		for _, v := range vl.Items {

			log.V(logLevel).Debugf("restore volume state: %s \n", v.SelfLink())
			s.Cluster.SetVolume(v)
		}

	}

	go s.watchPods(context.Background(), &pr.Storage.Revision)
	go s.watchDeployments(context.Background(), &dr.Storage.Revision)
	go s.watchTasks(context.Background(), &tr.Storage.Revision)
	go s.watchServices(context.Background(), &sr.Storage.Revision)
	go s.watchJobs(context.Background(), &jr.Storage.Revision)
	go s.watchVolumes(context.Background(), &vr.Storage.Revision)
	go s.watchSecrets(context.Background(), &scr.Storage.Revision)
	go s.watchConfigs(context.Background(), &cr.Storage.Revision)

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
					_, ok := s.Service[w.Data.SelfLink().String()]
					if ok {
						delete(s.Service, w.Data.SelfLink().String())
					}
					continue
				}

				_, ok := s.Service[w.Data.SelfLink().String()]
				if !ok {
					s.Service[w.Data.SelfLink().String()] = service.NewServiceState(s.Cluster, w.Data)
				}

				s.Service[w.Data.SelfLink().String()].SetService(w.Data)
			}
		}
	}()

	sm.Watch(svc, rev)
}

func (s *State) watchJobs(ctx context.Context, rev *int64) {

	var (
		svc = make(chan types.JobEvent)
	)

	sm := distribution.NewJobModel(ctx, envs.Get().GetStorage())

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
					_, ok := s.Job[w.Data.SelfLink().String()]
					if ok {
						delete(s.Job, w.Data.SelfLink().String())
					}
					continue
				}

				_, ok := s.Job[w.Data.SelfLink().String()]
				if !ok {
					s.Job[w.Data.SelfLink().String()] = job.NewJobState(s.Cluster, w.Data)
				}

				s.Job[w.Data.SelfLink().String()].SetJob(w.Data)
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

				_, sl := w.Data.SelfLink().Parent()
				if w.IsActionRemove() {
					_, ok := s.Service[sl.String()]
					if ok {
						s.Service[sl.String()].DelDeployment(w.Data)
					}
					continue
				}

				_, ok := s.Service[sl.String()]
				if !ok {
					break
				}

				s.Service[sl.String()].SetDeployment(w.Data)
			}
		}
	}()

	dm.Watch(d, rev)
}

func (s *State) watchTasks(ctx context.Context, rev *int64) {

	// Watch pods change
	var (
		d = make(chan types.TaskEvent)
	)

	tm := distribution.NewTaskModel(ctx, envs.Get().GetStorage())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w := <-d:

				if w.Data == nil {
					continue
				}

				_, sl := w.Data.SelfLink().Parent()
				if w.IsActionRemove() {
					_, ok := s.Job[w.Data.Meta.Job]
					if ok {
						s.Job[sl.String()].DelTask(w.Data)
					}
					continue
				}

				_, ok := s.Job[sl.String()]
				if !ok {
					break
				}

				s.Job[sl.String()].SetTask(w.Data)
			}
		}
	}()

	tm.Watch(d, rev)
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

				kind, sl := w.Data.SelfLink().Parent()
				switch kind {
				case types.KindDeployment:
					_, ok := s.Service[sl.String()]
					if !ok {
						break
					}

					if w.IsActionRemove() {
						_, ok := s.Service[sl.String()]
						if ok {
							s.Service[sl.String()].DelPod(w.Data)
						}
						continue
					}

					s.Service[sl.String()].SetPod(w.Data)
				case types.KindTask:
					_, ok := s.Job[sl.String()]
					if !ok {
						break
					}

					if w.IsActionRemove() {
						_, ok := s.Job[sl.String()]
						if ok {
							s.Job[sl.String()].DelPod(w.Data)
						}
						continue
					}

					s.Job[sl.String()].SetPod(w.Data)
				}

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

					ss.CheckDeps(types.StatusDependency{
						Name:   w.Data.Meta.Name,
						Type:   types.KindVolume,
						Status: w.Data.Status.State,
					})
				}

				for _, js := range s.Job {

					if js.Namespace() != w.Data.Meta.Namespace {
						continue
					}

					js.CheckJobDeps(types.StatusDependency{
						Name:   w.Data.Meta.Name,
						Type:   types.KindVolume,
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

				var dep = types.StatusDependency{
					Name:   w.Data.Meta.Name,
					Type:   types.KindSecret,
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

				for _, js := range s.Job {
					if js.Namespace() != w.Data.Meta.Namespace {
						continue
					}
					js.CheckJobDeps(dep)
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

				var dep = types.StatusDependency{
					Name:   w.Data.Meta.Name,
					Type:   types.KindConfig,
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

				for _, js := range s.Job {
					if js.Namespace() != w.Data.Meta.Namespace {
						continue
					}
					js.CheckJobDeps(dep)
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
	state.Job = make(map[string]*job.JobState)
	return state
}
