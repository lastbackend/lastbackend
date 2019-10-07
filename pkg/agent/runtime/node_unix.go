// +build linux darwin

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
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/agent/envs"
	"github.com/lastbackend/lastbackend/pkg/util/system"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

const MinContainerMemory uint64 = 32

func NodeInfo() types.NodeInfo {

	var (
		info = types.NodeInfo{}
	)

	osInfo := system.GetOsInfo()
	hostname, err := os.Hostname()
	if err != nil {
		_ = fmt.Errorf("get hostname err: %s", err)
	}

	ip, err := system.GetHostIP(envs.Get().GetConfig().Network.Interface)
	if err != nil {
		_ = fmt.Errorf("get ip err: %s", err)
	}

	info.Hostname = hostname
	info.ExternalIP = ip
	info.OSType = osInfo.GoOS
	info.OSName = fmt.Sprintf("%s %s", osInfo.OS, osInfo.Core)
	info.Architecture = osInfo.Platform

	net := envs.Get().GetNet()
	if net != nil {
		nt := net.Info(context.Background())
		info.InternalIP = nt.IP
		info.CIDR = nt.CIDR
	}

	return info
}

func NodeStatus() types.NodeStatus {

	var state = types.NodeStatus{}

	state.Capacity = NodeCapacity()
	state.Allocated = NodeAllocation()

	return state
}

func NodeCapacity() types.NodeResources {

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		_ = fmt.Errorf("get memory err: %s", err)
	}

	cpuStat, err := cpu.Info()
	if err != nil {
		_ = fmt.Errorf("get cpu err: %s", err)
	}

	var stat syscall.Statfs_t

	syscall.Statfs(envs.Get().GetConfig().Container.Csi.Dir.Root, &stat)

	// Available blocks * size per block = available space in bytes
	storage := stat.Blocks * uint64(stat.Bsize)
	m := vmStat.Total

	return types.NodeResources{
		Storage:    int64(storage),
		RAM:        int64(m),
		CPU:        int64(cpuStat[0].Mhz) * int64(1e6) * int64(cpuStat[0].Cores),
		Pods:       int(m / (MinContainerMemory * 1024 * 1024)),
		Containers: int(m / (MinContainerMemory * 1024 * 1024)),
	}
}

func NodeAllocation() types.NodeResources {

	s := envs.Get().GetState().Pods()

	return types.NodeResources{
		RAM:        0,
		CPU:        0, // TODO: need get cpu resource value
		Pods:       s.GetPodsCount(),
		Containers: s.GetContainersCount(),
	}
}
