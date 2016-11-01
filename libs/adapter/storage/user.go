package storage

import (
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
)

const (
	userTable = "users"
)

type IUser interface {
	GetByID(string) (*model.User, *error)
	Insert(model.User) (*model.User, *e.Err)
}

// Service User type for interface in interfaces folder
type User struct{
  IUser
}

func (User) GetByID(uuid string) (*model.User, *e.Err) {

	var err error
	var user = new(model.User)
	ctx := context.Get()

	res, err := r.Table(userTable).Get(uuid).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.User.NotFound(err)
	}
	res.One(user)

	defer res.Close()
	return user, nil
}

func (User) Insert(u model.User) (*model.User, *e.Err) {

	var err error
	var user = new(model.User)
	ctx := context.Get()

	res, err := r.Table(userTable).Insert(user).Run(ctx.Storage.Session)
	if err != nil {
		return nil, e.User.NotFound(err)
	}
	res.One(user)

	defer res.Close()
	return user, nil
}
