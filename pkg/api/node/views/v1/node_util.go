//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
)

func New(obj *types.Node) *Node {
	n := Node{}
	n.Alive = obj.Alive
	n.Meta = ToNodeMeta(obj.Meta)
	n.State = ToNodeState(obj.State)
	return &n
}

func ToNodeMeta(meta types.NodeMeta) NodeMeta {
	m := NodeMeta{
		Hostname:     meta.Hostname,
		OSType:       meta.OSType,
		OSName:       meta.OSName,
		Architecture: meta.Architecture,
		Created:      meta.Created,
		Updated:      meta.Updated,
	}
	return m
}

func ToNodeState(state types.NodeState) NodeState {
	s := NodeState{
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
	return s
}

func (obj *Node) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func NewList(obj types.NodeList) *NodeList {
	n := NodeList{}
	if obj == nil {
		return nil
	}
	for _, v := range obj {
		n = append(n, New(v))
	}
	return &n
}

func (obj *NodeList) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}
