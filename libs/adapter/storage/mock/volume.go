package mock

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
)

const VolumeTable string = "volumes"

// Volume Service type for interface in interfaces folder
type VolumeMock struct {
	Mock *r.Mock
	storage.IVolume
}

func (s *VolumeMock) GetByID(user, id string) (*model.Volume, *e.Err) {
	return nil, nil
}

func (s *VolumeMock) ListByProject(id string) (*model.VolumeList, *e.Err) {
	return nil, nil
}

// Insert new volume into storage
func (s *VolumeMock) Insert(project *model.Volume) (*model.Volume, *e.Err) {
	return nil, nil
}

// Remove volume model
func (s *VolumeMock) Remove(id string) *e.Err {
	return nil
}

func newVolumeMock(mock *r.Mock) *VolumeMock {
	s := new(VolumeMock)
	s.Mock = mock
	return s
}
