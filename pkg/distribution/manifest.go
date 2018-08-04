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
)

const (
	logManifestPrefix = "distribution:manifest"
)

type Manifest struct {
	context context.Context
	storage storage.Storage
}

func (m *Manifest) PodManifestMap(node string) (*types.PodManifestMap, error) {
	log.Debugf("%s:PodManifestMap:> ", logManifestPrefix)

	var (
		mf = types.NewPodManifestMap()
		qs = m.storage.Filter().Manifest().ByKindManifest(node, storage.PodKind)
	)

	if err := m.storage.Map(m.context, storage.ManifestKind, qs, mf, nil); err != nil {
		log.Errorf("%s:PodManifestMap:> err :%s", logManifestPrefix, err.Error())
		return nil, err
	}

	return mf, nil
}

func (m *Manifest) PodManifestGet(node, pod string) (*types.PodManifest, error) {
	log.Debugf("%s:PodManifestGet:> ", logManifestPrefix)

	var (
		mf = new(types.PodManifest)
		k  = m.storage.Key().Manifest(node, storage.PodKind, pod)
	)

	if err := m.storage.Get(m.context, storage.ManifestKind, k, &mf, nil); err != nil {
		log.Errorf("%s:PodManifestMap:> err :%s", logManifestPrefix, err.Error())

		if errors.Storage().IsErrEntityNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return mf, nil
}

func (m *Manifest) PodManifestAdd(node, pod string, manifest *types.PodManifest) error {
	log.Debugf("%s:PodManifestAdd:> ", logManifestPrefix)

	var (
		k = m.storage.Key().Manifest(node, storage.PodKind, pod)
	)

	if err := m.storage.Put(m.context, storage.ManifestKind, k, manifest, nil); err != nil {
		log.Errorf("%s:PodManifestAdd:> err :%s", logManifestPrefix, err.Error())
		return err
	}

	return nil
}

func (m *Manifest) PodManifestSet(node, pod string, manifest *types.PodManifest) error {
	log.Debugf("%s:PodManifestSet:> ", logManifestPrefix)

	var (
		k = m.storage.Key().Manifest(node, storage.PodKind, pod)
	)

	if err := m.storage.Set(m.context, storage.ManifestKind, k, manifest, nil); err != nil {
		log.Errorf("%s:PodManifestSet:> err :%s", logManifestPrefix, err.Error())
		return err
	}

	return nil
}

func (m *Manifest) PodManifestDel(node, pod string) error {
	log.Debugf("%s:PodManifestDel:> ", logManifestPrefix)

	var (
		k = m.storage.Key().Manifest(node, storage.PodKind, pod)
	)

	if err := m.storage.Del(m.context, storage.ManifestKind, k); err != nil {
		log.Errorf("%s:PodManifestDel:> err :%s", logManifestPrefix, err.Error())
		return err
	}

	return nil
}

func (m *Manifest) VolumeManifestMap(node string) (*types.VolumeManifestMap, error) {
	log.Debugf("%s:VolumeManifestMap:> ", logManifestPrefix)

	var (
		mf = types.NewVolumeManifestMap()
		qs = m.storage.Filter().Manifest().ByKindManifest(node, storage.VolumeKind)
	)

	if err := m.storage.Map(m.context, storage.ManifestKind, qs, mf, nil); err != nil {
		log.Errorf("%s:VolumeManifestMap:> err :%s", logManifestPrefix, err.Error())
		return nil, err
	}
	return mf, nil
}

func (m *Manifest) VolumeManifestGet(node, volume string) (*types.VolumeManifest, error) {
	log.Debugf("%s:VolumeManifestGet:> ", logManifestPrefix)

	var (
		mf = new(types.VolumeManifest)
		k  = m.storage.Key().Manifest(node, storage.VolumeKind, volume)
	)

	if err := m.storage.Get(m.context, storage.ManifestKind, k, &mf, nil); err != nil {
		log.Errorf("%s:VolumeManifestGet:> err :%s", logManifestPrefix, err.Error())

		if errors.Storage().IsErrEntityNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return mf, nil
}

func (m *Manifest) VolumeManifestAdd(node, volume string, manifest *types.VolumeManifest) error {
	log.Debugf("%s:VolumeManifestAdd:> ", logManifestPrefix)

	var (
		k = m.storage.Key().Manifest(node, storage.VolumeKind, volume)
	)

	if err := m.storage.Put(m.context, storage.ManifestKind, k, manifest, nil); err != nil {
		log.Errorf("%s:VolumeManifestAdd:> err :%s", logManifestPrefix, err.Error())
		return err
	}

	return nil
}

func (m *Manifest) VolumeManifestSet(node, volume string, manifest *types.VolumeManifest) error {
	log.Debugf("%s:VolumeManifestSet:> ", logManifestPrefix)

	var (
		k = m.storage.Key().Manifest(node, storage.VolumeKind, volume)
	)

	if err := m.storage.Set(m.context, storage.ManifestKind, k, manifest, nil); err != nil {
		log.Errorf("%s:VolumeManifestSet:> err :%s", logManifestPrefix, err.Error())
		return err
	}

	return nil
}

func (m *Manifest) VolumeManifestDel(node, volume string) error {
	log.Debugf("%s:DelVolumeManifest:> ", logManifestPrefix)

	var (
		k = m.storage.Key().Manifest(node, storage.VolumeKind, volume)
	)

	if err := m.storage.Del(m.context, storage.ManifestKind, k); err != nil {
		log.Errorf("%s:PodManifestDel:> err :%s", logManifestPrefix, err.Error())
		return err
	}

	return nil
}

func NewManifestModel(ctx context.Context, stg storage.Storage) *Manifest {
	return &Manifest{ctx, stg}
}
