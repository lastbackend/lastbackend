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

package etcd

import (
	"context"
	"errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"regexp"
	"time"
)

const (
	nodeStorage = "node"
	timeout     = 15
)

// Node Service type for interface in interfaces folder
type NodeStorage struct {
	storage.Node
}

func (s *NodeStorage) List(ctx context.Context) (map[string]*types.Node, error) {

	log.V(logLevel).Debugf("storage:etcd:node:> get list nodes")

	const filter = `\b.+` + nodeStorage + `\/.+\/(?:meta|info|status|online|network)\b`

	nodes := make(map[string]*types.Node, 0)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	key := keyCreate(nodeStorage)
	if err := client.MapList(ctx, key, filter, nodes); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> get nodes list err: %s", err.Error())
		return nil, err
	}

	log.V(logLevel).Debugf("storage:etcd:node:> get nodes list result: %d", len(nodes))

	return nodes, nil
}

func (s *NodeStorage) Get(ctx context.Context, name string) (*types.Node, error) {

	log.V(logLevel).Debugf("storage:etcd:node:> get by id: %s", name)

	if len(name) == 0 {
		err := errors.New("node can not be empty")
		log.V(logLevel).Errorf("storage:etcd:node:> get node err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + nodeStorage + `\/.+\/(?:meta|info|status|online|network)\b`

	client, destroy, err := getClient(ctx)
	if err != nil {
		return nil, err
	}
	defer destroy()

	node := new(types.Node)
	key := keyDirCreate(nodeStorage, name)
	if err := client.Map(ctx, key, filter, node); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> create client err: %s", err.Error())
		return nil, err
	}

	return node, nil
}

func (s *NodeStorage) GetSpec(ctx context.Context, node *types.Node) (*types.NodeSpec, error) {

	log.V(logLevel).Debugf("storage:etcd:node:> get node spec: %v", node)

	var (
		spec = new(types.NodeSpec)
	)

	spec.Pods = make(map[string]types.PodSpec)
	spec.Volumes = make(map[string]types.VolumeSpec)
	spec.Routes = make(map[string]types.RouteSpec)

	if err := s.checkNodeExists(ctx, node); err != nil {
		return nil, err
	}

	const filterPods = `\b.+` + nodeStorage + `\/.+\/spec\/pods\/(.+)\b`
	const filterVolumes = `\b.+` + nodeStorage + `\/.+\/spec\/volumes\/(.+)\b`
	const filterRoutes = `\b.+` + nodeStorage + `\/.+\/spec\/routes\/(.+)\b`

	client, destroy, err := getClient(ctx)
	if err != nil {
		return nil, err
	}
	defer destroy()

	keyPods := keyDirCreate(nodeStorage, node.Meta.Name, "spec", "pods")
	if err := client.Map(ctx, keyPods, filterPods, spec.Pods); err != nil && err.Error() != store.ErrEntityNotFound {
		log.V(logLevel).Errorf("storage:etcd:node:> get node spec: err: %s", err.Error())
		return nil, err
	}

	keyVolumes := keyDirCreate(nodeStorage, node.Meta.Name, "spec", "volumes")
	if err := client.Map(ctx, keyVolumes, filterVolumes, spec.Volumes); err != nil && err.Error() != store.ErrEntityNotFound {
		log.V(logLevel).Errorf("storage:etcd:node:> get node spec: err: %s", err.Error())
		return nil, err
	}

	keyRoutes := keyDirCreate(nodeStorage, node.Meta.Name, "spec", "routes")
	if err := client.Map(ctx, keyRoutes, filterRoutes, spec.Routes); err != nil && err.Error() != store.ErrEntityNotFound {
		log.V(logLevel).Errorf("storage:etcd:node:> get node spec: err: %s", err.Error())
		return nil, err
	}

	return spec, nil
}

func (s *NodeStorage) GetSpecPod(ctx context.Context, node, pod string) (*types.PodSpec, error) {

	log.V(logLevel).Debugf("storage:etcd:node:> get spec for pod: %s", pod)

	if len(node) == 0 {
		err := errors.New("node can not be empty")
		log.V(logLevel).Errorf("storage:etcd:node:> get spec for pod: %s", err.Error())
		return nil, err
	}

	if len(pod) == 0 {
		err := errors.New("pod can not be empty")
		log.V(logLevel).Errorf("storage:etcd:node:> get spec for pod: %s", err.Error())
		return nil, err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		return nil, err
	}
	defer destroy()

	spec := new(types.PodSpec)

	key := keyCreate(nodeStorage, node, "spec", "pods", pod)
	if err := client.Get(ctx, key, spec); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> get spec for pod err: %s", err.Error())
		return nil, err
	}

	return spec, nil
}

func (s *NodeStorage) GetSpecVolume(ctx context.Context, node, volume string) (*types.PodSpec, error) {

	log.V(logLevel).Debugf("storage:etcd:node:> get spec for volume: %s", volume)

	if len(node) == 0 {
		err := errors.New("node can not be empty")
		log.V(logLevel).Errorf("storage:etcd:node:> get spec for volume: %s", err.Error())
		return nil, err
	}

	if len(volume) == 0 {
		err := errors.New("volume can not be empty")
		log.V(logLevel).Errorf("storage:etcd:node:> get spec for volume: %s", err.Error())
		return nil, err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		return nil, err
	}
	defer destroy()

	spec := new(types.PodSpec)

	key := keyCreate(nodeStorage, node, "spec", "volumes", volume)
	if err := client.Get(ctx, key, spec); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> get spec for volume err: %s", err.Error())
		return nil, err
	}

	return spec, nil
}

func (s *NodeStorage) GetSpecRoute(ctx context.Context, node, route string) (*types.RouteSpec, error) {

	log.V(logLevel).Debugf("storage:etcd:node:> get spec for route: %s", route)

	if len(node) == 0 {
		err := errors.New("node can not be empty")
		log.V(logLevel).Errorf("storage:etcd:node:> get spec for route: %s", err.Error())
		return nil, err
	}

	if len(route) == 0 {
		err := errors.New("route can not be empty")
		log.V(logLevel).Errorf("storage:etcd:node:> get spec for route: %s", err.Error())
		return nil, err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		return nil, err
	}
	defer destroy()

	spec := new(types.RouteSpec)

	key := keyCreate(nodeStorage, node, "spec", "routes", route)
	if err := client.Get(ctx, key, spec); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> get spec for route err: %s", err.Error())
		return nil, err
	}

	return spec, nil
}

func (s *NodeStorage) Insert(ctx context.Context, node *types.Node) error {

	log.V(logLevel).Debugf("storage:etcd:node:> insert node: %#v", node)

	if err := s.checkNodeArgument(node); err != nil {
		return err
	}

	node.Meta.Created = time.Now()
	node.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> create client err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(nodeStorage, node.Meta.Name, "meta")
	if err := tx.Create(keyMeta, &node.Meta, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> insert node err: %s", err.Error())
		return err
	}

	keyInfo := keyCreate(nodeStorage, node.Meta.Name, "info")
	if err := tx.Create(keyInfo, &node.Info, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> insert node err: %s", err.Error())
		return err
	}

	keyStatus := keyCreate(nodeStorage, node.Meta.Name, "status")
	if err := tx.Create(keyStatus, &node.Status, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> insert node err: %s", err.Error())
		return err
	}

	keyOnline := keyCreate(nodeStorage, node.Meta.Name, "online")
	if err := tx.Create(keyOnline, true, timeout); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> insert node err: %s", err.Error())
		return err
	}

	keyNetwork := keyCreate(nodeStorage, node.Meta.Name, "network")
	if err := tx.Create(keyNetwork, &node.Network, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> insert node err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> insert node err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) Update(ctx context.Context, node *types.Node) error {

	log.V(logLevel).Debugf("storage:etcd:node:> update node: %#v", node)

	if err := s.checkNodeExists(ctx, node); err != nil {
		return err
	}

	node.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> update node err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(nodeStorage, node.Meta.Name, "meta")
	if err := tx.Update(keyMeta, &node.Meta, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> update node err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> update node err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) SetStatus(ctx context.Context, node *types.Node) error {

	log.V(logLevel).Debugf("storage:etcd:node:> update node status: %#v", node)

	if err := s.checkNodeExists(ctx, node); err != nil {
		return err
	}

	node.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> update node status err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	key := keyCreate(nodeStorage, node.Meta.Name, "status")
	if err := tx.Update(key, &node.Status, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> update node status err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> update node status err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) SetInfo(ctx context.Context, node *types.Node) error {

	log.V(logLevel).Debugf("storage:etcd:node:> update node info: %#v", node)

	if err := s.checkNodeExists(ctx, node); err != nil {
		return err
	}

	node.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> update node info err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	key := keyCreate(nodeStorage, node.Meta.Name, "info")
	if err := tx.Update(key, &node.Info, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> update node info err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> update node info err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) SetNetwork(ctx context.Context, node *types.Node) error {

	log.V(logLevel).Debugf("storage:etcd:node:> update node network: %#v", node)

	if err := s.checkNodeExists(ctx, node); err != nil {
		return err
	}

	node.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> update node network err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	key := keyCreate(nodeStorage, node.Meta.Name, "network")
	if err := tx.Update(key, &node.Network, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> update node network err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> update node network err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) SetOnline(ctx context.Context, node *types.Node) error {

	log.V(logLevel).Debugf("storage:etcd:node:> set node online: %#v", node)

	if err := s.checkNodeExists(ctx, node); err != nil {
		return err
	}

	node.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> set node online err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(nodeStorage, node.Meta.Name, "online")
	if err := tx.Upsert(keyMeta, true, timeout); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> set node online err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> set node online err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) SetOffline(ctx context.Context, node *types.Node) error {

	log.V(logLevel).Debugf("storage:etcd:node:> set node offline: %#v", node)

	if err := s.checkNodeExists(ctx, node); err != nil {
		return err
	}

	node.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> set node offline err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	key := keyCreate(nodeStorage, node.Meta.Name, "online")
	tx.Delete(key)

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> set node offline err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) InsertPod(ctx context.Context, node *types.Node, pod *types.Pod) error {

	log.V(logLevel).Debugf("storage:etcd:node:> insert pod: %#v", pod)

	if err := s.checkNodeExists(ctx, node); err != nil {
		return err
	}

	if err := s.checkPodArgument(pod); err != nil {
		return err
	}

	node.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(nodeStorage, node.Meta.Name, "meta")
	if err := client.Update(ctx, keyMeta, node.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> insert pod err: %s", err.Error())
		return err
	}

	keyPod := keyCreate(nodeStorage, node.Meta.Name, "spec", "pods", pod.SelfLink())
	if err := client.Create(ctx, keyPod, pod.Spec, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> insert pod err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) RemovePod(ctx context.Context, node *types.Node, pod *types.Pod) error {

	if err := s.checkNodeExists(ctx, node); err != nil {
		return err
	}

	if err := s.checkPodArgument(pod); err != nil {
		return err
	}

	if err := s.checkPodSpecExists(ctx, node, pod); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(nodeStorage, node.Meta.Name, "spec", "pods", pod.SelfLink())
	if err := client.Delete(ctx, key); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> remove node err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) InsertVolume(ctx context.Context, node *types.Node, volume *types.Volume) error {

	if err := s.checkNodeExists(ctx, node); err != nil {
		return err
	}

	if err := s.checkVolumeArgument(volume); err != nil {
		return err
	}

	node.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(nodeStorage, node.Meta.Name, "meta")
	if err := client.Update(ctx, keyMeta, node.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> insert volume err: %s", err.Error())
		return err
	}

	keyVolume := keyCreate(nodeStorage, node.Meta.Name, "spec", "volumes", volume.SelfLink())
	if err := client.Create(ctx, keyVolume, volume.Spec, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> insert volume err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) RemoveVolume(ctx context.Context, node *types.Node, volume *types.Volume) error {

	if err := s.checkNodeExists(ctx, node); err != nil {
		return err
	}

	if err := s.checkVolumeArgument(volume); err != nil {
		return err
	}

	if err := s.checkVolumeSpecExists(ctx, node, volume); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(nodeStorage, node.Meta.Name, "spec", "volumes", volume.SelfLink())
	if err := client.Delete(ctx, key); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> remove node err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) InsertRoute(ctx context.Context, node *types.Node, route *types.Route) error {

	if err := s.checkNodeExists(ctx, node); err != nil {
		return err
	}

	if err := s.checkRouteArgument(route); err != nil {
		return err
	}

	node.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(nodeStorage, node.Meta.Name, "meta")
	if err := client.Update(ctx, keyMeta, node.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> insert route err: %s", err.Error())
		return err
	}

	keyRoute := keyCreate(nodeStorage, node.Meta.Name, "spec", "routes", route.SelfLink())
	if err := client.Create(ctx, keyRoute, route.Spec, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> insert route err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) RemoveRoute(ctx context.Context, node *types.Node, route *types.Route) error {

	if err := s.checkNodeExists(ctx, node); err != nil {
		return err
	}

	if err := s.checkRouteArgument(route); err != nil {
		return err
	}

	if err := s.checkRouteSpecExists(ctx, node, route); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(nodeStorage, node.Meta.Name, "spec", "routes", route.SelfLink())
	if err := client.Delete(ctx, key); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> remove node err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) Remove(ctx context.Context, node *types.Node) error {

	log.V(logLevel).Debugf("storage:etcd:node:> remove node: %#v", node)

	if err := s.checkNodeExists(ctx, node); err != nil {
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(nodeStorage, node.Meta.Name)
	if err := client.DeleteDir(ctx, key); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> remove node err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) Watch(ctx context.Context, node chan *types.Node) error {

	log.V(logLevel).Debug("storage:etcd:node:> watch node")

	const filter = `\b.+` + nodeStorage + `\/(.+)\/alive\b`

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> create client err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(nodeStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 2 {
			return
		}

		n, _ := s.Get(ctx, keys[1])
		if n == nil {
			return
		}

		// TODO: check previous node alive status to prevent multi calls
		if action == "PUT" {
			n.Online = true
			node <- n
			return
		}

		if action == "DELETE" {
			n.Online = false
			node <- n
			return
		}

		return
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> watch node err: %s", err.Error())
		return err
	}

	return nil
}

// Clear node storage
func (s *NodeStorage) Clear(ctx context.Context) error {

	log.V(logLevel).Debugf("storage:etcd:node:> clear")

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> clear err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := client.DeleteDir(ctx, nodeStorage); err != nil {
		log.V(logLevel).Errorf("storage:etcd:node:> clear err: %s", err.Error())
		return err
	}

	return nil
}

func newNodeStorage() *NodeStorage {
	s := new(NodeStorage)
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

func (s *NodeStorage) checkPodArgument(pod *types.Pod) error {
	if pod == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if pod.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

func (s *NodeStorage) checkVolumeArgument(volume *types.Volume) error {
	if volume == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if volume.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

func (s *NodeStorage) checkRouteArgument(route *types.Route) error {
	if route == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	if route.Meta.Name == "" {
		return errors.New(store.ErrStructArgIsInvalid)
	}

	return nil
}

func (s *NodeStorage) checkNodeExists(ctx context.Context, node *types.Node) error {

	if err := s.checkNodeArgument(node); err != nil {
		return err
	}

	log.V(logLevel).Debugf("storage:etcd:node:> check node exists")

	if _, err := s.Get(ctx, node.Meta.Name); err != nil {
		log.V(logLevel).Debugf("storage:etcd:node:> check node exists err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) checkPodSpecExists(ctx context.Context, node *types.Node, pod *types.Pod) error {

	if err := s.checkNodeArgument(node); err != nil {
		return err
	}

	if err := s.checkPodArgument(pod); err != nil {
		return err
	}

	log.V(logLevel).Debugf("storage:etcd:node:> check pod spec exists")

	if _, err := s.GetSpecPod(ctx, node.Meta.Name, pod.SelfLink()); err != nil {
		log.V(logLevel).Debugf("storage:etcd:node:> check pod spec exists err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) checkVolumeSpecExists(ctx context.Context, node *types.Node, volume *types.Volume) error {

	if err := s.checkNodeArgument(node); err != nil {
		return err
	}

	if err := s.checkVolumeArgument(volume); err != nil {
		return err
	}

	log.V(logLevel).Debugf("storage:etcd:node:> check volume spec exists")

	if _, err := s.GetSpecVolume(ctx, node.Meta.Name, volume.SelfLink()); err != nil {
		log.V(logLevel).Debugf("storage:etcd:node:> check volume spec exists err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) checkRouteSpecExists(ctx context.Context, node *types.Node, route *types.Route) error {

	if err := s.checkNodeArgument(node); err != nil {
		return err
	}

	if err := s.checkRouteArgument(route); err != nil {
		return err
	}

	log.V(logLevel).Debugf("storage:etcd:node:> check route spec exists")

	if _, err := s.GetSpecRoute(ctx, node.Meta.Name, route.SelfLink()); err != nil {
		log.V(logLevel).Debugf("storage:etcd:node:> check route spec exists err: %s", err.Error())
		return err
	}

	return nil
}
