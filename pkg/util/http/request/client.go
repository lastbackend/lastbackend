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

package request

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"github.com/pkg/errors"
	"time"

	"github.com/lastbackend/lastbackend/pkg/util/serializer"
)

type IRESTClient interface {
	Do(verb string, path string) *Request
	Post(path string) *Request
	Put(path string) *Request
	Get(path string) *Request
	Delete(path string) *Request
}

type RESTClient struct {
	base        *url.URL
	serializers serializer.Codec
	Client      *http.Client
	BearerToken string
}

func NewRESTClient(uri string, cfg *Config) (*RESTClient, error) {

	if cfg == nil {
		return DefaultRESTClient(uri), nil
	}

	c := new(RESTClient)
	c.base = parseURL(uri)

	c.Client = new(http.Client)

	c.Client.Timeout = cfg.Timeout * time.Second
	c.Client.Transport = new(http.Transport)
	c.BearerToken = cfg.BearerToken

	if cfg.TLS != nil && !cfg.TLS.Insecure {
		if err := withTLSClientConfig(cfg)(c.Client); err != nil {
			return nil, err
		}
		c.base.Scheme = "https"
	}

	return c, nil
}

func DefaultRESTClient(uri string) *RESTClient {
	return &RESTClient{
		base: parseURL(uri),
		Client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: new(http.Transport),
		},
	}
}

// WithTLSClientConfig applies a tls config to the client transport.
func withTLSClientConfig(cfg *Config) func(*http.Client) error {
	return func(c *http.Client) error {

		tc, err := NewTLSConfig(cfg)
		if err != nil {
			return errors.Wrap(err, "failed to create tls config")
		}

		if transport, ok := c.Transport.(*http.Transport); ok {
			transport.TLSClientConfig = tc
			return nil
		}

		return errors.Errorf("cannot apply tls config to transport: %T", c.Transport)
	}
}

func (c *RESTClient) Do(verb string, path string) *Request {
	c.base.Path = path
	req := New(c.Client, verb, c.base)
	if len(c.BearerToken) != 0 {
		req.AddHeader("Authorization", fmt.Sprintf("Bearer %s", string(c.BearerToken)))
	}
	return req
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

func parseURL(u string) *url.URL {

	uri, err := url.Parse(u)
	if err != nil || uri.Scheme == "" || uri.Host == "" {
		scheme := "http://"
		uri, err = url.Parse(scheme + u)
		if err != nil {
			return nil
		}
	}

	if !strings.HasSuffix(uri.Path, "/") {
		uri.Path += "/"
	}

	return uri
}
