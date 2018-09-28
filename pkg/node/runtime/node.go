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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/node/envs"
	"github.com/lastbackend/lastbackend/pkg/util/system"
	"github.com/shirou/gopsutil/mem"
	"github.com/spf13/viper"

	"fmt"
	"os"
)

const MinContainerMemory = 32

func NodeInfo() types.NodeInfo {

	var (
		info = types.NodeInfo{}
	)

	osInfo := system.GetOsInfo()
	hostname, err := os.Hostname()
	if err != nil {
		_ = fmt.Errorf("get hostname err: %s", err)
	}

	link := viper.GetString("runtime.interface")
	ip, err := system.GetHostIP(link)
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

func NodeStatus() types.NodeStatus {

	var state = types.NodeStatus{}

	state.Capacity = NodeCapacity()
	state.Allocated = NodeAllocation()

	//state.Services.Router.Enabled = viper.GetBool("node.services.router.enabled")
	//state.Services.Router.ExternalIP = viper.GetString("node.services.router.external_ip")
	//state.Services.Builder = viper.GetBool("node.services.builder")

	return state
}

func NodeCapacity() types.NodeResources {

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Errorf("get memory err: %s", err)
	}

	m := vmStat.Total / 1024 / 1024

	return types.NodeResources{
		Memory:     int64(m),
		Pods:       int(m / MinContainerMemory),
		Containers: int(m / MinContainerMemory),
	}
}

func NodeAllocation() types.NodeResources {

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Errorf("get memory err: %s", err)
	}

	m := vmStat.Free / 1024 / 1024
	s := envs.Get().GetState().Pods()

	return types.NodeResources{
		Memory:     int64(m),
		Pods:       s.GetPodsCount(),
		Containers: s.GetContainersCount(),
	}
}
