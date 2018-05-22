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
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const (
	logEndpointPrefix = "distribution:endpoint"
)

type IEndpoint interface {
	Get(namespace, service string) (*types.Endpoint, error)
	ListByNamespace(namespace string) (map[string]*types.Endpoint, error)
	Create(namespace string, service string, opts *types.EndpointCreateOptions) (*types.Endpoint, error)
	Update(endpoint *types.Endpoint, opts *types.EndpointUpdateOptions) (*types.Endpoint, error)
	Remove(endpoint *types.Endpoint) error
}

type Endpoint struct {
	context context.Context
	storage storage.Storage
}

func (e *Endpoint) Get(namespace, service string) (*types.Endpoint, error) {

	log.V(logLevel).Debugf("%s:get:> get endpoint by name %s: %s", logEndpointPrefix, namespace, service)

	hook, err := e.storage.Endpoint().Get(e.context, namespace, service)
	if err != nil {

		if err.Error() == store.ErrEntityNotFound {
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> create endpoint err: %s", logEndpointPrefix, err.Error())
		return nil, err
	}

	return hook, nil
}

func (e *Endpoint) ListByNamespace(namespace string) (map[string]*types.Endpoint, error) {
	log.Debugf("%s:listbynamespace:> in namespace: %s", namespace)

	el, err := e.storage.Endpoint().ListByNamespace(e.context, namespace)
	if err != nil {
		log.Errorf("%s:listbynamespace:> in namespace: %s err: %s", logEndpointPrefix, namespace, err.Error())
		return nil, err
	}

	return el, nil
}

func (e *Endpoint) Create(namespace, service string, opts *types.EndpointCreateOptions) (*types.Endpoint, error) {
	endpoint := new(types.Endpoint)

	endpoint.Meta.Name = service
	endpoint.Meta.Namespace = namespace
	endpoint.Meta.SetDefault()
	endpoint.SelfLink()

	endpoint.Status.State = types.StateCreated
	endpoint.Status.Message = ""
	endpoint.Spec.PortMap = make(map[int]string, 0)
	endpoint.Spec.Upstreams = make([]string, 0)

	for k, v := range opts.Ports {
		endpoint.Spec.PortMap[k] = v
	}

	endpoint.Spec.Policy = opts.Policy
	endpoint.Spec.Strategy.Route = opts.RouteStrategy
	endpoint.Spec.Strategy.Bind = opts.BindStrategy

	ip, err := envs.Get().GetIPAM().Lease()
	if err != nil {
		log.Errorf("%s:create:> distribution create endpoint: %s err: %s", logEndpointPrefix, endpoint.SelfLink(), err.Error())
	}

	endpoint.Spec.IP = ip.String()

	if err := e.storage.Endpoint().Insert(e.context, endpoint); err != nil {
		log.Errorf("%s:create:> distribution create endpoint: %s err: %s", logEndpointPrefix, endpoint.SelfLink(), err.Error())
		return nil, err
	}

	return endpoint, nil
}

func (e *Endpoint) Update(endpoint *types.Endpoint, opts *types.EndpointUpdateOptions) (*types.Endpoint, error) {
	log.Debugf("%s:update:> endpoint: %s", logEndpointPrefix, endpoint.SelfLink())

	if len(opts.Ports) != 0 {
		endpoint.Spec.PortMap = make(map[int]string, 0)
		for k, v := range opts.Ports {
			endpoint.Spec.PortMap[k] = v
		}
	}

	endpoint.Spec.Policy = opts.Policy
	endpoint.Spec.Strategy.Route = opts.RouteStrategy
	endpoint.Spec.Strategy.Bind = opts.BindStrategy

	if err := e.storage.Endpoint().Update(e.context, endpoint); err != nil {
		log.Errorf("%s:create:> distribution update endpoint: %s err: %s", logEndpointPrefix, endpoint.SelfLink(), err.Error())
		return nil, err
	}

	return endpoint, nil
}

func (e *Endpoint) Remove(endpoint *types.Endpoint) error {
	log.Debugf("%s:remove:> remove endpoint %s", logEndpointPrefix, endpoint.Meta.Name)
	if err := e.storage.Endpoint().Remove(e.context, endpoint); err != nil {
		log.Debugf("%s:remove:> remove endpoint %s err: %s", logEndpointPrefix, endpoint.Meta.Name, err.Error())
		return err
	}

	return nil
}

func NewEndpointModel(ctx context.Context, stg storage.Storage) IEndpoint {
	return &Endpoint{ctx, stg}
}
