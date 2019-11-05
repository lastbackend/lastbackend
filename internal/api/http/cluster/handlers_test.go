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

package cluster_test

import (
	"context"
	"encoding/json"
	"github.com/lastbackend/lastbackend/internal/api/envs"
	"github.com/lastbackend/lastbackend/internal/api/http/cluster"
	"github.com/lastbackend/lastbackend/internal/api/types/v1"
	"github.com/lastbackend/lastbackend/internal/api/types/v1/views"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

// Testing ClusterInfoH handler
func TestClusterInfo(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	c := getClusterAsset(4096)

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
		want         *views.Cluster
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking success get cluster",
			args:         args{ctx, c},
			fields:       fields{stg},
			handler:      cluster.ClusterInfoH,
			want:         v1.View().Cluster().New(c),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Cluster(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := stg.Put(context.Background(), stg.Collection().Cluster(), types.EmptyString, c, nil)
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
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr {
				assert.Error(t, err, "err, should be not nil")
				assert.NotEqual(t, 200, res.Code, "err, should be not nil")
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				s := new(views.Cluster)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Status.Capacity.RAM, s.Status.Capacity.RAM, "memory not equal")
			}
		})
	}
}

func getClusterAsset(memory int64) *types.Cluster {
	var c = types.Cluster{}
	c.Meta.SelfLink = *types.NewClusterSelfLink(types.EmptyString)
	c.Status.Capacity.RAM = memory
	return &c
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}
