package model

import (
	"time"
)

type Hook struct {
	ID        string
	ServiceID NullString
	Hub       NullString
	Owner     NullString
	Repo      NullString
	HookID    NullString
	Active    bool
	Created   time.Time
	Updated   time.Time
}
