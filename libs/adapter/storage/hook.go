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
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	db "github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const HookTable string = "hooks"

// Service Build type for interface in interfaces folder
type HookStorage struct {
	storage.IHook
	Client func() (store.Interface, store.DestroyFunc, error)
}

// Get hooks by image
func (s *HookStorage) GetByToken(token string) (*model.Hook, error) {
	return nil, nil
}

// Get hooks by image
func (s *HookStorage) ListByUser(id string) (*model.HookList, error) {
	return nil, nil
}

// Get hooks by image
func (s *HookStorage) ListByImage(user, id string) (*model.HookList, error) {
	return nil, nil
}

// Get hooks by service
func (s *HookStorage) ListByService(user, id string) (*model.HookList, error) {
	return nil, nil
}

// Insert new hook into storage
func (s *HookStorage) Insert(hook *model.Hook) (*model.Hook, error) {
	return nil, nil
}

// Remove  hook by service id from storage
func (s *HookStorage) RemoveByService(id string) error {

	return nil
}

func newHookStorage(config store.Config) *HookStorage {
	s := new(HookStorage)
	s.Client = func() (store.Interface, store.DestroyFunc, error) {
		return db.Create(config)
	}
	return s
}
