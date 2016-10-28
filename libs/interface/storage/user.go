package storage

import (
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/libs/adapter"
	e "github.com/lastbackend/lastbackend/libs/errors"
)

type IUserService interface {
	Insert(db adapter.IDatabase, username, email, gravatar string) (*string, *e.Err)
	Get(db adapter.IDatabase, username string) (*model.User, *e.Err)
}
