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

package middleware_test

import (
	"testing"
	"net/http"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"bytes"
	"github.com/spf13/viper"
	"fmt"
	"io/ioutil"
	"github.com/lastbackend/lastbackend/pkg/util/http/middleware"
)

// GetTestHandler returns a http.HandlerFunc for testing http middleware
func GetTestHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

// Testing NamespaceInfoH handler of a successful request (status 200)
func TestAuthenticateMiddleware(t *testing.T) {

	const token = "demotoken"

	viper.Set("security.token", token)

	tests := []struct {
		description  string
		url          string
		token        string
		expectedBody string
		expectedCode int
	}{
		{
			description:  http.StatusText(http.StatusUnauthorized),
			url:          "/",
			token:        "",
			expectedBody: "{\"code\":401,\"status\":\"Unauthorized\",\"message\":\"Unauthorized\"}",
			expectedCode: http.StatusUnauthorized,
		},
		{
			description:  http.StatusText(http.StatusOK),
			url:          "/",
			token:        fmt.Sprintf("Bearer %s", token),
			expectedBody: "",
			expectedCode: http.StatusOK,
		},
	}

	handler := middleware.Authenticate(GetTestHandler())
	ts := httptest.NewServer(handler)
	defer ts.Close()

	for _, tc := range tests {

		var u bytes.Buffer
		u.WriteString(string(ts.URL))
		u.WriteString(tc.url)

		req := httptest.NewRequest("GET", u.String(), nil)

		if len(tc.token) != 0 {
			req.Header.Add("Authorization", tc.token)
		}

		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		b, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		assert.Equal(t, tc.expectedCode, res.Code, tc.description)
		assert.Equal(t, tc.expectedBody, string(b), tc.description)
	}
}
