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

package distribution

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
)

const (
	logVolumePrefix = "distribution:volume"
)

type IVolume interface {
	Get(namespace, volume string) (*types.Volume, error)
	ListByNamespace(namespace string) (map[string]*types.Volume, error)
	Create(namespace *types.Namespace, opts *types.VolumeCreateOptions) (*types.Volume, error)
	Update(volume *types.Volume, opts *types.VolumeUpdateOptions) (*types.Volume, error)
	Destroy(volume *types.Volume) error
	Remove(volume *types.Volume) error
	SetStatus(volume *types.Volume, status *types.VolumeStatus) error
	Watch(dt chan *types.Volume) error
	WatchSpec(dt chan *types.Volume) error
}

type Volume struct {
	context context.Context
	storage storage.Storage
}

func (v *Volume) Get(namespace, name string) (*types.Volume, error) {
	log.V(logLevel).Debugf("%s:get:> get volume by id %s/%s", logVolumePrefix, namespace, name)

	item, err := v.storage.Volume().Get(v.context, namespace, name)
	if err != nil {

		if err.Error() == store.ErrEntityNotFound {
			log.V(logLevel).Warnf("%s:get:> in namespace %s by name %s not found", logVolumePrefix, namespace, name)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> in namespace %s by name %s error: %s", logVolumePrefix, namespace, name, err.Error())
		return nil, err
	}

	return item, nil
}

func (v *Volume) ListByNamespace(namespace string) (map[string]*types.Volume, error) {
	log.V(logLevel).Debugf("%s:list:> get volumes list", logVolumePrefix)

	items, err := v.storage.Volume().ListByNamespace(v.context, namespace)
	if err != nil {
		log.V(logLevel).Error("%s:list:> get volumes list err: %s", logVolumePrefix, err.Error())
		return items, err
	}

	log.V(logLevel).Debugf("%s:list:> get volumes list result: %d", logVolumePrefix, len(items))

	return items, nil
}

func (v *Volume) Create(namespace *types.Namespace, opts *types.VolumeCreateOptions) (*types.Volume, error) {
	log.V(logLevel).Debugf("%s:crete:> create volume %#v", logVolumePrefix, opts)

	volume := new(types.Volume)
	volume.Meta.SetDefault()
	volume.Meta.Name = generator.GenerateRandomString(10)
	volume.Meta.Namespace = namespace.Meta.Name
	volume.Status.Stage = types.StateInitialized

	if err := v.storage.Volume().Insert(v.context, volume); err != nil {
		log.V(logLevel).Errorf("%s:crete:> insert volume err: %s", logVolumePrefix, err.Error())
		return nil, err
	}

	return volume, nil
}

func (v *Volume) Update(volume *types.Volume, opts *types.VolumeUpdateOptions) (*types.Volume, error) {
	log.V(logLevel).Debugf("%s:update:> update volume %s", logVolumePrefix, volume.Meta.Name)

	volume.Meta.SetDefault()
	volume.Status.Stage = types.StateProvision

	if err := v.storage.Volume().Update(v.context, volume); err != nil {
		log.V(logLevel).Errorf("%s:update:> update volume err: %s", logVolumePrefix, err.Error())
		return nil, err
	}

	return volume, nil
}

func (v *Volume) Destroy(volume *types.Volume) error {

	if volume == nil {
		log.V(logLevel).Warnf("%s:destroy:> invalid argument %v", logVolumePrefix, volume)
		return nil
	}

	log.V(logLevel).Debugf("%s:destroy:> volume %s", logVolumePrefix, volume.Meta.Name)

	volume.Spec.State.Destroy = true
	if err := v.storage.Volume().SetSpec(v.context, volume); err != nil {
		log.Errorf("%s:destroy:> volume err: %s", logVolumePrefix, err.Error())
		return err
	}

	return nil
}

func (v *Volume) Remove(volume *types.Volume) error {
	log.V(logLevel).Debugf("%s:remove:> remove volume %#v", logVolumePrefix, volume)

	if err := v.storage.Volume().Remove(v.context, volume); err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove volume  err: %s", logVolumePrefix, err.Error())
		return err
	}

	return nil
}

func (v *Volume) SetStatus(volume *types.Volume, status *types.VolumeStatus) error {
	if volume == nil {
		log.V(logLevel).Warnf("%s:setstatus:> invalid argument %v", logVolumePrefix, volume)
		return nil
	}

	log.V(logLevel).Debugf("%s:setstatus:> set state volume %s -> %#v", logVolumePrefix, volume.Meta.Name, status)

	volume.Status = *status
	if err := v.storage.Volume().SetStatus(v.context, volume); err != nil {
		log.Errorf("%s:setstatus:> pod set status err: %s", err.Error())
		return err
	}

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
