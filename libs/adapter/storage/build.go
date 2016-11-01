package storage

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"

	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	r "gopkg.in/dancannon/gorethink.v2"
)

const BuildTable string = "builds"

// Build storage interface
type IBuild interface {
	// Get build and builds
	GetByID(string) (*model.Build, *error)
	GetByImage(string) (*model.BuildList, *error)
	// Insert and replace build
	Insert(*model.Build) (*model.Build, *error)
	Replace(*model.Build) (*model.Build, *error)
}

// Service Build type for interface in interfaces folder
type BuildStorage struct {
	IBuild
}

// Get build model by id
func (BuildStorage) GetByID(id string) (*model.Build, *e.Err) {

	var err error
	var build = new(model.Build)
	ctx := context.Get()

	res, err := r.Table(BuildTable).Get(id).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Build.NotFound(err)
	}
	res.One(build)

	defer res.Close()
	return build, nil
}

// Get builds by image
func (BuildStorage) GetByImage(id string) (*model.BuildList, *e.Err) {

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

// Insert new build into storage
func (BuildStorage) Insert(build *model.Build) (*model.Build, *e.Err) {
	ctx := context.Get()

	res, err := r.Table(BuildTable).Insert(build, r.InsertOpts{ReturnChanges: true}).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Build.Unknown(err)
	}
	res.One(build)

	defer res.Close()
	return build, nil
}

// Replace build model
func (BuildStorage) Replace(build *model.Build) (*model.Build, *e.Err) {
	ctx := context.Get()

	res, err := r.Table(BuildTable).Get(build.ID).Replace(build, r.ReplaceOpts{ReturnChanges: true}).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Build.Unknown(err)
	}
	res.One(build)

	defer res.Close()
	return build, nil
}
