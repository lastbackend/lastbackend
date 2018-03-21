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

type ServiceCreateOptions struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Sources     *string             `json:"sources"`
	Replicas    *int                `json:"replicas"`
	Spec        *ServiceOptionsSpec `json:"spec"`
}

type ServiceUpdateOptions struct {
	Description *string             `json:"description"`
	Spec        *ServiceOptionsSpec `json:"spec"`
}

type ServiceRemoveOptions struct {
	Force bool `json:"force"`
}

type ServiceOptionsSpec struct {
	Memory     *int64                    `json:"memory,omitempty"`
	Entrypoint *string                   `json:"entrypoint,omitempty"`
	Command    *string                   `json:"command,omitempty"`
	EnvVars    *[]string                 `json:"env,omitempty"`
	Ports      *[]ServiceOptionsSpecPort `json:"ports,omitempty"`
}

type ServiceOptionsSpecPort struct {
	Internal int    `json:"internal"`
	Protocol string `json:"protocol"`
}
