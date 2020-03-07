//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/internal/util/http/request"
	"github.com/lastbackend/lastbackend/pkg/client/types"
)

type Client struct {
	client *request.RESTClient
}

func New(req *request.RESTClient) *Client {
	return &Client{client: req}
}

func (s *Client) Cluster() types.ClusterClientV1 {
	return newClusterClient(s.client)
}

func (s *Client) Namespace(args ...string) types.NamespaceClientV1 {
	name := ""
	// Get any parameters passed to us out of the args variable into "real"
	// variables we created for them.
	for i := range args {
		switch i {
		case 0: // hostname
			name = args[0]
		default:
			panic("Wrong parameters count: (is allowed from 0 to 1)")
		}
	}
	return newNamespaceClient(s.client, name)
}
