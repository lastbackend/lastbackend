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

package namespace

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/views/v1"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/util/http/middleware"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	token                  = "demotoken"
	namespaceExistsName    = "demo"
	namespaceNotExistsName = "notexistsname"
)

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}

// Testing NamespaceInfoH handler
func TestNamespaceGet(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	strg.Namespace().Insert()


	viper.Set("verbose", 0)

	tests := []struct {
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		expectedBody string
		expectedCode int
	}{
		{
			url:          fmt.Sprintf("/namespace/%s", namespaceExistsName),
			handler:      NamespaceInfoH,
			description:  "successfully",
			expectedCode: http.StatusOK,
		},
		{
			url: fmt.Sprintf("/namespace/%s", namespaceExistsName),
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			handler:      middleware.Authenticate(NamespaceInfoH),
			description:  "successfully",
			expectedCode: http.StatusOK,
		},
		{
			url:         fmt.Sprintf("/namespace/%s", namespaceNotExistsName),
			handler:     NamespaceInfoH,
			description: "namespace not found",
			//expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Not found\"}",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {

		// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("GET", tc.url, nil)
		assert.NoError(t, err)

		if tc.headers != nil {
			for key, val := range tc.headers {
				req.Header.Set(key, val)
			}
		}

		r := mux.NewRouter()
		r.HandleFunc("/namespace/{namespace}", tc.handler)

		setRequestVars(r, req)

		// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		res := httptest.NewRecorder()

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		r.ServeHTTP(res, req)

		// Check the status code is what we expect.
		assert.Equal(t, tc.expectedCode, res.Code, tc.description)

		if res.Code == http.StatusOK {

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			ns := new(v1.Namespace)
			err = json.Unmarshal(body, ns)
			assert.NoError(t, err)

			// TODO: check response data with expectedBody
		}
	}

}

// Testing NamespaceInfoH handler
func TestNamespaceList(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	tests := []struct {
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		expectedBody string
		expectedCode int
	}{
		{
			url:          "/namespace",
			handler:      NamespaceListH,
			description:  "successfully",
			expectedCode: http.StatusOK,
		},
		{
			url: fmt.Sprintf("/namespace/%s", namespaceNotExistsName),
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			handler:      NamespaceListH,
			description:  "successfully",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {

		// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("GET", tc.url, nil)
		assert.NoError(t, err)

		if tc.headers != nil {
			for key, val := range tc.headers {
				req.Header.Set(key, val)
			}
		}

		r := mux.NewRouter()
		r.HandleFunc("/namespace", tc.handler)

		setRequestVars(r, req)

		// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		res := httptest.NewRecorder()

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		r.ServeHTTP(res, req)

		// Check the status code is what we expect.
		assert.Equal(t, tc.expectedCode, res.Code, tc.description)

		if res.Code == http.StatusOK {

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			ns := new(v1.NamespaceList)
			err = json.Unmarshal(body, ns)
			assert.NoError(t, err)

			// TODO: check response data with expectedBody
		}
	}

}
