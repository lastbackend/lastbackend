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
	images    *ImageState
	networks  *NetworkState
	volumes   *VolumesState
	secrets   *SecretsState
	endpoints *EndpointState
	task      *TaskState
	configs   *ConfigState
}

func (s *State) Node() *NodeState {
	return s.node
}

func (s *State) Pods() *PodState {
	return s.pods
}

func (s *State) Images() *ImageState {
	return s.images
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

func (s *State) Configs() *ConfigState {
	return s.configs
}

type NodeState struct {
	Info   types.NodeInfo
	Status types.NodeStatus
}

func New() *State {

	state := State{
		node: new(NodeState),
		pods: &PodState{
			local:      make(map[string]bool),
			containers: make(map[string]*types.PodContainer, 0),
			pods:       make(map[string]*types.PodStatus, 0),
			watchers:   make(map[chan string]bool, 0),
		},
		images: &ImageState{
			images: make(map[string]*types.Image, 0),
		},
		networks: &NetworkState{
			subnets: make(map[string]types.NetworkState, 0),
		},
		volumes: &VolumesState{
			volumes:  make(map[string]types.VolumeStatus, 0),
			local:    make(map[string]bool),
			watchers: make(map[chan string]bool, 0),
		},
		secrets: &SecretsState{
			secrets: make(map[string]types.Secret, 0),
		},
		endpoints: &EndpointState{
			endpoints: make(map[string]*types.EndpointState, 0),
		},
		task: &TaskState{
			tasks: make(map[string]types.NodeTask, 0),
		},
		configs: &ConfigState{
			configs: make(map[string]*types.ConfigManifest, 0),
		},
	}

	return &state
}
