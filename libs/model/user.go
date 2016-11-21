package model

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"golang.org/x/crypto/bcrypt"
	"k8s.io/client-go/1.5/pkg/util/json"
	"time"
)

type User struct {
	ID           string    `json:"id" gorethink:"id"`
	Username     string    `json:"username" gorethink:"username"`
	Email        string    `json:"email" gorethink:"email"`
	Gravatar     string    `json:"gravatar" gorethink:"gravatar"`
	Balance      float32   `json:"balance" gorethink:"balance"`
	Organization bool      `json:"organization" gorethink:"organization"`
	Created      time.Time `json:"created" gorethink:"created"`
	Updated      time.Time `json:"updated" gorethink:"updated"`

	Password string `json:"-" gorethink:"password,omitempty"`
	Salt     string `json:"-" gorethink:"salt,omitempty"`

	Profile Profile `json:"profile" gorethink:"profile"`
}

type Profile struct {
	FirstName string `json:"first_name" gorethink:"first_name"`
	LastName  string `json:"last_name" gorethink:"last_name"`
	Company   string `json:"company" gorethink:"company"`
}

// Validation methods
func (u *User) ValidatePassword(password string) *e.Err {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password+string(u.Salt))); err != nil {
		return e.Account.AccessDenied(err)
	}

	return nil
}

func (u *User) ToJson() ([]byte, *e.Err) {
	buf, err := json.Marshal(u)
	if err != nil {
		return nil, e.User.Unknown(err)
	}

	return buf, nil
}
