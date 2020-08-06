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
	"github.com/lastbackend/lastbackend/internal/util/table"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type RouteList []*Route
type Route views.Route

func (rl *RouteList) Print() {

	t := table.New([]string{"NAMESPACE", "NAME", "DOMAIN", "HTTPS", "STATUS"})
	t.VisibleHeader = true

	for _, r := range *rl {
		var data = map[string]interface{}{}
		data["NAMESPACE"] = r.Meta.Namespace
		data["NAME"] = r.Meta.Name
		data["DOMAIN"] = r.Spec.Domain
		data["PORT"] = r.Spec.Port
		data["STATUS"] = r.Status.State
		t.AddRow(data)
	}

	println()
	t.Print()
	println()
}

func (r *Route) Print() {
	var data = map[string]interface{}{}
	data["NAME"] = r.Meta.Name
	data["NAMESPACE"] = r.Meta.Namespace
	data["DOMAIN"] = r.Spec.Domain
	data["PORT"] = r.Spec.Port
	data["STATUS"] = r.Status.State
	println()
	table.PrintHorizontal(data)
	println()

	t := table.New([]string{"PATH", "SERVICE", "ENDPOINT", "PORT"})
	t.VisibleHeader = true

	for _, r := range r.Spec.Rules {
		var data = map[string]interface{}{}
		data["PATH"] = r.Path
		data["SERVICE"] = r.Service
		data["ENDPOINT"] = r.Endpoint
		data["PORT"] = r.Port
		t.AddRow(data)
	}

	println()
	t.Print()
	println()
}

func FromApiRouteView(route *views.Route) *Route {

	if route == nil {
		return nil
	}

	item := Route(*route)
	return &item
}

func FromApiRouteListView(routes *views.RouteList) *RouteList {
	var items = make(RouteList, 0)
	for _, route := range *routes {
		items = append(items, FromApiRouteView(route))
	}
	return &items
}
