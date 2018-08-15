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

package types

import (
	"context"
	"fmt"
)

// swagger:ignore
// swagger:model types_node_map
type NodeMap struct {
	Runtime
	Items map[string]*Node
}

// swagger:ignore
// swagger:model types_node_list
type NodeList struct {
	Runtime
	Items []*Node
}

// swagger:ignore
// swagger:model types_node
type Node struct {
	Runtime
	Meta   NodeMeta   `json:"meta"`
	Status NodeStatus `json:"status"`
	Spec   NodeSpec   `json:"spec"`
}

// swagger:ignore
// swagger:model types_node_meta
type NodeMeta struct {
	Meta
	NodeInfo
	Subnet   string `json:"subnet"`
	Cluster  string `json:"cluster"`
	Token    string `json:"token"`
}

func (m *NodeMeta) Set(meta *NodeUpdateMetaOptions) {
	if meta.Description != nil {
		m.Description = *meta.Description
	}

	if meta.Token != nil {
		m.Token = *meta.Token
	}

	if meta.Labels != nil {
		m.Labels = meta.Labels
	}

	if meta.Hostname != nil {
		m.Hostname = *meta.Hostname
	}
	if meta.Architecture != nil {
		m.Architecture = *meta.Architecture
	}
	if meta.OSName != nil {
		m.OSName = *meta.OSName
	}
	if meta.OSType != nil {
		m.OSType = *meta.OSType
	}
	if meta.ExternalIP != nil {
		m.ExternalIP = *meta.ExternalIP
	}
	if meta.InternalIP != nil {
		m.InternalIP = *meta.InternalIP
	}

}

// swagger:model types_node_info
type NodeInfo struct {
	Hostname     string `json:"hostname"`
	Architecture string `json:"architecture"`

	OSName string `json:"os_name"`
	OSType string `json:"os_type"`

	// RewriteIP - need to set true if you want to use an external ip
	ExternalIP string `json:"external_ip"`
	InternalIP string `json:"internal_ip"`
}

// swagger:model types_node_status
type NodeStatus struct {
	// state
	State NodeStatusState `json:"state"`
	// node status online
	Online bool `json:"online"`
	// Node Capacity
	Capacity NodeResources `json:"capacity"`
	// Node Allocated
	Allocated NodeResources `json:"allocated"`
}

type NodeStatusState struct {
	CRI NodeStatusInterfaceState `json:"cri"`
	CNI NodeStatusInterfaceState `json:"cni"`
	CPI NodeStatusInterfaceState `json:"cpi"`
	CSI NodeStatusInterfaceState `json:"csi"`
}

type NodeStatusInterfaceState struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	State   string `json:"state"`
	Message string `json:"message"`
}

// swagger:ignore
// swagger:model types_node_spec
type NodeSpec struct {
	Network   map[string]SubnetSpec   `json:"network"`
	Pods      map[string]PodSpec      `json:"pods"`
	Volumes   map[string]VolumeSpec   `json:"volumes"`
	Endpoints map[string]EndpointSpec `json:"endpoints"`
}

// swagger:model types_node_resources
type NodeResources struct {
	// Node total containers
	Containers int `json:"containers"`
	// Node total pods
	Pods int `json:"pods"`
	// Node total memory
	Memory int64 `json:"memory"`
	// Node total cpu
	Cpu int `json:"cpu"`
	// Node storage
	Storage int `json:"storage"`
}

// swagger:ignore
// swagger:model types_node_role
type NodeRole struct {
	Router  NodeRoleRouter `json:"router"`
	Builder bool           `json:"builder"`
}

// swagger:ignore
// swagger:model types_node_role_router
type NodeRoleRouter struct {
	ExternalIP string `json:"external_ip"`
	Enabled    bool   `json:"enabled"`
}

// swagger:ignore
// swagger:model types_node_task
type NodeTask struct {
	Cancel context.CancelFunc
}

// swagger:ignore
// swagger:model types_node_meta_create
type NodeCreateMetaOptions struct {
	MetaCreateOptions
	Subnet   string `json:"subnet"`
	Token    string `json:"token"`
	Region   string `json:"region"`
	Provider string `json:"provider"`
}

// swagger:model types_node_meta_update
type NodeUpdateMetaOptions struct {
	MetaUpdateOptions
	NodeUpdateInfoOptions
	Token    *string `json:"token"`
	Region   *string `json:"region"`
	Provider *string `json:"provider"`
}

type NodeUpdateInfoOptions struct {
	Hostname     *string `json:"hostname"`
	Architecture *string `json:"architecture"`
	OSName       *string `json:"os_name"`
	OSType       *string `json:"os_type"`
	ExternalIP   *string `json:"external_ip"`
	InternalIP   *string `json:"internal_ip"`
}

func (o *NodeUpdateInfoOptions) Set(i NodeInfo) {
	o.Hostname = &i.Hostname
	o.Architecture = &i.Architecture
	o.OSName = &i.OSName
	o.OSType = &i.OSType
	o.ExternalIP = &i.ExternalIP
	o.InternalIP = &i.InternalIP
}

// swagger:ignore
// swagger:model types_node_create
type NodeCreateOptions struct {
	Meta   NodeCreateMetaOptions `json:"meta",yaml:"meta"`
	Info   NodeInfo              `json:"info",yaml:"info"`
	Status NodeStatus            `json:"status",yaml:"status"`
}

func (n *Node) SelfLink() string {
	if n.Meta.SelfLink == "" {
		n.Meta.SelfLink = fmt.Sprintf("%s", n.Meta.Name)
	}
	return n.Meta.SelfLink
}

func NewNodeList() *NodeList {
	dm := new(NodeList)
	dm.Items = make([]*Node, 0)
	return dm
}

func NewNodeMap() *NodeMap {
	dm := new(NodeMap)
	dm.Items = make(map[string]*Node)
	return dm
}
