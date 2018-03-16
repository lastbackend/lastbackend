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

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type RouteView struct{}

func (rv *RouteView) New(obj *types.Route) *Route {
	r := Route{}
	r.Meta = r.ToMeta(obj.Meta)
	r.Spec  = r.ToSpec(obj.Spec)
	r.State = r.ToState(obj.State)
	return &r
}

func (p *Route) ToJson() ([]byte, error) {
	return json.Marshal(p)
}

func (r *Route) ToMeta(obj types.RouteMeta) RouteMeta {
	meta := RouteMeta{}
	meta.Domain = obj.Domain
	meta.Namespace = obj.Namespace
	meta.Security = obj.Security
	meta.Updated = obj.Updated
	meta.Created = obj.Created

	return meta
}

func (r *Route) ToSpec(obj types.RouteSpec) RouteSpec {
	spec := RouteSpec{}
	return spec
}

func (r *Route) ToState(obj types.RouteState) RouteState {
	state := RouteState{}
	state.Destroy = obj.Destroy
	state.Provision = obj.Provision
	return state
}

func (rv RouteView) NewList(obj []*types.Route) *RouteList {
	if obj == nil {
		return nil
	}

	n := make(RouteList, 0)
	for _, v := range obj {
		n = append(n, rv.New(v))
	}
	return &n
}

func (n *RouteList) ToJson() ([]byte, error) {
	if n == nil {
		n = &RouteList{}
	}
	return json.Marshal(n)
}