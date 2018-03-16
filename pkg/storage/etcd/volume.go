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
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const volumeStorage string = "volumes"

// Volume Service type for interface in interfaces folder
type VolumeStorage struct {
	storage.Volume
}

func (s *VolumeStorage) GetByID(ctx context.Context, id string) (*types.Volume, error) {

	log.V(logLevel).Debugf("Storage: Volume: get by id: %s", id)

	if len(id) == 0 {
		err := errors.New("id can not be empty")
		log.V(logLevel).Errorf("Storage: Volume: get volume err: %s", err.Error())
		return nil, err
	}

	return nil, nil
}

func (s *VolumeStorage) ListByProject(ctx context.Context, id string) ([]*types.Volume, error) {
	volumes := []*types.Volume{}
	return volumes, nil
}

// Insert new volume into storage
func (s *VolumeStorage) Insert(ctx context.Context, volume *types.Volume) error {
	return nil
}

// Remove build model
func (s *VolumeStorage) Remove(ctx context.Context, id string) error {
	return nil
}

func newVolumeStorage(config store.Config) *VolumeStorage {
	s := new(VolumeStorage)
	return s
}
