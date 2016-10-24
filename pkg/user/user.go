package user

import (
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/account"
	"time"
)

type IUser interface {
	Create(username, email, password string) *User
	Info(username string) *User
}

type User struct {
	UUID         string    `json:"uuid,omitempty"`
	Username     string    `json:"username,omitempty"`
	Email        string    `json:"email,omitempty"`
	Gravatar     string    `json:"gravatar,omitempty"`
	Active       bool      `json:"active,omitempty"`
	Organization bool      `json:"organization,omitempty"`
	Balance      float64   `json:"balance,omitempty"`
	Created      time.Time `json:"created,omitempty"`
	Updated      time.Time `json:"updated,omitempty"`
}

func Create(username, email string) (*User, error) {

	var ctx = context.Get()

	ctx.K8S.LB().Accounts().Create(account.Account{})

	return nil, nil
}

func Info(username string) (*User, error) {
	return nil, nil
}
