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

const (
	// EndpointSpecRouteStrategyRR - round robin balancing strategy type
	EndpointSpecRouteStrategyRR = "rr"
	// EndpointSpecBindStrategyDefault - default scheduling endpoint across all nodes
	EndpointSpecBindStrategyDefault = "default"
)

// swagger:ignore
// Upstream - service endpoint
type Endpoint struct {
	System
	Meta   EndpointMeta   `json:"meta"`
	Status EndpointStatus `json:"status"`
	Spec   EndpointSpec   `json:"spec"`
}

type EndpointList struct {
	System
	Items []*Endpoint
}

type EndpointMap struct {
	System
	Items map[string]*Endpoint
}

// swagger:ignore
// EndpointMeta - endpoint meta data
type EndpointMeta struct {
	Meta
	// Environment name
	Namespace string `json:"namespace"`
	// SelfLink
	SelfLink EndpointSelfLink `json:"self_link"`
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
	// Upstream state
	State     string               `json:"state"`
	External  bool                 `json:"external"`
	IP        string               `json:"ip"`
	Domain    string               `json:"domain"`
	PortMap   map[uint16]string    `json:"port_map"`
	Strategy  EndpointSpecStrategy `json:"strategy"`
	Policy    string               `json:"policy"`
	Upstreams []string             `json:"upstreams"`
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
func (e *Endpoint) SelfLink() *EndpointSelfLink {
	return &e.Meta.SelfLink
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

func NewEndpointList() *EndpointList {
	dm := new(EndpointList)
	dm.Items = make([]*Endpoint, 0)
	return dm
}

func NewEndpointMap() *EndpointMap {
	dm := new(EndpointMap)
	dm.Items = make(map[string]*Endpoint)
	return dm
}
