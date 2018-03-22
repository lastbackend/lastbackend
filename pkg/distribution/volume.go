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
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
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
	log.V(logLevel).Debug("api:distribution:volume: get volume by id %s/%s", namespace, name)

	item, err := v.storage.Volume().Get(v.context, namespace, name)
	if err != nil {

		if err.Error() == store.ErrEntityNotFound {
			log.V(logLevel).Warnf("api:distribution:volume:get: in namespace %s by name %s not found", namespace, name)
			return nil, nil
		}

		log.V(logLevel).Errorf("api:distribution:volume:get: in namespace %s by name %s error: %s", namespace, name, err)
		return nil, err
	}

	return item, nil
}

func (v *Volume) ListByNamespace(namespace string) (map[string]*types.Volume, error) {
	log.V(logLevel).Debug("api:distribution:volume: list volume")

	items, err := v.storage.Volume().ListByNamespace(v.context, namespace)
	if err != nil {
		log.V(logLevel).Error("api:distribution:volume: list volume err: %s", err)
		return items, err
	}

	log.V(logLevel).Debugf("api:distribution:volume: list volume result: %d", len(items))

	return items, nil
}

func (v *Volume) Create(namespace *types.Namespace, opts *types.VolumeCreateOptions) (*types.Volume, error) {
	log.V(logLevel).Debugf("api:distribution:volume:crete create volume %#v", opts)

	volume := new(types.Volume)
	volume.Meta.SetDefault()
	volume.Meta.Name = generator.GenerateRandomString(10)
	volume.Meta.Namespace = namespace.Meta.Name
	volume.Status.Stage = types.StageInitialized

	if err := v.storage.Volume().Insert(v.context, volume); err != nil {
		log.V(logLevel).Errorf("api:distribution:volume:crete insert volume err: %s", err)
		return nil, err
	}

	return volume, nil
}

func (v *Volume) Update(volume *types.Volume, opts *types.VolumeUpdateOptions) (*types.Volume, error) {
	log.V(logLevel).Debugf("api:distribution:volume:update update volume %s", volume.Meta.Name)

	volume.Meta.SetDefault()
	volume.Status.Stage = types.StageProvision

	if err := v.storage.Volume().Update(v.context, volume); err != nil {
		log.V(logLevel).Errorf("api:distribution:volume:update update volume err: %s", err)
		return nil, err
	}

	return volume, nil
}

func (v *Volume) Destroy(volume *types.Volume) error {

	if volume == nil {
		log.V(logLevel).Warnf("api:distribution:volume:destroy: invalid argument %v", volume)
		return nil
	}

	log.V(logLevel).Debugf("api:distribution:volume:destroy volume %s", volume.Meta.Name)

	volume.Spec.State.Destroy = true
	if err := v.storage.Volume().SetSpec(v.context, volume); err != nil {
		log.Errorf("api:distribution:volume:destroy volume err: %s", err.Error())
		return err
	}

	return nil
}

func (v *Volume) Remove(volume *types.Volume) error {
	log.V(logLevel).Debugf("api:distribution:volume:remove remove volume %#v", volume)

	if err := v.storage.Volume().Remove(v.context, volume); err != nil {
		log.V(logLevel).Errorf("api:distribution:volume:remove remove volume  err: %s", err)
		return err
	}

	return nil
}

func (v *Volume) SetStatus(volume *types.Volume, status *types.VolumeStatus) error {
	if volume == nil {
		log.V(logLevel).Warnf("api:distribution:volume:setstatus: invalid argument %v", volume)
		return nil
	}

	log.V(logLevel).Debugf("api:distribution:volume:setstate set state volume %s -> %#v", volume.Meta.Name, status)

	volume.Status = *status
	if err := v.storage.Volume().SetStatus(v.context, volume); err != nil {
		log.Errorf("Pod set status err: %s", err.Error())
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
