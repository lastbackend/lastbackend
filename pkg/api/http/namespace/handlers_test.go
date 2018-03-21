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
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"errors"
)

// Testing NamespaceInfoH handler
func TestNamespaceInfo(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 7)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	nl := make(map[string]*types.Namespace)
	nl[ns1.Meta.Name] = ns1
	nl[ns2.Meta.Name] = ns2

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
		want         *views.Namespace
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get namespace if not exists",
			args:         args{ctx, ns2},
			fields:       fields{stg},
			handler:      namespace.NamespaceInfoH,
			description:  "namespace not found",
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking success get namespace",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      namespace.NamespaceInfoH,
			description:  "successfully",
			want:         v1.View().Namespace().New(ns1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := envs.Get().GetStorage().Namespace().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
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

				n := new(views.Namespace)
				err := json.Unmarshal(body, &n)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, n.Meta.Name, tc.description)
			}

		})
	}

}

// Testing NamespaceInfoH handler
func TestNamespaceList(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	nl := make(map[string]*types.Namespace)
	nl[ns1.Meta.Name] = ns1
	nl[ns2.Meta.Name] = ns2

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
		want         map[string]*types.Namespace
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking success get namespace list",
			args:         args{ctx, nil},
			fields:       fields{stg},
			handler:      namespace.NamespaceListH,
			description:  "successfully",
			want:         nl,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			err := envs.Get().GetStorage().Namespace().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns2)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", "/namespace", nil)
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

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), tc.description)
			} else {

				n := new(views.NamespaceList)
				err := json.Unmarshal(body, &n)
				assert.NoError(t, err)

				for _, item := range *n {
					if _, ok := tc.want[item.Meta.SelfLink]; !ok {
						assert.Error(t, errors.New("not equals"))
					}
				}
			}
		})
	}

}

type NamespaceCreateOptions struct {
	request.NamespaceCreateOptions
}

func createNamespaceCreateOptions(name, description string, quotas *request.NamespaceQuotasOptions) *NamespaceCreateOptions {
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

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	ns1 := getNamespaceAsset("demo", "")

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
		want         *views.Namespace
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking create namespace if name already exists",
			description:  "namespace already exists",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      namespace.NamespaceCreateH,
			data:         createNamespaceCreateOptions("demo", "", nil).toJson(),
			err:          "{\"code\":400,\"status\":\"Not Unique\",\"message\":\"Name is already in use\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking create namespace error if name bad parameter",
			description:  "successfully",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      namespace.NamespaceCreateH,
			data:         createNamespaceCreateOptions("__test", "", &request.NamespaceQuotasOptions{RAM: 2, Routes: 1}).toJson(),
			err:          "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad name parameter\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking create namespace if incorrect json",
			description:  "successfully",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      namespace.NamespaceCreateH,
			data:         "{name:demo}",
			err:          "{\"code\":400,\"status\":\"Incorrect json\",\"message\":\"Incorrect json\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking success create namespace",
			description:  "successfully",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      namespace.NamespaceCreateH,
			data:         createNamespaceCreateOptions("test", "", &request.NamespaceQuotasOptions{RAM: 2, Routes: 1}).toJson(),
			want:         v1.View().Namespace().New(ns1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			err := envs.Get().GetStorage().Namespace().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("POST", "/namespace", strings.NewReader(tc.data))
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

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), tc.description)
			} else {

				got, err := tc.fields.stg.Namespace().Get(tc.args.ctx, tc.args.namespace.Meta.Name)
				assert.NoError(t, err)

				assert.Equal(t, ns1.Meta.Name, got.Meta.Name, tc.description)
			}

		})
	}

}

type NamespaceUpdateOptions struct {
	request.NamespaceUpdateOptions
}

func createNamespaceUpdateOptions(description *string, quotas *request.NamespaceQuotasOptions) *NamespaceUpdateOptions {
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

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("demo", "")
	ns2.Resources.RAM = 512
	ns2.Resources.Routes = 2
	ns3 := getNamespaceAsset("empty", "")

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
		want         *views.Namespace
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking update namespace if not exists",
			args:         args{ctx, ns3},
			fields:       fields{stg},
			handler:      namespace.NamespaceUpdateH,
			description:  "namespace not exists",
			data:         createNamespaceUpdateOptions(nil, nil).toJson(),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update namespace if not exists",
			description:  "",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      namespace.NamespaceUpdateH,
			data:         "{description:demo}",
			err:          "{\"code\":400,\"status\":\"Incorrect json\",\"message\":\"Incorrect json\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking success update namespace",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      namespace.NamespaceUpdateH,
			description:  "successfully",
			data:         createNamespaceUpdateOptions(nil, &request.NamespaceQuotasOptions{RAM: ns2.Resources.RAM, Routes: ns2.Resources.Routes}).toJson(),
			want:         v1.View().Namespace().New(ns1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			err := envs.Get().GetStorage().Namespace().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("PUT", fmt.Sprintf("/namespace/%s", tc.args.namespace.Meta.Name), strings.NewReader(tc.data))
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

				n := new(views.Namespace)
				err := json.Unmarshal(body, &n)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, n.Meta.Name, tc.description)
				assert.Equal(t, tc.want.Meta.Description, n.Meta.Description, tc.description)
			}

		})
	}

}

// Testing NamespaceRemoveH handler
func TestNamespaceRemove(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

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
		want         string
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking success remove namespace",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      namespace.NamespaceRemoveH,
			description:  "successfully",
			want:         "",
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "checking remove namespace if name not exists",
			args:         args{ctx, ns2},
			fields:       fields{stg},
			handler:      namespace.NamespaceRemoveH,
			description:  "namespace not found",
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			err := envs.Get().GetStorage().Namespace().Clear(context.Background())
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/namespace/%s", tc.args.namespace.Meta.Name), nil)
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
				assert.Equal(t, tc.want, string(body), tc.description)
			}

		})
	}

}

func getNamespaceAsset(name, desc string) *types.Namespace {
	var n = types.Namespace{}

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
