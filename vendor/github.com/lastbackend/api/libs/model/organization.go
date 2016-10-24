package model

import (
	"time"
)

type Organization struct {
	ID       string
	Name     string
	Email    string
	Gravatar string
	Created  time.Time // created
	Updated  time.Time // updated
}

type Organizations []Organization
