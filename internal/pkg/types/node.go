//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
)

// swagger:ignore
// swagger:model types_node_map
type NodeMap struct {
	System
	Items map[string]*Node
}

// swagger:ignore
// swagger:model types_node_list
type NodeList struct {
	System
	Items []*Node
}

// swagger:ignore
// swagger:model types_node
type Node struct {
	System
	Meta   NodeMeta   `json:"meta"`
	Status NodeStatus `json:"status"`
	Spec   NodeSpec   `json:"spec"`
}

// swagger:ignore
// swagger:model types_node_meta
type NodeMeta struct {
	Meta
	NodeInfo
	SelfLink NodeSelfLink `json:"self_link"`
	Subnet   string       `json:"subnet"`
	Cluster  string       `json:"cluster"`
}

func (m *NodeMeta) Set(meta *NodeUpdateMetaOptions) {
	if meta.Description != nil {
		m.Description = *meta.Description
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
	if meta.CIDR != nil {
		m.CIDR = *meta.CIDR
	}

}

// swagger:model types_node_info
type NodeInfo struct {
	Version      string `json:"version"`
	Hostname     string `json:"hostname"`
	Architecture string `json:"architecture"`

	OSName string `json:"os_name"`
	OSType string `json:"os_type"`

	// RewriteIP - need to set true if you want to use an external ip
	ExternalIP string `json:"external_ip"`
	InternalIP string `json:"internal_ip"`
	CIDR       string `json:"cidr"`
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
	Security NodeSecurity `json:"security"`
}

type NodeSecurity struct {
	TLS bool     `json:"tls"`
	SSL *NodeSSL `json:"ssl"`
}

type NodeSSL struct {
	CA   []byte `json:"ca"`
	Cert []byte `json:"cert"`
	Key  []byte `json:"key"`
}

// swagger:model types_node_resources
type NodeResources struct {
	// Node total containers
	Containers int `json:"containers"`
	// Node total pods
	Pods int `json:"pods"`
	// Node total memory
	RAM int64 `json:"ram"`
	// Node total cpu
	CPU int64 `json:"cpu"`
	// Node storage
	Storage int64 `json:"storage"`
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
}

type NodeUpdateInfoOptions struct {
	Hostname     *string `json:"hostname"`
	Architecture *string `json:"architecture"`
	OSName       *string `json:"os_name"`
	OSType       *string `json:"os_type"`
	ExternalIP   *string `json:"external_ip"`
	InternalIP   *string `json:"internal_ip"`
	CIDR         *string `json:"cidr"`
}

func (o *NodeUpdateInfoOptions) Set(i NodeInfo) {
	o.Hostname = &i.Hostname
	o.Architecture = &i.Architecture
	o.OSName = &i.OSName
	o.OSType = &i.OSType
	o.ExternalIP = &i.ExternalIP
	o.InternalIP = &i.InternalIP
	o.CIDR = &i.CIDR
}

// swagger:ignore
// swagger:model types_node_create
type NodeCreateOptions struct {
	Meta     NodeCreateMetaOptions `json:"meta", yaml:"meta"`
	Info     NodeInfo              `json:"info", yaml:"info"`
	Status   NodeStatus            `json:"status", yaml:"status"`
	Security NodeSecurity          `json:"security", yaml:"security"`
}

func (n *Node) SelfLink() *NodeSelfLink {
	return &n.Meta.SelfLink
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
