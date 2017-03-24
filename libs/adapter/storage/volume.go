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

const VolumeTable string = "volumes"

// Volume Service type for interface in interfaces folder
type VolumeStorage struct {
	storage.IVolume
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *VolumeStorage) GetByID(user, id string) (*model.Volume, error) {
	return nil, nil
}

func (s *VolumeStorage) ListByProject(id string) (*model.VolumeList, error) {
	return nil, nil
}

// Insert new volume into storage
func (s *VolumeStorage) Insert(volume *model.Volume) (*model.Volume, error) {
	return nil, nil
}

// Remove build model
func (s *VolumeStorage) Remove(id string) error {
	return nil
}

func newVolumeStorage(config store.Config) *VolumeStorage {
	s := new(VolumeStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return db.Create(config)
	}
	return s
}
