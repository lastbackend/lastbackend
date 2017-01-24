package mock

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
)

const mockUserDB = "test"
const mockUserID = "mocked"
const userTable = "users"

// Service User type for interface in interfaces folder
type UserMock struct {
	Mock *r.Mock
	storage.IUser
}

var userMock = &model.User{
	ID:       mockUserID,
	Username: "mocked",
	Email:    "mocked@mocked.com",
	Gravatar: "a931b3bb185354ecfe43d736b7ad51cb",
	Password: "$2a$10$KBFpkZS0DGBaEOepyGopPur7jsGr.lnBb6JLiFyf5W9mPuyr7IM0q",
	Salt:     "2d044a28170ba09b9906620794cce1c4d329118fc7ed70a24bb9ff6453d8",
	Profile:  model.Profile{},
}

func (u *UserMock) GetByUsername(_ string) (*model.User, error) {

	var err error
	var user = new(model.User)

	u.Mock.On(r.DB(mockUserDB).Table(userTable).Get(mockUserID)).Return(userMock, nil)

	res, err := r.DB(mockUserDB).Table(userTable).Get(mockUserID).Run(u.Mock)
	if err != nil {
		return nil, err
	}

	if res.IsNil() {
		return nil, nil
	}

	res.One(user)

	return user, nil
}

func (u *UserMock) GetByEmail(_ string) (*model.User, error) {

	var err error
	var user = new(model.User)

	u.Mock.On(r.DB(mockUserDB).Table(userTable).Get(mockUserID)).Return(userMock, nil)

	res, err := r.DB(mockUserDB).Table(userTable).Get(mockUserID).Run(u.Mock)
	if err != nil {
		return nil, err
	}

	if res.IsNil() {
		return nil, nil
	}

	res.One(user)

	return user, nil
}

func (u *UserMock) GetByID(_ string) (*model.User, error) {

	var err error
	var user = new(model.User)

	u.Mock.On(r.DB(mockUserDB).Table(userTable).Get(mockUserID)).Return(userMock, nil)

	res, err := r.DB(mockUserDB).Table(userTable).Get(mockUserID).Run(u.Mock)
	if err != nil {
		return nil, err
	}

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

func (u *UserMock) Insert(_ *model.User) (*model.User, error) {

	var err error
	var opts = r.InsertOpts{ReturnChanges: true}

	u.Mock.On(r.DB(mockUserDB).Table(mockUserDB).Insert(userMock)).Return(nil, nil)

	res, err := r.DB(mockUserDB).Table(mockUserDB).Insert(userMock, opts).RunWrite(u.Mock)

	if err != nil {
		return nil, err
	}

	userMock.ID = res.GeneratedKeys[0]

	return userMock, nil
}

func newUserMock(mock *r.Mock) *UserMock {
	s := new(UserMock)
	s.Mock = mock
	return s
}
