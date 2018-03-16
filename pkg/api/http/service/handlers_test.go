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

package service_test

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/http/service"
	"github.com/lastbackend/lastbackend/pkg/api/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"encoding/json"
)

// Testing ServiceInfoH handler
func TestServiceInfo(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	ns1 := getDefaultNamespace("demo")
	ns2 := getDefaultNamespace("test")
	s1 := getDefaultService("demo")
	s2 := getDefaultService("test")

	err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
	assert.NoError(t, err)


	err = envs.Get().GetStorage().Service().Insert(context.Background(), s1)
	assert.NoError(t, err)

	v, err := views.V1().Service().New(s1, make([]*types.Deployment, 0), make([]*types.Pod, 0)).ToJson()
	assert.NoError(t, err)

	tests := []struct {
		name         string
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "checking get service if not exists",
			url:          fmt.Sprintf("/namespace/%s/service/%s", ns1.Meta.Name, s2.Meta.Name),
			handler:      service.ServiceInfoH,
			description:  "service not found",
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Service not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service if namespace not exists",
			url:          fmt.Sprintf("/namespace/%s/service/%s", ns2.Meta.Name, s2.Meta.Name),
			handler:      service.ServiceInfoH,
			description:  "namespace not found",
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service successfully",
			url:          fmt.Sprintf("/namespace/%s/service/%s", ns1.Meta.Name, s1.Meta.Name),
			handler:      service.ServiceInfoH,
			description:  "successfully",
			expectedBody: string(v),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

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
			r.HandleFunc("/namespace/{namespace}/service/{service}", tc.handler)

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

// Testing ServiceListH handler
func TestServiceList(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	ns1 := getDefaultNamespace("demo")
	ns2 := getDefaultNamespace("test")
	err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
	assert.NoError(t, err)

	s1 := getDefaultService("demo")
	s2 := getDefaultService("test")

	err = envs.Get().GetStorage().Service().Insert(context.Background(), s1)
	assert.NoError(t, err)

	err = envs.Get().GetStorage().Service().Insert(context.Background(), s2)
	assert.NoError(t, err)

	v, err := views.V1().Service().New(s2, make([]*types.Deployment, 0), make([]*types.Pod, 0)).ToJson()
	assert.NoError(t, err)

	tests := []struct {
		name         string
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "checking get services list if namespace not found",
			url:          fmt.Sprintf("/namespace/%s", ns2.Meta.Name),
			handler:      service.ServiceListH,
			description:  "namespace not found",
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get services list successfully",
			url:          fmt.Sprintf("/namespace/%s", ns1.Meta.Name),
			handler:      service.ServiceListH,
			description:  "successfully",
			expectedBody: string(v),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

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
		})
	}

}

type ServiceCreateOptions struct {
	types.ServiceCreateOptions
}

func createServiceCreateOptions(name, description, sources *string, replicas *int, spec *types.ServiceOptionsSpec) *ServiceCreateOptions {
	opts := new(ServiceCreateOptions)
	opts.Name = name
	opts.Description = description
	opts.Replicas = replicas
	opts.Sources = sources
	opts.Spec = spec
	return opts
}

func (s *ServiceCreateOptions) toJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

// Testing ServiceCreateH handler
func TestServiceCreate(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	srtPointer := func(s string) *string { return &s }
	intPointer := func(i int) *int { return &i }
	int64Pointer := func(i int64) *int64 { return &i }

	ns := getDefaultNamespace("demo")

	err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns)
	assert.NoError(t, err)

	s1 := getDefaultService("demo")
	s2 := getDefaultService("test")

	err = envs.Get().GetStorage().Service().Insert(context.Background(), s1)
	assert.NoError(t, err)

	v, err := views.V1().Service().New(s1, make([]*types.Deployment, 0), make([]*types.Pod, 0)).ToJson()
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
			name:         "checking create service if name already exists",
			description:  "service already exists",
			url:          fmt.Sprintf("/namespace/%s/service", s1.Meta.Name),
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("demo"), nil, srtPointer("redis"), intPointer(1), nil).toJson(),
			expectedBody: "{\"code\":400,\"status\":\"Not Unique\",\"message\":\"Name is already in use\"}",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking create service if namespace not found",
			description:  "namespace not found",
			url:          fmt.Sprintf("/namespace/%s/service", s2.Meta.Name),
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("test"), nil, srtPointer("redis"), intPointer(1), nil).toJson(),
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "check create service if failed incoming json data",
			description:  "incoming json data is failed",
			url:          fmt.Sprintf("/namespace/%s/service", s1.Meta.Name),
			handler:      service.ServiceCreateH,
			data:         "{name:demo}",
			expectedBody: "{\"code\":400,\"status\":\"Incorrect json\",\"message\":\"Incorrect json\"}",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check create service if bad parameter name",
			description:  "incorrect name parameter",
			url:          fmt.Sprintf("/namespace/%s/service", s1.Meta.Name),
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("___test"), nil, srtPointer("redis"), intPointer(1), nil).toJson(),
			expectedBody: "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad name parameter\"}",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check create service if bad parameter memory",
			description:  "incorrect memory parameter",
			url:          fmt.Sprintf("/namespace/%s/service", s1.Meta.Name),
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("test"), nil, srtPointer("redis"), intPointer(1), &types.ServiceOptionsSpec{Memory: int64Pointer(127)}).toJson(),
			expectedBody: "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad memory parameter\"}",
			expectedCode: http.StatusBadRequest,
		},
		// TODO: check another spec parameters
		{
			name:         "check create service if bad parameter replicas",
			description:  "incorrect replicas parameter",
			url:          fmt.Sprintf("/namespace/%s/service", s1.Meta.Name),
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("test"), nil, srtPointer("redis"), intPointer(-1), nil).toJson(),
			expectedBody: "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad replicas parameter\"}",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check create service success",
			description:  "successfully",
			url:          fmt.Sprintf("/namespace/%s/service", s1.Meta.Name),
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("test"), nil, srtPointer("redis"), intPointer(1), nil).toJson(),
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
			r.HandleFunc("/namespace/{namespace}/service", tc.handler)

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

type ServiceUpdateOptions struct {
	types.ServiceUpdateOptions
}

func createServiceUpdateOptions(description, sources *string, replicas *int, spec *types.ServiceOptionsSpec) *ServiceUpdateOptions {
	opts := new(ServiceUpdateOptions)
	opts.Description = description
	opts.Replicas = replicas
	opts.Sources = sources
	opts.Spec = spec
	return opts
}

func (s *ServiceUpdateOptions) toJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

// Testing ServiceUpdateH handler
func TestServiceUpdate(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	srtPointer := func(s string) *string { return &s }
	intPointer := func(i int) *int { return &i }
	int64Pointer := func(i int64) *int64 { return &i }

	ns1 := getDefaultNamespace("demo")
	ns2 := getDefaultNamespace("test")

	err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
	assert.NoError(t, err)

	s1 := getDefaultService("demo")
	s2 := getDefaultService("test")

	err = envs.Get().GetStorage().Service().Insert(context.Background(), s1)
	assert.NoError(t, err)

	v, err := views.V1().Service().New(s1, make([]*types.Deployment, 0), make([]*types.Pod, 0)).ToJson()
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
			name:         "checking update service if name not exists",
			description:  "service not exists",
			url:          fmt.Sprintf("/namespace/%s/service/%s", ns1.Meta.Name, s2.Meta.Name),
			handler:      service.ServiceUpdateH,
			data:         createServiceUpdateOptions(nil, srtPointer("redis"), intPointer(1), nil).toJson(),
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Service not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update service if namespace not found",
			description:  "namespace not found",
			url:          fmt.Sprintf("/namespace/%s/service/%s", ns2.Meta.Name, s1.Meta.Name),
			handler:      service.ServiceUpdateH,
			data:         createServiceUpdateOptions(nil, srtPointer("redis"), intPointer(1), nil).toJson(),
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "check update service if failed incoming json data",
			description:  "incoming json data is failed",
			url:          fmt.Sprintf("/namespace/%s/service/%s", ns1.Meta.Name, s1.Meta.Name),
			handler:      service.ServiceUpdateH,
			data:         "{name:demo}",
			expectedBody: "{\"code\":400,\"status\":\"Incorrect json\",\"message\":\"Incorrect json\"}",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check update service if bad parameter memory",
			description:  "incorrect memory parameter",
			url:          fmt.Sprintf("/namespace/%s/service/%s", ns1.Meta.Name, s1.Meta.Name),
			handler:      service.ServiceUpdateH,
			data:         createServiceUpdateOptions(nil, srtPointer("redis"), intPointer(1), &types.ServiceOptionsSpec{Memory: int64Pointer(127)}).toJson(),
			expectedBody: "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad memory parameter\"}",
			expectedCode: http.StatusBadRequest,
		},
		// TODO: check another spec parameters
		{
			name:         "check update service if bad parameter replicas",
			description:  "incorrect replicas parameter",
			url:          fmt.Sprintf("/namespace/%s/service/%s", ns1.Meta.Name, s1.Meta.Name),
			handler:      service.ServiceUpdateH,
			data:         createServiceUpdateOptions(nil, srtPointer("redis"), intPointer(-1), nil).toJson(),
			expectedBody: "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad replicas parameter\"}",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check update service success",
			description:  "successfully",
			url:          fmt.Sprintf("/namespace/%s/service/%s", ns1.Meta.Name, s1.Meta.Name),
			handler:      service.ServiceUpdateH,
			data:         createServiceUpdateOptions(nil, srtPointer("redis"), intPointer(1), nil).toJson(),
			expectedBody: string(v),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

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
			r.HandleFunc("/namespace/{namespace}/service/{service}", tc.handler)

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

// Testing ServiceRemoveH handler
func TestServiceRemove(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	ns1 := getDefaultNamespace("demo")
	ns2 := getDefaultNamespace("test")
	err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
	assert.NoError(t, err)

	s1 := getDefaultService("demo")
	s2 := getDefaultService("test")
	err = envs.Get().GetStorage().Service().Insert(context.Background(), s1)
	assert.NoError(t, err)

	tests := []struct {
		name         string
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "checking get service if not exists",
			url:          fmt.Sprintf("/namespace/%s/service/%s", ns1.Meta.Name, s2.Meta.Name),
			handler:      service.ServiceRemoveH,
			description:  "service not found",
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Service not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service if namespace not exists",
			url:          fmt.Sprintf("/namespace/%s/service/%s", ns2.Meta.Name, s1.Meta.Name),
			handler:      service.ServiceRemoveH,
			description:  "namespace not found",
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service successfully",
			url:          fmt.Sprintf("/namespace/%s/service/%s", ns1.Meta.Name, s1.Meta.Name),
			handler:      service.ServiceRemoveH,
			description:  "successfully",
			expectedBody: "",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

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
			r.HandleFunc("/namespace/{namespace}/service/{service}", tc.handler)

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
		})
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

func getDefaultService(name string) *types.Service {
	res := new(types.Service)
	switch name {
	case "demo":
		res.Meta.SetDefault()
		res.Meta.Name = "demo"
		res.Meta.Description = "demo description"
		res.Meta.Namespace = "demo"
	case "test":
		res.Meta.SetDefault()
		res.Meta.Name = "test"
		res.Meta.Description = "test description"
		res.Meta.Namespace = "demo"
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
