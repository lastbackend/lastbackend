package mock

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
)

const BuildTable string = "builds"

// Service Build type for interface in interfaces folder
type BuildMock struct {
	Mock *r.Mock
	storage.IBuild
}

// Get build model by id
func (s *BuildMock) GetByID(user, id string) (*model.Build, *e.Err) {
	return nil, nil
}

// Get builds by image
func (s *BuildMock) GetByImage(user, id string) (*model.BuildList, *e.Err) {
	return nil, nil
}

// Insert new build into storage
func (s *BuildMock) Insert(build *model.Build) (*model.Build, *e.Err) {
	return nil, nil
}

// Replace build model
func (s *BuildMock) Replace(build *model.Build) (*model.Build, *e.Err) {
	return nil, nil
}

func newBuildMock(mock *r.Mock) *BuildMock {
	s := new(BuildMock)
	s.Mock = mock
	return s
}
