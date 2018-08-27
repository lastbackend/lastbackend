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
	"github.com/lastbackend/lastbackend/pkg/log"
	"strconv"
)

func GetConfig(spec *types.SpecTemplateContainer, secrets map[string]*types.Secret) *container.Config {

	var (
		volumes = make(map[string]struct{}, 0)
		ports = make(nat.PortSet, 0)
	)

	for _, p := range spec.Ports {
		port := nat.Port(strconv.Itoa(int(p.ContainerPort)))
		ports[port] = struct{}{}
	}

	var envs []string
	for _, e := range spec.EnvVars {
		var env string
		if e.From.Name != types.EmptyString {
			if _, ok := secrets[e.From.Name]; !ok {
				continue
			}
			v, err := secrets[e.From.Name].DecodeSecretTextData(e.From.Key)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			env = fmt.Sprintf("%s=%s", e.Name, v)
		} else {
			env = fmt.Sprintf("%s=%s", e.Name, e.Value)
		}
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
		Memory:    spec.Resources.Request.RAM * 1024 * 1024,
		CPUShares: spec.Resources.Request.CPU,
	}

	var (
		ports  = make(nat.PortMap, 0)
		mounts []mount.Mount
		binds  []string
		links  []string
	)

	for _, v := range spec.Volumes {
		if v.Name == types.EmptyString || v.Path == types.EmptyString {
			continue
		}

		if v.Mode != "rw" {
			v.Mode = "ro"
		}

		binds = append(binds, fmt.Sprintf("%s:%s:%s", v.Name, v.Path, v.Mode))
	}

	for _, l := range spec.Links {
		links = append(links, fmt.Sprintf("%s:%s", l.Link, l.Alias))
	}

	logC := container.LogConfig{}
	ports = make(nat.PortMap)

	return &container.HostConfig{
		Binds:           binds,
		LogConfig:       logC,
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
	//	NetworkID: spec.Subnet.SubnetSpec,
	//}
	//cfg.EndpointsConfig["lo"] = endpoint

	return cfg
}
