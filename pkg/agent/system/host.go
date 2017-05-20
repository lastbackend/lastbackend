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

package system

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/agent/config"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/util/system"
	"github.com/shirou/gopsutil/mem"
	"time"
)

const MinContainerMemory = 32

func GetNodeMeta() types.NodeMeta {
	var cfg = config.Get()
	var meta = types.NodeMeta{}

	meta.Created = time.Now()
	meta.Updated = time.Now()

	info := system.GetOsInfo()

	if cfg.Host.Hostname != nil && *cfg.Host.Hostname != "" {
		meta.Hostname = *cfg.Host.Hostname
	} else {
		meta.Hostname = info.Hostname
	}

	meta.Port = *cfg.AgentServer.Port
	if ip, err := system.GetExternalIP(); err != nil {
		fmt.Errorf("Get external ip err: %s", err.Error())
	} else {
		meta.IP = ip
	}

	meta.OSType = info.GoOS
	meta.OSName = fmt.Sprintf("%s %s", info.OS, info.Core)
	meta.Architecture = info.Platform

	return meta
}

func GetNodeState() types.NodeState {
	var state = types.NodeState{}

	state.Capacity = GetNodeCapacity()
	state.Allocated = GetNodeAllocation()

	return state
}

func GetNodeCapacity() types.NodeResources {

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Errorf("Get memory err: %s", err.Error())
	}
	m := vmStat.Total / 1024 / 1024
	capacity := types.NodeResources{
		Memory:     int64(m),
		Pods:       int(m / MinContainerMemory),
		Containers: int(m / MinContainerMemory),
	}
	return capacity
}

func GetNodeAllocation() types.NodeResources {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Errorf("Get memory err: %s", err.Error())
	}
	m := vmStat.Free / 1024 / 1024
	s := context.Get().GetCache().Pods()
	allocation := types.NodeResources{
		Memory:     int64(m),
		Pods:       s.GetPodsCount(),
		Containers: s.GetContainersCount(),
	}
	return allocation
}
