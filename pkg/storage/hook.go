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
	"errors"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const hookStorage string = "hooks"

// Service Hook type for interface in interfaces folder
type HookStorage struct {
	IHook
	log    logger.ILogger
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get hooks by id
func (s *HookStorage) Get(ctx context.Context, id string) (*types.Hook, error) {

	s.log.V(debugLevel).Debugf("Storage: Hook: get hook by id: %s", id)

	if len(id) == 0 {
		err := errors.New("id can not be empty")
		s.log.V(debugLevel).Errorf("Storage: Hook: get hook by id err: %s", err.Error())
		return nil, err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(debugLevel).Errorf("Storage: Hook: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	hook := new(types.Hook)
	keyMeta := s.util.Key(ctx, hookStorage, id)
	if err := client.Get(ctx, keyMeta, &hook); err != nil {
		s.log.V(debugLevel).Errorf("Storage: Hook: get hook meta err: %s", err.Error())
		return nil, err
	}

	return hook, nil
}

// Insert new hook into storage
func (s *HookStorage) Insert(ctx context.Context, hook *types.Hook) error {

	s.log.V(debugLevel).Debugf("Storage: Hook: create hook: %#v", hook)

	if hook == nil {
		err := errors.New("hook can not be nil")
		s.log.V(debugLevel).Errorf("Storage: Hook: create hook err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(debugLevel).Errorf("Storage: Hook: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := s.util.Key(ctx, hookStorage, hook.Meta.ID)
	if err := client.Create(ctx, key, hook, nil, 0); err != nil {
		s.log.V(debugLevel).Errorf("Storage: Hook: create hook err: %s", err.Error())
		return err
	}

	return nil
}

// Remove hook by id from storage
func (s *HookStorage) Remove(ctx context.Context, id string) error {

	s.log.V(debugLevel).Debugf("Storage: Hook: remove hook by id %#v", id)

	if len(id) == 0 {
		err := errors.New("id can not be nil")
		s.log.V(debugLevel).Errorf("Storage: Hook: remove hook err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(debugLevel).Errorf("Storage: Hook: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(hookStorage, id)
	if err := client.DeleteDir(ctx, key); err != nil {
		s.log.V(debugLevel).Errorf("Storage: Hook: remove hook err: %s", err.Error())
		return err
	}
	return nil
}

func newHookStorage(config store.Config, log logger.ILogger, util IUtil) *HookStorage {
	s := new(HookStorage)
	s.log = log
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config, log)
	}
	return s
}
