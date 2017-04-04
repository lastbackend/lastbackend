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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"golang.org/x/net/context"
	"time"
)

const ServiceTable string = "services"

// Service Service type for interface in interfaces folder
type ServiceStorage struct {
	IService
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get project by name for user
func (s *ServiceStorage) GetByID(username, projectID, serviceID string) (*types.Service, error) {
	var (
		project        = new(types.Project)
		service        = new(types.Service)
		keyProjectMeta = fmt.Sprintf("%s/%s/%s/meta", ProjectTable, username, projectID)
		keyMeta        = fmt.Sprintf("%s/%s/%s/%s/%s/meta", ProjectTable, username, projectID, ServiceTable, serviceID)
		keyConfig      = fmt.Sprintf("%s/%s/%s/%s/%s/config", ProjectTable, username, projectID, ServiceTable, serviceID)
		keySource      = fmt.Sprintf("%s/%s/%s/%s/%s/source", ProjectTable, username, projectID, ServiceTable, serviceID)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Get(ctx, keyProjectMeta, &project.Meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	if err := client.Get(ctx, keyMeta, &service.Meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	service.Config = new(types.ServiceConfig)
	if err := client.Get(ctx, keyConfig, service.Config); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	service.Source = new(types.ServiceSource)
	if err := client.Get(ctx, keySource, service.Source); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	service.User = username
	service.Project = project.Name

	return service, nil
}

// Get project by name for user
func (s *ServiceStorage) GetByName(username, projectID, name string) (*types.Service, error) {
	var (
		id string
		// Key example: /helper/projects/<username>/<project id>/services/<name>
		key = fmt.Sprintf("/helper/%s/%s/%s/%s/%s", ProjectTable, username, projectID, ServiceTable, name)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Get(ctx, key, &id); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return s.GetByID(username, projectID, id)
}

// List project by username
func (s *ServiceStorage) ListByProject(username, projectID string) (*types.ServiceList, error) {
	var (
		projects    = make(map[string]*types.Project)
		keyProjects = fmt.Sprintf("%s/%s", ProjectTable, username)
		key         = fmt.Sprintf("%s/%s/%s/%s", ProjectTable, username, projectID, ServiceTable)
		filter      = `\b(.+)\/info\b`
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	metaList := []types.Meta{}

	if err := client.Map(ctx, keyProjects, ``, projects); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	if err := client.List(ctx, key, filter, &metaList); err != nil {
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
		*serviceList = append(*serviceList, types.Service{Meta: meta, User: username, Project: projects[projectID].Name})
	}

	return serviceList, nil
}

// Insert new service into storage
func (s *ServiceStorage) Insert(username, projectID, name, description string, source *types.ServiceSource, config *types.ServiceConfig) (*types.Service, error) {
	var (
		id             = generator.GetUUIDV4()
		project        = new(types.Project)
		service        = new(types.Service)
		keyProjectMeta = fmt.Sprintf("%s/%s/%s/meta", ProjectTable, username, projectID)
		keyHelper      = fmt.Sprintf("/helper/%s/%s/%s/%s/%s", ProjectTable, username, projectID, ServiceTable, name)
		keyMeta        = fmt.Sprintf("%s/%s/%s/%s/%s/meta", ProjectTable, username, projectID, ServiceTable, id)
		keyConfig      = fmt.Sprintf("%s/%s/%s/%s/%s/config", ProjectTable, username, projectID, ServiceTable, id)
		keySource      = fmt.Sprintf("%s/%s/%s/%s/%s/source", ProjectTable, username, projectID, ServiceTable, id)
	)

	service.ID = id
	service.Name = name
	service.Project = projectID
	service.User = username
	service.Description = description
	service.Config = config
	service.Source = source
	service.Updated = time.Now()
	service.Created = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return service, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Get(ctx, keyProjectMeta, &project.Meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	tx := client.Begin(ctx)

	if err := tx.Create(keyHelper, &id, 0); err != nil {
		return nil, err
	}

	if err := tx.Create(keyMeta, service, 0); err != nil {
		return nil, err
	}

	if err := tx.Create(keyConfig, config, 0); err != nil {
		return nil, err
	}

	if err := tx.Create(keySource, source, 0); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return service, nil
}

// Update service in storage
func (s *ServiceStorage) Update(username, projectID string, service *types.Service) (*types.Service, error) {
	var (
		keyMeta   = fmt.Sprintf("%s/%s/%s/%s/%s/meta", ProjectTable, username, projectID, ServiceTable, service.Name)
		keyConfig = fmt.Sprintf("%s/%s/%s/%s/%s/config", ProjectTable, username, projectID, ServiceTable, service.Name)
		keySource = fmt.Sprintf("%s/%s/%s/%s/%s/source", ProjectTable, username, projectID, ServiceTable, service.Name)
	)

	service.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return service, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx := client.Begin(ctx)

	if err := tx.Update(keyMeta, service, 0); err != nil {
		return nil, err
	}

	if err := tx.Update(keyConfig, service.Config, 0); err != nil {
		return nil, err
	}

	if err := tx.Update(keySource, service.Source, 0); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return service, nil
}

// Remove service model
func (s *ServiceStorage) Remove(username, projectID, serviceID string) error {

	var (
		project        = new(types.Project)
		keyProjectMeta = fmt.Sprintf("%s/%s/%s/meta", ProjectTable, username, projectID)
		key            = fmt.Sprintf("%s/%s/%s/%s/%s", ProjectTable, username, projectID, ServiceTable, serviceID)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Get(ctx, keyProjectMeta, &project.Meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil
		}
		return err
	}

	var keyHelper = fmt.Sprintf("/helper/%s/%s/%s/%s/%s", ProjectTable, username, projectID, ServiceTable, project.Name)

	tx := client.Begin(ctx)

	tx.Delete(keyHelper)
	tx.Delete(key)

	return tx.Commit()
}

// Remove service model
func (s *ServiceStorage) RemoveByProject(username, project string) error {

	var (
		key = fmt.Sprintf("%s/%s/%s/%s", ProjectTable, username, project, ServiceTable)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Delete(ctx, key, nil); err != nil {
		return err
	}

	return nil
}

func newServiceStorage(config store.Config) *ServiceStorage {
	s := new(ServiceStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
