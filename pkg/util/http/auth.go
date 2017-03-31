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

import (
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/errors"
	c "golang.org/x/net/context"
	"net/http"
	"strings"
)

// Auth - authentication middleware
func Authenticate(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var token string
		var params = mux.Vars(r)

		if _, ok := params["token"]; ok {
			token = params["token"]
		} else if r.Header.Get("Authorization") != "" {
			// Parse authorization header
			var auth = strings.SplitN(r.Header.Get("Authorization"), " ", 2)

			// Check authorization header parts length and authorization header format
			if len(auth) != 2 || auth[0] != "Bearer" {
				errors.HTTP.Unauthorized(w)
				return
			}

			token = auth[1]

		} else {
			w.Header().Set("Content-Type", "application/json")
			errors.HTTP.Unauthorized(w)
			return
		}

		s := new(types.Session)
		err := s.Decode(token)
		if err != nil {
			errors.HTTP.Unauthorized(w)
			return
		}

		ctx := c.WithValue(r.Context(), "token", token)
		ctx = c.WithValue(ctx, "session", s)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	}
}