//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package cluster

import (
	"context"
	"github.com/lastbackend/lastbackend/internal/master/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/model"
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logPrefixVolume = "observer:cluster:volume"
)

func volumeObserve(ss *ClusterState, d *types.Volume) error {

	log.V(logLevel).Debugf("%s:> observe start: %s > %s", logPrefixVolume, d.SelfLink().String(), d.Status.State)

	switch d.Status.State {
	case types.StateCreated:
		if err := handleVolumeStateCreated(ss, d); err != nil {
			log.Errorf("%s:> handle volume state create err: %s", logPrefixVolume, err.Error())
			return err
		}
		break
	case types.StateProvision:
		if err := handleVolumeStateProvision(ss, d); err != nil {
			log.Errorf("%s:> handle volume state provision err: %s", logPrefixVolume, err.Error())
			return err
		}
		break
	case types.StateReady:
		if err := handleVolumeStateReady(ss, d); err != nil {
			log.Errorf("%s:> handle volume state ready err: %s", logPrefixVolume, err.Error())
			return err
		}
		break
	case types.StateError:
		if err := handleVolumeStateError(ss, d); err != nil {
			log.Errorf("%s:> handle volume state error err: %s", logPrefixVolume, err.Error())
			return err
		}
		break
	case types.StateDestroy:
		if err := handleVolumeStateDestroy(ss, d); err != nil {
			log.Errorf("%s:> handle volume state destroy err: %s", logPrefixVolume, err.Error())
			return err
		}
		break
	case types.StateDestroyed:
		if err := handleVolumeStateDestroyed(ss, d); err != nil {
			log.Errorf("%s:> handle volume state destroyed err: %s", logPrefixVolume, err.Error())
			return err
		}
		break
	}

	if d.Status.State == types.StateDestroyed {
		delete(ss.volume.list, d.SelfLink().String())
	} else {
		ss.volume.list[d.SelfLink().String()] = d
	}

	log.V(logLevel).Debugf("%s:> observe state: %s > %s", logPrefixVolume, d.SelfLink().String(), d.Status.State)

	if err := clusterStatusState(ss); err != nil {
		return err
	}

	log.V(logLevel).Debugf("%s:> observe finish: %s > %s", logPrefixVolume, d.SelfLink().String(), d.Status.State)

	return nil
}

func handleVolumeStateCreated(cs *ClusterState, v *types.Volume) error {
	log.V(logLevel).Debugf("%s:> handleVolumeStateCreated: %s > %s", logPrefixVolume, v.SelfLink().String(), v.Status.State)

	if err := volumeProvision(cs, v); err != nil {
		return err
	}
	return nil
}

func handleVolumeStateProvision(cs *ClusterState, v *types.Volume) error {
	log.V(logLevel).Debugf("%s:> handleVolumeStateProvision: %s > %s", logPrefixVolume, v.SelfLink().String(), v.Status.State)

	if err := volumeProvision(cs, v); err != nil {
		return err
	}
	return nil
}

func handleVolumeStateReady(cs *ClusterState, v *types.Volume) error {
	log.V(logLevel).Debugf("%s:> handleVolumeStateReady: %s > %s", logPrefixVolume, v.SelfLink().String(), v.Status.State)
	return nil
}

func handleVolumeStateError(cs *ClusterState, v *types.Volume) error {
	log.V(logLevel).Debugf("%s:> handleVolumeStateError: %s > %s", logPrefixVolume, v.SelfLink().String(), v.Status.State)
	return nil
}

func handleVolumeStateDestroy(cs *ClusterState, v *types.Volume) error {
	log.V(logLevel).Debugf("%s:> handleVolumeStateDestroy: %s > %s", logPrefixVolume, v.SelfLink().String(), v.Status.State)

	if err := volumeDestroy(cs, v); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handleVolumeStateDestroyed(cs *ClusterState, v *types.Volume) error {
	log.V(logLevel).Debugf("%s:> handleVolumeStateDestroyed: %s > %s", logPrefixVolume, v.SelfLink().String(), v.Status.State)

	if err := volumeRemove(cs, v); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func volumeUpdate(v *types.Volume, timestamp time.Time) error {

	if timestamp.Before(v.Meta.Updated) {
		vm := model.NewVolumeModel(context.Background(), envs.Get().GetStorage())
		if err := vm.Update(v); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	return nil
}

func volumeProvision(cs *ClusterState, volume *types.Volume) (err error) {

	t := volume.Meta.Updated

	defer func() {
		if err == nil {
			err = volumeUpdate(volume, t)
		}
	}()

	if volume.Meta.Node != types.EmptyString {
		log.Debugf("%s:> volume manifest create: %s", logPrefixVolume, volume.SelfLink())
		vm := model.NewVolumeModel(context.Background(), envs.Get().GetStorage())
		mf, err := vm.ManifestGet(volume.Meta.Node, volume.SelfLink().String())
		if err != nil {
			log.Errorf("%s:> volume manifest create err: %s", logPrefixVolume, err.Error())
			return err
		}

		if mf != nil {
			if !volumeManifestCheckEqual(mf, volume) {
				if err := volumeManifestSet(volume); err != nil {
					log.Errorf("%s:> volume manifest set err: %s", logPrefixVolume, err.Error())
					return err
				}
			}
		} else {
			if err := volumeManifestAdd(volume); err != nil {
				log.Errorf("%s:> volume manifest add err: %s", logPrefixVolume, err.Error())
				return err
			}
		}

		if volume.Status.State != types.StateProvision {
			volume.Status.State = types.StateProvision
			volume.Meta.Updated = time.Now()
		}

		return nil
	}

	if volume.Meta.Node == types.EmptyString {

		log.Debugf("%s:> volume provision > find node: %s", logPrefixVolume, volume.SelfLink().String())

		node, err := cs.VolumeLease(volume)
		if err != nil {
			log.Errorf("%s:> volume manifest lease err: %s", logPrefixVolume, err.Error())
			return err
		}

		if node == nil {
			log.Debugf("%s:> volume provision > node not found: %s", logPrefixVolume, volume.SelfLink().String())
			volume.Status.State = types.StateError
			volume.Status.Message = errors.NodeNotFound
			volume.Meta.Updated = time.Now()
			return nil
		}

		log.Debugf("%s:> volume provision > node: %s found: %s", logPrefixVolume, node.SelfLink().String(), volume.SelfLink().String())

		volume.Meta.Node = node.SelfLink().String()
		volume.Meta.Updated = time.Now()
	}

	if err := volumeManifestAdd(volume); err != nil {
		log.Errorf("%s:> volume manifest set err: %s", logPrefixVolume, err.Error())
		return err
	}

	if volume.Status.State != types.StateProvision {
		volume.Status.State = types.StateProvision
		volume.Meta.Updated = time.Now()
	}

	return nil
}

func volumeDestroy(cs *ClusterState, volume *types.Volume) (err error) {

	t := volume.Meta.Updated

	defer func() {
		if err == nil {
			err = volumeUpdate(volume, t)
		}
	}()

	if volume.Spec.State.Destroy {
		if volume.Meta.Node == types.EmptyString {
			volume.Status.State = types.StateDestroyed
			volume.Meta.Updated = time.Now()
			return nil
		}
	} else {
		volume.Spec.State.Destroy = true
		volume.Status.State = types.StateDestroy
		volume.Meta.Updated = time.Now()
	}

	if volume.Status.State != types.StateDestroy {
		volume.Status.State = types.StateDestroy
		volume.Meta.Updated = time.Now()
	}

	if err = volumeManifestSet(volume); err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			if volume.Meta.Node != types.EmptyString {
				if _, err := cs.VolumeRelease(volume); err != nil {
					if !errors.Storage().IsErrEntityNotFound(err) {
						return err
					}
				}

				return nil
			}

			volume.Status.State = types.StateDestroyed
			volume.Meta.Updated = time.Now()
			return nil
		}

		return err
	}

	if volume.Meta.Node == types.EmptyString {
		volume.Status.State = types.StateDestroyed
		volume.Meta.Updated = time.Now()
	}

	return nil
}

func volumeRemove(cs *ClusterState, volume *types.Volume) (err error) {

	vm := model.NewVolumeModel(context.Background(), envs.Get().GetStorage())
	if _, err = cs.VolumeRelease(volume); err != nil {
		return err
	}

	if err = volumeManifestDel(volume); err != nil {
		return err
	}

	volume.Meta.Node = types.EmptyString
	volume.Meta.Updated = time.Now()

	if err = vm.Remove(volume); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	delete(cs.volume.list, volume.SelfLink().String())
	return nil
}

func volumeManifestAdd(vol *types.Volume) error {

	log.V(logLevel).Debugf("%s: create volume manifest for node: %s", logPrefixVolume, vol.SelfLink().String())

	var vm = new(types.VolumeManifest)

	vm.State.Destroy = false
	vm.Type = vol.Spec.Type
	vm.AccessMode = vol.Spec.AccessMode
	vm.HostPath = vol.Spec.HostPath
	vm.Capacity.Storage = vol.Spec.Capacity.Storage

	im := model.NewVolumeModel(context.Background(), envs.Get().GetStorage())
	if err := im.ManifestAdd(vol.Meta.Node, vol.SelfLink().String(), vm); err != nil {
		log.Errorf("%s:> volume manifest create err: %s", logPrefixVolume, err.Error())
		return err
	}

	return nil
}

func volumeManifestSet(vol *types.Volume) error {

	var (
		m   *types.VolumeManifest
		err error
	)

	im := model.NewVolumeModel(context.Background(), envs.Get().GetStorage())

	m, err = im.ManifestGet(vol.Meta.Node, vol.Meta.SelfLink.String())
	if err != nil {
		return err
	}

	// Update manifest
	if m == nil {
		log.V(logLevel).Debugf("%s: create volume for node: %s", logPrefixVolume, vol.SelfLink().String())
		ms := types.VolumeManifest(vol.Spec)
		m = &ms
	} else {
		*m = types.VolumeManifest(vol.Spec)
	}

	if err := im.ManifestSet(vol.Meta.Node, vol.SelfLink().String(), m); err != nil {
		log.Errorf("can not update volume manifest: %s", err.Error())
	}

	return nil
}

func volumeManifestDel(vol *types.Volume) error {

	if vol.Meta.Node == types.EmptyString {
		return nil
	}

	// Remove manifest
	vm := model.NewVolumeModel(context.Background(), envs.Get().GetStorage())
	err := vm.ManifestDel(vol.Meta.Node, vol.SelfLink().String())
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			return err
		}
	}

	return nil
}

func volumeManifestCheckEqual(mf *types.VolumeManifest, vol *types.Volume) bool {

	if mf.Capacity.Storage != vol.Spec.Capacity.Storage {
		return false
	}

	if mf.HostPath != vol.Spec.HostPath {
		return false
	}

	if mf.AccessMode != vol.Spec.AccessMode {
		return false
	}

	if mf.Type != vol.Spec.Type {
		return false
	}

	return true
}
