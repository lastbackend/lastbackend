package model

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
