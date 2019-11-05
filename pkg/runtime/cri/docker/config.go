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

package docker

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"strconv"
)

func GetConfig(manifest *types.ContainerManifest) *container.Config {

	var (
		volumes = make(map[string]struct{}, 0)
		ports   = make(nat.PortSet, 0)
	)

	for _, p := range manifest.Ports {
		port := nat.Port(strconv.Itoa(int(p.ContainerPort)))
		ports[port] = struct{}{}
	}

	return &container.Config{
		Hostname:     manifest.Network.Hostname,
		Domainname:   manifest.Network.Domain,
		Env:          manifest.Envs,
		ExposedPorts: ports,
		Volumes:      volumes,
		Labels:       manifest.Labels,
		Cmd:          manifest.Exec.Command,
		Entrypoint:   manifest.Exec.Entrypoint,
		Image:        manifest.Image,
		WorkingDir:   manifest.Exec.Workdir,
	}
}

func GetHostConfig(manifest *types.ContainerManifest) *container.HostConfig {

	rPolicy := container.RestartPolicy{
		Name:              manifest.RestartPolicy.Policy,
		MaximumRetryCount: manifest.RestartPolicy.Attempt,
	}

	resources := container.Resources{
		Memory:   manifest.Resources.Limits.RAM,
		NanoCPUs: manifest.Resources.Limits.CPU,
	}

	var (
		ports  = make(nat.PortMap, 0)
		mounts []mount.Mount
		links  []string
	)

	for _, l := range manifest.Links {
		links = append(links, fmt.Sprintf("%s:%s", l.Link, l.Alias))
	}

	logC := container.LogConfig{
		Type: "lastbackend",
	}
	ports = make(nat.PortMap)

	for _, p := range manifest.Ports {
		port := nat.Port(strconv.Itoa(int(p.ContainerPort)))
		ports[port] = make([]nat.PortBinding, 0)
		if p.HostPort != 0 {

			if p.HostIP == types.EmptyString {
				p.HostIP = "0.0.0.0"
			}

			ports[port] = append(ports[port], nat.PortBinding{
				HostIP:   p.HostIP,
				HostPort: fmt.Sprintf("%d", p.HostPort),
			})
		}

	}

	cfg := container.HostConfig{
		Binds:           manifest.Binds,
		LogConfig:       logC,
		NetworkMode:     container.NetworkMode(manifest.Network.Mode),
		PortBindings:    ports,
		DNS:             manifest.DNS.Server,
		DNSOptions:      manifest.DNS.Options,
		DNSSearch:       manifest.DNS.Search,
		RestartPolicy:   rPolicy,
		Resources:       resources,
		Mounts:          mounts,
		Links:           links,
		Privileged:      manifest.Security.Privileged,
		ExtraHosts:      manifest.ExtraHosts,
		PublishAllPorts: manifest.PublishAllPorts,
		AutoRemove:      manifest.AutoRemove,
	}

	if cfg.NetworkMode != types.EmptyString {
		cfg.DNS = make([]string, 0)
		cfg.DNSOptions = make([]string, 0)
		cfg.DNSSearch = make([]string, 0)
	}

	return &cfg
}

func GetNetworkConfig(manifest *types.ContainerManifest) *network.NetworkingConfig {

	cfg := &network.NetworkingConfig{}

	return cfg
}
