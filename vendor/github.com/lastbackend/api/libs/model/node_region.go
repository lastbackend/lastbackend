package model

import "time"

type NodeRegion struct {
	UUID      string
	UserID    string
	Name      string
	Code      string
	Available bool
	Created   time.Time
	Updated   time.Time
}

type NodeRegions []NodeRegion
