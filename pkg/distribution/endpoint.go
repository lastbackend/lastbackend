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
)

const (
	logEndpointPrefix = "distribution:endpoint"
)

type IEndpoint interface {
	Get(namespace, service, name string) (*types.Endpoint, error)
	ListByNamespace(namespace string) (map[string]*types.Endpoint, error)
	ListByService(namespace, service string) (map[string]*types.Endpoint, error)
	Create(namespace *types.Namespace, opts *types.EndpointCreateOptions) (*types.Endpoint, error)
	Update(endpoint *types.Endpoint, opts *types.EndpointUpdateOptions) (*types.Endpoint, error)
	SetStatus(endpoint *types.Endpoint, status *types.EndpointStatus) error
	Remove(endpoint *types.Endpoint) error
}

type Endpoint struct {
	context context.Context
	storage storage.Storage
}

func (h *Endpoint) Get(namespace, service, name string) (*types.Endpoint, error) {

	log.V(logLevel).Debugf("%s:get:> get endpoint by name %s: %s", logEndpointPrefix, namespace, name)

	hook, err := h.storage.Endpoint().Get(h.context, namespace, service, name)
	if err != nil {
		log.V(logLevel).Errorf("%s:get:> create endpoint err: %s", logEndpointPrefix, err.Error())
		return nil, err
	}

	return hook, nil
}



func NewEndpointModel(ctx context.Context, stg storage.Storage) IEndpoint {
	return &Endpoint{ctx, stg}
}
