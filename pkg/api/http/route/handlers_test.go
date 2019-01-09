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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/http/route"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/stretchr/testify/assert"
)

// Testing RouteInfoH handler
func TestRouteInfo(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	r1 := getRouteAsset(ns1.Meta.Name, "demo")
	r2 := getRouteAsset(ns2.Meta.Name, "test")

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
		err          string
		want         *views.Route
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get route if not exists",
			args:         args{ctx, ns1, r2},
			fields:       fields{stg},
			handler:      route.RouteInfoH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Route not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get route if namespace not exists",
			args:         args{ctx, ns2, r1},
			fields:       fields{stg},
			handler:      route.RouteInfoH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get route successfully",
			args:         args{ctx, ns1, r1},
			fields:       fields{stg},
			handler:      route.RouteInfoH,
			want:         v1.View().Route().New(r1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := stg.Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = stg.Del(context.Background(), stg.Collection().Route(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), tc.fields.stg.Key().Namespace(ns1.Meta.Name), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Route(), stg.Key().Route(r1.Meta.Namespace, r1.Meta.Name), r1, nil)
			assert.NoError(t, err)

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
			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				n := new(views.Route)
				err := json.Unmarshal(body, &n)
				assert.NoError(t, err)

			}
		})
	}

}

// Testing RouteListH handler
func TestRouteList(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	r1 := getRouteAsset(ns1.Meta.Name, "demo")
	r2 := getRouteAsset(ns1.Meta.Name, "test")

	rl := types.NewRouteMap()
	rl.Items[r1.SelfLink()] = r1
	rl.Items[r2.SelfLink()] = r2

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
		err          string
		want         *types.RouteMap
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get routes list if namespace not found",
			args:         args{ctx, ns2},
			fields:       fields{stg},
			handler:      route.RouteListH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get routes list successfully",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      route.RouteListH,
			want:         rl,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = stg.Del(context.Background(), stg.Collection().Route(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), tc.fields.stg.Key().Namespace(ns1.Meta.Name), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Route(), stg.Key().Route(r1.Meta.Namespace, r1.Meta.Name), r1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Route(), stg.Key().Route(r2.Meta.Namespace, r2.Meta.Name), r2, nil)
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
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				r := new(views.RouteList)
				err := json.Unmarshal(body, &r)
				assert.NoError(t, err)

				for _, item := range *r {
					if _, ok := tc.want.Items[item.Meta.SelfLink]; !ok {
						assert.Error(t, errors.New("not equals"))
					}
				}
			}
		})
	}

}

// Testing RouteCreateH handler
func TestRouteCreate(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	sv1 := getServiceAsset(ns1.Meta.Name, "demo", "")

	sl := new(types.ServiceList)
	sl.Items = append(sl.Items, sv1)

	// initial route in database
	r0 := getRouteAsset(ns1.Meta.Name, "route-0")
	mf0 := getRouteManifest(r0.Meta.Name, sv1.Meta.Name)
	mf0.Spec.Port = 1080
	mf0.Spec.Endpoint = "lastbackend.com"
	mf0.SetRouteSpec(r0, sl)
	mf0s, _ := mf0.ToJson()

	// check unknown namespace
	r1 := getRouteAsset(ns2.Meta.Name, "route-1")
	mf1 := getRouteManifest(r1.Meta.Name, sv1.Meta.Name)
	mf1.SetRouteSpec(r1, sl)
	mf1s, _ := mf1.ToJson()

	// check service not exists
	r2 := getRouteAsset(ns2.Meta.Name, "route-2")
	mf2s, _ := getRouteManifest(r2.Meta.Name, "not_found").ToJson()

	// check tcp port is reserved
	r3 := getRouteAsset(ns1.Meta.Name, "route-3")
	mf3 := getRouteManifest(r3.Meta.Name, sv1.Meta.Name)
	mf3.Spec.Port = mf0.Spec.Port
	mf3.SetRouteSpec(r3, sl)
	mf3s, _ := mf3.ToJson()

	// check endpoint is reserved
	r4 := getRouteAsset(ns1.Meta.Name, "route-4")
	mf4 := getRouteManifest(r4.Meta.Name, sv1.Meta.Name)
	mf4.Spec.Endpoint = mf0.Spec.Endpoint
	mf4.SetRouteSpec(r4, sl)
	mf4s, _ := mf4.ToJson()

	// check successful creation
	r5 := getRouteAsset(ns1.Meta.Name, "route-5")
	mf5 := getRouteManifest(r5.Meta.Name, sv1.Meta.Name)
	mf5.Spec.Endpoint = "lstbknd.net"
	mf5.Spec.Port = 2080
	mf5.SetRouteSpec(r5, sl)
	mf5s, _ := mf5.ToJson()

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
		data         string
		err          string
		want         *views.Route
		wantErr      bool
		expectedCode int
	}{
		// TODO: need checking for uniqueness
		{
			name:         "checking create route if namespace not found",
			args:         args{ctx, ns2},
			fields:       fields{stg},
			handler:      route.RouteCreateH,
			data:         string(mf1s),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking create route if service not found in rules",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      route.RouteCreateH,
			data:         string(mf2s),
			err:          "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad rules parameter\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking create route if exists",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      route.RouteCreateH,
			data:         string(mf0s),
			err:          "{\"code\":400,\"status\":\"Not Unique\",\"message\":\"Name is already in use\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check create route if failed incoming json data",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      route.RouteCreateH,
			data:         "{name:demo}",
			err:          "{\"code\":400,\"status\":\"Incorrect Json\",\"message\":\"Incorrect json\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check create route if tcp port is already use",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      route.RouteCreateH,
			data:         string(mf3s),
			err:          "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Port is already in use\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check create route if endpoint is already use",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      route.RouteCreateH,
			data:         string(mf4s),
			err:          "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Endpoint is already in use\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		// TODO: need checking incoming data for validity
		{
			name:         "check create route success",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      route.RouteCreateH,
			data:         string(mf5s),
			want:         v1.View().Route().New(r5),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = stg.Del(context.Background(), stg.Collection().Route(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Service(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), tc.fields.stg.Key().Namespace(ns1.Meta.Name), ns1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), stg.Key().Service(sv1.Meta.Namespace, sv1.Meta.Name), sv1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Route(), stg.Key().Service(r0.Meta.Namespace, r0.Meta.Name), r0, nil)
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
			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				t.Error(string(body))
				return
			}

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect code message")
			} else {

				t.Log("route self-link:", tc.want.Meta.Namespace, tc.want.Meta.Name)

				got := new(types.Route)
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Route(), stg.Key().Route(tc.want.Meta.Namespace, tc.want.Meta.Name), got, nil)
				if !assert.NoError(t, err) {
					return
				}

				if assert.NotEmpty(t, got, "route is empty") {
					return
				}

				if !assert.Equal(t, tc.want.Meta.Name, got.Meta.Name, "names mismatch") {
					return
				}

				if !assert.Equal(t, len(tc.want.Spec.Rules), len(got.Spec.Rules), "rules count mismatch") {
					return
				}

				assert.Equal(t, tc.want.Spec.Rules[0].Endpoint, got.Spec.Rules[0].Endpoint, "endpoints mismatch")

			}
		})
	}

}

// Testing RouteUpdateH handler
func TestRouteUpdate(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	sv1 := getServiceAsset(ns1.Meta.Name, "demo", "")
	sv2 := getServiceAsset(ns1.Meta.Name, "test1", "")

	sl := new(types.ServiceList)
	sl.Items = append(sl.Items, sv1)

	r0 := getRouteAsset(ns1.Meta.Name, "initial-0")
	r0.Spec.Port = 1080
	r0.Spec.Endpoint = "lstbknd.net"

	r1 := getRouteAsset(ns1.Meta.Name, "initial-1")
	r1.Spec.Port = 1000
	r1.Spec.Endpoint = "lastbackend.com"

	// route not exists
	r2 := getRouteAsset(ns1.Meta.Name, "demo")
	mf2 := getRouteManifest(r2.Meta.Name, sv1.Meta.Name)
	mf2.SetRouteSpec(r2, sl)
	mf2s, _ := mf2.ToJson()

	// invalid namespace
	r3 := getRouteAsset(ns2.Meta.Name, r1.Meta.Name)
	mf3 := getRouteManifest(r3.Meta.Name, sv1.Meta.Name)
	mf3.SetRouteSpec(r1, sl)
	mf3s, _ := mf3.ToJson()

	// invalid service
	r4 := getRouteAsset(ns1.Meta.Name, r1.Meta.Name)
	mf4s, _ := getRouteManifest(r4.Meta.Name, "not_found").ToJson()

	// port in use
	r5 := getRouteAsset(ns1.Meta.Name, r1.Meta.Name)
	mf5 := getRouteManifest(r1.Meta.Name, sv1.Meta.Name)
	mf5.Spec.Port = r0.Spec.Port
	mf5.SetRouteSpec(r5, sl)
	mf5s, _ := mf5.ToJson()

	// endpoint in use
	r6 := getRouteAsset(ns1.Meta.Name, r1.Meta.Name)
	mf6 := getRouteManifest(r1.Meta.Name, sv1.Meta.Name)
	mf6.Spec.Endpoint = r0.Spec.Endpoint
	mf6.SetRouteSpec(r6, sl)
	mf6s, _ := mf5.ToJson()

	// successful update
	r7 := getRouteAsset(ns1.Meta.Name, r1.Meta.Name)
	mf7 := getRouteManifest(r1.Meta.Name, sv1.Meta.Name)
	mf7.Spec.Endpoint = "lbdp.io"
	mf7.Spec.Port = 8080
	mf7.SetRouteSpec(r7, sl)
	mf7s, _ := mf7.ToJson()

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
		data         string
		err          string
		want         *views.Route
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking update route if not exists",
			args:         args{ctx, ns1, r2},
			fields:       fields{stg},
			handler:      route.RouteUpdateH,
			data:         string(mf2s),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Route not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update route if namespace not found",
			args:         args{ctx, ns2, r3},
			fields:       fields{stg},
			handler:      route.RouteUpdateH,
			data:         string(mf3s),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update route if service not found",
			args:         args{ctx, ns2, r4},
			fields:       fields{stg},
			handler:      route.RouteUpdateH,
			data:         string(mf4s),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "check update route if failed incoming json data",
			args:         args{ctx, ns1, r1},
			fields:       fields{stg},
			handler:      route.RouteUpdateH,
			data:         "{name:demo}",
			err:          "{\"code\":400,\"status\":\"Incorrect Json\",\"message\":\"Incorrect json\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking update route if port is in use",
			args:         args{ctx, ns2, r5},
			fields:       fields{stg},
			handler:      route.RouteUpdateH,
			data:         string(mf5s),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update route if endpoint is in use",
			args:         args{ctx, ns2, r6},
			fields:       fields{stg},
			handler:      route.RouteUpdateH,
			data:         string(mf6s),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "check update route success",
			args:         args{ctx, ns1, r7},
			fields:       fields{stg},
			handler:      route.RouteUpdateH,
			data:         string(mf7s),
			want:         v1.View().Route().New(r7),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := stg.Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = stg.Del(context.Background(), stg.Collection().Service(), types.EmptyString)
		assert.NoError(t, err)

		err = stg.Del(context.Background(), stg.Collection().Route(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), tc.fields.stg.Key().Namespace(ns1.Meta.Name), ns1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), stg.Key().Service(sv1.Meta.Namespace, sv1.Meta.Name), sv1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), stg.Key().Service(sv2.Meta.Namespace, sv2.Meta.Name), sv2, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Route(), stg.Key().Route(r0.Meta.Namespace, r0.Meta.Name), r0, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Route(), stg.Key().Route(r1.Meta.Namespace, r1.Meta.Name), r1, nil)
			assert.NoError(t, err)

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
			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				t.Error(string(body))
				return
			}

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect code message")
			} else {
				got := new(types.Route)
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Route(), stg.Key().Route(tc.args.namespace.Meta.Name, tc.want.Meta.Name), got, nil)
				assert.NoError(t, err)
				if assert.NotEmpty(t, got, "route is empty") {
					assert.Equal(t, tc.want.Meta.Name, got.Meta.Name, "names mismatch")
					assert.Equal(t, len(tc.want.Spec.Rules), len(got.Spec.Rules), "rules count mismatch")
					assert.Equal(t, tc.want.Spec.Rules[0].Endpoint, got.Spec.Rules[0].Endpoint, "endpoints mismatch")
				}
			}
		})
	}

}

// Testing RouteRemoveH handler
func TestRouteRemove(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

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
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Route not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get route if namespace not exists",
			args:         args{ctx, ns2, r1},
			fields:       fields{stg},
			handler:      route.RouteRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking del route successfully",
			args:         args{ctx, ns1, r1},
			fields:       fields{stg},
			handler:      route.RouteRemoveH,
			want:         "",
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = stg.Del(context.Background(), stg.Collection().Route(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), tc.fields.stg.Key().Namespace(ns1.Meta.Name), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Route(), stg.Key().Route(r1.Meta.Namespace, r1.Meta.Name), r1, nil)
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
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {
				var got = new(types.Route)
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Route(), stg.Key().Route(tc.args.namespace.Meta.Name, tc.args.route.Meta.Name), got, nil)
				if err != nil {
					assert.NoError(t, err)
				}

				if got != nil {
					assert.Equal(t, got.Status.State, types.StateDestroy, "can not be set to destroy")
				}

				assert.Equal(t, tc.want, string(body), "response not empty")
			}
		})
	}

}

func getNamespaceAsset(name, desc string) *types.Namespace {
	var n = types.Namespace{}
	n.Meta.SetDefault()
	n.Meta.Name = name
	n.Meta.Description = desc
	n.Meta.Endpoint = fmt.Sprintf("%s", name)
	return &n
}

func getServiceAsset(namespace, name, desc string) *types.Service {
	var s = types.Service{}
	s.Meta.SetDefault()
	s.Meta.Namespace = namespace
	s.Meta.Name = name
	s.Meta.Description = desc
	s.Meta.Endpoint = fmt.Sprintf("%s.%s", namespace, name)
	return &s
}

func getRouteAsset(namespace, name string) *types.Route {
	var r = types.Route{}
	r.Meta.SetDefault()
	r.Meta.Namespace = namespace
	r.Meta.Name = name
	r.Spec.Endpoint = fmt.Sprintf("%s.test-domain.com", name)
	r.Spec.Rules = make([]types.RouteRule, 0)
	return &r
}

func getRouteManifest(name, service string) *request.RouteManifest {
	var mf = new(request.RouteManifest)

	mf.Meta.Name = &name
	mf.Spec.Port = 80
	mf.Spec.Rules = make([]request.RouteManifestSpecRulesOption, 0)
	mf.Spec.Rules = append(mf.Spec.Rules, request.RouteManifestSpecRulesOption{
		Port:    80,
		Path:    "/",
		Service: service,
	})

	return mf
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}
