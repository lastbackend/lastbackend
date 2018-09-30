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

	"encoding/json"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

const (
	logNodePrefix = "distribution:node"
)

type Node struct {
	context context.Context
	storage storage.Storage
}

func (n *Node) List() (*types.NodeList, error) {
	log.V(logLevel).Debugf("%s:list:> get nodes list", logNodePrefix)

	nodes := types.NewNodeList()

	err := n.storage.List(n.context, n.storage.Collection().Node().Info(), "", nodes, nil)
	if err != nil {
		log.V(logLevel).Debugf("%s:list:> get nodes list err: %v", logNodePrefix, err)
		return nil, err
	}
	return nodes, nil
}

func (n *Node) Put(opts *types.NodeCreateOptions) (*types.Node, error) {

	log.V(logLevel).Debugf("%s:create:> create node in cluster", logNodePrefix)

	ni := new(types.Node)
	ni.Meta.SetDefault()

	ni.Meta.Name = opts.Meta.Name
	ni.Meta.NodeInfo = opts.Info
	ni.Status = opts.Status
	ni.Status.Online = true

	ni.Spec.Security.TLS = opts.Security.TLS

	if opts.Security.SSL != nil {
		ni.Spec.Security.SSL = new(types.NodeSSL)
		ni.Spec.Security.SSL.CA = opts.Security.SSL.CA
		ni.Spec.Security.SSL.Cert = opts.Security.SSL.Cert
		ni.Spec.Security.SSL.Key = opts.Security.SSL.Key
	}

	ni.SelfLink()

	if err := n.storage.Put(n.context, n.storage.Collection().Node().Info(), n.storage.Key().Node(ni.Meta.Name), ni, nil); err != nil {
		log.V(logLevel).Debugf("%s:create:> insert node err: %v", logNodePrefix, err)
		return nil, err
	}

	return ni, nil
}

func (n *Node) Get(hostname string) (*types.Node, error) {

	log.V(logLevel).Debugf("%s:get:> get by hostname %s", logNodePrefix, hostname)

	node := new(types.Node)

	err := n.storage.Get(n.context, n.storage.Collection().Node().Info(), n.storage.Key().Node(hostname), &node, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> get: node %s not found", logNodePrefix, hostname)
			return nil, nil
		}

		log.V(logLevel).Debugf("%s:get:> get node `%s` err: %v", logNodePrefix, hostname, err)
		return nil, err
	}

	return node, nil
}

func (n *Node) Set(node *types.Node) error {

	log.V(logLevel).Debugf("%s:setmeta:> update Node %#v", logNodePrefix, node)
	if err := n.storage.Set(n.context, n.storage.Collection().Node().Info(), n.storage.Key().Node(node.Meta.Name), node, nil); err != nil {
		log.V(logLevel).Errorf("%s:setmeta:> update Node meta err: %v", logNodePrefix, err)
		return err
	}

	return nil
}

func (n *Node) Remove(node *types.Node) error {

	log.V(logLevel).Debugf("%s:remove:> remove node %s", logNodePrefix, node.Meta.Name)

	if err := n.storage.Del(n.context, n.storage.Collection().Node().Info(), n.storage.Key().Node(node.Meta.Name)); err != nil {
		log.V(logLevel).Debugf("%s:remove:> remove node err: %v", logNodePrefix, err)
		return err
	}

	return nil
}

// Watch node changes
func (n *Node) Watch(ch chan types.NodeEvent, rev *int64) error {

	log.V(logLevel).Debugf("%s:watch:> watch node", logNodePrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-n.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := types.NodeEvent{}
				res.Action = e.Action
				res.Name = e.Name

				obj := new(types.Node)

				if err := json.Unmarshal(e.Data.([]byte), &obj); err != nil {
					log.Errorf("%s:watch:> parse json", logNodePrefix)
					continue
				}

				res.Data = obj

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	opts.Rev = rev

	if err := n.storage.Watch(n.context, n.storage.Collection().Node().Info(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewNodeModel(ctx context.Context, stg storage.Storage) *Node {
	return &Node{ctx, stg}
}
