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

package secret_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/lastbackend/lastbackend/internal/api/envs"
	"github.com/lastbackend/lastbackend/internal/master/http/secret"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
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

// Testing SecretListH handler
func TestSecretList(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")

	s1 := getSecretAsset(ns1, "demo")
	s2 := getSecretAsset(ns1, "test")

	s1.Spec.Data["demo"] = []byte("demo")
	s2.Spec.Data["test"] = []byte("test")

	rl := types.NewSecretMap()
	rl.Items[s1.SelfLink().String()] = s1
	rl.Items[s2.SelfLink().String()] = s2

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

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Secret(), s1.SelfLink().String(), &s1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Secret(), s2.SelfLink().String(), &s2, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", fmt.Sprintf("/namespace/%s/secret", ns1.Meta.Name), nil)
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

// Testing SecretCreateH handler
func TestSecretCreate(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")

	s1 := getSecretAsset(ns1, "demo")
	s1.Spec.Type = types.KindSecretOpaque
	s1.Spec.Data["demo"] = []byte(base64.StdEncoding.EncodeToString([]byte("demo")))

	mf1, _ := getSecretManifest(s1).ToJson()

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
			err:          "{\"code\":400,\"status\":\"Incorrect Json\",\"message\":\"Incorrect json\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		// TODO: need checking incoming data for validity
		{
			name:         "check create secret success",
			args:         args{ctx},
			fields:       fields{stg},
			handler:      secret.SecretCreateH,
			data:         string(mf1),
			want:         v1.View().Secret().New(s1),
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

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("POST", fmt.Sprintf("/namespace/%s/secret", s1.Meta.Namespace), strings.NewReader(tc.data))
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

				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Secret(), tc.want.Meta.SelfLink, got, nil)
				if !assert.NoError(t, err) {
					return
				}

				if !assert.Equal(t, s1.Spec.Type, got.Spec.Type, "secret type different") {
					return
				}

				for key, value := range s1.Spec.Data {
					if !assert.NotNil(t, got.Spec.Data[key], "secret key not exists") {
						return
					}

					if !assert.Equal(t, got.Spec.Data[key], value, "secret data not equal") {
						return
					}

				}
			}
		})
	}

}

// Testing SecretUpdateH handler
func TestSecretUpdate(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")

	s1 := getSecretAsset(ns1, "demo")
	s2 := getSecretAsset(ns1, "test")

	s1.Spec.Data["demo"] = []byte(base64.StdEncoding.EncodeToString([]byte("demo")))
	s2.Spec.Data["test"] = []byte(base64.StdEncoding.EncodeToString([]byte("test")))

	mf2, _ := getSecretManifest(s2).ToJson()

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx    context.Context
		secret *types.Secret
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
			args:         args{ctx, s1},
			fields:       fields{stg},
			handler:      secret.SecretUpdateH,
			data:         string(mf2),
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

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Secret(), tc.args.secret.SelfLink().String(), &tc.args.secret, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("PUT", fmt.Sprintf("/namespace/%s/secret/%s", tc.args.secret.Meta.Namespace, tc.args.secret.Meta.Name), strings.NewReader(tc.data))
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

				assert.Equal(t, tc.want.Spec.Data, s.Spec.Data, "secret data mismatch")
			}
		})
	}

}

// Testing SecretRemoveH handler
func TestSecretRemove(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")

	s1 := getSecretAsset(ns1, "demo")
	s2 := getSecretAsset(ns1, "test")

	s1.Spec.Data["demo"] = []byte("demo")
	s2.Spec.Data["test"] = []byte("test")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx    context.Context
		secret *types.Secret
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
			name:         "checking remove secret if not exists",
			args:         args{ctx, s2},
			fields:       fields{stg},
			handler:      secret.SecretRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Secret not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking remove secret successfully",
			args:         args{ctx, s1},
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

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = stg.Put(context.Background(), stg.Collection().Secret(), s1.SelfLink().String(), s1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/namespace/%s/secret/%s", tc.args.secret.Meta.Namespace, tc.args.secret.Meta.Name), nil)
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
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Secret(), tc.args.secret.SelfLink().String(), got, nil)
				if err != nil && !errors.Storage().IsErrEntityNotFound(err) {
					assert.NoError(t, err)
				}

				assert.Equal(t, tc.want, string(body), "response not empty")
			}
		})
	}

}

func getSecretManifest(s *types.Secret) *request.SecretManifest {

	smf := new(request.SecretManifest)

	smf.Meta.Name = &s.Meta.Name
	smf.Meta.Namespace = &s.Meta.Namespace
	smf.Spec.Data = make(map[string]string, 0)

	for key, val := range s.Spec.Data {
		str, _ := base64.StdEncoding.DecodeString(string(val))
		smf.Spec.Data[key] = string(str)
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

func getSecretAsset(namespace *types.Namespace, name string) *types.Secret {
	var s = types.Secret{}
	s.Meta.SetDefault()
	s.Meta.Name = name
	s.Meta.Namespace = namespace.Meta.Name
	s.Meta.SelfLink = *types.NewSecretSelfLink(namespace.Meta.Name, name)

	s.Spec.Type = types.KindSecretOpaque
	s.Spec.Data = make(map[string][]byte, 0)

	return &s
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}
