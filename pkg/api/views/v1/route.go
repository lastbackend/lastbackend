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

package v1

import "time"

type Route struct {
	Meta  RouteMeta    `json:"meta"`
	Spec  RouteSpec 	 `json:"spec"`
	State RouteState   `json:"state"`
}

type RouteMeta struct {
	Domain      string    `json:"domain"`
	Namespace   string    `json:"namespace"`
	Security    bool      `json:"security"`
	Updated     time.Time `json:"updated"`
	Created     time.Time `json:"created"`
}

type RouteSpec struct {

}

type RouteRule struct {
	Path     string `json:"path"`
	Service  string `json:"service"`
	Port     int    `json:"port"`
	Security bool   `json:"security"`
}

type RouteState struct {
	Destroy   bool `json:"destroy"`
	Provision bool `json:"provision"`
}

type RouteList map[string]*Route
