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
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"strings"
)

type VolumeStorage struct {
	storage.Volume
	data map[string]*types.Volume
}

// Get volume by name
func (s *VolumeStorage) Get(ctx context.Context, namespace, name string) (*types.Volume, error) {

	if ns, ok := s.data[s.keyCreate(namespace, name)]; ok {
		return ns, nil
	}
	return nil, errors.New(store.ErrEntityNotFound)
}

// Get volumes by namespace name
func (s *VolumeStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Volume, error) {
	list := make(map[string]*types.Volume, 0)

	prefix := fmt.Sprintf("%s:", namespace)
	for _, d := range s.data {

		if strings.HasPrefix(s.keyGet(d), prefix) {
			list[s.keyGet(d)] = d
		}
	}

	return list, nil
}

// Update volume state
func (s *VolumeStorage) SetStatus(ctx context.Context, volume *types.Volume) error {
	if err := s.checkVolumeExists(volume); err != nil {
		return err
	}

	s.data[s.keyGet(volume)].Status = volume.Status
	return nil
}

// Update volume state
func (s *VolumeStorage) SetSpec(ctx context.Context, volume *types.Volume) error {
	if err := s.checkVolumeExists(volume); err != nil {
		return err
	}

	s.data[s.keyGet(volume)].Spec = volume.Spec
	return nil
}

// Insert new volume
func (s *VolumeStorage) Insert(ctx context.Context, volume *types.Volume) error {

	if err := s.checkVolumeArgument(volume); err != nil {
		return err
	}

	s.data[s.keyGet(volume)] = volume

	return nil
}

// Update volume info
func (s *VolumeStorage) Update(ctx context.Context, volume *types.Volume) error {

	if err := s.checkVolumeExists(volume); err != nil {
		return err
	}

	s.data[s.keyGet(volume)] = volume

	return nil
}

// Remove volume from storage
func (s *VolumeStorage) Remove(ctx context.Context, volume *types.Volume) error {

	if err := s.checkVolumeExists(volume); err != nil {
		return err
	}

	delete(s.data, s.keyGet(volume))

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

// Watch volume state changes
func (s *VolumeStorage) WatchStatus(ctx context.Context, volume chan *types.Volume) error {
	return nil
}

// Clear volume storage
func (s *VolumeStorage) Clear(ctx context.Context) error {
	s.data = make(map[string]*types.Volume)
	return nil
}

// keyCreate util function
func (s *VolumeStorage) keyCreate(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

// keyGet util function
func (s *VolumeStorage) keyGet(v *types.Volume) string {
	return v.SelfLink()
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

	if _, ok := s.data[s.keyGet(volume)]; !ok {
		return errors.New(store.ErrEntityNotFound)
	}

	return nil
}
