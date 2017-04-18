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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/storage/store"
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
func (s *NamespaceStorage) GetByID(ctx context.Context, id string) (*types.Namespace, error) {

	const filter = `\b(.+)` + namespaceStorage + `\/[a-z0-9-]{36}\/(meta)\b`
	namespace := new(types.Namespace)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, namespaceStorage, id)
	if err := client.Map(ctx, key, filter, namespace); err != nil {
		return nil, err
	}

	return namespace, nil
}

// Get namespace by name
func (s *NamespaceStorage) GetByName(ctx context.Context, name string) (*types.Namespace, error) {

	var (
		id string
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, "helper", namespaceStorage, name)
	if err := client.Get(ctx, key, &id); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, id)
}

// List projects
func (s *NamespaceStorage) List(ctx context.Context) ([]*types.Namespace, error) {

	const filter = `\b(.+)` + namespaceStorage + `\/[a-z0-9-]{36}\/(meta)\b`

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, namespaceStorage)
	namespaces := []*types.Namespace{}
	if err := client.List(ctx, key, filter, &namespaces); err != nil {
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

	tx := client.Begin(ctx)

	keyHelper := s.util.Key(ctx, "helper", namespaceStorage, namespace.Meta.Name)
	if err := tx.Create(keyHelper, namespace.Meta.ID, 0); err != nil {
		return err
	}

	keyMeta := s.util.Key(ctx, namespaceStorage, namespace.Meta.ID, "meta")
	if err := tx.Create(keyMeta, namespace.Meta, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
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

	keyMeta := s.util.Key(ctx, namespaceStorage, namespace.Meta.ID, "meta")
	pmeta := new(types.Meta)
	if err := client.Get(ctx, keyMeta, pmeta); err != nil {
		return err
	}

	meta := types.Meta{}
	meta = namespace.Meta
	meta.Updated = time.Now()

	tx := client.Begin(ctx)

	if pmeta.Name != namespace.Meta.Name {
		keyHelper1 := s.util.Key(ctx, "helper", namespaceStorage, pmeta.Name)
		tx.Delete(keyHelper1)

		keyHelper2 := s.util.Key(ctx, "helper", namespaceStorage, namespace.Meta.Name)
		if err := tx.Create(keyHelper2, namespace.Meta.ID, 0); err != nil {
			return err
		}
	}

	keyMeta = s.util.Key(ctx, namespaceStorage, namespace.Meta.ID, "meta")
	if err := tx.Update(keyMeta, meta, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Remove namespace model
func (s *NamespaceStorage) Remove(ctx context.Context, id string) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	keyMeta := s.util.Key(ctx, namespaceStorage, id, "meta")
	meta := new(types.Meta)
	if err := client.Get(ctx, keyMeta, meta); err != nil {
		return err
	}

	tx := client.Begin(ctx)

	keyHelper := s.util.Key(ctx, "helper", namespaceStorage, meta.Name)
	tx.Delete(keyHelper)

	key := s.util.Key(ctx, namespaceStorage, id)
	tx.DeleteDir(key)

	if err := tx.Commit(); err != nil {
		return err
	}

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
