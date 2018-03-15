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
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/views/v1"
	"github.com/lastbackend/lastbackend/pkg/util/http/middleware"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"strings"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
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
	viper.Set("verbose", 0)

	tests := []struct {
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		expectedBody string
		expectedCode int
		want         *v1.Namespace
		wantErr      *errors.Http
		isErr        bool
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
			url:          fmt.Sprintf("/namespace/%s", namespaceNotExistsName),
			handler:      NamespaceInfoH,
			description:  "namespace not found",
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

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		if res.Code == http.StatusOK {

			ns := new(v1.Namespace)
			err = json.Unmarshal(body, ns)
			assert.NoError(t, err)

			fmt.Println(">>>>>>>>", ns.Meta.Name)

			// TODO: check response data with expectedBody
		} else {
			e := new(errors.Http)
			err = json.Unmarshal(body, e)
			assert.NoError(t, err)

			fmt.Println(">>>>>>>>", e.Code, e.Status, e.Message)
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

type NamespaceCreateOptions struct {
	types.NamespaceCreateOptions
}

func createNamespaceCreateOptions(name, description string, quotas *types.NamespaceQuotasOptions) *NamespaceCreateOptions {
	opts := new(NamespaceCreateOptions)
	opts.Name = name
	opts.Description = description
	opts.Quotas = quotas
	return opts
}

func (s *NamespaceCreateOptions) toJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

// Testing NamespaceCreateH handler
func TestNamespaceCreate(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	tests := []struct {
		name         string
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		data         string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "check create namespace success",
			description:  "successfully",
			url:          "/namespace",
			handler:      NamespaceCreateH,
			data:         createNamespaceCreateOptions("test", "", &types.NamespaceQuotasOptions{RAM: 2, Routes: 1}).toJson(),
			expectedCode: http.StatusOK,
		},
		{
			name:        "check create namespace success with auth middleware",
			description: "successfully",
			url:         "/namespace",
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			handler:      middleware.Authenticate(NamespaceCreateH),
			data:         createNamespaceCreateOptions("test", "", &types.NamespaceQuotasOptions{RAM: 2, Routes: 1}).toJson(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "check create namespace success if name already exists",
			description:  "namespace already exists",
			url:          "/namespace",
			handler:      NamespaceCreateH,
			data:         createNamespaceCreateOptions("demo", "", nil).toJson(),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("POST", tc.url, strings.NewReader(tc.data))
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

			// TODO: check response data with expectedBody

		})
	}

}

type NamespaceUpdateOptions struct {
	types.NamespaceUpdateOptions
}

func createNamespaceUpdateOptions(description *string, quotas *types.NamespaceQuotasOptions) *NamespaceUpdateOptions {
	opts := new(NamespaceUpdateOptions)
	opts.Description = description
	opts.Quotas = quotas
	return opts
}

func (s *NamespaceUpdateOptions) toJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

// Testing NamespaceUpdateH handler
func TestNamespaceUpdate(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	tests := []struct {
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		data         string
		expectedBody string
		expectedCode int
	}{
		{
			url:          fmt.Sprintf("/namespace/%s", namespaceExistsName),
			handler:      NamespaceUpdateH,
			description:  "successfully",
			data:         createNamespaceUpdateOptions(nil, &types.NamespaceQuotasOptions{RAM: 2, Routes: 1}).toJson(),
			expectedCode: http.StatusOK,
		},
		{
			url: fmt.Sprintf("/namespace/%s", namespaceExistsName),
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			handler:      middleware.Authenticate(NamespaceUpdateH),
			description:  "successfully",
			data:         createNamespaceUpdateOptions(nil, &types.NamespaceQuotasOptions{RAM: 2, Routes: 1}).toJson(),
			expectedCode: http.StatusOK,
		},
		{
			url:          fmt.Sprintf("/namespace/%s", namespaceNotExistsName),
			handler:      NamespaceUpdateH,
			description:  "namespace not exists",
			data:         createNamespaceUpdateOptions(nil, nil).toJson(),
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {

		// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("PUT", tc.url, strings.NewReader(tc.data))
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

		// TODO: check response data with expectedBody

	}

}

// Testing NamespaceRemoveH handler
func TestNamespaceRemove(t *testing.T) {

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
			url:          fmt.Sprintf("/namespace/%s", namespaceExistsName),
			handler:      NamespaceRemoveH,
			description:  "successfully",
			expectedCode: http.StatusOK,
		},
		{
			url: fmt.Sprintf("/namespace/%s", namespaceExistsName),
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			handler:      middleware.Authenticate(NamespaceRemoveH),
			description:  "successfully",
			expectedCode: http.StatusOK,
		},
		{
			url:          fmt.Sprintf("/namespace/%s", namespaceNotExistsName),
			handler:      NamespaceRemoveH,
			description:  "namespace not found",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {

		// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("DELETE", tc.url, nil)
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

		// TODO: check response data with expectedBody

	}

}
