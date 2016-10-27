package account

import (
	"time"
)

type IAccount interface {
	Create(userID, password string) (*Account, error)
	Info(UserID string) (*Account, error)
}

type Session struct {
	Token string `json:"token"`
}

type Account struct {
	UUID     string    `json:"uuid,omitempty"`
	Password string    `json:"password,omitempty"`
	Salt     string    `json:"salt,omitempty"`
	Balance  string    `json:"balance,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Updated  time.Time `json:"updated,omitempty"`
}

func Create(userID, password string) (*Session, error) {
	return nil, nil
}

func Auth(userID string) (*Session, error) {
	return nil, nil
}

func Info(userID string) (*Account, error) {
	return nil, nil
}
