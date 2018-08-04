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

// swagger:ignore
// Endpoint - service endpoint
type Endpoint struct {
	Runtime
	Meta   EndpointMeta   `json:"meta"`
	Status EndpointStatus `json:"status"`
	Spec   EndpointSpec   `json:"spec"`
}

type EndpointList struct {
	Runtime
	Items []*Endpoint
}

type EndpointMap struct {
	Runtime
	Items map[string]*Endpoint
}


// swagger:ignore
// EndpointMeta - endpoint meta data
type EndpointMeta struct {
	Meta
	// Namespace name
	Namespace string `json:"namespace"`
}

// swagger:ignore
// EndpointStatus - endpoint status
type EndpointStatus struct {
	State string          `json:"state"`
	Ready map[string]bool `json:"ready"`
}

// EndpointSpec spec data
// swagger:model types_endpoint_spec
type EndpointSpec struct {
	// Endpoint state
	State string `json:"state"`

	IP        string               `json:"ip"`
	Domain    string               `json:"domain"`
	PortMap   map[uint16]string    `json:"port_map"`
	Upstreams []string             `json:"upstreams"`
	Strategy  EndpointSpecStrategy `json:"strategy"`
	Policy    string               `json:"policy"`
}

type EndpointState struct {
	EndpointSpec
}

// EndpointSpecStrategy describes route and bind
// swagger:model types_endpoint_spec_strategy
type EndpointSpecStrategy struct {
	Route string `json:"route"`
	Bind  string `json:"bind"`
}

// swagger:ignore
// EndpointUpstream describe endpoint backend data
type EndpointUpstream struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// swagger:ignore
// SelfLink generates and returning link to object in platform
func (e *Endpoint) SelfLink() string {
	if e.Meta.SelfLink == "" {
		e.Meta.SelfLink = e.CreateSelfLink(e.Meta.Namespace, e.Meta.Name)
	}
	return e.Meta.SelfLink
}

func (e *Endpoint) CreateSelfLink(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}


// swagger:ignore
type EndpointCreateOptions struct {
	IP            string            `json:"ip"`
	Domain        string            `json:"domain"`
	Ports         map[uint16]string `json:"ports"`
	RouteStrategy string            `json:"route_strategy"`
	Policy        string            `json:"policy"`
	BindStrategy  string            `json:"bind_strategy"`
}

// swagger:ignore
type EndpointUpdateOptions struct {
	IP            *string           `json:"ip"`
	Ports         map[uint16]string `json:"ports"`
	RouteStrategy string            `json:"route_strategy"`
	Policy        string            `json:"policy"`
	BindStrategy  string            `json:"bind_strategy"`
}

func NewEndpointList () *EndpointList {
	dm := new(EndpointList)
	dm.Items = make([]*Endpoint, 0)
	return dm
}

func NewEndpointMap () *EndpointMap {
	dm := new(EndpointMap)
	dm.Items = make(map[string]*Endpoint)
	return dm
}
