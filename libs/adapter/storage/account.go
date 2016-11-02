package storage

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
)

const AccountTable = "accounts"

// Service Account type for interface in interfaces folder
type AccountStorage struct {
	Session *r.Session
	storage.IAccount
}

func (s *AccountStorage) GetByID(uuid string) (*model.Account, *e.Err) {

	var err error
	var account = new(model.Account)

	res, err := r.Table(AccountTable).Get(uuid).Run(s.Session)
	if err != nil {
		return nil, e.Account.NotFound(err)
	}
	res.One(account)

	defer res.Close()
	return account, nil
}

func (s *AccountStorage) Insert(account *model.Account) (*model.Account, *e.Err) {

	var err error

	res, err := r.Table(AccountTable).Insert(account).Run(s.Session)
	if err != nil {
		return nil, e.Account.NotFound(err)
	}
	res.One(account)

	defer res.Close()
	return account, nil
}

func newAccountStorage(session *r.Session) *AccountStorage {
	s := new(AccountStorage)
	s.Session = session
	return s
}
