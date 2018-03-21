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
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/http/service"
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
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

// Testing ServiceInfoH handler
func TestServiceInfo(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

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
		description  string
		want         *types.Service
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking get service if not exists",
			handler:      service.ServiceInfoH,
			description:  "service not found",
			args:         args{ctx, ns1, s2},
			fields:       fields{stg},
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Service not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service if namespace not exists",
			handler:      service.ServiceInfoH,
			description:  "namespace not found",
			args:         args{ctx, ns2, s1},
			fields:       fields{stg},
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service successfully",
			handler:      service.ServiceInfoH,
			description:  "successfully",
			args:         args{ctx, ns1, s1},
			fields:       fields{stg},
			want:         s1,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			err := envs.Get().GetStorage().Namespace().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Service().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Service().Insert(context.Background(), s1)
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
			assert.Equal(t, tc.expectedCode, res.Code, tc.description)

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), tc.description)
			} else {

				s := new(views.Service)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, s.Meta.Name, tc.description)
			}

		})
	}

}

// Testing ServiceListH handler
func TestServiceList(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	s1 := getServiceAsset(ns1.Meta.Name, "demo", "")
	s2 := getServiceAsset(ns1.Meta.Name, "test", "")

	sl := make(types.ServiceList, 0)
	sl[s1.SelfLink()] = s1
	sl[s2.SelfLink()] = s2

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
		description  string
		err          string
		want         types.ServiceList
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get services list if namespace not found",
			args:         args{ctx, ns2, nil},
			fields:       fields{stg},
			handler:      service.ServiceListH,
			description:  "namespace not found",
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get services list successfully",
			args:         args{ctx, ns1, nil},
			fields:       fields{stg},
			handler:      service.ServiceListH,
			description:  "successfully",
			want:         sl,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			err := envs.Get().GetStorage().Namespace().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Service().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Service().Insert(context.Background(), s1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Service().Insert(context.Background(), s2)
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

				s := new(views.ServiceList)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				for _, item := range *s {
					if _, ok := tc.want[item.Meta.SelfLink]; !ok {
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

func createServiceCreateOptions(name, description, sources *string, replicas *int, spec *request.ServiceOptionsSpec) *ServiceCreateOptions {
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

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

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
		description  string
		data         string
		want         *types.Service
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking create service if name already exists",
			description:  "service already exists",
			args:         args{ctx, ns1, s1},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("demo"), nil, srtPointer("redis"), intPointer(1), nil).toJson(),
			err:          "{\"code\":400,\"status\":\"Not Unique\",\"message\":\"Name is already in use\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking create service if namespace not found",
			description:  "namespace not found",
			args:         args{ctx, ns2, s2},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("test"), nil, srtPointer("redis"), intPointer(1), nil).toJson(),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "check create service if failed incoming json data",
			description:  "incoming json data is failed",
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
			description:  "incorrect name parameter",
			args:         args{ctx, ns1, s3},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("___test"), nil, srtPointer("redis"), intPointer(1), nil).toJson(),
			err:          "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad name parameter\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check create service if bad parameter memory",
			description:  "incorrect memory parameter",
			args:         args{ctx, ns1, s3},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("test"), nil, srtPointer("redis"), intPointer(1), &request.ServiceOptionsSpec{Memory: int64Pointer(127)}).toJson(),
			err:          "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad memory parameter\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		// TODO: check another spec parameters
		{
			name:         "check create service if bad parameter replicas",
			description:  "incorrect replicas parameter",
			args:         args{ctx, ns1, s3},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer("test"), nil, srtPointer("redis"), intPointer(-1), nil).toJson(),
			err:          "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad replicas parameter\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check create service success",
			description:  "successfully",
			args:         args{ctx, ns1, s3},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         createServiceCreateOptions(srtPointer(s3.Meta.Name), nil, srtPointer("redis"), intPointer(1), nil).toJson(),
			want:         s3,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			err := envs.Get().GetStorage().Namespace().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Service().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Service().Insert(context.Background(), s1)
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
			assert.Equal(t, tc.expectedCode, res.Code, tc.description)

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), tc.description)
			} else {

				got, _ := tc.fields.stg.Service().Get(tc.args.ctx, tc.args.namespace.Meta.Name, tc.args.service.Meta.Name)
				got.Meta.Updated = tc.want.Meta.Updated
				got.Meta.Created = tc.want.Meta.Created

				assert.Equal(t, s3.Meta.Name, got.Meta.Name, tc.description)
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

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

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
		description  string
		data         string
		want         *types.Service
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking update service if name not exists",
			description:  "service not exists",
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
			description:  "namespace not found",
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
			description:  "incoming json data is failed",
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
			description:  "incorrect memory parameter",
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
			description:  "successfully",
			fields:       fields{stg},
			args:         args{ctx, ns1, s1},
			handler:      service.ServiceUpdateH,
			data:         createServiceUpdateOptions(&s3.Meta.Description, nil).toJson(),
			want:         s3,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			err := envs.Get().GetStorage().Namespace().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Service().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Service().Insert(context.Background(), s1)
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
			assert.Equal(t, tc.expectedCode, res.Code, tc.description)

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), tc.description)
			} else {
				s := new(views.Service)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, s.Meta.Name, tc.description)
				assert.Equal(t, tc.want.Meta.Description, s.Meta.Description, tc.description)
				assert.Equal(t, tc.want.Spec.Replicas, s.Spec.Replicas, tc.description)
				// TODO: check all updated parameters
			}
		})
	}

}

// Testing ServiceRemoveH handler
func TestServiceRemove(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

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
		description  string
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
			description:  "service not found",
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Service not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service if namespace not exists",
			fields:       fields{stg},
			args:         args{ctx, ns2, s1},
			handler:      service.ServiceRemoveH,
			description:  "namespace not found",
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service successfully",
			fields:       fields{stg},
			args:         args{ctx, ns1, s1},
			handler:      service.ServiceRemoveH,
			description:  "successfully",
			want:         "",
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			err := envs.Get().GetStorage().Namespace().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Service().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Service().Insert(context.Background(), s1)
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

func getServiceAsset(namespace, name, desc string) *types.Service {
	var n = types.Service{}

	n.Meta.SetDefault()
	n.Meta.Namespace = namespace
	n.Meta.Name = name
	n.Meta.Description = desc
	return &n
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}
