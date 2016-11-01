package storage

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"

	"github.com/lastbackend/lastbackend/libs/interface/storage"
	r "gopkg.in/dancannon/gorethink.v2"
)

const ProjectTable string = "projects"

// Project Service type for interface in interfaces folder
type ProjectStorage struct {
	Session *r.Session
	storage.IProject
}

func (s *ProjectStorage) GetByID(uuid string) (*model.Project, *e.Err) {

	var err error
	var project = new(model.Project)

	res, err := r.Table(ProjectTable).Get(uuid).Run(s.Session)
	if err != nil {
		return nil, e.Project.NotFound(err)
	}
	res.One(project)

	defer res.Close()
	return project, nil
}

func (s *ProjectStorage) GetByUser(id string) (*model.ProjectList, *e.Err) {

	var err error
	var projects = new(model.ProjectList)

	res, err := r.Table(BuildTable).Get(id).Run(s.Session)
	if err != nil {
		return nil, e.Build.Unknown(err)
	}

	res.All(projects)

	defer res.Close()
	return projects, nil
}

// Insert new project into storage
func (s *ProjectStorage) Insert(project *model.Project) (*model.Project, *e.Err) {

	res, err := r.Table(ProjectTable).Insert(project, r.InsertOpts{ReturnChanges: true}).Run(s.Session)
	if err != nil {
		return nil, e.Project.Unknown(err)
	}
	res.One(project)

	defer res.Close()
	return project, nil
}

// Replace build model
func (s *ProjectStorage) Replace(project *model.Project) (*model.Project, *e.Err) {

	res, err := r.Table(ProjectTable).Get(project.ID).Replace(project, r.ReplaceOpts{ReturnChanges: true}).Run(s.Session)
	if err != nil {
		return nil, e.Build.Unknown(err)
	}
	res.One(project)

	defer res.Close()
	return project, nil
}

func newProjectStorage(session *r.Session) *ProjectStorage {
	s := new(ProjectStorage)
	s.Session = session
	return s
}
