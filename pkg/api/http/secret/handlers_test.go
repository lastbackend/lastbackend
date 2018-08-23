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

	r1 := getSecretAsset("demo")
	r2 := getSecretAsset("test")

	r1.Data["demo"]=[]byte("demo")
	r2.Data["test"]=[]byte("test")

	rl := types.NewSecretMap()
	rl.Items[r1.SelfLink()] = r1
	rl.Items[r2.SelfLink()] = r2

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
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
			name:         "checking get secrets list successfully",
			args:         args{ctx},
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

			err := stg.Put(context.Background(), stg.Collection().Secret(), stg.Key().Secret(r1.Meta.Name), &r1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Secret(), stg.Key().Secret(r2.Meta.Name), &r2, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", fmt.Sprintf("/secret"), nil)
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/secret", tc.handler)

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

func createSecretCreateOptions(s *types.Secret) *SecretCreateOptions {
	opts := new(SecretCreateOptions)
	opts.Name = s.Meta.Name
	opts.Kind = s.Meta.Kind
	opts.Data = s.Data
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

	r1 := getSecretAsset("demo")
	r1.Meta.Kind = types.KindSecretText
	r1.Data["demo"]=[]byte("demo")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
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
			name:         "check create secret if failed incoming json data",
			args:         args{ctx},
			fields:       fields{stg},
			handler:      secret.SecretCreateH,
			data:         "{name:demo}",
			err:          "{\"code\":400,\"status\":\"Incorrect json\",\"message\":\"Incorrect json\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		// TODO: need checking incoming data for validity
		{
			name:         "check create secret success",
			args:         args{ctx},
			fields:       fields{stg},
			handler:      secret.SecretCreateH,
			data:         createSecretCreateOptions(r1).toJson(),
			want:         v1.View().Secret().New(r1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Secret(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()
			
			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("POST", fmt.Sprintf("/secret"), strings.NewReader(tc.data))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/secret", tc.handler)

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
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Secret(), tc.fields.stg.Key().Secret(tc.want.Meta.Name), got, nil)
				assert.NoError(t, err)

				if !assert.Equal(t, r1.Meta.Kind, got.Meta.Kind, "secret kind different") {
					return
				}

				if !assert.Equal(t, r1.Data, got.Data, "secret kind different") {
					return
				}
			}
		})
	}

}

type SecretUpdateOptions struct {
	request.SecretUpdateOptions
}

func createSecretUpdateOptions(kind string, data map[string][]byte) *SecretUpdateOptions {
	opts := new(SecretUpdateOptions)
	opts.Kind = kind
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

	r1 := getSecretAsset("demo")
	r2 := getSecretAsset("test")

	r1.Data["demo"]=[]byte("demo")
	r2.Data["test"]=[]byte("test")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
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
			args:         args{ctx,  r1},
			fields:       fields{stg},
			handler:      secret.SecretUpdateH,
			data:         createSecretUpdateOptions(r2.Meta.Kind, r2.Data).toJson(),
			want:         v1.View().Secret().New(r2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Secret(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := stg.Put(context.Background(), stg.Collection().Secret(), stg.Key().Secret(r1.Meta.Name), &r1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("PUT", fmt.Sprintf("/secret/%s", tc.args.secret.Meta.Name), strings.NewReader(tc.data))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/secret/{secret}", tc.handler)

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

				assert.Equal(t, tc.want.Data, s.Data, "secret data mismatch")
			}
		})
	}

}

// Testing SecretRemoveH handler
func TestSecretRemove(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.Get("mock")
	envs.Get().SetStorage(stg)

	r1 := getSecretAsset("demo")
	r2 := getSecretAsset("test")

	r1.Data["demo"]=[]byte("demo")
	r2.Data["test"]=[]byte("test")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
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
			args:         args{ctx, r2},
			fields:       fields{stg},
			handler:      secret.SecretRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Secret not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get secret successfully",
			args:         args{ctx, r1},
			fields:       fields{stg},
			handler:      secret.SecretRemoveH,
			want:         "",
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Secret(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := stg.Put(context.Background(), stg.Collection().Secret(), stg.Key().Secret(r1.Meta.Name), &r1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/secret/%s", tc.args.secret.Meta.Name), nil)
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)

				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/secret/{secret}", tc.handler)

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
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Secret(), tc.fields.stg.Key().Secret(tc.args.secret.Meta.Name), got, nil)
				if err != nil && !errors.Storage().IsErrEntityNotFound(err) {
					assert.NoError(t, err)
				}

				assert.Equal(t, tc.want, string(body), "response not empty")
			}
		})
	}

}

func getSecretAsset(name string) *types.Secret {
	var r = types.Secret{}
	r.Meta.SetDefault()
	r.Meta.Name = name
	r.Data = make(map[string][]byte, 0)
	return &r
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}
