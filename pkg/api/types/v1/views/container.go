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

type Container struct {
	// Container ID
	ID string `json:"id"`
	// Spec ID
	Spec ContainerSpec `json:"-"`
	// Repo information
	Image string `json:"image"`
	// Container current state
	State string `json:"state"`
	// Container current state
	Status string `json:"status"`
	// Container ports mapping
	Ports map[string][]*types.Port `json:"ports"`
	// Container created time
	Created time.Time `json:"created"`
	// Container started time
	Started time.Time `json:"started"`
}

type ContainerSpec struct {
	ID   string            `json:"id"`
	Meta ContainerSpecMeta `json:"meta"`
	// Repo spec
	Image ContainerImageSpec `json:"image"`
	// Subnet spec
	Network ContainerNetworkSpec `json:"network"`
	// Ports configuration
	Ports []ContainerPortSpec `json:"ports"`
	// Labels list
	Labels map[string]string `json:"labels"`
	// Environments list
	Envs []string `json:"envs"`
	// Container enrtypoint
	Entrypoint []string `json:"entrypoint"`
	// Container run command
	Command []string `json:"command"`
	// Container run command arguments
	Args []string `json:"args"`
	// Container DNS configuration
	DNS ContainerDNSSpec `json:""`
	// Container resources quota
	Quota ContainerQuotaSpec `json:"quota"`
	// Container restart policy
	RestartPolicy ContainerRestartPolicySpec `json:"restart_policy"`
	// Container volumes mount
	Volumes []ContainerVolumeSpec `json:"volumes"`
}

type ContainerSpecMeta struct {
	// Meta id
	ID string `json:"id"`
	// Meta id
	Spec string `json:"spec"`
	// Meta name
	Name string `json:"name,omitempty"`
	// Meta description
	Description string `json:"description,omitempty"`
	// Meta labels
	Labels map[string]string `json:"lables,omitempty"`
	// Meta created time
	Created time.Time `json:"created"`
	// Meta updated time
	Updated time.Time `json:"updated"`
}

type ContainerNetworkSpec struct {
	// Container hostname
	Hostname string `json:"hostname"`
	// Container host domain
	Domain string `json:"domain"`
	// Subnet ID to use
	Network string `json:"network"`
	// Subnet Mode to use
	Mode string `json:"mode"`
}

type ContainerPortSpec struct {
	// Container port to expose
	ContainerPort int `json:"container_port"`
	// Containers protocol allowed on exposed port
	Protocol string `json:"protocol"`
}

type ContainerDNSSpec struct {
	// List of DNS servers
	Server []string `json:"server"`
	// DNS server search options
	Search []string `json:"search"`
	// DNS server other options
	Options []string `json:"options"`
}

type ContainerQuotaSpec struct {
	// Maximum memory allowed to use
	Memory int64 `json:"memory"`
	// CPU shares for container on one node
	CPUShares int64 `json:"cpu_shares"`
}

type ContainerRestartPolicySpec struct {
	// Restart policy name
	Name string `json:"name"`
	// Attempt to restart container
	Attempt int `json:"attempt"`
}

type ContainerVolumeSpec struct {
	// Volume name
	Volume string `json:"volume"`
	// Container mount path
	MountPath string `json:"mount_path"`
}

type ContainerImageSpec struct {
	// Name full name
	Name string `json:"name"`
	// Name pull provision flag
	Secret string `json:"secret"`
}
