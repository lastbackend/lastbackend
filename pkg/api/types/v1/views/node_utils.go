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
	nm := NodeMeta{
		Name:     meta.Name,
		SelfLink: meta.SelfLink,
		Created:  meta.Created,
		Updated:  meta.Updated,
	}
	nm.NodeInfo = nv.ToNodeInfo(meta.NodeInfo)

	return nm
}

func (nv *NodeView) ToNodeInfo(info types.NodeInfo) NodeInfo {
	ni := NodeInfo{
		Hostname:     info.Hostname,
		OSType:       info.OSType,
		OSName:       info.OSName,
		Architecture: info.Architecture,
	}
	ni.IP.External = info.ExternalIP
	ni.IP.Internal = info.InternalIP

	return ni
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

func (nv *NodeView) ToNodeSpec(status types.NodeSpec) NodeSpec {
	spec := NodeSpec{}

	return spec
}

func (obj *Node) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *NodeManifest) Decode() *types.NodeManifest {

	manifest := types.NodeManifest{
		Network:   make(map[string]*types.SubnetManifest, 0),
		Pods:      make(map[string]*types.PodManifest, 0),
		Volumes:   make(map[string]*types.VolumeManifest, 0),
		Endpoints: make(map[string]*types.EndpointManifest, 0),
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
		Network:   make(map[string]*types.SubnetManifest, 0),
		Pods:      make(map[string]*types.PodManifest, 0),
		Volumes:   make(map[string]*types.VolumeManifest, 0),
		Endpoints: make(map[string]*types.EndpointManifest, 0),
	}

	if obj == nil {
		return nil
	}

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
