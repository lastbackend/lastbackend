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

package volume_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/internal/api/envs"
	"github.com/lastbackend/lastbackend/internal/api/http/volume"
	"github.com/lastbackend/lastbackend/internal/api/types/v1"
	"github.com/lastbackend/lastbackend/internal/api/types/v1/request"
	"github.com/lastbackend/lastbackend/internal/api/types/v1/views"
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

// Testing VolumeInfoH handler
func TestVolumeInfo(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	vl1 := getVolumeAsset(ns1.Meta.Name, "demo")
	vl2 := getVolumeAsset(ns2.Meta.Name, "test")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		volume    *types.Volume
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		err          string
		want         *views.Volume
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get volume if not exists",
			args:         args{ctx, ns1, vl2},
			fields:       fields{stg},
			handler:      volume.VolumeInfoH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Volume not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get volume if namespace not exists",
			args:         args{ctx, ns2, vl1},
			fields:       fields{stg},
			handler:      volume.VolumeInfoH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get volume successfully",
			args:         args{ctx, ns1, vl1},
			fields:       fields{stg},
			handler:      volume.VolumeInfoH,
			want:         v1.View().Volume().New(vl1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := stg.Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = stg.Del(context.Background(), stg.Collection().Volume(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Volume(), vl1.SelfLink().String(), vl1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", fmt.Sprintf("/namespace/%s/volume/%s", tc.args.namespace.Meta.Name, tc.args.volume.Meta.Name), nil)
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/volume/{volume}", tc.handler)

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

				n := new(views.Volume)
				err := json.Unmarshal(body, &n)
				assert.NoError(t, err)

				if assert.Equal(t, tc.want.Meta.Name, n.Meta.Name, "volume name not match") {
					return
				}

				if assert.Equal(t, tc.want.Spec.Selector.Node, n.Spec.Selector.Node, "volume node selector not match") {
					return
				}
			}
		})
	}

}

// Testing VolumeListH handler
func TestVolumeList(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	vl1 := getVolumeAsset(ns1.Meta.Name, "demo")
	vl2 := getVolumeAsset(ns1.Meta.Name, "test")

	vl := types.NewVolumeList()
	vl.Items = append(vl.Items, vl1, vl2)

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
		want         *types.VolumeList
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get volumes list if namespace not found",
			args:         args{ctx, ns2},
			fields:       fields{stg},
			handler:      volume.VolumeListH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get volumes list successfully",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      volume.VolumeListH,
			want:         vl,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = stg.Del(context.Background(), stg.Collection().Volume(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Volume(), vl1.SelfLink().String(), vl1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Volume(), vl2.SelfLink().String(), vl2, nil)
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

				r := new(views.VolumeList)
				err := json.Unmarshal(body, &r)
				assert.NoError(t, err)
				assert.Equal(t, len(*r), len(tc.want.Items), "volumes count not equal")
			}
		})
	}

}

// Testing VolumeCreateH handler
func TestVolumeCreate(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	sv1 := getServiceAsset(ns1.Meta.Name, "demo", "")

	vl1 := getVolumeAsset(ns1.Meta.Name, "demo")

	mf := getVolumeManifest(sv1.Meta.Name)
	mf.SetVolumeSpec(vl1)

	mf1, _ := mf.ToJson()

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
		want         *views.Volume
		wantErr      bool
		expectedCode int
	}{
		// TODO: need checking for unique
		{
			name:         "checking create volume if namespace not found",
			args:         args{ctx, ns2},
			fields:       fields{stg},
			handler:      volume.VolumeCreateH,
			data:         string(mf1),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "check create volume if failed incoming json data",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      volume.VolumeCreateH,
			data:         "{name:demo}",
			err:          "{\"code\":400,\"status\":\"Incorrect Json\",\"message\":\"Incorrect json\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		// TODO: need checking incoming data for validity
		{
			name:         "check create volume success",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      volume.VolumeCreateH,
			data:         string(mf1),
			want:         v1.View().Volume().New(vl1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = stg.Del(context.Background(), stg.Collection().Volume(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Service(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), sv1.SelfLink().String(), sv1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("POST", fmt.Sprintf("/namespace/%s/volume", tc.args.namespace.Meta.Name), strings.NewReader(tc.data))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/volume", tc.handler)

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

				got := new(types.Volume)
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Volume(), tc.want.Meta.SelfLink, got, nil)
				assert.NoError(t, err)
				if assert.NotEmpty(t, got, "volume is empty") {

					if !assert.Equal(t, tc.want.Meta.Name, got.Meta.Name, "names mismatch") {
						return
					}

					if !assert.Equal(t, tc.want.Spec.Selector.Node, got.Spec.Selector.Node, "volume selector node mismatch") {
						return
					}
				}
			}
		})
	}

}

// Testing VolumeUpdateH handler
func TestVolumeUpdate(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	sv1 := getServiceAsset(ns1.Meta.Name, "demo", "")
	sv2 := getServiceAsset(ns1.Meta.Name, "test1", "")
	sv3 := getServiceAsset(ns1.Meta.Name, "test2", "")

	vl1 := getVolumeAsset(ns1.Meta.Name, "demo")
	vl2 := getVolumeAsset(ns1.Meta.Name, "test")
	vl3 := getVolumeAsset(ns1.Meta.Name, "demo")

	vl3.Spec.Selector.Node = "node"
	vl3.Spec.HostPath = "/"
	vl3.Spec.Capacity.Storage, _ = resource.DecodeMemoryResource("1GB")

	mf2, _ := getVolumeManifest(sv2.Meta.Name).ToJson()
	mf3, _ := getVolumeManifest(sv3.Meta.Name).ToJson()

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		volume    *types.Volume
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         string
		err          string
		want         *views.Volume
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking update volume if name not exists",
			args:         args{ctx, ns1, vl2},
			fields:       fields{stg},
			handler:      volume.VolumeUpdateH,
			data:         string(mf3),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Volume not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update volume if namespace not found",
			args:         args{ctx, ns2, vl1},
			fields:       fields{stg},
			handler:      volume.VolumeUpdateH,
			data:         string(mf2),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "check update volume if failed incoming json data",
			args:         args{ctx, ns1, vl1},
			fields:       fields{stg},
			handler:      volume.VolumeUpdateH,
			data:         "{name:demo}",
			err:          "{\"code\":400,\"status\":\"Incorrect Json\",\"message\":\"Incorrect json\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check update volume success",
			args:         args{ctx, ns1, vl1},
			fields:       fields{stg},
			handler:      volume.VolumeUpdateH,
			data:         string(mf2),
			want:         v1.View().Volume().New(vl3),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := stg.Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = stg.Del(context.Background(), stg.Collection().Service(), types.EmptyString)
		assert.NoError(t, err)

		err = stg.Del(context.Background(), stg.Collection().Volume(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), sv1.SelfLink().String(), sv1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), sv2.SelfLink().String(), sv2, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Volume(), vl1.SelfLink().String(), vl1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("PUT", fmt.Sprintf("/namespace/%s/volume/%s", tc.args.namespace.Meta.Name, tc.args.volume.Meta.Name), strings.NewReader(tc.data))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/volume/{volume}", tc.handler)

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
				got := new(types.Volume)
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Volume(), tc.want.Meta.SelfLink, got, nil)
				assert.NoError(t, err)
				if assert.NotEmpty(t, got, "volume is empty") {
					assert.Equal(t, tc.want.Meta.Name, got.Meta.Name, "names mismatch")
					if !assert.Equal(t, tc.want.Meta.Name, got.Meta.Name, "names mismatch") {
						return
					}

					if !assert.Equal(t, tc.want.Spec.Selector.Node, got.Spec.Selector.Node, "volume selector node mismatch") {
						return
					}

					assert.Equal(t, tc.want.Spec.HostPath, got.Spec.HostPath, "hostpath mismatch")
				}
			}
		})
	}

}

// Testing VolumeRemoveH handler
func TestVolumeRemove(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	vl1 := getVolumeAsset(ns1.Meta.Name, "demo")
	vl2 := getVolumeAsset(ns1.Meta.Name, "test")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		volume    *types.Volume
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
			name:         "checking get volume if not exists",
			args:         args{ctx, ns1, vl2},
			fields:       fields{stg},
			handler:      volume.VolumeRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Volume not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get volume if namespace not exists",
			args:         args{ctx, ns2, vl1},
			fields:       fields{stg},
			handler:      volume.VolumeRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking del volume successfully",
			args:         args{ctx, ns1, vl1},
			fields:       fields{stg},
			handler:      volume.VolumeRemoveH,
			want:         "",
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = stg.Del(context.Background(), stg.Collection().Volume(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Volume(), vl1.SelfLink().String(), vl1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/namespace/%s/volume/%s", tc.args.namespace.Meta.Name, tc.args.volume.Meta.Name), nil)
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)

				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/volume/{volume}", tc.handler)

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
				got := new(types.Volume)

				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Volume(), tc.args.volume.SelfLink().String(), got, nil)
				if err != nil && !errors.Storage().IsErrEntityNotFound(err) {
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
	n.Meta.SelfLink = *types.NewNamespaceSelfLink(name)
	return &n
}

func getServiceAsset(namespace, name, desc string) *types.Service {
	var s = types.Service{}
	s.Meta.SetDefault()
	s.Meta.Namespace = namespace
	s.Meta.Name = name
	s.Meta.Description = desc
	s.Meta.Endpoint = fmt.Sprintf("%s.%s", namespace, name)
	s.Meta.SelfLink = *types.NewServiceSelfLink(namespace, name)
	return &s
}

func getVolumeAsset(namespace, name string) *types.Volume {
	var r = types.Volume{}
	r.Meta.SetDefault()
	r.Meta.Namespace = namespace
	r.Meta.Name = name
	r.Spec.Selector.Node = ""
	r.Spec.HostPath = "/"
	r.Spec.Capacity.Storage, _ = resource.DecodeMemoryResource("128MB")
	r.Meta.SelfLink = *types.NewVolumeSelfLink(namespace, name)
	return &r
}

func getVolumeManifest(name string) *request.VolumeManifest {
	var mf = new(request.VolumeManifest)

	mf.Meta.Name = &name
	mf.Spec.Type = types.KindVolumeHostDir
	mf.Spec.Selector.Node = "node"
	mf.Spec.Capacity.Storage = "256MBi"

	return mf
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}
