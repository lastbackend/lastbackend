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
	"github.com/lastbackend/lastbackend/pkg/api/client/http"
)

type Client struct {
	client http.Interface
}

func (s *Client) Cluster() *ClusterClient {
	if s == nil {
		return nil
	}
	return newClusterClient(s.client)
}

func (s *Client) Namespace(name ...string) *NamespaceClient {
	if s == nil {
		return nil
	}
	n := ""
	if len(name) > 0 {
		n = name[0]
	}
	return newNamespaceClient(s.client, n)
}

func New(req http.Interface) *Client {
	return &Client{client: req}
}
