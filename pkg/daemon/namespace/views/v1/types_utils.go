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

package v1

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/util/table"
)

func New(obj *types.Namespace) *Namespace {
	p := new(Namespace)

	p.Name = obj.Meta.Name
	p.Description = obj.Meta.Description
	p.Updated = obj.Meta.Updated
	p.Created = obj.Meta.Created

	return p
}

func (obj *Namespace) ToJson() ([]byte, error) {
	return json.Marshal(obj)
}

func (p *Namespace) DrawTable() {
	table.PrintHorizontal(map[string]interface{}{
		"Name":        p.Name,
		"Description": p.Description,
		"Created":     p.Created,
		"Updated":     p.Updated,
	})
}

func NewList(obj *types.NamespaceList) *NamespaceList {
	p := new(NamespaceList)
	if obj == nil {
		return nil
	}
	for _, v := range *obj {
		*p = append(*p, *New(&v))
	}
	return p
}

func (obj *NamespaceList) ToJson() ([]byte, error) {
	if obj == nil || len(*obj) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(obj)
}

func (projects *NamespaceList) DrawTable() {
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
