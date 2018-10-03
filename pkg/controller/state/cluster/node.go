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

package cluster

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type NodeLease struct {
	done     chan bool
	sync     bool
	Request  NodeLeaseOptions
	Response struct {
		Err  error
		Node *types.Node
	}
}

type NodeLeaseOptions struct {
	Node     *string
	Memory   *int64
	Storage  *int64
	Selector map[string]string
}

func (nl *NodeLease) Wait() {
	<-nl.done
}

func handleNodeLease(cs *ClusterState, nl *NodeLease) error {

	defer func() {
		if !nl.sync {
			nl.done <- true
		}
	}()

	for _, n := range cs.node.list {


		var (
			node *types.Node
			allocated  = new(types.NodeResources)
		)

		if nl.Request.Memory != nil {
			if (n.Status.Capacity.Memory - n.Status.Allocated.Memory) > *nl.Request.Memory {
				node = n
				allocated.Pods++
				allocated.Memory += *nl.Request.Memory
			}
		}

		if nl.Request.Storage != nil {
			if (n.Status.Capacity.Storage - n.Status.Allocated.Storage) > *nl.Request.Storage {
				if node == nil {
					node = n
				}
				allocated.Storage += *nl.Request.Storage
			}
		}


		if node != nil {

			node.Status.Allocated.Pods += allocated.Pods
			node.Status.Allocated.Memory += allocated.Memory
			node.Status.Allocated.Storage += allocated.Storage

			nm := distribution.NewNodeModel(context.Background(), envs.Get().GetStorage())
			nm.Set(n)

			nl.Response.Node = n
			return nil
		}


	}

	return nil
}

func handleNodeRelease(cs *ClusterState, nl *NodeLease) error {

	defer func() {
		if !nl.sync {
			nl.done <- true
		}
	}()

	if _, ok := cs.node.list[*nl.Request.Node]; !ok {
		return nil
	}

	n := cs.node.list[*nl.Request.Node]

	if nl.Request.Memory != nil {
		n.Status.Allocated.Pods--
		n.Status.Allocated.Memory -= *nl.Request.Memory
	}

	if nl.Request.Storage != nil {
		n.Status.Allocated.Storage -= *nl.Request.Storage
	}

	nm := distribution.NewNodeModel(context.Background(), envs.Get().GetStorage())
	nm.Set(n)

	return nil
}
