package storage

import (
	"github.com/lastbackend/lastbackend/cmd/daemon/config"
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	r "gopkg.in/dancannon/gorethink.v2"
)

type Storage struct {
	*UserStorage
	*ProjectStorage
	*ImageStorage
	*BuildStorage
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

func Get() (*Storage, error) {

	store := new(Storage)

	session, err := r.Connect(config.GetRethinkDB())
	if err != nil {
		panic(err.Error())
	}

	store.UserStorage = newUserStorage(session)
	store.ProjectStorage = newProjectStorage(session)
	store.ImageStorage = newImageStorage(session)
	store.BuildStorage = newBuildStorage(session)

	return store, nil
}
