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
	ctx "github.com/lastbackend/lastbackend/pkg/api/context"
	"github.com/lastbackend/lastbackend/pkg/api/service/routes/request"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/satori/go.uuid"
	"strings"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logLevel = 3

type service struct {
	Context context.Context
	App     types.Meta
}

func New(ctx context.Context, app types.Meta) *service {
	return &service{
		Context: ctx,
		App:     app,
	}
}

func (s *service) List() (types.ServiceList, error) {
	var (
		storage = ctx.Get().GetStorage()
		list    = types.ServiceList{}
	)

	log.V(logLevel).Debug("Service: list service")

	items, err := storage.Service().ListByApp(s.Context, s.App.Name)
	if err != nil {
		log.V(logLevel).Error("Service: list service err: %s", err.Error())
		return nil, err
	}

	log.V(logLevel).Debugf("Service: list service result: %d", len(items))

	for _, item := range items {
		var srv = item
		list = append(list, srv)
	}

	return list, nil
}

func (s *service) Get(service string) (*types.Service, error) {

	var (
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Service: get service %s", service)

	svc, err := storage.Service().GetByName(s.Context, s.App.Name, service)
	if err != nil {
		if err.Error() == store.ErrKeyNotFound {
			log.V(logLevel).Warnf("Service: service by name `%s` not found", service)
			return nil, nil
		}
		log.V(logLevel).Errorf("Service: get service by name `%s` err: %s", service, err.Error())
		return nil, err
	}
	return svc, nil
}

// TODO: check available names
func (s *service) Create(rq *request.RequestServiceCreateS) (*types.Service, error) {

	var (
		storage = ctx.Get().GetStorage()
		svc     = types.Service{}
	)

	log.V(logLevel).Debugf("Service: create service %#v", rq)

	svc.Meta = types.ServiceMeta{}
	svc.Meta.SetDefault()
	svc.Meta.Name = rq.Name
	svc.Meta.Region = rq.Region
	svc.Meta.App = s.App.Name
	svc.Meta.Description = rq.Description
	svc.Meta.Replicas = 1
	svc.State.State = types.StateProvision
	svc.Pods = make(map[string]*types.Pod)

	if rq.Replicas != nil && *rq.Replicas > 0 {
		svc.Meta.Replicas = *rq.Replicas
	}

	spec := generateSpec(rq.Spec, nil)

	svc.Spec = make(map[string]*types.ServiceSpec)
	svc.Spec[uuid.NewV4().String()] = spec

	hook := types.Hook{}
	hook.Meta.SetDefault()
	hook.Meta.ID = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	hook.App = svc.Meta.App
	hook.Service = svc.Meta.Name

	svc.Meta.Hook = hook.Meta.ID

	if err := storage.Service().Insert(s.Context, &svc); err != nil {
		log.V(logLevel).Errorf("Service: insert service err: %s", err.Error())
		return nil, err
	}

	if err := storage.Hook().Insert(s.Context, &hook); err != nil {
		log.V(logLevel).Errorf("Service: insert service hook err: %s", err.Error())
		return nil, err
	}

	return &svc, nil
}

func (s *service) Update(service *types.Service, rq *request.RequestServiceUpdateS) error {

	var (
		err     error
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Service: update service %#v -> %#v", service, rq)

	if rq.Name != "" {
		service.Meta.Name = rq.Name
	}

	if rq.Description != nil {
		service.Meta.Description = *rq.Description
	}

	if rq.Replicas != nil && *rq.Replicas > 0 {
		log.Warnf("Service: set replicas: %d", *rq.Replicas)
		service.Meta.Replicas = *rq.Replicas
	}

	service.State.State = types.StateProvision

	if err = storage.Service().Update(s.Context, service); err != nil {
		log.V(logLevel).Errorf("Service: update service err: %s", err.Error())
		return err
	}

	if err = storage.Service().UpdateSpec(s.Context, service); err != nil {
		log.V(logLevel).Errorf("Service: update service spec err: %s", err.Error())
		return err
	}

	return nil
}

func (s *service) Remove(service *types.Service) error {
	var (
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Service: remove service %#v", service)

	service.State.State = types.StateDestroyed
	service.Meta.Replicas = int(0) // Delete all pods

	if err := storage.Service().Update(s.Context, service); err != nil {
		log.V(logLevel).Errorf("Service: update service state err: %s", err.Error())
		return err
	}

	if len(service.Pods) == 0 {
		if err := storage.Service().Remove(s.Context, service); err != nil {
			log.V(logLevel).Errorf("Service: remove service err: %s", err.Error())
			return err
		}
		return nil
	}

	if err := storage.Service().UpdateSpec(s.Context, service); err != nil {
		log.V(logLevel).Errorf("Service: update service spec err: %s", err.Error())
		return err
	}

	return nil
}

func (s *service) AddSpec(service *types.Service, rq *request.RequestServiceSpecS) error {

	var (
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Service: add service `%s` spec from data %#v", service.Meta.Name, rq)

	spec := generateSpec(rq, nil)
	service.Spec[spec.Meta.ID] = spec

	if err := storage.Service().Update(s.Context, service); err != nil {
		log.V(logLevel).Errorf("Service: update service err: %s", err.Error())
		return err
	}

	if err := storage.Service().UpdateSpec(s.Context, service); err != nil {
		log.V(logLevel).Errorf("Service: update service spec err: %s", err.Error())
		return err
	}

	return nil
}

func (s *service) SetSpec(service *types.Service, id string, rq *request.RequestServiceSpecS) error {

	var (
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Service: set service `%s` spec %s from data %#v", service.Meta.Name, id, rq)

	spec := generateSpec(rq, service.Spec[id])
	service.Spec[spec.Meta.ID] = spec
	service.State.State = types.StateProvision
	delete(service.Spec, spec.Meta.Parent)

	if err := storage.Service().Update(s.Context, service); err != nil {
		log.V(logLevel).Errorf("Service: update service err: %s", err.Error())
		return err
	}

	if err := storage.Service().UpdateSpec(s.Context, service); err != nil {
		log.V(logLevel).Errorf("Service: update service spec err: %s", err.Error())
		return err
	}

	return nil
}

func (s *service) DelSpec(service *types.Service, id string) error {

	var (
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Service: delete service `%s` spec %#v", service.Meta.Name, id)

	if _, ok := service.Spec[id]; !ok {
		return nil
	}

	delete(service.Spec, id)

	if err := storage.Service().Update(s.Context, service); err != nil {
		log.V(logLevel).Errorf("Service: update service err: %s", err.Error())
		return err
	}

	if err := storage.Service().UpdateSpec(s.Context, service); err != nil {
		log.V(logLevel).Errorf("Service: update service spec err: %s", err.Error())
		return err
	}

	return nil
}

func (s *service) Redeploy(service *types.Service) error {

	var (
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Service: redeploy service %#v", service)

	specs := make(map[string]*types.ServiceSpec)
	service.State.State = types.StateProvision
	for id := range service.Spec {
		sp := service.Spec[id]
		sp.Meta.Parent = id
		sp.Meta.ID = uuid.NewV4().String()
		delete(service.Spec, id)
		specs[sp.Meta.ID] = sp
	}

	service.Spec = specs

	if err := storage.Service().Update(s.Context, service); err != nil {
		log.V(logLevel).Errorf("Service: update service err: %s", err.Error())
		return err
	}

	if err := storage.Service().UpdateSpec(s.Context, service); err != nil {
		log.V(logLevel).Errorf("Service: update service spec err: %s", err.Error())
		return err
	}

	return nil
}

func generateSpec(opts *request.RequestServiceSpecS, spec *types.ServiceSpec) *types.ServiceSpec {

	s := types.ServiceSpec{}
	if spec != nil {
		s = *spec
		s.Meta.Parent = spec.Meta.ID
	}

	s.Meta.SetDefault()
	s.Meta.ID = uuid.NewV4().String()
	s.Memory = int64(32)

	if opts.Memory != nil {
		s.Memory = *opts.Memory
	}

	if opts.Command != nil {
		s.Command = strings.Split(*opts.Command, " ")
	}

	if opts.Image != nil {
		s.Image = *opts.Image
	}

	if opts.Entrypoint != nil {
		s.Entrypoint = strings.Split(*opts.Entrypoint, " ")
	}

	if opts.EnvVars != nil {
		s.EnvVars = *opts.EnvVars
	}

	if opts.Ports != nil {
		s.Ports = []types.Port{}
		for _, val := range *opts.Ports {
			s.Ports = append(s.Ports, types.Port{
				Protocol:  val.Protocol,
				Container: val.Internal,
				Host:      val.External,
				Published: val.Published,
			})
		}
	}

	return &s
}
