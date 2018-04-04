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

package cluster_test

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/http/cluster"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Testing ClusterInfoH handler
func TestClusterInfo(t *testing.T) {

	var ctx = context.Background()

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)

	c := getClusterAsset("demo", "")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx     context.Context
		cluster *types.Cluster
	}

	tests := []struct {
		name         string
		headers      map[string]string
		fields       fields
		args         args
		handler      func(http.ResponseWriter, *http.Request)
		want         *types.Cluster
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking success get cluster",
			args:         args{ctx, c},
			fields:       fields{stg},
			handler:      cluster.ClusterInfoH,
			want:         c,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Cluster().Clear(context.Background())
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := envs.Get().GetStorage().Cluster().Insert(context.Background(), c)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", "/cluster", nil)
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/cluster", tc.handler)

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

				s := new(views.Cluster)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, s.Meta.Name, "name not equal")
			}
		})
	}
}

type ClusterUpdateOptions struct {
	request.ClusterUpdateOptions
}

func createClusterUpdateOptions(description *string, quotas *request.ClusterQuotasOptions) *ClusterUpdateOptions {
	opts := new(ClusterUpdateOptions)
	opts.Description = description
	opts.Quotas = quotas
	return opts
}

func (s *ClusterUpdateOptions) toJson() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

// Testing NamespaceUpdateH handler
func TestClusterUpdate(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	strPointer := func(s string) *string { return &s }

	c1 := getClusterAsset("demo", "")
	c2 := getClusterAsset("demo", "new description")

	str := make([]string, 1024)
	for i := range str {
		str[i] = "a"
	}
	testDesc := strings.Join(str, "")

	tests := []struct {
		name         string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         string
		expectedCode int
		want         *types.Cluster
		wantErr      bool
		err          string
	}{
		{
			name:         "checking update cluster if bad description parameter",
			data:         createClusterUpdateOptions(&testDesc, &request.ClusterQuotasOptions{}).toJson(),
			handler:      cluster.ClusterUpdateH,
			err:          "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad description parameter\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		// TODO: checking quotas options
		{
			name:         "checking success update cluster",
			data:         createClusterUpdateOptions(strPointer(c2.Meta.Description), &request.ClusterQuotasOptions{}).toJson(),
			handler:      cluster.ClusterUpdateH,
			want:         c2,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Cluster().Clear(context.Background())
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := envs.Get().GetStorage().Cluster().Insert(context.Background(), c1)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("PUT", "/cluster", strings.NewReader(tc.data))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/cluster", tc.handler)

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

				n := new(views.Cluster)
				err := json.Unmarshal(body, &n)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Description, n.Meta.Description, "description not updated")
			}
		})
	}

}

func getClusterAsset(name, desc string) *types.Cluster {
	var c = types.Cluster{}
	c.Meta.Name = name
	c.Meta.Description = desc
	return &c
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}
