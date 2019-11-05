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
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/internal/util/proxy"
	"github.com/lastbackend/lastbackend/internal/util/system"
	"os"
)

func (r *Runtime) ExporterInfo() types.ExporterInfo {

	var (
		info = types.ExporterInfo{}
	)

	osInfo := system.GetOsInfo()
	hostname, err := os.Hostname()
	if err != nil {
		_ = fmt.Errorf("get hostname err: %s", err)
	}

	ip, err := system.GetHostIP(r.iface)
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

func (r *Runtime) ExporterStatus() types.ExporterStatus {

	var state = types.ExporterStatus{}

	ip, err := system.GetHostIP(r.iface)
	if err != nil {
		_ = fmt.Errorf("get ip err: %s", err)
	}

	hp := r.port
	if hp == 0 {
		hp = proxy.DefaultPort
	}

	state.Ready = true
	state.Listener.IP = ip
	state.Http.IP = ip
	state.Http.Port = hp

	if r.logger != nil {
		lp := r.logger.GetPort()
		if lp == 0 {
			lp = proxy.DefaultPort
		}
		state.Listener.Port = lp
	}

	return state
}
