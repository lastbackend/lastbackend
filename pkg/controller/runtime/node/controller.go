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

package node

import (
	"context"
	"reflect"

	"github.com/lastbackend/lastbackend/pkg/cache"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/storage"

	stgtypes "github.com/lastbackend/lastbackend/pkg/storage/etcd/types"
	"encoding/json"
)

const (
	logPrefix = "nodecontroller"
)

type Controller struct {
	node   chan *types.Node
	active bool

	cache *cache.NodeCache
}

func (nc *Controller) Watch(node chan *types.Node) {

	var (
		stg   = envs.Get().GetStorage()
		event = make(chan *stgtypes.WatcherEvent)
	)

	log.Debugf("%s:> start watch", logPrefix)

	go func() {
		for {
			select {
			case n := <-nc.node:
				{
					if !nc.active {
						log.Debugf("%s:> skip management cause it is in slave mode", logPrefix)
						continue
					}

					log.Debugf("%s:> node check state: %s", logPrefix, n.Meta.Name)

					item := nc.cache.Get(n.Info.Hostname)
					if item == nil || !reflect.DeepEqual(item, n) {
						nc.cache.Set(n.Info.Hostname, n)

						nodes := nc.cache.List()

						cl := new(types.Cluster)

						if err := envs.Get().GetStorage().Get(context.Background(), storage.ClusterKind, types.EmptyString, &cl); err != nil {
							log.Errorf("%s:> get cluster info err: %v", logPrefix, err)
							continue
						}

						cl.Status = *getClusterStatus(nodes)

						if err := envs.Get().GetStorage().Upsert(context.Background(), storage.ClusterKind, types.EmptyString, cl, nil); err != nil {
							log.Errorf("%s:> set cluster status err: %v", logPrefix, err)
							continue
						}

					}

					if n.Online {
						log.Debugf("%s:> node set alive, try to provision on it pods: %s", logPrefix, n.Meta.Name)
						node <- n
						continue
					}

					log.Debugf("%s:> node set offline, try to move all pods to another", logPrefix)

				}
			}
		}
	}()

	go func() {
		for {
			select {
			case e := <-event:
				if e.Data == nil {
					continue
				}

				node := new(types.Node)

				if err := json.Unmarshal(e.Data.([]byte), *node); err != nil {
					log.Errorf("%s:> parse data err: %v", logPrefix, err)
					continue
				}

				nc.node <- node
			}
		}
	}()

	stg.Watch(context.Background(), storage.NodeKind, event)
}

func (nc *Controller) Pause() {
	nc.active = false
}

func (nc *Controller) Resume() {
	nc.active = true
	log.Debugf("%s:> start check pods state", logPrefix)
}

func NewNodeController(ctx context.Context) *Controller {

	sc := new(Controller)
	sc.active = false
	sc.node = make(chan *types.Node)
	sc.cache = cache.NewNodeCache()

	nodes := make(map[string]*types.Node, 0)

	err := envs.Get().GetStorage().Map(ctx, storage.NodeKind, "", nodes)
	if err != nil {
		log.Fatalf("%s:> get nodes list err: %v", logPrefix, err)
	}

	for _, node := range nodes {
		sc.cache.Set(node.Info.Hostname, node)
	}

	cl := new(types.Cluster)

	err = envs.Get().GetStorage().Get(context.Background(), storage.ClusterKind, types.EmptyString, &cl)
	if err != nil {
		log.Fatalf("%s:> get cluster info err: %v", logPrefix, err)
	}

	cl.Status = *getClusterStatus(nodes)

	err = envs.Get().GetStorage().Upsert(context.Background(), storage.ClusterKind, types.EmptyString, cl, nil)
	if err != nil {
		log.Fatalf("%s:> set cluster status err: %v", logPrefix, err)
	}

	return sc
}
