//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

import (
	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

// swagger:model request_node_meta
type NodeMetaOptions struct {
	Meta *models.NodeUpdateMetaOptions `json:"meta"`
}

// swagger:model request_node_connect
type NodeConnectOptions struct {
	Info    models.NodeInfo      `json:"info"`
	Status  models.NodeStatus    `json:"status"`
	Network *models.NetworkState `json:"network"`
	TLS     bool                 `json:"tls"`
	SSL     *SSL                 `json:"ssl"`
}

type SSL struct {
	CA   []byte `json:"ca"`
	Cert []byte `json:"cert"`
	Key  []byte `json:"key"`
}

// swagger:model request_node_status
type NodeStatusOptions struct {
	// Node interface options
	State models.NodeStatusState `json:"state"`
	// Pods statuses
	Pods map[string]*NodePodStatusOptions `json:"pods"`
	// Volumes statuses
	Volumes map[string]*NodeVolumeStatusOptions `json:"volumes"`
	// Node resources
	Resources NodeResourcesOptions `json:"resources"`
}

// swagger:model request_node_resources
type NodeResourcesOptions struct {
	// Node state capacity
	Capacity models.NodeResources `json:"capacity"`
	// Node state allocated
	Allocated models.NodeResources `json:"allocated"`
}

// swagger:model request_node_pod_status
type NodePodStatusOptions struct {
	// Pod state
	State string `json:"state" yaml:"state"`
	// Pod state
	Status string `json:"status" yaml:"status"`
	// Pod state
	Running bool `json:"running" yaml:"running"`
	// Pod state message
	Message string `json:"message" yaml:"message"`
	// Pod steps
	Steps models.PodSteps `json:"steps" yaml:"steps"`
	// Pod network
	Network models.PodNetwork `json:"network" yaml:"network"`
	// Pod containers
	Runtime models.PodStatusRuntime `json:"runtime" yaml:"runtime"`
}

// swagger:model request_node_volume_status
type NodeVolumeStatusOptions struct {
	// route status state
	State string `json:"state" yaml:"state"`
	// route status message
	Message string `json:"message" yaml:"message"`
}

// swagger:model request_node_route_status
type NodeRouteStatusOptions struct {
	// route status state
	State string `json:"state" yaml:"state"`
	// route status message
	Message string `json:"message" yaml:"message"`
}

// swagger:ignore
// swagger:model request_node_remove
type NodeRemoveOptions struct {
	Force bool `json:"force"`
}

// swagger:ignore
// swagger:model request_node_logs
type NodeLogsOptions struct {
	Follow bool `json:"follow"`
}
