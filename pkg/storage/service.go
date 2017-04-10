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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/satori/go.uuid"
	"time"
)

const serviceStorage string = "services"

// Service Service type for interface in interfaces folder
type ServiceStorage struct {
	IService
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get project by name
func (s *ServiceStorage) GetByID(ctx context.Context, projectID, serviceID string) (*types.Service, error) {
	var service = new(types.Service)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	keyMeta := s.util.Key(ctx, projectStorage, projectID, serviceStorage, serviceID, "meta")
	service.Meta = types.Meta{}
	if err := client.Get(ctx, keyMeta, &service.Meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	keyConfig := s.util.Key(ctx, projectStorage, projectID, serviceStorage, serviceID, "config")
	service.Config = new(types.ServiceConfig)
	if err := client.Get(ctx, keyConfig, service.Config); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}
	return service, nil
}

// Get project by name
func (s *ServiceStorage) GetByName(ctx context.Context, projectID string, name string) (*types.Service, error) {
	var id string

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, "helper", projectStorage, projectID, serviceStorage, name)
	if err := client.Get(ctx, key, &id); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}
	return s.GetByID(ctx, projectID, id)
}

// List project
func (s *ServiceStorage) ListByProject(ctx context.Context, projectID string) (*types.ServiceList, error) {

	const filter = `\b(.+)services\/[a-z0-9-]{36}\/meta\b`

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	keyServiceList := s.util.Key(ctx, projectStorage, projectID, serviceStorage)
	metaList := []types.Meta{}
	if err := client.List(ctx, keyServiceList, filter, &metaList); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	if metaList == nil {
		return nil, nil
	}

	serviceList := new(types.ServiceList)
	for _, meta := range metaList {
		service := types.Service{}
		service.Meta.ID = meta.ID
		service.Meta.Name = meta.Name
		service.Meta.Description = meta.Description
		service.Meta.Labels = meta.Labels
		service.Meta.Created = meta.Created
		service.Meta.Updated = meta.Updated

		*serviceList = append(*serviceList, service)
	}

	return serviceList, nil
}

// Insert new service into storage
func (s *ServiceStorage) Insert(ctx context.Context, projectID string, name, description string, config *types.ServiceConfig) (*types.Service, error) {
	var (
		service = new(types.Service)
	)

	service.Meta = types.Meta{
		ID:          uuid.NewV4().String(),
		Name:        name,
		Description: description,
		Updated:     time.Now(),
		Created:     time.Now(),
	}

	service.Config = config

	client, destroy, err := s.Client()
	if err != nil {
		return service, err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyHelper := s.util.Key(ctx, "helper", projectStorage, projectID, serviceStorage, name)
	if err := tx.Create(keyHelper, &service.Meta.ID, 0); err != nil {
		return nil, err
	}

	keyMeta := s.util.Key(ctx, projectStorage, projectID, serviceStorage, service.Meta.ID, "meta")
	if err := tx.Create(keyMeta, &service.Meta, 0); err != nil {
		return nil, err
	}

	keyConfig := s.util.Key(ctx, projectStorage, projectID, serviceStorage, service.Meta.ID, "config")
	if err := tx.Create(keyConfig, config, 0); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return service, nil
}

// Update service in storage
func (s *ServiceStorage) Update(ctx context.Context, projectID string, service *types.Service) (*types.Service, error) {

	service.Meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return service, err
	}
	defer destroy()

	keyMeta := s.util.Key(ctx, projectStorage, projectID, serviceStorage, service.Meta.ID, "meta")
	smeta := new(types.Meta)
	if err := client.Get(ctx, keyMeta, smeta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	tx := client.Begin(ctx)

	if smeta.Name != service.Meta.Name {
		keyHelper1 := s.util.Key(ctx, "helper", projectStorage, projectID, serviceStorage, smeta.Name)
		tx.Delete(keyHelper1)

		keyHelper2 := s.util.Key(ctx, "helper", projectStorage, projectID, serviceStorage, service.Meta.Name)
		if err := tx.Create(keyHelper2, &service.Meta.ID, 0); err != nil {
			return nil, err
		}
	}

	keyMeta = s.util.Key(ctx, projectStorage, projectID, serviceStorage, service.Meta.ID, "meta")
	if err := tx.Update(keyMeta, service.Meta, 0); err != nil {
		return nil, err
	}

	keyMeta = s.util.Key(ctx, projectStorage, projectID, serviceStorage, service.Meta.ID, "config")
	if err := tx.Update(keyMeta, service.Config, 0); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return service, nil
}

// Remove service model
func (s *ServiceStorage) Remove(ctx context.Context, projectID string, service *types.Service) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	keyMeta := s.util.Key(ctx, projectStorage, projectID, serviceStorage, service.Meta.ID, "meta")
	meta := types.Meta{}
	if err := client.Get(ctx, keyMeta, &meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil
		}
		return err
	}

	fmt.Println(meta)
	tx := client.Begin(ctx)

	keyHelper := s.util.Key(ctx, "helper", projectStorage, projectID, serviceStorage, meta.Name)
	tx.Delete(keyHelper)

	keyService := s.util.Key(ctx, projectStorage, projectID, serviceStorage, service.Meta.ID)
	tx.DeleteDir(keyService)

	return tx.Commit()
}

// Remove services from project
func (s *ServiceStorage) RemoveByProject(ctx context.Context, projectID string) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	key := s.util.Key(ctx, projectStorage, projectID, serviceStorage)
	if err := client.DeleteDir(ctx, key); err != nil {
		return err
	}

	return nil
}

func newServiceStorage(config store.Config, util IUtil) *ServiceStorage {
	s := new(ServiceStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
