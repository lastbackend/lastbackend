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
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logLevel = 3

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

func (n *node) Get(id string) (*types.Node, error) {
	var (
		storage = ctx.Get().GetStorage()
	)

	log.V(logLevel).Debugf("Node: get node by id %s", id)

	node, err := storage.Node().Get(n.Context, id)
	if err != nil {
		log.V(logLevel).Debugf("Node: get node `%s` err: %s", id, err.Error())
		return nil, err
	}

	return node, nil
}

func (n *node) Update(node *types.Node) error {
	var (
		storage = ctx.Get().GetStorage()
	)
	return storage.Node().Update(n.Context, node)
}

func (n *node) Create(meta *types.NodeMeta, state *types.NodeState) (*types.Node, error) {

	var (
		storage = ctx.Get().GetStorage()
		node    = new(types.Node)
	)

	log.V(logLevel).Debugf("Node: create node with meta: %#v, state: %#v", meta, state)

	node.Meta = *meta
	node.State = *state

	if err := storage.Node().Insert(n.Context, node); err != nil {
		log.V(logLevel).Debugf("Node: insert node err: %s", err.Error())
		return nil, err
	}

	return node, nil
}
