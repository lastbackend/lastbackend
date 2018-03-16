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

package etcd

import (
	"context"
	"errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
)

const hookStorage string = "hooks"

// Service Hook type for interface in interfaces folder
type HookStorage struct {
	storage.Hook
}

// Get hooks by id
func (s *HookStorage) Get(ctx context.Context, id string) (*types.Hook, error) {

	log.V(logLevel).Debugf("Storage: Hook: get hook by id: %s", id)

	if len(id) == 0 {
		err := errors.New("id can not be empty")
		log.V(logLevel).Errorf("Storage: Hook: get hook by id err: %s", err.Error())
		return nil, err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Hook: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	hook := new(types.Hook)
	keyMeta := keyCreate(hookStorage, id)
	if err := client.Get(ctx, keyMeta, &hook); err != nil {
		log.V(logLevel).Errorf("Storage: Hook: get hook meta err: %s", err.Error())
		return nil, err
	}

	return hook, nil
}

// Insert new hook into storage
func (s *HookStorage) Insert(ctx context.Context, hook *types.Hook) error {

	log.V(logLevel).Debugf("Storage: Hook: create hook: %#v", hook)

	if hook == nil {
		err := errors.New("hook can not be nil")
		log.V(logLevel).Errorf("Storage: Hook: create hook err: %s", err.Error())
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Hook: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(hookStorage, hook.Meta.Name)
	if err := client.Create(ctx, key, hook, nil, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Hook: create hook err: %s", err.Error())
		return err
	}

	return nil
}

// Remove hook by id from storage
func (s *HookStorage) Remove(ctx context.Context, id string) error {

	log.V(logLevel).Debugf("Storage: Hook: remove hook by id %#v", id)

	if len(id) == 0 {
		err := errors.New("id can not be nil")
		log.V(logLevel).Errorf("Storage: Hook: remove hook err: %s", err.Error())
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Hook: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(hookStorage, id)
	if err := client.DeleteDir(ctx, key); err != nil {
		log.V(logLevel).Errorf("Storage: Hook: remove hook err: %s", err.Error())
		return err
	}
	return nil
}

func newHookStorage() *HookStorage {
	s := new(HookStorage)
	return s
}
