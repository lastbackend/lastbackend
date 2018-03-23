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

package deployment_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/http/deployment"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Testing DeploymentInfoH handler
func TestDeploymentInfo(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	s1 := getServiceAsset(ns1.Meta.Name, "demo", "")
	s2 := getServiceAsset(ns1.Meta.Name, "test", "")
	d1 := getDeploymentAsset(ns1.Meta.Name, s1.Meta.Name, "demo")
	d2 := getDeploymentAsset(ns1.Meta.Name, s2.Meta.Name, "test")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx        context.Context
		namespace  *types.Namespace
		service    *types.Service
		deployment *types.Deployment
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		want         *views.Deployment
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking get deployment if not exists",
			handler:      deployment.DeploymentInfoH,
			args:         args{ctx, ns1, s2, d1},
			fields:       fields{stg},
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Service not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get deployment if namespace not exists",
			handler:      deployment.DeploymentInfoH,
			args:         args{ctx, ns2, s1, d1},
			fields:       fields{stg},
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get deployment if not exists",
			handler:      deployment.DeploymentInfoH,
			args:         args{ctx, ns1, s1, d2},
			fields:       fields{stg},
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Deployment not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get deployment successfully",
			handler:      deployment.DeploymentInfoH,
			args:         args{ctx, ns1, s1, d1},
			fields:       fields{stg},
			want:         v1.View().Deployment().New(d1, nil),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Namespace().Clear(context.Background())
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Service().Clear(context.Background())
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Deployment().Clear(context.Background())
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Service().Insert(context.Background(), s1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Deployment().Insert(context.Background(), d1)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", fmt.Sprintf("/namespace/%s/service/%s/deployment/%s", tc.args.namespace.Meta.Name, tc.args.service.Meta.Name, tc.args.deployment.Meta.Name), nil)
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/service/{service}/deployment/{deployment}", tc.handler)

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

				d := new(views.Deployment)
				err := json.Unmarshal(body, &d)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, d.Meta.Name, "name not equal")
			}

		})
	}

}

// Testing ServiceListH handler
func TestServiceList(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	s1 := getServiceAsset(ns1.Meta.Name, "demo", "")
	s2 := getServiceAsset(ns2.Meta.Name, "test", "")
	d1 := getDeploymentAsset(ns1.Meta.Name, s1.Meta.Name, "demo")
	d2 := getDeploymentAsset(ns1.Meta.Name, s2.Meta.Name, "test")

	dl := make(types.DeploymentMap, 0)
	dl[d1.SelfLink()] = d1
	dl[d2.SelfLink()] = d2

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
		want         types.DeploymentMap
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get deployment list if service not exists",
			handler:      deployment.DeploymentListH,
			args:         args{ctx, ns1, s2},
			fields:       fields{stg},
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Service not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get deployment list if namespace not exists",
			handler:      deployment.DeploymentListH,
			args:         args{ctx, ns2, s1},
			fields:       fields{stg},
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get deployments list successfully",
			args:         args{ctx, ns1, s1},
			fields:       fields{stg},
			handler:      deployment.DeploymentListH,
			want:         dl,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Namespace().Clear(context.Background())
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Service().Clear(context.Background())
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Deployment().Clear(context.Background())
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := envs.Get().GetStorage().Namespace().Insert(context.Background(), ns1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Service().Insert(context.Background(), s1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Service().Insert(context.Background(), s2)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Deployment().Insert(context.Background(), d1)
			assert.NoError(t, err)

			err = envs.Get().GetStorage().Deployment().Insert(context.Background(), d2)
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
			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				d := new(views.DeploymentList)
				err := json.Unmarshal(body, &d)
				assert.NoError(t, err)

				for _, item := range *d {
					if _, ok := tc.want[item.Meta.SelfLink]; !ok {
						assert.Error(t, errors.New("not equals"))
					}
				}
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

func getDeploymentAsset(namespace, service, name string) *types.Deployment {
	var d = types.Deployment{}
	d.Meta.SetDefault()
	d.Meta.Namespace = namespace
	d.Meta.Service = service
	d.Meta.Name = name
	return &d
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}
