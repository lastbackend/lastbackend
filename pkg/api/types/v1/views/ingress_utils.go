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
)

type IngressView struct{}

func (nv *IngressView) New(obj *types.Ingress) *Ingress {
	n := Ingress{}
	n.Meta = nv.ToIngressMeta(obj.Meta)
	n.Status = nv.ToIngressStatus(obj.Status)
	return &n
}

func (nv *IngressView) ToIngressMeta(meta types.IngressMeta) IngressMeta {
	return IngressMeta{
		Name:        meta.Name,
		Description: meta.Description,
		Created:     meta.Created,
		Updated:     meta.Updated,
	}
}

func (nv *IngressView) ToIngressStatus(status types.IngressStatus) IngressStatus {
	return IngressStatus{
		Ready: status.Ready,
	}
}

func (obj *Ingress) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (nv *IngressView) NewList(obj *types.IngressList) *IngressList {
	if obj == nil {
		return nil
	}
	ingresses := make(IngressList, 0)
	for _, v := range obj.Items {
		nn := nv.New(v)
		ingresses[nn.Meta.Name] = nn
	}

	return &ingresses
}

func (nv *IngressView) NewSpec(obj *types.IngressSpec) *IngressSpec {

	spec := IngressSpec{}

	if obj == nil {
		return nil
	}

	spec.Routes = obj.Routes

	return &spec
}

func (obj *IngressSpec) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *IngressSpec) Decode() *types.IngressSpec {

	spec := types.IngressSpec{
		Routes: make(map[string]types.RouteSpec, 0),
	}

	for i, s := range obj.Routes {
		spec.Routes[i] = s
	}

	return &spec
}

func (obj *IngressList) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
