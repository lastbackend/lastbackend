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

package middleware

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
	"strings"
)

// Auth - authentication middleware
func Authenticate(ctx context.Context, h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var accessToken = types.EmptyString
		var token = types.EmptyString
		var params = utils.Vars(r)

		if ctx.Value("access_token") != nil {
			accessToken = ctx.Value("access_token").(string)
		}

		if _, ok := r.URL.Query()["x-lastbackend-token"]; ok {
			token = r.URL.Query().Get("x-lastbackend-token")
		} else if _, ok := params["x-lastbackend-token"]; ok {
			token = params["x-lastbackend-token"]
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
			errors.HTTP.Unauthorized(w)
			return
		}

		if token != accessToken {
			errors.HTTP.Unauthorized(w)
			return
		}

		h.ServeHTTP(w, r)
	}
}
