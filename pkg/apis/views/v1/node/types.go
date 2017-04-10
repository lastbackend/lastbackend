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

package node

import "time"

type Spec struct {
	// Pods spec
	Pods []Pod `json:"pods"`
}

type Pod struct {
	// Pod Meta
	Meta PodMeta `json:"meta"`
	// Container spec
	Spec PodSpec `json:"spec"`
	// Pod state
	State PodState `json:"state"`
}

type PodState struct {
	// Pod current state
	State string `json:"state"`
	// Pod current status
	Status string `json:"status"`
}

type PodMeta struct {
	// Meta id
	ID string `json:"id"`
	// Meta labels
	Labels map[string]string `json:"lables"`
	// Pod owner
	Owner string `json:"owner"`
	// Pod project
	Project string `json:"project"`
	// Pod service
	Service string `json:"service"`
	// Current Spec ID
	Spec string `json:"spec"`
	// Meta created time
	Created time.Time `json:"created"`
	// Meta updated time
	Updated time.Time `json:"updated"`
}

type PodSpec struct {
	// Provision ID
	ID string `json:"id"`
	// Provision state
	State string `json:"state"`
	// Provision status
	Status string `json:"status"`

	// Containers spec for pod
	Containers []ContainerSpec

	// Provision create time
	Created time.Time `json:"created"`
	// Provision update time
	Updated time.Time `json:"updated"`
}

type ContainerSpec struct {
	// Image spec
	Image ImageSpec `json:"image"`
	// Network spec
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
	DNS ContainerDNSSpec `json:"dns"`
	// Container resources quota
	Quota ContainerQuotaSpec `json:"quota"`
	// Container restart policy
	RestartPolicy ContainerRestartPolicySpec `json:"restart_policy"`
	// Container volumes mount
	Volumes []ContainerVolumeSpec `json:"volumes"`
}

type ImageSpec struct {
	// Image full name
	Name string `json:"name"`
	// Image pull provision flag
	Pull bool `json:"image-pull"`
	// Image Auth base64 encoded string
	Auth string `json:"auth"`
}

type ContainerNetworkSpec struct {
	// Container hostname
	Hostname string `json:"hostname"`
	// Container host domain
	Domain string `json:"domain"`
	// Network ID to use
	Network string `json:"network"`
	// Network Mode to use
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
