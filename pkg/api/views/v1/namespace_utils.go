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

package v1

import (
	"encoding/json"
	"fmt"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type NamespaceView struct{}

func (nv *NamespaceView) New(obj *types.Namespace) *Namespace {
	n := Namespace{}
	n.ID = obj.Meta.Name
	n.Meta = n.ToMeta(obj.Meta)
	n.Env = n.ToEnv(obj.Env)
	n.Resources = n.ToResources(obj.Resources)
	n.Quotas = n.ToQuotas(obj.Quotas)
	return &n
}

func (r *Namespace) ToMeta(obj types.NamespaceMeta) NamespaceMeta {
	meta := NamespaceMeta{}
	meta.Name = obj.Name
	meta.Description = obj.Description
	meta.SelfLink = obj.SelfLink
	meta.Endpoint = obj.Endpoint
	meta.Created = obj.Created
	meta.Updated = obj.Updated
	meta.Labels = make(map[string]string, 0)
	if len(obj.Labels) != 0 {
		meta.Labels = obj.Meta.Labels
	}

	return meta
}

func (r *Namespace) ToEnv(obj types.NamespaceEnvs) NamespaceEnvs {
	envs := make(NamespaceEnvs, 0)
	for _, env := range obj {
		envs = append(envs, fmt.Sprintf("%s=%s", env.Name, env.Value))
	}
	return envs
}

func (r *Namespace) ToResources(obj types.NamespaceResources) NamespaceResources {
	return NamespaceResources{
		RAM:    obj.RAM,
		Routes: obj.Routes,
	}
}

func (r *Namespace) ToQuotas(obj types.NamespaceQuotas) NamespaceQuotas {
	return NamespaceQuotas{
		Disabled: obj.Disabled,
		RAM:      obj.RAM,
		Routes:   obj.Routes,
	}
}

func (p *Namespace) ToJson() ([]byte, error) {
	return json.Marshal(p)
}

func (nv NamespaceView) NewList(obj map[string]*types.Namespace) *NamespaceList {
	if obj == nil {
		return nil
	}

	n := make(NamespaceList, 0)
	for _, v := range obj {
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
