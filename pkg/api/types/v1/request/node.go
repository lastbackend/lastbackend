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

package request

import "github.com/lastbackend/lastbackend/pkg/distribution/types"

type NodeMetaOptions struct {
	Meta *types.NodeUpdateMetaOptions `json:"meta"`
}

type NodeConnectOptions struct {
	Info   types.NodeInfo   `json:"info"`
	Status types.NodeStatus `json:"status"`
}

type NodeStatusOptions struct {
	// Node state capacity
	Capacity types.NodeResources `json:"capacity"`
	// Node state allocated
	Allocated types.NodeResources `json:"allocated"`
}

type NodePodStatusOptions struct {
	// Pod state
	State string `json:"state" yaml:"state"`
	// Pod state message
	Message string `json:"message" yaml:"message"`
	// Pod steps
	Steps types.PodSteps `json:"steps" yaml:"steps"`
	// Pod network
	Network types.PodNetwork `json:"network" yaml:"network"`
	// Pod containers
	Containers map[string]*types.PodContainer `json:"containers" yaml:"containers"`
}

type NodeVolumeStatusOptions struct {
	// route status state
	State string `json:"state" yaml:"state"`
	// route status message
	Message string `json:"message" yaml:"message"`
}

type NodeRouteStatusOptions struct {
	// route status state
	State string `json:"state" yaml:"state"`
	// route status message
	Message string `json:"message" yaml:"message"`
}

type NodeRemoveOptions struct {
	Force bool `json:"force"`
}

type NodeLogsOptions struct {
	Follow bool `json:"follow"`
}
