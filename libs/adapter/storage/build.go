package storage

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"

	"github.com/lastbackend/lastbackend/libs/interface/storage"
	r "gopkg.in/dancannon/gorethink.v2"
)

const BuildTable string = "builds"

// Service Build type for interface in interfaces folder
type BuildStorage struct {
	Session *r.Session
	storage.IBuild
}

// Get build model by id
func (s *BuildStorage) GetByID(id string) (*model.Build, *e.Err) {

	var err error
	var build = new(model.Build)

	res, err := r.Table(BuildTable).Get(id).Run(s.Session)
	if err != nil {
		return nil, e.Build.NotFound(err)
	}
	res.One(build)

	defer res.Close()
	return build, nil
}

// Get builds by image
func (s *BuildStorage) GetByImage(id string) (*model.BuildList, *e.Err) {

	var err error
	var builds = new(model.BuildList)

	res, err := r.Table(BuildTable).Get(id).Run(s.Session)
	if err != nil {
		return nil, e.Build.Unknown(err)
	}

	res.All(builds)

	defer res.Close()
	return builds, nil
}

// Insert new build into storage
func (s *BuildStorage) Insert(build *model.Build) (*model.Build, *e.Err) {

	res, err := r.Table(BuildTable).Insert(build, r.InsertOpts{ReturnChanges: true}).Run(s.Session)
	if err != nil {
		return nil, e.Build.Unknown(err)
	}
	res.One(build)

	defer res.Close()
	return build, nil
}

// Replace build model
func (s *BuildStorage) Replace(build *model.Build) (*model.Build, *e.Err) {

	res, err := r.Table(BuildTable).Get(build.ID).Replace(build, r.ReplaceOpts{ReturnChanges: true}).Run(s.Session)
	if err != nil {
		return nil, e.Build.Unknown(err)
	}
	res.One(build)

	defer res.Close()
	return build, nil
}

func newBuildStorage(session *r.Session) *BuildStorage {
	s := new(BuildStorage)
	s.Session = session
	return s
}
