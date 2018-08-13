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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/http/service"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/stretchr/testify/assert"
)

// Testing ServiceInfoH handler
func TestServiceInfo(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	s1 := getServiceAsset(ns1.Meta.Name, "demo", "")
	s2 := getServiceAsset(ns1.Meta.Name, "test", "")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		service   *types.Service
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		want         *views.Service
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking get service if not exists",
			handler:      service.ServiceInfoH,
			args:         args{ctx, ns1, s2},
			fields:       fields{stg},
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Service not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service if namespace not exists",
			handler:      service.ServiceInfoH,
			args:         args{ctx, ns2, s1},
			fields:       fields{stg},
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service successfully",
			handler:      service.ServiceInfoH,
			args:         args{ctx, ns1, s1},
			fields:       fields{stg},
			want:         v1.View().Service().NewWithDeployment(s1, nil, nil),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
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

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), tc.fields.stg.Key().Service(s1.Meta.Namespace, s1.Meta.Name), s1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", fmt.Sprintf("/namespace/%s/service/%s", tc.args.namespace.Meta.Name, tc.args.service.Meta.Name), nil)
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
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				s := new(views.Service)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, s.Meta.Name, "name not equal")
			}

		})
	}

}

// Testing ServiceListH handler
func TestServiceList(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	s1 := getServiceAsset(ns1.Meta.Name, "demo", "")
	s2 := getServiceAsset(ns1.Meta.Name, "test", "")

	sl := types.NewServiceMap()
	sl.Items[s1.SelfLink()] = s1
	sl.Items[s2.SelfLink()] = s2

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		service   *types.Service
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		err          string
		want         *types.ServiceMap
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get services list if namespace not found",
			args:         args{ctx, ns2, nil},
			fields:       fields{stg},
			handler:      service.ServiceListH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get services list successfully",
			args:         args{ctx, ns1, nil},
			fields:       fields{stg},
			handler:      service.ServiceListH,
			want:         sl,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
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

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), tc.fields.stg.Key().Service(s1.Meta.Namespace, s1.Meta.Name), s1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), tc.fields.stg.Key().Service(s2.Meta.Namespace, s2.Meta.Name), s2, nil)
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

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				s := new(views.RouteList)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				for _, item := range *s {
					if _, ok := tc.want.Items[item.Meta.SelfLink]; !ok {
						assert.Error(t, errors.New("not equals"))
					}
				}
			}

		})
	}

}

type ServiceCreateOptions struct {
	request.ServiceCreateOptions
}

func createServiceCreateOptions(name, description, image *string, spec *request.ServiceOptionsSpec) *ServiceCreateOptions {
	opts := new(ServiceCreateOptions)
	opts.Name = name
	opts.Description = description
	opts.Image = image
	opts.Spec = spec
	return opts
}

func (s *ServiceCreateOptions) toJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

// Testing ServiceCreateH handler
func TestServiceCreate(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	srtPointer := func(s string) *string { return &s }
	intPointer := func(i int) *int { return &i }
	int64Pointer := func(i int64) *int64 { return &i }

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	s1 := getServiceAsset(ns1.Meta.Name, "demo", "")
	s2 := getServiceAsset(ns1.Meta.Name, "test", "")
	s3 := getServiceAsset(ns1.Meta.Name, "new_demo", "")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		service   *types.Service
	}

	tests := []struct {
		name         string
		args         args
		fields       fields
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         string
		want         *views.Service
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking create service if name already exists",
			args:         args{ctx, ns1, s1},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("demo"), nil, srtPointer("redis"), nil).toJson(),
			err:          "{\"code\":400,\"status\":\"Not Unique\",\"message\":\"Name is already in use\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking create service if namespace not found",
			args:         args{ctx, ns2, s2},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("test"), nil, srtPointer("redis"), nil).toJson(),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "check create service if failed incoming json data",
			args:         args{ctx, ns1, s3},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         "{name:demo}",
			err:          "{\"code\":400,\"status\":\"Incorrect json\",\"message\":\"Incorrect json\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check create service if bad parameter name",
			args:         args{ctx, ns1, s3},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("___test"), nil, srtPointer("redis"), nil).toJson(),
			err:          "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad name parameter\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check create service if bad parameter memory",
			args:         args{ctx, ns1, s3},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("test"), nil, srtPointer("redis"), &request.ServiceOptionsSpec{Replicas: intPointer(1), Memory: int64Pointer(127)}).toJson(),
			err:          "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad memory parameter\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		// TODO: check another spec parameters
		{
			name:         "check create service if bad parameter replicas",
			args:         args{ctx, ns1, s3},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("test"), nil, srtPointer("redis"), &request.ServiceOptionsSpec{Replicas: intPointer(-1)}).toJson(),
			err:          "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad replicas parameter\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check create service success",
			args:         args{ctx, ns1, s3},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer(s3.Meta.Name), nil, srtPointer("redis"), nil).toJson(),
			want:         v1.View().Service().NewWithDeployment(s3, nil, nil),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
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

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), tc.fields.stg.Key().Service(s1.Meta.Namespace, s1.Meta.Name), s1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("POST", fmt.Sprintf("/namespace/%s/service", tc.args.namespace.Meta.Name), strings.NewReader(tc.data))
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
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				got := new(types.Service)
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Service(), stg.Key().Service(tc.args.namespace.Meta.Name, tc.args.service.Meta.Name), got, nil)
				assert.NoError(t, err)

				if got == nil {
					t.Error("can not be not nil")
					return
				}

				assert.Equal(t, s3.Meta.Name, got.Meta.Name, "name not equal")
				assert.Equal(t, s3.Meta.Description, got.Meta.Description, "description not equal")
			}
		})
	}

}

type ServiceUpdateOptions struct {
	request.ServiceUpdateOptions
}

func createServiceUpdateOptions(description *string, spec *request.ServiceOptionsSpec) *ServiceUpdateOptions {
	opts := new(ServiceUpdateOptions)
	opts.Description = description
	opts.Spec = spec
	return opts
}

func (s *ServiceUpdateOptions) toJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

// Testing ServiceUpdateH handler
func TestServiceUpdate(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	int64Pointer := func(i int64) *int64 { return &i }

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	s1 := getServiceAsset(ns1.Meta.Name, "demo", "")
	s2 := getServiceAsset(ns1.Meta.Name, "test", "")
	s3 := getServiceAsset(ns1.Meta.Name, "demo", "demo description")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		service   *types.Service
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         string
		want         *views.Service
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking update service if name not exists",
			fields:       fields{stg},
			args:         args{ctx, ns1, s2},
			handler:      service.ServiceUpdateH,
			data:         createServiceUpdateOptions(nil, nil).toJson(),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Service not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update service if namespace not found",
			fields:       fields{stg},
			args:         args{ctx, ns2, s1},
			handler:      service.ServiceUpdateH,
			data:         createServiceUpdateOptions(nil, nil).toJson(),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "check update service if failed incoming json data",
			fields:       fields{stg},
			args:         args{ctx, ns1, s1},
			handler:      service.ServiceUpdateH,
			data:         "{name:demo}",
			err:          "{\"code\":400,\"status\":\"Incorrect json\",\"message\":\"Incorrect json\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check update service if bad parameter memory",
			fields:       fields{stg},
			args:         args{ctx, ns1, s1},
			handler:      service.ServiceUpdateH,
			data:         createServiceUpdateOptions(nil, &request.ServiceOptionsSpec{Memory: int64Pointer(127)}).toJson(),
			err:          "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad memory parameter\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		// TODO: check another spec parameters
		{
			name:         "check update service success",
			fields:       fields{stg},
			args:         args{ctx, ns1, s1},
			handler:      service.ServiceUpdateH,
			data:         createServiceUpdateOptions(&s3.Meta.Description, nil).toJson(),
			want:         v1.View().Service().NewWithDeployment(s3, nil, nil),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
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

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), tc.fields.stg.Key().Service(s1.Meta.Namespace, s1.Meta.Name), s1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("PUT", fmt.Sprintf("/namespace/%s/service/%s", tc.args.namespace.Meta.Name, tc.args.service.Meta.Name), strings.NewReader(tc.data))
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
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {
				s := new(views.Service)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, s.Meta.Name, "description not equal")
				assert.Equal(t, tc.want.Meta.Description, s.Meta.Description, "description not equal")
				assert.Equal(t, tc.want.Spec.Replicas, s.Spec.Replicas, "replicas not equal")
				// TODO: check all updated parameters
			}
		})
	}

}

// Testing ServiceRemoveH handler
func TestServiceRemove(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	s1 := getServiceAsset(ns1.Meta.Name, "demo", "")
	s2 := getServiceAsset(ns1.Meta.Name, "test", "")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		service   *types.Service
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
			name:         "checking get service if not exists",
			fields:       fields{stg},
			args:         args{ctx, ns1, s2},
			handler:      service.ServiceRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Service not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service if namespace not exists",
			fields:       fields{stg},
			args:         args{ctx, ns2, s1},
			handler:      service.ServiceRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service successfully",
			fields:       fields{stg},
			args:         args{ctx, ns1, s1},
			handler:      service.ServiceRemoveH,
			want:         "",
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
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

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), tc.fields.stg.Key().Service(s1.Meta.Namespace, s1.Meta.Name), s1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/namespace/%s/service/%s", tc.args.namespace.Meta.Name, tc.args.service.Meta.Name), nil)
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
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				got := new(types.Service)
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Service(), stg.Key().Service(tc.args.namespace.Meta.Name, tc.args.service.Meta.Name), got, nil)
				if err != nil && !errors.Storage().IsErrEntityNotFound(err) {
					assert.NoError(t, err)
				}

				if got != nil {
					assert.Equal(t, got.Status.State, types.StateDestroy, "status not destroy")
				}

				assert.Equal(t, tc.want, string(body), "response not equal with want")
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

func getServiceAsset(namespace, name, desc string) *types.Service {
	var s = types.Service{}
	s.Meta.SetDefault()
	s.Meta.Namespace = namespace
	s.Meta.Name = name
	s.Meta.Description = desc
	return &s
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}
