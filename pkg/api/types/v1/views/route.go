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

import "time"

// swagger:model views_route
type Route struct {
	Meta   RouteMeta   `json:"meta"`
	Spec   RouteSpec   `json:"spec"`
	Status RouteStatus `json:"status"`
}

// swagger:ignore
type RouteMap map[string]*Route

// swagger:model views_route_list
type RouteList []*Route

// swagger:model views_route_meta
type RouteMeta struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	SelfLink  string    `json:"self_link"`
	Updated   time.Time `json:"updated"`
	Created   time.Time `json:"created"`
}

// swagger:model views_route_spec
type RouteSpec struct {
	Domain string       `json:"domain"`
	Port   uint16       `json:"port"`
	Rules  []*RouteRule `json:"rules"`
}

// swagger:model views_route_rule
type RouteRule struct {
	Service  string `json:"service"`
	Path     string `json:"path"`
	Endpoint string `json:"endpoint"`
	Port     int    `json:"port"`
}

// swagger:model views_route_status
type RouteStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
}
