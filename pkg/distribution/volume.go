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

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
)

const (
	logVolumePrefix = "distribution:volume"
)

type Volume struct {
	context context.Context
	storage storage.Storage
}

func (v *Volume) Get(namespace, name string) (*types.Volume, error) {
	log.V(logLevel).Debugf("%s:get:> get volume by id %s/%s", logVolumePrefix, namespace, name)

	item := new(types.Volume)

	err := v.storage.Get(v.context, storage.VolumeKind, v.storage.Key().Volume(namespace, name), &item, nil)
	if err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> in namespace %s by name %s not found", logVolumePrefix, namespace, name)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> in namespace %s by name %s error: %v", logVolumePrefix, namespace, name, err)
		return nil, err
	}

	return item, nil
}

func (v *Volume) ListByNamespace(namespace string) (*types.VolumeList, error) {
	log.V(logLevel).Debugf("%s:list:> get volumes list", logVolumePrefix)

	list := types.NewVolumeList()
	filter := v.storage.Filter().Volume().ByNamespace(namespace)
	err := v.storage.Map(v.context, storage.VolumeKind, filter, list, nil)
	if err != nil {
		log.V(logLevel).Error("%s:list:> get volumes list err: %v", logVolumePrefix, err)
		return list, err
	}

	log.V(logLevel).Debugf("%s:list:> get volumes list result: %d", logVolumePrefix, len(list.Items))

	return list, nil
}

func (v *Volume) Create(namespace *types.Namespace, opts *types.VolumeCreateOptions) (*types.Volume, error) {
	log.V(logLevel).Debugf("%s:crete:> create volume %#v", logVolumePrefix, opts)

	volume := new(types.Volume)
	volume.Meta.SetDefault()
	volume.Meta.Name = generator.GenerateRandomString(10)
	volume.Meta.Namespace = namespace.Meta.Name
	volume.Status.State = types.StateInitialized

	if err := v.storage.Put(v.context, storage.VolumeKind,
		v.storage.Key().Volume(volume.Meta.Namespace, volume.Meta.Name), volume, nil); err != nil {
		log.V(logLevel).Errorf("%s:crete:> insert volume err: %v", logVolumePrefix, err)
		return nil, err
	}

	return volume, nil
}

func (v *Volume) Update(volume *types.Volume, opts *types.VolumeUpdateOptions) (*types.Volume, error) {
	log.V(logLevel).Debugf("%s:update:> update volume %s", logVolumePrefix, volume.Meta.Name)

	volume.Meta.SetDefault()
	volume.Status.State = types.StateProvision

	if err := v.storage.Set(v.context, storage.VolumeKind,
		v.storage.Key().Volume(volume.Meta.Namespace, volume.Meta.Name), volume, nil); err != nil {
		log.V(logLevel).Errorf("%s:update:> update volume err: %v", logVolumePrefix, err)
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

	if err := v.storage.Set(v.context, storage.VolumeKind,
		v.storage.Key().Volume(volume.Meta.Namespace, volume.Meta.Name), volume, nil); err != nil {
		log.Errorf("%s:destroy:> volume err: %v", logVolumePrefix, err)
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

	if err := v.storage.Set(v.context, storage.VolumeKind,
		v.storage.Key().Volume(volume.Meta.Namespace, volume.Meta.Name), volume, nil); err != nil {
		log.Errorf("%s:setstatus:> pod set status err: %v", err)
		return err
	}

	return nil
}

func (v *Volume) Remove(volume *types.Volume) error {
	log.V(logLevel).Debugf("%s:remove:> remove volume %#v", logVolumePrefix, volume)

	if err := v.storage.Del(v.context, storage.VolumeKind,
		v.storage.Key().Volume(volume.Meta.Namespace, volume.Meta.Name)); err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove volume  err: %v", logVolumePrefix, err)
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

func NewVolumeModel(ctx context.Context, stg storage.Storage) *Volume {
	return &Volume{ctx, stg}
}
