package storage

import (
	"github.com/lastbackend/lastbackend/libs/adapter"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
)

type IAccountService interface {
	Insert(db adapter.IStorage, username, userID, password, salt string) (*string, *e.Err)
	Get(db adapter.IStorage, username string) (*model.Account, *e.Err)
}
