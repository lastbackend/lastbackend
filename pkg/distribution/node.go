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

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const (
	logNodePrefix = "distribution:node"
)

type INode interface {
	List() (map[string]*types.Node, error)
	Create(opts *types.NodeCreateOptions) (*types.Node, error)

	Get(name string) (*types.Node, error)
	GetSpec(node *types.Node) (*types.NodeSpec, error)

	SetMeta(node *types.Node, meta *types.NodeUpdateMetaOptions) error
	SetStatus(node *types.Node, state types.NodeStatus) error
	SetInfo(node *types.Node, info types.NodeInfo) error
	SetNetwork(node *types.Node, network types.NetworkSpec) error
	SetOnline(node *types.Node) error
	SetOffline(node *types.Node) error

	InsertPod(node *types.Node, pod *types.Pod) error
	UpdatePod(node *types.Node, pod *types.Pod) error
	RemovePod(node *types.Node, pod *types.Pod) error
	InsertVolume(node *types.Node, volume *types.Volume) error
	RemoveVolume(node *types.Node, volume *types.Volume) error
	InsertRoute(node *types.Node, route *types.Route) error
	RemoveRoute(node *types.Node, route *types.Route) error
	Remove(node *types.Node) error
}

type Node struct {
	context context.Context
	storage storage.Storage
}

func (n *Node) List() (map[string]*types.Node, error) {
	return n.storage.Node().List(n.context)
}

func (n *Node) Create(opts *types.NodeCreateOptions) (*types.Node, error) {

	log.Debugf("%s:create:> create node in cluster", logNodePrefix)

	ni := new(types.Node)
	ni.Meta.SetDefault()

	ni.Meta.Name = opts.Meta.Name
	ni.Meta.Token = opts.Meta.Token
	ni.Meta.Region = opts.Meta.Region
	ni.Meta.Provider = opts.Meta.Provider

	ni.Info = opts.Info
	ni.Status = opts.Status

	if ni.Meta.Token == "" {
		ni.Meta.Token = generator.GenerateRandomString(32)
	}

	ni.Online = true

	if err := n.storage.Node().Insert(n.context, ni); err != nil {
		log.Debugf("%s:create:> insert node err: %s", logNodePrefix, err.Error())
		return nil, err
	}

	return ni, nil
}

func (n *Node) Get(name string) (*types.Node, error) {

	log.V(logLevel).Debugf("%s:get:> get by name %s", logNodePrefix, name)

	node, err := n.storage.Node().Get(n.context, name)
	if err != nil {

		if err.Error() == store.ErrEntityNotFound {
			log.V(logLevel).Warnf("%s:get:> get: node %s not found", logNodePrefix, name)
			return nil, nil
		}

		log.V(logLevel).Debugf("%s:get:> get node `%s` err: %s", logNodePrefix, name, err.Error())
		return nil, err
	}

	return node, nil
}

func (n *Node) GetSpec(node *types.Node) (*types.NodeSpec, error) {

	log.V(logLevel).Debugf("%s:getspec:> get node spec: %s", logNodePrefix, node.Meta.Name)

	spec, err := n.storage.Node().GetSpec(n.context, node)
	if err != nil {
		log.V(logLevel).Debugf("%s:getspec:> get Node `%s` err: %s", logNodePrefix, node.Meta.Name, err.Error())
		return nil, err
	}

	return spec, nil
}

func (n *Node) SetMeta(node *types.Node, meta *types.NodeUpdateMetaOptions) error {

	log.V(logLevel).Debugf("%s:setmeta:> update Node %#v", logNodePrefix, meta)
	if meta == nil {
		log.V(logLevel).Errorf("%s:setmeta:> update Node err: %s", logNodePrefix, errors.New(errors.ArgumentIsEmpty))
		return errors.New(errors.ArgumentIsEmpty)
	}

	node.Meta.Set(meta)

	if err := n.storage.Node().Update(n.context, node); err != nil {
		log.V(logLevel).Errorf("%s:setmeta:> update Node meta err: %s", logNodePrefix, err.Error())
		return err
	}

	return nil
}

func (n *Node) SetOnline(node *types.Node) error {

	if err := n.storage.Node().SetOnline(n.context, node); err != nil {
		log.Errorf("%s:setonline:> set node online state error: %s", logNodePrefix, err.Error())
		return err
	}

	return nil
}

func (n *Node) SetOffline(node *types.Node) error {

	if err := n.storage.Node().SetOffline(n.context, node); err != nil {
		log.Errorf("%s:setoffline:> set node offline state error: %s", logNodePrefix, err.Error())
		return err
	}

	return nil

}

func (n *Node) SetStatus(node *types.Node, status types.NodeStatus) error {

	node.Status = status

	if err := n.storage.Node().SetStatus(n.context, node); err != nil {
		log.Errorf("%s:setstatus:> set node offline state error: %s", logNodePrefix, err.Error())
		return err
	}

	return nil
}

func (n *Node) SetInfo(node *types.Node, info types.NodeInfo) error {

	node.Info = info
	if err := n.storage.Node().SetInfo(n.context, node); err != nil {
		log.Errorf("%s:setinfo:> set node info error: %s", logNodePrefix, err.Error())
		return err
	}

	return nil
}

func (n *Node) SetNetwork(node *types.Node, network types.NetworkSpec) error {

	node.Network = network
	if err := n.storage.Node().SetNetwork(n.context, node); err != nil {
		log.Errorf("%s:setnetwork:> set node network error: %s", logNodePrefix, err.Error())
		return err
	}

	return nil
}

func (n *Node) InsertPod(node *types.Node, pod *types.Pod) error {

	if err := n.storage.Node().InsertPod(n.context, node, pod); err != nil {
		log.Errorf("%s:insertpod:> set node network error: %s", logNodePrefix, err.Error())
		return err
	}

	return nil
}

func (n *Node) UpdatePod(node *types.Node, pod *types.Pod) error {

	if err := n.storage.Node().UpdatePod(n.context, node, pod); err != nil {
		log.Errorf("%s:updatepod:> set node network error: %s", logNodePrefix, err.Error())
		return err
	}

	return nil
}

func (n *Node) RemovePod(node *types.Node, pod *types.Pod) error {

	if err := n.storage.Node().RemovePod(n.context, node, pod); err != nil {
		log.Errorf("%s:removepod:> set node network error: %s", logNodePrefix, err.Error())
		return err
	}

	return nil
}

func (n *Node) InsertVolume(node *types.Node, volume *types.Volume) error {
	return nil
}

func (n *Node) RemoveVolume(node *types.Node, volume *types.Volume) error {
	return nil
}

func (n *Node) InsertRoute(node *types.Node, route *types.Route) error {
	return nil
}

func (n *Node) RemoveRoute(node *types.Node, route *types.Route) error {
	return nil
}

func (n *Node) Remove(node *types.Node) error {

	log.V(logLevel).Debugf("%s:remove:> remove node %s", logNodePrefix, node.Meta.Name)

	if err := n.storage.Node().Remove(n.context, node); err != nil {
		log.V(logLevel).Debugf("%s:remove:> remove node err: %s", logNodePrefix, err.Error())
		return err
	}

	return nil
}

func NewNodeModel(ctx context.Context, stg storage.Storage) INode {
	return &Node{ctx, stg}
}
