//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

package view

import (
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/util/table"
)

type NamespaceList []*Namespace
type Namespace struct {
	Meta *NamespaceMeta `json:"meta"`
}

type NamespaceMeta struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ClusterID   string `json:"cluster"`
	AccountID   string `json:"account"`
	Owner       string `json:"owner"`
	Endpoint    string `json:"endpoint"`
}

func (nl *NamespaceList) Print() {

	t := table.New([]string{"NAME", "DESCRIPTION", "OWNER", "ENDPOINT"})
	t.VisibleHeader = true

	for _, n := range *nl {
		var data = map[string]interface{}{}

		data["NAME"] = n.Meta.Name
		data["DESCRIPTION"] = n.Meta.Description
		data["OWNER"] = n.Meta.Owner
		data["ENDPOINT"] = n.Meta.Endpoint

		t.AddRow(data)
	}

	println()
	t.Print()
	println()
}

func (n *Namespace) Print() {

	println()
	table.PrintHorizontal(map[string]interface{}{
		"NAME":        n.Meta.Name,
		"DESCRIPTION": n.Meta.Description,
		"OWNER":       n.Meta.Owner,
		"ENDPOINT":    n.Meta.Endpoint,
	})
	println()
}

func FromApiNamespaceView(namespace *views.Namespace) *Namespace {
	var ns = new(Namespace)
	ns.Meta.Name = namespace.Meta.Name
	return ns
}

func FromApiNamespaceListView(namespaces *views.NamespaceList) *NamespaceList {
	var nls = make(NamespaceList, 0)
	for _, namespace := range *namespaces {
		ns := FromApiNamespaceView(namespace)
		nls = append(nls, ns)
	}
	return &nls
}
