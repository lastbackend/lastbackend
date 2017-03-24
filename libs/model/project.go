//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package model

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/util/table"
	"time"
)

type ProjectList []Project

type Project struct {
	// Project user
	User string `json:"user"`
	// Project name
	Name string `json:"name"`
	// Project description
	Description string `json:"description"`
	// Project labels
	Labels map[string]string `json:"labels,omitempty"`
	// Project created time
	Created time.Time `json:"created"`
	// Project updated time
	Updated time.Time `json:"updated"`
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
			"Name":        p.Name,
			"Description": p.Description,
			"Created":     p.Created.String()[:10],
			"Updated":     p.Updated.String()[:10],
		})
	}

	t.AddRow(map[string]interface{}{})

	t.Print()
}
