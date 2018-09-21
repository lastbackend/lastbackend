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
	m := IngressMeta{}
	m.Name = meta.Name
	m.Description = meta.Description
	m.Created = meta.Created
	m.Updated = meta.Updated
	return m
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

func (obj *IngressList) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (nv *IngressView) NewManifest(obj *types.IngressManifest) *IngressManifest {

	manifest := IngressManifest{
		Endpoints: make(map[string]*types.EndpointManifest, 0),
		Routes:    make(map[string]*types.RouteManifest, 0),
		Subnets:   make(map[string]*types.SubnetManifest, 0),
	}

	if obj == nil {
		return nil
	}

	manifest.Meta.Initial = obj.Meta.Initial
	manifest.Meta.Discovery = obj.Meta.Discovery
	manifest.Endpoints = obj.Endpoints
	manifest.Routes = obj.Routes
	manifest.Subnets = obj.Network

	return &manifest
}

func (obj *IngressManifest) Decode() *types.IngressManifest {

	manifest := types.IngressManifest{
		Routes: make(map[string]*types.RouteManifest, 0),
		Endpoints: make(map[string]*types.EndpointManifest, 0),
		Network: make(map[string]*types.SubnetManifest, 0),
	}

	manifest.Meta.Initial = obj.Meta.Initial
	manifest.Meta.Discovery = obj.Meta.Discovery

	for i, r := range obj.Routes {
		manifest.Routes[i] = r
	}

	for i, e := range obj.Endpoints {
		manifest.Endpoints[i] = e
	}

	for i, e := range obj.Subnets {
		manifest.Network[i] = e
	}

	return &manifest
}

func (obj *IngressManifest) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
