//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/http2"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Request struct {
	// required
	client HTTPClient
	verb   string

	baseURL *url.URL

	pathPrefix string
	params     url.Values
	headers    http.Header

	err  error
	body io.Reader

	ctx context.Context
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func New(client HTTPClient, method string, baseURL *url.URL) *Request {
	pathPrefix := "/"

	if baseURL != nil {
		pathPrefix = path.Join(pathPrefix, baseURL.Path)
	}

	r := &Request{
		client:     client,
		verb:       method,
		baseURL:    baseURL,
		pathPrefix: pathPrefix,
	}

	return r
}

func (r *Request) Body(data []byte) *Request {
	r.body = bytes.NewReader(data)
	return r
}

func (r *Request) AddHeader(key, val string) *Request {
	if r.headers == nil {
		r.headers = make(map[string][]string)
	}
	r.headers.Add(key, val)
	return r
}

func (r *Request) Do() Result {
	var result Result

	if r.err != nil {
		return Result{err: r.err}
	}

	client := r.client
	if client == nil {
		panic("client not initialized")
	}

	u := r.URL().String()
	req, err := http.NewRequest(r.verb, u, r.body)
	if err != nil {
		return Result{err: r.err}
	}

	if r.ctx != nil {
		req = req.WithContext(r.ctx)
	}

	req.Header = r.headers

	resp, err := client.Do(req)
	if err != nil {
		return Result{err: r.err}
	}

	result = r.transformResponse(resp, req)

	return result
}

func (r *Request) JSON(success interface{}, failure interface{}) error {

	client := r.client
	if client == nil {
		panic("client not initialized")
	}

	u := r.URL().String()
	req, err := http.NewRequest(r.verb, u, r.body)
	if err != nil {
		return err
	}

	if r.ctx != nil {
		req = req.WithContext(r.ctx)
	}

	req.Header = r.headers

	resp, err := client.Do(req)
	if err != nil {
		if io.EOF == err || strings.Contains(err.Error(), "EOF") {
			return nil
		}
		return err
	}

	return decodeJSON(resp, success, failure)
}

type Result struct {
	body        []byte
	contentType string
	err         error
	statusCode  int
}

// Raw returns the raw result.
func (r Result) Raw() ([]byte, error) {
	return r.body, r.err
}

// Raw returns the raw result.
func decodeJSON(r *http.Response, success interface{}, failure interface{}) error {

	if code := r.StatusCode; 200 > code || code > 299 {
		if failure == nil {
			return nil
		}
		return decodeResponseJSON(r, failure)
	}

	if success == nil {
		return nil
	}

	return decodeResponseJSON(r, success)
}

func (r Result) StatusCode() int {
	return r.statusCode
}

func (r Result) Error() error {
	return r.err
}

func (r *Request) Stream() (io.ReadCloser, *http.Response, error) {
	if r.err != nil {
		return nil, nil, r.err
	}

	u := r.URL().String()

	req, err := http.NewRequest(r.verb, u, nil)
	if err != nil {
		return nil, nil, err
	}
	if r.ctx != nil {
		req = req.WithContext(r.ctx)
	}

	req.Header = r.headers

	client := r.client
	if client == nil {
		client = http.DefaultClient
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	switch {
	case (res.StatusCode >= 200) && (res.StatusCode < 400):
		return res.Body, res, nil

	default:
		defer res.Body.Close()

		result := r.transformResponse(res, req)
		err := result.Error()
		if err == nil {
			err = fmt.Errorf("%d while accessing %v: %s", result.statusCode, u, string(result.body))
		}
		return nil, res, err
	}
}

func (r *Request) Param(name, value string) *Request {
	if r.params == nil {
		r.params = make(url.Values)
	}
	r.params[name] = append(r.params[name], value)
	return r
}

func (r *Request) transformResponse(resp *http.Response, req *http.Request) Result {
	var body []byte

	if resp.Body != nil {
		data, err := ioutil.ReadAll(resp.Body)
		switch err.(type) {
		case nil:
			body = data
		case http2.StreamError:
			return Result{
				err: fmt.Errorf("stream error %#v when reading", err),
			}
		default:
			return Result{
				err: fmt.Errorf("unexpected error %#v", err),
			}
		}
	}

	return Result{
		body:        body,
		contentType: resp.Header.Get("Content-Type"),
		statusCode:  resp.StatusCode,
	}
}

func decodeResponseJSON(r *http.Response, v interface{}) error {

	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil && io.EOF == err {
		return nil
	}
	return err
}

func (r *Request) URL() *url.URL {
	p := r.pathPrefix

	finalURL := &url.URL{}
	if r.baseURL != nil {
		*finalURL = *r.baseURL
	}
	finalURL.Path = p

	query := url.Values{}
	for key, values := range r.params {
		for _, value := range values {
			query.Add(key, value)
		}
	}

	finalURL.RawQuery = query.Encode()

	return finalURL
}
