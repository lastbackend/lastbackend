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

package state

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

const logLevel = 3

type State struct {
	node      *NodeState
	pods      *PodState
	networks  *NetworkState
	volumes   *VolumesState
	secrets   *SecretsState
	endpoints *EndpointState
	task      *TaskState
}

func (s *State) Node() *NodeState {
	return s.node
}

func (s *State) Pods() *PodState {
	return s.pods
}

func (s *State) Networks() *NetworkState {
	return s.networks
}

func (s *State) Volumes() *VolumesState {
	return s.volumes
}

func (s *State) Secrets() *SecretsState {
	return s.secrets
}

func (s *State) Endpoints() *EndpointState {
	return s.endpoints
}

func (s *State) Tasks() *TaskState {
	return s.task
}

type NodeState struct {
	Info   types.NodeInfo
	Status types.NodeStatus
}

func New() *State {

	state := State{
		node: new(NodeState),
		pods: &PodState{
			containers: make(map[string]*types.PodContainer, 0),
			pods:       make(map[string]*types.PodStatus, 0),
		},
		networks: &NetworkState{
			subnets: make(map[string]types.NetworkSpec, 0),
		},
		volumes: &VolumesState{
			volumes: make(map[string]types.VolumeSpec, 0),
		},
		secrets: &SecretsState{
			secrets: make(map[string]types.Secret, 0),
		},
		endpoints: &EndpointState{
			endpoints: make(map[string]*types.EndpointStatus, 0),
		},
		task: &TaskState{
			tasks: make(map[string]types.NodeTask, 0),
		},
	}

	return &state
}
