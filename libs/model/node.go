package model

import "time"

type NodeList []Node

type Node struct {
	// Node uuid, generated automatically
	ID string `json:"id"`
	// Node user
	Owner string `json:"user"`
	// Node hostname
	Hostname string `json:"hostname"`

	// Node tag lists
	Labels map[string]string `json:"labels"`
	// Node created time
	Created time.Time `json:"created"`
	// Node updated time
	Updated time.Time `json:"updated"`
}

type NodeMemory struct {
	// Total node memory in MB
	Total int `json:"total"`
	// Used node memory in MB
	Used int `json:"used"`
}

type NodeContainers struct {
	// Total node memory in MB
	Total int `json:"total"`
	// Running containers count
	Running int `json:"running"`
	// Error containers count
	Error int `json:"error"`
}

