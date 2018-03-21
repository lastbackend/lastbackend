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
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"io"
	"io/ioutil"
	"strings"
)

type Route struct {
	Meta   RouteMeta   `json:"meta" yaml:"meta"`
	State  RouteState  `json:"state" yaml:"state"`
	Spec   RouteSpec   `json:"spec" yaml:"spec"`
	Status RouteStatus `json:"status" yaml:"status"`
}

type RouteList map[string]*Route

type RouteMeta struct {
	Meta      `yaml:",inline"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Security  bool   `json:"security" yaml:"security"`
}

type RouteState struct {
	Destroy   bool `json:"destroy" yaml:"destroy"`
	Provision bool `json:"provision" yaml:"provision"`
}

type RouteSpec struct {
	Domain string       `json:"domain" yaml:"domain"`
	Rules  []*RouteRule `json:"rules" yaml:"rules"`
}

type RouteStatus struct {
	// Pod stage
	Stage string `json:"stage" yaml:"stage"`
	// Pod state message
	Message string `json:"message" yaml:"message"`
}

type RouteRule struct {
	Path     string `json:"path" yaml:"path"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	Port     int    `json:"port" yaml:"port"`
}

type RouteOptions struct {
	Subdomain string        `json:"subdomain" yaml:"subdomain"`
	Domain    string        `json:"domain" yaml:"domain"`
	Custom    bool          `json:"custom" yaml:"custom"`
	Security  bool          `json:"security" yaml:"security"`
	Rules     []RulesOption `json:"rules" yaml:"rules"`
}

type RouteStateEvent struct {
	ID     string `json:"id" yaml:"id"`
	Status string `json:"status" yaml:"status"`
}

type RouterConfig struct {
	ID        string            `json:"id" yaml:"id"`
	Hash      string            `json:"hash" yaml:"hash"`
	RootPath  string            `json:"-" yaml:"-"`
	State     RouteState        `json:"state" yaml:"state"`
	Upstreams []*UpstreamServer `json:"upstreams" yaml:"upstreams"`
	Server    RouteServer       `json:"server" yaml:"server"`
}

type RouteServer struct {
	Hostname  string          `json:"hostname" yaml:"hostname"`
	Port      int             `json:"port" yaml:"port"`
	Protocol  string          `json:"protocol" yaml:"protocol"`
	Locations []*RoteLocation `json:"locations" yaml:"locations"`
}

type UpstreamServer struct {
	Name    string `json:"name" yaml:"name"`
	Address string `json:"address" yaml:"address"`
}

type RoteLocation struct {
	Path      string `json:"path" yaml:"path"`
	ProxyPass string `json:"proxy_pass" yaml:"proxy_pass"`
}

type RulesOption struct {
	Service *string `json:"service" yaml:"service"`
	Path    string  `json:"path" yaml:"path"`
	Port    *int    `json:"port" yaml:"port"`
}

func (r *Route) SelfLink() string {
	if r.Meta.SelfLink == "" {
		r.Meta.SelfLink = fmt.Sprintf("%s:%s", r.Meta.Namespace, r.Meta.Name)
	}
	return r.Meta.SelfLink
}

func (r *Route) GetRouteConfig() *RouterConfig {
	var RouterConfig = new(RouterConfig)

	RouterConfig.ID = r.Meta.Name
	RouterConfig.State = r.State

	RouterConfig.Server.Hostname = r.Spec.Domain
	RouterConfig.Server.Protocol = "http"
	RouterConfig.Server.Port = 80

	if r.Meta.Security {
		RouterConfig.Server.Protocol = "https"
		RouterConfig.Server.Port = 443
	}

	RouterConfig.Upstreams = make([]*UpstreamServer, 0)
	RouterConfig.Server.Locations = make([]*RoteLocation, 0)
	for _, rule := range r.Spec.Rules {

		id := generator.GetUUIDV4()

		RouterConfig.Upstreams = append(RouterConfig.Upstreams, &UpstreamServer{
			Name:    id,
			Address: strings.ToLower(fmt.Sprintf("%s:%d", rule.Endpoint, rule.Port)),
		})

		RouterConfig.Server.Locations = append(RouterConfig.Server.Locations, &RoteLocation{
			Path:      rule.Path,
			ProxyPass: strings.ToLower(fmt.Sprintf("http://%s", id)),
		})
	}

	return RouterConfig
}

func (s *RouteOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	log.V(logLevel).Debug("Request: Route: decode and validate data")

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.V(logLevel).Errorf("Request: Route: decode and validate data for creating err: %s", err)
		return errors.New("route").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.V(logLevel).Errorf("Request: Route: convert struct from json err: %s", err)
		return errors.New("route").IncorrectJSON(err)
	}

	for _, rule := range s.Rules {
		if rule.Path == "" {
			rule.Path = "/"
		}

		if rule.Service == nil || len(*rule.Service) == 0 {
			log.V(logLevel).Error("Request: Route: parameter service can not be empty")
			return errors.New("route").BadParameter("service")
		}

		if rule.Port == nil {
			log.V(logLevel).Error("Request: Route: parameter port can not be empty")
			return errors.New("route").BadParameter("port")
		}
	}

	return nil
}
