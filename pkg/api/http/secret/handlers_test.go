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

package secret_test

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
	"github.com/lastbackend/lastbackend/pkg/api/http/secret"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/stretchr/testify/assert"
)

// Testing SecretListH handler
func TestSecretList(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	r1 := getSecretAsset(ns1.Meta.Name, "demo", "demo")
	r2 := getSecretAsset(ns1.Meta.Name, "test", "test")

	rl := types.NewSecretMap()
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
		want         *types.SecretMap
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get secrets list if namespace not found",
			args:         args{ctx, ns2},
			fields:       fields{stg},
			handler:      secret.SecretListH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get secrets list successfully",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      secret.SecretListH,
			want:         rl,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Secret(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), tc.fields.stg.Key().Namespace(ns1.Meta.Name), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Secret(), stg.Key().Secret(r1.Meta.Namespace, r1.Meta.Name), &r1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Secret(), stg.Key().Secret(r2.Meta.Namespace, r2.Meta.Name), &r2, nil)
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

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				r := new(views.SecretList)
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

type SecretCreateOptions struct {
	request.SecretCreateOptions
}

func createSecretCreateOptions(data *string) *SecretCreateOptions {
	opts := new(SecretCreateOptions)
	opts.Data = data
	return opts
}

func (s *SecretCreateOptions) toJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

// Testing SecretCreateH handler
func TestSecretCreate(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	srtPointer := func(s string) *string { return &s }

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	r1 := getSecretAsset(ns1.Meta.Name, "demo", "demo")

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
		want         *views.Secret
		wantErr      bool
		expectedCode int
	}{
		// TODO: need checking for unique
		{
			name:         "checking create secret if namespace not found",
			args:         args{ctx, ns2},
			fields:       fields{stg},
			handler:      secret.SecretCreateH,
			data:         createSecretCreateOptions(srtPointer(r1.Data)).toJson(),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "check create secret if failed incoming json data",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      secret.SecretCreateH,
			data:         "{name:demo}",
			err:          "{\"code\":400,\"status\":\"Incorrect Json\",\"message\":\"Incorrect json\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		// TODO: need checking incoming data for validity
		{
			name:         "check create secret success",
			args:         args{ctx, ns1},
			fields:       fields{stg},
			handler:      secret.SecretCreateH,
			data:         createSecretCreateOptions(srtPointer(r1.Data)).toJson(),
			want:         v1.View().Secret().New(r1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Secret(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), tc.fields.stg.Key().Namespace(ns1.Meta.Name), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Secret(), stg.Key().Secret(r1.Meta.Namespace, r1.Meta.Name), &r1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("POST", fmt.Sprintf("/namespace/%s/secret", tc.args.namespace.Meta.Name), strings.NewReader(tc.data))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/secret", tc.handler)

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

				got := new(types.Secret)
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Secret(), tc.fields.stg.Key().Secret(tc.args.namespace.Meta.Name, tc.want.Meta.Name), got, nil)
				assert.NoError(t, err)

				assert.Equal(t, ns1.Meta.Name, got.Meta.Name, "it was not be create")
			}
		})
	}

}

type SecretUpdateOptions struct {
	request.SecretUpdateOptions
}

func createSecretUpdateOptions(data *string) *SecretUpdateOptions {
	opts := new(SecretUpdateOptions)
	opts.Data = data
	return opts
}

func (s *SecretUpdateOptions) toJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

// Testing SecretUpdateH handler
func TestSecretUpdate(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	srtPointer := func(s string) *string { return &s }

	ns1 := getNamespaceAsset("demo", "")
	s1 := getSecretAsset(ns1.Meta.Name, "demo", "demo")
	s2 := getSecretAsset(ns1.Meta.Name, "test", "test")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		secret    *types.Secret
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         string
		err          string
		want         *views.Secret
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking update secret if name not exists",
			args:         args{ctx, ns1, s1},
			fields:       fields{stg},
			handler:      secret.SecretUpdateH,
			data:         createSecretUpdateOptions(srtPointer(s2.Data)).toJson(),
			want:         v1.View().Secret().New(s2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Secret(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), tc.fields.stg.Key().Namespace(ns1.Meta.Name), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Secret(), stg.Key().Secret(s1.Meta.Namespace, s1.Meta.Name), &s1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("PUT", fmt.Sprintf("/namespace/%s/secret/%s", tc.args.namespace.Meta.Name, tc.args.secret.Meta.Name), strings.NewReader(tc.data))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/secret/{secret}", tc.handler)

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

				s := new(views.Secret)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Data, s.Data, "it was not be update")
			}
		})
	}

}

// Testing SecretRemoveH handler
func TestSecretRemove(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	s1 := getSecretAsset(ns1.Meta.Name, "demo", "demo")
	s2 := getSecretAsset(ns1.Meta.Name, "test", "test")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		secret    *types.Secret
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
			name:         "checking get secret if not exists",
			args:         args{ctx, ns1, s2},
			fields:       fields{stg},
			handler:      secret.SecretRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Secret not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get secret if namespace not exists",
			args:         args{ctx, ns2, s1},
			fields:       fields{stg},
			handler:      secret.SecretRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get secret successfully",
			args:         args{ctx, ns1, s1},
			fields:       fields{stg},
			handler:      secret.SecretRemoveH,
			want:         "",
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Secret(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), tc.fields.stg.Key().Namespace(ns1.Meta.Name), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Secret(), stg.Key().Secret(s1.Meta.Namespace, s1.Meta.Name), &s1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/namespace/%s/secret/%s", tc.args.namespace.Meta.Name, tc.args.secret.Meta.Name), nil)
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)

				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/secret/{secret}", tc.handler)

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

				got := new(types.Secret)
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Secret(), tc.fields.stg.Key().Secret(tc.args.namespace.Meta.Name, tc.args.secret.Meta.Name), got, nil)
				if err != nil && !errors.Storage().IsErrEntityNotFound(err) {
					assert.NoError(t, err)
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
	return &n
}

func getSecretAsset(namespace, name, data string) *types.Secret {
	var r = types.Secret{}
	r.Meta.SetDefault()
	r.Meta.Namespace = namespace
	r.Meta.Name = name
	r.Data = data
	return &r
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}
