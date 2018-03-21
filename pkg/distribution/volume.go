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

package distribution

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

type IVolume interface {
	Get(namespace, volume string) (*types.Volume, error)
	List(namespace string) (map[string]*types.Volume, error)
	Create(namespace *types.Namespace, opts *types.VolumeCreateOptions) (*types.Volume, error)
	Update(volume *types.Volume, opts *types.VolumeUpdateOptions) (*types.Volume, error)
	Destroy(volume *types.Volume) error
	Remove(volume *types.Volume) error
	SetState(volume *types.Volume) error
	SetStatus(volume *types.Volume, status *types.VolumeStatus) error
	Watch(dt chan *types.Volume) error
	WatchSpec(dt chan *types.Volume) error
}

type Volume struct {
	context context.Context
	storage storage.Storage
}

func (v *Volume) Get(namespace, volume string) (*types.Volume, error) {
	return nil, nil
}

func (v *Volume) List(namespace string) (map[string]*types.Volume, error) {
	return nil, nil
}

func (v *Volume) Create(namespace *types.Namespace, opts *types.VolumeCreateOptions) (*types.Volume, error) {
	return nil, nil
}

func (v *Volume) Update(volume *types.Volume, opts *types.VolumeUpdateOptions) (*types.Volume, error) {
	return nil, nil
}

func (v *Volume) Destroy(volume *types.Volume) error {
	return nil
}

func (v *Volume) Remove(volume *types.Volume) error {
	return nil
}

func (v *Volume) SetState(volume *types.Volume) error {
	return nil
}

func (v *Volume) SetStatus(volume *types.Volume, status *types.VolumeStatus) error {
	return nil
}

func (v *Volume) Watch(dt chan *types.Volume) error {
	return nil
}

func (v *Volume) WatchSpec(dt chan *types.Volume) error {
	return nil
}

func NewVolumeModel(ctx context.Context, stg storage.Storage) IVolume {
	return &Volume{ctx, stg}
}
