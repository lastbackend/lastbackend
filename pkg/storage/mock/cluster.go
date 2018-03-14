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

package mock

import (
	"context"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
)

// ClusterStorage database mock structure
type ClusterStorage struct {
	storage.Cluster
}

// Info - return cluster info from storage
func (*ClusterStorage) Info(ctx context.Context) (*types.Cluster, error) {
	panic("implement me")
}

// Update - update cluster info in storage
func (*ClusterStorage) Update(ctx context.Context, cluster *types.Cluster) error {
	panic("implement me")
}
