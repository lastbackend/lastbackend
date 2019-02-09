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
// patents in process, and are protected by trade config or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package config_test

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
	"github.com/lastbackend/lastbackend/pkg/api/http/config"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/stretchr/testify/assert"
)

// Testing ConfigListH handler
func TestConfigList(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")

	c1 := getConfigAsset(ns1, "demo")
	c2 := getConfigAsset(ns1, "test")

	c1.Spec.Data["test.txt"] = "test1"
	c2.Spec.Data["test.txt"] = "test2"

	cl := types.NewConfigMap()
	cl.Items[c1.SelfLink().String()] = c1
	cl.Items[c2.SelfLink().String()] = c2

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		err          string
		want         *types.ConfigMap
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get configs list successfully",
			args:         args{ctx},
			fields:       fields{stg},
			handler:      config.ConfigListH,
			want:         cl,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Config(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Config(), c1.SelfLink().String(), &c1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Config(), c2.SelfLink().String(), &c2, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", fmt.Sprintf("/namespace/%s/config", ns1.Meta.Name), nil)
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/config", tc.handler)

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

				r := new(views.ConfigList)
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

// Testing ConfigCreateH handler
func TestConfigCreate(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")

	c1 := getConfigAsset(ns1, "demo")
	c1.Meta.Kind = types.KindConfigText
	c1.Spec.Data["test.txt"] = "test1"

	mf1, _ := getConfigManifest(c1).ToJson()

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         string
		err          string
		want         *views.Config
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "check create config if failed incoming json data",
			args:         args{ctx},
			fields:       fields{stg},
			handler:      config.ConfigCreateH,
			data:         "{name:demo}",
			err:          "{\"code\":400,\"status\":\"Incorrect Json\",\"message\":\"Incorrect json\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		// TODO: need checking incoming data for validity
		{
			name:         "check create config success",
			args:         args{ctx},
			fields:       fields{stg},
			handler:      config.ConfigCreateH,
			data:         string(mf1),
			want:         v1.View().Config().New(c1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {

		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Config(), types.EmptyString)
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
			req, err := http.NewRequest("POST", fmt.Sprintf("/namespace/%s/config", c1.Meta.Namespace), strings.NewReader(tc.data))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/config", tc.handler)

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

				got := new(types.Config)
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Config(), tc.want.Meta.SelfLink, got, nil)
				assert.NoError(t, err)

				if !assert.Equal(t, c1.Spec.Type, got.Spec.Type, "config kind different") {
					return
				}

				if !assert.Equal(t, len(c1.Spec.Data), len(got.Spec.Data), "config kind different") {
					return
				}
			}
		})
	}

}

// Testing ConfigUpdateH handler
func TestConfigUpdate(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")

	c1 := getConfigAsset(ns1, "demo")
	c1.Meta.Kind = types.KindConfigText
	c1.Spec.Data["test.txt"] = "test1"

	c2 := getConfigAsset(ns1, "demo")
	c2.Meta.Kind = types.KindConfigText
	c2.Spec.Data["test.txt"] = "test2"
	c2.Spec.Data["test.cfg"] = "cfg"

	mf2, _ := getConfigManifest(c2).ToJson()

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx    context.Context
		config *types.Config
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         string
		err          string
		want         *views.Config
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking update config if name not exists",
			args:         args{ctx, c1},
			fields:       fields{stg},
			handler:      config.ConfigUpdateH,
			data:         string(mf2),
			want:         v1.View().Config().New(c2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Config(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Config(),
				tc.args.config.SelfLink().String(), &tc.args.config, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("PUT", fmt.Sprintf("/namespace/%s/config/%s",
				tc.args.config.Meta.Namespace, tc.args.config.Meta.Name), strings.NewReader(tc.data))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/config/{config}", tc.handler)

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

				s := new(views.Config)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				assert.Equal(t, len(tc.want.Spec.Data), len(s.Spec.Data), "config data mismatch")
			}
		})
	}

}

// Testing ConfigRemoveH handler
func TestConfigRemove(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")

	c1 := getConfigAsset(ns1, "demo")
	c1.Meta.Kind = types.KindConfigText
	c1.Spec.Data["test.txt"] = "test1"

	c2 := getConfigAsset(ns1, "test")
	c2.Meta.Kind = types.KindConfigText
	c2.Spec.Data["test.txt"] = "test2"
	c2.Spec.Data["test.cfg"] = "cfg"

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx    context.Context
		config *types.Config
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
			name:         "checking remove config if not exists",
			args:         args{ctx, c2},
			fields:       fields{stg},
			handler:      config.ConfigRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Config not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking remove config successfully",
			args:         args{ctx, c1},
			fields:       fields{stg},
			handler:      config.ConfigRemoveH,
			want:         "",
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Config(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Config(),
				c1.SelfLink().String(), c1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/namespace/%s/config/%s",
				tc.args.config.Meta.Namespace, tc.args.config.Meta.Name), nil)
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)

				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/config/{config}", tc.handler)

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

				got := new(types.Config)
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Config(),
					tc.args.config.SelfLink().String(), got, nil)
				if err != nil && !errors.Storage().IsErrEntityNotFound(err) {
					assert.NoError(t, err)
				}

				assert.Equal(t, tc.want, string(body), "response not empty")
			}
		})
	}

}

func getConfigManifest(s *types.Config) *request.ConfigManifest {

	smf := new(request.ConfigManifest)

	smf.Meta.Name = &s.Meta.Name
	smf.Meta.Namespace = &s.Meta.Namespace
	smf.Spec.Data = make(map[string]string, 0)

	for key, val := range s.Spec.Data {
		smf.Spec.Data[key] = val
	}

	smf.Spec.Type = s.Spec.Type

	return smf
}

func getNamespaceAsset(name, desc string) *types.Namespace {
	var n = types.Namespace{}
	n.Meta.SetDefault()
	n.Meta.Name = name
	n.Meta.Description = desc
	n.Meta.SelfLink = *types.NewNamespaceSelfLink(name)
	return &n
}

func getConfigAsset(namespace *types.Namespace, name string) *types.Config {
	var c = types.Config{}
	c.Meta.SetDefault()
	c.Meta.Name = name
	c.Meta.Namespace = namespace.Meta.Name
	c.Meta.SelfLink = *types.NewConfigSelfLink(namespace.Meta.Name, name)
	c.Spec.Type = types.KindConfigText
	c.Spec.Data = make(map[string]string, 0)
	return &c
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}
