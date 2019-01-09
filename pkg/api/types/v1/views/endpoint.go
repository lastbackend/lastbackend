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

type Endpoint struct {
	Meta   EndpointMeta   `json:"meta"`
	Spec   EndpointSpec   `json:"spec"`
	Status EndpointStatus `json:"status"`
}

type EndpointMeta struct {
	Name     string    `json:"name"`
	SelfLink string    `json:"self_link"`
	Updated  time.Time `json:"updated"`
	Created  time.Time `json:"created"`
}

type EndpointSpec struct {
	// Upstream state
	State string `json:"state"`

	IP       string               `json:"ip"`
	Domain   string               `json:"domain"`
	PortMap  map[uint16]string    `json:"port_map"`
	Strategy EndpointSpecStrategy `json:"strategy"`
	Policy   string               `json:"policy"`
}

type EndpointSpecStrategy struct {
	Route string `json:"route"`
	Bind  string `json:"bind"`
}

type EndpointStatus struct {
}

type EndpointList map[string]*Endpoint
