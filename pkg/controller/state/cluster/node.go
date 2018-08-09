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
				"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"context"
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	)

type NodeLease struct {
	done    chan bool
	Request  NodeLeaseOptions
	Response struct {
		Err  error
		Node *types.Node
	}
}

type NodeLeaseOptions struct {
	Node    *string
	Memory  *int64
	Storage *int64
}

func (nl *NodeLease) Wait() {
	<- nl.done
}

func handleNodeLease (cs *ClusterState, nl *NodeLease) error {



	defer func() {
		nl.done <- true
	}()

	for _, n := range cs.node.list {

		if (n.Status.Capacity.Memory-n.Status.Allocated.Memory) > *nl.Request.Memory {

			n.Status.Allocated.Pods++
			n.Status.Allocated.Memory += *nl.Request.Memory

			nm := distribution.NewNodeModel(context.Background(), envs.Get().GetStorage())
			nm.SetStatus(n, n.Status)

			nl.Response.Node = n
			return nil
		}

	}

	return nil
}

func handleNodeRelease (cs *ClusterState, nl *NodeLease) error {

	defer func() {
		nl.done <- true
	}()

	if _, ok := cs.node.list[*nl.Request.Node]; !ok {
		return nil
	}

	node := cs.node.list[*nl.Request.Node]
	node.Status.Allocated.Pods--
	node.Status.Allocated.Memory-=*nl.Request.Memory

	nm := distribution.NewNodeModel(context.Background(), envs.Get().GetStorage())
	nm.SetStatus(node, node.Status)

	return nil
}