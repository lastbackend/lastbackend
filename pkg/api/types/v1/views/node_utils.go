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

type NodeView struct{}

func (nv *NodeView) New(obj *types.Node) *Node {
	n := Node{}
	n.Meta = nv.ToNodeMeta(obj.Meta)
	n.Status = nv.ToNodeStatus(obj.Status)
	n.Spec = nv.ToNodeSpec(obj.Spec)
	return &n
}

func (nv *NodeView) ToNodeMeta(meta types.NodeMeta) NodeMeta {
	nm := NodeMeta{}
	nm.Name = meta.Name
	nm.SelfLink = meta.SelfLink
	nm.Created = meta.Created
	nm.Updated = meta.Updated
	nm.Hostname = meta.Hostname
	nm.OSName = meta.OSName
	nm.OSType = meta.OSType
	nm.Architecture = meta.Architecture
	nm.IP.External = meta.ExternalIP
	nm.IP.Internal = meta.InternalIP
	nm.CIDR = meta.CIDR
	return nm
}

func (nv *NodeView) ToNodeStatus(status types.NodeStatus) NodeStatus {
	ns := NodeStatus{}

	ns.Online = status.Online

	ns.Capacity.Containers = status.Capacity.Containers
	ns.Capacity.Pods = status.Capacity.Pods
	ns.Capacity.Memory = status.Capacity.Memory
	ns.Capacity.Cpu = status.Capacity.Cpu
	ns.Capacity.Storage = status.Capacity.Storage

	ns.Allocated.Containers = status.Allocated.Containers
	ns.Allocated.Pods = status.Allocated.Pods
	ns.Allocated.Memory = status.Allocated.Memory
	ns.Allocated.Cpu = status.Allocated.Cpu
	ns.Allocated.Storage = status.Allocated.Storage

	ns.State.CNI.Type = status.State.CNI.Type
	ns.State.CNI.State = status.State.CNI.State
	ns.State.CNI.Version = status.State.CNI.Version
	ns.State.CNI.Message = status.State.CNI.Message

	ns.State.CPI.Type = status.State.CPI.Type
	ns.State.CPI.State = status.State.CPI.State
	ns.State.CPI.Version = status.State.CPI.Version
	ns.State.CPI.Message = status.State.CPI.Message

	ns.State.CRI.Type = status.State.CRI.Type
	ns.State.CRI.State = status.State.CRI.State
	ns.State.CRI.Version = status.State.CRI.Version
	ns.State.CRI.Message = status.State.CRI.Message

	ns.State.CSI.Type = status.State.CSI.Type
	ns.State.CSI.State = status.State.CSI.State
	ns.State.CSI.Version = status.State.CSI.Version
	ns.State.CSI.Message = status.State.CSI.Message

	return ns
}

func (nv *NodeView) ToNodeSpec(spec types.NodeSpec) NodeSpec {
	ns := NodeSpec{}
	ns.Security.TLS = spec.Security.TLS
	return ns
}

func (obj *Node) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *NodeManifest) Decode() *types.NodeManifest {

	manifest := types.NodeManifest{
		Secrets:   make(map[string]*types.SecretManifest, 0),
		Configs:   make(map[string]*types.ConfigManifest, 0),
		Network:   make(map[string]*types.SubnetManifest, 0),
		Pods:      make(map[string]*types.PodManifest, 0),
		Volumes:   make(map[string]*types.VolumeManifest, 0),
		Endpoints: make(map[string]*types.EndpointManifest, 0),
	}

	manifest.Meta.Initial = obj.Meta.Initial
	manifest.Resolvers = make(map[string]*types.ResolverManifest, 0)

	for i, s := range obj.Discovery {
		manifest.Resolvers[i] = s
	}

	for i, s := range obj.Network {
		manifest.Network[i] = s
	}

	for i, s := range obj.Pods {
		manifest.Pods[i] = s
	}

	for i, s := range obj.Volumes {
		manifest.Volumes[i] = s
	}

	for i, s := range obj.Endpoints {
		manifest.Endpoints[i] = s
	}

	for i, s := range obj.Configs {
		manifest.Configs[i] = s
	}

	for i, s := range obj.Secrets {
		manifest.Secrets[i] = s
	}

	return &manifest
}

func (nv *NodeView) NewList(obj *types.NodeList) *NodeList {
	if obj == nil {
		return nil
	}
	nodes := make(NodeList, 0)
	for _, v := range obj.Items {
		nn := nv.New(v)
		nodes[nn.Meta.Name] = nn
	}

	return &nodes
}

func (nv *NodeView) NewManifest(obj *types.NodeManifest) *NodeManifest {

	manifest := NodeManifest{
		Configs:   make(map[string]*types.ConfigManifest, 0),
		Secrets:   make(map[string]*types.SecretManifest, 0),
		Network:   make(map[string]*types.SubnetManifest, 0),
		Pods:      make(map[string]*types.PodManifest, 0),
		Volumes:   make(map[string]*types.VolumeManifest, 0),
		Endpoints: make(map[string]*types.EndpointManifest, 0),
	}

	if obj == nil {
		return nil
	}

	manifest.Meta.Initial = obj.Meta.Initial
	manifest.Discovery = obj.Resolvers

	manifest.Configs = obj.Configs
	manifest.Secrets = obj.Secrets
	manifest.Network = obj.Network
	manifest.Pods = obj.Pods
	manifest.Volumes = obj.Volumes
	manifest.Endpoints = obj.Endpoints

	return &manifest
}

func (obj *NodeManifest) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *NodeList) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
