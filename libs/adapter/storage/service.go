package storage

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
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

func (s *ServiceStorage) CheckExistsByName(user, project, name string) (bool, error) {

	var (
		err            error
		service_filter = map[string]string{
			"name":    name,
			"user":    user,
			"project": project,
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

func (s *ServiceStorage) GetByNameOrID(user, project, nameOrID string) (*model.Service, *e.Err) {

	var (
		err     error
		service = new(model.Service)
	)

	res, err := r.Table(ServiceTable).Filter(func(talk r.Term) r.Term {
		return r.And(
			talk.Field("user").Eq(user),
			r.Or(
				r.And(talk.Field("project").Eq(project), talk.Field("id").Eq(nameOrID)),
				r.And(talk.Field("project").Eq(project), talk.Field("name").Eq(nameOrID)),
			),
		)
	}).Run(s.Session)

	if err != nil {
		return nil, e.New("service").NotFound(err)
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(service)

	return service, nil
}

func (s *ServiceStorage) GetByName(user, project, name string) (*model.Service, *e.Err) {

	var (
		err            error
		service        = new(model.Service)
		project_filter = map[string]interface{}{
			"name":    name,
			"project": project,
			"user":    user,
		}
	)

	res, err := r.Table(ServiceTable).Filter(project_filter).Run(s.Session)

	if err != nil {
		return nil, e.New("service").NotFound(err)
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(service)

	return service, nil
}

func (s *ServiceStorage) GetByID(user, project, id string) (*model.Service, *e.Err) {

	var (
		err            error
		service        = new(model.Service)
		project_filter = map[string]interface{}{
			"id":      id,
			"project": project,
			"user":    user,
		}
	)

	res, err := r.Table(ServiceTable).Filter(project_filter).Run(s.Session)

	if err != nil {
		return nil, e.New("service").NotFound(err)
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(service)

	return service, nil
}

func (s *ServiceStorage) ListByProject(user, project string) (*model.ServiceList, *e.Err) {

	var (
		err            error
		projects       = new(model.ServiceList)
		project_filter = r.Row.Field("project").Eq(project)
	)

	res, err := r.Table(ServiceTable).Filter(project_filter).Run(s.Session)
	if err != nil {
		return nil, e.New("service").Unknown(err)
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.All(projects)

	return projects, nil
}

// Insert new service into storage
func (s *ServiceStorage) Insert(service *model.Service) (*model.Service, *e.Err) {

	var (
		err  error
		opts = r.InsertOpts{ReturnChanges: true}
	)

	service.Created = time.Now()
	service.Updated = time.Now()

	res, err := r.Table(ServiceTable).Insert(service, opts).RunWrite(s.Session)
	if err != nil {
		return nil, e.New("service").Unknown(err)
	}

	service.ID = res.GeneratedKeys[0]

	return service, nil
}

// Update build model
func (s *ServiceStorage) Update(service *model.Service) (*model.Service, *e.Err) {

	service.Updated = time.Now()

	var (
		err  error
		opts = r.UpdateOpts{ReturnChanges: true}
		data = map[string]interface{}{
			"description": service.Description,
			"updated":     service.Updated,
		}
	)

	_, err = r.Table(ServiceTable).Get(service.ID).Update(data, opts).RunWrite(s.Session)
	if err != nil {
		return nil, e.New("service").Unknown(err)
	}

	return service, nil
}

// Remove build model
func (s *ServiceStorage) Remove(user, project, id string) *e.Err {

	var (
		err            error
		opts           = r.DeleteOpts{ReturnChanges: true}
		project_filter = map[string]interface{}{
			"user":    user,
			"project": project,
			"id":      id,
		}
	)

	_, err = r.Table(ServiceTable).Filter(project_filter).Delete(opts).RunWrite(s.Session)
	if err != nil {
		return e.New("service").Unknown(err)
	}

	return nil
}

// Remove build model
func (s *ServiceStorage) RemoveByProject(user, project string) *e.Err {

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
		return e.New("service").Unknown(err)
	}

	return nil
}

func newServiceStorage(session *r.Session) *ServiceStorage {
	r.TableCreate(ServiceTable, r.TableCreateOpts{}).Run(session)
	s := new(ServiceStorage)
	s.Session = session
	return s
}
