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

// SetStatus - set cluster status object into mock storage
func (s *ClusterStorage) SetStatus(ctx context.Context, status *types.ClusterStatus) error {

	if status == nil {
		return errors.New(store.ErrStructArgIsNil)
	}

	s.data.Status = *status
	return nil
}

// Get - return  cluster info from mock storage
func (s *ClusterStorage) Get(ctx context.Context) (*types.Cluster, error) {
	return &s.data, nil
}

// Watch cluster changes
func (s *ClusterStorage) Watch(ctx context.Context, event chan *types.Event) error {
	return nil
}

// Clear database stare
func (s *ClusterStorage) Clear(ctx context.Context) error {
	s.data = types.Cluster{}
	return nil
}

// newClusterStorage - return new mock cluster interface
func newClusterStorage() *ClusterStorage {
	s := new(ClusterStorage)
	return s
}
