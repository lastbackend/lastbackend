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
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/util/table"
)

func New(obj *types.Namespace) *Namespace {
	p := Namespace{}
	p.Meta.Name = obj.Meta.Name
	p.Meta.Description = obj.Meta.Description
	p.Meta.Labels = obj.Meta.Labels
	p.Meta.Created = obj.Meta.Created
	p.Meta.Updated = obj.Meta.Updated
	return &p
}

func (n *Namespace) ToJson() ([]byte, error) {
	return json.Marshal(n)
}

func (n *Namespace) DrawTable() {
	var labels []string

	for _, v := range n.Meta.Labels {
		labels = append(labels, v)
	}

	table.PrintHorizontal(map[string]interface{}{
		"NAME":        n.Meta.Name,
		"DESCRIPTION": n.Meta.Description,
		"LABELS":      labels,
		"CREATED":     n.Meta.Created.String()[:10],
		"UPDATED":     n.Meta.Updated.String()[:10],
	})
}

func NewList(obj types.NamespaceList) *NamespaceList {
	p := NamespaceList{}
	if obj == nil {
		return nil
	}

	for _, v := range obj {
		p = append(p, New(v))
	}
	return &p
}

func (ns *NamespaceList) ToJson() ([]byte, error) {

	if ns == nil || len(*ns) == 0 {
		return make([]byte, 0), nil
	}

	return json.Marshal(ns)
}

func (ns *NamespaceList) DrawTable() {
	t := table.New([]string{"NAME", "DESCRIPTION", "LABELS", "CREATED", "UPDATED"})
	t.VisibleHeader = true

	for _, n := range *ns {
		var labels []string

		for _, v := range n.Meta.Labels {
			labels = append(labels, v)
		}

		t.AddRow(map[string]interface{}{
			"NAME":        n.Meta.Name,
			"DESCRIPTION": n.Meta.Description,
			"LABELS":      labels,
			"CREATED":     n.Meta.Created.String()[:10],
			"UPDATED":     n.Meta.Updated.String()[:10],
		})
	}

	t.AddRow(map[string]interface{}{})

	t.Print()
}
