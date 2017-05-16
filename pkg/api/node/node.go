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
