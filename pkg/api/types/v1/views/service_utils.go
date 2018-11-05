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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"strconv"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
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

	s.Deployments = make(DeploymentMap, 0)
	if d != nil {
		s.Deployments = s.ToDeployments(d, p)
	}
	return s
}

func (sv *Service) ToMeta(obj types.ServiceMeta) ServiceMeta {
	sm := ServiceMeta{
		Name:        obj.Name,
		Description: obj.Description,
		SelfLink:    obj.SelfLink,
		Endpoint:    obj.Endpoint,
		Namespace:   obj.Namespace,
		Labels:      obj.Labels,
		Updated:     obj.Updated,
		Created:     obj.Created,
	}

	sm.Labels = make(map[string]string, 0)
	if obj.Labels != nil {
		sm.Labels = obj.Labels
	}

	return sm
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

		for _, p := range s.Ports {
			c.Ports = append(c.Ports, string(p.ContainerPort))
		}

		for _, env := range s.EnvVars {
			c.Env = append(c.Env, ManifestSpecTemplateContainerEnv{
				Name:  env.Name,
				Value: env.Value,
				Secret: ManifestSpecTemplateContainerEnvSecret{
					Name: env.Secret.Name,
					Key:  env.Secret.Key,
				},
				Config: ManifestSpecTemplateContainerEnvConfig{
					Name: env.Config.Name,
					Key:  env.Config.Key,
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
			Volume: ManifestSpecTemplateVolumeClaim{
				Name:    s.Volume.Name,
				Subpath: s.Volume.Subpath,
			},
			Secret: ManifestSpecTemplateSecretVolume{
				Name:  s.Secret.Name,
				Binds: make([]ManifestSpecTemplateSecretVolumeBind, 0),
			},
			Config: ManifestSpecTemplateConfigVolume{
				Name:  s.Config.Name,
				Binds: make([]ManifestSpecTemplateConfigVolumeBind, 0),
			},
		}

		for _, b := range s.Secret.Binds {
			v.Secret.Binds = append(v.Secret.Binds, ManifestSpecTemplateSecretVolumeBind{
				Key:  b.Key,
				File: b.File,
			})
		}

		for _, b := range s.Config.Binds {
			v.Config.Binds = append(v.Config.Binds, ManifestSpecTemplateConfigVolumeBind{
				Key:  b.Key,
				File: b.File,
			})
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

func (sv Service) ToRequestManifest() *request.ServiceManifest {
	sm := new(request.ServiceManifest)

	sm.Meta.Name = &sv.Meta.Name
	sm.Meta.Description = &sv.Meta.Description

	if sv.Meta.Labels == nil {
		sm.Meta.Labels = make(map[string]string, 0)
	} else {
		sm.Meta.Labels = sv.Meta.Labels
	}

	sm.Spec.Replicas = &sv.Spec.Replicas

	sm.Spec.Selector = new(request.ManifestSpecSelector)
	sm.Spec.Selector.Node = sv.Spec.Selector.Node
	sm.Spec.Selector.Labels = sv.Spec.Selector.Labels
	if sm.Spec.Selector.Labels == nil {
		sm.Spec.Selector.Labels = make(map[string]string, 0)
	}

	sm.Spec.Strategy = new(request.ManifestSpecStrategy)
	sm.Spec.Strategy.Type = &sv.Spec.Strategy.Type

	sm.Spec.Network = new(request.ManifestSpecNetwork)
	sm.Spec.Network.IP = &sv.Spec.Network.IP
	sm.Spec.Network.Ports = make([]string, 0)

	if sv.Spec.Network.Ports != nil {
		// k - port, v - port/protocol
		for k, v := range sv.Spec.Network.Ports {
			p := strconv.Itoa(int(k))
			if p == v {
				sm.Spec.Network.Ports = append(sm.Spec.Network.Ports, v)
			} else {
				sm.Spec.Network.Ports = append(sm.Spec.Network.Ports, fmt.Sprintf("%s:%s", p, v))
			}
		}
	}

	sm.Spec.Template = new(request.ManifestSpecTemplate)
	sm.Spec.Template.Volumes = make([]request.ManifestSpecTemplateVolume, 0)

	if sv.Spec.Template.Volumes != nil {
		for _, v := range sm.Spec.Template.Volumes {

			data := request.ManifestSpecTemplateVolume{
				Name: v.Name,
				Type: v.Type,
			}

			data.Secret.Name = v.Secret.Name
			data.Secret.Binds = v.Secret.Binds
			if data.Secret.Binds == nil {
				data.Secret.Binds = make([]request.ManifestSpecTemplateSecretVolumeBind, 0)
			}

			data.Config.Name = v.Config.Name
			data.Config.Binds = v.Config.Binds
			if data.Config.Binds == nil {
				data.Config.Binds = make([]request.ManifestSpecTemplateConfigVolumeBind, 0)
			}

			sm.Spec.Template.Volumes = append(sm.Spec.Template.Volumes, data)
		}
	}

	sm.Spec.Template.Containers = make([]request.ManifestSpecTemplateContainer, 0)
	if sv.Spec.Template.Containers != nil {
		for _, v := range sv.Spec.Template.Containers {

			data := request.ManifestSpecTemplateContainer{
				Name:       v.Name,
				Command:    v.Command,
				Workdir:    v.Workdir,
				Entrypoint: v.Entrypoint,
				Ports:      v.Ports,
			}

			data.Args = v.Args
			if data.Args == nil {
				data.Args = make([]string, 0)
			}

			data.Env = make([]request.ManifestSpecTemplateContainerEnv, 0)
			if v.Env != nil {
				for _, v := range v.Env {
					item := request.ManifestSpecTemplateContainerEnv{
						Name:  v.Name,
						Value: v.Value,
						Secret: request.ManifestSpecTemplateContainerEnvSecret{
							Name: v.Secret.Name,
							Key:  v.Secret.Key,
						},
						Config: request.ManifestSpecTemplateContainerEnvConfig{
							Name: v.Config.Name,
							Key:  v.Config.Key,
						},
					}
					data.Env = append(data.Env, item)
				}
			}

			data.Volumes = make([]request.ManifestSpecTemplateContainerVolume, 0)
			if v.Volumes != nil {
				for _, v := range v.Volumes {
					item := request.ManifestSpecTemplateContainerVolume{
						Name: v.Name,
						Mode: v.Mode,
						Path: v.Path,
					}
					data.Volumes = append(data.Volumes, item)
				}
			}

			data.Image.Name = v.Image.Name
			data.Image.Secret = v.Image.Secret

			data.Resources.Request.RAM = v.Resources.Request.RAM
			data.Resources.Request.CPU = v.Resources.Request.CPU
			data.Resources.Limits.RAM = v.Resources.Limits.RAM
			data.Resources.Limits.CPU = v.Resources.Limits.CPU

			data.RestartPolicy.Policy = v.RestartPolicy.Policy
			data.RestartPolicy.Attempt = v.RestartPolicy.Attempt

			sm.Spec.Template.Containers = append(sm.Spec.Template.Containers, data)
		}
	}

	return sm
}

func (sv *ServiceView) NewList(obj *types.ServiceList, d *types.DeploymentList, pl *types.PodList) *ServiceList {
	if obj == nil {
		return nil
	}

	s := make(ServiceList, 0)
	slv := ServiceView{}
	for _, v := range obj.Items {
		s = append(s, slv.NewWithDeployment(v, d, pl))
	}
	return &s
}

func (sv *ServiceList) ToJson() ([]byte, error) {
	return json.Marshal(sv)
}
