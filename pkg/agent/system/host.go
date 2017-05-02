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
	"github.com/lastbackend/lastbackend/pkg/agent/utils/system"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/shirou/gopsutil/mem"
	"time"
)

const MinContainerMemory = 32

func GetNodeMeta() types.NodeMeta {
	var cfg = config.Get().Host
	var meta = types.NodeMeta{}

	meta.Created = time.Now()
	meta.Updated = time.Now()

	if cfg.Hostname != "" {
		meta.Hostname = cfg.Hostname
	} else {
		meta.Hostname, _ = system.GetHostname()
	}

	meta.State.Capacity = GetNodeCapacity()
	meta.State.Allocated = GetNodeAllocation()

	return meta
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
	s := context.Get().GetStorage().Pods()
	allocation := types.NodeResources{
		Memory:     int64(m),
		Pods:       s.GetPodsCount(),
		Containers: s.GetContainersCount(),
	}
	return allocation
}
