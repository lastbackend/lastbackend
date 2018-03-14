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

package mock

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
)

const (
	nodeStorage = "node"
	timeout     = 15
)

// Node Service type for interface in interfaces folder
type NodeStorage struct {
	storage.Node
}

func (s *NodeStorage) List(ctx context.Context) ([]*types.Node, error) {
	return make([]*types.Node, 0), nil
}

func (s *NodeStorage) Get(ctx context.Context, id string) (*types.Node, error) {
	return new(types.Node), nil
}

func (s *NodeStorage) Insert(ctx context.Context, node *types.Node) error {
	return nil
}

func (s *NodeStorage) Update(ctx context.Context, node *types.Node) error {
	return nil
}

func (s *NodeStorage) InsertPod(ctx context.Context, meta *types.NodeMeta, pod *types.Pod) error {
	return nil
}

func (s *NodeStorage) UpdatePod(ctx context.Context, meta *types.NodeMeta, pod *types.Pod) error {
	return nil
}

func (s *NodeStorage) RemovePod(ctx context.Context, meta *types.NodeMeta, pod *types.Pod) error {
	return nil
}

func (s *NodeStorage) Remove(ctx context.Context, name string) error {
	return nil
}

func (s *NodeStorage) Watch(ctx context.Context, node chan *types.Node) error {
	return nil
}

func newNodeStorage() *NodeStorage {
	s := new(NodeStorage)
	return s
}
