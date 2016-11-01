package rethinkdb

import (
	"github.com/lastbackend/lastbackend/libs/adapter"
	"github.com/lastbackend/lastbackend/libs/model"
)

type Interface {
	GetByID(storage.Session, string) (*model.Build, *error)
}

// Service Service type for interface in interfaces folder
type BuildService struct{}

func (BuildService) GetByID(storage adapter.IStorage, id string) (*model.Build, *e.Err) {

	var err error

	return build, nil
}