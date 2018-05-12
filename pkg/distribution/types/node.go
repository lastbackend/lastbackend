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
// swagger:model types_node_status_map
type NodeMapStatus map[string]*NodeStatus

// swagger:ignore
// swagger:model types_node_map
type NodeMap map[string]*Node

// swagger:ignore
// swagger:model types_node_list
type NodeList []*Node

// swagger:ignore
// swagger:model types_node
type Node struct {
	Meta    NodeMeta    `json:"meta"`
	Info    NodeInfo    `json:"info"`
	Status  NodeStatus  `json:"status"`
	Spec    NodeSpec    `json:"spec"`
	Roles   NodeRole    `json:"roles"`
	Network NetworkSpec `json:"network"`
	Online  bool        `json:"online"`
}

// swagger:ignore
// swagger:model types_node_meta
type NodeMeta struct {
	Meta
	Cluster  string `json:"cluster"`
	Token    string `json:"token"`
	Region   string `json:"region"`
	Provider string `json:"provider"`
}

func (m *NodeMeta) Set(meta *NodeUpdateMetaOptions) {
	if meta.Description != nil {
		m.Description = *meta.Description
	}

	if meta.Token != nil {
		m.Token = *meta.Token
	}

	if meta.Region != nil {
		m.Region = *meta.Region
	}

	if meta.Provider != nil {
		m.Provider = *meta.Provider
	}

	if meta.Labels != nil {
		m.Labels = meta.Labels
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
	// Node Capacity
	Capacity NodeResources `json:"capacity"`
	// Node Allocated
	Allocated NodeResources `json:"allocated"`
}

// swagger:ignore
// swagger:model types_node_spec
type NodeSpec struct {
	Network map[string]NetworkSpec `json:"network"`
	Pods    map[string]PodSpec     `json:"pods"`
	Volumes map[string]VolumeSpec  `json:"volumes"`
}

// swagger:ignore
// swagger:model types_node_namespace
type NodeNamespace struct {
	Meta NamespaceMeta     `json:"meta",yaml:"meta"`
	Spec NodeNamespaceSpec `json:"spec",yaml:"spec"`
}

// swagger:ignore
// swagger:model types_node_namespace_spec
type NodeNamespaceSpec struct {
	Routes  []*Route  `json:"routes",yaml:"routes"`
	Pods    []*Pod    `json:"pods",yaml:"pods"`
	Volumes []*Volume `json:"volumes",yaml:"volumes"`
	Secrets []*Secret `json:"secrets",yaml:"secrets"`
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
	Token    string `json:"token"`
	Region   string `json:"region"`
	Provider string `json:"provider"`
}

// swagger:model types_node_meta_update
type NodeUpdateMetaOptions struct {
	MetaUpdateOptions
	Token    *string `json:"token"`
	Region   *string `json:"region"`
	Provider *string `json:"provider"`
}

// swagger:ignore
// swagger:model types_node_create
type NodeCreateOptions struct {
	Meta    NodeCreateMetaOptions `json:"meta",yaml:"meta"`
	Info    NodeInfo              `json:"info",yaml:"info"`
	Status  NodeStatus            `json:"status",yaml:"status"`
	Network NetworkSpec           `json:"network"`
}

func (n *Node) SelfLink() string {
	if n.Meta.SelfLink == "" {
		n.Meta.SelfLink = fmt.Sprintf("%s", n.Meta.Name)
	}
	return n.Meta.SelfLink
}
