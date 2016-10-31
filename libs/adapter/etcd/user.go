package etcd

import (
	"github.com/lastbackend/lastbackend/libs/adapter"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/utils"
	"golang.org/x/net/context"
	"time"
)

// User Service type for interface in interfaces folder
type UserService struct{}

func (UserService) Insert(db adapter.IStorage, username, email, gravatar string) (*string, *e.Err) {

	var err error
	var uuid = utils.GetUUIDV4()
	var t = time.Now()

	user := model.User{
		UUID:     uuid,
		Username: username,
		Email:    email,
		Gravatar: gravatar,
		Updated:  t,
		Created:  t,
	}

	err = db.Create(context.Background(), "user", user, 0)
	if err != nil {
		return nil, e.User.Unknown(err)
	}

	return &uuid, nil
}

func (UserService) Get(db adapter.IStorage, username string) (*model.User, *e.Err) {
	var err error

	var user = new(model.User)
	err = db.Get(context.Background(), "", user)
	if err != nil {
		return nil, e.User.Unknown(err)
	}

	return user, nil
}
