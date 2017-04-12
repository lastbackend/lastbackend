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
