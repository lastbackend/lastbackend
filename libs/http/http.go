package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	rawURL   string
	method   string
	header   http.Header
	bodyJSON interface{}
	body     io.ReadCloser
}

func New(host string) *RawReq {
	return &RawReq{
		host:   host,
		header: http.Header{},
	}
}

func (r *RawReq) Request(successV, failureV interface{}) (req *http.Request, resp *http.Response, err error) {

	req, err = r.getRequest()
	if err != nil {
		return nil, nil, err
	}

	req.Cookies()

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	err = decodeResponseJSON(resp, successV, failureV)
	if err != nil {
		return nil, nil, err
	}

	return req, resp, err
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
				s := buf.String()
				fmt.Println(s)
			case "text/plain":
				buf := new(bytes.Buffer)
				buf.ReadFrom(resp.Body)
				s := buf.String()
				fmt.Println(s)
			case "application/json":
				return decodeResponseBodyJSON(resp, failureV)
			//fmt.Printf("%+v", failureV)
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

func (r RawReq) clear() (err error) {

	err = r.body.Close()
	if err != nil {
		return err
	}

	r.method = ""
	r.rawURL = ""
	r.bodyJSON = nil
	r.header = nil
	r.body = nil

	return nil
}
