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
	"time"
)

type ClusterStorage struct {
	storage.Cluster
}

const clusterStorage = "cluster"

// Insert - insert new cluster object into mock storage
func (s *ClusterStorage) Insert(ctx context.Context, cluster *types.Cluster) error {

	log.V(logLevel).Debugf("storage:etcd:cluster:> insert: %v", cluster)

	if cluster == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:cluster:> insert err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(clusterStorage, "meta")
	if err := client.Create(ctx, keyMeta, cluster.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:cluster:> insert err: %s", err.Error())
		return err
	}

	keyState := keyCreate(clusterStorage, "state")
	if err := client.Create(ctx, keyState, cluster.State, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:cluster:> insert err: %s", err.Error())
		return err
	}

	return nil
}

// Get - return  cluster info from mock storage
func (s *ClusterStorage) Get(ctx context.Context) (*types.Cluster, error) {

	log.V(logLevel).Debug("storage:etcd:cluster:> get meta|state")

	const filter = `\b.+` + clusterStorage + `\/(?:meta|state)\b`

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

// Update cluster info into mock storage
func (s *ClusterStorage) Update(ctx context.Context, cluster *types.Cluster) error {

	if cluster == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	log.V(logLevel).Debugf("storage:etcd:cluster:> update: %v", cluster)

	cluster.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:cluster:> update err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMeta := keyCreate(clusterStorage, "meta")
	if err := client.Update(ctx, keyMeta, cluster.Meta, nil, 0); err != nil {
		log.V(logLevel).Errorf("storage:etcd:cluster:> update err: %s", err.Error())
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

// newClusterStorage - return new mock cluster interface
func newClusterStorage() *ClusterStorage {
	s := new(ClusterStorage)
	return s
}
