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
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	db "github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	st "github.com/lastbackend/lastbackend/pkg/storage/store"
	"golang.org/x/net/context"
	"time"
)

const ServiceTable string = "services"

// Service Service type for interface in interfaces folder
type ServiceStorage struct {
	storage.IService
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get project by name for user
func (s *ServiceStorage) GetByName(username, project, name string) (*model.Service, error) {
	var (
		service = new(model.Service)
		key     = fmt.Sprintf("/%s/%s/%s/%s/%s/info", ProjectTable, username, project, name, ServiceTable)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Get(ctx, key, service); err != nil {
		if err.Error() == st.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return service, nil
}

// List project by username
func (s *ServiceStorage) ListByProject(username, project string) (*model.ServiceList, error) {
	var (
		serviceList = new(model.ServiceList)
		key         = fmt.Sprintf("/%s/%s/%s/%s", ProjectTable, username, project, ServiceTable)
		filter      = `\b(.+)\/info\b`
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.List(ctx, key, filter, serviceList); err != nil {
		if err.Error() == st.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return serviceList, nil
}

// Insert new service into storage
func (s *ServiceStorage) Insert(username, name, description string) (*model.Service, error) {
	var (
		service = new(model.Service)
		keyInfo = fmt.Sprintf("%s/%s/%s/info", ProjectTable, username, name)
	)

	service.Name = name
	service.User = username
	service.Description = description
	service.Updated = time.Now()
	service.Created = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return service, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Create(ctx, keyInfo, service, nil, 0); err != nil {
		return service, err
	}

	return service, nil
}

// Update service model
func (s *ServiceStorage) Update(service *model.Service) (*model.Service, error) {
	return nil, nil
}

// Remove service model
func (s *ServiceStorage) Remove(username, project, name string) error {

	var (
		key = fmt.Sprintf("%s/%s/%s/%s/%s", ProjectTable, username, project, ServiceTable, name)
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
		return db.Create(config)
	}
	return s
}
