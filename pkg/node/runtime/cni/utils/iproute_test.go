//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseBridgeFDB(t *testing.T) {

	var assets = make(map[string]FDBRule)

	assets["ee:72:88:1e:a1:f4 dev lb.1 dst 10.220.7.12 self permanent"] = FDBRule{
		Mac:       "ee:72:88:1e:a1:f4",
		Device:    "lb.1",
		DST:       "10.220.7.12",
		Self:      true,
		Permanent: true,
	}

	assets["ee:72:88:1e:a1:f5 dev lb.1 dst 10.220.7.12 self permanent"] = FDBRule{
		Mac:       "ee:72:88:1e:a1:f5",
		Device:    "lb.1",
		DST:       "10.220.7.12",
		Self:      true,
		Permanent: true,
	}

	assets["33:33:00:00:00:01 dev eth0 self permanent"] = FDBRule{
		Mac:       "33:33:00:00:00:01",
		Device:    "eth0",
		Self:      true,
		Permanent: true,
	}

	assets["01:00:5e:00:00:01 dev eth0 self permanent"] = FDBRule{
		Mac:       "01:00:5e:00:00:01",
		Device:    "eth0",
		Self:      true,
		Permanent: true,
	}

	assets["33:33:00:00:00:01 dev docker0 self permanent"] = FDBRule{
		Mac:       "33:33:00:00:00:01",
		Device:    "docker0",
		Self:      true,
		Permanent: true,
	}

	assets["02:42:91:fb:3f:70 dev docker0 master docker0 permanent"] = FDBRule{
		Mac:       "02:42:91:fb:3f:70",
		Device:    "docker0",
		Master:    "docker0",
		Permanent: true,
	}

	assets["02:42:91:fb:3f:70 dev docker0 vlan 1 master docker0 permanent"] = FDBRule{
		Mac:       "02:42:91:fb:3f:70",
		Vlan:      "1",
		Device:    "docker0",
		Master:    "docker0",
		Permanent: true,
	}

	for fdbr, valid := range assets {
		i := BridgeFDBParse(fdbr)
		assert.Equal(t, valid.Mac, i.Mac, "mac")
		assert.Equal(t, valid.Device, i.Device, "device")
		assert.Equal(t, valid.DST, i.DST, "dst")
		assert.Equal(t, valid.Vlan, i.Vlan, "vlan")
		assert.Equal(t, valid.Master, i.Master, "master")
		assert.Equal(t, valid.Self, i.Self, "self")
		assert.Equal(t, valid.Permanent, i.Permanent, "permanent")
	}

}
