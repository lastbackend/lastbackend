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
	"context"
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	ctx "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/daemon/node"
	"github.com/lastbackend/lastbackend/pkg/daemon/service/routes/request"
	"github.com/lastbackend/lastbackend/pkg/daemon/storage/store"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

type service struct {
	Context   context.Context
	Namespace types.Meta
}

func New(ctx context.Context, namespace types.Meta) *service {
	return &service{
		Context:   ctx,
		Namespace: namespace,
	}
}

func (s *service) List() (types.ServiceList, error) {
	var (
		storage = ctx.Get().GetStorage()
		list    = types.ServiceList{}
	)

	items, err := storage.Service().ListByNamespace(s.Context, s.Namespace.ID)
	if err != nil {
		return list, err
	}

	for _, item := range items {
		var service = item
		list = append(list, service)
	}

	return list, nil
}

func (s *service) Get(service string) (*types.Service, error) {

	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	svc, err := storage.Service().GetByName(s.Context, s.Namespace.ID, service)
	if err != nil {
		log.Errorf("Error: find service by name: %s", err.Error())
		return nil, err
	}

	return svc, nil
}

func (s *service) Create(rq *request.RequestServiceCreateS) (*types.Service, error) {

	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		svc     = types.Service{}
	)

	log.Debug("Service: create new service")

	svc.Meta = types.ServiceMeta{}
	svc.Meta.SetDefault()

	svc.Meta.Name = rq.Name
	svc.Meta.Region = rq.Region
	svc.Meta.Namespace = s.Namespace.Name
	svc.Meta.Description = rq.Description

	svc.Meta.Replicas = 1
	svc.Pods = make(map[string]*types.Pod)

	if rq.Replicas != nil && *rq.Replicas > 0 {
		svc.Meta.Replicas = *rq.Replicas
	}

	spec, err := createSpec(rq.Spec)
	if err != nil {
		log.Errorf("Error: create spec from request opts : %s", err.Error())
		return &svc, err
	}

	svc.Spec = make(map[string]*types.ServiceSpec)
	svc.Spec[uuid.NewV4().String()] = spec

	log.Debugf("Service: Create: add pods : %d", svc.Meta.Replicas)
	for i := 0; i < svc.Meta.Replicas; i++ {
		log.Debug("Service: Create: add new pod")
		if err := s.AddPod(&svc); err != nil {
			log.Errorf("Service: Create: add new pod error: %s", err.Error())
			return &svc, err
		}
	}

	s.StateUpdate(&svc)
	s.ResourcesUpdate(&svc)

	if err = storage.Service().Insert(s.Context, &svc); err != nil {
		log.Errorf("Error: insert service to db : %s", err.Error())
		return &svc, err
	}

	return &svc, nil
}

func (s *service) Update(service *types.Service, rq *request.RequestServiceUpdateS) error {

	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	log.Debug("Service: update service info and config")

	if rq.Name != "" {
		service.Meta.Name = rq.Name
	}

	if rq.Description != nil {
		service.Meta.Description = *rq.Description
	}

	if rq.Replicas != nil {
		log.Debugf("Service: Update: set replicas: %d", *rq.Replicas)
		service.Meta.Replicas = *rq.Replicas
		s.Scale(service)
	}

	s.StateUpdate(service)
	s.ResourcesUpdate(service)

	if err = storage.Service().Update(s.Context, service); err != nil {
		log.Error("Error: insert service to db", err)
		return err
	}

	return nil
}

func (s *service) Remove(service *types.Service) error {
	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	service.State.State = types.StateDestroy

	if len(service.Pods) == 0 {
		if err := storage.Service().Remove(s.Context, service); err != nil {
			log.Error("Error: insert service to db", err)
			return err
		}
		return nil
	}

	for _, pod := range service.Pods {
		pod.State.State = types.StateDestroy
	}

	if err := storage.Service().Update(s.Context, service); err != nil {
		log.Error("Error: insert service to db", err)
		return err
	}
	return nil
}

func (s *service) AddPod(service *types.Service) error {

	var (
		log = ctx.Get().GetLogger()
	)

	log.Debug("Create new pod state on service")

	pod := types.Pod{}

	pod.Meta.ID = uuid.NewV4().String()
	pod.Meta.Created = time.Now()
	pod.Meta.Updated = time.Now()
	pod.State.Provision = true
	pod.Spec.State = types.StateStarted

	if len(service.Pods) > 0 {
		for _, p := range service.Pods {
			pod.Spec = p.Spec
			break
		}
	} else {
		pod.Spec = s.GenerateSpec(service)
	}

	n, err := node.New().Allocate(s.Context, pod.Spec)
	if err != nil {
		return err
	}

	log.Debugf("Service: Add pod: Node meta: %s", n.Meta)
	pod.Meta.Hostname = n.Meta.Hostname
	service.Pods[pod.Meta.ID] = &pod

	return nil
}

func (s *service) DelPod(service *types.Service) error {

	var (
		log = ctx.Get().GetLogger()
	)

	log.Debug("Delete pod service")

	for _, pod := range service.Pods {
		if pod.Spec.State != types.StateDestroy {
			log.Debugf("Mark pod for deletion: %s", pod.Meta.ID)
			pod.State.Provision = true
			pod.State.Ready = false
			pod.Spec.State = types.StateDestroy

			for _, c := range pod.Containers {
				c.State = types.StateProvision
			}

			break
		}
	}

	return nil
}

func (s *service) SetPods(pods []types.Pod) error {

	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	for _, pod := range pods {
		log.Debugf("update pod %s state: %s", pod.Meta.ID, pod.State.State)
		svc, err := storage.Service().GetByPodID(s.Context, pod.Meta.ID)
		if err != nil {
			log.Errorf("Error: get service by pod ID %s from db: %s", pod.Meta.ID, err.Error())
			if err.Error() == store.ErrKeyNotFound {
				log.Debugf("Pod %s not found", pod.Meta.ID)
				continue
			}
			return err
		}

		p, e := storage.Pod().GetByID(s.Context, svc.Meta.Namespace, svc.Meta.ID, pod.Meta.ID)

		if e != nil {
			log.Errorf("Error: get pod from db: %s", e.Error())
			continue
		}

		p.Containers = pod.Containers
		p.State = pod.State

		if p.State.State == types.StateDestroyed {
			log.Debugf("Service: Set pods: remove deleted pod: %s", p.Meta.ID)
			if err := storage.Pod().Remove(s.Context, svc.Meta.Namespace, svc.Meta.ID, p); err != nil {
				log.Errorf("Error: set pod to db: %s", err)
				return err
			}
			delete(svc.Pods, p.Meta.ID)

			if len(svc.Pods) == 0 && svc.State.State == types.StateDestroy {
				storage.Service().Remove(s.Context, svc)
			}

			return nil
		}

		if err := storage.Pod().Update(s.Context, svc.Meta.Namespace, svc.Meta.ID, p); err != nil {
			log.Errorf("Error: set pod to db: %s", err)
			return err
		}

	}

	return nil
}

func (s *service) StateUpdate(service *types.Service) {

	service.State.Replicas = types.ServiceReplicasState{}

	for _, p := range service.Pods {
		service.State.Replicas.Total++
		switch p.State.State {
		case types.StateCreated: service.State.Replicas.Created++
		case types.StateStarted: service.State.Replicas.Running++
		case types.StateStopped: service.State.Replicas.Stopped++
		case types.StateError: service.State.Replicas.Errored++
		}

		if p.State.Provision {
			service.State.Replicas.Provision++
		}

		if p.State.Ready {
			service.State.Replicas.Ready++
		}
	}

}

func (s *service) ResourcesUpdate(service *types.Service) {

	service.State.Resources = types.ServiceResourcesState{}

	for _, s := range service.Spec {
		service.State.Resources.Memory += int(s.Memory)*service.Meta.Replicas
	}

}

func (s *service) Scale(service *types.Service) error {
	var (
		log      = ctx.Get().GetLogger()
		replicas int
	)

	for _, pod := range service.Pods {
		if pod.Spec.State != types.StateDestroy {
			replicas++
		}
	}

	log.Debugf("Service: Scale: current replicas: %d", replicas)

	if replicas == service.Meta.Replicas {
		log.Debug("Service: Replicas not needed, replicas equal")
		return nil
	}

	if replicas < service.Meta.Replicas {
		log.Debug("Service: Replicas: create a new replicas")
		for i := 0; i < (service.Meta.Replicas - replicas); i++ {
			if err := s.AddPod(service); err != nil {
				return err
			}
		}
	}

	if replicas > service.Meta.Replicas {
		log.Debug("Service: Replicas: remove  unneeded replicas")
		for i := 0; i < (replicas - service.Meta.Replicas); i++ {
			if err := s.DelPod(service); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *service) AddSpec(service *types.Service, rq *request.RequestServiceSpecCreateS) error {

	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	log.Debug("Add spec service")

	spec, err := createSpec(rq)
	if err != nil {
		log.Errorf("Error: create spec from request opts : %s", err.Error())
		return err
	}

	service.Spec[spec.Meta.ID] = spec

	for _, pod := range service.Pods {
		pod.Spec = s.GenerateSpec(service)
		pod.State.Provision = true
		pod.State.Ready = false
		service.Pods[pod.Meta.ID] = pod
	}

	if err := storage.Service().Update(s.Context, service); err != nil {
		log.Errorf("Error: AddSpec: update service spec to db : %s", err.Error())
		return err
	}

	return nil
}

func (s *service) SetSpec(service *types.Service, id string, rq *request.RequestServiceSpecUpdateS) error {

	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	log.Debug("Set spec service")

	spec := service.Spec[id]

	updateSpec(rq, spec)

	delete(service.Spec, id)
	service.Spec[spec.Meta.ID] = spec

	for _, pod := range service.Pods {

		if pod.Spec.State == types.StateDestroy {
			continue
		}

		pod.Spec = s.GenerateSpec(service)
		pod.State.Provision = true
		pod.State.Ready = false

		for _, c := range pod.Containers {
			c.State = types.StateProvision
		}

		service.Pods[pod.Meta.ID] = pod
	}

	if err := storage.Service().Update(s.Context, service); err != nil {
		log.Errorf("Error: AddSpec: update service spec to db : %s", err.Error())
		return err
	}

	return nil
}

func (s *service) DelSpec(service *types.Service, id string) error {

	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	buf, _ := json.Marshal(service)

	log.Info(string(buf))

	log.Debug("Delete spec service")

	if _, ok := service.Spec[id]; !ok {
		return nil
	}

	delete(service.Spec, id)

	for _, pod := range service.Pods {
		pod.Spec = s.GenerateSpec(service)
		pod.State.Provision = true
		pod.State.Ready = false
		service.Pods[pod.Meta.ID] = pod

		for _, c := range pod.Containers {
			c.State = types.StateProvision
		}
	}

	if err := storage.Service().Update(s.Context, service); err != nil {
		log.Errorf("Error: AddSpec: update service spec to db : %s", err.Error())
		return err
	}

	return nil
}

func (s *service) GenerateSpec(service *types.Service) types.PodSpec {

	var (
		log = ctx.Get().GetLogger()
	)

	log.Debug("Generate new node pod spec")

	var spec = types.PodSpec{}
	spec.ID = uuid.NewV4().String()
	spec.Created = time.Now()
	spec.Updated = time.Now()
	spec.Containers = make(map[string]*types.ContainerSpec)

	for _, spc := range service.Spec {

		cs := new(types.ContainerSpec)

		cs.Meta.SetDefault()
		cs.Meta.ID = spc.Meta.ID
		cs.Meta.Labels = spc.Meta.Labels
		cs.Meta.Created = time.Now()
		cs.Meta.Updated = time.Now()

		cs.Image = types.ImageSpec{
			Name: spc.Image,
			Pull: true,
		}

		for _, port := range spc.Ports {
			cs.Ports = append(cs.Ports, types.ContainerPortSpec{
				ContainerPort: port.Container,
				Protocol:      port.Protocol,
			})
		}

		cs.Command = spc.Command
		cs.Entrypoint = spc.Entrypoint
		cs.Envs = spc.EnvVars
		cs.Quota = types.ContainerQuotaSpec{
			Memory: spc.Memory,
		}

		cs.RestartPolicy = types.ContainerRestartPolicySpec{
			Name:    "always",
			Attempt: 0,
		}

		spec.Containers[cs.Meta.ID] = cs
	}

	spec.State = types.StateStarted

	return spec
}

func createSpec(opts *request.RequestServiceSpecCreateS) (*types.ServiceSpec, error) {

	spec := new(types.ServiceSpec)
	spec.Meta.SetDefault()

	spec.Memory = int64(32)

	if opts.Memory != nil {
		spec.Memory = *opts.Memory
	}

	if opts.Command != nil {
		spec.Command = strings.Split(*opts.Command, " ")
	}

	if opts.Image != nil {
		spec.Image = *opts.Image
	}

	if opts.Entrypoint != nil {
		spec.Entrypoint = strings.Split(*opts.Entrypoint, " ")
	}

	if opts.EnvVars != nil {
		spec.EnvVars = *opts.EnvVars
	}

	if opts.Ports != nil {
		spec.Ports = []types.Port{}
		for _, val := range *opts.Ports {
			spec.Ports = append(spec.Ports, types.Port{
				Protocol:  val.Protocol,
				Container: val.Internal,
				Host:      val.External,
				Published: val.Published,
			})
		}
	}

	return spec, nil
}

func updateSpec(opts *request.RequestServiceSpecUpdateS, spec *types.ServiceSpec) error {

	if spec == nil {
		return errors.New("Error: spec is nil")
	}

	spec.Meta.Parent = spec.Meta.ID
	spec.Meta.ID = uuid.NewV4().String()
	spec.Meta.Revision++

	spec.Memory = int64(32)

	if opts.Memory != nil {
		spec.Memory = *opts.Memory
	}

	if opts.Command != nil {
		spec.Command = strings.Split(*opts.Command, " ")
	}

	if opts.Image != nil {
		spec.Image = *opts.Image
	}

	if opts.Entrypoint != nil {
		spec.Entrypoint = strings.Split(*opts.Entrypoint, " ")
	}

	if opts.EnvVars != nil {
		spec.EnvVars = *opts.EnvVars
	}

	if opts.Ports != nil {
		spec.Ports = []types.Port{}
		for _, val := range *opts.Ports {
			spec.Ports = append(spec.Ports, types.Port{
				Protocol:  val.Protocol,
				Container: val.Internal,
				Host:      val.External,
				Published: val.Published,
			})
		}
	}

	return nil
}
