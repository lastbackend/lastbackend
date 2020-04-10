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

package models_test

import (
	"testing"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {

	var network = &models.SubnetSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: models.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	var asset = &models.SubnetSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: models.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assert.Equal(t, true, models.SubnetSpecEqual(network, asset), "equal")
}

func TestNotEqual(t *testing.T) {

	var network = &models.SubnetSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: models.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	var assets = make(map[string]*models.SubnetSpec)

	assets["type"] = &models.SubnetSpec{
		Type: "vlan",
		CIDR: "10.0.0.0/24",
		IFace: models.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["subnet"] = &models.SubnetSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/22",
		IFace: models.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["subnet & type"] = &models.SubnetSpec{
		Type: "vlan",
		CIDR: "10.0.0.0/22",
		IFace: models.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface index"] = &models.SubnetSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: models.NetworkInterface{
			Index: 2,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface index & type"] = &models.SubnetSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: models.NetworkInterface{
			Index: 2,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface index & subnet"] = &models.SubnetSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: models.NetworkInterface{
			Index: 2,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface name"] = &models.SubnetSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: models.NetworkInterface{
			Index: 1,
			Name:  "lb.0",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface addr"] = &models.SubnetSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: models.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.2",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["iface haddr"] = &models.SubnetSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: models.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:c8:fe",
		},
		Addr: "10.0.0.0",
	}

	assets["addr"] = &models.SubnetSpec{
		Type: "vxlan",
		CIDR: "10.0.0.0/24",
		IFace: models.NetworkInterface{
			Index: 1,
			Name:  "lb.1",
			Addr:  "10.0.0.1",
			HAddr: "b6:3c:b9:62:e8:fe",
		},
		Addr: "10.0.0.1",
	}

	for attr, asset := range assets {
		assert.Equal(t, false, models.SubnetSpecEqual(network, asset), attr)
	}
}
