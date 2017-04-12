package service

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	ctx "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/daemon/service/routes/request"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
)

type Service struct {
}

func New() *Service {
	return new(Service)
}

func (s *Service) List(c context.Context, namespace string) (*types.ServiceList, error) {
	var storage = ctx.Get().GetStorage()
	return storage.Service().ListByProject(c, namespace)
}

func (s *Service) Get(c context.Context, namespace, service string) (*types.Service, error) {

	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		svc     *types.Service
	)

	if validator.IsUUID(service) {
		svc, err = storage.Service().GetByID(c, namespace, service)
	} else {
		svc, err = storage.Service().GetByName(c, namespace, service)
	}

	if err != nil {
		log.Error("Error: find service by name", err.Error())
		return svc, err
	}

	return svc, nil
}

func (s *Service) Create(c context.Context, namespace string, rq *request.RequestServiceCreateS) (*types.Service, error) {

	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		svc     *types.Service
	)

	svc, err = storage.Service().Insert(c, namespace, rq.Name, rq.Description, rq.Config)
	if err != nil {
		log.Errorf("Error: insert service to db : %s", err.Error())
		return svc, err
	}

	return svc, nil
}

func (s *Service) Update(c context.Context, namespace string, service *types.Service) (*types.Service, error) {
	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		svc     *types.Service
	)

	svc, err = storage.Service().Update(c, namespace, service)
	if err != nil {
		log.Error("Error: insert service to db", err)
		return svc, err
	}

	return svc, nil
}

func (s *Service) Remove(c context.Context, namespace string, service *types.Service) error {
	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	if err := storage.Service().Remove(c, namespace, service); err != nil {
		log.Error("Error: insert service to db", err)
		return err
	}
	return nil
}
