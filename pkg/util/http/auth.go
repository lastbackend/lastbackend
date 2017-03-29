package http

import (
	"github.com/gorilla/mux"
	"strings"
	"net/http"
	"github.com/lastbackend/lastbackend/pkg/api/types"
	"github.com/gorilla/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
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

		// Add session and token to context
		context.Set(r, "token", token)
		context.Set(r, "session", s)

		h.ServeHTTP(w, r)
	}
}