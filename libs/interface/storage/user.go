package storage

import (
	"github.com/lastbackend/lastbackend/libs/adapter"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
)

type IUserService interface {
	Insert(db adapter.IStorage, username, email, gravatar string) (*string, *e.Err)
	Get(db adapter.IStorage, username string) (*model.User, *e.Err)
}
