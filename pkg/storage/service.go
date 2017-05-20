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

package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"regexp"
	"strings"
	"time"
)

const serviceStorage string = "services"

// Service Service type for interface in interfaces folder
type ServiceStorage struct {
	IService
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get service by name
func (s *ServiceStorage) GetByName(ctx context.Context, namespace, name string) (*types.Service, error) {

	const filter = `\b.+` + serviceStorage + `\/.+\/(?:meta|state)\b`
	var service = new(types.Service)
	service.Spec = make(map[string]*types.ServiceSpec)
	service.Pods = make(map[string]*types.Pod)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	keyService := keyCreate(serviceStorage, namespace, name)
	if err := client.Map(ctx, keyService, filter, service); err != nil {
		return nil, err
	}

	keySpec := keyCreate(serviceStorage, namespace, name, "specs")
	if err := client.Map(ctx, keySpec, "", &service.Spec); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return service, nil
		}
		return nil, err
	}

	keyPods := keyCreate(podStorage, namespace, fmt.Sprintf("%s:%s", namespace, name))
	if err := client.Map(ctx, keyPods, "", &service.Pods); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return service, nil
		}
		return nil, err
	}

	return service, nil
}

// Get service by pod name
func (s *ServiceStorage) GetByPodName(ctx context.Context, name string) (*types.Service, error) {
	parts := strings.Split(name, ":")
	if len(parts) < 3 {
		return nil, errors.New(fmt.Sprintf("can not parse pod name: %s", name))
	}
	return s.GetByName(ctx, parts[0], parts[1])
}

// List services
func (s *ServiceStorage) ListByNamespace(ctx context.Context, namespace string) ([]*types.Service, error) {

	const filter = `\b.+` + serviceStorage + `\/.+\/(?:meta|state)\b`

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	var services = []*types.Service{}

	keyServices := keyCreate(serviceStorage, namespace)
	if err := client.List(ctx, keyServices, filter, &services); err != nil {
		return nil, err
	}

	return services, nil
}

// Insert new service into storage
func (s *ServiceStorage) Insert(ctx context.Context, service *types.Service) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	if err := s.updateState(ctx, service); err != nil {
		return err
	}

	tx := client.Begin(ctx)

	keyMeta := keyCreate(serviceStorage, service.Meta.Namespace, service.Meta.Name, "meta")
	if err := tx.Create(keyMeta, &service.Meta, 0); err != nil {
		return err
	}

	for _, spec := range service.Spec {
		keyConfig := keyCreate(serviceStorage, service.Meta.Namespace, service.Meta.Name, "specs", spec.Meta.ID)
		if err := tx.Create(keyConfig, &spec, 0); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Update service in storage
func (s *ServiceStorage) Update(ctx context.Context, service *types.Service) error {

	service.Meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	if err := s.updateState(ctx, service); err != nil {
		return err
	}

	tx := client.Begin(ctx)

	keyMeta := keyCreate(serviceStorage, service.Meta.Namespace, service.Meta.Name, "meta")
	if err := tx.Update(keyMeta, service.Meta, 0); err != nil {
		return err
	}

	keyState := keyCreate(serviceStorage, service.Meta.Namespace, service.Meta.Name, "state")
	if err := client.Upsert(ctx, keyState, service.State, nil, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Update service spec in storage
func (s *ServiceStorage) UpdateSpec(ctx context.Context, service *types.Service) error {

	service.Meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	if err := s.updateState(ctx, service); err != nil {
		return err
	}

	tx := client.Begin(ctx)

	specs := make(map[string]*types.ServiceSpec)
	keySpecs := keyCreate(serviceStorage, service.Meta.Namespace, service.Meta.Name, "specs")
	if err := client.Map(ctx, keySpecs, "", specs); err != nil {
		if err.Error() != store.ErrKeyNotFound {
			return err
		}
	}

	for id := range specs {
		if _, ok := service.Spec[id]; ok {
			continue
		}
		keySpec := keyCreate(serviceStorage, service.Meta.Namespace, service.Meta.Name, "specs", id)
		tx.DeleteDir(keySpec)
	}

	for id, spec := range service.Spec {
		keySpec := keyCreate(serviceStorage, service.Meta.Namespace, service.Meta.Name, "specs", id)
		if err := tx.Upsert(keySpec, &spec, 0); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Remove service model
func (s *ServiceStorage) Remove(ctx context.Context, service *types.Service) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyService := keyCreate(serviceStorage, service.Meta.Namespace, service.Meta.Name)
	tx.DeleteDir(keyService)

	keyServiceController := s.util.Key(ctx, systemStorage, types.KindController, "services", fmt.Sprintf("%s:%s", service.Meta.Namespace, service.Meta.Name))
	tx.Delete(keyServiceController)

	return tx.Commit()
}

// Remove services from namespace
func (s *ServiceStorage) RemoveByNamespace(ctx context.Context, namespace string) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	keyAll := keyCreate(serviceStorage, namespace)
	if err := client.DeleteDir(ctx, keyAll); err != nil {
		return err
	}

	return nil
}

func (s *ServiceStorage) Watch(ctx context.Context, service chan *types.Service) error {
	const filter = `\/` + systemStorage + `\/` + types.KindController + `\/services\/([a-z0-9_-]+):([a-z0-9_-]+)\b`
	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(systemStorage, types.KindController, "services")
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if svc, err := s.GetByName(ctx, keys[1], keys[2]); err == nil {
			service <- svc
		}
	}

	client.Watch(ctx, key, filter, cb)
	return nil
}

func (s *ServiceStorage) SpecWatch(ctx context.Context, service chan *types.Service) error {
	const filter = `\b\/` + serviceStorage + `\/(.+)\/(.+)\/specs/.+\b`
	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(serviceStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if svc, err := s.GetByName(ctx, keys[1], keys[2]); err == nil {
			s.updateState(ctx, svc)
			service <- svc
		}

	}

	client.Watch(ctx, key, filter, cb)
	return nil
}

func (s *ServiceStorage) PodsWatch(ctx context.Context, service chan *types.Service) error {
	const filter = `\b\/` + podStorage + `\/(.+)/(.+)\b`
	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(podStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if svc, err := s.GetByPodName(ctx, keys[2]); err == nil {
			s.updateState(ctx, svc)
			service <- svc
		}
	}

	client.Watch(ctx, key, filter, cb)
	return nil
}

func (s *ServiceStorage) BuildsWatch(ctx context.Context, service chan *types.Service) error {
	const filter = `\b.+` + buildStorage + `\/(.+)\/(.+)\/.+\b`
	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(serviceStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if svc, err := s.GetByName(ctx, keys[1], keys[2]); err == nil {
			s.updateState(ctx, svc)
			service <- svc
		}

	}

	client.Watch(ctx, key, filter, cb)
	return nil
}

// Update service state
func (s *ServiceStorage) updateState(ctx context.Context, service *types.Service) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	// Update service resource info
	service.State.Resources = types.ServiceResourcesState{}

	for _, s := range service.Spec {
		service.State.Resources.Memory += int(s.Memory) * service.Meta.Replicas
	}

	// Update service state info
	service.State.Replicas = types.ServiceReplicasState{}

	for _, p := range service.Pods {
		service.State.Replicas.Total++
		switch p.State.State {
		case types.StateCreated:
			service.State.Replicas.Created++
		case types.StateStarted:
			service.State.Replicas.Running++
		case types.StateStopped:
			service.State.Replicas.Stopped++
		case types.StateError:
			service.State.Replicas.Errored++
		}

		if p.State.Provision {
			service.State.Replicas.Provision++
		}

		if p.State.Ready {
			service.State.Replicas.Ready++
		}
	}

	keyState := keyCreate(serviceStorage, service.Meta.Namespace, service.Meta.Name, "state")
	if err := client.Upsert(ctx, keyState, service.State, nil, 0); err != nil {
		return err
	}

	keyServiceController := s.util.Key(ctx, systemStorage, types.KindController, "services", fmt.Sprintf("%s:%s", service.Meta.Namespace, service.Meta.Name))
	if err := client.Upsert(ctx, keyServiceController, &service.State.State, nil, 0); err != nil {
		return err
	}

	return nil
}

func newServiceStorage(config store.Config, util IUtil) *ServiceStorage {
	s := new(ServiceStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
