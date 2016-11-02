package storage

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"

	"github.com/lastbackend/lastbackend/libs/interface/storage"
	r "gopkg.in/dancannon/gorethink.v2"
	"time"
)

const ProjectTable string = "projects"

// Project Service type for interface in interfaces folder
type ProjectStorage struct {
	Session *r.Session
	storage.IProject
}

func (s *ProjectStorage) GetByID(user, id string) (*model.Project, *e.Err) {

	var err error
	var project = new(model.Project)

	var user_filter = r.Row.Field("user").Eq(user)
	res, err := r.Table(ProjectTable).Get(id).Filter(user_filter).Run(s.Session)
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

	res, err := r.Table(ProjectTable).Filter(r.Row.Field("user").Eq(id)).Run(s.Session)
	if err != nil {
		return nil, e.Build.Unknown(err)
	}

	res.All(projects)

	defer res.Close()
	return projects, nil
}

// Insert new project into storage
func (s *ProjectStorage) Insert(project *model.Project) (*model.Project, *e.Err) {

	project.Created = time.Now()
	project.Updated = time.Now()
	res, err := r.Table(ProjectTable).Insert(project, r.InsertOpts{ReturnChanges: true}).RunWrite(s.Session)
	if err != nil {
		return nil, e.Project.Unknown(err)
	}
	project.ID = res.GeneratedKeys[0]
	return project, nil
}

// Replace build model
func (s *ProjectStorage) Replace(project *model.Project) (*model.Project, *e.Err) {
	project.Updated = time.Now()
	var user_filter = r.Row.Field("user").Eq(project.User)
	_, err := r.Table(ProjectTable).Get(project.ID).Filter(user_filter).Replace(project, r.ReplaceOpts{ReturnChanges: true}).RunWrite(s.Session)
	if err != nil {
		return nil, e.Build.Unknown(err)
	}
	return project, nil
}

func newProjectStorage(session *r.Session) *ProjectStorage {
	r.TableCreate(ProjectTable, r.TableCreateOpts{}).Run(session)
	s := new(ProjectStorage)
	s.Session = session
	return s
}
