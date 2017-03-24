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
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	db "github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	st "github.com/lastbackend/lastbackend/pkg/storage/store"
	"golang.org/x/net/context"
	"time"
)

const UserTable = "users"

// Service User type for interface in interfaces folder
type UserStorage struct {
	storage.IUser
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *UserStorage) GetByUsername(username string) (*model.User, error) {

	var (
		keyInfo     = fmt.Sprintf("%s/%s/info", UserTable, username)
		keyProfile  = fmt.Sprintf("%s/%s/profile", UserTable, username)
		keyPassword = fmt.Sprintf("%s/%s/security/password", UserTable, username)
		keyEmails   = fmt.Sprintf("%s/%s/emails", UserTable, username)
		user        = new(model.User)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info := new(model.UserInfo)
	if err := client.Get(ctx, keyInfo, info); err != nil {
		if err.Error() == st.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	profile := new(model.UserProfile)
	if err := client.Get(ctx, keyProfile, profile); err != nil {
		if err.Error() == st.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	password := new(model.UserPassword)
	if err := client.Get(ctx, keyPassword, password); err != nil {
		if err.Error() == st.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	emails := new(model.UserEmails)
	if err := client.Get(ctx, keyEmails, emails); err != nil {
		if err.Error() == st.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	user.Username = username
	user.Profile = *profile
	user.Emails = *emails
	user.Gravatar = info.Gravatar
	user.Created = info.Created
	user.Updated = info.Updated
	user.Security.Pass.Salt = password.Salt
	user.Security.Pass.Password = password.Password

	return user, nil
}

func (s *UserStorage) GetByEmail(email string) (*model.User, error) {

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
		if err.Error() == st.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return s.GetByUsername(username)
}

func (s *UserStorage) SetInfo(username string, info *model.UserInfo) error {

	var key = fmt.Sprintf("%s/%s/info", UserTable, username)

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Create(ctx, key, info, nil, 0)
}

func (s *UserStorage) SetProfile(username string, profile *model.UserProfile) error {

	var key = fmt.Sprintf("%s/%s/profile", UserTable, username)

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Create(ctx, key, profile, nil, 0)
}

func (s *UserStorage) SetEmail(username, email string, asDefault bool) error {
	var (
		keyList   = fmt.Sprintf("%s/%s/emails", UserTable, username)
		keyCreate = fmt.Sprintf("%s/%s/emails/%s", UserTable, username, email)
		//keyUpdate = fmt.Sprintf("%s/%s/emails", UserTable, username)
		list = make(map[string]bool)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.List(ctx, keyList, ``, &list); err != nil {
		return err
	}

	if len(email) == 0 {
		return client.Create(ctx, keyCreate, asDefault, nil, 0)
	}

	if asDefault {
		for email := range list {
			list[email] = false
		}
	}

	list[email] = asDefault

	return nil // TODO: Need implement update method to storage -> client.Update(ctx, keyUpdate, list)
}

func (s *UserStorage) DeleteEmail(username, email string) error {
	return nil
}

func newUserStorage(config store.Config) *UserStorage {
	s := new(UserStorage)
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return db.Create(config)
	}
	return s
}
