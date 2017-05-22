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

// Service Hook type for interface in interfaces folder
type HookStorage struct {
	IHook
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get hooks by id
func (s *HookStorage) Get(ctx context.Context, id string) (*types.Hook, error) {
	var hook = new(types.Hook)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	keyMeta := s.util.Key(ctx, hookStorage, id)
	if err := client.Get(ctx, keyMeta, &hook); err != nil {
		return nil, err
	}

	return hook, nil
}

// Insert new hook into storage
func (s *HookStorage) Insert(ctx context.Context, hook *types.Hook) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	key := s.util.Key(ctx, hookStorage, hook.Meta.ID)
	if err := client.Create(ctx, key, hook, nil, 0); err != nil {
		return err
	}

	return nil
}

// Remove hook by id from storage
func (s *HookStorage) Remove(ctx context.Context, id string) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	key := keyCreate(hookStorage, id)
	client.DeleteDir(ctx, key)

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
