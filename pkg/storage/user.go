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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"golang.org/x/net/context"
	"time"
)

const UserTable = "users"

// Service User type for interface in interfaces folder
type UserStorage struct {
	IUser
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *UserStorage) GetByUsername(username string) (*types.User, error) {

	var (
		keyInfo     = fmt.Sprintf("%s/%s/info", UserTable, username)
		keyProfile  = fmt.Sprintf("%s/%s/profile", UserTable, username)
		keyPassword = fmt.Sprintf("%s/%s/security/password", UserTable, username)
		keyEmails   = fmt.Sprintf("%s/%s/emails", UserTable, username)
		keyVendors  = fmt.Sprintf("%s/%s/vendors", UserTable, username)
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
	if err := client.Get(ctx, keyInfo, info); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	profile := new(types.UserProfile)
	if err := client.Get(ctx, keyProfile, profile); err != nil && err.Error() != store.ErrKeyNotFound {
		return nil, err
	}

	password := new(types.UserPassword)
	if err := client.Get(ctx, keyPassword, password); err != nil && err.Error() != store.ErrKeyNotFound {
		return nil, err
	}

	emails := new(types.UserEmails)
	if err := client.Get(ctx, keyEmails, emails); err != nil && err.Error() != store.ErrKeyNotFound {
		return nil, err
	}

	vendors := new(types.UserVendors)
	if err := client.Get(ctx, keyVendors, vendors); err != nil && err.Error() != store.ErrKeyNotFound {
		return nil, err
	}

	user.Username = username
	user.Profile = *profile
	user.Emails = *emails
	user.Vendors = *vendors
	user.Gravatar = info.Gravatar
	user.Created = info.Created
	user.Updated = info.Updated
	user.Security.Pass.Salt = password.Salt
	user.Security.Pass.Password = password.Password

	return user, nil
}

func (s *UserStorage) GetByEmail(email string) (*types.User, error) {

	var (
		key      = fmt.Sprintf("helper/emails/%s", email)
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

	return s.GetByUsername(username)
}

func newUserStorage(config store.Config) *UserStorage {
	s := new(UserStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
