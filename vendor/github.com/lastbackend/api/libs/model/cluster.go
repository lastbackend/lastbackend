package model

import (
	"time"
)

type Cluster struct {
	UUID    string
	UserID  string
	Name    string
	Created time.Time
	Updated time.Time
}

type Clusters []Cluster
