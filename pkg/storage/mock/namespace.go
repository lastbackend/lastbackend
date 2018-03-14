//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
	storage.Namespace
}

// Get namespace by name
func (s *NamespaceStorage) GetByName(ctx context.Context, name string) (*types.Namespace, error) {
	return new(types.Namespace), nil
}

// List projects
func (s *NamespaceStorage) List(ctx context.Context) ([]*types.Namespace, error) {
	return make([]*types.Namespace, 0), nil
}

// Insert new namespace into storage
func (s *NamespaceStorage) Insert(ctx context.Context, namespace *types.Namespace) error {
	return nil
}

// Update namespace model
func (s *NamespaceStorage) Update(ctx context.Context, namespace *types.Namespace) error {
	return nil
}

// Remove namespace model
func (s *NamespaceStorage) Remove(ctx context.Context, name string) error {
	return nil
}

func newNamespaceStorage() *NamespaceStorage {
	s := new(NamespaceStorage)
	return s
}
