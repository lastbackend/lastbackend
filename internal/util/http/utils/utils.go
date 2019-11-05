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

package utils

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/util/converter"
	"net/http"
)

func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func SetContext(r *http.Request, name string, val interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), name, val))
}

func QueryString(r *http.Request, param string) string {
	return r.URL.Query().Get(param)
}

func QueryFloat(r *http.Request, param string) float64 {
	return converter.StringToFloat(r.URL.Query().Get(param))
}

func QueryInt(r *http.Request, param string) int64 {
	return converter.StringToInt64(r.URL.Query().Get(param))
}

func QueryBool(r *http.Request, param string) bool {
	return converter.StringToBool(r.URL.Query().Get(param))
}
