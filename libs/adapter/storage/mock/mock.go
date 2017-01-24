package mock

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	r "gopkg.in/dancannon/gorethink.v2"
)

type Mock struct {
	*UserMock
	*ProjectMock
	*ServiceMock
	*ImageMock
	*BuildMock
	*HookMock
	*VolumeMock
	*ActivityMock
}

func (s *Mock) User() storage.IUser {
	if s == nil {
		return nil
	}
	return s.UserMock
}

func (s *Mock) Project() storage.IProject {
	if s == nil {
		return nil
	}
	return s.ProjectMock
}

func (s *Mock) Service() storage.IService {
	if s == nil {
		return nil
	}
	return s.ServiceMock
}

func (s *Mock) Image() storage.IImage {
	if s == nil {
		return nil
	}
	return s.ImageMock
}

func (s *Mock) Build() storage.IBuild {
	if s == nil {
		return nil
	}
	return s.BuildMock
}

func (s *Mock) Hook() storage.IHook {
	if s == nil {
		return nil
	}
	return s.HookMock
}

func (s *Mock) Volume() storage.IVolume {
	if s == nil {
		return nil
	}
	return s.VolumeMock
}

func (s *Mock) Activity() storage.IActivity {
	if s == nil {
		return nil
	}
	return s.ActivityMock
}

func Get() (*Mock, error) {

	var store = new(Mock)
	var mock = r.NewMock()

	store.UserMock = newUserMock(mock)
	store.ProjectMock = newProjectMock(mock)
	store.ServiceMock = newServiceMock(mock)
	store.ImageMock = newImageMock(mock)
	store.BuildMock = newBuildMock(mock)
	store.HookMock = newHookMock(mock)
	store.VolumeMock = newVolumeMock(mock)
	store.ActivityMock = newActivityMock(mock)

	return store, nil
}
