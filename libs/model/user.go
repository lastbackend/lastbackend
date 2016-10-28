package model

import (
	"k8s.io/client-go/1.5/pkg/util/json"
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

// Convert to json
func (u *User) ToJson() (string, error) {
	buf, err := json.Marshal(u)
	return string(buf), err
}
