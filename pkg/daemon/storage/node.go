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

package storage

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/storage/store"
	"github.com/satori/go.uuid"
	"time"
)

const nodeStorage = "node"
const nodeStateTTL = 3000

// Namespace Service type for interface in interfaces folder
type NodeStorage struct {
	INode
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

func (s *NodeStorage) List(ctx context.Context) ([]*types.Node, error) {

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	key := s.util.Key(ctx, nodeStorage)
	list := []*types.Node{}
	if err := client.List(ctx, key, "", &list); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return list, nil
}

func (s *NodeStorage) Get(ctx context.Context, hostname string) (*types.Node, error) {
	node := new(types.Node)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	keyMeta := s.util.Key(ctx, nodeStorage, hostname, "meta")
	if err := client.Get(ctx, keyMeta, &node.Meta); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	keyState := s.util.Key(ctx, nodeStorage, hostname, "state")
	if err := client.Get(ctx, keyState, &node.State); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	keySpec := s.util.Key(ctx, nodeStorage, hostname, "spec")
	if err := client.List(ctx, keySpec, "", &node.Spec.Pods); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return node, nil
}

func (s *NodeStorage) Insert(ctx context.Context, meta *types.NodeMeta, state *types.NodeState) (*types.Node, error) {

	var (
		id   = uuid.NewV4().String()
		node = new(types.Node)
	)

	client, destroy, err := s.Client()
	if err != nil {
		return nil, err
	}
	defer destroy()

	tx := client.Begin(ctx)

	node.Meta = *meta
	node.Meta.ID = id
	node.Meta.Labels = map[string]string{"tier": "node"}
	node.Meta.Updated = time.Now()
	node.Meta.Created = time.Now()

	node.State = *state

	keyMeta := s.util.Key(ctx, nodeStorage, node.Meta.Hostname, "meta")
	if err := tx.Create(keyMeta, node.Meta, 0); err != nil {
		fmt.Println("meta", err.Error())
		return nil, err
	}

	keyState := s.util.Key(ctx, nodeStorage, node.Meta.Hostname, "state")
	if err := tx.Create(keyState, node.State, nodeStateTTL); err != nil {
		fmt.Println("state", err.Error())
		return nil, err
	}

	keySpec := s.util.Key(ctx, nodeStorage, node.Meta.Hostname, "spec")
	if err := tx.Create(keySpec, node.Spec, 0); err != nil {
		fmt.Println("spec", err.Error())
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("commit", err.Error())
		return nil, err
	}

	return node, nil
}

func (s *NodeStorage) UpdateMeta(ctx context.Context, meta *types.NodeMeta) error {
	meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)
	keyMeta := s.util.Key(ctx, nodeStorage, meta.Hostname, "meta")
	if err := tx.Update(keyMeta, meta, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (s *NodeStorage) UpdateState(ctx context.Context, meta *types.NodeMeta, state *types.NodeState) error {
	meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)
	keyMeta := s.util.Key(ctx, nodeStorage, meta.Hostname, "meta")
	if err := tx.Update(keyMeta, meta, 0); err != nil {
		return err
	}

	keyState := s.util.Key(ctx, nodeStorage, meta.Hostname, "state")
	if err := tx.Update(keyState, state, nodeStateTTL); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *NodeStorage) InsertPod(ctx context.Context, meta *types.NodeMeta, pod *types.PodNodeSpec) error {
	meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)
	keyMeta := s.util.Key(ctx, nodeStorage, meta.Hostname, "meta")
	if err := tx.Update(keyMeta, meta, 0); err != nil {
		return err
	}

	keyPod := s.util.Key(ctx, nodeStorage, meta.Hostname, "pod", pod.Meta.ID)
	if err := tx.Create(keyPod, pod, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *NodeStorage) UpdatePod(ctx context.Context, meta *types.NodeMeta, pod *types.PodNodeSpec) error {
	meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)
	keyMeta := s.util.Key(ctx, nodeStorage, meta.Hostname, "meta")
	if err := tx.Update(keyMeta, meta, 0); err != nil {
		return err
	}

	keyPod := s.util.Key(ctx, nodeStorage, meta.Hostname, "pod", pod.Meta.ID)
	if err := tx.Update(keyPod, pod, 0); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *NodeStorage) RemovePod(ctx context.Context, meta *types.NodeMeta, pod *types.PodNodeSpec) error {
	meta.Updated = time.Now()

	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)
	keyMeta := s.util.Key(ctx, nodeStorage, meta.Hostname, "meta")
	if err := tx.Update(keyMeta, meta, 0); err != nil {
		return err
	}

	keyPod := s.util.Key(ctx, nodeStorage, meta.Hostname, "pod", pod.Meta.ID)
	tx.Delete(keyPod)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *NodeStorage) Remove(ctx context.Context, meta *types.NodeMeta) error {
	client, destroy, err := s.Client()
	if err != nil {
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)
	key := s.util.Key(ctx, nodeStorage, meta.Hostname)
	tx.DeleteDir(key)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func newNodeStorage(config store.Config, util IUtil) *NodeStorage {
	s := new(NodeStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
