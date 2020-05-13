// +build linux
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

package runc

import (
	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/intelrdt"
	"github.com/opencontainers/runc/types"
)

func convertLibcontainerStats(ls *libcontainer.Stats) *types.Stats {
	cg := ls.CgroupStats
	if cg == nil {
		return nil
	}
	var s types.Stats
	s.Pids.Current = cg.PidsStats.Current
	s.Pids.Limit = cg.PidsStats.Limit

	s.CPU.Usage.Kernel = cg.CpuStats.CpuUsage.UsageInKernelmode
	s.CPU.Usage.User = cg.CpuStats.CpuUsage.UsageInUsermode
	s.CPU.Usage.Total = cg.CpuStats.CpuUsage.TotalUsage
	s.CPU.Usage.Percpu = cg.CpuStats.CpuUsage.PercpuUsage
	s.CPU.Throttling.Periods = cg.CpuStats.ThrottlingData.Periods
	s.CPU.Throttling.ThrottledPeriods = cg.CpuStats.ThrottlingData.ThrottledPeriods
	s.CPU.Throttling.ThrottledTime = cg.CpuStats.ThrottlingData.ThrottledTime

	s.Memory.Cache = cg.MemoryStats.Cache
	s.Memory.Kernel = convertMemoryEntry(cg.MemoryStats.KernelUsage)
	s.Memory.KernelTCP = convertMemoryEntry(cg.MemoryStats.KernelTCPUsage)
	s.Memory.Swap = convertMemoryEntry(cg.MemoryStats.SwapUsage)
	s.Memory.Usage = convertMemoryEntry(cg.MemoryStats.Usage)
	s.Memory.Raw = cg.MemoryStats.Stats

	s.Blkio.IoServiceBytesRecursive = convertBlkioEntry(cg.BlkioStats.IoServiceBytesRecursive)
	s.Blkio.IoServicedRecursive = convertBlkioEntry(cg.BlkioStats.IoServicedRecursive)
	s.Blkio.IoQueuedRecursive = convertBlkioEntry(cg.BlkioStats.IoQueuedRecursive)
	s.Blkio.IoServiceTimeRecursive = convertBlkioEntry(cg.BlkioStats.IoServiceTimeRecursive)
	s.Blkio.IoWaitTimeRecursive = convertBlkioEntry(cg.BlkioStats.IoWaitTimeRecursive)
	s.Blkio.IoMergedRecursive = convertBlkioEntry(cg.BlkioStats.IoMergedRecursive)
	s.Blkio.IoTimeRecursive = convertBlkioEntry(cg.BlkioStats.IoTimeRecursive)
	s.Blkio.SectorsRecursive = convertBlkioEntry(cg.BlkioStats.SectorsRecursive)

	s.Hugetlb = make(map[string]types.Hugetlb)
	for k, v := range cg.HugetlbStats {
		s.Hugetlb[k] = convertHugtlb(v)
	}

	if is := ls.IntelRdtStats; is != nil {
		if intelrdt.IsCatEnabled() {
			s.IntelRdt.L3CacheInfo = convertL3CacheInfo(is.L3CacheInfo)
			s.IntelRdt.L3CacheSchemaRoot = is.L3CacheSchemaRoot
			s.IntelRdt.L3CacheSchema = is.L3CacheSchema
		}
		if intelrdt.IsMbaEnabled() {
			s.IntelRdt.MemBwInfo = convertMemBwInfo(is.MemBwInfo)
			s.IntelRdt.MemBwSchemaRoot = is.MemBwSchemaRoot
			s.IntelRdt.MemBwSchema = is.MemBwSchema
		}
	}

	s.NetworkInterfaces = ls.Interfaces
	return &s
}

func convertHugtlb(c cgroups.HugetlbStats) types.Hugetlb {
	return types.Hugetlb{
		Usage:   c.Usage,
		Max:     c.MaxUsage,
		Failcnt: c.Failcnt,
	}
}

func convertMemoryEntry(c cgroups.MemoryData) types.MemoryEntry {
	return types.MemoryEntry{
		Limit:   c.Limit,
		Usage:   c.Usage,
		Max:     c.MaxUsage,
		Failcnt: c.Failcnt,
	}
}

func convertBlkioEntry(c []cgroups.BlkioStatEntry) []types.BlkioEntry {
	var out []types.BlkioEntry
	for _, e := range c {
		out = append(out, types.BlkioEntry{
			Major: e.Major,
			Minor: e.Minor,
			Op:    e.Op,
			Value: e.Value,
		})
	}
	return out
}

func convertL3CacheInfo(i *intelrdt.L3CacheInfo) *types.L3CacheInfo {
	return &types.L3CacheInfo{
		CbmMask:    i.CbmMask,
		MinCbmBits: i.MinCbmBits,
		NumClosids: i.NumClosids,
	}
}

func convertMemBwInfo(i *intelrdt.MemBwInfo) *types.MemBwInfo {
	return &types.MemBwInfo{
		BandwidthGran: i.BandwidthGran,
		DelayLinear:   i.DelayLinear,
		MinBandwidth:  i.MinBandwidth,
		NumClosids:    i.NumClosids,
	}
}
