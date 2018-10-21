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

type DiscoveryView struct{}

func (nv *DiscoveryView) New(obj *types.Discovery) *Discovery {
	n := Discovery{}
	n.Meta = nv.ToDiscoveryMeta(obj.Meta)
	n.Status = nv.ToDiscoveryStatus(obj.Status)
	return &n
}

func (nv *DiscoveryView) ToDiscoveryMeta(meta types.DiscoveryMeta) DiscoveryMeta {
	m := DiscoveryMeta{}
	m.Name = meta.Name
	m.Description = meta.Description
	m.Created = meta.Created
	m.Updated = meta.Updated
	return m
}

func (nv *DiscoveryView) ToDiscoveryStatus(status types.DiscoveryStatus) DiscoveryStatus {
	return DiscoveryStatus{
		Ready: status.Ready,
	}
}

func (obj *Discovery) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (nv *DiscoveryView) NewList(obj *types.DiscoveryList) *DiscoveryList {
	if obj == nil {
		return nil
	}
	ingresses := make(DiscoveryList, 0)
	for _, v := range obj.Items {
		nn := nv.New(v)
		ingresses[nn.Meta.Name] = nn
	}

	return &ingresses
}

func (obj *DiscoveryList) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (nv *DiscoveryView) NewManifest(obj *types.DiscoveryManifest) *DiscoveryManifest {

	manifest := DiscoveryManifest{
		Subnets: make(map[string]*types.SubnetManifest, 0),
	}

	if obj == nil {
		return nil
	}

	manifest.Meta.Initial = obj.Meta.Initial
	manifest.Subnets = obj.Network

	return &manifest
}

func (obj *DiscoveryManifest) Decode() *types.DiscoveryManifest {

	manifest := types.DiscoveryManifest{
		Network: make(map[string]*types.SubnetManifest, 0),
	}

	manifest.Meta.Initial = obj.Meta.Initial

	for i, e := range obj.Subnets {
		manifest.Network[i] = e
	}

	return &manifest
}

func (obj *DiscoveryManifest) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
