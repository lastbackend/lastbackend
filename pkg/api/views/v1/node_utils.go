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

type NodeView struct{}

func (nv *NodeView) New(obj *types.Node) *Node {
	n := Node{}
	n.Meta = nv.ToNodeMeta(obj.Meta)
	n.State = nv.ToNodeState(obj.State)
	n.Info = nv.ToNodeInfo(obj.Info)
	return &n
}

func (nv *NodeView) ToNodeMeta(meta types.NodeMeta) NodeMeta {
	return NodeMeta{
		ID:          meta.Name,
		Description: meta.Description,
		Created:     meta.Created,
		Updated:     meta.Updated,
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

func (nv *NodeView) ToNodeState(state types.NodeState) NodeState {
	return NodeState{
		Capacity: NodeResources{
			Containers: state.Capacity.Containers,
			Pods:       state.Capacity.Pods,
			Memory:     state.Capacity.Memory,
			Cpu:        state.Capacity.Cpu,
			Storage:    state.Capacity.Storage,
		},
		Allocated: NodeResources{
			Containers: state.Allocated.Containers,
			Pods:       state.Allocated.Pods,
			Memory:     state.Allocated.Memory,
			Cpu:        state.Allocated.Cpu,
			Storage:    state.Allocated.Storage,
		},
	}
}

func (obj *Node) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (nv *NodeView) NewList(obj map[string]*types.Node) *NodeList {
	if obj == nil {
		return nil
	}
	nodes := make(NodeList, 0)
	for _, v := range obj {
		nn := nv.New(v)
		nodes[nn.Meta.ID] = nn
	}

	return &nodes
}

func (nv *NodeView) NewSpec(obj *types.Node) *NodeSpec {

	spec := NodeSpec {}

	if obj == nil {
		return nil
	}

	return &spec
}

func (obj *NodeSpec) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *NodeList) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
