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
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"fmt"
	"strings"
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

type ServiceStorage struct {
	storage.Service
	data map[string]*types.Service
}

// Get service by name
func (s *ServiceStorage) Get(ctx context.Context, namespace, name string) (*types.Service, error) {
	if ns, ok := s.data[s.keyCreate(namespace, name)]; ok {
		return ns, nil
	}
	return nil, errors.New(store.ErrEntityNotFound)
}

// Get services by namespace name
func (s *ServiceStorage) ListByNamespace(ctx context.Context, namespace string) (map[string]*types.Service, error) {
	list := make(map[string]*types.Service, 0)

	prefix := fmt.Sprintf("%s:", namespace)
	for _, d := range s.data {

		if strings.HasPrefix(s.keyGet(d), prefix) {
			list[s.keyGet(d)] = d
		}
	}

	return list, nil
}

// Update service state
func (s *ServiceStorage) SetState(ctx context.Context, service *types.Service) error {
	if err := s.checkServiceExists(service); err != nil {
		return err
	}

	s.data[s.keyCreate(service.Meta.Namespace, service.Meta.Name)].State = service.State
	return nil
}

// Update service state
func (s *ServiceStorage) SetSpec(ctx context.Context, service *types.Service) error {
	if err := s.checkServiceExists(service); err != nil {
		return err
	}

	s.data[s.keyCreate(service.Meta.Namespace, service.Meta.Name)].Spec = service.Spec
	return nil
}

// Insert new service
func (s *ServiceStorage) Insert(ctx context.Context, service *types.Service) error {

	if err := s.checkServiceArgument(service); err != nil {
		return err
	}

	s.data[s.keyCreate(service.Meta.Namespace, service.Meta.Name)] = service

	return nil
}

// Update service info
func (s *ServiceStorage) Update(ctx context.Context, service *types.Service) error {

	if err := s.checkServiceExists(service); err != nil {
		return err
	}

	s.data[s.keyCreate(service.Meta.Namespace, service.Meta.Name)] = service

	return nil
}

// Remove service from storage
func (s *ServiceStorage) Remove(ctx context.Context, service *types.Service) error {

	if err := s.checkServiceExists(service); err != nil {
		return err
	}

	delete(s.data, s.keyCreate(service.Meta.Namespace, service.Meta.Name))

	return nil
}

// Watch service changes
func (s *ServiceStorage) Watch(ctx context.Context, service chan *types.Service) error {
	return nil
}

// Watch service spec changes
func (s *ServiceStorage) WatchSpec(ctx context.Context, service chan *types.Service) error {
	return nil
}

// Clear service storage
func (s *ServiceStorage) Clear(ctx context.Context) error {
	s.data = make(map[string]*types.Service)
	return nil
}

// keyCreate util function
func (s *ServiceStorage) keyCreate (namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

// keyGet util function
func (s *ServiceStorage) keyGet(svc *types.Service) string {
	return svc.SelfLink()
}

// newServiceStorage returns new storage
func newServiceStorage() *ServiceStorage {
	s := new(ServiceStorage)
	s.data = make(map[string]*types.Service)
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
func (s *ServiceStorage) checkServiceExists(service *types.Service) error {

	if err := s.checkServiceArgument(service); err != nil {
		return err
	}

	if _, ok := s.data[s.keyCreate(service.Meta.Namespace, service.Meta.Name)]; !ok {
		return errors.New(store.ErrEntityNotFound)
	}

	return nil
}