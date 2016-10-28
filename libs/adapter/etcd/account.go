package etcd

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/lastbackend/lastbackend/libs/adapter"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/utils"
	"golang.org/x/net/context"
	"time"
)

// Account Service type for interface in interfaces folder
type AccountService struct{}

func (AccountService) Insert(db adapter.IDatabase, username, userID, password, salt string) (*string, *e.Err) {

	var err error
	var uuid = utils.GetUUIDV4()
	var t = time.Now()

	user := model.Account{
		UUID:     uuid,
		UserID:   userID,
		Password: password,
		Salt:     salt,
		Updated:  t,
		Created:  t,
	}

	data, err := user.ToJson()
	if err != nil {
		return nil, e.Account.Unknown(err)
	}

	_, err = client.NewKeysAPI(db).Set(context.Background(), fmt.Sprintf("account:%s", username), string(data), nil)
	if err != nil {
		return nil, e.Account.Unknown(err)
	}

	return &uuid, nil
}

func (AccountService) Get(db adapter.IDatabase, username string) (*model.Account, *e.Err) {
	var err error

	resp, err := client.NewKeysAPI(db).Get(context.Background(), fmt.Sprintf("account:%s", username), nil)
	if err != nil {
		return nil, e.Account.Unknown(err)
	}

	var account = new(model.Account)
	err = json.Unmarshal([]byte(resp.Node.Value), account)
	if err != nil {
		return nil, e.Account.Unknown(err)
	}

	return account, nil
}
