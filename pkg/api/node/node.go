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
	"errors"
	ctx "github.com/lastbackend/lastbackend/pkg/api/context"
	"github.com/lastbackend/lastbackend/pkg/common/types"
)

type node struct {
	Context context.Context
}

func New(ctx context.Context) *node {
	return &node{Context: ctx}
}

func (n *node) List() ([]*types.Node, error) {
	var storage = ctx.Get().GetStorage()
	return storage.Node().List(n.Context)
}

func (n *node) Get(hostname string) (*types.Node, error) {
	var (
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	log.Debug("Node: Get node info")
	node, err := storage.Node().Get(n.Context, hostname)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (n *node) SetMeta(node *types.Node) error {
	var (
		storage = ctx.Get().GetStorage()
	)

	return storage.Node().UpdateMeta(n.Context, node)
}

func (n *node) SetState(node *types.Node) error {
	var (
		storage = ctx.Get().GetStorage()
	)

	return storage.Node().UpdateState(n.Context, node)
}

func (n *node) Create(meta *types.NodeMeta, state *types.NodeState) (*types.Node, error) {

	var (
		storage = ctx.Get().GetStorage()
		node    = new(types.Node)
		log     = ctx.Get().GetLogger()
	)

	log.Debug("Create new Node")

	node.Meta = *meta
	node.State = *state

	if err := storage.Node().Insert(n.Context, node); err != nil {
		return node, err
	}

	return node, nil
}

func (n *node) PodSpecRemove(hostname string, spec *types.PodNodeSpec) error {

	var (
		storage = ctx.Get().GetStorage()
		log     = ctx.Get().GetLogger()
	)

	node, err := n.Get(hostname)
	if err != nil {
		log.Errorf("Node: Pod spec remove: remove pod spec err: %s", err.Error())
		return err
	}

	log.Debug("Remove pod spec from node")
	if err := storage.Node().RemovePod(n.Context, &node.Meta, spec); err != nil {
		log.Errorf("Node: Pod spec remove: remove pod spec err: %s", err.Error())
		return err
	}

	// Update pod node spec

	return nil
}

func (n *node) PodSpecUpdate(hostname string, spec *types.PodNodeSpec) error {
	// Get node by hostname
	// Update pod node spec
	var (
		storage = ctx.Get().GetStorage()
		log     = ctx.Get().GetLogger()
	)

	node, err := n.Get(hostname)
	if err != nil {
		log.Errorf("Node: Pod spec remove: remove pod spec err: %s", err.Error())
		return err
	}

	log.Debug("Remove pod spec from node")
	if err := storage.Node().UpdatePod(n.Context, &node.Meta, spec); err != nil {
		log.Errorf("Node: Pod spec remove: remove pod spec err: %s", err.Error())
		return err
	}

	// Update pod node spec

	return nil
}

func (n *node) Allocate(spec types.PodSpec) (*types.Node, error) {

	var (
		node    *types.Node
		storage = ctx.Get().GetStorage()
		log     = ctx.Get().GetLogger()
		memory  = int64(0)
	)

	log.Debug("Allocate Pod to Node")

	nodes, err := storage.Node().List(n.Context)
	if err != nil {
		log.Errorf("Node: allocate: get nodes error: %s", err.Error())
		return nil, err
	}

	for _, c := range spec.Containers {
		memory += c.Quota.Memory
	}

	for _, node = range nodes {
		log.Debugf("Node: Allocate: available memory %d", node.State.Capacity)
		if node.State.Capacity.Memory > memory {
			break
		}
	}

	if node == nil {
		log.Error("Node: Allocate: Available node not found")
		return nil, errors.New("Available node not found")
	}

	return node, nil
}
