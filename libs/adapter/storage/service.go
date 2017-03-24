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
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	db "github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const ServiceTable string = "services"

// Service Service type for interface in interfaces folder
type ServiceStorage struct {
	storage.IService
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *ServiceStorage) CheckExistsByName(user, name string) (bool, error) {
	return false, nil
}

func (s *ServiceStorage) GetByNameOrID(user, nameOrID string) (*model.Service, error) {
	return nil, nil
}

func (s *ServiceStorage) GetByName(user, name string) (*model.Service, error) {
	return nil, nil
}

func (s *ServiceStorage) GetByID(user, id string) (*model.Service, error) {
	return nil, nil
}

func (s *ServiceStorage) ListByProject(user, project string) (*model.ServiceList, error) {
	return nil, nil
}

// Insert new service into storage
func (s *ServiceStorage) Insert(service *model.Service) (*model.Service, error) {
	return nil, nil
}

// Update service model
func (s *ServiceStorage) Update(service *model.Service) (*model.Service, error) {
	return nil, nil
}

// Remove service model
func (s *ServiceStorage) Remove(user, id string) error {
	return nil
}

// Remove service model
func (s *ServiceStorage) RemoveByProject(user, project string) error {
	return nil
}

func newServiceStorage(config store.Config) *ServiceStorage {
	s := new(ServiceStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return db.Create(config)
	}
	return s
}
