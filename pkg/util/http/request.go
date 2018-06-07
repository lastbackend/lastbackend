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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	_url "github.com/lastbackend/lastbackend/pkg/util/url"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	contentType     = "Content-Type"
	jsonContentType = "application/json"
	formContentType = "application/x-www-form-urlencoded"
)

type RawReq struct {
	host     string
	port     int
	rawURL   string
	method   string
	tls      bool
	header   http.Header
	bodyJSON interface{}
	body     io.ReadCloser
}

type ReqOpts struct {
	TLS bool
}

func New(host string, opts *ReqOpts) (*RawReq, error) {
	raw := &RawReq{
		host:   host,
		header: http.Header{},
	}

	if opts != nil {
		raw.tls = opts.TLS
	}

	u, err := _url.Parse(raw.host)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	raw.host = u.String()

	return raw, nil
}

func (r *RawReq) Request(successV, failureV interface{}) (req *http.Request, resp *http.Response, err error) {
	req, err = r.getRequest()
	if err != nil {
		return nil, nil, err
	}

	req.Cookies()

	client := http.DefaultClient
	if r.tls {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	err = decodeResponseJSON(resp, successV, failureV)
	if err != nil {
		return nil, nil, err
	}

	return req, resp, nil
}

func (r *RawReq) Do() (req *http.Request, resp *http.Response, err error) {

	req, err = r.getRequest()
	if err != nil {
		return nil, nil, err
	}

	req.Cookies()

	client := http.DefaultClient
	if r.tls {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	return req, resp, nil
}

func (r *RawReq) AddHeader(key, value string) *RawReq {
	r.header.Set(key, value)
	return r
}

func (r *RawReq) getRequest() (*http.Request, error) {

	var err error

	reqURL, err := url.Parse(r.rawURL)
	if err != nil {
		return nil, err
	}

	reqMethod := r.method

	reqBody, err := r.getRequestBody()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(reqMethod, reqURL.String(), reqBody)
	if err != nil {
		return nil, err
	}

	req.Header = r.header

	return req, nil
}

func (r RawReq) clear() (err error) {

	err = r.body.Close()
	if err != nil {
		return err
	}

	r.tls = false
	r.method = ""
	r.rawURL = ""
	r.bodyJSON = nil
	r.header = nil
	r.body = nil

	return nil
}

func decodeResponseJSON(resp *http.Response, successV, failureV interface{}) error {

	if code := resp.StatusCode; 200 <= code && code <= 299 {
		if successV != nil {
			return decodeResponseBodyJSON(resp, successV)
		}
	} else {
		if failureV != nil {
			switch strings.Split(resp.Header.Get("Content-type"), ";")[0] {
			case "text/html":
				buf := new(bytes.Buffer)
				buf.ReadFrom(resp.Body)
			case "text/plain":
				buf := new(bytes.Buffer)
				buf.ReadFrom(resp.Body)
			case "application/json":
				return decodeResponseBodyJSON(resp, failureV)
			default:
				return errors.New(fmt.Sprintf("Unknown content-type (%+v)", strings.Split(resp.Header.Get("Content-type"), ";")))
			}
		}
	}

	return nil
}

func decodeResponseBodyJSON(resp *http.Response, v interface{}) error {
	err := json.NewDecoder(resp.Body).Decode(v)
	if err != nil && io.EOF == err {
		return nil
	}
	return err
}
