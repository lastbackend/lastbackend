package model

import (
	"time"
)

type Image struct {
	UUID      string     `json:"uuid,omitempty"`
	UserID    string     `json:"user,omitempty"`
	ServiceID string     `json:"service,omitempty"`
	StorageID string     `json:"storage,omitempty"`
	BuildID   string     `json:"build,omitempty"`
	Hub       string     `json:"hub,omitempty"`
	Owner     string     `json:"owner,omitempty"`
	Repo      string     `json:"repo,omitempty"`
	Tag       string     `json:"tag,omitempty"`
	Auth      *ImageAuth `json:"auth,omitempty"`
	Created   time.Time  `json:"created,omitempty"`
	Updated   time.Time  `json:"updated,omitempty"`
}

type ImageAuth struct {
	Username string `json:"user"`
	Email    string `json:"email"`
	Host     string `json:"host"`
}
