package storage

import (
	"github.com/lastbackend/lastbackend/libs/adapter"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
)

type IAccountService interface {
	Insert(db adapter.IDatabase, username, userID, password, salt string) (*string, *e.Err)
	Get(db adapter.IDatabase, username string) (*model.Account, *e.Err)
}
