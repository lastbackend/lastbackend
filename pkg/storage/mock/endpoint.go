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

package mock

import (
	"context"
	"errors"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"strings"
)

// Service Endpoint type for interface in interfaces folder
type EndpointStorage struct {
	storage.Endpoint
	data map[string]*types.Endpoint
}

// Get endpoints by id
func (s *EndpointStorage) Get(ctx context.Context, namespace, service, name string) (*types.Endpoint, error) {

	if n, ok := s.data[s.keyCreate(namespace, service, name)]; ok {
		return n, nil
	}

	return nil, errors.New(store.ErrEntityNotFound)
}

// Get endpoint by namespace name
func (s *EndpointStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Endpoint, error) {

	list := make(map[string]*types.Endpoint, 0)

	prefix := fmt.Sprintf("%s:", namespace)
	for _, d := range s.data {

		if strings.HasPrefix(s.keyGet(d), prefix) {
			list[s.keyGet(d)] = d
		}
	}

	return list, nil
}

// Get endpoint by service name
func (s *EndpointStorage) ListByService(ctx context.Context, namespace, service string) (map[string]*types.Endpoint, error) {

	list := make(map[string]*types.Endpoint, 0)

	prefix := fmt.Sprintf("%s:%s:", namespace, service)

	for _, d := range s.data {
		if strings.HasPrefix(s.keyGet(d), prefix) {
			list[s.keyGet(d)] = d
		}
	}

	return list, nil
}

// Update endpoint status
func (s *EndpointStorage) SetStatus(ctx context.Context, endpoint *types.Endpoint) error {

	if err := s.checkEndpointExists(endpoint); err != nil {
		return err
	}

	s.data[s.keyGet(endpoint)].Status = endpoint.Status
	return nil
}

// Update endpoint status
func (s *EndpointStorage) SetSpec(ctx context.Context, endpoint *types.Endpoint) error {

	if err := s.checkEndpointExists(endpoint); err != nil {
		return err
	}

	s.data[s.keyGet(endpoint)].Spec = endpoint.Spec
	return nil
}

// Insert new endpoint into storage
func (s *EndpointStorage) Insert(ctx context.Context, endpoint *types.Endpoint) error {

	if err := s.checkEndpointArgument(endpoint); err != nil {
		return err
	}

	s.data[s.keyGet(endpoint)] = endpoint

	return nil
}

// Update endpoint info
func (s *EndpointStorage) Update(ctx context.Context, endpoint *types.Endpoint) error {

	if err := s.checkEndpointExists(endpoint); err != nil {
		return err
	}

	s.data[s.keyGet(endpoint)] = endpoint

	return nil
}

// Remove endpoint by id from storage
func (s *EndpointStorage) Remove(ctx context.Context, endpoint *types.Endpoint) error {

	if err := s.checkEndpointExists(endpoint); err != nil {
		return err
	}

	delete(s.data, s.keyGet(endpoint))

	return nil
}

// Watch endpoint changes
func (s *EndpointStorage) Watch(ctx context.Context, endpoint chan *types.Endpoint) error {

	return nil
}

// Watch endpoint spec changes
func (s *EndpointStorage) WatchSpec(ctx context.Context, endpoint chan *types.Endpoint) error {

	return nil
}

// Watch endpoint status changes
func (s *EndpointStorage) WatchStatus(ctx context.Context, endpoint chan *types.Endpoint) error {

	return nil
}

// Clear endpoint storage
func (s *EndpointStorage) Clear(ctx context.Context) error {

	s.data = make(map[string]*types.Endpoint)
	return nil
}

// keyCreate util function
func (s *EndpointStorage) keyCreate(namespace, service, name string) string {
	return fmt.Sprintf("%s:%s:%s", namespace, service, name)
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
func (s *EndpointStorage) checkEndpointExists(endpoint *types.Endpoint) error {

	if err := s.checkEndpointArgument(endpoint); err != nil {
		return err
	}

	if _, ok := s.data[s.keyGet(endpoint)]; !ok {
		return errors.New(store.ErrEntityNotFound)
	}

	return nil
}
