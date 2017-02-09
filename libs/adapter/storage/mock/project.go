package mock

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
)

const ProjectTable string = "projects"

// Project Service type for interface in interfaces folder
type ProjectMock struct {
	Mock *r.Mock
	storage.IProject
}

func (s *ProjectMock) GetByNameOrID(user, nameOrID string) (*model.Project, error) {
	return nil, nil
}

func (s *ProjectMock) GetByID(user, id string) (*model.Project, error) {
	return nil, nil
}

func (s *ProjectMock) ListByUser(user string) (*model.ProjectList, error) {
	return nil, nil
}

// Insert new project into storage
func (s *ProjectMock) Insert(project *model.Project) (*model.Project, error) {
	return nil, nil
}

// Update project model
func (s *ProjectMock) Update(project *model.Project) (*model.Project, error) {
	return nil, nil
}

// Remove project model
func (s *ProjectMock) Remove(user, id string) error {
	return nil
}

func newProjectMock(mock *r.Mock) *ProjectMock {
	s := new(ProjectMock)
	s.Mock = mock
	return s
}
