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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"golang.org/x/net/context"
	"time"
)

const UserTable = "users"

// Service User type for interface in interfaces folder
type UserStorage struct {
	IUser
	Helper IHelper
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *UserStorage) GetByUsername(ctx context.Context, username string) (*types.User, error) {

	var (
		keyMeta     = s.Helper.KeyDecorator(ctx,"meta")
		keyProfile  = s.Helper.KeyDecorator(ctx,"profile")
		keyPassword = s.Helper.KeyDecorator(ctx,"security", "password")
		keyEmails   = s.Helper.KeyDecorator(ctx,"emails")
		keyVendors  = s.Helper.KeyDecorator(ctx,"vendors")
		user        = new(types.User)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info := new(types.UserInfo)
	if err := client.Get(ctx, keyMeta, info); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	if err := client.Get(ctx, keyProfile, &user.Profile); err != nil && err.Error() != store.ErrKeyNotFound {
		return nil, err
	}

	password := new(types.UserPassword)
	if err := client.Get(ctx, keyPassword, password); err != nil && err.Error() != store.ErrKeyNotFound {
		return nil, err
	}

	user.Emails = make(map[string]bool)
	if err := client.Map(ctx, keyEmails, ``, user.Emails); err != nil && err.Error() != store.ErrKeyNotFound {
		return nil, err
	}

	user.Vendors = make(map[string]*types.Vendor)
	if err := client.Map(ctx, keyVendors, ``, user.Vendors); err != nil && err.Error() != store.ErrKeyNotFound {
		return nil, err
	}

	user.Username = username
	user.Gravatar = info.Gravatar
	user.Created = info.Created
	user.Updated = info.Updated
	user.Security.Pass.Salt = password.Salt
	user.Security.Pass.Password = password.Password

	return user, nil
}

func (s *UserStorage) GetByEmail(ctx context.Context, email string) (*types.User, error) {

	var (
		key      = s.Helper.KeyDecorator(ctx,"helper", "emails", email)
		username string
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Get(ctx, key, &username); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return s.GetByUsername(ctx, username)
}

func NewUserStorage(config store.Config, helper IHelper) *UserStorage {
	s := new(UserStorage)
	s.Helper = helper
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
