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
	"github.com/lastbackend/lastbackend/pkg/cache"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/scheduler/envs"
	"reflect"
)

type Controller struct {
	node   chan *types.Node
	active bool

	cache *cache.NodeCache
}

func (nc *Controller) Watch(node chan *types.Node) {

	var (
		stg   = envs.Get().GetStorage()
		event = make(chan *types.Event)
	)

	log.Debug("PodController: start watch")
	go func() {
		for {
			select {
			case n := <-nc.node:
				{
					if !nc.active {
						log.Debug("NodeController: skip management cause it is in slave mode")
						continue
					}

					log.Debugf("Node check state: %s", n.Meta.Name)

					item := nc.cache.Get(n.Info.Hostname)
					if item == nil || !reflect.DeepEqual(item, n) {
						nc.cache.Set(n.Info.Hostname, n)

						nodes := nc.cache.List()

						err := stg.Cluster().SetStatus(context.Background(), getClusterStatus(nodes))
						if err != nil {
							log.Debug("NodeController: set cluster status err: %s", err.Error())
							continue
						}

					}

					if n.Online {
						log.Debugf("Node set alive, try to provision on it pods: %s", n.Meta.Name)
						node <- n
						continue
					}

					log.Debugf("Node set offline, try to move all pods to another")

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

				nc.node <- e.Data.(*types.Node)
			}
		}
	}()

	stg.Node().Watch(context.Background(), event)
}

func (nc *Controller) Pause() {
	nc.active = false
}

func (nc *Controller) Resume() {
	nc.active = true
	log.Debug("NodeController: start check pods state")
}

func NewNodeController(ctx context.Context) *Controller {

	sc := new(Controller)
	sc.active = false
	sc.node = make(chan *types.Node)
	sc.cache = cache.NewNodeCache()

	nodes, err := envs.Get().GetStorage().Node().List(ctx)
	if err != nil {
		log.Fatalf("NodeController: get nodes list err: %s", err.Error())
	}

	for _, node := range nodes {
		sc.cache.Set(node.Info.Hostname, node)
	}

	err = envs.Get().GetStorage().Cluster().SetStatus(context.Background(), getClusterStatus(nodes))
	if err != nil {
		log.Fatalf("NodeController: set cluster status err: %s", err.Error())
	}

	return sc
}
