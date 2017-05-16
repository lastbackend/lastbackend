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
	"github.com/lastbackend/lastbackend/pkg/api/container/views/v1"
	"github.com/lastbackend/lastbackend/pkg/common/types"
)

func ToPodInfo(pod *types.Pod) PodInfo {
	info := PodInfo{
		Meta:  ToPodMeta(pod.Meta),
		State: ToPodState(pod.State),
	}

	if len(pod.Containers) == 0 {
		info.Containers = make([]v1.Container, 0)
		return info
	}

	for _, c := range pod.Containers {
		info.Containers = append(info.Containers, v1.ToContainer(c))
	}

	return info
}

func ToPodMeta(meta types.PodMeta) PodMeta {
	m := PodMeta{
		Name:    meta.Name,
		Labels:  meta.Labels,
		Created: meta.Created,
		Updated: meta.Updated,
	}

	if len(m.Labels) == 0 {
		m.Labels = make(map[string]string)
	}

	return m
}

func ToPodSpec(spec types.PodSpec) PodSpec {
	s := PodSpec{
		ID:         spec.ID,
		State:      spec.State,
		Status:     spec.Status,
		Created:    spec.Created,
		Updated:    spec.Updated,
		Containers: make(map[string]v1.ContainerSpec),
	}

	for _, c := range spec.Containers {
		s.Containers[c.Meta.ID] = v1.ToContainerSpec(c)
	}

	return s
}

func ToPodState(state types.PodState) PodState {
	return PodState{
		State:     state.State,
		Status:    state.Status,
		Provision: state.Provision,
		Ready:     state.Ready,
	}
}

func FromPodMeta(meta PodMeta) types.PodMeta {
	m := types.PodMeta{}
	m.Name = meta.Name
	m.Labels = meta.Labels
	m.Created = meta.Created
	m.Updated = meta.Updated
	return m
}

func FromPodSpec(spec PodSpec) types.PodSpec {
	s := types.PodSpec{
		ID:         spec.ID,
		State:      spec.State,
		Status:     spec.Status,
		Containers: make(map[string]*types.ContainerSpec),
		Created:    spec.Created,
		Updated:    spec.Updated,
	}

	for _, c := range spec.Containers {
		s.Containers[c.Meta.ID] = v1.FromContainerSpec(c)
	}

	return s
}

func FromPodState(state PodState) types.PodState {
	return types.PodState{
		State:     state.State,
		Status:    state.Status,
		Provision: state.Provision,
		Ready:     state.Ready,
	}
}
