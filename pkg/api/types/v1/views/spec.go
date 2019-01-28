//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

import "time"

// SpecTemplateVolumeList is a list of spec template volumes
// swagger:model types_spec_template_volume_list
type SpecTemplateVolumeList []*SpecTemplateVolume

// swagger:model types_spec_template_volume
type SpecTemplateVolume struct {
	// Template volume name
	Name string `json:"name"`
	// Template volume name
	Type string `json:"type"`
	// Template volume from persistent volume
	Volume SpecTemplateVolumeClaim `json:"volume,omitempty"`
	// Template volume from secret type
	Secret SpecTemplateSecretVolume `json:"secret,omitempty"`
	// Template volume from config type
	Config SpecTemplateConfigVolume `json:"config,omitempty"`
}

// SpecTemplateVolumeClaim - volume bind to use persistent volume in pod
// swagger:model types_spec_template_volume_claim
type SpecTemplateVolumeClaim struct {
	// Persistent volume name to mount
	Name string `json:"name"`
	// Persistent Volume Subpath
	Subpath string `json:"subpath"`
}

// SpecTemplateSecretVolume - use secret as volume in pod
type SpecTemplateSecretVolume struct {
	// Secret name to mount
	Name string `json:"name"`
	// Secret file key
	Binds []SpecTemplateSecretVolumeBind `json:"binds"`
}

// SpecTemplateSecretVolumeBind - files bindings.
// Get secret value by key and create file
type SpecTemplateSecretVolumeBind struct {
	Key  string `json:"key"`
	File string `json:"file"`
}

type SpecTemplateConfigVolume struct {
	// Secret name to mount
	Name string `json:"name"`
	// Config file binding
	Binds []SpecTemplateConfigVolumeBind `json:"binds"`
}

type SpecTemplateConfigVolumeBind struct {
	// Config key
	Key string `json:"key"`
	// File to create
	File string `json:"file"`
}

// swagger:ignore
// swagger:model types_spec_template_volume_mounts
type SpecTemplateVolumeMounts struct {
	// Template volume mounts name
	Name string `json:"name"`
}

// SpecTemplateContainers is a list of spec template containers
// swagger:model types_spec_template_container_list
type SpecTemplateContainers []*SpecTemplateContainer

type SpecTemplateContainer struct {
	// Template container id
	ID string `json:"id" yaml:"id"`
	// Template container name
	Name string `json:"name" yaml:"name"`
	// Template container role
	Role string `json:"role" yaml:"role"`
	// Automatically remove container when it exits
	AutoRemove bool `json:"autoremove" yaml:"autoremove"`
	// Labels list
	Labels map[string]string `json:"labels" yaml:"labels"`
	// Template container image
	Image SpecTemplateContainerImage `json:"image" yaml:"image"`
	// Template container ports binding
	Ports SpecTemplateContainerPorts `json:"ports" yaml:"ports"`
	// Template container envs
	EnvVars SpecTemplateContainerEnvs `json:"env" yaml:"env"`
	// Template container resources
	Resources SpecTemplateContainerResources `json:"resources" yaml:"resources"`
	// Template container exec options
	Exec SpecTemplateContainerExec `json:"exec" yaml:"exec"`
	// Template container volumes
	Volumes SpecTemplateContainerVolumes `json:"volumes" yaml:"volumes"`
	// Template container probes
	Probes SpecTemplateContainerProbes `json:"probes" yaml:"probes"`
	// Template container security
	Security SpecTemplateContainerSecurity `json:"security" yaml:"security"`
	// Subnet container settings
	Network SpecTemplateContainerNetwork `json:"network" yaml:"network"`
	// Container DNS configuration
	DNS SpecTemplateContainerDNS `json:"dns" yaml:"dns"`
	// List of extra hosts
	ExtraHosts []string `json:"extra_hosts" yaml:"extra_hosts"`
	// Should docker publish all exposed port for the container
	PublishAllPorts bool `json:"publish" yaml:"publish"`
	// Links to another containers
	Links []SpecTemplateContainerLink `json:"links" yaml:"links"`
	// Restart Policy
	RestartPolicy SpecTemplateRestartPolicy `json:"restart" yaml:"restart"`
}

type SpecTemplateContainerImage struct {
	Name   string `json:"name" yaml:"name"`
	Secret string `json:"secret" yaml:"secret"`
	Policy string `json:"policy" yaml:"policy"`
}

type SpecTemplateContainerPorts []*SpecTemplateContainerPort

type SpecTemplateContainerPort struct {
	// Container port
	ContainerPort uint16 `json:"container_port"`
	// Host port
	HostPort uint16 `json:"host_port"`
	// Host port
	HostIP string `json:"host_port"`
	// Binding protocol
	Protocol string `json:"protocol"`
}

// SpecTemplateContainerPorts is a list of spec template container env vars
// swagger:model types_spec_template_container_env_list
type SpecTemplateContainerEnvs []*SpecTemplateContainerEnv

// swagger:model types_spec_template_container_env
type SpecTemplateContainerEnv struct {
	Name   string                         `json:"name"`
	Value  string                         `json:"value,omitempty"`
	Secret SpecTemplateContainerEnvSecret `json:"secret,omitempty"`
	Config SpecTemplateContainerEnvConfig `json:"config,omitempty"`
}

// swagger:model types_spec_template_container_env_secret
type SpecTemplateContainerEnvSecret struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type SpecTemplateContainerEnvConfig struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type SpecTemplateContainerResources struct {
	Limits  SpecTemplateContainerResource `json:"limits"`
	Request SpecTemplateContainerResource `json:"request"`
}

// swagger:model types_spec_template_container_resource
type SpecTemplateContainerResource struct {
	// CPU resource option
	CPU string `json:"cpu"`
	// RAM resource option
	RAM string `json:"ram"`
}

type SpecTemplateContainerExec struct {
	Command []string `json:"command"`
	// Container enrtypoint
	Entrypoint []string `json:"entrypoint"`
	// Container run workdir option
	Workdir string `json:"workdir"`
	// Container run command arguments
	Args []string `json:"args"`
}

type SpecTemplateContainerVolumes []*SpecTemplateContainerVolume

type SpecTemplateContainerVolume struct {
	Name string `json:"name"`
	Mode string `json:"mode"`
	Path string `json:"path"`
}

type SpecTemplateContainerProbes struct {
	LiveProbe SpecTemplateContainerProbe `json:"live_probe"`
	ReadProbe SpecTemplateContainerProbe `json:"read_probe"`
}

type SpecTemplateContainerProbe struct {
	Exec struct {
		Command []string `json:"command"`
	} `json:"exec"`

	Socket struct {
		Protocol string `json:"protocol"`
		Port     int    `json:"port"`
	} `json:"socket"`

	InitialDelaySeconds int `json:"initial_delay"`
	TimeoutSeconds      int `json:"timeout_seconds"`
	PeriodSeconds       int `json:"period_seconds"`
	ThresholdSuccess    int `json:"threshold_success"`
	ThresholdFailure    int `json:"threshold_failure"`
}

type SpecTemplateContainerSecurity struct {
	// Start container in priveleged mode
	Privileged bool `json:"privileged"`
	// Add linux security options
	LinuxOptions SpecTemplateContainerSecurityLinuxOptions `json:"linux_options"`
	// Run container as particular user
	User int `json:"user"`
}

// swagger:model types_spec_template_container_security_linux
type SpecTemplateContainerSecurityLinuxOptions struct {
	Level string `json:"level"`
}

// swagger:model types_spec_template_container_network
type SpecTemplateContainerNetwork struct {
	// Container hostname
	Hostname string `json:"hostname"`
	// Container host domain
	Domain string `json:"domain"`
	// Subnet ID to use
	Network string `json:"network"`
	// Subnet Mode to use
	Mode string `json:"mode"`
}

// swagger:model types_spec_template_container_dns
type SpecTemplateContainerDNS struct {
	// List of DNS servers
	Server []string `json:"server"`
	// DNS server search options
	Search []string `json:"search"`
	// DNS server other options
	Options []string `json:"options"`
}

// swagger:model types_spec_template_container_link
type SpecTemplateContainerLink struct {
	// Link name
	Link string `json:"link"`
	// Container alias
	Alias string `json:"alias"`
}

// swagger:model types_spec_template_policy
type SpecTemplateRestartPolicy struct {
	// Restart policy name
	Policy string `json:"policy" yaml:"policy"`
	// Attempt period
	Attempt int `json:"attempt" yaml:"attempt"`
}

// swagger:model types_spec_strategy
type SpecStrategy struct {
	Type           string                     `json:"type"` // Rolling
	RollingOptions SpecStrategyRollingOptions `json:"rollingOptions"`
	Resources      SpecStrategyResources      `json:"resources"`
	Deadline       int                        `json:"deadline"`
	// Spec updated time
	Updated time.Time `json:"updated"`
}

// swagger:model types_spec_strategy_resources
type SpecStrategyResources struct {
}

// swagger:model types_spec_strategy_rolling
type SpecStrategyRollingOptions struct {
	PeriodUpdate   int `json:"period_update"`
	Interval       int `json:"interval"`
	Timeout        int `json:"timeout"`
	MaxUnavailable int `json:"max_unavailable"`
	MaxSurge       int `json:"max_surge"`
}

// SpecTriggers is a list of spec triggers
// swagger:model types_spec_trigger_list
type SpecTriggers []SpecTrigger

// swagger:model types_spec_trigger
type SpecTrigger struct {
}

// swagger:model types_spec_selector
type SpecSelector struct {
	Labels map[string]string `json:"labels"`

	Node string `json:"node"`
	// Spec updated time
	Updated time.Time `json:"updated"`
}
