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

	"encoding/json"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

const (
	logEndpointPrefix = "distribution:endpoint"
)

type Endpoint struct {
	context context.Context
	storage storage.Storage
}

func (e *Endpoint) Get(namespace, service string) (*types.Endpoint, error) {

	log.V(logLevel).Debugf("%s:get:> get endpoint by namespace %s and service %s", logEndpointPrefix, namespace, service)

	item := new(types.Endpoint)

	err := e.storage.Get(e.context, storage.EndpointKind, e.storage.Key().Endpoint(namespace, service), &item)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> get endpoint err: %v", logEndpointPrefix, err)
		return nil, err
	}

	return item, nil
}

func (e *Endpoint) ListByNamespace(namespace string) (map[string]*types.Endpoint, error) {
	log.Debugf("%s:listbynamespace:> in namespace: %s", namespace)

	list := make(map[string]*types.Endpoint, 0)

	err := e.storage.Map(e.context, storage.EndpointKind, e.storage.Filter().Endpoint().ByNamespace(namespace), &list)
	if err != nil {
		log.Errorf("%s:listbynamespace:> in namespace: %s err: %v", logEndpointPrefix, namespace, err)
		return nil, err
	}

	return list, nil
}

func (e *Endpoint) Create(namespace, service string, opts *types.EndpointCreateOptions) (*types.Endpoint, error) {
	endpoint := new(types.Endpoint)

	endpoint.Meta.Name = service
	endpoint.Meta.Namespace = namespace
	endpoint.Meta.SetDefault()
	endpoint.SelfLink()

	endpoint.Status.State = types.StateCreated
	endpoint.Spec.PortMap = make(map[uint16]string, 0)
	endpoint.Spec.Upstreams = make([]string, 0)

	for k, v := range opts.Ports {
		endpoint.Spec.PortMap[k] = v
	}

	endpoint.Spec.Policy = opts.Policy
	endpoint.Spec.Strategy.Route = opts.RouteStrategy
	endpoint.Spec.Strategy.Bind = opts.BindStrategy

	endpoint.Spec.IP = opts.IP
	endpoint.Spec.Domain = opts.Domain

	key := e.storage.Key().Endpoint(namespace, service)
	if err := e.storage.Put(e.context, storage.EndpointKind, key, endpoint, nil); err != nil {
		log.Errorf("%s:create:> distribution create endpoint: %s err: %v", logEndpointPrefix, endpoint.SelfLink(), err)
		return nil, err
	}

	return endpoint, nil
}

func (e *Endpoint) Update(endpoint *types.Endpoint, opts *types.EndpointUpdateOptions) (*types.Endpoint, error) {
	log.Debugf("%s:update:> endpoint: %s", logEndpointPrefix, endpoint.SelfLink())

	if len(opts.Ports) != 0 {
		endpoint.Spec.PortMap = make(map[uint16]string, 0)
		for k, v := range opts.Ports {
			endpoint.Spec.PortMap[k] = v
		}
	}

	if opts.IP != nil {
		endpoint.Spec.IP = *opts.IP
	}

	endpoint.Spec.Policy = opts.Policy
	endpoint.Spec.Strategy.Route = opts.RouteStrategy
	endpoint.Spec.Strategy.Bind = opts.BindStrategy

	if err := e.storage.Set(e.context, storage.EndpointKind,
		e.storage.Key().Endpoint(endpoint.Meta.Namespace, endpoint.Meta.Name), endpoint, nil); err != nil {
		log.Errorf("%s:create:> distribution update endpoint: %s err: %v", logEndpointPrefix, endpoint.SelfLink(), err)
		return nil, err
	}

	return endpoint, nil
}

func (e *Endpoint) SetSpec(endpoint *types.Endpoint, spec *types.EndpointSpec) (*types.Endpoint, error) {
	endpoint.Spec = *spec
	if err := e.storage.Set(e.context, storage.EndpointKind,
		e.storage.Key().Endpoint(endpoint.Meta.Namespace, endpoint.Meta.Name), endpoint, nil); err != nil {
		log.Errorf("%s:create:> distribution update endpoint spec: %s err: %v", logEndpointPrefix, endpoint.SelfLink(), err)
		return nil, err
	}
	return endpoint, nil
}

func (e *Endpoint) Remove(endpoint *types.Endpoint) error {
	log.Debugf("%s:remove:> remove endpoint %s", logEndpointPrefix, endpoint.Meta.Name)
	if err := e.storage.Del(e.context, storage.EndpointKind,
		e.storage.Key().Endpoint(endpoint.Meta.Namespace, endpoint.Meta.Name)); err != nil {
		log.Debugf("%s:remove:> remove endpoint %s err: %v", logEndpointPrefix, endpoint.Meta.Name, err)
		return err
	}

	return nil
}

// Watch endpoint changes
func (e *Endpoint) Watch(ch chan types.EndpointEvent) error {

	log.Debugf("%s:watch:> watch endpoint", logEndpointPrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-e.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := types.EndpointEvent{}
				res.Action = e.Action
				res.Name = e.Name

				endpoint := new(types.Endpoint)

				if err := json.Unmarshal(e.Data.([]byte), *endpoint); err != nil {
					log.Errorf("%s:> parse data err: %v", logEndpointPrefix, err)
					continue
				}

				res.Data = endpoint

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	if err := e.storage.Watch(e.context, storage.EndpointKind, watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewEndpointModel(ctx context.Context, stg storage.Storage) *Endpoint {
	return &Endpoint{ctx, stg}
}
