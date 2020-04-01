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

package views

import (
	"fmt"

	"github.com/lastbackend/lastbackend/internal/util/table"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type NamespaceList []*Namespace
type Namespace views.Namespace
type NamespaceApplyStatus views.NamespaceApplyStatus

func (n *Namespace) Print() {

	println()
	table.PrintHorizontal(map[string]interface{}{
		"NAME":        n.Meta.Name,
		"DESCRIPTION": n.Meta.Description,
		"ENDPOINT":    n.Meta.Endpoint,
	})
	println()
}

func (nl *NamespaceList) Print() {

	t := table.New([]string{"NAME", "DESCRIPTION", "ENDPOINT"})
	t.VisibleHeader = true

	for _, n := range *nl {
		var data = map[string]interface{}{}

		data["NAME"] = n.Meta.Name
		data["DESCRIPTION"] = n.Meta.Description
		data["ENDPOINT"] = n.Meta.Endpoint

		t.AddRow(data)
	}

	println()
	t.Print()
	println()
}

func (ns *NamespaceApplyStatus) Print() {

	var printEntity = func(kind string, status map[string]bool) {
		fmt.Printf("%s:\n", kind)
		t := table.New([]string{"NAME", "STATUS"})
		t.VisibleHeader = false

		for n, s := range status {
			var data = map[string]interface{}{}
			data["NAME"] = n
			if s {
				data["STATUS"] = "Provisioned"
			} else {
				data["STATUS"] = "Error"
			}
			t.AddRow(data)
		}
		t.Print()
		println()
	}

	if len(ns.Services) > 0 {
		printEntity("Services", ns.Services)
	}

	if len(ns.Configs) > 0 {
		printEntity("Configs", ns.Configs)
	}

	if len(ns.Secrets) > 0 {
		printEntity("Secrets", ns.Secrets)
	}

	if len(ns.Volumes) > 0 {
		printEntity("Volumes", ns.Volumes)
	}

	if len(ns.Routes) > 0 {
		printEntity("Routes", ns.Routes)
	}

	if len(ns.Jobs) > 0 {
		printEntity("Jobs", ns.Jobs)
	}

}

func FromApiNamespaceView(namespace *views.Namespace) *Namespace {

	if namespace == nil {
		return nil
	}

	item := Namespace(*namespace)
	return &item
}

func FromApiNamespaceListView(namespaces *views.NamespaceList) *NamespaceList {
	var items = make(NamespaceList, 0)
	for _, namespace := range *namespaces {
		items = append(items, FromApiNamespaceView(namespace))
	}
	return &items
}

func FromApiNamespaceStatusView(status *views.NamespaceApplyStatus) *NamespaceApplyStatus {
	ns := new(NamespaceApplyStatus)
	ns.Configs = status.Configs
	ns.Secrets = status.Secrets
	ns.Services = status.Services
	ns.Volumes = status.Volumes
	ns.Routes = status.Routes
	ns.Jobs = status.Jobs
	return ns
}
