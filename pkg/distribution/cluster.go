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
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const (
	logClusterPrefix = "distribution:cluster"
)

// Cluster - distribution model
type Cluster struct {
	context context.Context
	storage storage.Storage
}

// Info - get cluster info
func (c *Cluster) Get() (*types.Cluster, error) {

	log.V(logLevel).Debugf("%s:get:> get info", logClusterPrefix)

	cluster, err := c.storage.Cluster().Get(c.context)
	if err != nil {
		if err.Error() == store.ErrEntityNotFound {
			log.V(logLevel).Warnf("%s:get:> cluster not found", logClusterPrefix)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> get cluster err: %s", logClusterPrefix, err.Error())
		return nil, err
	}

	return cluster, nil
}

// Update - update cluster stats data and meta information
func (c *Cluster) Update(cluster *types.Cluster, opts *request.ClusterUpdateOptions) error {

	log.V(logLevel).Debugf("%s:update:> update cluster %#v", logClusterPrefix, cluster)

	if opts.Description != nil {
		cluster.Meta.Description = *opts.Description
	}

	if err := c.storage.Cluster().Update(c.context, cluster); err != nil {
		log.V(logLevel).Errorf("%s:update:> update cluster err: %s", logClusterPrefix, err.Error())
		return err
	}

	return nil
}

// NewClusterModel - return new cluster model
func NewClusterModel(ctx context.Context, stg storage.Storage) *Cluster {
	return &Cluster{ctx, stg}
}
