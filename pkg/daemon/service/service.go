package service

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	ctx "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/daemon/service/routes/request"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
)

type service struct {
	Context   context.Context
	Namespace string
}

func New(ctx context.Context, namespace string) *service {
	return &service{
		Context:   ctx,
		Namespace: namespace,
	}
}

func (s *service) List() (*types.ServiceList, error) {
	var storage = ctx.Get().GetStorage()
	return storage.Service().ListByProject(s.Context, s.Namespace)
}

func (s *service) Get(service string) (*types.Service, error) {

	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		svc     *types.Service
	)

	if validator.IsUUID(service) {
		svc, err = storage.Service().GetByID(s.Context, s.Namespace, service)
	} else {
		svc, err = storage.Service().GetByName(s.Context, s.Namespace, service)
	}

	if err != nil {
		log.Error("Error: find service by name", err.Error())
		return svc, err
	}

	return svc, nil
}

func (s *service) Create(rq *request.RequestServiceCreateS) (*types.Service, error) {

	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		svc     *types.Service
	)

	svc, err = storage.Service().Insert(s.Context, s.Namespace, rq.Name, rq.Description, rq.Config)
	if err != nil {
		log.Errorf("Error: insert service to db : %s", err.Error())
		return svc, err
	}

	return svc, nil
}

func (s *service) Update(service *types.Service) (*types.Service, error) {
	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		svc     *types.Service
	)

	svc, err = storage.Service().Update(s.Context, s.Namespace, service)
	if err != nil {
		log.Error("Error: insert service to db", err)
		return svc, err
	}

	return svc, nil
}

func (s *service) Remove(service *types.Service) error {
	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	if err := storage.Service().Remove(s.Context, s.Namespace, service); err != nil {
		log.Error("Error: insert service to db", err)
		return err
	}
	return nil
}
