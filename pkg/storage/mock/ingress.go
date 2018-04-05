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

// IngressStorage Service type for interface in interfaces folder
type IngressStorage struct {
	storage.Ingress
	data map[string]*types.Ingress
}

func (s *IngressStorage) List(ctx context.Context) (map[string]*types.Ingress, error) {
	return s.data, nil
}

func (s *IngressStorage) Get(ctx context.Context, name string) (*types.Ingress, error) {

	if n, ok := s.data[name]; ok {
		return n, nil
	}

	return nil, errors.New(store.ErrEntityNotFound)
}

func (s *IngressStorage) GetSpec(ctx context.Context, ingress *types.Ingress) (*types.IngressSpec, error) {

	if err := s.checkIngressExists(ingress); err != nil {
		return nil, err
	}

	return &s.data[ingress.Meta.Name].Spec, nil
}

func (s *IngressStorage) Insert(ctx context.Context, ingress *types.Ingress) error {

	if err := s.checkIngressArgument(ingress); err != nil {
		return err
	}

	ingress.Spec.Routes = make(map[string]types.RouteSpec)

	s.data[ingress.Meta.Name] = ingress

	return nil
}

func (s *IngressStorage) Update(ctx context.Context, ingress *types.Ingress) error {

	if err := s.checkIngressExists(ingress); err != nil {
		return err
	}

	s.data[ingress.Meta.Name].Meta = ingress.Meta
	return nil
}

func (s *IngressStorage) SetStatus(ctx context.Context, ingress *types.Ingress) error {

	if err := s.checkIngressExists(ingress); err != nil {
		return err
	}

	s.data[ingress.Meta.Name].Status = ingress.Status
	return nil
}

func (s *IngressStorage) Remove(ctx context.Context, ingress *types.Ingress) error {

	if err := s.checkIngressExists(ingress); err != nil {
		return err
	}

	delete(s.data, ingress.Meta.Name)
	return nil
}

func (s *IngressStorage) Watch(ctx context.Context, ingress chan *types.Ingress) error {
	return nil
}

// Clear ingress storage
func (s *IngressStorage) Clear(ctx context.Context) error {
	s.data = make(map[string]*types.Ingress)
	return nil
}

func newIngressStorage() *IngressStorage {
	s := new(IngressStorage)
	s.data = make(map[string]*types.Ingress)
	return s
}

func (s *IngressStorage) checkIngressArgument(ingress *types.Ingress) error {
	if ingress == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if ingress.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

func (s *IngressStorage) checkIngressExists(ingress *types.Ingress) error {

	if err := s.checkIngressArgument(ingress); err != nil {
		return err
	}

	if _, ok := s.data[ingress.Meta.Name]; !ok {
		return errors.New(store.ErrEntityNotFound)
	}

	return nil
}
