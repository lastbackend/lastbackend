//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

import "net/http"

const (
	MethodGet     = http.MethodGet
	MethodHead    = http.MethodHead
	MethodPost    = http.MethodPost
	MethodPut     = http.MethodPut
	MethodPatch   = http.MethodPatch
	MethodDelete  = http.MethodDelete
	MethodConnect = http.MethodConnect
	MethodOptions = http.MethodOptions
	MethodTrace   = http.MethodTrace
)

func (r *RawReq) POST(pathURL string) *RawReq {
	r.method = MethodPost
	r.rawURL = r.host + pathURL
	return r
}

func (r *RawReq) GET(pathURL string) *RawReq {
	r.method = MethodGet
	r.rawURL = r.host + pathURL
	return r
}

func (r *RawReq) PUT(pathURL string) *RawReq {
	r.method = MethodPut
	r.rawURL = r.host + pathURL
	return r
}

func (r *RawReq) DELETE(pathURL string) *RawReq {
	r.method = MethodDelete
	r.rawURL = r.host + pathURL
	return r
}
