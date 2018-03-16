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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEqual(t *testing.T) {

	var network = &Subnet{
		Type:   "vxlan",
		Subnet: "10.0.0.0/24",
		IFace: NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	var asset = &Subnet{
		Type:   "vxlan",
		Subnet: "10.0.0.0/24",
		IFace: NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assert.Equal(t, true, network.Equal(asset), "equal")
}

func TestNotEqual(t *testing.T) {

	var network = &Subnet{
		Type:   "vxlan",
		Subnet: "10.0.0.0/24",
		IFace: NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	var assets = make(map[string]*Subnet)

	assets["type"] = &Subnet{
		Type:   "vlan",
		Subnet: "10.0.0.0/24",
		IFace: NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["subnet"] = &Subnet{
		Type:   "vxlan",
		Subnet: "10.0.0.0/22",
		IFace: NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["subnet & type"] = &Subnet{
		Type:   "vlan",
		Subnet: "10.0.0.0/22",
		IFace: NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface index"] = &Subnet{
		Type:   "vxlan",
		Subnet: "10.0.0.0/24",
		IFace: NetworkInterface{
			Index: 2,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface index & type"] = &Subnet{
		Type:   "vxlan",
		Subnet: "10.0.0.0/24",
		IFace: NetworkInterface{
			Index: 2,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface index & subnet"] = &Subnet{
		Type:   "vxlan",
		Subnet: "10.0.0.0/24",
		IFace: NetworkInterface{
			Index: 2,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface name"] = &Subnet{
		Type:   "vxlan",
		Subnet: "10.0.0.0/24",
		IFace: NetworkInterface{
			Index: 1,
			Name:  "lb.0",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface addr"] = &Subnet{
		Type:   "vxlan",
		Subnet: "10.0.0.0/24",
		IFace: NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.2",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface haddr"] = &Subnet{
		Type:   "vxlan",
		Subnet: "10.0.0.0/24",
		IFace: NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:c8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["addr"] = &Subnet{
		Type:   "vxlan",
		Subnet: "10.0.0.0/24",
		IFace: NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.1",
	}

	for attr, asset := range assets {
		assert.Equal(t, false, network.Equal(asset), attr)
	}
}
