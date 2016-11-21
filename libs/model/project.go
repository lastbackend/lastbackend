package model

import (
	"encoding/json"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"time"
)

type ProjectList []Project

type Project struct {
	// Project uuid, incremented automatically
	ID string `json:"id" gorethink:"id"`
	// Project user
	User string `json:"user" gorethink:"user"`
	// Project name
	Name string `json:"name" gorethink:"name"`
	// Project description
	Description string `json:"description" gorethink:"description"`
	// Project created time
	Created time.Time `json:"created" gorethink:"created"`
	// Project updated time
	Updated time.Time `json:"updated" gorethink:"updated"`
}

func (p *Project) ToJson() ([]byte, *e.Err) {
	buf, err := json.Marshal(p)
	if err != nil {
		return nil, e.Project.Unknown(err)
	}

	return buf, nil
}

func (p *ProjectList) ToJson() ([]byte, *e.Err) {
	buf, err := json.Marshal(p)
	if err != nil {
		return nil, e.Project.Unknown(err)
	}

	return buf, nil
}
