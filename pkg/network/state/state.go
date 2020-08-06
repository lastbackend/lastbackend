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

package state

import (
	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

const logLevel = 3

type State struct {
	subnets   *SubnetState
	endpoints *EndpointState
	resolvers *ResolverState
}

func (s *State) Subnets() *SubnetState {
	return s.subnets
}

func (s *State) Endpoints() *EndpointState {
	return s.endpoints
}

func (s *State) Resolvers() *ResolverState {
	return s.resolvers
}

func New() *State {

	state := State{
		subnets: &SubnetState{
			subnets: make(map[string]models.NetworkState, 0),
		},
		endpoints: &EndpointState{
			endpoints: make(map[string]*models.EndpointState, 0),
		},
		resolvers: &ResolverState{
			resolvers: make(map[string]*models.ResolverManifest, 0),
		},
	}

	return &state
}
