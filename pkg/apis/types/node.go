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

package types

import "time"

type NodeList []Node

type Node struct {
	// Node metadata
	Meta NodeMeta `json:"meta"`
	// Node spec info
	Spec NodeSpec `json:"spec"`
	// Node Addresses
	Addresses []NodeAddress `json:"addresses"`
	// Node Capacity
	Capacity NodeState `json:"capacity"`
	// Node Allocated
	Allocated NodeState `json:"allocated"`
	// Node labels list
	Labels map[string]string `json:"labels"`
	// Node images list
	Images []NodeImage `json:"images"`
}

type NodeMeta struct {
	// Node unique ID
	ID string `json:"ID"`
	// Node cluster link
	Cluster string `json:"cluster"`
	// Node hostname
	HostName string `json:"hostname"`
	// Node created time
	Created time.Time `json:"created"`
	// Node updated time
	Updated time.Time `json:"updated"`
}

type NodeSpec struct {
	Network string `json:"network"`
	CIDR    string `json:"cidr"`
}

type NodeAddress struct {
	Type    string `json:"type"`
	Address string `json:"address"`
}

type NodeState struct {
	Containers int    `json:"containers"`
	Pods       int    `json:"pods"`
	Memory     string `json:"memory"`
	Cpu        int    `json:"cpu"`
}

type NodeImage struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
	Sha  string `json:"sha"`
}
