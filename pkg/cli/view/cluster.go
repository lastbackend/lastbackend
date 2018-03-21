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

package view

import (
	"time"

	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/util/table"
)

type ClusterList []*Cluster
type Cluster struct {
	Meta  ClusterMeta  `json:"meta"`
	State ClusterState `json:"state"`
}

type ClusterMeta struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Region      string            `json:"region"`
	Provider    string            `json:"provider"`
	Labels      map[string]string `json:"labels"`
	Created     time.Time         `json:"created"`
	Updated     time.Time         `json:"updated"`
}

type ClusterState struct {
	Nodes struct {
		Total   int `json:"total"`
		Online  int `json:"online"`
		Offline int `json:"offline"`
	} `json:"nodes"`
	Capacity  ClusterResources `json:"capacity"`
	Allocated ClusterResources `json:"allocated"`
	Deleted   bool             `json:"deleted"`
}

type ClusterResources struct {
	Containers int   `json:"containers"`
	Pods       int   `json:"pods"`
	Memory     int64 `json:"memory"`
	Cpu        int   `json:"cpu"`
	Storage    int   `json:"storage"`
}

func (c *Cluster) Print() {

	println()
	table.PrintHorizontal(map[string]interface{}{
		"NAME":        c.Meta.Name,
		"DESCRIPTION": c.Meta.Description,
		"REGION":      c.Meta.Region,
		"PROVIDER":    c.Meta.Provider,
	})
	println()
}

func FromApiClusterView(cluster *views.Cluster) *Cluster {
	var item = new(Cluster)
	item.Meta.Name = cluster.Meta.Name
	item.Meta.Description = cluster.Meta.Description
	item.Meta.Region = cluster.Meta.Region
	item.Meta.Provider = cluster.Meta.Provider
	item.Meta.Created = cluster.Meta.Created
	item.Meta.Updated = cluster.Meta.Updated

	item.Meta.Labels = cluster.Meta.Labels
	if item.Meta.Labels == nil {
		item.Meta.Labels = make(map[string]string, 0)
	}

	item.State.Nodes.Total = cluster.Status.Nodes.Total
	item.State.Nodes.Online = cluster.Status.Nodes.Online
	item.State.Nodes.Offline = cluster.Status.Nodes.Offline
	item.State.Capacity.Containers = cluster.Status.Capacity.Containers
	item.State.Capacity.Pods = cluster.Status.Capacity.Pods
	item.State.Capacity.Memory = cluster.Status.Capacity.Memory
	item.State.Capacity.Cpu = cluster.Status.Capacity.Cpu
	item.State.Capacity.Storage = cluster.Status.Capacity.Storage
	item.State.Allocated.Containers = cluster.Status.Allocated.Containers
	item.State.Allocated.Pods = cluster.Status.Allocated.Pods
	item.State.Allocated.Memory = cluster.Status.Allocated.Memory
	item.State.Allocated.Cpu = cluster.Status.Allocated.Cpu
	item.State.Allocated.Storage = cluster.Status.Allocated.Storage
	item.State.Deleted = cluster.Status.Deleted

	return item
}
