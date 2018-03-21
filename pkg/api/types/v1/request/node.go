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

type NodeUpdateOptions struct {
	Meta *types.NodeUpdateMetaOptions
}

type NodeInfoOptions struct {
	// Node hostname
	Hostname string `json:"hostname"`
	// Linux architecture
	Architecture string `json:"architecture"`
	// OS information
	OSName string `json:"os_name"`
	OSType string `json:"os_type"`
	// RewriteIP - need to set true if you want to use an external ip
	ExternalIP string `json:"external_ip"`
	InternalIP string `json:"internal_ip"`
}

type NodeStateOptions struct {
	// Node state capacity
	Capacity types.NodeResources `json:"capacity"`
	// Node state allocated
	Allocated types.NodeResources `json:"allocated"`
}

type NodePodStatusOptions struct {
	// Pod stage
	Stage string `json:"stage" yaml:"stage"`
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
	// route status stage
	Stage string `json:"stage" yaml:"stage"`
	// route status message
	Message string `json:"message" yaml:"message"`
}

type NodeRouteStatusOptions struct {
	// route status stage
	Stage string `json:"stage" yaml:"stage"`
	// route status message
	Message string `json:"message" yaml:"message"`
}

type NodeRemoveOptions struct {
	Force bool `json:"force"`
}
