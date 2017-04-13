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

package node

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	ctx "github.com/lastbackend/lastbackend/pkg/daemon/context"
)

type Node struct {
}

func New() *Node {
	return new(Node)
}

func (n *Node) List(c context.Context) (*types.NodeList, error) {
	var storage = ctx.Get().GetStorage()
	return storage.Node().List(c)
}

func (n *Node) Get(c context.Context, hostname string) (*types.Node, error) {
	var (
		storage = ctx.Get().GetStorage()
	)

	return storage.Node().Get(c, hostname)
}

func (n *Node) SetMeta(c context.Context, node *types.Node) error {
	var (
		storage = ctx.Get().GetStorage()
	)

	return storage.Node().UpdateMeta(c, &node.Meta)
}

func (n *Node) SetState(c context.Context, node *types.Node) error {
	var (
		storage = ctx.Get().GetStorage()
	)

	return storage.Node().UpdateState(c, &node.Meta, &node.State)
}

func (n *Node) Create(c context.Context, meta *types.NodeMeta, state *types.NodeState) (*types.Node, error) {

	var (
		storage = ctx.Get().GetStorage()
		node    = new(types.Node)
		log     = ctx.Get().GetLogger()
	)

	log.Debug("Create new Node")

	node.Meta = *meta
	node.State = *state

	return storage.Node().Insert(c, &node.Meta, &node.State)
}

func (n *Node) PodSpecRemove (c context.Context, hostname string, spec *types.PodNodeSpec) error {

	var (
		storage = ctx.Get().GetStorage()
		log     = ctx.Get().GetLogger()
	)

	node, err := n.Get(c, hostname)
	if err !=nil {
		log.Errorf("Node: Pod spec remove: remove pod spec err: %s", err.Error())
		return err
	}

	if node == nil {
		log.Debug("Node: Pod spec remove: node not found")
		return nil
	}

	log.Debug("Remove pod spec from node")
	if err := storage.Node().RemovePod(c, &node.Meta, spec); err != nil {
		log.Errorf("Node: Pod spec remove: remove pod spec err: %s", err.Error())
		return err
	}

	// Update pod node spec

	return nil
}

func (n *Node) PodSpecUpdate (c context.Context, hostname string, spec *types.PodNodeSpec) error {
	// Get node by hostname
	// Update pod node spec
	var (
		storage = ctx.Get().GetStorage()
		log     = ctx.Get().GetLogger()
	)

	node, err := n.Get(c, hostname)
	if err !=nil {
		log.Errorf("Node: Pod spec remove: remove pod spec err: %s", err.Error())
		return err
	}

	if node == nil {
		log.Debug("Node: Pod spec remove: node not found")
		return nil
	}

	log.Debug("Remove pod spec from node")
	if err := storage.Node().UpdatePod(c, &node.Meta, spec); err != nil {
		log.Errorf("Node: Pod spec remove: remove pod spec err: %s", err.Error())
		return err
	}

	// Update pod node spec

	return nil
}

func (n *Node) Allocate(c context.Context, spec *types.PodNodeSpec) error {

	var (
		node    types.Node
		storage = ctx.Get().GetStorage()
		log     = ctx.Get().GetLogger()
		memory  = int64(0)
	)

	log.Debug("Allocate Pod to Node")

	nodes, err := storage.Node().List(c)
	if err != nil {
		log.Errorf("Node: allocate: get nodes error: %s", err.Error())
		return err
	}

	for _, c := range spec.Spec.Containers {
		memory += c.Quota.Memory
	}

	for _, node = range *nodes {
		if node.State.Capacity.Memory > memory {
			break
		}
	}

	node.Spec.Pods = append(node.Spec.Pods, spec)

	return nil
}