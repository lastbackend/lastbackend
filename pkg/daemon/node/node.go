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
