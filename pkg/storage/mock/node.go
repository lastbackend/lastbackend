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
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

// Node Service type for interface in interfaces folder
type NodeStorage struct {
	storage.Node
	data map[string]*types.Node
}

func (s *NodeStorage) List(ctx context.Context) ([]*types.Node, error) {

	nl := make([]*types.Node, 0)

	for _, n := range s.data {
		nl = append(nl, n)
	}

	return nl, nil
}

func (s *NodeStorage) Get(ctx context.Context, name string) (*types.Node, error) {

	if n, ok := s.data[name]; ok {
		return n, nil
	}

	return nil, nil
}

func (s *NodeStorage) Insert(ctx context.Context, node *types.Node) error {

	if err := s.checkNodeArgument(node); err != nil {
		return err
	}

	s.data[node.Meta.Name] = node

	return nil
}

func (s *NodeStorage) Update(ctx context.Context, node *types.Node) error {

	if err := s.checkNodeExists(node); err != nil {
		return err
	}

	s.data[node.Meta.Name].Meta = node.Meta
	return nil
}

func (s *NodeStorage) SetState (ctx context.Context, node *types.Node) error {

	if err := s.checkNodeExists(node); err != nil {
		return err
	}

	s.data[node.Meta.Name].State = node.State
	return nil
}


func (s *NodeStorage) SetInfo(ctx context.Context, node *types.Node) error {
	if err := s.checkNodeExists(node); err != nil {
		return err
	}

	s.data[node.Meta.Name].Info = node.Info
	return nil
}

func (s *NodeStorage) SetNetwork(ctx context.Context, node *types.Node) error {
	if err := s.checkNodeExists(node); err != nil {
		return err
	}

	s.data[node.Meta.Name].Network = node.Network
	return nil
}

func (s *NodeStorage) SetAvailable(ctx context.Context, node *types.Node) error {

	if err := s.checkNodeExists(node); err != nil {
		return err
	}

	s.data[node.Meta.Name].Alive = true

	return nil
}

func (s *NodeStorage) SetUnavailable(ctx context.Context, node *types.Node) error {

	if err := s.checkNodeExists(node); err != nil {
		return err
	}

	s.data[node.Meta.Name].Alive = false

	return nil
}

func (s *NodeStorage) InsertPod(ctx context.Context, node *types.Node, pod *types.Pod) error {

	if err := s.checkNodeExists(node); err != nil {
		return err
	}

	s.data[node.Meta.Name].Spec.Pods[pod.Meta.Name] = pod

	return nil
}

func (s *NodeStorage) RemovePod(ctx context.Context, node *types.Node, pod *types.Pod) error {

	if err := s.checkNodeExists(node); err != nil {
		return err
	}

	if _, ok :=  s.data[node.Meta.Name].Spec.Pods[pod.Meta.Name]; !ok {
		return errors.New(store.ErrKeyNotFound)
	}

	delete(s.data[node.Meta.Name].Spec.Pods, pod.Meta.Name)

	return nil
}

func (s *NodeStorage) InsertVolume(ctx context.Context, node *types.Node, volume *types.Volume) error {

	if err := s.checkNodeExists(node); err != nil {
		return err
	}

	s.data[node.Meta.Name].Spec.Volumes[volume.Meta.Name] = volume

	return nil
}

func (s *NodeStorage) RemoveVolume(ctx context.Context, node *types.Node, volume *types.Volume)  error {


	if err := s.checkNodeExists(node); err != nil {
		return err
	}

	if _, ok :=  s.data[node.Meta.Name].Spec.Volumes[volume.Meta.Name]; !ok {
		return errors.New(store.ErrKeyNotFound)
	}

	delete(s.data[node.Meta.Name].Spec.Volumes, volume.Meta.Name)

	return nil
}

func (s *NodeStorage) InsertRoute(ctx context.Context, node *types.Node, route *types.Route)  error {

	if err := s.checkNodeExists(node); err != nil {
		return err
	}

	s.data[node.Meta.Name].Spec.Routes[route.Meta.Name] = route

	return nil
}

func (s *NodeStorage) RemoveRoute(ctx context.Context, node *types.Node, route *types.Route)  error {

	if err := s.checkNodeExists(node); err != nil {
		return err
	}

	if _, ok :=  s.data[node.Meta.Name].Spec.Routes[route.Meta.Name]; !ok {
		return errors.New(store.ErrKeyNotFound)
	}

	delete(s.data[node.Meta.Name].Spec.Routes, route.Meta.Name)

	return nil
}

func (s *NodeStorage) Remove(ctx context.Context, node *types.Node) error {

	if err := s.checkNodeExists(node); err != nil {
		return err
	}

	return nil
}

func (s *NodeStorage) Watch(ctx context.Context, node chan *types.Node) error {
	return nil
}

func newNodeStorage() *NodeStorage {
	s := new(NodeStorage)
	s.data = make(map[string]*types.Node)
	return s
}


func (s *NodeStorage) checkNodeArgument(node *types.Node) error {
	if node == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if node.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

func (s *NodeStorage) checkNodeExists(node *types.Node) error {

	if err := s.checkNodeArgument(node); err != nil {
		return err
	}

	if _, ok := s.data[node.Meta.Name]; !ok {
		return errors.New(store.ErrKeyNotFound)
	}

	return nil
}