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
	"time"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

// Node - default node structure
// swagger:model views_node
type Node struct {
	Meta   NodeMeta   `json:"meta"`
	Status NodeStatus `json:"status"`
	Spec   NodeSpec   `json:"spec"`
}

// NodeList - node map list
// swagger:model views_node_list
type NodeList map[string]*Node

// NodeMeta - node metadata structure
// swagger:model views_node_meta
type NodeMeta struct {
	NodeInfo
	Name     string    `json:"name"`
	SelfLink string    `json:"self_link"`
	Subnet  string `json:"subnet"`
	Cluster string `json:"cluster"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

// NodeInfo - node info struct
// swagger:model views_node_info
type NodeInfo struct {
	Hostname     string `json:"hostname"`
	OSName       string `json:"os_name"`
	OSType       string `json:"os_type"`
	Architecture string `json:"architecture"`
	IP           struct {
		External string `json:"external"`
		Internal string `json:"internal"`
	} `json:"ip"`
}

// NodeStatus - node state struct
// swagger:model views_node_status
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

type NodeSpec struct {
	Network   map[string]SubnetSpec   `json:"network"`
	Pods      map[string]PodSpec      `json:"pods"`
	Volumes   map[string]VolumeSpec   `json:"volumes"`
	Endpoints map[string]EndpointSpec `json:"endpoints"`
	Security  NodeSecurity            `json:"security"`
}

type SubnetSpec struct {
	// Subnet state
	State string `json:"state"`
	// Node network type
	Type string `json:"type"`
	// Node Subnet subnet info
	CIDR string `json:"cidr"`
	// Node Subnet interface
	IFace NetworkInterface `json:"iface"`
	// Node Public IP
	Addr string `json:"addr"`
}

type NetworkInterface struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
	Addr  string `json:"addr"`
	HAddr string `json:"HAddr"`
}

type NodeSecurity struct {
	TLS bool `json:"tls"`
}

// NodeResources - node resources structure
// swagger:model views_node_resources
type NodeResources struct {
	Containers int   `json:"containers"`
	Pods       int   `json:"pods"`
	Memory     int64 `json:"memory"`
	Cpu        int   `json:"cpu"`
	Storage    int   `json:"storage"`
}

// swagger:model views_node_spec
type NodeManifest struct {
	Network   map[string]*types.SubnetManifest   `json:"network, omitempty"`
	Pods      map[string]*types.PodManifest      `json:"pods, omitempty"`
	Volumes   map[string]*types.VolumeManifest   `json:"volumes, omitempty"`
	Endpoints map[string]*types.EndpointManifest `json:"endpoints, omitempty"`
}
