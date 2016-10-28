package etcd

import (
	"encoding/json"
	"github.com/coreos/etcd/client"
	"github.com/lastbackend/lastbackend/libs/adapter"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/utils"
	"golang.org/x/net/context"
)

// User Service type for interface in interfaces folder
type UserService struct{}

func (UserService) Insert(db adapter.IDatabase, username, email, gravatar string) (*string, *e.Err) {

	var err error
	var uuid = utils.GetUUIDV4()

	user := model.User{
		UUID:     uuid,
		Username: username,
		Email:    email,
		Gravatar: gravatar,
	}

	data, err := user.ToJson()
	if err != nil {
		return nil, e.User.Unknown(err)
	}

	_, err = client.NewKeysAPI(db).Set(context.Background(), username, data, nil)
	if err != nil {
		return nil, e.User.Unknown(err)
	}

	return &uuid, nil
}

func (UserService) Get(db adapter.IDatabase, username string) (*model.User, *e.Err) {
	var err error

	resp, err := client.NewKeysAPI(db).Get(context.Background(), username, nil)
	if err != nil {
		return nil, e.User.Unknown(err)
	}

	var user = new(model.User)
	err = json.Unmarshal([]byte(resp.Node.Value), user)
	if err != nil {
		return nil, e.User.Unknown(err)
	}

	return user, nil
}
