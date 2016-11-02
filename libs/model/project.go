package model

import "time"

type ProjectList []Project

type Project struct {
	// Project uuid, incremented automatically
	ID string `json:"id, omitempty" gorethink:"id,omitempty"`
	// Project user
	User string `json:"user, omitempty" gorethink:"user,omitempty"`
	// Project name
	Name string `json:"name, omitempty" gorethink:"name,omitempty"`
	// Project description
	Description string `json:"description, omitempty" gorethink:"description,omitempty"`
	// Project created time
	Created time.Time `json:"created, omitempty" gorethink:"created,omitempty"`
	// Project updated time
	Updated time.Time `json:"updated, omitempty" gorethink:"updated,omitempty"`
}
