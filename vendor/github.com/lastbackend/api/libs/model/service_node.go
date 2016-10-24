package model

import (
	"time"
)

type ServiceNode struct {
	UUID      string
	ServiceID string
	NodeID    string
	Created   time.Time
	Updated   time.Time
}

type ServiceNodes []ServiceNode
