package middleware

import (
	"net/http"
	"strings"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/util/http/util"
	"github.com/lastbackend/lastbackend/tools/logger"
)

// Authenticate - authentication middleware
func (m Middleware) Authenticate(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ctx := logger.NewContext(r.Context(), nil)
		log := logger.WithContext(ctx)

		var (
			secret = m.viper.GetString("security.token")
			params = util.Vars(r)
			token  string
		)

		if _, ok := r.URL.Query()["x-access-token"]; ok {
			token = r.URL.Query().Get("x-access-token")
		} else if _, ok := params["x-access-token"]; ok {
			token = params["x-access-token"]
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
			log.Errorf("%s:authenticate:> token not set", logPrefix)
			errors.HTTP.Unauthorized(w)
			return
		}

		if token != secret {
			errors.HTTP.Unauthorized(w)
			return
		}

		h.ServeHTTP(w, r)
	}
}
