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

package v1

import (
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/image/views/v1"
)

func ToContainer(c *types.Container) Container {
	container := Container{
		ID:      c.ID,
		Spec:    c.Spec,
		State:   c.State,
		Status:  c.Status,
		Image:   c.Image,
		Created: c.Created,
		Started: c.Started,
	}

	container.Ports = make(map[string]int, len(c.Ports))
	if len(c.Ports) != 0 {
		container.Ports = c.Ports
	}

	return container
}

func ToContainerSpec(spec types.ContainerSpec) ContainerSpec {
	s := ContainerSpec{
		Meta:          ToContainerSpecMeta(spec.Meta),
		Image:         v1.ToImageSpec(spec.Image),
		Network:       ToContainerNetworkSpec(spec.Network),
		Labels:        spec.Labels,
		Envs:          spec.Envs,
		Entrypoint:    spec.Entrypoint,
		Command:       spec.Command,
		Args:          spec.Args,
		DNS:           ToContainerDNSSpec(spec.DNS),
		Quota:         ToContainerQuotaSpec(spec.Quota),
		RestartPolicy: ToContainerRestartPolicySpec(spec.RestartPolicy),
	}

	for _, port := range spec.Ports {
		s.Ports = append(s.Ports, ToContainerPortSpec(port))
	}

	for _, volume := range spec.Volumes {
		s.Volumes = append(s.Volumes, ToContainerVolumeSpec(volume))
	}

	return s
}

func ToContainerSpecMeta(meta types.ContainerSpecMeta) ContainerSpecMeta {
	return ContainerSpecMeta{
		ID: meta.ID,
	}
}

func ToContainerNetworkSpec(spec types.ContainerNetworkSpec) ContainerNetworkSpec {
	return ContainerNetworkSpec{
		Hostname: spec.Hostname,
		Domain:   spec.Domain,
		Network:  spec.Network,
		Mode:     spec.Mode,
	}
}

func ToContainerPortSpec(spec types.ContainerPortSpec) ContainerPortSpec {
	return ContainerPortSpec{
		ContainerPort: spec.ContainerPort,
		Protocol:      spec.Protocol,
	}
}

func ToContainerDNSSpec(spec types.ContainerDNSSpec) ContainerDNSSpec {
	return ContainerDNSSpec{
		Server:  spec.Server,
		Search:  spec.Search,
		Options: spec.Options,
	}
}

func ToContainerQuotaSpec(spec types.ContainerQuotaSpec) ContainerQuotaSpec {
	return ContainerQuotaSpec{
		Memory:    spec.Memory,
		CPUShares: spec.CPUShares,
	}
}

func ToContainerRestartPolicySpec(spec types.ContainerRestartPolicySpec) ContainerRestartPolicySpec {
	return ContainerRestartPolicySpec{
		Name:    spec.Name,
		Attempt: spec.Attempt,
	}
}

func ToContainerVolumeSpec(spec types.ContainerVolumeSpec) ContainerVolumeSpec {
	return ContainerVolumeSpec{
		Volume:    spec.Volume,
		MountPath: spec.MountPath,
	}
}

func FromContainerSpec(spec ContainerSpec) *types.ContainerSpec {
	s := &types.ContainerSpec{
		Meta:          FromContainerSpecMeta(spec.Meta),
		Image:         v1.FromImageSpec(spec.Image),
		Network:       FromContainerNetworkSpec(spec.Network),
		Labels:        spec.Labels,
		Envs:          spec.Envs,
		Entrypoint:    spec.Entrypoint,
		Command:       spec.Command,
		Args:          spec.Args,
		DNS:           FromContainerDNSSpec(spec.DNS),
		Quota:         FromContainerQuotaSpec(spec.Quota),
		RestartPolicy: FromContainerRestartPolicySpec(spec.RestartPolicy),
	}

	for _, port := range spec.Ports {
		s.Ports = append(s.Ports, FromContainerPortSpec(port))
	}

	for _, volume := range spec.Volumes {
		s.Volumes = append(s.Volumes, FromContainerVolumeSpec(volume))
	}

	return s
}

func FromContainerSpecMeta(meta ContainerSpecMeta) types.ContainerSpecMeta {
	m := types.ContainerSpecMeta{}
	meta.ID = m.ID
	return m
}

func FromContainerNetworkSpec(spec ContainerNetworkSpec) types.ContainerNetworkSpec {
	return types.ContainerNetworkSpec{
		Hostname: spec.Hostname,
		Domain:   spec.Domain,
		Network:  spec.Network,
		Mode:     spec.Mode,
	}
}

func FromContainerPortSpec(spec ContainerPortSpec) types.ContainerPortSpec {
	return types.ContainerPortSpec{
		ContainerPort: spec.ContainerPort,
		Protocol:      spec.Protocol,
	}
}

func FromContainerDNSSpec(spec ContainerDNSSpec) types.ContainerDNSSpec {
	return types.ContainerDNSSpec{
		Server:  spec.Server,
		Search:  spec.Search,
		Options: spec.Options,
	}
}

func FromContainerQuotaSpec(spec ContainerQuotaSpec) types.ContainerQuotaSpec {
	return types.ContainerQuotaSpec{
		Memory:    spec.Memory,
		CPUShares: spec.CPUShares,
	}
}

func FromContainerRestartPolicySpec(spec ContainerRestartPolicySpec) types.ContainerRestartPolicySpec {
	return types.ContainerRestartPolicySpec{
		Name:    spec.Name,
		Attempt: spec.Attempt,
	}
}

func FromContainerVolumeSpec(spec ContainerVolumeSpec) types.ContainerVolumeSpec {
	return types.ContainerVolumeSpec{
		Volume:    spec.Volume,
		MountPath: spec.MountPath,
	}
}
