package pgsql

import (
	"github.com/lastbackend/lastbackend/libs/adapter"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
)

// User Service type for interface in interfaces folder
type UserService struct{}

type userModel struct {
	id           nullString
	username     nullString
	email        nullString
	serviceID    nullString
	gravatar     nullString
	password     nullString
	salt         nullString
	gaclientid   nullString
	active       nullBool
	organization nullBool
	balance      nullFloat64
	created      nullTime
	updated      nullTime
}

func (um *userModel) convert() *model.User {

	var u = new(model.User)

	u.UUID = um.id.String
	u.Username = um.username.String
	u.Email = um.email.String
	u.Gravatar = um.gravatar.String
	u.Password = um.password.String
	u.Salt = um.salt.String
	u.GaClientID = um.gaclientid.String
	u.Active = um.active.Bool
	u.Organization = um.organization.Bool
	u.Balance = um.balance.Float64
	u.Created = um.created.Time
	u.Updated = um.updated.Time

	return u
}

// all arguments should be started with lowercase: email, gravatar, password
func (UserService) Insert(db adapter.IDatabase, username, email, gravatar, password, salt string) (*string, *e.Err) {

	var err error
	var id nullString
	var uname, uemail *string

	const sqlstr = `
		INSERT INTO users (username, email, gravatar, password, salt)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;`

	uname = &username
	if username == "" {
		uname = nil
	}

	uemail = &email
	if username == "" {
		uname = nil
	}

	err = db.QueryRow(sqlstr, uname, uemail, gravatar, password, salt).Scan(&id)

	if err != nil {
		return nil, e.User.Unknown(err)
	}

	return &id.String, nil
}
