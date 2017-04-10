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

package node

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
)

func (obj *Spec) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func ToNodeSpec(obj []*types.Pod) Spec {
	spec := Spec{}
	for _, pod := range obj {
		spec.Pods = append(spec.Pods, Pod{
			Meta:  ToPodMeta(pod.Meta),
			Spec:  ToPodSpec(pod.Spec),
			State: ToPodState(pod.State),
		})
	}
	return spec
}

func ToPodMeta(meta types.PodMeta) PodMeta {
	return PodMeta{
		ID:      meta.ID,
		Labels:  meta.Labels,
		Owner:   meta.Owner,
		Project: meta.Project,
		Service: meta.Service,
		Spec:    meta.Spec,
		Created: meta.Created,
		Updated: meta.Updated,
	}
}

func ToPodSpec(spec types.PodSpec) PodSpec {
	s := PodSpec{
		ID:      spec.ID,
		State:   spec.State,
		Status:  spec.Status,
		Created: spec.Created,
		Updated: spec.Updated,
	}

	for _, c := range spec.Containers {
		s.Containers = append(s.Containers, ToContainerSpec(*c))
	}

	return s
}

func ToPodState(state types.PodState) PodState {
	return PodState{
		State:  state.State,
		Status: state.Status,
	}
}

func ToContainerSpec(spec types.ContainerSpec) ContainerSpec {
	s := ContainerSpec{
		Image:         ToImageSpec(spec.Image),
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

func ToImageSpec(spec types.ImageSpec) ImageSpec {
	return ImageSpec{
		Name: spec.Name,
		Pull: spec.Pull,
		Auth: spec.Auth,
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

func FromNodeSpec(spec Spec) []*types.Pod {
	var pods []*types.Pod
	for _, s := range spec.Pods {
		pod := types.NewPod()
		pod.Meta = FromPodMeta(s.Meta)
		pod.Spec = FromPodSpec(s.Spec)
		pod.State = FromPodState(s.State)
		pods = append(pods, pod)
	}
	return pods
}

func FromPodMeta(meta PodMeta) types.PodMeta {
	m := types.PodMeta{}
	m.ID = meta.ID
	m.Labels = meta.Labels
	m.Owner = meta.Owner
	m.Project = meta.Project
	m.Service = meta.Service
	m.Spec = meta.Spec
	m.Created = meta.Created
	m.Updated = meta.Updated
	return m
}

func FromPodSpec(spec PodSpec) types.PodSpec {
	s := types.PodSpec{
		ID:      spec.ID,
		State:   spec.State,
		Status:  spec.Status,
		Created: spec.Created,
		Updated: spec.Updated,
	}

	for _, c := range spec.Containers {
		s.Containers = append(s.Containers, FromContainerSpec(c))
	}

	return s
}

func FromPodState(state PodState) types.PodState {
	return types.PodState{
		State:  state.State,
		Status: state.Status,
	}
}

func FromContainerSpec(spec ContainerSpec) *types.ContainerSpec {
	s := &types.ContainerSpec{
		Image:         FromImageSpec(spec.Image),
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

func FromImageSpec(spec ImageSpec) types.ImageSpec {
	return types.ImageSpec{
		Name: spec.Name,
		Pull: spec.Pull,
		Auth: spec.Auth,
	}
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
