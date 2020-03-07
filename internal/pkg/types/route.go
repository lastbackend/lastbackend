//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package types

import (
	"time"
)

// Route
// swagger:ignore
// swagger:model types_route
type Route struct {
	System
	Meta   RouteMeta   `json:"meta" yaml:"meta"`
	Status RouteStatus `json:"status" yaml:"status"`
	Spec   RouteSpec   `json:"spec" yaml:"spec"`
}

// swagger:ignore
type RouteMap struct {
	System
	Items map[string]*Route
}

// swagger:ignore
type RouteList struct {
	System
	Items []*Route
}

// swagger:ignore
// swagger:model types_route_meta
type RouteMeta struct {
	Meta      `yaml:",inline"`
	SelfLink  RouteSelfLink `json:"self_link"`
	Namespace string        `json:"namespace" yaml:"namespace"`
	Ingress   string        `json:"ingress" yaml:"ingress"`
}

// swagger:model types_route_spec
type RouteSpec struct {
	Selector RouteSelector `json:"selector" yaml:"selector"`
	State    string        `json:"state" yaml:"state"`
	Endpoint string        `json:"endpoint" yaml:"endpoint"`
	Port     uint16        `json:"port" yaml:"port"`
	Rules    []RouteRule   `json:"rules" yaml:"rules"`
	Updated  time.Time     `json:"updated"`
}

type RouteSelector struct {
	Ingress string            `json:"ingress" yaml:"ingress"`
	Label   map[string]string `json:"label" yaml:"label"`
}

// swagger:ignore
// swagger:model types_route_status
// RouteStatus - status of current route state
type RouteStatus struct {
	State   string `json:"state" yaml:"state"`
	Message string `json:"message" yaml:"message"`
}

// swagger:model types_route_rule
type RouteRule struct {
	Service  string `json:"service" yaml:"service"`
	Path     string `json:"path" yaml:"path"`
	Upstream string `json:"upstream" yaml:"upstream"`
	Port     int    `json:"port" yaml:"port"`
}

func (r *Route) SelfLink() *RouteSelfLink {
	return &r.Meta.SelfLink
}

type RouteManifest struct {
	State    string      `json:"state"`
	Endpoint string      `json:"endpoint"`
	Port     uint16      `json:"port"`
	Rules    []RouteRule `json:"rules"`
}

type RouteManifestList struct {
	System
	Items []*RouteManifest
}

type RouteManifestMap struct {
	System
	Items map[string]*RouteManifest
}

func (r *RouteManifest) Set(route *Route) {
	r.State = route.Spec.State
	r.Endpoint = route.Spec.Endpoint
	r.Rules = route.Spec.Rules
	r.Port = route.Spec.Port
}

func NewRouteList() *RouteList {
	dm := new(RouteList)
	dm.Items = make([]*Route, 0)
	return dm
}

func NewRouteMap() *RouteMap {
	dm := new(RouteMap)
	dm.Items = make(map[string]*Route)
	return dm
}

func NewRouteManifestList() *RouteManifestList {
	dm := new(RouteManifestList)
	dm.Items = make([]*RouteManifest, 0)
	return dm
}

func NewRouteManifestMap() *RouteManifestMap {
	dm := new(RouteManifestMap)
	dm.Items = make(map[string]*RouteManifest, 0)
	return dm
}
