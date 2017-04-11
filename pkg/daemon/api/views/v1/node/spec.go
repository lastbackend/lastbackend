package node

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"time"
)

type Spec struct {
	// Pods spec
	Pods []Pod `json:"pods"`
}

type Pod struct {
	// Pod Meta
	Meta PodMeta `json:"meta"`
	// Container spec
	Spec PodSpec `json:"spec"`
	// Pod state
	State PodState `json:"state"`
}

type PodState struct {
	// Pod current state
	State string `json:"state"`
	// Pod current status
	Status string `json:"status"`
}

type PodMeta struct {
	// Meta id
	ID string `json:"id"`
	// Meta labels
	Labels map[string]string `json:"lables"`
	// Pod owner
	Owner string `json:"owner"`
	// Pod project
	Project string `json:"project"`
	// Pod service
	Service string `json:"service"`
	// Current Spec ID
	Spec string `json:"spec"`
	// Meta created time
	Created time.Time `json:"created"`
	// Meta updated time
	Updated time.Time `json:"updated"`
}

type PodSpec struct {
	// Provision ID
	ID string `json:"id"`
	// Provision state
	State string `json:"state"`
	// Provision status
	Status string `json:"status"`

	// Containers spec for pod
	Containers []ContainerSpec `json:"containers"`

	// Provision create time
	Created time.Time `json:"created"`
	// Provision update time
	Updated time.Time `json:"updated"`
}

type ContainerSpec struct {
	// Image spec
	Image ImageSpec `json:"image"`
	// Network spec
	Network ContainerNetworkSpec `json:"network"`
	// Ports configuration
	Ports []ContainerPortSpec `json:"ports"`
	// Labels list
	Labels map[string]string `json:"labels"`
	// Environments list
	Envs []string `json:"envs"`
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
}

type ImageSpec struct {
	// Image full name
	Name string `json:"name"`
	// Image pull provision flag
	Pull bool `json:"pull"`
	// Image Auth base64 encoded string
	Auth string `json:"auth"`
}

type ContainerNetworkSpec struct {
	// Container hostname
	Hostname string `json:"hostname"`
	// Container host domain
	Domain string `json:"domain"`
	// Network ID to use
	Network string `json:"network"`
	// Network Mode to use
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
