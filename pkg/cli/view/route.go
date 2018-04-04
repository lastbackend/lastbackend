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
	"time"
)

type RouteList []*Route
type Route struct {
	Meta   RouteMeta   `json:"meta"`
	Spec   RouteSpec   `json:"spec"`
	Status RouteStatus `json:"status"`
}

type RouteMeta struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	SelfLink  string    `json:"self_link"`
	Security  bool      `json:"security"`
	Updated   time.Time `json:"updated"`
	Created   time.Time `json:"created"`
}

type RouteSpec struct {
	Domain string       `json:"domain"`
	Rules  []*RouteRule `json:"rules"`
}

type RouteRule struct {
	Service  string `json:"service"`
	Path     string `json:"path"`
	Endpoint string `json:"endpoint"`
	Port     int    `json:"port"`
}

type RouteStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

func (rl *RouteList) Print() {

	t := table.New([]string{"NAMESPACE", "NAME", "DOMAIN", "HTTPS", "STATUS"})
	t.VisibleHeader = true

	for _, r := range *rl {
		var data = map[string]interface{}{}
		data["NAMESPACE"] = r.Meta.Namespace
		data["NAME"] = r.Meta.Name
		data["DOMAIN"] = r.Spec.Domain
		data["HTTPS"] = r.Meta.Security
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
	data["HTTPS"] = r.Meta.Security
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
	var item = new(Route)
	item.Meta.Name = route.Meta.Name
	item.Meta.Namespace = route.Meta.Namespace
	item.Meta.Security = route.Meta.Security
	item.Meta.Created = route.Meta.Created
	item.Meta.Updated = route.Meta.Updated

	item.Status.State = route.Status.State
	item.Status.Message = route.Status.Message

	item.Spec.Domain = route.Spec.Domain

	for _, rule := range route.Spec.Rules {
		item.Spec.Rules = append(item.Spec.Rules, &RouteRule{
			Service:  rule.Service,
			Path:     rule.Path,
			Endpoint: rule.Endpoint,
			Port:     rule.Port,
		})
	}

	return item
}

func FromApiRouteListView(routes *views.RouteList) *RouteList {
	var items = make(RouteList, 0)
	for _, route := range *routes {
		items = append(items, FromApiRouteView(route))
	}
	return &items
}
