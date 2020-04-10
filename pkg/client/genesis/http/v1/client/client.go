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

package client

import (
	"github.com/lastbackend/lastbackend/pkg/client/genesis/types"
	"github.com/lastbackend/lastbackend/internal/util/http/request"
)

type Client struct {
	client *request.RESTClient
}

func New(req *request.RESTClient) *Client {
	return &Client{client: req}
}

func (s *Client) Registry() types.RegistryClientV1 {
	return newRegistryClient(s.client)
}

func (s *Client) Account() types.AccountClientV1 {
	return newAccountClient(s.client)
}

func (s *Client) Cluster() types.ClusterClientV1 {
	return newClusterClient(s.client)
}
