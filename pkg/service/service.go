package service

import (
	"errors"
	"github.com/lastbackend/lastbackend/libs/adapter/k8s/converter"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/interface/k8s"
	"github.com/lastbackend/lastbackend/pkg/service/resource/deployment"
	"github.com/unloop/gopipe"
	"io"
	"k8s.io/client-go/1.5/pkg/api/unversioned"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/pkg/apis/extensions/v1beta1"
	"time"
)

type Service struct {
	deployment.Deployment
	config *v1beta1.Deployment
}

func Get(client k8s.IK8S, namespace, name string) (*Service, *e.Err) {

	var er error

	detail, er := deployment.Get(client, namespace, name)
	if er != nil {
		return nil, e.New("service").Unknown(er)
	}

	return &Service{*detail, nil}, nil
}

func List(client k8s.IK8S, namespace string) (map[string]*Service, *e.Err) {

	var (
		er          error
		serviceList = make(map[string]*Service)
	)

	detailList, er := deployment.List(client, namespace)
	if er != nil {
		return nil, e.New("service").Unknown(er)
	}

	for _, val := range detailList {
		serviceList[val.ObjectMeta.Name] = &Service{val, nil}
	}

	return serviceList, nil
}

func Create(name string, config interface{}) (*Service, *e.Err) {

	var s = new(Service)

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

	s.config.Name = name

	return s, nil
}

func Update(client k8s.IK8S, namespace, name string, config *ServiceConfig) *e.Err {

	var er error

	dp, er := client.Extensions().Deployments(namespace).Get(name)
	if er != nil {
		return e.New("service").Unknown(er)
	}

	config.update(dp)

	er = deployment.Update(client, namespace, dp)
	if er != nil {
		return e.New("service").Unknown(er)
	}

	return nil
}

type ServiceLogsOption struct {
	Stream       io.Writer
	Pod          string
	Container    string
	Follow       bool
	Previous     bool
	Timestamps   bool
	SinceSeconds *int64
	SinceTime    *time.Time
	TailLines    *int64
	LimitBytes   *int64
}

func Logs(client k8s.IK8S, namespace string, opts *ServiceLogsOption, close chan bool) *e.Err {

	var (
		er     error
		s      = stream.New(opts.Stream)
		option = v1.PodLogOptions{
			Container:  opts.Container,
			Follow:     opts.Follow,
			Previous:   opts.Previous,
			Timestamps: opts.Timestamps,
		}
	)

	if opts.SinceSeconds != nil {
		option.SinceSeconds = opts.SinceSeconds
	}

	if opts.SinceTime != nil {
		t := unversioned.Time{}
		t.Time = *opts.SinceTime
		option.SinceTime = &t
	}

	if opts.TailLines != nil {
		option.TailLines = opts.TailLines
	}

	if opts.LimitBytes != nil {
		option.LimitBytes = opts.LimitBytes
	}

	req := client.Core().Pods(namespace).GetLogs(opts.Pod, &option)

	readCloser, err := req.Stream()
	if err != nil {
		return e.New("service").Unknown(er)
	}
	defer readCloser.Close()

	go s.Pipe(&readCloser)

	<-close

	s.Close()

	return nil
}

func (s *Service) Deploy(client k8s.IK8S, namespace string) *e.Err {

	var er error

	_, er = client.Extensions().Deployments(namespace).Create(s.config)
	if er != nil {
		return e.New("service").Unknown(er)
	}

	detail, er := deployment.Get(client, namespace, s.ObjectMeta.Name)
	if er != nil {
		return e.New("service").Unknown(er)
	}

	s = &Service{*detail, nil}

	return nil
}
