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
)

type NodeMapState map[string]*NodeState

type NodeList []*Node

type Node struct {
	Meta    NodeMeta  `json:"meta"`
	Info    NodeInfo  `json:"info"`
	State   NodeState `json:"state"`
	Spec    NodeSpec  `json:"spec"`
	Roles   NodeRole  `json:"roles"`
	Network Subnet    `json:"network"`
	Online  bool      `json:"online"`
}

type NodeMeta struct {
	Meta
	Token    string `json:"token"`
	Region   string `json:"region"`
	Provider string `json:"provider"`
}

type NodeInfo struct {
	Hostname     string `json:"hostname"`
	Architecture string `json:"architecture"`

	OSName string `json:"os_name"`
	OSType string `json:"os_type"`

	// RewriteIP - need to set true if you want to use an external ip
	ExternalIP string `json:"external_ip"`
	InternalIP string `json:"internal_ip"`
}

type NodeState struct {
	// Node Capacity
	Capacity NodeResources `json:"capacity"`
	// Node Allocated
	Allocated NodeResources `json:"allocated"`
}

type NodeSpec struct {
	Routes  map[string]RouteSpec  `json:"routes"`
	Network map[string]Subnet     `json:"network"`
	Pods    map[string]PodSpec    `json:"pods"`
	Volumes map[string]VolumeSpec `json:"volumes"`
}

type NodeNamespace struct {
	Meta NamespaceMeta     `json:"meta",yaml:"meta"`
	Spec NodeNamespaceSpec `json:"spec",yaml:"spec"`
}

type NodeNamespaceSpec struct {
	Routes  []*Route  `json:"routes",yaml:"routes"`
	Pods    []*Pod    `json:"pods",yaml:"pods"`
	Volumes []*Volume `json:"volumes",yaml:"volumes"`
	Secrets []*Secret `json:"secrets",yaml:"secrets"`
}

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

type NodeRole struct {
	Router  NodeRoleRouter `json:"router"`
	Builder bool           `json:"builder"`
}

type NodeRoleRouter struct {
	ExternalIP string `json:"external_ip"`
	Enabled    bool   `json:"enabled"`
}

type NodeTask struct {
	Cancel context.CancelFunc
}
