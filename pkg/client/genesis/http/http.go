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

package http

import (
	"github.com/lastbackend/lastbackend/pkg/client/genesis/types"
	"github.com/lastbackend/lastbackend/internal/util/http/request"
	"github.com/lastbackend/lastbackend/pkg/client/genesis/config"
	"github.com/lastbackend/lastbackend/pkg/client/genesis/http/v1/client"
	rr "github.com/lastbackend/lastbackend/pkg/client/genesis/http/v1/request"
	"github.com/lastbackend/lastbackend/pkg/client/genesis/http/v1/views"
	"github.com/lastbackend/lastbackend/tools/log"
)

type Client struct {
	client *request.RESTClient
}

func View() views.IView {
	return &views.View{}
}

func Request() rr.IRequest {
	return &rr.Request{}
}

func New(endpoint string, cfg *config.Config) (*Client, error) {

	cl := new(Client)

	if cfg == nil {
		cl.client = request.DefaultRESTClient(endpoint)
		return cl, nil
	}

	opts := new(request.Config)
	opts.BearerToken = cfg.BearerToken
	opts.Timeout = cfg.Timeout

	cli, err := request.NewRESTClient(endpoint, opts)
	if err != nil {
		log.Errorf("can not initialize client: %s", err.Error())
	}

	cl.client = cli

	return cl, nil
}

func (c Client) V1() types.ClientV1 {
	return client.New(c.client)
}
