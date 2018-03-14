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

package docker

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"strconv"
)

func GetConfig(spec *types.SpecTemplateContainer) *container.Config {

	var volumes map[string]struct{}
	var ports nat.PortSet

	ports = make(map[nat.Port]struct{})
	volumes = make(map[string]struct{})

	for _, p := range spec.Ports {
		port := nat.Port(strconv.Itoa(p.ContainerPort))
		ports[port] = struct{}{}
	}

	var envs []string
	for _, e := range spec.EnvVars {
		env := fmt.Sprintf("%s=%s", e.Name, e.Value)
		envs = append(envs, env)
	}

	return &container.Config{
		Hostname:     spec.Network.Hostname,
		Domainname:   spec.Network.Domain,
		Env:          envs,
		ExposedPorts: ports,
		Volumes:      volumes,
		Labels:       spec.Labels,
		Cmd:          strslice.StrSlice(spec.Exec.Command),
		Entrypoint:   strslice.StrSlice(spec.Exec.Entrypoint),
		Image:        spec.Image.Name,
		WorkingDir:   spec.Exec.Workdir,
	}
}

func GetHostConfig(spec *types.SpecTemplateContainer) *container.HostConfig {

	rPolicy := container.RestartPolicy{
		Name:              spec.RestartPolicy.Policy,
		MaximumRetryCount: spec.RestartPolicy.Attempt,
	}

	resources := container.Resources{
		Memory:    spec.Resources.Quota.RAM * 1024 * 1024,
		CPUShares: spec.Resources.Quota.CPU,
	}

	var (
		ports  nat.PortMap
		mounts []mount.Mount
		binds  []string
		links  []string
	)

	for _, v := range spec.Volumes {
		binds = append(binds, fmt.Sprintf("%s:%s:%s", v.Name, v.Path, v.Mode))
	}

	for _, l := range spec.Links {
		links = append(links, fmt.Sprintf("%s:%s", l.Link, l.Alias))
	}

	log := container.LogConfig{}
	ports = make(nat.PortMap)

	return &container.HostConfig{
		Binds:           binds,
		LogConfig:       log,
		NetworkMode:     container.NetworkMode(spec.Network.Mode),
		PortBindings:    ports,
		DNS:             spec.DNS.Server,
		DNSOptions:      spec.DNS.Options,
		DNSSearch:       spec.DNS.Search,
		RestartPolicy:   rPolicy,
		Resources:       resources,
		Mounts:          mounts,
		Links:           links,
		Privileged:      spec.Security.Privileged,
		ExtraHosts:      spec.ExtraHosts,
		PublishAllPorts: spec.PublishAllPorts,
		AutoRemove:      spec.AutoRemove,
	}
}

func GetNetworkConfig(spec *types.SpecTemplateContainer) *network.NetworkingConfig {

	cfg := &network.NetworkingConfig{
		//EndpointsConfig: make(map[string]*network.EndpointSettings),
	}

	//endpoint := &network.EndpointSettings{
	//	NetworkID: spec.Network.Subnet,
	//}
	//cfg.EndpointsConfig["lo"] = endpoint

	return cfg
}
