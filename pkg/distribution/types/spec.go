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

package types

import (
	"fmt"
	"io"
)

const ContainerRolePrimary = "primary"
const ContainerRoleSlave = "slave"

type SpecState struct {
	Destroy     bool `json:"destroy"`
	Maintenance bool `json:"maintenance"`
}

type SpecTemplate struct {
	// Template Volume
	Volumes SpecTemplateVolumes `json:"volumes"`
	// Template main container
	Containers SpecTemplateContainers `json:"container"`
	// Termination period
	Termination int `json:"termination"`
}

type SpecTemplateVolumes []SpecTemplateVolume

type SpecTemplateVolume struct {
	// Template volume name
	Name string `json:"name"`
}

type SpecTemplateVolumeMounts struct {
	// Template volume mounts name
	Name string `json:"name"`
}

type SpecTemplateContainers []SpecTemplateContainer

type SpecTemplateContainer struct {
	// Template container name
	Name string `json:"name"`
	// Template container role
	Role string `json:"role"`
	// Automatically remove container when it exits
	AutoRemove bool `json:"autoremove"`
	// Labels list
	Labels map[string]string `json:"labels"`
	// Template container image
	Image SpecTemplateContainerImage `json:"image"`
	// Template container ports binding
	Ports SpecTemplateContainerPorts `json:"ports"`
	// Template container envs
	EnvVars SpecTemplateContainerEnvs `json:"env"`
	// Template container resources
	Resources SpecTemplateContainerResources `json:"resources"`
	// Template container exec options
	Exec SpecTemplateContainerExec `json:"exec"`
	// Template container volumes
	Volumes SpecTemplateContainerVolumes `json:"volumes"`
	// Template container probes
	Probes SpecTemplateContainerProbes `json:"probes"`
	// Template container security
	Security SpecTemplateContainerSecurity `json:"security"`
	// Network container settings
	Network SpecTemplateContainerNetwork `json:"network"`
	// Container DNS configuration
	DNS SpecTemplateContainerDNS `json:"dns"`
	// List of extra hosts
	ExtraHosts []string `json:"extra_hosts"`
	// Should docker publish all exposed port for the container
	PublishAllPorts bool `json:"publish"`
	// Links to another containers
	Links []SpecTemplateContainerLink `json:"links"`
	// Restart Policy
	RestartPolicy SpecTemplateRestartPolicy `json:"restart"`
}

type SpecTemplateContainerImage struct {
	Name   string `json:"name"`
	Auth   string `json:"auth"`
	Policy string `json:"policy"`
}

type SpecBuildImage struct {
	Tags           []string
	NoCache        bool
	Memory         int64
	Dockerfile     string
	SuppressOutput bool
	AuthConfigs    map[string]AuthConfig
	Context        io.Reader
	ExtraHosts     []string // List of extra hosts
}

// AuthConfig contains authorization information for connecting to a Registry
type AuthConfig struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Auth     string `json:"auth,omitempty"`
	// Email is an optional value associated with the username.
	// This field is deprecated and will be removed in a later
	// version of docker.
	Email         string `json:"email,omitempty"`
	ServerAddress string `json:"serveraddress,omitempty"`
	// IdentityToken is used to authenticate the user and get
	// an access token for the registry.
	IdentityToken string `json:"identitytoken,omitempty"`
	// RegistryToken is a bearer token to be sent to a registry
	RegistryToken string `json:"registrytoken,omitempty"`
}

type SpecTemplateContainerPorts []SpecTemplateContainerPort

type SpecTemplateContainerPort struct {
	// Container port
	ContainerPort int `json:"container_port"`
	// Binding protocol
	Protocol string `json:"protocol"`
}

type SpecTemplateContainerEnvs []SpecTemplateContainerEnv

type SpecTemplateContainerEnv struct {
	Name  string                         `json:"name"`
	Value string                         `json:"value"`
	From  SpecTemplateContainerEnvSecret `json:"from"`
}

type SpecTemplateContainerEnvSecret struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type SpecTemplateContainerResources struct {
	// Limit resources
	Limits SpecTemplateContainerResource `json:"limits"`
	// Quota resources
	Quota SpecTemplateContainerResource `json:"quota"`
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

type SpecTemplateContainerResource struct {
	// CPU resource option
	CPU int64 `json:"cpu"`
	// RAM resource option
	RAM int64 `json:"ram"`
}

type SpecTemplateContainerVolumes []SpecTemplateContainerVolume

type SpecTemplateContainerVolume struct {
	// Volume name
	Name string `json:"name"`
	// Volume mount mode
	Mode string `json:"mode"`
	// Volume mount path
	Path string `json:"path"`
}

type SpecTemplateContainerProbes struct {
	LiveProbe SpecTemplateContainerProbe `json:"live_probe"`
	ReadProbe SpecTemplateContainerProbe `json:"read_probe"`
}

type SpecTemplateContainerProbe struct {
	// Exec command to check container liveness
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

type SpecTemplateContainerSecurityLinuxOptions struct {
	Level string `json:"level"`
}

type SpecTemplateContainerNetwork struct {
	// Container hostname
	Hostname string `json:"hostname"`
	// Container host domain
	Domain string `json:"domain"`
	// Network ID to use
	Network string `json:"network"`
	// Network Mode to use
	Mode string `json:"mode"`
}

type SpecTemplateContainerDNS struct {
	// List of DNS servers
	Server []string `json:"server"`
	// DNS server search options
	Search []string `json:"search"`
	// DNS server other options
	Options []string `json:"options"`
}

type SpecTemplateContainerLink struct {
	// Link name
	Link string `json:"link"`
	// Container alias
	Alias string `json:"alias"`
}

type SpecTemplateRestartPolicy struct {
	// Restart policy name
	Policy string `json:"policy"`
	// Attempt period
	Attempt int `json:"attempt"`
}

type SpecStrategy struct {
	Type           string                     `json:"type"` // Rolling
	RollingOptions SpecStrategyRollingOptions `json:"rollingOptions"`
	Resources      SpecStrategyResources      `json:"resources"`
	Deadline       int                        `json:"deadline"`
}

type SpecStrategyResources struct {
}

type SpecStrategyRollingOptions struct {
	PeriodUpdate   int `json:"period_update"`
	Interval       int `json:"interval"`
	Timeout        int `json:"timeout"`
	MaxUnavailable int `json:"max_unavailable"`
	MaxSurge       int `json:"max_surge"`
}

type SpecTriggers []SpecTrigger

type SpecTrigger struct {
}

type SpecSelector struct {
}

func (s *SpecTemplateContainerEnvs) ToLinuxFormat() []string {
	env := make([]string, 0)

	for _, e := range *s {
		env = append(env, fmt.Sprintf("%s=%s", e.Name, e.Value))
	}

	return env
}

func (s *SpecTemplate) SetDefault() {
	// Set default configurations

	s.Containers = make(SpecTemplateContainers, 1)
	s.Volumes = make(SpecTemplateVolumes, 0)
}

func (s *SpecTemplateContainer) SetDefault() {
	s.Resources.Limits.RAM = int64(128)
	s.EnvVars = make(SpecTemplateContainerEnvs, 0)
	s.Ports = make(SpecTemplateContainerPorts, 0)
	s.Volumes = make(SpecTemplateContainerVolumes, 0)
	s.Exec.Command = make([]string, 0)
	s.Exec.Entrypoint = make([]string, 0)
	s.Probes.LiveProbe.Exec.Command = make([]string, 0)
	s.Probes.ReadProbe.Exec.Command = make([]string, 0)
	s.RestartPolicy = SpecTemplateRestartPolicy{
		Policy: "always",
	}
}
