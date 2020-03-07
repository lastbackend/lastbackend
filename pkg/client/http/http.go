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

package http

import (
	"github.com/lastbackend/lastbackend/internal/util/http/request"
	"github.com/lastbackend/lastbackend/pkg/client/config"
	"github.com/lastbackend/lastbackend/pkg/client/http/v1"
	"github.com/lastbackend/lastbackend/pkg/client/types"
	"github.com/lastbackend/lastbackend/tools/log"
)

type Client struct {
	client *request.RESTClient
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
	opts.Headers = make(map[string]string, 0)

	if cfg.TLS != nil {
		opts.TLS = new(request.TLSConfig)
		opts.TLS.Insecure = !cfg.TLS.Verify
		opts.TLS.ServerName = cfg.TLS.ServerName
		opts.TLS.CertFile = cfg.TLS.CertFile
		opts.TLS.KeyFile = cfg.TLS.KeyFile
		opts.TLS.CAFile = cfg.TLS.CAFile
		opts.TLS.CAData = cfg.TLS.CAData
		opts.TLS.CertData = cfg.TLS.CertData
		opts.TLS.KeyData = cfg.TLS.KeyData
	}

	if cfg.Headers == nil {
		cfg.Headers = make(map[string]string, 0)
	}

	for k, v := range cfg.Headers {
		opts.Headers[k] = v
	}

	client, err := request.NewRESTClient(endpoint, opts)
	if err != nil {
		log.Errorf("can not initialize client: %s", err.Error())
	}

	cl.client = client

	return cl, nil
}

func (c Client) V1() types.ClientV1 {
	return v1.New(c.client)
}
