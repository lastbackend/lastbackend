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

package node

import "github.com/lastbackend/lastbackend/pkg/distribution/types"

func getClusterStatus(nodes map[string]*types.Node) *types.ClusterStatus {
	status := new(types.ClusterStatus)
	status.Nodes.Total = len(nodes)

	for _, node := range nodes {
		status.Allocated.Containers += node.Status.Allocated.Containers
		status.Allocated.Pods += node.Status.Allocated.Pods
		status.Allocated.Memory += node.Status.Allocated.Memory
		status.Allocated.Cpu += node.Status.Allocated.Cpu
		status.Allocated.Storage += node.Status.Allocated.Storage

		status.Capacity.Containers += node.Status.Capacity.Containers
		status.Capacity.Pods += node.Status.Capacity.Pods
		status.Capacity.Memory += node.Status.Capacity.Memory
		status.Capacity.Cpu += node.Status.Capacity.Cpu
		status.Capacity.Storage += node.Status.Capacity.Storage

		if node.Online {
			status.Nodes.Online++
		} else {
			status.Nodes.Offline++
		}
	}

	return status
}
