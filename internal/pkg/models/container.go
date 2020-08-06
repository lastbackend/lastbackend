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

package models

import (
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/internal/util/generator"
	"github.com/lastbackend/lastbackend/internal/util/resource"
	"github.com/lastbackend/lastbackend/tools/logger"
	"strings"
	"context"
	"time"
)

const (
	ContainerTypeLBC            = "LBC"
	ContainerTypeRuntime        = "LBC:Container"
	ContainerTypeRuntimeService = "service"
	ContainerTypeRuntimeTask    = "task"
)

type Container struct {
	// Container CID
	ID string `json:"id"`
	// Container Pod ID
	Pod string `json:"pod"`
	// Container Deployment ID
	Deployment string `json:"deployment"`
	// Container Environment ID
	Namespace string `json:"namespace"`
	// Container envs
	Envs []string `json:"envs"`
	// Container binds
	Binds []string `json:"binds"`
	// Container name
	Name string `json:"name"`
	// Name information
	Image string `json:"image"`
	// Container current state
	State string `json:"state"`
	// Container error message
	Error string `json:"state"`
	// ExitCode of the container
	ExitCode int `json:"exit_code"`
	// Container current state
	Status string `json:"status,omitempty"`
	// Container network settings
	Network NetworkSettings `json:"network"`
	// Container labels
	Labels map[string]string `json:"labels"`
	// Contaienr restart policy
	Restart struct {
		Policy string `json:"policy"`
		Retry  int    `json:"count"`
	} `json:"restart"`
	// Container created time
	Created time.Time `json:"created"`
	// Container started time
	Started time.Time `json:"started"`

	Exec SpecTemplateContainerExec `json:"exec"`
}

type Port struct {
	// HostIP is the host IP Address
	HostIP string `json:"host_ip"`
	// HostPort is the host port number
	HostPort string `json:"host_port"`
}

type NetworkSettings struct {
	// Container gatway ip address
	Gateway string `json:"gateway"`
	// Container ip address
	IPAddress string `json:"ip"`
	// Container ports mapping
	Ports []*SpecTemplateContainerPort `json:"ports"`
}

type ContainerSpec struct {
	ID string `json:"id"`
	// Container meta spec
	Meta ContainerSpecMeta `json:"meta"`
	// Name spec
	Image ImageSpec `json:"image"`
	// Subnet spec
	Network ContainerNetworkSpec `json:"network"`
	// Ports configuration
	Ports []ContainerPortSpec `json:"ports"`
	// Labels list
	Labels map[string]string `json:"labels"`
	// Environments list
	EnvVars []string `json:"environments"`
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
	// Links to another containers
	Links []ContainerLinkSpec `json:"links"`
	// Container in privileged mode
	Privileged bool `json:"privileged"`
	// PWD where the commands will be run
	Workdir string `json:"workdir"`
	// List of extra hosts
	ExtraHosts []string `json:"extra_hosts"`
	// Should docker publish all exposed port for the container
	PublishAllPorts bool `json:"publish_all_ports"`
}

type ContainerManifest struct {
	// Template container id
	ID string `json:"id" yaml:"id"`
	// Template container name
	Name string `json:"name" yaml:"name"`
	// Labels list
	Labels map[string]string `json:"labels" yaml:"labels"`
	// Template container image
	Image string `json:"image" yaml:"image"`
	// Template container ports binding
	Ports SpecTemplateContainerPorts `json:"ports" yaml:"ports"`
	// Template container envs
	Envs []string `json:"env" yaml:"env"`
	// Template container resources
	Resources SpecTemplateContainerResources `json:"resources" yaml:"resources"`
	// Template container exec options
	Exec SpecTemplateContainerExec `json:"exec" yaml:"exec"`
	// Template container binds
	Binds []string `json:"volumes" yaml:"volumes"`
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
	// AutoRemove flag
	AutoRemove bool `json:"autoremove" yaml:"autoremove"`
}

type ContainerSpecMeta struct {
	Meta
	// Applications id
	Service string `json:"service"`
	// Applications spec id
	Spec string `json:"spec"`
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

type ContainerLinkSpec struct {
	// Link name
	Link string `json:"link"`
	// Container alias
	Alias string `json:"alias"`
}

type ContainerStatusInfo struct {
	// Container ID on host
	ID string `json:"cid"`
	// Name ID
	Image string `json:"image"`
	// Container current state
	State string `json:"state"`
	// Container current state
	Status string `json:"status"`
	// Container ports mapping
	Ports map[string][]ContainerStatusInfoPort `json:"ports"`
	// Container created time
	Created time.Time `json:"created"`
	// Container updated time
	Updated time.Time `json:"updated"`
}

type ContainerStatusInfoPort struct {
	HostIP   string `json:"host_ip"`
	HostPort string `json:"host_port"`
}

func (cs *ContainerSpec) CommandToString() string {
	log := logger.WithContext(context.Background())
	res, err := convertSliceToString(cs.Command)
	if err != nil {
		log.Errorf("Can-not convert command value to string: %s", err)
		return EmptyStringSlice
	}
	return res
}

func (cs *ContainerSpec) CommandFromString(command string) {
	log := logger.WithContext(context.Background())
	if err := json.Unmarshal([]byte(command), &cs.Command); err != nil {
		log.Errorf("Can-not convert command value from string: %s", err)
	}
}

func (cs *ContainerSpec) EntrypointToString() string {
	log := logger.WithContext(context.Background())
	res, err := convertSliceToString(cs.Entrypoint)
	if err != nil {
		log.Errorf("Can-not convert entrypoint value to string: %s", err)
		return EmptyStringSlice
	}
	return res
}

func (cs *ContainerSpec) EntrypointFromString(entrypoint string) {
	log := logger.WithContext(context.Background())
	if err := json.Unmarshal([]byte(entrypoint), &cs.Entrypoint); err != nil {
		log.Errorf("Can-not convert entrypoint value from string: %s", err)
	}
}

func (cs *ContainerSpec) DNSServerToString() string {
	log := logger.WithContext(context.Background())
	res, err := convertSliceToString(cs.DNS.Server)
	if err != nil {
		log.Errorf("Can-not convert dns server value to string: %s", err)
		return EmptyStringSlice
	}
	return res
}

func (cs *ContainerSpec) DNSServerFromString(server string) {
	log := logger.WithContext(context.Background())
	if err := json.Unmarshal([]byte(server), &cs.DNS.Server); err != nil {
		log.Errorf("Can-not convert dns server value from string: %s", err)
	}
}

func (cs *ContainerSpec) DNSSearchToString() string {
	log := logger.WithContext(context.Background())
	res, err := convertSliceToString(cs.DNS.Search)
	if err != nil {
		log.Errorf("Can-not convert dns search value to string: %s", err)
		return EmptyStringSlice
	}
	return res
}

func (cs *ContainerSpec) DNSSearchFromString(search string) {
	log := logger.WithContext(context.Background())
	if err := json.Unmarshal([]byte(search), &cs.DNS.Search); err != nil {
		log.Errorf("Can-not convert dns search value from string: %s", err)
	}
}

func (cs *ContainerSpec) DNSOptionsToString() string {
	log := logger.WithContext(context.Background())
	res, err := convertSliceToString(cs.DNS.Options)
	if err != nil {
		log.Errorf("Can-not convert dns options value to string: %s", err)
		return EmptyStringSlice
	}
	return res
}

func (cs *ContainerSpec) DNSOptionsFromString(options string) {
	log := logger.WithContext(context.Background())
	if err := json.Unmarshal([]byte(options), &cs.DNS.Options); err != nil {
		log.Errorf("Can-not convert dns options value from string: %s", err)
	}
}

func (cs *ContainerSpec) VolumesToString() string {
	log := logger.WithContext(context.Background())
	b, err := json.Marshal(cs.Volumes)
	if err != nil {
		log.Errorf("Can-not convert volumes value to string: %s", err)
		return EmptyStringSlice
	}
	if string(b) == "null" {
		return EmptyStringSlice
	}
	return string(b)
}

func (cs *ContainerSpec) VolumesFromString(volumes string) {
	log := logger.WithContext(context.Background())
	if err := json.Unmarshal([]byte(volumes), &cs.Volumes); err != nil {
		log.Errorf("Can-not convert volumes value from string: %s", err)
	}
}

func (cs *ContainerSpec) ENVsToString() string {
	log := logger.WithContext(context.Background())
	res, err := convertSliceToString(cs.EnvVars)
	if err != nil {
		log.Errorf("Can-not convert envs value to string: %s", err)
		return EmptyStringSlice
	}
	return res
}

func (cs *ContainerSpec) ENVsFromString(envs string) {
	log := logger.WithContext(context.Background())
	if err := json.Unmarshal([]byte(envs), &cs.EnvVars); err != nil {
		log.Errorf("Can-not convert envs value from string: %s", err)
	}
}

func (cs *ContainerSpec) PortsToString() string {
	log := logger.WithContext(context.Background())
	if cs == nil {
		return EmptyStringSlice
	}
	res, err := json.Marshal(cs.Ports)
	if err != nil {
		log.Errorf("Can-not convert ports value to string: %s", err)
		return EmptyStringSlice
	}
	if string(res) == "null" {
		return EmptyStringSlice
	}
	return string(res)
}

func (cs *ContainerSpec) PortsFromString(ports string) {
	log := logger.WithContext(context.Background())
	if err := json.Unmarshal([]byte(ports), &cs.Ports); err != nil {
		log.Errorf("Can-not convert ports value from string: %s", err)
	}
}

func convertSliceToString(slice []string) (string, error) {
	log := logger.WithContext(context.Background())
	if slice == nil {
		return EmptyStringSlice, nil
	}
	res, err := json.Marshal(slice)
	if err != nil {
		log.Errorf("Can-not convert ports value to string: %s", err)
		return EmptyString, err
	}
	if string(res) == "null" {
		return EmptyStringSlice, nil
	}
	return string(res), nil
}

func NewContainerManifest(spec *SpecTemplateContainer) *ContainerManifest {

	mf := new(ContainerManifest)
	mf.Resources = spec.Resources
	mf.Labels = spec.Labels
	mf.Ports = spec.Ports
	mf.Network = spec.Network
	mf.Exec = spec.Exec
	mf.DNS = spec.DNS
	mf.ExtraHosts = spec.ExtraHosts
	mf.RestartPolicy = spec.RestartPolicy
	mf.PublishAllPorts = spec.PublishAllPorts
	mf.Security = spec.Security
	mf.AutoRemove = spec.AutoRemove

	mf.Image = spec.Image.Name
	if len(spec.Image.Sha) != 0 {
		mf.Image = fmt.Sprintf("%s@%s", strings.Split(spec.Image.Name, ":")[0], spec.Image.Sha)
	}

	var envs []string
	for _, e := range spec.EnvVars {
		var env string
		if e.Value == EmptyString {
			continue
		}
		env = fmt.Sprintf("%s=%s", e.Name, e.Value)
		envs = append(envs, env)
	}

	mf.Envs = envs
	return mf
}

func GetPauseContainerTemplate() *SpecTemplateContainer {
	tpl := new(SpecTemplateContainer)
	tpl.Name = generator.GenerateRandomString(6)
	tpl.Image = SpecTemplateContainerImage{Name: "busybox"}
	tpl.Exec.Entrypoint = []string{"/bin/sh"}
	tpl.Exec.Command = []string{"-c", "while true; do sleep 30m; done;"}
	tpl.Resources.Request.RAM, _ = resource.DecodeMemoryResource("32MB")
	tpl.Resources.Request.CPU, _ = resource.DecodeCpuResource("0.1")
	tpl.Resources.Limits.RAM, _ = resource.DecodeMemoryResource("32MB")
	tpl.Resources.Limits.CPU, _ = resource.DecodeCpuResource("0.1")
	return tpl
}
