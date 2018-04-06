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

package etcd

import (
	"context"
	"time"

	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const namespaceStorage = "namespace"

// Namespace Service type for interface in interfaces folder
type NamespaceStorage struct {
	storage.Namespace
}

// Get namespace by name
func (s *NamespaceStorage) Get(ctx context.Context, name string) (*types.Namespace, error) {

	log.V(logLevel).Debugf("storage:etcd:namespace:> get by name: %s", name)

	const filter = `\b.+` + namespaceStorage + `\/.+\/(meta|state|spec)\b`

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("storage:etcd:namespace:> get by name err: %s", err.Error())
		return nil, err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:namespace:> get by name err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	namespace := new(types.Namespace)
	key := keyDirCreate(namespaceStorage, name)

	if err := client.Map(ctx, key, filter, namespace); err != nil {
		log.V(logLevel).Errorf("storage:etcd:namespace:> get by name err: %s", err.Error())
		return nil, err
	}

	return namespace, nil
}

// List projects
func (s *NamespaceStorage) List(ctx context.Context) (map[string]*types.Namespace, error) {

	log.V(logLevel).Debugf("storage:etcd:namespace:> get list")

	const filter = `\b.+` + namespaceStorage + `\/(.+)\/(meta|state|spec)\b`

	var (
		namespaces = make(map[string]*types.Namespace)
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:namespace:>  get list err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	key := keyDirCreate(namespaceStorage)
	if err := client.MapList(ctx, key, filter, namespaces); err != nil {
		log.V(logLevel).Errorf("storage:etcd:namespace:>  get list err: %s", err.Error())
		return nil, err
	}

	return namespaces, nil
}

// Insert new namespace into storage
func (s *NamespaceStorage) Insert(ctx context.Context, namespace *types.Namespace) error {

	log.V(logLevel).Debug("storage:etcd:namespace:> insert namespace: %#v", namespace)

	if err := s.checkNamespaceArgument(namespace); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Namespace: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	key := keyCreate(namespaceStorage, namespace.Meta.Name, "meta")
	if err := tx.Create(key, namespace.Meta, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:namespace:> insert namespace meta err: %s", err.Error())
		return err
	}

	keySpec := keyCreate(namespaceStorage, namespace.Meta.Name, "spec")
	if err := tx.Create(keySpec, namespace.Spec, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:namespace:> insert namespace spec err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:namespace:> insert namespace err: %s", err.Error())
		return err
	}

	return nil
}

// Update namespace model
func (s *NamespaceStorage) Update(ctx context.Context, namespace *types.Namespace) error {

	log.V(logLevel).Debug("storage:etcd:namespace:> update namespace: %#v", namespace)

	if err := s.checkNamespaceExists(ctx, namespace); err != nil {
		return err
	}

	namespace.Meta.Updated = time.Now()
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:namespace:> update namespace err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	key := keyCreate(namespaceStorage, namespace.Meta.Name, "meta")
	if err := client.Update(ctx, key, namespace.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:namespace:> update namespace err: %s", err.Error())
		return err
	}

	keySpec := keyCreate(namespaceStorage, namespace.Meta.Name, "spec")
	if err := tx.Update(keySpec, namespace.Spec, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:namespace:> update namespace spec err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:namespace:> update namespace err: %s", err.Error())
		return err
	}

	return nil
}

// Remove namespace model
func (s *NamespaceStorage) Remove(ctx context.Context, namespace *types.Namespace) error {

	log.V(logLevel).Debug("storage:etcd:namespace:> remove namespace: %#v", namespace)

	if err := s.checkNamespaceExists(ctx, namespace); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> remove namespace err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyDirCreate(namespaceStorage, s.keyGet(namespace))

	if err := client.DeleteDir(ctx, keyMeta); err != nil {
		log.V(logLevel).Errorf("storage:etcd:deployment:> remove namespace err: %s", err.Error())
		return err
	}

	return nil
}

// Clear namespace storage
func (s *NamespaceStorage) Clear(ctx context.Context) error {

	log.V(logLevel).Debugf("storage:etcd:namespace:> clear")

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:namespace:> clear err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := client.DeleteDir(ctx, namespaceStorage); err != nil {
		log.V(logLevel).Errorf("storage:etcd:namespace:> clear err: %s", err.Error())
		return err
	}

	return nil
}

// keyCreate util function
func (s *NamespaceStorage) keyCreate(name string) string {
	return fmt.Sprintf("%s", name)
}

// keyGet util function
func (s *NamespaceStorage) keyGet(namespace *types.Namespace) string {
	return namespace.SelfLink()
}

func newNamespaceStorage() *NamespaceStorage {
	s := new(NamespaceStorage)
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

func (s *NamespaceStorage) checkNamespaceExists(ctx context.Context, namespace *types.Namespace) error {

	if err := s.checkNamespaceArgument(namespace); err != nil {
		return err
	}

	log.V(logLevel).Debugf("storage:etcd:namespace:> check namespace exists")

	if _, err := s.Get(ctx, namespace.Meta.Name); err != nil {
		log.V(logLevel).Debugf("storage:etcd:namespace:> check namespace exists err: %s", err.Error())
		return err
	}

	return nil
}
