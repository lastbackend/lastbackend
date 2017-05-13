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

package storage

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"time"
)

const namespaceStorage = "namespace"

// Namespace Service type for interface in interfaces folder
type NamespaceStorage struct {
	INamespace
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get namespace by name
func (s *NamespaceStorage) GetByName(ctx context.Context, name string) (*types.Namespace, error) {

	const filter = `\b(.+)` + namespaceStorage + `\/.+\/(meta)\b`
	namespace := new(types.Namespace)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	keyNamespace := keyPrepare(namespaceStorage, name)
	if err := client.Map(ctx, keyNamespace, filter, namespace); err != nil {
		return nil, err
	}

	return namespace, nil
}

// List projects
func (s *NamespaceStorage) List(ctx context.Context) ([]*types.Namespace, error) {

	const filter = `\b(.+)` + namespaceStorage + `\/.+\/(meta)\b`

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	keyNamespaces := keyPrepare(namespaceStorage)
	namespaces := []*types.Namespace{}
	if err := client.List(ctx, keyNamespaces, filter, &namespaces); err != nil {
		return nil, err
	}

	return namespaces, nil
}

// Insert new namespace into storage
func (s *NamespaceStorage) Insert(ctx context.Context, namespace *types.Namespace) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	keyMeta := keyPrepare(namespaceStorage, namespace.Meta.Name, "meta")
	if err := client.Create(ctx, keyMeta, namespace.Meta, nil, 0); err != nil {
		return err
	}

	return nil
}

// Update namespace model
func (s *NamespaceStorage) Update(ctx context.Context, namespace *types.Namespace) error {

	namespace.Meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	meta := types.Meta{}
	meta = namespace.Meta
	meta.Updated = time.Now()

	keyMeta := keyPrepare(namespaceStorage, namespace.Meta.Name, "meta")
	if err := client.Update(ctx, keyMeta, meta, nil, 0); err != nil {
		return err
	}

	return nil
}

// Remove namespace model
func (s *NamespaceStorage) Remove(ctx context.Context, name string) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	keyNamespace := keyPrepare(namespaceStorage, name)
	client.DeleteDir(ctx, keyNamespace)

	return nil
}

func newNamespaceStorage(config store.Config, util IUtil) *NamespaceStorage {
	s := new(NamespaceStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
