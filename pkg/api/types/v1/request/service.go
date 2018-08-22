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

// swagger:model request_service_create
type ServiceCreateOptions struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Image       *ServiceImageSpec   `json:"image"`
	Spec        *ServiceOptionsSpec `json:"spec"`
}

// swagger:model request_service_update
type ServiceUpdateOptions struct {
	Description *string             `json:"description"`
	Image       *ServiceImageSpec   `json:"image"`
	Spec        *ServiceOptionsSpec `json:"spec"`
}

// swagger:ignore
// swagger:model request_service_remove
type ServiceRemoveOptions struct {
	Force bool `json:"force"`
}

type ServiceImageSpec struct {
	Name   *string `json:"name"`
	Secret *string `json:"secret"`
}

// swagger:model request_service_spec
type ServiceOptionsSpec struct {
	Replicas   *int              `json:"replicas"`
	Memory     *int64            `json:"memory,omitempty"`
	Entrypoint *string           `json:"entrypoint,omitempty"`
	Command    *string           `json:"command,omitempty"`
	EnvVars    *[]string         `json:"env,omitempty"`
	Ports      map[uint16]string `json:"ports,omitempty"`
}

// swagger:ignore
// swagger:model request_service_logs
type ServiceLogsOptions struct {
	Deployment string `json:"deployment"`
	Pod        string `json:"pod"`
	Container  string `json:"container"`
	Follow     bool   `json:"follow"`
}
