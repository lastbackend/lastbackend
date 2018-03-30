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

const serviceStorage = "services"

// Service Service type for interface in interfaces folder
type ServiceStorage struct {
	storage.Service
}

// Get service by name
func (s *ServiceStorage) Get(ctx context.Context, namespace, name string) (*types.Service, error) {

	log.V(logLevel).Debugf("storage:etcd:service:> get by name: %s", name)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:service:> get by name err: %s", err.Error())
		return nil, err
	}

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("storage:etcd:service:> get by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + serviceStorage + `\/.+\/(meta|status|spec)\b`
	var (
		service = new(types.Service)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> get by name err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyMeta := keyDirCreate(serviceStorage, s.keyCreate(namespace, name))
	if err := client.Map(ctx, keyMeta, filter, service); err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> get by name %s err: %s", name, err.Error())
		return nil, err
	}

	if service.Meta.Name == "" {
		return nil, errors.New(store.ErrEntityNotFound)
	}

	return service, nil
}

// Get services by namespace name
func (s *ServiceStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Service, error) {

	log.V(logLevel).Debugf("storage:etcd:service:> get list by namespace: %s", namespace)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:service:> get list by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + serviceStorage + `\/(.+)\/(meta|status|spec)\b`

	var (
		services = make(map[string]*types.Service)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> get list by namespace err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	key := keyCreate(serviceStorage, fmt.Sprintf("%s:", namespace))
	if err := client.MapList(ctx, key, filter, services); err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> err: %s", namespace, err.Error())
		return nil, err
	}

	return services, nil
}

// Update service status
func (s *ServiceStorage) SetStatus(ctx context.Context, service *types.Service) error {

	log.V(logLevel).Debugf("storage:etcd:service:> update service status: %#v", service)

	if err := s.checkServiceExists(ctx, service); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:>: update service err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(serviceStorage, s.keyGet(service), "status")
	if err := client.Upsert(ctx, key, service.Status, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:>: update service err: %s", err.Error())
		return err
	}

	return nil
}

// Update service spec in storage
func (s *ServiceStorage) SetSpec(ctx context.Context, service *types.Service) error {

	log.V(logLevel).Debugf("storage:etcd:service:> update service spec: %#v", service)

	if err := s.checkServiceExists(ctx, service); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:>: update service err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(serviceStorage, s.keyGet(service), "spec")
	if err := client.Upsert(ctx, key, service.Spec, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:>: update service err: %s", err.Error())
		return err
	}

	return nil
}

// Insert new service into storage
func (s *ServiceStorage) Insert(ctx context.Context, service *types.Service) error {

	log.V(logLevel).Debugf("storage:etcd:service:> insert service: %#v", service)

	if err := s.checkServiceArgument(service); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> insert service err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(serviceStorage, s.keyGet(service), "meta")
	if err := tx.Create(keyMeta, service.Meta, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> insert service err: %s", err.Error())
		return err
	}

	keyStatus := keyCreate(serviceStorage, s.keyGet(service), "status")
	if err := tx.Create(keyStatus, service.Status, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> insert service err: %s", err.Error())
		return err
	}

	keySpec := keyCreate(serviceStorage, s.keyGet(service), "spec")
	if err := tx.Create(keySpec, service.Spec, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> insert service err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> insert service err: %s", err.Error())
		return err
	}

	return nil
}

// Update service in storage
func (s *ServiceStorage) Update(ctx context.Context, service *types.Service) error {

	log.V(logLevel).Debugf("storage:etcd:service:> update service: %#v", service)

	if err := s.checkServiceExists(ctx, service); err != nil {
		return err
	}

	service.Meta.Updated = time.Now()
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> update service err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(serviceStorage, s.keyGet(service), "meta")
	if err := client.Upsert(ctx, keyMeta, service.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> update service err: %s", err.Error())
		return err
	}

	return nil
}

// Remove service from storage
func (s *ServiceStorage) Remove(ctx context.Context, service *types.Service) error {

	log.V(logLevel).Debugf("storage:etcd:service:> remove service: %#v", service)

	if err := s.checkServiceExists(ctx, service); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> remove err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(serviceStorage, s.keyGet(service))
	if err := client.DeleteDir(ctx, key); err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> remove service err: %s", err.Error())
		return err
	}

	return nil
}

// Watch service changes
func (s *ServiceStorage) Watch(ctx context.Context, service chan *types.Service) error {

	log.V(logLevel).Debug("storage:etcd:service:> watch service")

	const filter = `\b\/` + serviceStorage + `\/(.+):(.+)/.+\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> watch service err: %s", err.Error())
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

		if action == types.STORAGEDELEVENT {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2]); err == nil {
			service <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> watch service err: %s", err.Error())
		return err
	}

	return nil
}

// Watch service spec changes
func (s *ServiceStorage) WatchSpec(ctx context.Context, service chan *types.Service) error {

	log.V(logLevel).Debug("storage:etcd:service:> watch service")

	const filter = `\b\/` + serviceStorage + `\/(.+):(.+)/spec\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> watch service err: %s", err.Error())
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

		if action == types.STORAGEDELEVENT {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2]); err == nil {
			service <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> watch service err: %s", err.Error())
		return err
	}

	return nil
}

// Watch service status changes
func (s *ServiceStorage) WatchStatus(ctx context.Context, service chan *types.Service) error {

	log.V(logLevel).Debug("storage:etcd:service:> watch service")

	const filter = `\b\/` + serviceStorage + `\/(.+):(.+)/status\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> watch service err: %s", err.Error())
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

		if action == types.STORAGEDELEVENT {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2]); err == nil {
			service <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> watch service err: %s", err.Error())
		return err
	}

	return nil
}

// Clear service storage
func (s *ServiceStorage) Clear(ctx context.Context) error {

	log.V(logLevel).Debugf("storage:etcd:service:> clear")

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> clear err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := client.DeleteDir(ctx, serviceStorage); err != nil {
		log.V(logLevel).Errorf("storage:etcd:service:> clear err: %s", err.Error())
		return err
	}

	return nil
}

// keyCreate util function
func (s *ServiceStorage) keyCreate(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

// keyGet util function
func (s *ServiceStorage) keyGet(svc *types.Service) string {
	return svc.SelfLink()
}

// newServiceStorage returns new storage
func newServiceStorage() *ServiceStorage {
	s := new(ServiceStorage)
	return s
}

// checkServiceArgument - check if argument is valid for manipulations
func (s *ServiceStorage) checkServiceArgument(service *types.Service) error {

	if service == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if service.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

// checkServiceArgument - check if service exists in store
func (s *ServiceStorage) checkServiceExists(ctx context.Context, service *types.Service) error {

	if err := s.checkServiceArgument(service); err != nil {
		return err
	}

	log.V(logLevel).Debugf("storage:etcd:service:> check service exists")

	if _, err := s.Get(ctx, service.Meta.Namespace, service.Meta.Name); err != nil {
		log.V(logLevel).Debugf("storage:etcd:service:> check service exists err: %s", err.Error())
		return err
	}

	return nil

	return nil
}
