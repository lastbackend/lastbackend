package storage

import (
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
)

const (
	accountTable = "accounts"
)

type IAccount interface {
	GetByID(string) (*model.Account, *error)
	Insert(model.Account) (*model.Account, *e.Err)
}

// Service Account type for interface in interfaces folder
type Account struct{
  IAccount
}

func (Account) GetByID(uuid string) (*model.Account, *e.Err) {

	var err error
	var account = new(model.Account)
	ctx := context.Get()

	res, err := r.Table(accountTable).Get(uuid).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Account.NotFound(err)
	}
	res.One(account)

	defer res.Close()
	return account, nil
}

func (Account) Insert(u model.Account) (*model.Account, *e.Err) {

	var err error
	var account = new(model.Account)
	ctx := context.Get()

	res, err := r.Table(accountTable).Insert(account).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.Account.NotFound(err)
	}
	res.One(account)

	defer res.Close()
	return account, nil
}
