package mock

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
)

const ServiceTable string = "services"

// Service Service type for interface in interfaces folder
type ServiceMock struct {
	Mock *r.Mock
	storage.IService
}

func (s *ServiceMock) GetByNameOrID(user, project, nameOrID string) (*model.Service, *e.Err) {
	return nil, nil
}

func (s *ServiceMock) GetByID(user, project, id string) (*model.Service, *e.Err) {
	return nil, nil
}

func (s *ServiceMock) ListByProject(user, project string) (*model.ServiceList, *e.Err) {
	return nil, nil
}

// Insert new service into storage
func (s *ServiceMock) Insert(service *model.Service) (*model.Service, *e.Err) {
	return nil, nil
}

// Update service model
func (s *ServiceMock) Update(service *model.Service) (*model.Service, *e.Err) {
	return nil, nil
}

// Remove service model
func (s *ServiceMock) Remove(user, project, id string) *e.Err {
	return nil
}

// Remove service model
func (s *ServiceMock) RemoveByProject(user, project string) *e.Err {
	return nil
}

func newServiceMock(mock *r.Mock) *ServiceMock {
	s := new(ServiceMock)
	s.Mock = mock
	return s
}
