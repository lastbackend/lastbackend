package model

import (
	"time"
)

type VolumeList []Volume

type Volume struct {
	// Volume uuid, incremented automatically
	ID string `json:"id" gorethink:"id,omitempty"`
	// Volume uuid, incremented automatically
	Project string `json:"project" gorethink:"project,omitempty"`
	// Volume user
	User string `json:"user" gorethink:"user,omitempty"`
	// Volume name
	Name string `json:"name" gorethink:"name,omitempty"`
	// Volume tag lists
	Created time.Time `json:"created" gorethink:"created,omitempty"`
	// Volume updated time
	Updated time.Time `json:"updated" gorethink:"updated,omitempty"`
}
