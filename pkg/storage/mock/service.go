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
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const serviceStorage string = "services"

// Service Service type for interface in interfaces folder
type ServiceStorage struct {
	data map[string]map[string]*types.Service
	storage.Service
}

// Get service by name
func (s *ServiceStorage) Get(ctx context.Context, namespace, name string) (*types.Service, error) {
	if _, ok := s.data[namespace]; !ok {
		return nil, nil
	}
	if srv, ok := s.data[namespace][name]; ok {
		return srv, nil
	}
	return nil, nil
}

// Get service by pod name
func (s *ServiceStorage) GetByPodName(ctx context.Context, name string) (*types.Service, error) {
	return new(types.Service), nil
}

// List services
func (s *ServiceStorage) ListByNamespace(ctx context.Context, namespace string) ([]*types.Service, error) {
	return make([]*types.Service, 0), nil
}

// Count services
func (s *ServiceStorage) CountByNamespace(ctx context.Context, namespace string) (int, error) {
	return 0, nil
}

// Insert new service into storage
func (s *ServiceStorage) Insert(ctx context.Context, service *types.Service) error {
	if _, ok := s.data[service.Meta.Namespace]; !ok {
		s.data[service.Meta.Namespace] = make(map[string]*types.Service)
	}
	if _, ok := s.data[service.Meta.Namespace][service.Meta.Name]; ok {
		return errors.New(store.ErrEntityNotFound)
	} else {
		s.data[service.Meta.Namespace][service.Meta.Name] = service
	}
	return nil
}

// Update service in storage
func (s *ServiceStorage) Update(ctx context.Context, service *types.Service) error {
	return nil
}

// Update service spec in storage
func (s *ServiceStorage) UpdateSpec(ctx context.Context, service *types.Service) error {
	return nil
}

// Remove service model
func (s *ServiceStorage) Remove(ctx context.Context, service *types.Service) error {
	return nil
}

// Remove services from namespace
func (s *ServiceStorage) RemoveByNamespace(ctx context.Context, namespace string) error {
	return nil
}

func (s *ServiceStorage) Watch(ctx context.Context, service chan *types.Service) error {
	return nil
}

func (s *ServiceStorage) SpecWatch(ctx context.Context, service chan *types.Service) error {
	return nil
}

func (s *ServiceStorage) PodsWatch(ctx context.Context, service chan *types.Service) error {
	return nil
}

// Update service state
func (s *ServiceStorage) updateState(ctx context.Context, service *types.Service) error {
	return nil
}

func newServiceStorage() *ServiceStorage {
	s := new(ServiceStorage)
	s.data = make(map[string]map[string]*types.Service)
	return s
}
