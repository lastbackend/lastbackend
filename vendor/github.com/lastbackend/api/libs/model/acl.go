package model

import (
	"time"
)

type Acl struct {
	ID      string
	UserID  string
	OrgID   string
	Role    string
	Created time.Time
	Updated time.Time
}
