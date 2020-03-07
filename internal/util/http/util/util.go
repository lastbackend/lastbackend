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

package util

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/util/converter"
)

// Vars is a helper function that returns URL parameters as map[string]string from request
func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

// QueryString is a helper function that returns querystring parameter as string value
func QueryString(r *http.Request, param string) string {
	return r.URL.Query().Get(param)
}

// QueryFloat is a helper function that returns querystring parameter as float64 value
func QueryFloat(r *http.Request, param string) float64 {
	return converter.StringToFloat(r.URL.Query().Get(param))
}

// QueryInt is a helper function that returns querystring parameter as int64 value
func QueryInt64(r *http.Request, param string) int64 {
	return converter.StringToInt64(r.URL.Query().Get(param))
}

// QueryInt is a helper function that returns querystring parameter as int64 value
func QueryInt(r *http.Request, param string) int {
	return converter.StringToInt(r.URL.Query().Get(param))
}

// QueryBool is a helper function that returns querystring parameter as bool value
func QueryBool(r *http.Request, param string) bool {
	return converter.StringToBool(r.URL.Query().Get(param))
}

// SetContext is a helper function that returns
func SetContext(r *http.Request, name string, val interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), name, val))
}

// GetURL is a helper function that returns the URL as a string
// (the scheme + hostname; without the path)
func GetURL(r *http.Request) string {
	return GetScheme(r) + "://" + GetHost(r)
}

// GetScheme is a helper function returns the scheme, HTTP or HTTPS.
// It is able to detect, using the X-Forwarded-Proto,
// if the original request was HTTPS and routed
// through a reverse proxy with TLS termination.
func GetScheme(r *http.Request) string {
	switch {
	case r.URL.Scheme == "https":
		return "https"
	case r.TLS != nil:
		return "https"
	case strings.HasPrefix(r.Proto, "HTTPS"):
		return "https"
	case r.Header.Get("X-Forwarded-Proto") == "https":
		return "https"
	default:
		return "http"
	}
}

// GetHost is a helper function returns the hostname.
// It is able to detect, using the X Forwarded-For header,
// the original hostname when routed through a reverse proxy.
func GetHost(r *http.Request) string {
	switch {
	case len(r.Host) != 0:
		return r.Host
	case len(r.URL.Host) != 0:
		return r.URL.Host
	case len(r.Header.Get("X-Forwarded-For")) != 0:
		return r.Header.Get("X-Forwarded-For")
	case len(r.Header.Get("X-Host")) != 0:
		return r.Header.Get("X-Host")
	case len(r.Header.Get("XFF")) != 0:
		return r.Header.Get("XFF")
	case len(r.Header.Get("X-Real-IP")) != 0:
		return r.Header.Get("X-Real-IP")
	default:
		return "localhost:8080"
	}
}
