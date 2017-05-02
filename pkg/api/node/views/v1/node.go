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

type Node struct {
	Meta  NodeMeta  `json:"meta"`
	State NodeState `json:"state"`
}

type NodeMeta struct {
	Hostname     string `json:"hostname"`
	OSName       string `json:"os_name"`
	OSType       string `json:"os_type"`
	Architecture string `json:"architecture"`

	CPU     NodeCPU     `json:"cpu"`
	Memory  NodeMemory  `json:"memory"`
	Network NodeNetwork `json:"network"`
	Storage NodeStorage `json:"storage"`
}

type NodeCPU struct {
	Name  string `json:"name"`
	Cores int64  `json:"cores"`
}

type NodeMemory struct {
	Total     int64 `json:"total"`
	Used      int64 `json:"used"`
	Available int64 `json:"available"`
}

type NodeNetwork struct {
	Interface string   `json:"interface,omitempty"`
	IP        []string `json:"ip,omitempty"`
}

type NodeStorage struct {
	Available string `json:"available"`
	Used      string `json:"used"`
	Total     string `json:"total"`
}

type NodeState struct {
}

type NodeList []*Node
