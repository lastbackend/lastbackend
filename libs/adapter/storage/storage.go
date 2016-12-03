package storage

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/pkg/daemon/config"
	r "gopkg.in/dancannon/gorethink.v2"
)

type Storage struct {
	*UserStorage
	*ProjectStorage
	*ServiceStorage
	*ImageStorage
	*BuildStorage
	*HookStorage
	*VolumeStorage
}

func (s *Storage) User() storage.IUser {
	if s == nil {
		return nil
	}
	return s.UserStorage
}

func (s *Storage) Project() storage.IProject {
	if s == nil {
		return nil
	}
	return s.ProjectStorage
}

func (s *Storage) Service() storage.IService {
	if s == nil {
		return nil
	}
	return s.ServiceStorage
}

func (s *Storage) Image() storage.IImage {
	if s == nil {
		return nil
	}
	return s.ImageStorage
}

func (s *Storage) Build() storage.IBuild {
	if s == nil {
		return nil
	}
	return s.BuildStorage
}

func (s *Storage) Hook() storage.IHook {
	if s == nil {
		return nil
	}
	return s.HookStorage
}

func (s *Storage) Volume() storage.IVolume {
	if s == nil {
		return nil
	}
	return s.VolumeStorage
}

func Get() (*Storage, error) {

	store := new(Storage)

	session, err := r.Connect(config.GetRethinkDB())
	if err != nil {
		return nil, err
	}

	r.DBCreate(config.Get().RethinkDB.Database).Run(session)

	store.UserStorage = newUserStorage(session)
	store.ProjectStorage = newProjectStorage(session)
	store.ServiceStorage = newServiceStorage(session)
	store.ImageStorage = newImageStorage(session)
	store.BuildStorage = newBuildStorage(session)
	store.HookStorage = newHookStorage(session)
	store.VolumeStorage = newVolumeStorage(session)

	return store, nil
}
