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
	"strings"
)

const NetworkTypeVxLAN = "vxlan"

type Network struct {
	System
	Meta   Meta          `json:"meta"`
	Status NetworkStatus `json:"status"`
	Spec   NetworkSpec   `json:"spec"`
}

type NetworkStatus struct {
	State string `json:"state"`
}

type NetworkSpec struct {
}

type NetworkState struct {
	SubnetSpec
}

type NodeMode struct {
	Ingress bool `json:"ingress"`
}

// swagger:model types_network_interface
type NetworkInterface struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
	Addr  string `json:"addr"`
	HAddr string `json:"HAddr"`
}

type Subnet struct {
	System
	Meta   SubnetMeta   `json:"meta"`
	Status SubnetStatus `json:"status"`
	Spec   SubnetSpec   `json:"spec"`
}

type SubnetStatus struct {
	State string `json:"state"`
}

type SubnetMeta struct {
	Meta
	Node string `json:"node"`
}

// swagger:model subnet_spec
type SubnetSpec struct {
	// Subnet state
	State string `json:"state"`
	// Node network type
	Type string `json:"type"`
	// Node Subnet subnet info
	CIDR string `json:"cidr"`
	// Node Subnet interface
	IFace NetworkInterface `json:"iface"`
	// Node Public IP
	Addr string `json:"addr"`
	// Node Internal IP
	IP string `json:"ip"`
}

func SubnetGetNameFromCIDR(CIDR string) string {
	return strings.Replace(CIDR, "/", ":", -1)
}

func SubnetSpecEqual(n *SubnetSpec, nt *SubnetSpec) bool {

	switch false {
	case n.Type == nt.Type:
		return false
	case n.CIDR == nt.CIDR:
		return false
	case n.IFace.Index == nt.IFace.Index:
		return false
	case n.IFace.Name == nt.IFace.Name:
		return false
	case n.IFace.Addr == nt.IFace.Addr:
		return false
	case n.IFace.HAddr == nt.IFace.HAddr:
		return false
	case n.Addr == nt.Addr:
		return false
	}
	return true
}

func (s Subnet) SelfLink() string {
	if s.Meta.SelfLink == EmptyString {
		s.Meta.SelfLink = s.Meta.Name
	}
	return s.Meta.SelfLink
}
