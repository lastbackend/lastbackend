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

package storage

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const hookStorage string = "hooks"

// Service Build type for interface in interfaces folder
type HookStorage struct {
	IHook
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get hooks by image
func (s *HookStorage) GetByToken(ctx context.Context, token string) (*types.Hook, error) {
	var hook = new(types.Hook)
	return hook, nil
}

// Get hooks by image
func (s *HookStorage) List(ctx context.Context, id string) ([]*types.Hook, error) {
	var hooks = []*types.Hook{}
	return hooks, nil
}

// Get hooks by image
func (s *HookStorage) ListByImage(ctx context.Context, id string) ([]*types.Hook, error) {
	var hooks = []*types.Hook{}
	return hooks, nil
}

// Get hooks by service
func (s *HookStorage) ListByService(ctx context.Context, id string) ([]*types.Hook, error) {
	var hooks = []*types.Hook{}
	return hooks, nil
}

// Insert new hook into storage
func (s *HookStorage) Insert(ctx context.Context, hook *types.Hook) error {
	return nil
}

// Remove  hook by service id from storage
func (s *HookStorage) RemoveByService(ctx context.Context, id string) error {

	return nil
}

func newHookStorage(config store.Config, util IUtil) *HookStorage {
	s := new(HookStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
