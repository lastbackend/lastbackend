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
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"fmt"
	"strings"
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

type VolumeStorage struct {
	storage.Volume
	data map[string]*types.Volume
}

// Get volume by name
func (s *VolumeStorage) Get(ctx context.Context, name string) (*types.Volume, error) {
	if ns, ok := s.data[name]; ok {
		return ns, nil
	}
	return nil, errors.New(store.ErrEntityNotFound)
}

// Get volumes by namespace name
func (s *VolumeStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Volume, error) {
	list := make(map[string]*types.Volume, 0)

	prefix := fmt.Sprintf("%s:", namespace)
	for _, d := range s.data {

		if strings.HasPrefix(d.Meta.Name, prefix) {
			list[d.Meta.Name] = d
		}
	}

	return list, nil
}

// Update volume state
func (s *VolumeStorage) SetState(ctx context.Context, volume *types.Volume) error {
	if err := s.checkVolumeExists(volume); err != nil {
		return err
	}

	s.data[volume.Meta.Name].State = volume.State
	return nil
}

// Insert new volume
func (s *VolumeStorage) Insert(ctx context.Context, volume *types.Volume) error {

	if err := s.checkVolumeArgument(volume); err != nil {
		return err
	}

	s.data[volume.Meta.Name] = volume

	return nil
}

// Update volume info
func (s *VolumeStorage) Update(ctx context.Context, volume *types.Volume) error {

	if err := s.checkVolumeExists(volume); err != nil {
		return err
	}

	s.data[volume.Meta.Name] = volume

	return nil
}

// Remove volume from storage
func (s *VolumeStorage) Remove(ctx context.Context, volume *types.Volume) error {

	if err := s.checkVolumeExists(volume); err != nil {
		return err
	}

	delete(s.data, volume.Meta.Name)

	return nil
}

// Watch volume changes
func (s *VolumeStorage) Watch(ctx context.Context, volume chan *types.Volume) error {
	return nil
}

// Watch volume spec changes
func (s *VolumeStorage) WatchSpec(ctx context.Context, volume chan *types.Volume) error {
	return nil
}

// newVolumeStorage returns new storage
func newVolumeStorage() *VolumeStorage {
	s := new(VolumeStorage)
	s.data = make(map[string]*types.Volume)
	return s
}

// checkVolumeArgument - check if argument is valid for manipulations
func (s *VolumeStorage) checkVolumeArgument(volume *types.Volume) error {

	if volume == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if volume.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

// checkVolumeArgument - check if volume exists in store
func (s *VolumeStorage) checkVolumeExists(volume *types.Volume) error {

	if err := s.checkVolumeArgument(volume); err != nil {
		return err
	}

	if _, ok := s.data[volume.Meta.Name]; !ok {
		return errors.New(store.ErrEntityNotFound)
	}

	return nil
}

