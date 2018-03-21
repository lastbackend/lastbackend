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
	"github.com/lastbackend/lastbackend/pkg/util/serializer"
	"net/http"
	"net/url"
	"strings"
	"time"
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
}

func NewRESTClient(baseURL *url.URL) (*RESTClient, error) {
	base := *baseURL

	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}

	return &RESTClient{
		base: &base,
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
	}, nil
}

func (c *RESTClient) Do(verb string, path string) *Request {
	c.base.Path = path
	return NewRequest(c.Client, verb, c.base)
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
