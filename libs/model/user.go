package model

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID           string    `json:"id" gorethink:"id,omitempty"`
	Username     string    `json:"username" gorethink:"username,omitempty"`
	Email        string    `json:"email" gorethink:"email,omitempty"`
	Gravatar     string    `json:"gravatar" gorethink:"gravatar,omitempty"`
	Balance      float32   `json:"balance" gorethink:"balance,omitempty"`
	Organization bool      `json:"organization" gorethink:"organization,omitempty"`
	Created      time.Time `json:"created" gorethink:"created,omitempty"`
	Updated      time.Time `json:"updated" gorethink:"updated,omitempty"`

	Password string `json:"-" gorethink:"password,omitempty,omitempty"`
	Salt     string `json:"-" gorethink:"salt,omitempty,omitempty"`

	Profile Profile `json:"profile" gorethink:"profile,omitempty"`
}

type Profile struct {
	FirstName string `json:"first_name" gorethink:"first_name,omitempty"`
	LastName  string `json:"last_name" gorethink:"last_name,omitempty"`
	Company   string `json:"company" gorethink:"company,omitempty"`
}

// Validation methods
func (u *User) ValidatePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password+string(u.Salt))); err != nil {
		return err
	}

	return nil
}

func (u *User) ToJson() ([]byte, error) {
	buf, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
