package model

import (
	"encoding/json"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Account struct {
	UUID     string    `json:"uuid"`
	UserID   string    `json:"user_id"`
	Password string    `json:"password"`
	Salt     string    `json:"salt"`
	Updated  time.Time `json:"updated"`
	Created  time.Time `json:"created"`
}

type AccountView struct {
	UUID    string    `json:"uuid"`
	Updated time.Time `json:"updated"`
	Created time.Time `json:"created"`
}

// Validation methods
func (a *Account) ValidatePassword(password string) *e.Err {
	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password+string(a.Salt))); err != nil {
		return e.Account.AccessDenied(err)
	}

	return nil
}

// Convert to json
func (a *Account) ToJson() ([]byte, error) {
	buf, err := json.Marshal(a)
	return buf, err
}

// Convert to view
func (a *Account) View() *AccountView {
	var view = new(AccountView)

	view.UUID = a.UUID
	view.Created = a.Created
	view.Updated = a.Updated

	return view
}

// Convert to json
func (a *AccountView) ToJson() ([]byte, error) {
	buf, err := json.Marshal(a)
	return buf, err
}
