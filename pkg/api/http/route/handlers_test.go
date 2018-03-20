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

package route_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/http/route"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
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

// Testing RouteInfoH handler
func TestRouteInfo(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	err := envs.Get().GetStorage().Namespace().Clear(context.Background())
	assert.NoError(t, err)

	err = envs.Get().GetStorage().Route().Clear(context.Background())
	assert.NoError(t, err)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	r1 := getRouteAsset(ns1.Meta.Name, "demo")
	r2 := getRouteAsset(ns2.Meta.Name, "test")

	err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
	assert.NoError(t, err)

	err = envs.Get().GetStorage().Route().Insert(context.Background(), r1)
	assert.NoError(t, err)

	v, err := v1.View().Route().New(r1).ToJson()
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
			name:         "checking get route if not exists",
			url:          fmt.Sprintf("/namespace/%s/route/%s", ns1.Meta.Name, r2.Meta.Name),
			handler:      route.RouteInfoH,
			description:  "route not found",
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Route not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get route if namespace not exists",
			url:          fmt.Sprintf("/namespace/%s/route/%s", ns2.Meta.Name, r1.Meta.Name),
			handler:      route.RouteInfoH,
			description:  "namespace not found",
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get route successfully",
			url:          fmt.Sprintf("/namespace/%s/route/%s", ns1.Meta.Name, r1.Meta.Name),
			handler:      route.RouteInfoH,
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
			r.HandleFunc("/namespace/{namespace}/route/{route}", tc.handler)

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

// Testing RouteListH handler
func TestRouteList(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	err := envs.Get().GetStorage().Namespace().Clear(context.Background())
	assert.NoError(t, err)

	err = envs.Get().GetStorage().Route().Clear(context.Background())
	assert.NoError(t, err)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
	assert.NoError(t, err)

	r1 := getRouteAsset(ns1.Meta.Name, "demo")
	r2 := getRouteAsset(ns1.Meta.Name, "test")

	err = envs.Get().GetStorage().Route().Insert(context.Background(), r1)
	assert.NoError(t, err)

	err = envs.Get().GetStorage().Route().Insert(context.Background(), r2)
	assert.NoError(t, err)

	rl := make(types.RouteList, 0)
	rl[r1.SelfLink()] = r1
	rl[r2.SelfLink()] = r2

	v, err := v1.View().Route().NewList(rl).ToJson()
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
			name:         "checking get routes list if namespace not found",
			url:          fmt.Sprintf("/namespace/%s", ns2.Meta.Name),
			handler:      route.RouteListH,
			description:  "namespace not found",
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get routes list successfully",
			url:          fmt.Sprintf("/namespace/%s", ns1.Meta.Name),
			handler:      route.RouteListH,
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

type RouteCreateOptions struct {
	request.RouteCreateOptions
}

func createRouteCreateOptions(subdomain, domain string, custom, security bool, rules []request.RulesOption) *RouteCreateOptions {
	opts := new(RouteCreateOptions)
	opts.Subdomain = subdomain
	opts.Security = security
	opts.Domain = domain
	opts.Custom = custom
	opts.Rules = rules
	return opts
}

func (s *RouteCreateOptions) toJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

// Testing RouteCreateH handler
func TestRouteCreate(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	srtPointer := func(s string) *string { return &s }
	intPointer := func(i int) *int { return &i }

	err := envs.Get().GetStorage().Namespace().Clear(context.Background())
	assert.NoError(t, err)

	err = envs.Get().GetStorage().Route().Clear(context.Background())
	assert.NoError(t, err)

	ns := getNamespaceAsset("demo", "")

	err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns)
	assert.NoError(t, err)

	r1 := getRouteAsset(ns.Meta.Name, "demo")
	r2 := getRouteAsset(ns.Meta.Name, "test")

	err = envs.Get().GetStorage().Route().Insert(context.Background(), r1)
	assert.NoError(t, err)

	v, err := v1.View().Route().New(r1).ToJson()
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
		// TODO: need checking for unique
		{
			name:         "checking create route if namespace not found",
			description:  "namespace not found",
			url:          fmt.Sprintf("/namespace/%s/route", r2.Meta.Name),
			handler:      route.RouteCreateH,
			data:         createRouteCreateOptions("demo", "", false, false, []request.RulesOption{{Endpoint: srtPointer("route.test-domain.com"), Path: "/", Port: intPointer(80)}}).toJson(),
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "check create route if failed incoming json data",
			description:  "incoming json data is failed",
			url:          fmt.Sprintf("/namespace/%s/route", r1.Meta.Name),
			handler:      route.RouteCreateH,
			data:         "{name:demo}",
			expectedBody: "{\"code\":400,\"status\":\"Incorrect json\",\"message\":\"Incorrect json\"}",
			expectedCode: http.StatusBadRequest,
		},
		// TODO: need checking incoming data for validity
		{
			name:         "check create route success",
			description:  "successfully",
			url:          fmt.Sprintf("/namespace/%s/route", r1.Meta.Name),
			handler:      route.RouteCreateH,
			data:         createRouteCreateOptions("demo", "", false, false, []request.RulesOption{{Endpoint: srtPointer("route.test-domain.com"), Path: "/", Port: intPointer(80)}}).toJson(),
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
			r.HandleFunc("/namespace/{namespace}/route", tc.handler)

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

type RouteUpdateOptions struct {
	request.RouteUpdateOptions
}

func createRouteUpdateOptions(subdomain, domain string, custom, security bool, rules []request.RulesOption) *RouteUpdateOptions {
	opts := new(RouteUpdateOptions)
	opts.Subdomain = subdomain
	opts.Security = security
	opts.Domain = domain
	opts.Custom = custom
	opts.Rules = rules
	return opts
}

func (s *RouteUpdateOptions) toJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

// Testing RouteUpdateH handler
func TestRouteUpdate(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	srtPointer := func(s string) *string { return &s }
	intPointer := func(i int) *int { return &i }

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
	assert.NoError(t, err)

	s1 := getRouteAsset(ns1.Meta.Name, "demo")
	s2 := getRouteAsset(ns1.Meta.Name, "test")

	err = envs.Get().GetStorage().Route().Insert(context.Background(), s1)
	assert.NoError(t, err)

	v, err := v1.View().Route().New(s1).ToJson()
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
			name:         "checking update route if name not exists",
			description:  "route not exists",
			url:          fmt.Sprintf("/namespace/%s/route/%s", ns1.Meta.Name, s2.Meta.Name),
			handler:      route.RouteUpdateH,
			data:         createRouteCreateOptions("demo", "", false, false, []request.RulesOption{{Endpoint: srtPointer("route.test-domain.com"), Path: "/", Port: intPointer(80)}}).toJson(),
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Route not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update route if namespace not found",
			description:  "namespace not found",
			url:          fmt.Sprintf("/namespace/%s/route/%s", ns2.Meta.Name, s1.Meta.Name),
			handler:      route.RouteUpdateH,
			data:         createRouteCreateOptions("demo", "", false, false, []request.RulesOption{{Endpoint: srtPointer("route.test-domain.com"), Path: "/", Port: intPointer(80)}}).toJson(),
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "check update route if failed incoming json data",
			description:  "incoming json data is failed",
			url:          fmt.Sprintf("/namespace/%s/route/%s", ns1.Meta.Name, s1.Meta.Name),
			handler:      route.RouteUpdateH,
			data:         "{name:demo}",
			expectedBody: "{\"code\":400,\"status\":\"Incorrect json\",\"message\":\"Incorrect json\"}",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check update route success",
			description:  "successfully",
			url:          fmt.Sprintf("/namespace/%s/route/%s", ns1.Meta.Name, s1.Meta.Name),
			handler:      route.RouteUpdateH,
			data:         createRouteCreateOptions("demo", "", false, false, []request.RulesOption{{Endpoint: srtPointer("route.test-domain.com"), Path: "/", Port: intPointer(80)}}).toJson(),
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
			r.HandleFunc("/namespace/{namespace}/route/{route}", tc.handler)

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

// Testing RouteRemoveH handler
func TestRouteRemove(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
	assert.NoError(t, err)

	r1 := getRouteAsset(ns1.Meta.Name, "demo")
	r2 := getRouteAsset(ns1.Meta.Name, "test")

	err = envs.Get().GetStorage().Route().Insert(context.Background(), r1)
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
			name:         "checking get route if not exists",
			url:          fmt.Sprintf("/namespace/%s/route/%s", ns1.Meta.Name, r2.Meta.Name),
			handler:      route.RouteRemoveH,
			description:  "route not found",
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Route not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get route if namespace not exists",
			url:          fmt.Sprintf("/namespace/%s/route/%s", ns2.Meta.Name, r1.Meta.Name),
			handler:      route.RouteRemoveH,
			description:  "namespace not found",
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get route successfully",
			url:          fmt.Sprintf("/namespace/%s/route/%s", ns1.Meta.Name, r1.Meta.Name),
			handler:      route.RouteRemoveH,
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
			r.HandleFunc("/namespace/{namespace}/route/{route}", tc.handler)

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

func getNamespaceAsset(name, desc string) *types.Namespace {
	var n = types.Namespace{}
	n.Meta.SetDefault()
	n.Meta.Name = name
	n.Meta.Description = desc
	return &n
}

func getRouteAsset(namespace, name string) *types.Route {
	var r = types.Route{}
	r.Meta.SetDefault()
	r.Meta.Namespace = namespace
	r.Meta.Name = name
	r.Meta.Security = true
	r.Spec.Domain = fmt.Sprintf("%s.test-domain.com", name)
	r.Spec.Rules = append(r.Spec.Rules, &types.RouteRule{
		Path:     "/",
		Endpoint: "route.test-domain.com",
		Port:     80,
	})
	return &r
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}
