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
)

const namespaceStorage = "namespace"

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
	return nil, nil
}

// List projects
func (s *NamespaceStorage) List(ctx context.Context) ([]*types.Namespace, error) {
	list := make([]*types.Namespace, 0)
	for _, ns := range s.data {
		list = append(list, ns)
	}
	return list, nil
}

// Insert new namespace into storage
func (s *NamespaceStorage) Insert(ctx context.Context, namespace *types.Namespace) error {
	if _, ok := s.data[namespace.Meta.Name]; !ok {
		s.data[namespace.Meta.Name] = namespace
	}
	return nil
}

// Update namespace model
func (s *NamespaceStorage) Update(ctx context.Context, namespace *types.Namespace) error {
	if _, ok := s.data[namespace.Meta.Name]; ok {
		s.data[namespace.Meta.Name] = namespace
	}
	return nil
}

// Remove namespace model
func (s *NamespaceStorage) Remove(ctx context.Context, name string) error {
	if _, ok := s.data[name]; ok {
		delete(s.data, name)
	}
	return nil
}

func newNamespaceStorage() *NamespaceStorage {
	s := new(NamespaceStorage)
	s.data = make(map[string]*types.Namespace)
	return s
}
