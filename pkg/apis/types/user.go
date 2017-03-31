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

package types

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	userInfo

	Username string       `json:"username"`
	Security UserSecurity `json:"security"`
	Emails   UserEmails   `json:"emails"`
	Profile  UserProfile  `json:"profile"`
	Vendors  UserVendors  `json:"vendors"`
}

type userInfo struct {
	Gravatar string    `json:"gravatar"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type userPass struct {
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type UserEmails map[string]bool
type UserVendors map[string]string
type UserInfo struct{ userInfo }
type UserPassword struct{ userPass }

type UserProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserSecurity struct {
	Pass UserPassword `json:"pass"`
	SSH  []UserSSH    `json:"ssh"`
}

type UserSSH struct {
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	Key         string `json:"key"`
}

// Validation methods
func (p *UserPassword) ValidatePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password+string(p.Salt))); err != nil {
		return err
	}

	return nil
}

// Get primary email
func (p *UserEmails) GetDefault() string {
	for k, v := range *p {
		if v == true {
			return k
		}
	}

	var email string
	for email = range *p {
		break
	}

	return email
}

func (u *User) ToJson() ([]byte, error) {
	buf, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
