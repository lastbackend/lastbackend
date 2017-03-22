//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

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
	Organization bool      `json:"organization" gorethink:"organization,omitempty"`
	Balance      float64   `json:"balance" gorethink:"balance,omitempty"`
	Created      time.Time `json:"created" gorethink:"created,omitempty"`
	Updated      time.Time `json:"updated" gorethink:"updated,omitempty"`

	Password string `json:"-" gorethink:"password,omitempty,omitempty"`
	Salt     string `json:"-" gorethink:"salt,omitempty,omitempty"`

	Profile      Profile           `json:"profile" gorethink:"profile,omitempty"`
	Integrations map[string]string `json:"integrations" gorethink:"integrations,omitempty"`
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
