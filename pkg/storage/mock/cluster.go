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

package mock

import (
	"context"
	"errors"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

// ClusterStorage - mock storage for cluster
type ClusterStorage struct {
	storage.Cluster
	data types.Cluster
}

// Insert - insert new cluster object into mock storage
func (s *ClusterStorage) Insert(ctx context.Context, cluster *types.Cluster) error {

	if cluster == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	s.data = *cluster
	return nil
}

// Info - return  cluster info from mock storage
func (s *ClusterStorage) Get(ctx context.Context) (*types.Cluster, error) {
	return &s.data, nil
}

// Update cluster info into mock storage
func (s *ClusterStorage) Update(ctx context.Context, cluster *types.Cluster) error {
	if cluster == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	s.data = *cluster
	return nil
}

// newClusterStorage - return new mock cluster interface
func newClusterStorage() *ClusterStorage {
	s := new(ClusterStorage)
	return s
}
