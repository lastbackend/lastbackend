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

package views

import (
	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

type ContainerView struct{}

func (cv *ContainerView) New(c *models.Container) Container {
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

func (cv *ContainerView) ToContainerSpec(spec *models.ContainerSpec) ContainerSpec {
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

func (cv *ContainerView) ToContainerSpecMeta(meta models.ContainerSpecMeta) ContainerSpecMeta {
	return ContainerSpecMeta{
		Spec: meta.Spec,
	}
}

func (cv *ContainerView) ToContainerNetworkSpec(spec models.ContainerNetworkSpec) ContainerNetworkSpec {
	return ContainerNetworkSpec{
		Hostname: spec.Hostname,
		Domain:   spec.Domain,
		Network:  spec.Network,
		Mode:     spec.Mode,
	}
}

func (cv *ContainerView) ToContainerPortSpec(spec models.ContainerPortSpec) ContainerPortSpec {
	return ContainerPortSpec{
		ContainerPort: spec.ContainerPort,
		Protocol:      spec.Protocol,
	}
}

func (cv *ContainerView) ToContainerDNSSpec(spec models.ContainerDNSSpec) ContainerDNSSpec {
	return ContainerDNSSpec{
		Server:  spec.Server,
		Search:  spec.Search,
		Options: spec.Options,
	}
}

func (cv *ContainerView) ToContainerQuotaSpec(spec models.ContainerQuotaSpec) ContainerQuotaSpec {
	return ContainerQuotaSpec{
		Memory:    spec.Memory,
		CPUShares: spec.CPUShares,
	}
}

func (cv *ContainerView) ToContainerRestartPolicySpec(spec models.ContainerRestartPolicySpec) ContainerRestartPolicySpec {
	return ContainerRestartPolicySpec{
		Name:    spec.Name,
		Attempt: spec.Attempt,
	}
}

func (cv *ContainerView) ToContainerVolumeSpec(spec models.ContainerVolumeSpec) ContainerVolumeSpec {
	return ContainerVolumeSpec{
		Volume:    spec.Volume,
		MountPath: spec.MountPath,
	}
}

func (cv *ContainerView) FromContainerSpec(spec ContainerSpec) *models.ContainerSpec {
	s := &models.ContainerSpec{
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

func (cv *ContainerView) FromContainerSpecMeta(meta ContainerSpecMeta) models.ContainerSpecMeta {
	m := models.ContainerSpecMeta{}
	return m
}

func (cv *ContainerView) FromContainerNetworkSpec(spec ContainerNetworkSpec) models.ContainerNetworkSpec {
	return models.ContainerNetworkSpec{
		Hostname: spec.Hostname,
		Domain:   spec.Domain,
		Network:  spec.Network,
		Mode:     spec.Mode,
	}
}

func (cv *ContainerView) FromContainerPortSpec(spec ContainerPortSpec) models.ContainerPortSpec {
	return models.ContainerPortSpec{
		ContainerPort: spec.ContainerPort,
		Protocol:      spec.Protocol,
	}
}

func (cv *ContainerView) FromContainerDNSSpec(spec ContainerDNSSpec) models.ContainerDNSSpec {
	return models.ContainerDNSSpec{
		Server:  spec.Server,
		Search:  spec.Search,
		Options: spec.Options,
	}
}

func (cv *ContainerView) FromContainerQuotaSpec(spec ContainerQuotaSpec) models.ContainerQuotaSpec {
	return models.ContainerQuotaSpec{
		Memory:    spec.Memory,
		CPUShares: spec.CPUShares,
	}
}

func (cv *ContainerView) FromContainerRestartPolicySpec(spec ContainerRestartPolicySpec) models.ContainerRestartPolicySpec {
	return models.ContainerRestartPolicySpec{
		Name:    spec.Name,
		Attempt: spec.Attempt,
	}
}

func (cv *ContainerView) FromContainerVolumeSpec(spec ContainerVolumeSpec) models.ContainerVolumeSpec {
	return models.ContainerVolumeSpec{
		Volume:    spec.Volume,
		MountPath: spec.MountPath,
	}
}

func (cv *ContainerView) ToImageSpec(spec models.ImageSpec) ContainerImageSpec {
	return ContainerImageSpec{
		Name:   spec.Name,
		Secret: spec.Secret,
	}
}

func (cv *ContainerView) FromImageSpec(spec ContainerImageSpec) models.ImageSpec {
	return models.ImageSpec{
		Name:   spec.Name,
		Secret: spec.Secret,
	}
}

func (cv *ContainerView) NewPodContainer(c *models.PodContainer) PodContainer {

	container := PodContainer{
		ID:    c.ID,
		Pod:   c.Pod,
		Name:  c.Name,
		Ready: c.Ready,
	}

	container.State.Error.Error = c.State.Error.Error
	container.State.Error.Message = c.State.Error.Message
	container.State.Created.Timestamp = c.State.Created.Created
	container.State.Started.Started = c.State.Started.Started
	container.State.Started.Timestamp = c.State.Started.Timestamp
	container.State.Stopped.Stopped = c.State.Stopped.Stopped
	container.State.Stopped.Exit.Timestamp = c.State.Stopped.Exit.Timestamp
	container.State.Stopped.Exit.Code = c.State.Stopped.Exit.Code

	container.Image.Name = c.Image.Name
	container.Image.Sha = c.Image.Sha
	container.Image.ID = c.Image.ID

	return container
}
