package storage

import (
	"github.com/lastbackend/lastbackend/libs/adapter"
	e "github.com/lastbackend/lastbackend/libs/errors"
)

type IUserService interface {
	Insert(db adapter.IDatabase, username, email, gravatar, password, salt string) (*string, *e.Err)
}
