//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	envs2 "github.com/lastbackend/lastbackend/internal/api/envs"
	"github.com/lastbackend/lastbackend/internal/api/http/namespace"
	v12 "github.com/lastbackend/lastbackend/internal/api/types/v1"
	request2 "github.com/lastbackend/lastbackend/internal/api/types/v1/request"
	views2 "github.com/lastbackend/lastbackend/internal/api/types/v1/views"
	"github.com/lastbackend/lastbackend/internal/util/resource"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

// Testing NamespaceInfoH handler
func TestNamespaceInfo(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs2.Get().SetStorage(stg)

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
		err          string
		want         *views2.Namespace
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get namespace if not exists",
			args:         args{ctx, ns2},
			fields:       fields{stg},
			handler:      namespace.NamespaceInfoH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			want:         v12.View().Namespace().New(ns1),
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking success get namespace",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      namespace.NamespaceInfoH,
			want:         v12.View().Namespace().New(ns1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs2.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
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
			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			// Check the status code is what we expect.
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				n := new(views2.Namespace)
				err := json.Unmarshal(body, &n)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, n.Meta.Name, "name not equal")
				assert.Equal(t, tc.want.Meta.SelfLink, n.Meta.SelfLink, "name not equal")
			}

		})
	}

}

// Testing NamespaceInfoH handler
func TestNamespaceList(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs2.Get().SetStorage(stg)

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
			want:         nl,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs2.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns2.SelfLink().String(), ns2, nil)
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
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				n := new(views2.NamespaceList)
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

// Testing NamespaceCreateH handler
func TestNamespaceCreate(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs2.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
	}

	nsm1, _ := createNamespaceManifest("demo", "", nil).ToJson()
	nsm2, _ := createNamespaceManifest("__test", "", &request2.NamespaceResourcesOptions{Request: &request2.NamespaceResourceOptions{RAM: getStrPtr("1GB")}}).ToJson()
	nsm3, _ := createNamespaceManifest("test", "", &request2.NamespaceResourcesOptions{Request: &request2.NamespaceResourceOptions{RAM: getStrPtr("1GB")}}).ToJson()

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         string
		err          string
		want         *views2.Namespace
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking create namespace if name already exists",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      namespace.NamespaceCreateH,
			data:         string(nsm1),
			err:          "{\"code\":400,\"status\":\"Not Unique\",\"message\":\"Name is already in use\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking create namespace error if name bad parameter",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      namespace.NamespaceCreateH,
			data:         string(nsm2),
			err:          "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad name parameter\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking create namespace if incorrect json",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      namespace.NamespaceCreateH,
			data:         "{name:demo}",
			err:          "{\"code\":400,\"status\":\"Incorrect Json\",\"message\":\"Incorrect json\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking success create namespace",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      namespace.NamespaceCreateH,
			data:         string(nsm3),
			want:         v12.View().Namespace().New(ns1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs2.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
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
			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			// Check the status code is what we expect.
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				got := new(types.Namespace)
				err = tc.fields.stg.Get(context.Background(), stg.Collection().Namespace(), tc.args.namespace.Meta.Name, got, nil)
				assert.NoError(t, err)
				assert.Equal(t, ns1.Meta.Name, got.Meta.Name, "name not equal")
			}

		})
	}

}

// Testing NamespaceUpdateH handler
func TestNamespaceUpdate(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs2.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("empty", "new description")
	ns3 := getNamespaceAsset("demo", "")
	ns3.Spec.Resources.Request.RAM, _ = resource.DecodeMemoryResource("512MB")

	nsm1, _ := createNamespaceManifest("test", "nil", nil).ToJson()
	nsm3, _ := createNamespaceManifest(ns1.Meta.Name, ns3.Meta.Description, &request2.NamespaceResourcesOptions{Request: &request2.NamespaceResourceOptions{RAM: getStrPtr("512MB")}}).ToJson()

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
		want         *views2.Namespace
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking update namespace if not exists",
			args:         args{ctx, ns2},
			fields:       fields{stg},
			handler:      namespace.NamespaceUpdateH,
			data:         string(nsm1),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			want:         v12.View().Namespace().New(ns1),
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update namespace if not exists",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      namespace.NamespaceUpdateH,
			err:          "{\"code\":400,\"status\":\"Incorrect Json\",\"message\":\"Incorrect json\"}",
			want:         v12.View().Namespace().New(ns1),
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking success update namespace",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      namespace.NamespaceUpdateH,
			data:         string(nsm3),
			want:         v12.View().Namespace().New(ns3),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs2.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
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
			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			// Check the status code is what we expect.
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				n := new(views2.Namespace)
				err := json.Unmarshal(body, &n)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, n.Meta.Name, "name not equal")
				assert.Equal(t, tc.want.Meta.Description, n.Meta.Description, "description not equal")

				if tc.want.Spec.Resources.Request != nil && n.Spec.Resources.Request != nil {
					assert.Equal(t, tc.want.Spec.Resources.Request.RAM, n.Spec.Resources.Request.RAM, "ram not equal")
				}

			}

		})
	}

}

// Testing NamespaceRemoveH handler
func TestNamespaceRemove(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs2.Get().SetStorage(stg)

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
			want:         "",
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "checking remove namespace if name not exists",
			args:         args{ctx, ns2},
			fields:       fields{stg},
			handler:      namespace.NamespaceRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
	}

	clear := func() {
		err := envs2.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
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
			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			// Check the status code is what we expect.
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "")
			} else {

				got := new(types.Namespace)
				err = tc.fields.stg.Get(context.Background(), stg.Collection().Namespace(), tc.args.namespace.SelfLink().String(), got, nil)
				if err != nil && !errors.Storage().IsErrEntityNotFound(err) {
					assert.NoError(t, err)
				}

				if errors.Storage().IsErrEntityNotFound(err) {
					return
				}

				assert.Equal(t, tc.want, string(body), "response not empty")
			}

		})
	}

}

func createNamespaceManifest(name, description string, resources *request2.NamespaceResourcesOptions) *request2.NamespaceManifest {
	opts := new(request2.NamespaceManifest)
	opts.Meta.Name = &name
	opts.Meta.Description = &description
	opts.Spec.Resources = resources
	return opts
}

func getStrPtr(a string) *string {
	return &a
}

func getNamespaceAsset(name, desc string) *types.Namespace {
	var n = types.Namespace{}

	n.Meta.Name = name
	n.Meta.Description = desc

	n.Meta.SelfLink = *types.NewNamespaceSelfLink(name)

	n.Status.Resources.Allocated.RAM, _ = resource.DecodeMemoryResource("512 MB")
	n.Status.Resources.Allocated.CPU, _ = resource.DecodeCpuResource("0.1")

	n.Spec.Resources.Limits.RAM, _ = resource.DecodeMemoryResource("1 GB")
	n.Spec.Resources.Limits.CPU, _ = resource.DecodeCpuResource("0.1")

	n.Spec.Resources.Request.RAM, _ = resource.DecodeMemoryResource("1 GB")
	n.Spec.Resources.Request.CPU, _ = resource.DecodeCpuResource("0.1")

	return &n
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}
