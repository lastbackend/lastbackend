//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package runtime

import (
	"fmt"
	"os"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/util/system"
	"github.com/spf13/viper"
)

func DiscoveryInfo() types.DiscoveryInfo {

	var (
		info = types.DiscoveryInfo{}
	)

	osInfo := system.GetOsInfo()
	hostname, err := os.Hostname()
	if err != nil {
		_ = fmt.Errorf("get hostname err: %s", err)
	}

	iface := viper.GetString("runtime.interface")
	ip, err := system.GetHostIP(iface)
	if err != nil {
		_ = fmt.Errorf("get ip err: %s", err)
	}

	info.Hostname = hostname
	info.InternalIP = ip
	info.OSType = osInfo.GoOS
	info.OSName = fmt.Sprintf("%s %s", osInfo.OS, osInfo.Core)
	info.Architecture = osInfo.Platform

	return info
}

func DiscoveryStatus() types.DiscoveryStatus {

	var state = types.DiscoveryStatus{}

	iface := viper.GetString("runtime.interface")
	ip, err := system.GetHostIP(iface)
	if err != nil {
		_ = fmt.Errorf("get ip err: %s", err)
	}

	state.Port = uint16(viper.GetInt("dns.port"))
	state.IP = ip

	return state
}
