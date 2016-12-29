package service

import (
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
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Scale     int32             `json:"scale"`
	Template  struct {
		ContainerList []Container `json:"containers"`
	} `json:"tempalte"`
	PodList []Pod `json:"pods"`
}

type Pod struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Status        string            `json:"status"`
	Labels        map[string]string `json:"labels"`
	ContainerList []Container       `json:"containers"`
}

type Container struct {
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	WorkingDir string   `json:"workdir"`
	Command    []string `json:"command"`
	Args       []string `json:"args"`
	PortList   []Port   `json:"ports"`
	EnvList    []Env    `json:"env"`
	VolumeList []Volume `json:"volumes"`
}

type Port struct {
	Name      string `json:"name"`
	Container int32  `json:"container"`
	Protocol  string `json:"protocol"`
}

type Env struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Volume struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Readonly bool   `json:"readonly"`
}

func Get(client k8s.IK8S, namespace, name string) (*Service, *e.Err) {

	var er error

	dp, er := deployment.Get(client, namespace, name)
	if er != nil {
		return nil, e.New("service").Unknown(er)
	}

	service := convert_deployment_to_service(dp)

	return service, nil
}

func List(client k8s.IK8S, namespace string) (map[string]*Service, *e.Err) {

	var (
		er          error
		serviceList = make(map[string]*Service)
	)

	deploymentList, er := deployment.List(client, namespace)
	if er != nil {
		return nil, e.New("service").Unknown(er)
	}

	for _, dp := range deploymentList {
		service := convert_deployment_to_service(&dp)
		serviceList[dp.ObjectMeta.Name] = service
	}

	return serviceList, nil
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

func Deploy(client k8s.IK8S, namespace string, config *v1beta1.Deployment) (*Service, *e.Err) {

	var er error

	_, er = client.Extensions().Deployments(namespace).Create(config)
	if er != nil {
		return nil, e.New("service").Unknown(er)
	}

	dp, er := deployment.Get(client, namespace, config.Name)
	if er != nil {
		return nil, e.New("service").Unknown(er)
	}

	service := convert_deployment_to_service(dp)

	return service, nil
}

func convert_deployment_to_service(dp *deployment.Deployment) *Service {

	var s = new(Service)

	s.Name = dp.ObjectMeta.Name
	s.Namespace = dp.ObjectMeta.Namespace
	s.Labels = dp.ObjectMeta.Labels
	s.Scale = dp.Spec.Replicas

	for _, container := range dp.Spec.Template.Spec.Containers {
		c := Container{}
		c.Name = container.Name
		c.Image = container.Image
		c.WorkingDir = container.WorkingDir
		c.Command = container.Command
		c.Args = container.Args

		for _, port := range container.Ports {
			cp := Port{}
			cp.Name = port.Name
			cp.Protocol = string(port.Protocol)
			cp.Container = port.ContainerPort

			c.PortList = append(c.PortList, cp)
		}

		for _, env := range container.Env {
			ce := Env{}
			ce.Name = env.Name
			ce.Value = env.Value

			c.EnvList = append(c.EnvList, ce)
		}

		for _, volume := range container.VolumeMounts {
			cv := Volume{}
			cv.Name = volume.Name
			cv.Path = volume.MountPath
			cv.Readonly = volume.ReadOnly

			c.VolumeList = append(c.VolumeList, cv)
		}

		s.Template.ContainerList = append(s.Template.ContainerList, c)
	}

	for _, pod := range dp.PodList.Pods {
		p := Pod{}
		p.Name = pod.ObjectMeta.Name
		p.Namespace = pod.ObjectMeta.Namespace
		p.Status = string(pod.PodStatus.PodPhase)
		p.Labels = pod.ObjectMeta.Labels

		for _, container := range pod.ContainerList.Containers {
			c := Container{}
			c.Name = container.Name
			c.Image = container.Image
			c.WorkingDir = container.WorkingDir
			c.Command = container.Command
			c.Args = container.Args

			for _, port := range container.Ports {
				cp := Port{}
				cp.Name = port.Name
				cp.Protocol = port.Protocol
				cp.Container = port.ContainerPort

				c.PortList = append(c.PortList, cp)
			}

			for _, env := range container.Env {
				ce := Env{}
				ce.Name = env.Name
				ce.Value = env.Value

				c.EnvList = append(c.EnvList, ce)
			}

			for _, volume := range container.Volumes {
				cv := Volume{}
				cv.Name = volume.Name
				cv.Path = volume.MountPath
				cv.Readonly = volume.ReadOnly

				c.VolumeList = append(c.VolumeList, cv)
			}

			p.ContainerList = append(p.ContainerList, c)
		}

		s.PodList = append(s.PodList, p)
	}

	return s
}
