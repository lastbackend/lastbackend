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

package namespace_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/http/namespace"
	"github.com/lastbackend/lastbackend/pkg/api/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Testing NamespaceInfoH handler
func TestNamespaceInfo(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	ns1 := getDefaultNamespace("demo")
	ns2 := getDefaultNamespace("test")

	err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
	assert.NoError(t, err)

	v, err := views.V1().Namespace().New(ns1).ToJson()
	assert.NoError(t, err)

	tests := []struct {
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		expectedBody string
		expectedCode int
	}{
		{
			url:          fmt.Sprintf("/namespace/%s", ns2.Meta.Name),
			handler:      namespace.NamespaceInfoH,
			description:  "namespace not found",
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			url:          fmt.Sprintf("/namespace/%s", ns1.Meta.Name),
			handler:      namespace.NamespaceInfoH,
			description:  "successfully",
			expectedBody: string(v),
			expectedCode: http.StatusOK,
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
			assert.Equal(t, tc.expectedBody, string(v), tc.description)
		} else {
			assert.Equal(t, tc.expectedBody, string(body), tc.description)
		}
	}

}

// Testing NamespaceInfoH handler
func TestNamespaceList(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	nsl := make(types.NamespaceList, 0)
	ns1 := getDefaultNamespace("demo")
	ns2 := getDefaultNamespace("test")

	err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
	assert.NoError(t, err)
	nsl = append(nsl, ns1)
	err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns2)
	assert.NoError(t, err)
	nsl = append(nsl, ns2)

	v, err := views.V1().Namespace().NewList(nsl).ToJson()
	assert.NoError(t, err)

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
			handler:      namespace.NamespaceListH,
			description:  "successfully",
			expectedBody: string(v),
			expectedCode: http.StatusOK,
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

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		if res.Code == http.StatusOK {
			assert.Equal(t, tc.expectedBody, string(v), tc.description)
		} else {
			assert.Equal(t, tc.expectedBody, string(body), tc.description)
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

	ns := getDefaultNamespace("demo")

	err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns)
	assert.NoError(t, err)

	v, err := views.V1().Namespace().New(ns).ToJson()
	assert.NoError(t, err)

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
			name:         "check create namespace success if name already exists",
			description:  "namespace already exists",
			url:          "/namespace",
			handler:      namespace.NamespaceCreateH,
			data:         createNamespaceCreateOptions("demo", "", nil).toJson(),
			expectedBody: "{\"code\":400,\"status\":\"Not Unique\",\"message\":\"Name is already in use\"}",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check create namespace error",
			description:  "successfully",
			url:          "/namespace",
			handler:      namespace.NamespaceCreateH,
			data:         createNamespaceCreateOptions("__test", "", &types.NamespaceQuotasOptions{RAM: 2, Routes: 1}).toJson(),
			expectedBody: "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad name parameter\"}",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check create namespace error",
			description:  "successfully",
			url:          "/namespace",
			handler:      namespace.NamespaceCreateH,
			data:         "{name:demo}",
			expectedBody: "{\"code\":400,\"status\":\"Incorrect json\",\"message\":\"Incorrect json\"}",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check create namespace success",
			description:  "successfully",
			url:          "/namespace",
			handler:      namespace.NamespaceCreateH,
			data:         createNamespaceCreateOptions("test", "", &types.NamespaceQuotasOptions{RAM: 2, Routes: 1}).toJson(),
			expectedBody: string(v),
			expectedCode: http.StatusOK,
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

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if res.Code == http.StatusOK {
				assert.Equal(t, tc.expectedBody, string(v), tc.description)
			} else {
				assert.Equal(t, tc.expectedBody, string(body), tc.description)
			}

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

	ns1 := getDefaultNamespace("demo")
	ns2 := getDefaultNamespace("test")

	err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
	assert.NoError(t, err)

	v, err := views.V1().Namespace().New(ns2).ToJson()
	assert.NoError(t, err)

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
			url:          fmt.Sprintf("/namespace/%s", ns2.Meta.Name),
			handler:      namespace.NamespaceUpdateH,
			description:  "namespace not exists",
			data:         createNamespaceUpdateOptions(nil, nil).toJson(),
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			description:  "successfully",
			url:          fmt.Sprintf("/namespace/%s", ns1.Meta.Name),
			handler:      namespace.NamespaceUpdateH,
			data:         "{description:demo}",
			expectedBody: "{\"code\":400,\"status\":\"Incorrect json\",\"message\":\"Incorrect json\"}",
			expectedCode: http.StatusBadRequest,
		},
		{
			url:          fmt.Sprintf("/namespace/%s", ns1.Meta.Name),
			handler:      namespace.NamespaceUpdateH,
			description:  "successfully",
			data:         createNamespaceUpdateOptions(nil, &types.NamespaceQuotasOptions{RAM: ns2.Resources.RAM, Routes: ns2.Resources.Routes}).toJson(),
			expectedBody: string(v),
			expectedCode: http.StatusOK,
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

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		if res.Code == http.StatusOK {
			assert.Equal(t, tc.expectedBody, string(v), tc.description)
		} else {
			assert.Equal(t, tc.expectedBody, string(body), tc.description)
		}

	}

}

// Testing NamespaceRemoveH handler
func TestNamespaceRemove(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	ns := getDefaultNamespace("demo")

	err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns)
	assert.NoError(t, err)

	tests := []struct {
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		expectedBody string
		expectedCode int
	}{
		{
			url:          fmt.Sprintf("/namespace/%s", ns.Meta.Name),
			handler:      namespace.NamespaceRemoveH,
			description:  "successfully",
			expectedCode: http.StatusOK,
		},
		{
			url:          fmt.Sprintf("/namespace/%s", ns.Meta.Name),
			handler:      namespace.NamespaceRemoveH,
			description:  "namespace not found",
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
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

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		assert.Equal(t, tc.expectedBody, string(body), tc.description)

	}

}

func getDefaultNamespace(name string) *types.Namespace {
	res := new(types.Namespace)
	switch name {
	case "demo":
		res.Meta.SetDefault()
		res.Meta.Name = "demo"
		res.Meta.Description = "demo description"
		res.Quotas.Routes = int(2)
		res.Quotas.RAM = int64(256)
	case "test":
		res.Meta.SetDefault()
		res.Meta.Name = "test"
		res.Meta.Description = "test description"
		res.Quotas.Routes = int(1)
		res.Quotas.RAM = int64(128)
	default:
		res = nil
	}
	return res
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}
