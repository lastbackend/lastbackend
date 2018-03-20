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

package v1

import (
	"context"

	"github.com/lastbackend/lastbackend/pkg/api/client/http"
	"github.com/lastbackend/lastbackend/pkg/api/client/interfaces"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type NodeClient struct {
	interfaces.Node
	client http.Interface
}

func (s *NodeClient) List(ctx context.Context) (*vv1.NodeList, error) {
	return nil, nil
}

func (s *NodeClient) Get(ctx context.Context, name string) (*vv1.Node, error) {
	return nil, nil
}

func (s *NodeClient) Update(ctx context.Context, name string, opts rv1.NodeUpdateOptions) (*vv1.Node, error) {
	return nil, nil
}

func (s *NodeClient) Remove(ctx context.Context, name string, opts rv1.NodeRemoveOptions) error {
	return nil
}

func newNodeClient(req http.Interface) *NodeClient {
	return &NodeClient{client: req}
}
