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
	n.Info = nv.ToNodeInfo(obj.Info)
	n.Status.Online = obj.Online
	return &n
}

func (nv *NodeView) ToNodeMeta(meta types.NodeMeta) NodeMeta {
	return NodeMeta{
		Name:     meta.Name,
		SelfLink: meta.SelfLink,
		Created:  meta.Created,
		Updated:  meta.Updated,
	}
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
	return NodeStatus{
		Capacity: NodeResources{
			Containers: status.Capacity.Containers,
			Pods:       status.Capacity.Pods,
			Memory:     status.Capacity.Memory,
			Cpu:        status.Capacity.Cpu,
			Storage:    status.Capacity.Storage,
		},
		Allocated: NodeResources{
			Containers: status.Allocated.Containers,
			Pods:       status.Allocated.Pods,
			Memory:     status.Allocated.Memory,
			Cpu:        status.Allocated.Cpu,
			Storage:    status.Allocated.Storage,
		},
	}
}

func (obj *Node) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *NodeManifest) Decode() *types.NodeManifest {

	manifest := types.NodeManifest{
		Network:   make(map[string]types.NetworkManifest, 0),
		Pods:      make(map[string]types.PodManifest, 0),
		Volumes:   make(map[string]types.VolumeManifest, 0),
		Endpoints: make(map[string]types.EndpointManifest, 0),
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

func (nv *NodeView) NewList(obj map[string]*types.Node) *NodeList {
	if obj == nil {
		return nil
	}
	nodes := make(NodeList, 0)
	for _, v := range obj {
		nn := nv.New(v)
		nodes[nn.Meta.Name] = nn
	}

	return &nodes
}

func (nv *NodeView) NewManifest(obj *types.NodeManifest) *NodeManifest {

	manifest := NodeManifest{}

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
