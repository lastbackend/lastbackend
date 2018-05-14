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
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"strings"
)

// Route
// swagger:ignore
// swagger:model types_route
type Route struct {
	Meta   RouteMeta   `json:"meta" yaml:"meta"`
	Spec   RouteSpec   `json:"spec" yaml:"spec"`
	Status RouteStatus `json:"status" yaml:"status"`
}
// swagger:ignore
type RouteMap map[string]*Route
// swagger:ignore
type RouteList []*Route

// swagger:ignore
// swagger:model types_route_meta
type RouteMeta struct {
	Meta      `yaml:",inline"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Security  bool   `json:"security" yaml:"security"`
}

// swagger:model types_route_spec
type RouteSpec struct {
	Domain string       `json:"domain" yaml:"domain"`
	Rules  []*RouteRule `json:"rules" yaml:"rules"`
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

// swagger:ignore
type RouterConfig struct {
	Name      string      `json:"id" yaml:"id"`
	Hash      string      `json:"hash" yaml:"hash"`
	RootPath  string      `json:"-" yaml:"-"`
	Upstreams []*Upstream `json:"upstreams" yaml:"upstreams"`
	Server    RouteServer `json:"server" yaml:"server"`
}

// swagger:ignore
type RouteServer struct {
	Hostname  string          `json:"hostname" yaml:"hostname"`
	Port      int             `json:"port" yaml:"port"`
	Protocol  string          `json:"protocol" yaml:"protocol"`
	Locations []*RoteLocation `json:"locations" yaml:"locations"`
}

// swagger:ignore
type Upstream struct {
	Name    string `json:"name" yaml:"name"`
	Address string `json:"address" yaml:"address"`
}

// swagger:ignore
type RoteLocation struct {
	Path      string `json:"path" yaml:"path"`
	ProxyPass string `json:"proxy_pass" yaml:"proxy_pass"`
}

func (r *Route) SelfLink() string {
	if r.Meta.SelfLink == "" {
		r.Meta.SelfLink = fmt.Sprintf("%s:%s", r.Meta.Namespace, r.Meta.Name)
	}
	return r.Meta.SelfLink
}

func (r *Route) GetRouteConfig() *RouterConfig {
	var RouterConfig = new(RouterConfig)

	RouterConfig.Name = r.Meta.Name

	RouterConfig.Server.Hostname = r.Spec.Domain
	RouterConfig.Server.Protocol = "http"
	RouterConfig.Server.Port = 80

	if r.Meta.Security {
		RouterConfig.Server.Protocol = "https"
		RouterConfig.Server.Port = 443
	}

	RouterConfig.Upstreams = make([]*Upstream, 0)
	RouterConfig.Server.Locations = make([]*RoteLocation, 0)
	for _, rule := range r.Spec.Rules {

		name := generator.GetUUIDV4()

		RouterConfig.Upstreams = append(RouterConfig.Upstreams, &Upstream{
			Name:    name,
			Address: strings.ToLower(fmt.Sprintf("%s:%d", rule.Endpoint, rule.Port)),
		})

		RouterConfig.Server.Locations = append(RouterConfig.Server.Locations, &RoteLocation{
			Path:      rule.Path,
			ProxyPass: strings.ToLower(fmt.Sprintf("http://%s", name)),
		})
	}

	return RouterConfig
}

// swagger:ignore
type RouteCreateOptions struct {
	Name     string       `json:"name"`
	Domain   string       `json:"domain"`
	Security bool         `json:"security"`
	Rules    []RuleOption `json:"rules"`
}

// swagger:ignore
type RouteUpdateOptions struct {
	Security bool         `json:"security"`
	Rules    []RuleOption `json:"rules"`
}

// swagger:ignore
type RouteRemoveOptions struct {
	Force bool `json:"force"`
}

// swagger:ignore
type RuleOption struct {
	Service  string `json:"service"`
	Endpoint string `json:"endpoint"`
	Path     string `json:"path"`
	Port     int    `json:"port"`
}
