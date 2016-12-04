package model

import (
	"encoding/json"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/pkg/util/table"
	"time"
)

type ProjectList []Project

type Project struct {
	// Project uuid, incremented automatically
	ID string `json:"id" gorethink:"id,omitempty"`
	// Project user
	User string `json:"user" gorethink:"user,omitempty"`
	// Project name
	Name string `json:"name" gorethink:"name,omitempty"`
	// Project description
	Description string `json:"description" gorethink:"description,omitempty"`
	// Project created time
	Created time.Time `json:"created" gorethink:"created,omitempty"`
	// Project updated time
	Updated time.Time `json:"updated" gorethink:"updated,omitempty"`
}

func (p *Project) ToJson() ([]byte, *e.Err) {
	buf, err := json.Marshal(p)
	if err != nil {
		return nil, e.Project.Unknown(err)
	}

	return buf, nil
}

func (p *Project) DrawTable() {
	t := table.New([]string{"ID", "Name", "Created", "Updated"})
	t.AddRow(map[string]interface{}{
		"ID":      p.ID,
		"Name":    p.Name,
		"Created": p.Created.String()[:10],
		"Updated": p.Updated.String()[:10],
	})
	t.Print()
}

func (p *ProjectList) ToJson() ([]byte, *e.Err) {

	if p == nil {
		return []byte("[]"), nil
	}

	buf, err := json.Marshal(p)
	if err != nil {
		return nil, e.Project.Unknown(err)
	}

	return buf, nil
}

func (p *ProjectList) DrawTable() {
	t := table.New([]string{"ID", "Name", "Created", "Updated"})

	for _, project := range *p {
		t.AddRow(map[string]interface{}{
			"ID":      project.ID,
			"Name":    project.Name,
			"Created": project.Created.String()[:10],
			"Updated": project.Updated.String()[:10],
		})
	}

	t.Print()
}
