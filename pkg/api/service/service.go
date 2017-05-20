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

	items, err := storage.Service().ListByNamespace(s.Context, s.Namespace.Name)
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

	svc, err := storage.Service().GetByName(s.Context, s.Namespace.Name, service)
	if err != nil && err.Error() != store.ErrKeyNotFound {
		log.Errorf("Error: find service by name: %s", err.Error())
		return nil, err
	}
	return svc, nil
}

// TODO: check available names
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
	svc.State.State = types.StateProvision
	svc.Pods = make(map[string]*types.Pod)

	if rq.Replicas != nil && *rq.Replicas > 0 {
		svc.Meta.Replicas = *rq.Replicas
	}

	spec := generateSpec(rq.Spec, nil)

	svc.Spec = make(map[string]*types.ServiceSpec)
	svc.Spec[uuid.NewV4().String()] = spec

	if err := storage.Service().Insert(s.Context, &svc); err != nil {
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

	if rq.Replicas != nil && *rq.Replicas > 0 {
		log.Debugf("Service: Update: set replicas: %d", *rq.Replicas)
		service.Meta.Replicas = *rq.Replicas
	}

	service.State.State = types.StateProvision

	if err = storage.Service().Update(s.Context, service); err != nil {
		log.Error("Error: update service info to db", err)
		return err
	}

	if err = storage.Service().UpdateSpec(s.Context, service); err != nil {
		log.Error("Error: update service spec to db", err)
		return err
	}

	return nil
}

func (s *service) Remove(service *types.Service) error {
	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	service.State.State = types.StateDestroyed
	service.Meta.Replicas = int(0) // Delete all pods

	if err := storage.Service().Update(s.Context, service); err != nil {
		log.Error("Error: update service info to db", err)
		return err
	}

	if len(service.Pods) == 0 {
		if err := storage.Service().Remove(s.Context, service); err != nil {
			log.Error("Error: insert service to db", err)
			return err
		}
		return nil
	}

	if err := storage.Service().UpdateSpec(s.Context, service); err != nil {
		log.Error("Error: update service spec info to db", err)
		return err
	}

	return nil
}

func (s *service) AddSpec(service *types.Service, rq *request.RequestServiceSpecS) error {

	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	log.Debug("Add spec service")

	spec := generateSpec(rq, nil)
	service.Spec[spec.Meta.ID] = spec

	if err := storage.Service().Update(s.Context, service); err != nil {
		log.Errorf("Error: DelSpec: update service info to db : %s", err.Error())
		return err
	}

	if err := storage.Service().UpdateSpec(s.Context, service); err != nil {
		log.Errorf("Error: AddSpec: update service spec to db : %s", err.Error())
		return err
	}

	return nil
}

func (s *service) SetSpec(service *types.Service, id string, rq *request.RequestServiceSpecS) error {

	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	log.Debug("Set spec service")

	spec := generateSpec(rq, service.Spec[id])
	service.Spec[spec.Meta.ID] = spec
	service.State.State = types.StateProvision
	delete(service.Spec, spec.Meta.Parent)

	if err := storage.Service().Update(s.Context, service); err != nil {
		log.Errorf("Error: SetSpec: update service info to db : %s", err.Error())
		return err
	}

	if err := storage.Service().UpdateSpec(s.Context, service); err != nil {
		log.Errorf("Error: SetSpec: update service spec to db : %s", err.Error())
		return err
	}

	return nil
}

func (s *service) DelSpec(service *types.Service, id string) error {

	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	log.Debug("Delete spec service")

	if _, ok := service.Spec[id]; !ok {
		return nil
	}

	delete(service.Spec, id)

	if err := storage.Service().Update(s.Context, service); err != nil {
		log.Errorf("Error: DelSpec: update service info to db : %s", err.Error())
		return err
	}

	if err := storage.Service().UpdateSpec(s.Context, service); err != nil {
		log.Errorf("Error: DelSpec: update service spec to db : %s", err.Error())
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
