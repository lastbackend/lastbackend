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

package node_test

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/http/node"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Testing NodeList handler
func TestNodeListH(t *testing.T) {

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	var (
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test2", "", false)
		nl  = make(map[string]*types.Node, 0)
	)

	nl[n1.Meta.Name] = &n1
	nl[n2.Meta.Name] = &n2

	v, err := v1.View().Node().NewList(nl).ToJson()
	assert.NoError(t, err)

	tests := []struct {
		name         string
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "checking get node list successfully",
			url:          fmt.Sprintf("/cluster/node"),
			handler:      node.NodeListH,
			description:  "successfully",
			expectedBody: string(v),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		err = envs.Get().GetStorage().Node().Clear(ctx)
		assert.NoError(t, err)

		for _, n := range nl {
			err = envs.Get().GetStorage().Node().Insert(ctx, n)
			assert.NoError(t, err)
		}

		t.Run(tc.name, func(t *testing.T) {

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
			r.HandleFunc("/cluster/node", tc.handler)

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
		})
	}

}

func TestNodeGetH(t *testing.T) {
	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	var (
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test2", "", true)
	)

	v, err := v1.View().Node().New(&n1).ToJson()
	assert.NoError(t, err)

	tests := []struct {
		name         string
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "checking get node failed: not found",
			url:          fmt.Sprintf("/cluster/node/%s", n2.Meta.Name),
			handler:      node.NodeGetH,
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Node not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get node successfully",
			url:          fmt.Sprintf("/cluster/node/%s", n1.Meta.Name),
			handler:      node.NodeGetH,
			description:  "successfully",
			expectedBody: string(v),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		err = envs.Get().GetStorage().Node().Clear(ctx)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Node().Insert(ctx, &n1)
		assert.NoError(t, err)

		t.Run(tc.name, func(t *testing.T) {

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
			r.HandleFunc("/cluster/node/{node}", tc.handler)

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
			assert.Equal(t, tc.expectedBody, string(body), tc.description)

		})
	}
}

func TestNodeGetSpecH(t *testing.T) {
	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	var (
		ns  = "ns"
		svc = "svc"
		dp  = "dp"
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test2", "", true)
		p1  = getPodAsset(ns, svc, dp, "test1", "")
		p2  = getPodAsset(ns, svc, dp, "test2", "")
	)

	n1.Spec.Pods = make(map[string]types.PodSpec)
	n1.Spec.Volumes = make(map[string]types.VolumeSpec)
	n1.Spec.Routes = make(map[string]types.RouteSpec)

	n1.Spec.Pods[p1.SelfLink()] = p1.Spec
	n1.Spec.Pods[p2.SelfLink()] = p2.Spec

	v, err := v1.View().Node().NewSpec(&n1.Spec).ToJson()
	assert.NoError(t, err)

	tests := []struct {
		name         string
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "checking get node spec failed: not found",
			url:          fmt.Sprintf("/cluster/node/%s/spec", n2.Meta.Name),
			handler:      node.NodeGetSpecH,
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Node not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get node spec successfully",
			url:          fmt.Sprintf("/cluster/node/%s/spec", n1.Meta.Name),
			handler:      node.NodeGetSpecH,
			description:  "successfully",
			expectedBody: string(v),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		err = envs.Get().GetStorage().Node().Clear(ctx)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Node().Insert(ctx, &n1)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Node().InsertPod(ctx, &n1, &p1)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Node().InsertPod(ctx, &n1, &p2)
		assert.NoError(t, err)

		t.Run(tc.name, func(t *testing.T) {

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
			r.HandleFunc("/cluster/node/{node}/spec", tc.handler)

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

			assert.Equal(t, tc.expectedBody, string(body), tc.description)
		})
	}
}

func TestNodeRemoveH(t *testing.T) {
	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	var (
		err error
		ctx = context.Background()
		n1  = getNodeAsset("test1", "", true)
		n2  = getNodeAsset("test2", "", true)
	)

	tests := []struct {
		name         string
		url          string
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		description  string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "checking remove failed: not found",
			url:          fmt.Sprintf("/cluster/node/%s", n2.Meta.Name),
			handler:      node.NodeRemoveH,
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Node not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking remove node successfully",
			url:          fmt.Sprintf("/cluster/node/%s", n1.Meta.Name),
			handler:      node.NodeRemoveH,
			expectedBody: "",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		err = envs.Get().GetStorage().Node().Clear(ctx)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Node().Insert(ctx, &n1)
		assert.NoError(t, err)

		t.Run(tc.name, func(t *testing.T) {

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
			r.HandleFunc("/cluster/node/{node}", tc.handler)

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

			assert.Equal(t, tc.expectedBody, string(body), tc.description)
		})
	}
}

func TestNodeUpdateH(t *testing.T) {

}

func TestNodeSetInfoH(t *testing.T) {

}

func TestNodeSetStateH(t *testing.T) {

}

func TestNodeSetPodStatusH(t *testing.T) {

}

func TestNodeSetVolumeStatusH(t *testing.T) {

}

func TestNodeSetRouteStatusH(t *testing.T) {

}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}

func getNodeAsset(name, desc string, online bool) types.Node {
	var n = types.Node{
		Meta: types.NodeMeta{
			Region:   "local",
			Token:    "token",
			Provider: "local",
		},
		Info: types.NodeInfo{
			Hostname: name,
		},
		State: types.NodeState{
			Capacity: types.NodeResources{
				Containers: 2,
				Pods:       2,
				Memory:     1024,
				Cpu:        2,
				Storage:    512,
			},
			Allocated: types.NodeResources{
				Containers: 1,
				Pods:       1,
				Memory:     512,
				Cpu:        1,
				Storage:    256,
			},
		},
		Spec: types.NodeSpec{
			Pods:    make(map[string]types.PodSpec),
			Volumes: make(map[string]types.VolumeSpec),
			Routes:  make(map[string]types.RouteSpec),
		},
		Roles: types.NodeRole{},
		Network: types.Subnet{
			Type:   types.NetworkTypeVxLAN,
			Subnet: "10.0.0.1",
			IFace: types.NetworkInterface{
				Index: 1,
				Name:  "lb",
				Addr:  "10.0.0.1",
				HAddr: "dc:a9:04:83:0d:eb",
			},
		},
		Online: online,
	}

	n.Meta.Name = name
	n.Meta.Description = desc

	return n
}

func getPodAsset(namespace, service, deployment, name, desc string) types.Pod {
	p := types.Pod{}

	p.Meta.Name = name
	p.Meta.Description = desc
	p.Meta.Namespace = namespace
	p.Meta.Service = service
	p.Meta.Deployment = deployment
	p.SelfLink()

	return p
}
