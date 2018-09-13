//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package views

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"strings"
)

type ServiceView struct{}

// ***************************************************
// SERVICE INFO MODEL
// ***************************************************

func (sv *ServiceView) New(srv *types.Service) *Service {
	s := new(Service)
	s.Meta = s.ToMeta(srv.Meta)
	s.Status = s.ToStatus(srv.Status)
	s.Spec = s.ToSpec(srv.Spec)
	return s
}

func (sv *ServiceView) NewWithDeployment(srv *types.Service, d *types.DeploymentList, p *types.PodList) *Service {
	s := new(Service)
	s.Meta = s.ToMeta(srv.Meta)
	s.Status = s.ToStatus(srv.Status)
	s.Spec = s.ToSpec(srv.Spec)
	if d != nil {
		s.Deployments = s.ToDeployments(d, p)
	}
	return s
}

func (sv *Service) ToMeta(obj types.ServiceMeta) ServiceMeta {
	return ServiceMeta{
		Name:        obj.Name,
		Description: obj.Description,
		SelfLink:    obj.SelfLink,
		Endpoint:    obj.Endpoint,
		Namespace:   obj.Namespace,
		Labels:      obj.Labels,
		Updated:     obj.Updated,
		Created:     obj.Created,
	}
}

func (sv *Service) ToStatus(obj types.ServiceStatus) ServiceStatus {
	return ServiceStatus{
		State:   obj.State,
		Message: obj.Message,
	}
}

func (sv *Service) ToSpec(obj types.ServiceSpec) ServiceSpec {

	var spec = ServiceSpec{
		Replicas: obj.Replicas,
		Template: ManifestSpecTemplate{
			Containers: make([]ManifestSpecTemplateContainer, 0),
			Volumes:    make([]ManifestSpecTemplateVolume, 0),
		},
		Selector: ManifestSpecSelector{
			Node:   obj.Selector.Node,
			Labels: obj.Selector.Labels,
		},
		Network: ManifestSpecNetwork{
			IP:    obj.Network.IP,
			Ports: obj.Network.Ports,
		},
		Strategy: ManifestSpecStrategy{
			Type: obj.Strategy.Type,
		},
	}

	for _, s := range obj.Template.Containers {

		c := ManifestSpecTemplateContainer{
			Name:       s.Name,
			Command:    strings.Join(s.Exec.Command, " "),
			Workdir:    s.Exec.Workdir,
			Args:       s.Exec.Args,
			Entrypoint: strings.Join(s.Exec.Entrypoint, " "),
		}

		for _, env := range s.EnvVars {
			c.Env = append(c.Env, ManifestSpecTemplateContainerEnv{
				Name:  env.Name,
				Value: env.Value,
				From: ManifestSpecTemplateContainerEnvSecret{
					Name: env.Secret.Name,
					Key:  env.Secret.Key,
				},
			})
		}

		c.Image.Name = s.Image.Name
		c.Image.Secret = s.Image.Secret

		for _, volume := range s.Volumes {
			c.Volumes = append(c.Volumes, ManifestSpecTemplateContainerVolume{
				Name: volume.Name,
				Mode: volume.Mode,
				Path: volume.Path,
			})
		}

		c.Resources.Limits.RAM = s.Resources.Limits.RAM
		c.Resources.Limits.CPU = s.Resources.Limits.CPU
		c.Resources.Request.RAM = s.Resources.Request.RAM
		c.Resources.Request.CPU = s.Resources.Request.CPU

		spec.Template.Containers = append(spec.Template.Containers, c)
	}

	for _, s := range obj.Template.Volumes {
		v := ManifestSpecTemplateVolume{
			Name: s.Name,
			Type: s.Type,
			From: ManifestSpecTemplateSecretVolume{
				Name:  s.Secret.Name,
				Files: s.Secret.Files,
			},
		}

		spec.Template.Volumes = append(spec.Template.Volumes, v)
	}

	return spec
}

func (sv *Service) ToDeployments(obj *types.DeploymentList, pods *types.PodList) DeploymentMap {
	deployments := make(DeploymentMap, 0)
	for _, d := range obj.Items {
		if d.Meta.Service == sv.Meta.Name {
			dv := new(DeploymentView)
			dp := dv.New(d, pods)
			deployments[dp.Meta.Name] = dp
		}
	}
	return deployments
}

func (sv *Service) ToJson() ([]byte, error) {
	return json.Marshal(sv)
}

func (sv *ServiceView) NewList(obj *types.ServiceList, d *types.DeploymentList) *ServiceList {
	if obj == nil {
		return nil
	}

	s := make(ServiceList, 0)
	slv := ServiceView{}
	for _, v := range obj.Items {
		s = append(s, slv.NewWithDeployment(v, d, nil))
	}
	return &s
}

func (sv *ServiceList) ToJson() ([]byte, error) {
	return json.Marshal(sv)
}
