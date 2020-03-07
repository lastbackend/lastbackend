//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/internal/util/resource"

	"github.com/lastbackend/lastbackend/internal/pkg/types"
)

type NamespaceView struct{}

func (nv *NamespaceView) New(obj *types.Namespace) *Namespace {
	n := Namespace{}
	n.Meta = n.ToMeta(obj.Meta)
	n.Status = n.ToStatus(obj.Status)
	n.Spec = n.ToSpec(obj.Spec)
	return &n
}

func (nv *NamespaceView) NewApplyStatus(status struct {
	Configs  map[string]bool
	Secrets  map[string]bool
	Volumes  map[string]bool
	Services map[string]bool
	Jobs     map[string]bool
	Routes   map[string]bool
}) *NamespaceApplyStatus {
	n := NamespaceApplyStatus{}
	n.Secrets = make(map[string]bool, 0)
	n.Configs = make(map[string]bool, 0)
	n.Volumes = make(map[string]bool, 0)
	n.Services = make(map[string]bool, 0)
	n.Routes = make(map[string]bool, 0)
	n.Jobs = make(map[string]bool, 0)

	for name, status := range status.Secrets {
		n.Secrets[name] = status
	}

	for name, status := range status.Configs {
		n.Configs[name] = status
	}

	for name, status := range status.Volumes {
		n.Volumes[name] = status
	}

	for name, status := range status.Services {
		n.Services[name] = status
	}

	for name, status := range status.Routes {
		n.Routes[name] = status
	}

	for name, status := range status.Jobs {
		n.Jobs[name] = status
	}

	return &n
}

func (r *Namespace) ToMeta(obj types.NamespaceMeta) NamespaceMeta {
	meta := NamespaceMeta{}
	meta.Name = obj.Name
	meta.Description = obj.Description
	meta.SelfLink = obj.SelfLink.String()
	meta.Endpoint = obj.Endpoint
	meta.Created = obj.Created
	meta.Updated = obj.Updated
	meta.Labels = make(map[string]string, 0)

	if obj.Labels != nil {
		meta.Labels = obj.Meta.Labels
	}

	return meta
}

func (r *Namespace) ToSpec(spec types.NamespaceSpec) NamespaceSpec {
	return NamespaceSpec{
		Resources: NamespaceResources{
			Limits:  r.ToResources(spec.Resources.Limits),
			Request: r.ToResources(spec.Resources.Request),
		},
		Env: r.ToEnv(spec.Env),
		Domain: NamespaceDomain{
			Internal: spec.Domain.Internal,
			External: spec.Domain.External,
		},
	}
}

func (r *Namespace) ToStatus(status types.NamespaceStatus) NamespaceStatus {
	return NamespaceStatus{
		Resources: NamespaceStatusResources{
			Allocated: r.ToResources(status.Resources.Allocated),
		},
	}
}

func (r *Namespace) ToEnv(obj types.NamespaceEnvs) NamespaceEnvs {
	envs := make(NamespaceEnvs, 0)
	for _, env := range obj {
		envs = append(envs, fmt.Sprintf("%s=%s", env.Name, env.Value))
	}
	return envs
}

func (r *Namespace) ToResources(obj types.ResourceItem) *NamespaceResource {

	if obj.RAM == 0 || obj.CPU == 0 || obj.Storage == 0 {
		return nil
	}

	return &NamespaceResource{
		RAM:     resource.EncodeMemoryResource(obj.RAM),
		CPU:     resource.EncodeCpuResource(obj.CPU),
		Storage: resource.EncodeMemoryResource(obj.Storage),
	}
}

func (r *Namespace) ToJson() ([]byte, error) {
	return json.Marshal(r)
}

func (nv NamespaceView) NewList(obj *types.NamespaceList) *NamespaceList {
	if obj == nil {
		return nil
	}

	n := make(NamespaceList, 0)
	for _, v := range obj.Items {
		n = append(n, nv.New(v))
	}
	return &n
}

func (n *NamespaceList) ToJson() ([]byte, error) {
	if n == nil {
		n = &NamespaceList{}
	}
	return json.Marshal(n)
}

func (s *NamespaceApplyStatus) ToJson() ([]byte, error) {
	return json.Marshal(s)
}
