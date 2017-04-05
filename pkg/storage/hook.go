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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"golang.org/x/net/context"
)

const HookTable string = "hooks"

// Service Build type for interface in interfaces folder
type HookStorage struct {
	IHook
	Helper IHelper
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get hooks by image
func (s *HookStorage) GetByToken(ctx context.Context, token string) (*types.Hook, error) {
	return nil, nil
}

// Get hooks by image
func (s *HookStorage) ListByUser(ctx context.Context, id string) (*types.HookList, error) {
	return nil, nil
}

// Get hooks by image
func (s *HookStorage) ListByImage(ctx context.Context, user, id string) (*types.HookList, error) {
	return nil, nil
}

// Get hooks by service
func (s *HookStorage) ListByService(ctx context.Context, user, id string) (*types.HookList, error) {
	return nil, nil
}

// Insert new hook into storage
func (s *HookStorage) Insert(ctx context.Context, hook *types.Hook) (*types.Hook, error) {
	return nil, nil
}

// Remove  hook by service id from storage
func (s *HookStorage) RemoveByService(ctx context.Context, id string) error {

	return nil
}

func NewHookStorage(config store.Config, helper IHelper) *HookStorage {
	s := new(HookStorage)
	s.Helper = helper
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
