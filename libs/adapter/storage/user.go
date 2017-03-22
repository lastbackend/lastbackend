//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package storage

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	db "github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"golang.org/x/net/context"
	"time"
)

const UserTable = "users"

// Service User type for interface in interfaces folder
type UserStorage struct {
	storage.IUser
	Client func() (store.Interface, store.DestroyFunc, error)
}

func (s *UserStorage) GetByUsername(username string) (*model.User, error) {
	return nil, nil
}

func (s *UserStorage) GetByEmail(email string) (*model.User, error) {
	var user = new(model.User)
	return user, nil
}

func (s *UserStorage) GetByID(id string) (*model.User, error) {

	var user = new(model.User)

	client, close, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Get(ctx, UserTable+"/"+id, user)
	cancel()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStorage) GetByUsernameOrEmail(usernameOrEmail string) (*model.User, error) {
	var user = new(model.User)
	return user, nil
}

func (s *UserStorage) Insert(user *model.User) (*model.User, error) {
	return user, nil
}

func newUserStorage(config store.Config) *UserStorage {
	s := new(UserStorage)
	s.Client = func() (store.Interface, store.DestroyFunc, error) {
		return db.Create(config)
	}
	return s
}
