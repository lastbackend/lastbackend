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
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
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

	keyProjectMeta := s.util.Key(ctx, projectStorage, projectID, "meta")
	pmeta := new(types.ProjectMeta)
	if err := client.Get(ctx, keyProjectMeta, pmeta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	keyMeta := s.util.Key(ctx, projectStorage, projectID, serviceStorage, serviceID, "meta")
	smeta := new(types.ServiceMeta)
	if err := client.Get(ctx, keyMeta, smeta); err != nil {
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

	service.ID = smeta.ID
	service.Name = smeta.Name
	service.Image = smeta.Image
	service.Description = smeta.Description
	service.Labels = smeta.Labels
	service.Created = smeta.Created
	service.Updated = smeta.Updated
	service.Project = pmeta.Name

	return service, nil
}

// Get project by name
func (s *ServiceStorage) GetByName(ctx context.Context, projectID, name string) (*types.Service, error) {
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

	keyProjectMeta := s.util.Key(ctx, projectStorage, projectID, "meta")
	pmeta := new(types.ProjectMeta)
	if err := client.Get(ctx, keyProjectMeta, pmeta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	keyServiceList := s.util.Key(ctx, projectStorage, projectID, serviceStorage)
	metaList := []types.ServiceMeta{}
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
		service.ID = meta.ID
		service.Name = meta.Name
		service.Description = meta.Description
		service.Image = meta.Image
		service.Labels = meta.Labels
		service.Created = meta.Created
		service.Updated = meta.Updated
		service.Project = pmeta.Name

		*serviceList = append(*serviceList, service)
	}

	return serviceList, nil
}

// Insert new service into storage
func (s *ServiceStorage) Insert(ctx context.Context, projectID, name, description, image string, config *types.ServiceConfig) (*types.Service, error) {
	var (
		id      = generator.GetUUIDV4()
		service = new(types.Service)
	)

	service.ID = id
	service.Name = name
	service.Project = projectID
	service.Description = description
	service.Config = config
	service.Image = image
	service.Updated = time.Now()
	service.Created = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return service, err
	}
	defer destroy()

	keyProjectMeta := s.util.Key(ctx, projectStorage, projectID, "meta")
	pmeta := new(types.ProjectMeta)
	if err := client.Get(ctx, keyProjectMeta, pmeta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	tx := client.Begin(ctx)

	keyHelper := s.util.Key(ctx, "helper", projectStorage, projectID, serviceStorage, name)
	if err := tx.Create(keyHelper, &id, 0); err != nil {
		return nil, err
	}

	keyMeta := s.util.Key(ctx, projectStorage, projectID, serviceStorage, service.ID, "meta")
	if err := tx.Create(keyMeta, service, 0); err != nil {
		return nil, err
	}

	keyConfig := s.util.Key(ctx, projectStorage, projectID, serviceStorage, service.ID, "config")
	if err := tx.Create(keyConfig, config, 0); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	service.Project = pmeta.Name

	return service, nil
}

// Update service in storage
func (s *ServiceStorage) Update(ctx context.Context, projectID string, service *types.Service) (*types.Service, error) {

	service.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return service, err
	}
	defer destroy()

	keyMeta := s.util.Key(ctx, projectStorage, projectID, serviceStorage, service.ID, "meta")
	smeta := new(types.ProjectMeta)
	if err := client.Get(ctx, keyMeta, smeta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	tx := client.Begin(ctx)

	if smeta.Name != service.Name {
		keyHelper1 := s.util.Key(ctx, "helper", projectStorage, projectID, serviceStorage, smeta.Name)
		tx.Delete(keyHelper1)

		keyHelper2 := s.util.Key(ctx, "helper", projectStorage, projectID, serviceStorage, service.Name)
		if err := tx.Create(keyHelper2, &service.ID, 0); err != nil {
			return nil, err
		}
	}

	keyMeta = s.util.Key(ctx, projectStorage, service.ID, "meta")
	if err := tx.Update(keyMeta, service.ServiceMeta, 0); err != nil {
		return nil, err
	}

	keyMeta = s.util.Key(ctx, projectStorage, service.ID, "config")
	if err := tx.Update(keyMeta, service.Config, 0); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return service, nil
}

// Remove service model
func (s *ServiceStorage) Remove(ctx context.Context, projectID, serviceID string) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	keyMeta := s.util.Key(ctx, projectStorage, projectID, serviceStorage, serviceID, "meta")
	meta := new(types.ProjectMeta)
	if err := client.Get(ctx, keyMeta, meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil
		}
		return err
	}

	tx := client.Begin(ctx)

	keyHelper := s.util.Key(ctx, "helper", projectStorage, projectID, serviceStorage, meta.Name)
	tx.Delete(keyHelper)

	keyService := s.util.Key(ctx, projectStorage, projectID, serviceStorage, serviceID)
	tx.DeleteDir(keyService)

	return tx.Commit()
}

// Remove services from project
func (s *ServiceStorage) RemoveByProject(ctx context.Context, project string) error {

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	key := s.util.Key(ctx, projectStorage, project, serviceStorage)
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
