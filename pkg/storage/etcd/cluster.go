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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

type ClusterStorage struct {
	storage.Cluster
}

const clusterStorage = "cluster"

func (s *ClusterStorage) Info(ctx context.Context) (*types.Cluster, error) {

	log.V(logLevel).Debug("Storage: cluster: info")

	const filter = `\b.+` + clusterStorage + `\/.+\/(?:meta|state)\b`

	client, destroy, err := getClient(ctx)
	if err != nil {
		return nil, err
	}
	defer destroy()

	cluster := new(types.Cluster)
	key := keyCreate(clusterStorage)
	if err := client.Map(ctx, key, filter, cluster); err != nil {
		if err.Error() == store.ErrEntityNotFound {
			return nil, nil
		}
		log.V(logLevel).Errorf("Storage: Node: create client err: %s", err.Error())
		return nil, err
	}

	return cluster, nil
}

func (s *ClusterStorage) Update(ctx context.Context, cluster *types.Cluster) error {

	log.V(logLevel).Debugf("Storage: cluster: update: #v", cluster)

	return nil
}

func newClusterStorage() *ClusterStorage {
	s := new(ClusterStorage)
	return s
}
