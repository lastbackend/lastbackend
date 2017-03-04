package storage

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
	"time"
)

const ServiceTable string = "services"

// Service Service type for interface in interfaces folder
type ServiceStorage struct {
	Session *r.Session
	storage.IService
}

func (s *ServiceStorage) CheckExistsByName(user, name string) (bool, error) {

	var (
		err            error
		service_filter = map[string]string{
			"name": name,
			"user": user,
		}
	)

	res, err := r.Table(ServiceTable).Filter(service_filter).Run(s.Session)

	if err != nil {
		return true, err
	}

	if !res.IsNil() {
		return true, nil
	}

	return !res.IsNil(), nil
}

func (s *ServiceStorage) GetByNameOrID(user, nameOrID string) (*model.Service, error) {

	var (
		err     error
		service = new(model.Service)
	)

	res, err := r.Table(ServiceTable).Filter(func(talk r.Term) r.Term {
		return r.And(
			talk.Field("user").Eq(user),
			r.Or(talk.Field("id").Eq(nameOrID), talk.Field("name").Eq(nameOrID)),
		)
	}).Run(s.Session)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(service)

	return service, nil
}

func (s *ServiceStorage) GetByName(user, name string) (*model.Service, error) {

	var (
		err            error
		service        = new(model.Service)
		project_filter = map[string]interface{}{
			"name": name,
			"user": user,
		}
	)

	res, err := r.Table(ServiceTable).Filter(project_filter).Run(s.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(service)

	return service, nil
}

func (s *ServiceStorage) GetByID(user, id string) (*model.Service, error) {

	var (
		err            error
		service        = new(model.Service)
		project_filter = map[string]interface{}{
			"id":   id,
			"user": user,
		}
	)

	res, err := r.Table(ServiceTable).Filter(project_filter).Run(s.Session)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(service)

	return service, nil
}

func (s *ServiceStorage) ListByProject(user, project string) (*model.ServiceList, error) {

	var (
		err            error
		projects       = new(model.ServiceList)
		project_filter = r.Row.Field("project").Eq(project)
	)

	res, err := r.Table(ServiceTable).Filter(project_filter).Run(s.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.All(projects)

	return projects, nil
}

// Insert new service into storage
func (s *ServiceStorage) Insert(service *model.Service) (*model.Service, error) {

	var (
		err  error
		opts = r.InsertOpts{ReturnChanges: true}
	)

	service.Created = time.Now()
	service.Updated = time.Now()

	res, err := r.Table(ServiceTable).Insert(service, opts).RunWrite(s.Session)
	if err != nil {
		return nil, err
	}

	service.ID = res.GeneratedKeys[0]

	return service, nil
}

// Update service model
func (s *ServiceStorage) Update(service *model.Service) (*model.Service, error) {

	service.Updated = time.Now()

	var (
		err  error
		opts = r.UpdateOpts{ReturnChanges: true}
		data = map[string]interface{}{
			"name":        service.Name,
			"description": service.Description,
			"updated":     service.Updated,
		}
	)

	_, err = r.Table(ServiceTable).Get(service.ID).Update(data, opts).RunWrite(s.Session)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// Remove service model
func (s *ServiceStorage) Remove(user, id string) error {

	var (
		err            error
		opts           = r.DeleteOpts{ReturnChanges: true}
		project_filter = map[string]interface{}{
			"user": user,
			"id":   id,
		}
	)

	_, err = r.Table(ServiceTable).Filter(project_filter).Delete(opts).RunWrite(s.Session)
	if err != nil {
		return err
	}

	return nil
}

// Remove service model
func (s *ServiceStorage) RemoveByProject(user, project string) error {

	var (
		err            error
		opts           = r.DeleteOpts{ReturnChanges: true}
		project_filter = map[string]interface{}{
			"user":    user,
			"project": project,
		}
	)

	_, err = r.Table(ServiceTable).Filter(project_filter).Delete(opts).RunWrite(s.Session)
	if err != nil {
		return err
	}

	return nil
}

func newServiceStorage(session *r.Session) *ServiceStorage {
	r.TableCreate(ServiceTable, r.TableCreateOpts{}).Run(session)
	s := new(ServiceStorage)
	s.Session = session
	return s
}
