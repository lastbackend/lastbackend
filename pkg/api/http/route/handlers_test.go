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
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"errors"
)

// Testing RouteInfoH handler
func TestRouteInfo(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
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

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		route     *types.Route
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		err          string
		want         *types.Route
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get route if not exists",
			args:         args{ctx, ns1, r2},
			fields:       fields{stg},
			handler:      route.RouteInfoH,
			description:  "route not found",
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Route not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get route if namespace not exists",
			args:         args{ctx, ns2, r1},
			fields:       fields{stg},
			handler:      route.RouteInfoH,
			description:  "namespace not found",
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get route successfully",
			args:         args{ctx, ns1, r1},
			fields:       fields{stg},
			handler:      route.RouteInfoH,
			description:  "successfully",
			want:         r1,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", fmt.Sprintf("/namespace/%s/route/%s", tc.args.namespace.Meta.Name, tc.args.route.Meta.Name), nil)
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

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), tc.description)
			} else {

				n := new(views.Route)
				err := json.Unmarshal(body, &n)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, n.Meta.Name, tc.description)
			}
		})
	}

}

// Testing RouteListH handler
func TestRouteList(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	r1 := getRouteAsset(ns1.Meta.Name, "demo")
	r2 := getRouteAsset(ns1.Meta.Name, "test")

	rl := make(types.RouteList, 0)
	rl[r1.SelfLink()] = r1
	rl[r2.SelfLink()] = r2

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		err          string
		want         types.RouteList
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get routes list if namespace not found",
			args:         args{ctx, ns2},
			fields:       fields{stg},
			handler:      route.RouteListH,
			description:  "namespace not found",
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get routes list successfully",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      route.RouteListH,
			description:  "successfully",
			want:         rl,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			err := envs.Get().GetStorage().Namespace().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Route().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Route().Insert(context.Background(), r1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Route().Insert(context.Background(), r2)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", fmt.Sprintf("/namespace/%s", tc.args.namespace.Meta.Name), nil)
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

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), tc.description)
			} else {

				r := new(views.RouteList)
				err := json.Unmarshal(body, &r)
				assert.NoError(t, err)

				for _, item := range *r {
					if _, ok := tc.want[item.Meta.SelfLink]; !ok {
						assert.Error(t, errors.New("not equals"))
					}
				}
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

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	srtPointer := func(s string) *string { return &s }
	intPointer := func(i int) *int { return &i }

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	r1 := getRouteAsset(ns1.Meta.Name, "demo")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		data         string
		err          string
		want         *types.Route
		wantErr      bool
		expectedCode int
	}{
		// TODO: need checking for unique
		{
			name:         "checking create route if namespace not found",
			description:  "namespace not found",
			args:         args{ctx, ns2},
			fields:       fields{stg},
			handler:      route.RouteCreateH,
			data:         createRouteCreateOptions("demo", "", false, false, []request.RulesOption{{Endpoint: srtPointer("route.test-domain.com"), Path: "/", Port: intPointer(80)}}).toJson(),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "check create route if failed incoming json data",
			description:  "incoming json data is failed",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      route.RouteCreateH,
			data:         "{name:demo}",
			err:          "{\"code\":400,\"status\":\"Incorrect json\",\"message\":\"Incorrect json\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		// TODO: need checking incoming data for validity
		{
			name:         "check create route success",
			description:  "successfully",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      route.RouteCreateH,
			data:         createRouteCreateOptions("demo", "", false, false, []request.RulesOption{{Endpoint: srtPointer("route.test-domain.com"), Path: "/", Port: intPointer(80)}}).toJson(),
			want:         r1,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			err := envs.Get().GetStorage().Namespace().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Route().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Route().Insert(context.Background(), r1)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("POST", fmt.Sprintf("/namespace/%s/route", tc.args.namespace.Meta.Name), strings.NewReader(tc.data))
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

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), tc.description)
			} else {

				got, err := tc.fields.stg.Route().Get(tc.args.ctx, tc.args.namespace.Meta.Name, tc.want.Meta.Name)
				assert.NoError(t, err)

				assert.Equal(t, ns1.Meta.Name, got.Meta.Name, tc.description)
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

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	srtPointer := func(s string) *string { return &s }
	intPointer := func(i int) *int { return &i }

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
	assert.NoError(t, err)

	r1 := getRouteAsset(ns1.Meta.Name, "demo")
	r2 := getRouteAsset(ns1.Meta.Name, "test")
	r3 := getRouteAsset(ns1.Meta.Name, "demo")

	err = envs.Get().GetStorage().Route().Insert(context.Background(), r1)
	assert.NoError(t, err)

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		route     *types.Route
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		data         string
		err          string
		want         *types.Route
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking update route if name not exists",
			description:  "route not exists",
			args:         args{ctx, ns1, r2},
			fields:       fields{stg},
			handler:      route.RouteUpdateH,
			data:         createRouteUpdateOptions("demo", "", false, false, []request.RulesOption{{Endpoint: srtPointer("route.test-domain.com"), Path: "/", Port: intPointer(80)}}).toJson(),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Route not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update route if namespace not found",
			description:  "namespace not found",
			args:         args{ctx, ns2, r1},
			fields:       fields{stg},
			handler:      route.RouteUpdateH,
			data:         createRouteUpdateOptions("demo", "", false, false, []request.RulesOption{{Endpoint: srtPointer("route.test-domain.com"), Path: "/", Port: intPointer(80)}}).toJson(),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "check update route if failed incoming json data",
			description:  "incoming json data is failed",
			args:         args{ctx, ns1, r1},
			fields:       fields{stg},
			handler:      route.RouteUpdateH,
			data:         "{name:demo}",
			err:          "{\"code\":400,\"status\":\"Incorrect json\",\"message\":\"Incorrect json\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check update route success",
			description:  "successfully",
			args:         args{ctx, ns1, r1},
			fields:       fields{stg},
			handler:      route.RouteUpdateH,
			data:         createRouteUpdateOptions("demo", "", false, false, []request.RulesOption{{Endpoint: srtPointer("route.test-domain.com"), Path: "/", Port: intPointer(80)}}).toJson(),
			want:         r3,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("PUT", fmt.Sprintf("/namespace/%s/route/%s", tc.args.namespace.Meta.Name, tc.args.route.Meta.Name), strings.NewReader(tc.data))
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

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), tc.description)
			} else {

				n := new(views.Route)
				err := json.Unmarshal(body, &n)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, n.Meta.Name, tc.description)
			}
		})
	}

}

// Testing RouteRemoveH handler
func TestRouteRemove(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	r1 := getRouteAsset(ns1.Meta.Name, "demo")
	r2 := getRouteAsset(ns1.Meta.Name, "test")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		route     *types.Route
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		err          string
		want         string
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get route if not exists",
			args:         args{ctx, ns1, r2},
			fields:       fields{stg},
			handler:      route.RouteRemoveH,
			description:  "route not found",
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Route not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get route if namespace not exists",
			args:         args{ctx, ns2, r1},
			fields:       fields{stg},
			handler:      route.RouteRemoveH,
			description:  "namespace not found",
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get route successfully",
			args:         args{ctx, ns1, r1},
			fields:       fields{stg},
			handler:      route.RouteRemoveH,
			description:  "successfully",
			want:         "",
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Route().Insert(context.Background(), r1)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/namespace/%s/route/%s", tc.args.namespace.Meta.Name, tc.args.route.Meta.Name), nil)
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

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), tc.description)
			} else {
				assert.Equal(t, tc.want, string(body), tc.description)
			}
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
