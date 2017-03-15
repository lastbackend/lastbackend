package model

import (
	"encoding/json"
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
	// Project labels
	Labels map[string]string  `json:"labels,omitempty" gorethink:"-"`
	// Project created time
	Created time.Time `json:"created" gorethink:"created,omitempty"`
	// Project updated time
	Updated time.Time `json:"updated" gorethink:"updated,omitempty"`
}

func (p *Project) ToJson() ([]byte, error) {
	buf, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (p *Project) DrawTable() {
	table.PrintHorizontal(map[string]interface{}{
		"ID":          p.ID,
		"Name":        p.Name,
		"Description": p.Description,
		"Created":     p.Created,
		"Updated":     p.Updated,
	})
}

func (p *ProjectList) ToJson() ([]byte, error) {

	if p == nil {
		return []byte("[]"), nil
	}

	buf, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (projects *ProjectList) DrawTable() {
	t := table.New([]string{"ID", "Name", "Description", "Created", "Updated"})
	t.VisibleHeader = true

	for _, p := range *projects {
		t.AddRow(map[string]interface{}{
			"ID":          p.ID,
			"Name":        p.Name,
			"Description": p.Description,
			"Created":     p.Created.String()[:10],
			"Updated":     p.Updated.String()[:10],
		})
	}

	t.AddRow(map[string]interface{}{})

	t.Print()
}
