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

package types

import (
	"fmt"
	"time"
)

// Route
// swagger:ignore
// swagger:model types_route
type Route struct {
	Runtime
	Meta   RouteMeta   `json:"meta" yaml:"meta"`
	Spec   RouteSpec   `json:"spec" yaml:"spec"`
	Status RouteStatus `json:"status" yaml:"status"`
}

// swagger:ignore
type RouteMap struct {
	Runtime
	Items map[string]*Route
}

// swagger:ignore
type RouteList struct {
	Runtime
	Items []*Route
}

// swagger:ignore
// swagger:model types_route_meta
type RouteMeta struct {
	Meta      `yaml:",inline"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Security  bool   `json:"security" yaml:"security"`
}

// swagger:model types_route_spec
type RouteSpec struct {
	Security bool        `json:"security" yaml:"security"`
	Domain   string      `json:"domain" yaml:"domain"`
	Port     uint16      `json:"port" yaml:"port"`
	Rules    []RouteRule `json:"rules" yaml:"rules"`
	Updated  time.Time   `json:"updated"`
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
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	Port     int    `json:"port" yaml:"port"`
}

func (r *Route) SelfLink() string {
	if r.Meta.SelfLink == "" {
		r.Meta.SelfLink = r.CreateSelfLink(r.Meta.Namespace, r.Meta.Name)
	}
	return r.Meta.SelfLink
}

func (r *Route) CreateSelfLink(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

type RouteManifest struct {
	State    string      `json:"state"`
	Domain   string      `json:"domain"`
	Port     uint16      `json:"port"`
	Endpoint string      `json:"endpoint"`
	Rules    []RouteRule `json:"rules"`
}

type RouteManifestList struct {
	Runtime
	Items []*RouteManifest
}

type RouteManifestMap struct {
	Runtime
	Items map[string]*RouteManifest
}

func (r *RouteManifest) Set(route *Route) {
	r.Domain = route.Spec.Domain
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
	dm.Items = make(map[string]*RouteManifest)
	return dm
}
