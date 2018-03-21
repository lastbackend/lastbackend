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
type Node struct {
	Meta  NodeMeta  `json:"meta"`
	Info  NodeInfo  `json:"info"`
	State NodeState `json:"state"`
}

// NodeList - node map list
type NodeList map[string]*Node

// NodeMeta - node metadata structure
type NodeMeta struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

// NodeInfo - node info struct
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

// NodeState - node state struct
type NodeState struct {
	Online    bool          `json:"online"`
	Capacity  NodeResources `json:"capacity"`
	Allocated NodeResources `json:"allocated"`
}

// NodeResources - node resources structure
type NodeResources struct {
	Containers int   `json:"containers"`
	Pods       int   `json:"pods"`
	Memory     int64 `json:"memory"`
	Cpu        int   `json:"cpu"`
	Storage    int   `json:"storage"`
}

type NodeSpec struct {
	Routes  map[string]types.RouteSpec  `json:"routes"`
	Network map[string]types.Subnet     `json:"network"`
	Pods    map[string]types.PodSpec    `json:"pods"`
	Volumes map[string]types.VolumeSpec `json:"volumes"`
}
