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
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/http/cluster"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"encoding/json"
)

// Testing ClusterInfoH handler
func TestClusterInfo(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 7)

	c := getClusterAsset("demo", "")
	err := envs.Get().GetStorage().Cluster().Insert(context.Background(), c)
	assert.NoError(t, err)

	v, err := v1.View().Cluster().New(c).ToJson()
	assert.NoError(t, err)

	tests := []struct {
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		expectedBody string
		expectedCode int
	}{
		{
			url:          "/cluster",
			handler:      cluster.ClusterInfoH,
			description:  "successfully",
			expectedBody: string(v),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("GET", tc.url, nil)
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
		assert.Equal(t, tc.expectedCode, res.Code, tc.description)

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		if res.Code == http.StatusOK {
			assert.Equal(t, tc.expectedBody, string(v), tc.description)
		} else {
			assert.Equal(t, tc.expectedBody, string(body), tc.description)
		}
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
func TestNamespaceUpdate(t *testing.T) {

	strg, _ := storage.GetMock()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	c := getClusterAsset("demo", "")

	err := envs.Get().GetStorage().Cluster().Insert(context.Background(), c)
	assert.NoError(t, err)

	v, err := v1.View().Cluster().New(c).ToJson()
	assert.NoError(t, err)

	str := make([]string, 1024)
	for i := range str {
		str[i] = "a"
	}
	testDesc := strings.Join(str, "")

	tests := []struct {
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		data         string
		expectedBody string
		expectedCode int
	}{
		{
			url:          "/cluster",
			description:  "successfully",
			data:         createClusterUpdateOptions(&testDesc, &request.ClusterQuotasOptions{}).toJson(),
			handler:      cluster.ClusterUpdateH,
			expectedBody: "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad description parameter\"}",
			expectedCode: http.StatusBadRequest,
		},
		// TODO: checking quotas options
		{
			url:          "/cluster",
			description:  "successfully",
			data:         createClusterUpdateOptions(nil, &request.ClusterQuotasOptions{}).toJson(),
			handler:      cluster.ClusterUpdateH,
			expectedBody: string(v),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("PUT", tc.url, strings.NewReader(tc.data))
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
		assert.Equal(t, tc.expectedCode, res.Code, tc.description)

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		if res.Code == http.StatusOK {
			assert.Equal(t, tc.expectedBody, string(v), tc.description)
		} else {
			assert.Equal(t, tc.expectedBody, string(body), tc.description)
		}

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
