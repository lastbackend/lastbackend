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

const endpointStorage = "endpoint"

// Service Endpoint type for interface in interfaces folder
type EndpointStorage struct {
	storage.Endpoint
}

// Get endpoints by id
func (s *EndpointStorage) Get(ctx context.Context, namespace, service string) (*types.Endpoint, error) {

	log.V(logLevel).Debugf("storage:etcd:endpoint:> get by service: %s", service)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:endpoint:> get by name err: %s", err.Error())
		return nil, err
	}

	if len(service) == 0 {
		err := errors.New("service can not be empty")
		log.V(logLevel).Errorf("storage:etcd:endpoint:> get by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + endpointStorage + `\/.+\/(meta|status|spec)\b`

	var (
		endpoint = new(types.Endpoint)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyEndpoint := keyDirCreate(endpointStorage, s.keyCreate(namespace, service))
	if err := client.Map(ctx, keyEndpoint, filter, endpoint); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> err: %s", service, err.Error())
		return nil, err
	}

	if endpoint.Meta.Name == "" {
		return nil, errors.New(store.ErrEntityNotFound)
	}

	return endpoint, nil
}

// Get endpoints
func (s *EndpointStorage) List(ctx context.Context) (map[string]*types.Endpoint, error) {

	log.V(logLevel).Debug("storage:etcd:endpoint:> get list")

	const filter = `\b.+` + endpointStorage + `\/(.+)\/(meta|status|spec)\b`

	var (
		endpoints = make(map[string]*types.Endpoint)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyEndpoint := keyDirCreate(endpointStorage)
	if err := client.MapList(ctx, keyEndpoint, filter, endpoints); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> err: %s", err.Error())
		return nil, err
	}

	return endpoints, nil
}

// Get endpoint by namespace name
func (s *EndpointStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Endpoint, error) {

	log.V(logLevel).Debugf("storage:etcd:endpoint:> get list by namespace: %s", namespace)

	if len(namespace) == 0 {
		err := errors.New("namespace can not be empty")
		log.V(logLevel).Errorf("storage:etcd:endpoint:> get list by name err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + endpointStorage + `\/(.+)\/(meta|status|spec)\b`

	var (
		endpoints = make(map[string]*types.Endpoint)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	keyEndpoint := keyCreate(endpointStorage, fmt.Sprintf("%s:", namespace))
	if err := client.MapList(ctx, keyEndpoint, filter, endpoints); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> err: %s", namespace, err.Error())
		return nil, err
	}

	return endpoints, nil
}

// Update endpoint status
func (s *EndpointStorage) SetStatus(ctx context.Context, endpoint *types.Endpoint) error {

	log.V(logLevel).Debugf("storage:etcd:endpoint:> update endpoint status: %#v", endpoint)

	if err := s.checkEndpointExists(ctx, endpoint); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:>: update endpoint err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(endpointStorage, s.keyGet(endpoint), "status")
	if err := client.Upsert(ctx, key, endpoint.Status, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:>: update endpoint err: %s", err.Error())
		return err
	}

	return nil
}

// Update endpoint status
func (s *EndpointStorage) SetSpec(ctx context.Context, endpoint *types.Endpoint) error {

	log.V(logLevel).Debugf("storage:etcd:endpoint:> update endpoint spec: %#v", endpoint)

	if err := s.checkEndpointExists(ctx, endpoint); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:>: update endpoint err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(endpointStorage, s.keyGet(endpoint), "spec")
	if err := client.Upsert(ctx, key, endpoint.Spec, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:>: update endpoint err: %s", err.Error())
		return err
	}

	return nil
}

// Insert new endpoint into storage
func (s *EndpointStorage) Insert(ctx context.Context, endpoint *types.Endpoint) error {

	log.V(logLevel).Debugf("storage:etcd:endpoint:> insert endpoint: %#v", endpoint)

	if err := s.checkEndpointArgument(endpoint); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> insert endpoint err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(endpointStorage, s.keyGet(endpoint), "meta")
	if err := tx.Create(keyMeta, endpoint.Meta, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> insert endpoint err: %s", err.Error())
		return err
	}

	keyStatus := keyCreate(endpointStorage, s.keyGet(endpoint), "status")
	if err := tx.Create(keyStatus, endpoint.Status, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> insert endpoint err: %s", err.Error())
		return err
	}

	keySpec := keyCreate(endpointStorage, s.keyGet(endpoint), "spec")
	if err := tx.Create(keySpec, endpoint.Spec, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> insert endpoint err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> insert endpoint err: %s", err.Error())
		return err
	}

	return nil
}

// Update endpoint info
func (s *EndpointStorage) Update(ctx context.Context, endpoint *types.Endpoint) error {

	if err := s.checkEndpointExists(ctx, endpoint); err != nil {
		return err
	}

	endpoint.Meta.Updated = time.Now()
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> update endpoint err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(endpointStorage, s.keyGet(endpoint), "meta")
	if err := tx.Update(keyMeta, endpoint.Meta, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> update endpoint err: %s", err.Error())
		return err
	}

	keyStatus := keyCreate(endpointStorage, s.keyGet(endpoint), "status")
	if err := tx.Update(keyStatus, endpoint.Status, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> update endpoint err: %s", err.Error())
		return err
	}

	keySpec := keyCreate(endpointStorage, s.keyGet(endpoint), "spec")
	if err := tx.Update(keySpec, endpoint.Spec, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> update endpoint err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> update endpoint err: %s", err.Error())
		return err
	}

	return nil
}

// Remove endpoint by id from storage
func (s *EndpointStorage) Remove(ctx context.Context, endpoint *types.Endpoint) error {

	if err := s.checkEndpointExists(ctx, endpoint); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> remove err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(endpointStorage, s.keyGet(endpoint))
	if err := client.DeleteDir(ctx, keyMeta); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> remove endpoint err: %s", err.Error())
		return err
	}

	return nil
}

// Watch endpoint changes
func (s *EndpointStorage) Watch(ctx context.Context, endpoint chan *types.Endpoint) error {

	log.V(logLevel).Debug("storage:etcd:endpoint:> watch endpoint")

	const filter = `\b\/` + endpointStorage + `\/(.+):(.+)/.+\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> watch endpoint err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(endpointStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if action == store.STORAGEDELETEEVENT {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2]); err == nil {
			endpoint <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> watch endpoint err: %s", err.Error())
		return err
	}

	return nil
}

// Watch endpoint spec changes
func (s *EndpointStorage) WatchSpec(ctx context.Context, endpoint chan *types.Endpoint) error {

	log.V(logLevel).Debug("storage:etcd:endpoint:> watch endpoint by spec")

	const filter = `\b\/` + endpointStorage + `\/(.+):(.+)/spec\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> watch endpoint by spec err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(endpointStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if action == store.STORAGEDELETEEVENT {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2]); err == nil {
			endpoint <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> watch endpoint by spec err: %s", err.Error())
		return err
	}

	return nil
}

// Watch endpoint status changes
func (s *EndpointStorage) WatchStatus(ctx context.Context, endpoint chan *types.Endpoint) error {

	log.V(logLevel).Debug("storage:etcd:endpoint:> watch endpoint by spec")

	const filter = `\b\/` + endpointStorage + `\/(.+):(.+)/status\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> watch endpoint by spec err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(endpointStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 3 {
			return
		}

		if action == store.STORAGEDELETEEVENT {
			return
		}

		if d, err := s.Get(ctx, keys[1], keys[2]); err == nil {
			endpoint <- d
		}
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> watch endpoint by spec err: %s", err.Error())
		return err
	}

	return nil
}

func (s *EndpointStorage) EventSpec(ctx context.Context, event chan *types.EndpointSpecEvent) error {

	log.V(logLevel).Debug("storage:etcd:endpoint:> watch spec")

	const filter = `\b.+` + endpointStorage + `\/(.+)\/spec\b`

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> create client err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(endpointStorage)
	cb := func(action, key string, val []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 2 {
			return
		}

		e := new(types.EndpointSpecEvent)
		e.Event = action
		e.Name = keys[1]

		if err := client.Decode(ctx, val, &e.Spec); err != nil {
			log.Warnf("storage:etcd:endpoint:> decode err: %s", err.Error())
		}

		event <- e

		return
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> watch node err: %s", err.Error())
		return err
	}

	return nil
}

// Clear endpoint storage
func (s *EndpointStorage) Clear(ctx context.Context) error {

	log.V(logLevel).Debugf("storage:etcd:endpoint:> clear")

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> clear err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := client.DeleteDir(ctx, endpointStorage); err != nil {
		log.V(logLevel).Errorf("storage:etcd:endpoint:> clear err: %s", err.Error())
		return err
	}

	return nil
}

// keyCreate util function
func (s *EndpointStorage) keyCreate(namespace, service string) string {
	return fmt.Sprintf("%s:%s", namespace, service)
}

// keyGet util function
func (s *EndpointStorage) keyGet(t *types.Endpoint) string {
	return t.SelfLink()
}

func newEndpointStorage() *EndpointStorage {
	s := new(EndpointStorage)
	return s
}

// checkEndpointArgument - check if argument is valid for manipulations
func (s *EndpointStorage) checkEndpointArgument(endpoint *types.Endpoint) error {

	if endpoint == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if endpoint.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

// checkEndpointArgument - check if endpoint exists in store
func (s *EndpointStorage) checkEndpointExists(ctx context.Context, endpoint *types.Endpoint) error {

	if err := s.checkEndpointArgument(endpoint); err != nil {
		return err
	}

	log.V(logLevel).Debugf("storage:etcd:endpoint:> check endpoint exists")

	if _, err := s.Get(ctx, endpoint.Meta.Namespace, endpoint.Meta.Name); err != nil {
		log.V(logLevel).Debugf("storage:etcd:endpoint:> check endpoint exists err: %s", err.Error())
		return err
	}

	return nil
}
