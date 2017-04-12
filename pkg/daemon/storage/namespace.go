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
	"github.com/satori/go.uuid"
	"time"
)

const namespaceStorage = "namespace"

// Namespace Service type for interface in interfaces folder
type NamespaceStorage struct {
	INamespace
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get project by name
func (s *NamespaceStorage) GetByID(ctx context.Context, id string) (*types.Namespace, error) {

	namespace := new(types.Namespace)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, namespaceStorage, id, "meta")
	meta := &types.Meta{}

	if err := client.Get(ctx, key, meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	namespace.Meta = *meta

	return namespace, nil
}

// Get project by name
func (s *NamespaceStorage) GetByName(ctx context.Context, name string) (*types.Namespace, error) {

	var id string

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, "helper", namespaceStorage, name)
	if err := client.Get(ctx, key, &id); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return s.GetByID(ctx, id)
}

// List projects
func (s *NamespaceStorage) List(ctx context.Context) (*types.NamespaceList, error) {

	const filter = `\b(.+)projects\/[a-z0-9-]{36}\/meta\b`

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, namespaceStorage)
	metaList := []types.Meta{}
	if err := client.List(ctx, key, filter, &metaList); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	if metaList == nil {
		return nil, nil
	}

	namespaceList := new(types.NamespaceList)
	for _, meta := range metaList {
		project := types.Namespace{}
		project.Meta = meta
		*namespaceList = append(*namespaceList, project)
	}

	return namespaceList, nil
}

// Insert new project into storage
func (s *NamespaceStorage) Insert(ctx context.Context, name, description string) (*types.Namespace, error) {
	var (
		id        = uuid.NewV4().String()
		namespace = new(types.Namespace)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyHelper := s.util.Key(ctx, "helper", namespaceStorage, name)
	if err := tx.Create(keyHelper, id, 0); err != nil {
		return nil, err
	}

	namespace.Meta = types.Meta{
		ID:          id,
		Name:        name,
		Description: description,
		Labels:      map[string]string{"tier": "namespace"},
		Updated:     time.Now(),
		Created:     time.Now(),
	}

	keyMeta := s.util.Key(ctx, namespaceStorage, id, "meta")
	if err := tx.Create(keyMeta, namespace.Meta, 0); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return namespace, nil
}

// Update project model
func (s *NamespaceStorage) Update(ctx context.Context, namespace *types.Namespace) (*types.Namespace, error) {

	namespace.Meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return namespace, err
	}
	defer destroy()

	keyMeta := s.util.Key(ctx, namespaceStorage, namespace.Meta.ID, "meta")
	pmeta := new(types.Meta)
	if err := client.Get(ctx, keyMeta, pmeta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
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
			return namespace, err
		}
	}

	keyMeta = s.util.Key(ctx, namespaceStorage, namespace.Meta.ID, "meta")
	if err := tx.Update(keyMeta, meta, 0); err != nil {
		return namespace, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return namespace, nil
}

// Remove project model
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
