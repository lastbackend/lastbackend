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
		//
		//// Add session and token to context
		//context.Set(r, "token", token)
		//context.Set(r, "session", s)

		ctx := c.WithValue(r.Context(), "token", token)
		ctx = c.WithValue(ctx, "session", s)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	}
}
