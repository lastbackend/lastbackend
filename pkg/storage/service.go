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
	"github.com/lastbackend/lastbackend/pkg/log"
)

const serviceStorage string = "services"

// Service Service type for interface in interfaces folder
type ServiceStorage struct {
	IService
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get service by name
func (s *ServiceStorage) GetByName(ctx context.Context, app, name string) (*types.Service, error) {

	log.V(logLevel).Debugf("Storage: Service: get by name: %s in app: %s", name, app)

	if len(app) == 0 {
		err := errors.New("app can not be empty")
		log.V(logLevel).Errorf("Storage: Service: get service err: %s", err.Error())
		return nil, err
	}

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("Storage: Service: get service err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + serviceStorage + `\/.+\/(?:meta|state)\b`

	var (
		filterServiceEndpoint = `\b.+` + endpointStorage + `\/` + name + `-` + app + `\..+\b`
		endpoints             = make(map[string][]string)
		service               = new(types.Service)
	)

	service.Spec = make(map[string]*types.ServiceSpec)
	service.Pods = make(map[string]*types.Pod)

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Service: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyService := keyCreate(serviceStorage, app, name)
	if err := client.Map(ctx, keyService, filter, service); err != nil {
		log.V(logLevel).Errorf("Storage: Service: map service `%s` err: %s", name, err.Error())
		return nil, err
	}

	if service.Meta.Name == "" {
		return nil, errors.New(store.ErrKeyNotFound)
	}

	keySpec := keyCreate(serviceStorage, app, name, "specs")
	if err := client.Map(ctx, keySpec, "", &service.Spec); err != nil && err.Error() != store.ErrKeyNotFound {
		log.V(logLevel).Errorf("Storage: Service: Map service `%s` spec err: %s", name, err.Error())
		return nil, err
	}

	keyPods := keyCreate(podStorage, app, fmt.Sprintf("%s:%s", app, name))
	if err := client.Map(ctx, keyPods, "", &service.Pods); err != nil && err.Error() != store.ErrKeyNotFound {
		log.V(logLevel).Errorf("Storage: Service: Map service `%s` pods err: %s", name, err.Error())
		return nil, err
	}

	for _, pod := range service.Pods {
		name := strings.Replace(pod.Meta.Name, ":", "-", -1)
		filterPodEndpoint := `\b.+` + endpointStorage + `\/` + name + `\..+\b`
		endpoints := make(map[string][]string)
		keyEndpoints := keyCreate(endpointStorage)
		if err := client.Map(ctx, keyEndpoints, filterPodEndpoint, endpoints); err != nil && err.Error() != store.ErrKeyNotFound {
			log.V(logLevel).Errorf("Storage: Service: map endpoints for pod err: %s", name, err.Error())
			return nil, err
		}

		for pod.Meta.Endpoint = range endpoints {
			break
		}
	}

	keyEndpoints := keyCreate(endpointStorage)
	if err := client.Map(ctx, keyEndpoints, filterServiceEndpoint, endpoints); err != nil && err.Error() != store.ErrKeyNotFound {
		log.V(logLevel).Errorf("Storage: Service: map service endpoint `%s` meta err: %s", name, err.Error())
		return nil, err
	}

	for service.DNS.Primary = range endpoints {
		break
	}

	return service, nil
}

// Get service by pod name
func (s *ServiceStorage) GetByPodName(ctx context.Context, name string) (*types.Service, error) {

	log.V(logLevel).Debugf("Storage: Service: get by pod name: %s in app: %s", name)

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("Storage: Service: get service by pod name err: %s", err.Error())
		return nil, err
	}

	parts := strings.Split(name, ":")
	if len(parts) < 3 {
		err := errors.New(fmt.Sprintf("can not parse pod name: %s", name))
		log.V(logLevel).Errorf("Storage: Service: get service by pod name err: %s", err.Error())
		return nil, err
	}

	return s.GetByName(ctx, parts[0], parts[1])
}

// List services
func (s *ServiceStorage) ListByApp(ctx context.Context, app string) ([]*types.Service, error) {

	log.V(logLevel).Debugf("Storage: Service: get service list in app: %s", app)

	if len(app) == 0 {
		err := errors.New("app can not be empty")
		log.V(logLevel).Errorf("Storage: Service: get service list err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + serviceStorage + `\/.+\/(?:meta|state)\b`

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Service: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	var services = []*types.Service{}

	keyServices := keyCreate(serviceStorage, app)
	if err := client.List(ctx, keyServices, filter, &services); err != nil {
		log.V(logLevel).Errorf("Storage: Service: list services err: %s", err.Error())
		return nil, err
	}

	return services, nil
}

// Count services
func (s *ServiceStorage) CountByApp(ctx context.Context, app string) (int, error) {

	log.V(logLevel).Debugf("Storage: Service: count service list in app: %s", app)

	if len(app) == 0 {
		err := errors.New("app can not be empty")
		log.V(logLevel).Errorf("Storage: Service: count service list err: %s", err.Error())
		return int(0), err
	}

	const filter = `\b.+` + serviceStorage + `\/.+\/meta\b`

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Service: create client err: %s", err.Error())
		return int(0), err
	}
	defer destroy()

	keyServices := keyCreate(serviceStorage, app)
	count, err := client.Count(ctx, keyServices, filter)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Service: list services err: %s", err.Error())
		return int(0), err
	}

	return count, nil
}

// Insert new service into storage
func (s *ServiceStorage) Insert(ctx context.Context, service *types.Service) error {

	log.V(logLevel).Debugf("Storage: Service: insert service: %#v", service)

	if service == nil {
		err := errors.New("service can not be nil")
		log.V(logLevel).Errorf("Storage: Service: insert service err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Service: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := s.updateState(ctx, service); err != nil {
		log.V(logLevel).Errorf("Storage: Service: update state err: %s", err.Error())
		return err
	}

	tx := client.Begin(ctx)

	keyMeta := keyCreate(serviceStorage, service.Meta.App, service.Meta.Name, "meta")
	if err := tx.Create(keyMeta, &service.Meta, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Service: create meta err: %s", err.Error())
		return err
	}

	for _, spec := range service.Spec {
		keyConfig := keyCreate(serviceStorage, service.Meta.App, service.Meta.Name, "specs", spec.Meta.ID)
		if err := tx.Create(keyConfig, &spec, 0); err != nil {
			log.V(logLevel).Errorf("Storage: Service: create spec err: %s", err.Error())
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("Storage: Service: commit transaction err: %s", err.Error())
		return err
	}

	return nil
}

// Update service in storage
func (s *ServiceStorage) Update(ctx context.Context, service *types.Service) error {

	log.V(logLevel).Debugf("Storage: Service: update service: %#v", service)

	if service == nil {
		err := errors.New("service can not be nil")
		log.V(logLevel).Errorf("Storage: Service: update service err: %s", err.Error())
		return err
	}

	service.Meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Service: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := s.updateState(ctx, service); err != nil {
		log.V(logLevel).Errorf("Storage: Service: update state err: %s", err.Error())
		return err
	}

	keyMeta := keyCreate(serviceStorage, service.Meta.App, service.Meta.Name, "meta")
	if err := client.Update(ctx, keyMeta, service.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Service: update meta err: %s", err.Error())
		return err
	}

	return nil
}

// Update service spec in storage
func (s *ServiceStorage) UpdateSpec(ctx context.Context, service *types.Service) error {

	log.V(logLevel).Debugf("Storage: Service: update spec service: %#v", service)

	if service == nil {
		err := errors.New("service can not be nil")
		log.V(logLevel).Errorf("Storage: Service: update spec service err: %s", err.Error())
		return err
	}

	service.Meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Service: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := s.updateState(ctx, service); err != nil {
		log.V(logLevel).Errorf("Storage: Service: update state err: %s", err.Error())
		return err
	}

	tx := client.Begin(ctx)

	specs := make(map[string]*types.ServiceSpec)
	keySpecs := keyCreate(serviceStorage, service.Meta.App, service.Meta.Name, "specs")
	if err := client.Map(ctx, keySpecs, "", specs); err != nil {
		if err.Error() != store.ErrKeyNotFound {
			log.V(logLevel).Errorf("Storage: Service: map specs err: %s", err.Error())
			return err
		}
	}

	for id := range specs {
		if _, ok := service.Spec[id]; ok {
			continue
		}
		keySpec := keyCreate(serviceStorage, service.Meta.App, service.Meta.Name, "specs", id)
		tx.DeleteDir(keySpec)
	}

	for id, spec := range service.Spec {
		keySpec := keyCreate(serviceStorage, service.Meta.App, service.Meta.Name, "specs", id)
		if err := tx.Upsert(keySpec, &spec, 0); err != nil {
			log.V(logLevel).Errorf("Storage: Service: upsert specs err: %s", err.Error())
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("Storage: Service: commit transaction err: %s", err.Error())
		return err
	}

	return nil
}

// Remove service model
func (s *ServiceStorage) Remove(ctx context.Context, service *types.Service) error {

	log.V(logLevel).Debugf("Storage: Service: remove service: %#v", service)

	if service == nil {
		err := errors.New("service can not be nil")
		log.V(logLevel).Errorf("Storage: Service: remove service err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Service: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyService := keyCreate(serviceStorage, service.Meta.App, service.Meta.Name)
	tx.DeleteDir(keyService)

	keyServiceController := keyCreate(systemStorage, types.KindController, "services", fmt.Sprintf("%s:%s", service.Meta.App, service.Meta.Name))
	tx.Delete(keyServiceController)

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("Storage: Service: commit transaction err: %s", err.Error())
		return err
	}

	return nil
}

// Remove services from app
func (s *ServiceStorage) RemoveByApp(ctx context.Context, app string) error {

	log.V(logLevel).Debugf("Storage: Service: remove services in app: %s", app)

	if len(app) == 0 {
		err := errors.New("app can not be nil")
		log.V(logLevel).Errorf("Storage: Service: remove services in app err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Service: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyAll := keyCreate(serviceStorage, app)
	if err := client.DeleteDir(ctx, keyAll); err != nil {
		log.V(logLevel).Errorf("Storage: Service: delete dir err: %s", err.Error())
		return err
	}

	return nil
}

func (s *ServiceStorage) Watch(ctx context.Context, service chan *types.Service) error {

	log.V(logLevel).Debug("Storage: Service: watch service")

	const filter = `\/` + systemStorage + `\/` + types.KindController + `\/services\/([a-z0-9_-]+):([a-z0-9_-]+)\b`
	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Service: create client err: %s", err.Error())
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

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("Storage: Service: watch service err: %s", err.Error())
		return err
	}

	return nil
}

func (s *ServiceStorage) SpecWatch(ctx context.Context, service chan *types.Service) error {

	log.V(logLevel).Debug("Storage: Service: watch service by spec")

	const filter = `\b\/` + serviceStorage + `\/(.+)\/(.+)\/specs/.+\b`
	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Service: create client err: %s", err.Error())
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

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("Storage: Service: watch service spec err: %s", err.Error())
		return err
	}

	return nil
}

func (s *ServiceStorage) PodsWatch(ctx context.Context, service chan *types.Service) error {

	log.V(logLevel).Debug("Storage: Service: watch service by pod")

	const filter = `\b\/` + podStorage + `\/(.+)/(.+)\b`
	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Service: create client err: %s", err.Error())
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

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("Storage: Service: watch service pod err: %s", err.Error())
		return err
	}

	return nil
}

func (s *ServiceStorage) BuildsWatch(ctx context.Context, service chan *types.Service) error {

	log.V(logLevel).Debug("Storage: Service: watch service by build")

	const filter = `\b.+` + buildStorage + `\/(.+)\/(.+)\/.+\b`
	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Service: create client err: %s", err.Error())
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

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("Storage: Service: watch service build err: %s", err.Error())
		return err
	}

	return nil
}

// Update service state
func (s *ServiceStorage) updateState(ctx context.Context, service *types.Service) error {

	log.V(logLevel).Debugf("Storage: Service: update service state: %#v", service)

	if service == nil {
		err := errors.New("service can not be nil")
		log.V(logLevel).Errorf("Storage: Service: update service state err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		log.V(logLevel).Errorf("Storage: Service: create client err: %s", err.Error())
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

	keyState := keyCreate(serviceStorage, service.Meta.App, service.Meta.Name, "state")
	if err := client.Upsert(ctx, keyState, service.State, nil, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Service: upsert state err: %s", err.Error())
		return err
	}

	keyServiceController := keyCreate(systemStorage, types.KindController, "services", fmt.Sprintf("%s:%s", service.Meta.App, service.Meta.Name))
	if err := client.Upsert(ctx, keyServiceController, &service.State.State, nil, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Service: upsert services controller err: %s", err.Error())
		return err
	}

	return nil
}

func newServiceStorage(config store.Config) *ServiceStorage {
	s := new(ServiceStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
