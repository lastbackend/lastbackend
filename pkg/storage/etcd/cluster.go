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
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd/cache"
	"time"
	"regexp"
	"encoding/json"
)

type ClusterStorage struct {
	storage.Cluster
	cache *cache.Cache
}

const clusterStorage = "cluster"

// SetStatus - set cluster status object into storage
func (s *ClusterStorage) SetStatus(ctx context.Context, status *types.ClusterStatus) error {

	log.V(logLevel).Debugf("storage:etcd:cluster:> set status: %v", status)

	if status == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:cluster:> set status err: %s", err.Error())
		return err
	}
	defer destroy()

	keyStatus := keyCreate(clusterStorage, "status")
	if err := client.Upsert(ctx, keyStatus, status, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:cluster:> set status err: %s", err.Error())
		return err
	}

	return nil
}

// Get - return  cluster info from storage
func (s *ClusterStorage) Get(ctx context.Context) (*types.Cluster, error) {

	log.V(logLevel).Debug("storage:etcd:cluster:> get status")

	const filter = `\b.+` + clusterStorage + `\/(status)\b`

	client, destroy, err := getClient(ctx)
	if err != nil {
		return nil, err
	}
	defer destroy()

	cluster := new(types.Cluster)
	key := keyCreate(clusterStorage)
	if err := client.Map(ctx, key, filter, cluster); err != nil {
		log.V(logLevel).Errorf("storage:etcd:cluster:> get err: %s", err.Error())
		return nil, err
	}

	return cluster, nil
}

// Watch cluster changes
func (s *ClusterStorage) Watch(ctx context.Context, event chan *types.Event) error {

	log.V(logLevel).Debug("storage:etcd:service:> watch cluster")

	const filter = `\b.+` + clusterStorage + `\/(.+)\b`
	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:cluster:> watch cluster err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(clusterStorage)
	cb := func(action, key string, data []byte) {

		keys := r.FindStringSubmatch(key)
		if len(keys) < 1 {
			return
		}

		e := new(types.Event)
		e.Action = action
		e.Name = "lb"

		if action == store.STORAGEDELETEEVENT {
			e.Data = nil
			event <- e
			return
		}

		item := s.cache.Get("cluster")

		if item == nil {
			if data, err := s.Get(ctx); err == nil {
				s.cache.Set("cluster", data)
				e.Data = data
				event <- e
			}
			return
		}

		cl := item.(*types.Cluster)

		switch keys[1] {
		case "status":
			var status types.ClusterStatus
			if err := json.Unmarshal(data, &status); err != nil {
				log.V(logLevel).Errorf("storage:etcd:cluster:> parse cluster status err: %s", err.Error())
				return
			}
			cl.Status = status
		}

		s.cache.Set("cluster", cl)

		e.Data = cl

		event <- e
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("storage:etcd:cluster:> watch cluster err: %s", err.Error())
		return err
	}

	return nil
}

// Clear database stare
func (s *ClusterStorage) Clear(ctx context.Context) error {

	log.V(logLevel).Debugf("storage:etcd:cluster:> clear")

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:cluster:> clear err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := client.DeleteDir(ctx, clusterStorage); err != nil {
		log.V(logLevel).Errorf("storage:etcd:cluster:> clear err: %s", err.Error())
		return err
	}

	return nil
}

// newClusterStorage - return new cluster interface
func newClusterStorage() *ClusterStorage {
	s := new(ClusterStorage)
	s.cache = cache.NewCache(24 * time.Hour)
	return s
}
