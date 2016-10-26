package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	UUID         string    // id
	Username     string    // username
	Email        string    // email
	ServiceID    string    // serviceID
	Gravatar     string    // gravatar
	Password     string    // password
	Salt         string    // salt
	GaClientID   string    // ga_client_id
	Active       bool      // active
	Organization bool      // organization
	Balance      float64   // balance
	Created      time.Time // created
	Updated      time.Time // updated
}

// Validation methods
func (u *User) ValidatePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password+string(u.Salt))); err != nil {
		return errors.New("INCORRECT_USER_PASSWORD")
	}

	return nil
}
