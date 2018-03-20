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

package v1

import (
	"context"

	"github.com/lastbackend/lastbackend/pkg/api/client/http"
	"github.com/lastbackend/lastbackend/pkg/api/client/interfaces"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type ClusterClient struct {
	interfaces.Cluster
	client http.Interface
}

func (s *ClusterClient) Node(hostname ...string) *NodeClient {
	hst := ""
	if len(hostname) > 0 {
		hst = hostname[0]
	}
	return newNodeClient(s.client, hst)
}

func (s *ClusterClient) Get(ctx context.Context) (*vv1.ClusterList, error) {
	return nil, nil
}

func (s *ClusterClient) Update(ctx context.Context, opts *rv1.ClusterUpdateOptions) (*vv1.Cluster, error) {
	return nil, nil
}

func newClusterClient(req http.Interface) *ClusterClient {
	return &ClusterClient{client: req}
}
