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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

// Namespace Service type for interface in interfaces folder
type NamespaceStorage struct {
	data map[string]*types.Namespace
	storage.Namespace
}

// Get namespace by name
func (s *NamespaceStorage) Get(ctx context.Context, name string) (*types.Namespace, error) {
	if ns, ok := s.data[name]; ok {
		return ns, nil
	}
	return nil, errors.New(store.ErrEntityNotFound)
}

// List projects
func (s *NamespaceStorage) List(ctx context.Context) (map[string]*types.Namespace, error) {
	return s.data, nil
}

// Insert new namespace into storage
func (s *NamespaceStorage) Insert(ctx context.Context, namespace *types.Namespace) error {

	if err := s.checkNamespaceArgument(namespace); err != nil {
		return err
	}

	s.data[namespace.Meta.Name] = namespace

	return nil
}

// Update namespace model
func (s *NamespaceStorage) Update(ctx context.Context, namespace *types.Namespace) error {

	if err := s.checkNamespaceExists(namespace); err != nil {
		return err
	}

	s.data[namespace.Meta.Name] = namespace
	return nil
}

// Remove namespace model
func (s *NamespaceStorage) Remove(ctx context.Context, namespace *types.Namespace) error {

	if err := s.checkNamespaceExists(namespace); err != nil {
		return err
	}

	delete(s.data, namespace.Meta.Name)
	return nil
}

func newNamespaceStorage() *NamespaceStorage {
	s := new(NamespaceStorage)
	s.data = make(map[string]*types.Namespace)
	return s
}

func (s *NamespaceStorage) checkNamespaceArgument(namespace *types.Namespace) error {
	if namespace == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if namespace.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

func (s *NamespaceStorage) checkNamespaceExists(namespace *types.Namespace) error {

	if err := s.checkNamespaceArgument(namespace); err != nil {
		return err
	}

	if _, ok := s.data[namespace.Meta.Name]; !ok {
		return errors.New(store.ErrEntityNotFound)
	}

	return nil
}
