package model

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"golang.org/x/crypto/bcrypt"
	"k8s.io/client-go/1.5/pkg/util/json"
	"time"
)

type User struct {
	ID           string    `json:"id,omitempty" gorethink:"id,omitempty"`
	Username     string    `json:"username,omitempty" gorethink:"username,omitempty"`
	Email        string    `json:"email,omitempty" gorethink:"email,omitempty"`
	Gravatar     string    `json:"gravatar,omitempty" gorethink:"gravatar,omitempty"`
	Balance      int       `json:"balance,omitempty" gorethink:"balance,omitempty"`
	Organization bool      `json:"organization,omitempty" gorethink:"organization,omitempty"`
	Created      time.Time `json:"created,omitempty" gorethink:"created,omitempty"`
	Updated      time.Time `json:"updated,omitempty" gorethink:"updated,omitempty"`

	Password string `json:"-" gorethink:"password,omitempty"`
	Salt     string `json:"-" gorethink:"salt,omitempty"`

	Profile Profile `json:"profile" gorethink:"profile"`
}

type Profile struct {
	FirstName string `json:"first_name,omitempty" gorethink:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty" gorethink:"last_name,omitempty"`
	Company   string `json:"company,omitempty" gorethink:"company,omitempty"`
}

// Validation methods
func (u *User) ValidatePassword(password string) *e.Err {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password+string(u.Salt))); err != nil {
		return e.Account.AccessDenied(err)
	}

	return nil
}

func (u *User) ToJson() ([]byte, *e.Err) {
	byte, err := json.Marshal(u)
	if err != nil {
		return nil, e.User.Unknown(err)
	}

	return byte, nil
}
