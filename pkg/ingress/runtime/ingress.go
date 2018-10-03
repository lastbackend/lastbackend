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

package runtime

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/util/system"
	"github.com/spf13/viper"
	"os"
)

func IngressInfo() types.IngressInfo {

	var (
		info = types.IngressInfo{}
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

func IngressStatus() types.IngressStatus {

	var state = types.IngressStatus{}
	return state
}

