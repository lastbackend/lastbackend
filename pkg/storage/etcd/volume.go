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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"regexp"
	"time"
)

const volumeStorage = "volumes"

// Volume Service type for interface in interfaces folder
type VolumeStorage struct {
	storage.Volume
}

func (s *VolumeStorage) Get(ctx context.Context, namespace, name string) (*types.Volume, error) {

	log.V(logLevel).Debugf("storage:etcd:volume:> get by name: %s", name)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:volume:> get by name err: %s", err.Error())
		return nil, err
	}

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("storage:etcd:volume:> get by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + volumeStorage + `\/.+\/(?:meta|status|spec)\b`

	var (
		volume = new(types.Volume)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyVolume := keyCreate(volumeStorage, s.keyCreate(namespace, name))
	if err := client.Map(ctx, keyVolume, filter, volume); err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> err: %s", name, err.Error())
		return nil, err
	}

	if volume.Meta.Name == "" {
		return nil, errors.New(store.ErrEntityNotFound)
	}

	return volume, nil
}

// Get volume by namespace name
func (s *VolumeStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Volume, error) {

	log.V(logLevel).Debugf("storage:etcd:volume:> get list by namespace: %s", namespace)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:volume:> get list by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + volumeStorage + `\/.+\/(?:meta|status|spec)\b`

	var (
		volumes = make(map[string]*types.Volume)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyVolume := keyCreate(volumeStorage, fmt.Sprintf("%s:", namespace))
	if err := client.MapList(ctx, keyVolume, filter, volumes); err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> err: %s", namespace, err.Error())
		return nil, err
	}

	return volumes, nil
}

// Get volume by service name
func (s *VolumeStorage) ListByService(ctx context.Context, namespace, service string) (map[string]*types.Volume, error) {

	log.V(logLevel).Debugf("storage:etcd:volume:> get list by namespace and service: %s:%s", namespace, service)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:volume:> get list by namespace and service err: %s", err.Error())
		return nil, err
	}

	if len(service) == 0 {
		err := errors.New("service can not be empty")
		log.V(logLevel).Errorf("storage:etcd:volume:> get list by namespace and service err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + volumeStorage + `\/.+\/(?:meta|status|spec)\b`

	var (
		volumes = make(map[string]*types.Volume)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:>  get list by namespace and service err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyVolume := keyCreate(volumeStorage, fmt.Sprintf("%s:%s:", namespace, service))
	if err := client.MapList(ctx, keyVolume, filter, volumes); err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> err: %s", namespace, err.Error())
		return nil, err
	}

	return volumes, nil
}

// Update volume status
func (s *VolumeStorage) SetStatus(ctx context.Context, volume *types.Volume) error {

	log.V(logLevel).Debugf("storage:etcd:volume:> update volume status: %#v", volume)

	if err := s.checkVolumeExists(ctx, volume); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:>: update volume err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(volumeStorage, s.keyGet(volume), "status")
	if err := client.Upsert(ctx, key, volume.Status, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:>: update volume err: %s", err.Error())
		return err
	}

	return nil
}

// Update volume status
func (s *VolumeStorage) SetSpec(ctx context.Context, volume *types.Volume) error {

	log.V(logLevel).Debugf("storage:etcd:volume:> update volume spec: %#v", volume)

	if err := s.checkVolumeExists(ctx, volume); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:>: update volume err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(volumeStorage, s.keyGet(volume), "spec")
	if err := client.Upsert(ctx, key, volume.Spec, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:>: update volume err: %s", err.Error())
		return err
	}

	return nil
}

// Insert new volume into storage
func (s *VolumeStorage) Insert(ctx context.Context, volume *types.Volume) error {

	log.V(logLevel).Debugf("storage:etcd:volume:> insert volume: %#v", volume)

	if err := s.checkVolumeArgument(volume); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> insert volume err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(volumeStorage, s.keyGet(volume), "meta")
	if err := tx.Create(keyMeta, volume.Meta, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> insert volume err: %s", err.Error())
		return err
	}

	keyStatus := keyCreate(volumeStorage, s.keyGet(volume), "status")
	if err := tx.Create(keyStatus, volume.Status, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> insert volume err: %s", err.Error())
		return err
	}

	keySpec := keyCreate(volumeStorage, s.keyGet(volume), "spec")
	if err := tx.Create(keySpec, volume.Spec, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> insert volume err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> insert volume err: %s", err.Error())
		return err
	}

	return nil
}

// Update volume info
func (s *VolumeStorage) Update(ctx context.Context, volume *types.Volume) error {

	if err := s.checkVolumeExists(ctx, volume); err != nil {
		return err
	}

	volume.Meta.Updated = time.Now()
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> update volume err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(volumeStorage, s.keyGet(volume), "meta")
	if err := client.Upsert(ctx, keyMeta, volume.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> update volume err: %s", err.Error())
		return err
	}

	return nil
}

// Remove volume by id from storage
func (s *VolumeStorage) Remove(ctx context.Context, volume *types.Volume) error {

	if err := s.checkVolumeExists(ctx, volume); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> remove err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(volumeStorage, s.keyGet(volume))
	if err := client.DeleteDir(ctx, keyMeta); err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> remove volume err: %s", err.Error())
		return err
	}

	return nil
}

// Watch volume changes
func (s *VolumeStorage) Watch(ctx context.Context, volume chan *types.Volume) error {

	log.V(logLevel).Debug("storage:etcd:volume:> watch volume")

	const filter = `\b\/` + volumeStorage + `\/(.+):(.+):(.+)/.+\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> watch volume err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(volumeStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if action == types.STORAGEDELEVENT {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2]); err == nil {
			volume <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> watch volume err: %s", err.Error())
		return err
	}

	return nil
}

// Watch volume spec changes
func (s *VolumeStorage) WatchSpec(ctx context.Context, volume chan *types.Volume) error {

	log.V(logLevel).Debug("storage:etcd:volume:> watch volume by spec")

	const filter = `\b\/` + volumeStorage + `\/(.+):(.+):(.+)/spec\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> watch volume by spec err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(volumeStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if action == types.STORAGEDELEVENT {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2]); err == nil {
			volume <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> watch volume by spec err: %s", err.Error())
		return err
	}

	return nil
}

// Watch volume status changes
func (s *VolumeStorage) WatchStatus(ctx context.Context, volume chan *types.Volume) error {

	log.V(logLevel).Debug("storage:etcd:volume:> watch volume by spec")

	const filter = `\b\/` + volumeStorage + `\/(.+):(.+):(.+)/status\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> watch volume by spec err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(volumeStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if action == types.STORAGEDELEVENT {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2]); err == nil {
			volume <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> watch volume by spec err: %s", err.Error())
		return err
	}

	return nil
}

// Clear volume storage
func (s *VolumeStorage) Clear(ctx context.Context) error {

	log.V(logLevel).Debugf("storage:etcd:volume:> clear")

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> clear err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := client.DeleteDir(ctx, volumeStorage); err != nil {
		log.V(logLevel).Errorf("storage:etcd:volume:> clear err: %s", err.Error())
		return err
	}

	return nil
}

// keyCreate util function
func (s *VolumeStorage) keyCreate(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

// keyGet util function
func (s *VolumeStorage) keyGet(t *types.Volume) string {
	return t.SelfLink()
}

func newVolumeStorage() *VolumeStorage {
	s := new(VolumeStorage)
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
func (s *VolumeStorage) checkVolumeExists(ctx context.Context, volume *types.Volume) error {

	if err := s.checkVolumeArgument(volume); err != nil {
		return err
	}

	log.V(logLevel).Debugf("storage:etcd:volume:> check volume exists")

	if _, err := s.Get(ctx, volume.Meta.Namespace, volume.Meta.Name); err != nil {
		log.V(logLevel).Debugf("storage:etcd:volume:> check volume exists err: %s", err.Error())
		return err
	}

	return nil
}
