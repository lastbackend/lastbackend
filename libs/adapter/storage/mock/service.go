package mock

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
	"time"
)

const mockServiceID string = "mocked"
const ServiceTable string = "services"

// Service Service type for interface in interfaces folder
type ServiceMock struct {
	Mock *r.Mock
	storage.IService
}

var serviceMock = model.Service{
	ID:          mockProjectID,
	User:        mockUserID,
	Project:     mockProjectID,
	Image:       mockImageID,
	Name:        "mockedname",
	Description: "mockeddesc",
	Spec:        nil,
	Created:     time.Now(),
	Updated:     time.Now(),
}

func (s *ServiceMock) GetByNameOrID(user, nameOrID string) (*model.Service, error) {

	var (
		err            error
		service        = new(model.Service)
		service_filter = func(talk r.Term) r.Term {
			return r.And(
				talk.Field("user").Eq(user),
				r.Or(talk.Field("id").Eq(nameOrID), talk.Field("name").Eq(nameOrID)),
			)
		}
	)

	s.Mock.On(r.DB("test").Table(ServiceTable).Filter(service_filter)).Return(serviceMock, nil)

	res, err := r.DB("test").Table(ServiceTable).Filter(service_filter).Run(s.Mock)
	if err != nil {
		return nil, err
	}

	if res.IsNil() {
		return nil, nil
	}

	err = res.One(service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceMock) GetByID(user, id string) (*model.Service, error) {

	var (
		err            error
		service        = new(model.Service)
		service_filter = map[string]interface{}{
			"user": user,
			"id":   id,
		}
	)

	s.Mock.On(r.DB("test").Table(ServiceTable).Filter(service_filter)).Return(serviceMock, nil)

	res, err := r.DB("test").Table(ServiceTable).Filter(service_filter).Run(s.Mock)
	if err != nil {
		return nil, err
	}

	if res.IsNil() {
		return nil, nil
	}

	err = res.One(service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceMock) ListByProject(user, project string) (*model.ServiceList, error) {
	return &model.ServiceList{}, nil
}

// Insert new service into storage
func (s *ServiceMock) Insert(service *model.Service) (*model.Service, error) {

	var (
		err  error
		opts = r.InsertOpts{ReturnChanges: true}
	)

	s.Mock.On(r.DB("test").Table(ServiceTable).Insert(serviceMock, opts)).Return(nil, nil)

	err = r.DB("test").Table(ServiceTable).Insert(serviceMock, opts).Exec(s.Mock)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// Update service model
func (s *ServiceMock) Update(service *model.Service) (*model.Service, error) {

	var (
		err  error
		data = map[string]interface{}{
			"name": service.Name,
		}
		opts = r.ReplaceOpts{ReturnChanges: true}
	)

	s.Mock.On(r.DB("test").Table(ServiceTable).Replace(data, opts)).Return(nil, nil)

	err = r.DB("test").Table(ServiceTable).Replace(data, opts).Exec(s.Mock)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// Remove service model
func (s *ServiceMock) Remove(user, id string) error {

	var (
		err            error
		opts           = r.DeleteOpts{ReturnChanges: true}
		service_filter = map[string]interface{}{
			"user": user,
			"id":   id,
		}
	)

	s.Mock.On(r.DB("test").Table(ServiceTable).Filter(service_filter).Delete(opts)).Return(nil, nil)

	err = r.DB("test").Table(ServiceTable).Filter(service_filter).Delete(opts).Exec(s.Mock)
	if err != nil {
		return err
	}

	return nil
}

// Remove service model
func (s *ServiceMock) RemoveByProject(user, project string) error {

	var (
		err            error
		opts           = r.DeleteOpts{ReturnChanges: true}
		service_filter = map[string]interface{}{
			"user":    user,
			"project": project,
		}
	)

	s.Mock.On(r.DB("test").Table(ServiceTable).Filter(service_filter).Delete(opts)).Return(nil, nil)

	err = r.DB("test").Table(ServiceTable).Filter(service_filter).Delete(opts).Exec(s.Mock)
	if err != nil {
		return err
	}

	return nil
}

func newServiceMock(mock *r.Mock) *ServiceMock {
	s := new(ServiceMock)
	s.Mock = mock
	return s
}
