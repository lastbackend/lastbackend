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

package distribution

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/util/generator"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

type INode interface {
	List() ([]*types.Node, error)
	Create() (*types.Node, error)
	Get(name string) (*types.Node, error)
	Update(node *types.Node, opts *types.NodeUpdateOptions) error
	SetAvailable(name string) error
	SetUnavailable(name string) error
	SetState(node *types.Node, state types.NodeState) (*types.Node, error)
	SetInfo(node *types.Node, info types.NodeInfo) error
	SetNetwork(node *types.Node, network types.Subnet) error
	GetSpec(name string) (types.NodeSpec, error)
	Remove(name string) error
}

type Node struct {
	context context.Context
	storage storage.Storage
}

func (n *Node) List() ([]*types.Node, error) {
	return n.storage.Node().List(n.context)
}

func (n *Node) Create() (*types.Node, error) {

	log.Debug("Node: create node in cluster")

	ni := new(types.Node)
	ni.Meta.Labels = make(map[string]string, 0)
	ni.Meta.Token = generator.GenerateRandomString(32)
	ni.State.Online = true

	if err := n.storage.Node().Insert(n.context, ni); err != nil {
		log.Debugf("Node: insert node err: %s", err)
		return nil, err
	}

	return ni, nil
}

func (n *Node) Get(name string) (*types.Node, error) {

	log.V(logLevel).Debugf("Node: get Node by name %s", name)

	node, err := n.storage.Node().Get(n.context, name)
	if err != nil {
		log.V(logLevel).Debugf("Node: get Node `%s` err: %s", name, err)
		return nil, err
	}

	return node, nil
}

func (n *Node) Update(node *types.Node, opts *types.NodeUpdateOptions) error {

	var (
		err error
	)

	log.V(logLevel).Debugf("Node: update Node %#v", node)

	if opts.Description != nil {
		log.V(logLevel).Debug("Node: update Node meta")
		node.Meta.Description = *opts.Description
	}

	if opts.ExternalIP != nil {
		log.V(logLevel).Debug("Node: update Node external ip")
		node.Info.ExternalIP = opts.ExternalIP.IP
	}

	if err = n.storage.Node().Update(n.context, node); err != nil {
		log.V(logLevel).Errorf("Node: update Node meta err: %s", err)
		return err
	}
	return nil
}

func (n *Node) SetAvailable(name string) error {

	node, err := n.storage.Node().Get(n.context, name)
	if err != nil {
		log.Errorf("Get node id by name error: %s", err)
		return err
	}

	if node == nil {
		log.Debugf("Node not found by provided name: %s", name)
		return nil
	}

	if err := n.storage.Node().SetAvailable(n.context, node); err != nil {
		log.Errorf("Set node online state error: %s", err)
		return err
	}

	return nil
}

func (n *Node) SetUnavailable(name string) error {

	if name == "" {
		log.Debugf("Node not found: %s", name)
		return nil
	}

	node, err := n.storage.Node().Get(n.context, name)
	if err != nil {
		log.Errorf("Get node id by token error: %s", err)
		return err
	}

	if node == nil {
		log.Debugf("Node not found: %s", name)
		return nil
	}

	if err := n.storage.Node().SetUnavailable(n.context, node); err != nil {
		log.Errorf("Set node offline state error: %s", err)
		return err
	}

	return nil

}

func (n *Node) SetState(node *types.Node, state types.NodeState) (*types.Node, error) {

	node.State = state

	if err := n.storage.Node().SetState(n.context, node); err != nil {
		log.Errorf("Set node offline state error: %s", err)
		return nil, err
	}

	return node, nil
}

func (n *Node) SetInfo(node *types.Node, info types.NodeInfo) error {

	node.Info = info
	if err := n.storage.Node().SetInfo(n.context, node); err != nil {
		log.Errorf("Set node info error: %s", err)
		return err
	}

	return nil
}

func (n *Node) SetNetwork(node *types.Node, network types.Subnet) error {

	node.Network = network
	if err := n.storage.Node().SetNetwork(n.context, node); err != nil {
		log.Errorf("Set node network error: %s", err)
		return err
	}

	return nil
}

func (n *Node) GetSpec(name string) (types.NodeSpec, error) {

	var spec types.NodeSpec

	return spec, nil
}

func (n *Node) Remove(name string) error {

	log.V(logLevel).Debugf("Node: remove Node %s", name)

	if err := n.storage.Node().Remove(n.context, name); err != nil {
		log.V(logLevel).Debugf("Node: remove Node err: %s", err)
		return err
	}

	return nil
}

func NewNodeModel(ctx context.Context, stg storage.Storage) INode {
	return &Node{ctx, stg}
}
