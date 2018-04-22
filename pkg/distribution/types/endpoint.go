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

import "fmt"

const (
	// EndpointSpecRouteStrategyRR - round robin balancing strategy type
	EndpointSpecRouteStrategyRR = "rr"
	// EndpointSpecBindStrategyDefault - default scheduling endpoint across all nodes
	EndpointSpecBindStrategyDefault = "default"
)

// Endpoint - service endpoint
type Endpoint struct {
	Meta   EndpointMeta   `json:"meta"`
	Status EndpointStatus `json:"status"`
	Spec   EndpointSpec   `json:"spec"`
}

// EndpointMeta - endpoint meta data
type EndpointMeta struct {
	Meta
	// Namespace name
	Namespace string `json:"namespace"`
	// Service name
	Service string `json:"service"`
}

// EndpointStatus - endpoint status
type EndpointStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

// EndpointSpec spec data
type EndpointSpec struct {
	IP            string                         `json:"ip"`
	Ports         []int                          `json:"ports"`
	Backends      map[string]EndpointSpecBackend `json:"backends"`
	RouteStrategy string                         `json:"route_strategy"`
	Policy        string                         `json:"policy"`
	BindStrategy  string                         `json:"bind_strategy"`
}

// EndpointSpecBackend describe endpoint backend data
type EndpointSpecBackend struct {
	IP      string      `json:"ip"`
	PortMap map[int]int `json:"port_map"`
}

// SelfLink generates and returning link to object in platform
func (e *Endpoint) SelfLink() string {
	if e.Meta.SelfLink == "" {
		e.Meta.SelfLink = fmt.Sprintf("%s:%s:%s", e.Meta.Namespace, e.Meta.Service, e.Meta.Name)
	}
	return e.Meta.SelfLink
}


type EndpointCreateOptions struct {

}

type EndpointUpdateOptions struct {

}