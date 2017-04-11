package v1

import (
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/container/views/v1"
)

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
		s.Containers = append(s.Containers, v1.ToContainerSpec(*c))
	}

	return s
}

func ToPodState(state types.PodState) PodState {
	return PodState{
		State:  state.State,
		Status: state.Status,
	}
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
		s.Containers = append(s.Containers, v1.FromContainerSpec(c))
	}

	return s
}

func FromPodState(state PodState) types.PodState {
	return types.PodState{
		State:  state.State,
		Status: state.Status,
	}
}
