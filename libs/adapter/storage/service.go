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

func (s *ServiceStorage) CheckExistsByName(userID, projectID, name string) (bool, error) {

	var err error
	var service_filter = map[string]string{
		"name":    name,
		"user":    userID,
		"project": projectID,
	}

	res, err := r.Table(ServiceTable).Filter(service_filter).Run(s.Session)

	if err != nil {
		return true, err
	}

	if !res.IsNil() {
		return true, nil
	}

	return !res.IsNil(), nil
}

func (s *ServiceStorage) GetByName(user, name string) (*model.Service, *e.Err) {

	var err error
	var service = new(model.Service)
	var project_filter = map[string]interface{}{
		"name": name,
		"user": user,
	}

	res, err := r.Table(ServiceTable).Filter(project_filter).Run(s.Session)

	if err != nil {
		return nil, e.Service.NotFound(err)
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(service)

	return service, nil
}

func (s *ServiceStorage) GetByID(user, id string) (*model.Service, *e.Err) {

	var err error
	var service = new(model.Service)
	var project_filter = map[string]interface{}{
		"id":   id,
		"user": user,
	}

	res, err := r.Table(ServiceTable).Filter(project_filter).Run(s.Session)

	if err != nil {
		return nil, e.Service.NotFound(err)
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(service)

	return service, nil
}

func (s *ServiceStorage) GetByUser(id string) (*model.ServiceList, *e.Err) {

	var err error
	var projects = new(model.ServiceList)
	var project_filter = r.Row.Field("user").Eq(id)

	res, err := r.Table(ServiceTable).Filter(project_filter).Run(s.Session)
	if err != nil {
		return nil, e.Service.Unknown(err)
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

	var err error
	var opts = r.InsertOpts{ReturnChanges: true}

	service.Created = time.Now()
	service.Updated = time.Now()

	res, err := r.Table(ServiceTable).Insert(service, opts).RunWrite(s.Session)
	if err != nil {
		return nil, e.Service.Unknown(err)
	}

	service.ID = res.GeneratedKeys[0]

	return service, nil
}

// Update build model
func (s *ServiceStorage) Update(service *model.Service) (*model.Service, *e.Err) {

	var err error
	var opts = r.UpdateOpts{ReturnChanges: true}
	var project_filter = map[string]interface{}{
		"name": service.Name,
	}

	service.Updated = time.Now()

	_, err = r.Table(ServiceTable).Update(project_filter, opts).RunWrite(s.Session)
	if err != nil {
		return nil, e.Service.Unknown(err)
	}

	return service, nil
}

// Remove build model
func (s *ServiceStorage) Remove(id string) *e.Err {

	var err error
	var opts = r.DeleteOpts{ReturnChanges: true}

	_, err = r.Table(ServiceTable).Get(id).Delete(opts).RunWrite(s.Session)
	if err != nil {
		return e.Service.Unknown(err)
	}

	return nil
}

// Remove build model
func (s *ServiceStorage) RemoveByProject(id string) *e.Err {

	var err error
	var opts = r.DeleteOpts{ReturnChanges: true}
	var project_filter = map[string]interface{}{
		"project": id,
	}

	_, err = r.Table(ServiceTable).Filter(project_filter).Delete(opts).RunWrite(s.Session)
	if err != nil {
		return e.Service.Unknown(err)
	}

	return nil
}

func newServiceStorage(session *r.Session) *ServiceStorage {
	r.TableCreate(ServiceTable, r.TableCreateOpts{}).Run(session)
	s := new(ServiceStorage)
	s.Session = session
	return s
}
