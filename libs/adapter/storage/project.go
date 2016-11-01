package storage

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"

	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	r "gopkg.in/dancannon/gorethink.v2"
)

const ProjectTable string = "projects"

type IProject interface {
	GetByID(string) (*model.Project, *error)
	GetByUser(string) (*model.ProjectList, *error)
	Insert(*model.Project) (*model.Project, *error)
	Replace(*model.Project) (*model.Project, *error)
}

// Project Service type for interface in interfaces folder
type ProjectStorage struct {
	IProject
}

func (ProjectStorage) GetByID(uuid string) (*model.Build, *e.Err) {

	var err error
	var build = new(model.Build)
	ctx := context.Get()

	res, err := r.Table(ProjectTable).Get(uuid).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Project.NotFound(err)
	}
	res.One(build)

	defer res.Close()
	return build, nil
}

func (ProjectStorage) GetByUser(id string) (*model.BuildList, *e.Err) {

	var err error
	var builds = new(model.BuildList)
	ctx := context.Get()

	res, err := r.Table(BuildTable).Get(id).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Build.Unknown(err)
	}

	res.All(builds)

	defer res.Close()
	return builds, nil
}

// Insert new project into storage
func (ProjectStorage) Insert(project *model.Project) (*model.Project, *e.Err) {
	ctx := context.Get()

	res, err := r.Table(ProjectTable).Insert(project, r.InsertOpts{ReturnChanges: true}).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Project.Unknown(err)
	}
	res.One(project)

	defer res.Close()
	return project, nil
}

// Replace build model
func (ProjectStorage) Replace(project *model.Project) (*model.Project, *e.Err) {
	ctx := context.Get()

	res, err := r.Table(ProjectTable).Get(project.ID).Replace(project, r.ReplaceOpts{ReturnChanges: true}).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Build.Unknown(err)
	}
	res.One(project)

	defer res.Close()
	return project, nil
}
