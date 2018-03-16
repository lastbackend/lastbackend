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

package mock

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
)

const hookStorage string = "hooks"

// Service Hook type for interface in interfaces folder
type HookStorage struct {
	storage.Hook
}

// Get hooks by id
func (s *HookStorage) Get(ctx context.Context, id string) (*types.Hook, error) {
	return new(types.Hook), nil
}

// Insert new hook into storage
func (s *HookStorage) Insert(ctx context.Context, hook *types.Hook) error {
	return nil
}

// Remove hook by id from storage
func (s *HookStorage) Remove(ctx context.Context, id string) error {
	return nil
}

func newHookStorage() *HookStorage {
	s := new(HookStorage)
	return s
}
