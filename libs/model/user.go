package model

import (
	"encoding/json"
	"time"
)

type User struct {
	UUID     string    `json:"uuid"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Gravatar string    `json:"gravatar"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type UserView struct {
	UUID     string    `json:"uuid"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Gravatar string    `json:"gravatar"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

// Convert to json
func (u *User) ToJson() ([]byte, error) {
	buf, err := json.Marshal(u)
	return buf, err
}

// Convert to view
func (u *User) View() *UserView {
	var view = new(UserView)

	view.UUID = u.UUID
	view.Username = u.Username
	view.Email = u.Email
	view.Gravatar = u.Gravatar
	view.Created = u.Created
	view.Updated = u.Updated

	return view
}

// Convert to json
func (u *UserView) ToJson() ([]byte, error) {
	buf, err := json.Marshal(u)
	return buf, err
}
