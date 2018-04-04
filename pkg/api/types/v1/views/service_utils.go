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
	"strings"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type ServiceView struct{}

// ***************************************************
// SERVICE INFO MODEL
// ***************************************************

func (sv *ServiceView) New(srv *types.Service, d map[string]*types.Deployment, p map[string]*types.Pod) *Service {
	s := new(Service)
	s.Meta = s.ToMeta(srv.Meta)
	s.Status = s.ToStatus(srv.Status)
	s.Spec = s.ToSpec(srv.Spec)
	s.Deployments = s.ToDeployments(d, p)
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

func (sv *Service) ToSource(obj types.ServiceSources) ServiceSources {
	s := ServiceSources{}
	s.Repo = new(ServiceSourcesRepo)
	s.Image = new(ServiceSourcesImage)

	if obj.Repo.ID == "" || obj.Repo.Build == "" {
		s.Image = new(ServiceSourcesImage)
		s.Image.Namespace = obj.Image.Namespace
		s.Image.Tag = obj.Image.Tag
		s.Image.Hash = obj.Image.Hash
		return s
	}

	return s
}

func (sv *Service) ToSpec(obj types.ServiceSpec) ServiceSpec {

	var spec = ServiceSpec{
		Replicas: obj.Replicas,
	}

	for _, s := range obj.Template.Containers {
		if s.Role == types.ContainerRolePrimary {
			spec.Memory = s.Resources.Limits.RAM
			spec.Image = s.Image.Name
			spec.Entrypoint = strings.Join(s.Exec.Entrypoint, " ")
			spec.Command = strings.Join(s.Exec.Command, " ")
			spec.EnvVars = s.EnvVars.ToLinuxFormat()
			spec.Ports = make([]*ServiceSpecPort, 0)
			for _, port := range s.Ports {
				p := new(ServiceSpecPort)
				p.Container = port.ContainerPort
				p.Protocol = port.Protocol
				spec.Ports = append(spec.Ports, p)
			}
			break
		}
	}

	return spec
}

func (sv *Service) ToDeployments(obj map[string]*types.Deployment, pods map[string]*types.Pod) DeploymentMap {
	deployments := make(DeploymentMap, 0)
	for _, d := range obj {
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

func (sv *ServiceView) NewList(obj map[string]*types.Service, d map[string]*types.Deployment) *ServiceList {
	if obj == nil {
		return nil
	}

	s := make(ServiceList, 0)
	slv := ServiceView{}
	for _, v := range obj {
		s = append(s, slv.New(v, d, nil))
	}
	return &s
}

func (sv *ServiceList) ToJson() ([]byte, error) {
	return json.Marshal(sv)
}
