package storage

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
	"time"
)

const UserTable = "users"

// Service User type for interface in interfaces folder
type UserStorage struct {
	Session *r.Session
	storage.IUser
}

func (s *UserStorage) GetByUsername(username string) (*model.User, error) {

	var err error
	var user = new(model.User)
	var username_filter = r.Row.Field("username").Eq(username)

	res, err := r.Table(UserTable).Filter(username_filter).Run(s.Session)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(user)

	return user, nil
}

func (s *UserStorage) GetByEmail(email string) (*model.User, error) {

	var err error
	var user = new(model.User)
	var email_filter = r.Row.Field("email").Eq(email)

	res, err := r.Table(UserTable).Filter(email_filter).Run(s.Session)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(user)

	return user, nil
}

func (s *UserStorage) GetByID(id string) (*model.User, error) {

	var err error
	var user = new(model.User)

	res, err := r.Table(UserTable).Get(id).Run(s.Session)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(user)

	return user, nil
}

func (s *UserStorage) Insert(user *model.User) (*model.User, error) {

	var err error
	var opts = r.InsertOpts{ReturnChanges: true}

	user.Created = time.Now()
	user.Updated = time.Now()

	res, err := r.Table(UserTable).Insert(user, opts).RunWrite(s.Session)

	if err != nil {
		return nil, err
	}

	user.ID = res.GeneratedKeys[0]

	return user, nil
}

func newUserStorage(session *r.Session) *UserStorage {
	r.TableCreate(UserTable, r.TableCreateOpts{}).Run(session)
	s := new(UserStorage)
	s.Session = session
	return s
}
