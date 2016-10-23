package model

import (
	"time"
)

type Build struct {
	UUID  string
	Image struct {
		Hub   string
		Owner string
		Repo  string
		Tag   string
		Auth  struct {
			Username string
			Password string
			Email    string
			Host     string
		}
	}
	Sources struct {
		Hub    string
		Owner  string
		Repo   string
		Branch string
		Token  string
	}
	Commit struct {
		Hash    string
		Message string
		Owner   string
		Email   string
		Date    time.Time
	}

	StorageID string
	ServiceID string
	UserID    string

	Machine    string
	Status     string
	Error      string
	Processing bool
	Done       bool
	Started    time.Time
	Created    time.Time
	Updated    time.Time
}

type Builds []Build
