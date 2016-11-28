package storage

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
	"time"

)

const ProjectTable string = "projects"

// Project Service type for interface in interfaces folder
type ProjectStorage struct {
	Session *r.Session
	storage.IProject
}

func (s *ProjectStorage) GetByName(user, name string) (*model.Project, *e.Err) {

	var err error
	var project = new(model.Project)
	var project_filter = map[string]string{
		"name": name,
		"user": user,
	}

	res, err := r.Table(ProjectTable).Filter(project_filter).Run(s.Session)

	if err != nil {
		return nil, e.Project.NotFound(err)
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(project)

	return project, nil
}



func (s *ProjectStorage) ExistByName(userID, name string) (bool, error) {
	var project_filter = map[string]string{
		"name": name,
		"user": userID,
	}
	res, err := r.Table(ProjectTable).Filter(project_filter).Run(s.Session)

	if err != nil {
		return true, err
	}

	if !res.IsNil() {
		return true, nil
	}
	return !res.IsNil(), nil
}

func (s *ProjectStorage) GetByID(user, id string) (*model.Project, *e.Err) {

	var err error
	var project = new(model.Project)
	var project_filter = map[string]interface{}{
		"id":   id,
		"user": user,
	}

	res, err := r.Table(ProjectTable).Filter(project_filter).Run(s.Session)

	if err != nil {
		return nil, e.Project.NotFound(err)
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(project)

	return project, nil
}

func (s *ProjectStorage) GetByUser(id string) (*model.ProjectList, *e.Err) {

	var err error
	var projects = new(model.ProjectList)
	var project_filter = r.Row.Field("user").Eq(id)

	res, err := r.Table(ProjectTable).Filter(project_filter).Run(s.Session)
	if err != nil {
		return nil, e.Project.Unknown(err)
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.All(projects)

	return projects, nil
}

// Insert new project into storage
func (s *ProjectStorage) Insert(project *model.Project) (*model.Project, *e.Err) {

	var err error
	var opts = r.InsertOpts{ReturnChanges: true}

	project.Created = time.Now()
	project.Updated = time.Now()

	res, err := r.Table(ProjectTable).Insert(project, opts).RunWrite(s.Session)
	if err != nil {
		return nil, e.Project.Unknown(err)
	}

	project.ID = res.GeneratedKeys[0]

	return project, nil
}

// Update build model
func (s *ProjectStorage) Update(project *model.Project) (*model.Project, *e.Err) {

	var err error
	var opts = r.UpdateOpts{ReturnChanges: true}

	project.Updated = time.Now()

	_, err = r.Table(ProjectTable).Get(project.ID).Update(map[string]string{
		"name": project.Name,
		"description": project.Description,
	}, opts).RunWrite(s.Session)

	if err != nil {
		return nil, e.Project.Unknown(err)
	}

	return project, nil
}

// Remove build model
func (s *ProjectStorage) Remove(id string) *e.Err {

	var err error
	var opts = r.DeleteOpts{ReturnChanges: true}

	_, err = r.Table(ProjectTable).Get(id).Delete(opts).RunWrite(s.Session)
	if err != nil {
		return e.Project.Unknown(err)
	}

	return nil
}

func newProjectStorage(session *r.Session) *ProjectStorage {
	r.TableCreate(ProjectTable, r.TableCreateOpts{}).Run(session)
	s := new(ProjectStorage)
	s.Session = session
	return s
}
