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
func (s *BuildStorage) GetByID(user, id string) (*model.Build, *e.Err) {

	var err error
	var build = new(model.Build)
	var user_filter = r.Row.Field("user").Eq(user)
	res, err := r.Table(BuildTable).Get(id).Filter(user_filter).Run(s.Session)
	if err != nil {
		return nil, e.New("build").NotFound(err)
	}
	res.One(build)

	defer res.Close()
	return build, nil
}

// Get builds by image
func (s *BuildStorage) ListByImage(user, id string) (*model.BuildList, *e.Err) {

	var err error
	var builds = new(model.BuildList)
	var image_filter = r.Row.Field("image").Field("id").Eq(id)
	var user_filter = r.Row.Field("user").Eq(user)
	res, err := r.Table(BuildTable).Filter(image_filter).Filter(user_filter).Run(s.Session)
	if err != nil {
		return nil, e.New("build").Unknown(err)
	}

	res.All(builds)

	defer res.Close()
	return builds, nil
}

// Insert new build into storage
func (s *BuildStorage) Insert(build *model.Build) (*model.Build, *e.Err) {

	res, err := r.Table(BuildTable).Insert(build, r.InsertOpts{ReturnChanges: true}).RunWrite(s.Session)
	if err != nil {
		return nil, e.New("build").Unknown(err)
	}

	build.ID = res.GeneratedKeys[0]

	return build, nil
}

// Replace build model
func (s *BuildStorage) Replace(build *model.Build) (*model.Build, *e.Err) {
	var user_filter = r.Row.Field("user").Eq(build.User)
	_, err := r.Table(BuildTable).Get(build.ID).Filter(user_filter).Replace(build, r.ReplaceOpts{ReturnChanges: true}).RunWrite(s.Session)
	if err != nil {
		return nil, e.New("build").Unknown(err)
	}

	return build, nil
}

func newBuildStorage(session *r.Session) *BuildStorage {
	r.TableCreate(BuildTable, r.TableCreateOpts{}).Run(session)
	s := new(BuildStorage)
	s.Session = session
	return s
}
