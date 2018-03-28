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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/util/serializer"
	"net/http"
	"net/url"
	"strings"
	"time"
	"crypto/tls"
)

// Interface captures the set of operations for generically interacting with Kubernetes REST apis.
type Interface interface {
	Do(verb, path string) *Request
	Post(path string) *Request
	Put(path string) *Request
	Get(path string) *Request
	Delete(path string) *Request
}

type RESTClient struct {
	base        *url.URL
	serializers serializer.Codec
	Client      *http.Client
	config      *Config
	token       string
}

func NewRESTClient(baseURL *url.URL, config *Config) (*RESTClient, error) {

	if config == nil {
		return nil, errors.New("config not set")
	}

	base := *baseURL

	if base.Scheme == "" {
		base.Scheme = "http"
	}

	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}

	return &RESTClient{
		base:   &base,
		config: config,
		Client: &http.Client{
			Timeout: time.Second * config.Timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Insecure},
			},
		},
	}, nil
}

func (c *RESTClient) Do(verb string, path string) *Request {
	c.base.Path = path
	request := NewRequest(c.Client, verb, c.base)
	if len(c.config.BearerToken) != 0 {
		request.AddHeader("Authorization", fmt.Sprintf("Bearer %s", string(c.config.BearerToken)))
	}
	return request
}

func (c *RESTClient) Post(path string) *Request {
	return c.Do(http.MethodPost, path)
}

func (c *RESTClient) Put(path string) *Request {
	return c.Do(http.MethodPut, path)
}

func (c *RESTClient) Get(path string) *Request {
	return c.Do(http.MethodGet, path)
}

func (c *RESTClient) Delete(path string) *Request {
	return c.Do(http.MethodDelete, path)
}
