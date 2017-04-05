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

const VolumeTable string = "volumes"

// Volume Service type for interface in interfaces folder
type VolumeStorage struct {
	IVolume
	Helper IHelper
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *VolumeStorage) GetByID(ctx context.Context, user, id string) (*types.Volume, error) {
	return nil, nil
}

func (s *VolumeStorage) ListByProject(ctx context.Context, id string) (*types.VolumeList, error) {
	return nil, nil
}

// Insert new volume into storage
func (s *VolumeStorage) Insert(ctx context.Context, volume *types.Volume) (*types.Volume, error) {
	return nil, nil
}

// Remove build model
func (s *VolumeStorage) Remove(ctx context.Context, id string) error {
	return nil
}

func NewVolumeStorage(config store.Config, helper IHelper) *VolumeStorage {
	s := new(VolumeStorage)
	s.Helper = helper
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
