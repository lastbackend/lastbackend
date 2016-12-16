package service

import (
	"errors"
	"fmt"
	"github.com/lastbackend/lastbackend/libs/adapter/k8s/converter"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/service/resource/deployment"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/apis/extensions/v1beta1"
)

type Service struct {
	deployment.Deployment
	config *v1beta1.Deployment
}

func Get(namespace, name string) (*Service, *e.Err) {

	var (
		er  error
		ctx = context.Get()
	)

	detail, er := deployment.Get(ctx.K8S, namespace, name)
	if er != nil {
		return nil, e.New("service").Unknown(er)
	}

	return &Service{*detail, nil}, nil
}

func List(namespace string) (map[string]*Service, *e.Err) {

	var (
		er          error
		ctx         = context.Get()
		serviceList = make(map[string]*Service)
	)

	detailList, er := deployment.List(ctx.K8S, namespace)
	if er != nil {
		return nil, e.New("service").Unknown(er)
	}

	for _, val := range detailList {
		serviceList[val.ObjectMeta.Name] = &Service{val, nil}
	}

	return serviceList, nil
}

func Create(user, project string, config interface{}) (*Service, *e.Err) {

	var (
		ctx     = context.Get()
		s       = new(Service)
		service = new(model.Service)
	)

	switch config.(type) {
	case *v1beta1.Deployment:
		s.config = config.(*v1beta1.Deployment)
	case *v1.ReplicationController:
		s.config = converter.Convert_ReplicationController_to_Deployment(config.(*v1.ReplicationController))
	case *v1.Pod:
		s.config = converter.Convert_Pod_to_Deployment(config.(*v1.Pod))
	default:
		return nil, e.New("service").Unknown(errors.New("unknown type config"))
	}

	service.User = user
	service.Project = project
	service.Name = fmt.Sprintf("%s-%s", s.config.Name, generator.GetUUIDV4()[0:12])

	service, err := ctx.Storage.Service().Insert(service)
	if err != nil {
		return nil, err
	}

	s.config.Name = service.Name

	return s, nil
}

func (s Service) Deploy(namespace string) *e.Err {

	var (
		er  error
		ctx = context.Get()
	)

	_, er = ctx.K8S.Extensions().Deployments(namespace).Create(s.config)
	if er != nil {
		return e.New("service").Unknown(er)
	}

	return nil
}
