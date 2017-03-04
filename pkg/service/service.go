package service

import (
	"io"
	"time"
	"github.com/lastbackend/lastbackend/libs/interface/k8s"
	"github.com/lastbackend/lastbackend/pkg/service/resource/deployment"
	"github.com/unloop/gopipe"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type Service struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Scale     int32             `json:"scale"`
	Template  struct {
		ContainerList []Container `json:"containers"`
	} `json:"template"`
	PodList []Pod `json:"pods"`
}

type Pod struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Status        string            `json:"status"`
	RestartCount  int32             `json:"restart_count"`
	RestartPolicy string            `json:"restart_policy"`
	Labels        map[string]string `json:"labels"`
	ContainerList []Container       `json:"containers"`
	StartTime     time.Time         `json:"uptime"`
}

type Container struct {
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	WorkingDir string   `json:"workdir"`
	Status     string   `json:"status,omitempty"`
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

func Get(client k8s.IK8S, namespace, name string) (*Service, error) {

	var err error

	dp, err := deployment.Get(client, namespace, name)
	if err != nil {
		return nil, err
	}

	s := convert_deployment_to_service(dp)

	return s, nil
}

func List(client k8s.IK8S, namespace string) (map[string]*Service, error) {

	var (
		err         error
		serviceList = make(map[string]*Service)
	)

	deploymentList, err := deployment.List(client, namespace)
	if err != nil {
		return nil, err
	}

	for _, dp := range deploymentList {
		service := convert_deployment_to_service(&dp)
		serviceList[dp.ObjectMeta.Name] = service
	}

	return serviceList, nil
}

func Update(client k8s.IK8S, namespace, name string, config *ServiceConfig) error {

	var err error

	dp, err := client.ExtensionsV1beta1().Deployments(namespace).Get(name)
	if err != nil {
		return err
	}

	config.patch(dp)

	return  deployment.Update(client, namespace, dp)
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

func Logs(client k8s.IK8S, namespace string, opts *ServiceLogsOption, close chan bool) error {

	var (
		err    error
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

	req := client.CoreV1().Pods(namespace).GetLogs(opts.Pod, &option)

	readCloser, err := req.Stream()
	if err != nil {
		return err
	}
	defer readCloser.Close()

	go s.Pipe(&readCloser)

	<-close

	s.Close()

	return nil
}

func Deploy(client k8s.IK8S, namespace string, config *v1beta1.Deployment) (*Service, error) {

	var err error

	_, err = client.ExtensionsV1beta1().Deployments(namespace).Create(config)
	if err != nil {
		return nil, err
	}

	dp, err := deployment.Get(client, namespace, config.Name)
	if err != nil {
		return nil, err
	}

	service := convert_deployment_to_service(dp)

	return service, nil
}

func UpdateImage(client k8s.IK8S, namespace, name string) error {

	dp, err := deployment.Get(client, namespace, name)
	if err != nil {
		return err
	}

	for index := range dp.PodList.Pods {
		if err := dp.PodList.Pods[index].Remove(client); err != nil {
			return err
		}
	}

	return nil
}

func convert_deployment_to_service(dp *deployment.Deployment) *Service {

	var s = new(Service)

	s.Name = dp.ObjectMeta.Name
	s.Namespace = dp.ObjectMeta.Namespace
	s.Labels = dp.ObjectMeta.Labels
	s.Scale = dp.Spec.Replicas
	s.Template.ContainerList = []Container{}
	s.PodList = []Pod{}

	for _, container := range dp.Spec.Template.Spec.Containers {
		c := Container{}

		c.Name = container.Name
		c.Image = container.Image
		c.WorkingDir = container.WorkingDir
		c.Command = container.Command
		c.Args = container.Args
		c.PortList = []Port{}
		c.EnvList = []Env{}
		c.VolumeList = []Volume{}

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
		p.Status = string(pod.Status.PodPhase)
		p.Labels = pod.ObjectMeta.Labels
		p.RestartCount = pod.RestartCount
		p.RestartPolicy = string(pod.RestartPolicy)
		p.StartTime = pod.StartTime
		p.ContainerList = []Container{}

		for _, container := range pod.ContainerList.Containers {
			c := Container{}

			c.Name = container.Name
			c.Image = container.Image
			c.Status = container.Status
			c.WorkingDir = container.WorkingDir
			c.Command = container.Command
			c.Args = container.Args
			c.PortList = []Port{}
			c.EnvList = []Env{}
			c.VolumeList = []Volume{}

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
