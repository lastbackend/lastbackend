package mock

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
	"time"
)

const mockProjectID string = "mocked"
const ProjectTable string = "projects"

// Project Service type for interface in interfaces folder
type ProjectMock struct {
	Mock *r.Mock
	storage.IProject
}

var projectMock = model.Project{
	ID:          mockProjectID,
	User:        mockUserID,
	Name:        "mockedname",
	Description: "mockeddesc",
	Created:     time.Now(),
	Updated:     time.Now(),
}

func (s *ProjectMock) GetByNameOrID(user, nameOrID string) (*model.Project, error) {

	var (
		err            error
		project        = new(model.Project)
		project_filter = func(talk r.Term) r.Term {
			return r.And(
				talk.Field("user").Eq(user),
				r.Or(
					talk.Field("name").Eq(nameOrID),
					talk.Field("id").Eq(nameOrID)),
			)
		}
	)

	s.Mock.On(r.DB("test").Table(userTable).Filter(project_filter)).Return(projectMock, nil)

	res, err := r.DB("test").Table(userTable).Filter(project_filter).Run(s.Mock)
	if err != nil {
		return nil, err
	}

	if res.IsNil() {
		return nil, nil
	}

	err = res.One(project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectMock) GetByID(user, id string) (*model.Project, error) {

	var (
		err            error
		project        = new(model.Project)
		project_filter = map[string]interface{}{
			"user":    project,
			"project": id,
		}
	)

	s.Mock.On(r.DB("test").Table(userTable).Filter(project_filter)).Return(projectMock, nil)

	res, err := r.DB("test").Table(userTable).Filter(project_filter).Run(s.Mock)
	if err != nil {
		return nil, err
	}

	if res.IsNil() {
		return nil, nil
	}

	err = res.One(project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectMock) ListByUser(user string) (*model.ProjectList, error) {

	var (
		err             error
		projectListMock = new(model.ProjectList)
		project_filter  = r.Row.Field("user").Eq(user)
	)

	s.Mock.On(r.DB("test").Table(userTable).Filter(project_filter)).Return(projectMock, nil)

	res, err := r.DB("test").Table(userTable).Filter(project_filter).Run(s.Mock)
	if err != nil {
		return nil, err
	}

	if res.IsNil() {
		return nil, nil
	}

	err = res.All(projectListMock)
	if err != nil {
		return nil, err
	}

	return projectListMock, nil
}

// Insert new project into storage
func (s *ProjectMock) Insert(project *model.Project) (*model.Project, error) {

	var (
		err  error
		opts = r.InsertOpts{ReturnChanges: true}
	)

	s.Mock.On(r.DB("test").Table(ProjectTable).Insert(projectMock, opts)).Return(nil, nil)

	err = r.DB("test").Table(ProjectTable).Insert(projectMock, opts).Exec(s.Mock)
	if err != nil {
		return nil, err
	}

	return project, nil
}

// Update project model
func (s *ProjectMock) Update(project *model.Project) (*model.Project, error) {

	var (
		err  error
		data = map[string]interface{}{
			"name": project.Name,
		}
		opts = r.ReplaceOpts{ReturnChanges: true}
	)

	s.Mock.On(r.DB("test").Table(ProjectTable).Replace(data, opts)).Return(nil, nil)

	err = r.DB("test").Table(ProjectTable).Replace(data, opts).Exec(s.Mock)
	if err != nil {
		return nil, err
	}

	return project, nil
}

// Remove project model
func (s *ProjectMock) Remove(user, id string) error {

	var (
		err            error
		opts           = r.DeleteOpts{ReturnChanges: true}
		project_filter = map[string]interface{}{
			"user": user,
			"id":   id,
		}
	)

	s.Mock.On(r.DB("test").Table(ProjectTable).Filter(project_filter).Delete(opts)).Return(nil, nil)

	err = r.DB("test").Table(ProjectTable).Filter(project_filter).Delete(opts).Exec(s.Mock)
	if err != nil {
		return err
	}

	return nil
}

func newProjectMock(mock *r.Mock) *ProjectMock {
	s := new(ProjectMock)
	s.Mock = mock
	return s
}
