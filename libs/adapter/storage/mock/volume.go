package mock

import (
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

func (s *VolumeMock) GetByID(user, id string) (*model.Volume, error) {
	return nil, nil
}

func (s *VolumeMock) ListByProject(id string) (*model.VolumeList, error) {
	return nil, nil
}

// Insert new volume into storage
func (s *VolumeMock) Insert(project *model.Volume) (*model.Volume, error) {
	return nil, nil
}

// Remove volume model
func (s *VolumeMock) Remove(id string) error {
	return nil
}

func newVolumeMock(mock *r.Mock) *VolumeMock {
	s := new(VolumeMock)
	s.Mock = mock
	return s
}
