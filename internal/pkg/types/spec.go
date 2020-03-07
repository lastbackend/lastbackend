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

package types

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"
)

const ContainerRolePrimary = "primary"
const ContainerRoleSlave = "slave"

// SpecState is a state of the spec
// swagger:model types_spec_state
type SpecState struct {
	Destroy     bool `json:"destroy"`
	Cancel      bool `json:"cancel"`
	Maintenance bool `json:"maintenance"`
}

// SpecRuntime is a runtime of the spec
// swagger:model types_spec_runtime
type SpecRuntime struct {
	Services []string          `json:"services"`
	Tasks    []SpecRuntimeTask `json:"tasks"`
	Updated  time.Time         `json:"updated"`
}

// SpecRuntimeTask is a runtime task to execute in runtime
// swagger:model types_spec_runtime_task
type SpecRuntimeTask struct {
	Name      string                    `json:"name"`
	Container string                    `json:"container" yaml:"container"`
	EnvVars   SpecTemplateContainerEnvs `json:"env" yaml:"env"`
	Commands  []string                  `json:"commands" yaml:"commands"`
}

// SpecTemplate is a template of the spec
// swagger:model types_spec_template
type SpecTemplate struct {
	// Template spec for volume
	Volumes SpecTemplateVolumeList `json:"volumes" yaml:"volumes"`
	// Template main container
	Containers SpecTemplateContainers `json:"containers" yaml:"containers"`
	// Termination period
	Termination int `json:"termination" yaml:"termination"`
	// Spec updated time
	Updated time.Time `json:"updated" yaml:"updated"`
}

// SpecNetwork is a map of spec template for network
// swagger:model types_spec_template_network
type SpecNetwork struct {
	IP       string               `json:"ip"`
	Ports    map[uint16]string    `json:"ports"`
	Strategy EndpointSpecStrategy `json:"strategy"`
	Policy   string               `json:"policy"`
	// Spec updated time
	Updated time.Time `json:"updated"`
}

// swagger:ignore
// SpecTemplateVolumeMap is a map of spec template volumes
// swagger:model types_spec_template_volume_map
type SpecTemplateVolumeMap map[string]*SpecTemplateVolume

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

// swagger:model types_spec_template_container
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

// swagger:model types_spec_template_container_image
type SpecTemplateContainerImage struct {
	Name   string                           `json:"name" yaml:"name"`
	Sha    string                           `json:"sha" yaml:"sha"`
	Secret SpecTemplateContainerImageSecret `json:"secret,omitempty" yaml:"secret"`
	Policy string                           `json:"policy,omitempty" yaml:"policy"`
}

// swagger:model types_spec_template_container_image
type SpecTemplateContainerImageSecret struct {
	Name string `json:"name" yaml:"name"`
	Key  string `json:"key" yaml:"key"`
}

// swagger:ignore
// SpecBuildImage is an image of the spec build
// swagger:model types_spec_build_image
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

// swagger:ignore
// AuthConfig contains authorization information for connecting to a Registry
// swagger:model types_authConfig
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

// SpecTemplateContainerPorts is a list of spec template container ports
// swagger:model types_spec_template_container_port_list
type SpecTemplateContainerPorts []*SpecTemplateContainerPort

// SpecTemplateContainerPort is a port of the spec template container
// swagger:model types_spec_template_container_port
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

// swagger:model types_spec_template_container_resources
type SpecTemplateContainerResources struct {
	// Limit resources
	Limits SpecTemplateContainerResource `json:"limits"`
	// Request resources
	Request SpecTemplateContainerResource `json:"quota"`
}

// swagger:model types_spec_volume_resources
type SpecVolumeCapacity struct {
	// Limit resources
	Storage int64 `json:"storage"`
}

// swagger:model types_spec_volume_resource
type SpecVolumeResource struct {
	// Size resource option
	Size int64 `json:"size"`
}

// swagger:model types_spec_template_container_exec
type SpecTemplateContainerExec struct {
	Command []string `json:"command"`
	// Container enrtypoint
	Entrypoint []string `json:"entrypoint"`
	// Container run workdir option
	Workdir string `json:"workdir"`
	// Container run command arguments
	Args []string `json:"args"`
}

// swagger:model types_spec_template_container_resource
type SpecTemplateContainerResource struct {
	// CPU resource option
	CPU int64 `json:"cpu"`
	// RAM resource option
	RAM int64 `json:"ram"`
}

// SpecTemplateContainerVolumes is a list of spec template container volumes
// swagger:model types_spec_template_container_volume_list
type SpecTemplateContainerVolumes []*SpecTemplateContainerVolume

// swagger:model types_spec_template_container_volume
type SpecTemplateContainerVolume struct {
	// Volume name
	Name string `json:"name"`
	// Volume mount mode
	Mode string `json:"mode"`
	// Volume mount path
	MountPath string `json:"path"`
	// Volume sub path
	SubPath string `json:"sub_path"`
}

// swagger:model types_spec_template_container_probes
type SpecTemplateContainerProbes struct {
	LiveProbe SpecTemplateContainerProbe `json:"live_probe"`
	ReadProbe SpecTemplateContainerProbe `json:"read_probe"`
}

// swagger:model types_spec_template_container_probe
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

// swagger:model types_spec_template_container_security
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

func (s *SpecTemplateContainerEnvs) ToLinuxFormat() []string {
	env := make([]string, 0)

	for _, e := range *s {
		env = append(env, fmt.Sprintf("%s=%s", e.Name, e.Value))
	}

	return env
}

func (ss *SpecSelector) SetDefault() {
	if ss.Node != EmptyString {
		ss.Node = EmptyString
		ss.Updated = time.Now()
	}

	if len(ss.Labels) > 0 {
		ss.Labels = make(map[string]string)
		ss.Updated = time.Now()
	}
}

func (s *SpecTemplate) SetDefault() {
	// Set default configurations

	s.Containers = make(SpecTemplateContainers, 1)
	s.Volumes = make(SpecTemplateVolumeList, 0)
}

func (s *SpecTemplateContainer) SetDefault() {
	s.Labels = make(map[string]string, 0)
	s.Resources.Limits.RAM = int64(128)
	s.Resources.Request.RAM = int64(128)
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

func (s *SpecTemplateContainerPort) Parse(p string) {

	var (
		base = 10
		size = 16
	)

	reg, _ := regexp.Compile(`((?P<host>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\:)?((?P<hport>\d{1,5})?)?(\:(?P<cport>\d{1,5}))?(\/(?P<proto>\w{3}))?`)

	match := reg.FindStringSubmatch(p)
	pm := make(map[string]string)
	for i, name := range reg.SubexpNames() {
		if i > 0 && i <= len(match) {
			pm[name] = match[i]
		}
	}

	if _, ok := pm["host"]; ok {
		if pm["host"] != EmptyString {
			s.HostIP = pm["host"]
		} else {
			s.HostIP = "127.0.0.1"
		}
	} else {
		s.HostIP = "127.0.0.1"
	}

	if _, ok := pm["proto"]; ok {

		if pm["proto"] != EmptyString {
			s.Protocol = pm["proto"]
		} else {
			s.Protocol = "tcp"
		}
	} else {
		s.Protocol = "tcp"
	}

	if _, ok := pm["hport"]; ok {
		if pt, err := strconv.ParseUint(pm["hport"], base, size); err == nil {
			s.HostPort = uint16(pt)
		}
	}

	if _, ok := pm["cport"]; ok {

		if pm["cport"] != EmptyString {
			if pt, err := strconv.ParseUint(pm["cport"], base, size); err == nil {
				s.ContainerPort = uint16(pt)
			}
		} else {
			s.ContainerPort = s.HostPort
		}
	} else {
		s.ContainerPort = s.HostPort
	}
}
