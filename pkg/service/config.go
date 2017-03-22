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
