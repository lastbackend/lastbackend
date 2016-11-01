package storage

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
)

const UserTable = "users"

// Service User type for interface in interfaces folder
type UserStorage struct {
	Session *r.Session
	storage.IUser
}

func (s *UserStorage) GetByID(uuid string) (*model.User, *e.Err) {

	var err error
	var user = new(model.User)

	res, err := r.Table(UserTable).Get(uuid).Run(s.Session)
	if err != nil {
		return nil, e.User.NotFound(err)
	}
	res.One(user)

	defer res.Close()
	return user, nil
}

func (s *UserStorage) Insert(user *model.User) (*model.User, *e.Err) {

	var err error

	res, err := r.Table(UserTable).Insert(user).Run(s.Session)
	if err != nil {
		return nil, e.User.NotFound(err)
	}
	res.One(user)

	defer res.Close()
	return user, nil
}

func newUserStorage(session *r.Session) *UserStorage {
	s := new(UserStorage)
	s.Session = session
	return s
}
