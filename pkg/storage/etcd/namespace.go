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

package etcd

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/pkg/errors"
	"time"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
)

const namespaceStorage = "namespace"

// Namespace Service type for interface in interfaces folder
type NamespaceStorage struct {
	storage.Namespace
}

// Get namespace by name
func (s *NamespaceStorage) GetByName(ctx context.Context, name string) (*types.Namespace, error) {

	log.V(logLevel).Debugf("Storage: Namespace: get by name: %s", name)

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("Storage: Namespace: get namespace err: %s", err.Error())
		return nil, err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Namespace: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	namespace := new(types.Namespace)
	keyMeta := keyCreate(namespaceStorage, name, "meta")
	if err := client.Get(ctx, keyMeta, &namespace.Meta); err != nil {
		log.V(logLevel).Errorf("Storage: Namespace: get namespace `%s` meta err: %s", name, err.Error())
		return nil, err
	}

	return namespace, nil
}

// List projects
func (s *NamespaceStorage) List(ctx context.Context) ([]*types.Namespace, error) {

	log.V(logLevel).Debug("Storage: Namespace: get namespace list")

	const filter = `\b(.+)` + namespaceStorage + `\/.+\/(meta)\b`

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Namespace: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	namespaces := []*types.Namespace{}
	keyNamespaces := keyCreate(namespaceStorage)
	if err := client.List(ctx, keyNamespaces, filter, &namespaces); err != nil {
		log.V(logLevel).Errorf("Storage: Namespace: get namespaces list err: %s", err.Error())
		return nil, err
	}

	log.V(logLevel).Debugf("Storage: Namespace: get namespace list result: %d", len(namespaces))

	return namespaces, nil
}

// Insert new namespace into storage
func (s *NamespaceStorage) Insert(ctx context.Context, namespace *types.Namespace) error {

	log.V(logLevel).Debug("Storage: Namespace: insert namespace: %#v", namespace)

	if namespace == nil {
		err := errors.New("namespace can not be nil")
		log.V(logLevel).Errorf("Storage: Namespace: insert namespace err: %s", err.Error())
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Namespace: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(namespaceStorage, namespace.Meta.Name, "meta")
	if err := client.Create(ctx, keyMeta, namespace.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Namespace: insert namespace err: %s", err.Error())
		return err
	}

	return nil
}

// Update namespace model
func (s *NamespaceStorage) Update(ctx context.Context, namespace *types.Namespace) error {

	log.V(logLevel).Debugf("Storage: Namespace: update namespace: %#v", namespace)

	if namespace == nil {
		err := errors.New("namespace can not be nil")
		log.V(logLevel).Errorf("Storage: Namespace: update namespace err: %s", err.Error())
		return err
	}

	namespace.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Namespace: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(namespaceStorage, namespace.Meta.Name, "meta")
	if err := client.Update(ctx, keyMeta, namespace.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Namespace: update namespace meta err: %s", err.Error())
		return err
	}

	return nil
}

// Remove namespace model
func (s *NamespaceStorage) Remove(ctx context.Context, name string) error {

	log.V(logLevel).Debugf("Storage: Namespace: remove namespace: %s", name)

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		log.V(logLevel).Errorf("Storage: Namespace: remove namespace err: %s", err.Error())
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Namespace: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyNamespace := keyCreate(namespaceStorage, name)
	if err := client.DeleteDir(ctx, keyNamespace); err != nil {
		log.V(logLevel).Errorf("Storage: Namespace: remove namespace `%s` err: %s", name, err.Error())
		return err
	}

	return nil
}

func newNamespaceStorage() *NamespaceStorage {
	s := new(NamespaceStorage)
	return s
}
