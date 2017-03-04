package service

import (
	"github.com/lastbackend/lastbackend/pkg/service/resource/container"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type ServiceConfig struct {
	Scale      int32
	Containers *[]container.Container
}

func (s *ServiceConfig) patch(dp *v1beta1.Deployment) {

	if s.Scale < 1 {
		s.Scale = 1
	}

	dp.Spec.Replicas = &s.Scale

	if s.Containers != nil {
		dp.Spec.Template.Spec.Containers = []v1.Container{}
		for _, item := range *s.Containers {
			c := v1.Container{}

			c.Name = item.Name
			c.Image = item.Image
			c.WorkingDir = item.WorkingDir
			c.Command = item.Command
			c.Args = item.Args

			for _, val := range item.Ports {
				c.Ports = append(c.Ports, v1.ContainerPort{
					Name:          val.Name,
					ContainerPort: val.ContainerPort,
					Protocol:      v1.Protocol(val.Protocol),
				})

				for _, val := range item.Env {
					c.Env = append(c.Env, v1.EnvVar{
						Name:  val.Name,
						Value: val.Value,
					})
				}

				dp.Spec.Template.Spec.Containers = append(dp.Spec.Template.Spec.Containers, c)
			}
		}
	}
}
