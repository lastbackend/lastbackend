package service

import (
	"errors"
	"fmt"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/service/resource/deployment"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/apis/extensions/v1beta1"
)

type Service struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Replicas  int32             `json:"replicas"`
	config    *v1beta1.Deployment
}

type ServiceList []Service

func Get(namespace, name string) (*Service, *e.Err) {

	var (
		er      error
		ctx     = context.Get()
		service = new(Service)
	)

	detail, er := deployment.GetDeployment(ctx.K8S, namespace, name)
	if er != nil {
		return nil, e.New("service").Unknown(er)
	}

	service.Name = detail.ObjectMeta.Name
	service.Namespace = detail.ObjectMeta.Namespace
	service.Labels = detail.ObjectMeta.Labels

	return service, nil
}

func List(namespace string) (*ServiceList, *e.Err) {
	return &ServiceList{}, nil
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
		s.config = convertReplicationControllerToDeployment(config.(*v1.ReplicationController))
	case *v1.Pod:
		s.config = convertPodToDeployment(config.(*v1.Pod))
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

func convertReplicationControllerToDeployment(config *v1.ReplicationController) *v1beta1.Deployment {

	var deployment = new(v1beta1.Deployment)

	deployment.APIVersion = "extensions/v1beta1"
	deployment.Kind = "Deployment"
	deployment.ObjectMeta = config.ObjectMeta
	deployment.Spec.Strategy.Type = "Recreate"
	deployment.Spec.Replicas = config.Spec.Replicas
	deployment.Spec.Template.Spec = config.Spec.Template.Spec
	deployment.Spec.Template.ObjectMeta = config.Spec.Template.ObjectMeta

	for key, val := range config.Spec.Selector {
		deployment.Spec.Selector.MatchLabels[key] = val
	}

	return deployment
}

func convertPodToDeployment(config *v1.Pod) *v1beta1.Deployment {

	var (
		replicas   int32 = 1
		deployment       = new(v1beta1.Deployment)
	)

	deployment.APIVersion = "extensions/v1beta1"
	deployment.Kind = "Deployment"
	deployment.ObjectMeta = config.ObjectMeta
	deployment.Spec.Strategy.Type = "Recreate"
	deployment.Spec.Replicas = &replicas
	deployment.Spec.Template.Spec = config.Spec

	return deployment
}
