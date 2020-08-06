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
	"context"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
)

const (
	logLevel  = 3
	logPrefix = "state"
)

// State main structure
type State struct {
	storage storage.IStorage

	Namespace *NamespaceController
	Resources map[string]ResourceController
}

// Resource wrapper for namespace
func (s *State) Resource(kind string) ResourceController {

	if c, ok := s.Resources[kind]; ok {
		return c
	}

	return nil
}

// NewState function returns new instance of state
func NewState(ctx context.Context, stg storage.IStorage) *State {

	var state = new(State)

	state.storage = stg

	state.Namespace = NewNamespaceController(ctx)
	state.Resources = make(map[string]ResourceController, 0)

	state.Resources[models.KindService] = NewServiceController(ctx)

	return state
}
