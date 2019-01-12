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

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type ContainerView struct{}

func (cv *ContainerView) New(c *types.Container) Container {
	container := Container{
		ID:      c.ID,
		State:   c.State,
		Status:  c.Status,
		Image:   c.Image,
		Created: c.Created,
		Started: c.Started,
	}

	return container
}

func (cv *ContainerView) ToContainerSpec(spec *types.ContainerSpec) ContainerSpec {
	s := ContainerSpec{
		ID:            spec.ID,
		Meta:          cv.ToContainerSpecMeta(spec.Meta),
		Image:         cv.ToImageSpec(spec.Image),
		Network:       cv.ToContainerNetworkSpec(spec.Network),
		Labels:        spec.Labels,
		Envs:          spec.EnvVars,
		Entrypoint:    spec.Entrypoint,
		Command:       spec.Command,
		Args:          spec.Args,
		DNS:           cv.ToContainerDNSSpec(spec.DNS),
		Quota:         cv.ToContainerQuotaSpec(spec.Quota),
		RestartPolicy: cv.ToContainerRestartPolicySpec(spec.RestartPolicy),
	}

	s.Ports = make([]ContainerPortSpec, 0)
	s.Volumes = make([]ContainerVolumeSpec, 0)
	for _, port := range spec.Ports {
		s.Ports = append(s.Ports, cv.ToContainerPortSpec(port))
	}

	for _, volume := range spec.Volumes {
		s.Volumes = append(s.Volumes, cv.ToContainerVolumeSpec(volume))
	}

	return s
}

func (cv *ContainerView) ToContainerSpecMeta(meta types.ContainerSpecMeta) ContainerSpecMeta {
	return ContainerSpecMeta{
		Spec: meta.Spec,
	}
}

func (cv *ContainerView) ToContainerNetworkSpec(spec types.ContainerNetworkSpec) ContainerNetworkSpec {
	return ContainerNetworkSpec{
		Hostname: spec.Hostname,
		Domain:   spec.Domain,
		Network:  spec.Network,
		Mode:     spec.Mode,
	}
}

func (cv *ContainerView) ToContainerPortSpec(spec types.ContainerPortSpec) ContainerPortSpec {
	return ContainerPortSpec{
		ContainerPort: spec.ContainerPort,
		Protocol:      spec.Protocol,
	}
}

func (cv *ContainerView) ToContainerDNSSpec(spec types.ContainerDNSSpec) ContainerDNSSpec {
	return ContainerDNSSpec{
		Server:  spec.Server,
		Search:  spec.Search,
		Options: spec.Options,
	}
}

func (cv *ContainerView) ToContainerQuotaSpec(spec types.ContainerQuotaSpec) ContainerQuotaSpec {
	return ContainerQuotaSpec{
		Memory:    spec.Memory,
		CPUShares: spec.CPUShares,
	}
}

func (cv *ContainerView) ToContainerRestartPolicySpec(spec types.ContainerRestartPolicySpec) ContainerRestartPolicySpec {
	return ContainerRestartPolicySpec{
		Name:    spec.Name,
		Attempt: spec.Attempt,
	}
}

func (cv *ContainerView) ToContainerVolumeSpec(spec types.ContainerVolumeSpec) ContainerVolumeSpec {
	return ContainerVolumeSpec{
		Volume:    spec.Volume,
		MountPath: spec.MountPath,
	}
}

func (cv *ContainerView) FromContainerSpec(spec ContainerSpec) *types.ContainerSpec {
	s := &types.ContainerSpec{
		Meta:          cv.FromContainerSpecMeta(spec.Meta),
		Image:         cv.FromImageSpec(spec.Image),
		Network:       cv.FromContainerNetworkSpec(spec.Network),
		Labels:        spec.Labels,
		EnvVars:       spec.Envs,
		Entrypoint:    spec.Entrypoint,
		Command:       spec.Command,
		Args:          spec.Args,
		DNS:           cv.FromContainerDNSSpec(spec.DNS),
		Quota:         cv.FromContainerQuotaSpec(spec.Quota),
		RestartPolicy: cv.FromContainerRestartPolicySpec(spec.RestartPolicy),
	}

	for _, port := range spec.Ports {
		s.Ports = append(s.Ports, cv.FromContainerPortSpec(port))
	}

	for _, volume := range spec.Volumes {
		s.Volumes = append(s.Volumes, cv.FromContainerVolumeSpec(volume))
	}

	return s
}

func (cv *ContainerView) FromContainerSpecMeta(meta ContainerSpecMeta) types.ContainerSpecMeta {
	m := types.ContainerSpecMeta{}
	return m
}

func (cv *ContainerView) FromContainerNetworkSpec(spec ContainerNetworkSpec) types.ContainerNetworkSpec {
	return types.ContainerNetworkSpec{
		Hostname: spec.Hostname,
		Domain:   spec.Domain,
		Network:  spec.Network,
		Mode:     spec.Mode,
	}
}

func (cv *ContainerView) FromContainerPortSpec(spec ContainerPortSpec) types.ContainerPortSpec {
	return types.ContainerPortSpec{
		ContainerPort: spec.ContainerPort,
		Protocol:      spec.Protocol,
	}
}

func (cv *ContainerView) FromContainerDNSSpec(spec ContainerDNSSpec) types.ContainerDNSSpec {
	return types.ContainerDNSSpec{
		Server:  spec.Server,
		Search:  spec.Search,
		Options: spec.Options,
	}
}

func (cv *ContainerView) FromContainerQuotaSpec(spec ContainerQuotaSpec) types.ContainerQuotaSpec {
	return types.ContainerQuotaSpec{
		Memory:    spec.Memory,
		CPUShares: spec.CPUShares,
	}
}

func (cv *ContainerView) FromContainerRestartPolicySpec(spec ContainerRestartPolicySpec) types.ContainerRestartPolicySpec {
	return types.ContainerRestartPolicySpec{
		Name:    spec.Name,
		Attempt: spec.Attempt,
	}
}

func (cv *ContainerView) FromContainerVolumeSpec(spec ContainerVolumeSpec) types.ContainerVolumeSpec {
	return types.ContainerVolumeSpec{
		Volume:    spec.Volume,
		MountPath: spec.MountPath,
	}
}

func (cv *ContainerView) ToImageSpec(spec types.ImageSpec) ContainerImageSpec {
	return ContainerImageSpec{
		Name:   spec.Name,
		Secret: spec.Secret,
	}
}

func (cv *ContainerView) FromImageSpec(spec ContainerImageSpec) types.ImageSpec {
	return types.ImageSpec{
		Name:   spec.Name,
		Secret: spec.Secret,
	}
}

func (cv *ContainerView) NewPodContainer(c *types.PodContainer) PodContainer {

	container := PodContainer{
		ID:      c.ID,
		Pod:     c.Pod,
		Name:    c.Name,
		Ready:   c.Ready,
	}

	container.State.Error.Error = c.State.Error.Error
	container.State.Error.Message = c.State.Error.Message
	container.State.Created.Created = c.State.Created.Created
	container.State.Started.Started = c.State.Started.Started
	container.State.Started.Timestamp = c.State.Started.Timestamp
	container.State.Stopped.Stopped = c.State.Stopped.Stopped
	container.State.Stopped.Exit.Timestamp = c.State.Stopped.Exit.Timestamp
	container.State.Stopped.Exit.Code = c.State.Stopped.Exit.Code

	container.Image.Name = c.Image.Name
	container.Image.ID = c.Image.ID

	return container
}
