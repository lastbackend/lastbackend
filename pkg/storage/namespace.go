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
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/pkg/errors"
	"time"
)

const namespaceStorage = "namespaces"

// Namespace Service type for interface in interfaces folder
type NamespaceStorage struct {
	INamespace
	log    logger.ILogger
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get namespace by name
func (s *NamespaceStorage) GetByName(ctx context.Context, name string) (*types.Namespace, error) {

	s.log.V(debugLevel).Debugf("Storage: Namespace: get by name: %s", name)

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		s.log.V(debugLevel).Errorf("Storage: Namespace: get namespace err: %s", err.Error())
		return nil, err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(debugLevel).Errorf("Storage: Namespace: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	namespace := new(types.Namespace)
	keyMeta := s.util.Key(ctx, namespaceStorage, name, "meta")
	if err := client.Get(ctx, keyMeta, &namespace.Meta); err != nil {
		s.log.V(debugLevel).Errorf("Storage: Namespace: get namespace `%s` meta err: %s", name, err.Error())
		return nil, err
	}
	return namespace, nil
}

// List projects
func (s *NamespaceStorage) List(ctx context.Context) ([]*types.Namespace, error) {

	s.log.V(debugLevel).Debug("Storage: Namespace: get namespace list")

	const filter = `\b(.+)` + namespaceStorage + `\/.+\/(meta)\b`

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(debugLevel).Errorf("Storage: Namespace: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	namespaces := []*types.Namespace{}
	keyNamespaces := keyCreate(namespaceStorage)
	if err := client.List(ctx, keyNamespaces, filter, &namespaces); err != nil {
		s.log.V(debugLevel).Errorf("Storage: Namespace: get namespaces list err: %s", err.Error())
		return nil, err
	}

	s.log.V(debugLevel).Debugf("Storage: Namespace: get namespace list result: %d", len(namespaces))

	return namespaces, nil
}

// Insert new namespace into storage
func (s *NamespaceStorage) Insert(ctx context.Context, namespace *types.Namespace) error {

	s.log.V(debugLevel).Debug("Storage: Namespace: insert namespace: %#v", namespace)

	if namespace == nil {
		err := errors.New("namespace can not be nil")
		s.log.V(debugLevel).Errorf("Storage: Namespace: insert namespace err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(debugLevel).Errorf("Storage: Namespace: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(namespaceStorage, namespace.Meta.Name, "meta")
	if err := client.Create(ctx, keyMeta, namespace.Meta, nil, 0); err != nil {
		s.log.V(debugLevel).Errorf("Storage: Namespace: insert namespace err: %s", err.Error())
		return err
	}

	return nil
}

// Update namespace model
func (s *NamespaceStorage) Update(ctx context.Context, namespace *types.Namespace) error {

	s.log.V(debugLevel).Debug("Storage: Namespace: update namespace: %#v", namespace)

	if namespace == nil {
		err := errors.New("namespace can not be nil")
		s.log.V(debugLevel).Errorf("Storage: Namespace: update namespace err: %s", err.Error())
		return err
	}

	namespace.Meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(debugLevel).Errorf("Storage: Namespace: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	meta := types.Meta{}
	meta = namespace.Meta
	meta.Updated = time.Now()

	keyMeta := keyCreate(namespaceStorage, namespace.Meta.Name, "meta")
	if err := client.Update(ctx, keyMeta, meta, nil, 0); err != nil {
		s.log.V(debugLevel).Errorf("Storage: Namespace: update namespace meta err: %s", err.Error())
		return err
	}

	return nil
}

// Remove namespace model
func (s *NamespaceStorage) Remove(ctx context.Context, name string) error {

	s.log.V(debugLevel).Debug("Storage: Namespace: remove namespace: %s", name)

	if len(name) == 0 {
		err := errors.New("name can not be empty")
		s.log.V(debugLevel).Errorf("Storage: Namespace: remove namespace err: %s", err.Error())
		return err
	}

	client, destroy, err := s.Client()
	if err != nil {
		s.log.V(debugLevel).Errorf("Storage: Namespace: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyNamespace := keyCreate(namespaceStorage, name)
	if err := client.DeleteDir(ctx, keyNamespace); err != nil {
		s.log.V(debugLevel).Errorf("Storage: Namespace: remove namespace `%s` err: %s", name, err.Error())
		return err
	}

	return nil
}

func newNamespaceStorage(config store.Config, log logger.ILogger, util IUtil) *NamespaceStorage {
	s := new(NamespaceStorage)
	s.log = log
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config, log)
	}
	return s
}
