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
	"strings"
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
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		expectedBody string
		expectedCode int
	}{
		{
			name:         "checking get node list successfully",
			handler:      node.NodeListH,
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
			req, err := http.NewRequest("GET", "/cluster/node", nil)
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
			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if res.Code == http.StatusOK {
				assert.Equal(t, tc.expectedBody, string(v), "status code not error")
			} else {
				assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")
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
		expectedBody string
		expectedCode int
	}{
		{
			name:         "checking get node failed: not found",
			url:          fmt.Sprintf("/cluster/node/%s", n2.Meta.Name),
			handler:      node.NodeInfoH,
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Node not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get node successfully",
			url:          fmt.Sprintf("/cluster/node/%s", n1.Meta.Name),
			handler:      node.NodeInfoH,
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
			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")

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
			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")
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
			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")
		})
	}
}

func TestNodeSetMetaH(t *testing.T) {
	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)
	strPointer := func(s string) *string { return &s }

	var (
		ctx = context.Background()

		n1 = getNodeAsset("test1", "", true)
		n2 = getNodeAsset("test2", "", true)
		uo = v1.Request().Node().UpdateOptions()
	)

	uo.Meta = &types.NodeUpdateMetaOptions{
		Provider: strPointer("test"),
	}

	v, err := v1.View().Node().New(&n1).ToJson()
	assert.NoError(t, err)

	type args struct {
		ctx  context.Context
		node string
	}

	tests := []struct {
		name         string
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "checking update node failed: not found",
			args:         args{ctx, n2.Meta.Name},
			handler:      node.NodeSetMetaH,
			data:         uo.ToJson(),
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Node not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update node successfully",
			args:         args{ctx, n1.Meta.Name},
			handler:      node.NodeSetMetaH,
			data:         uo.ToJson(),
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
			req, err := http.NewRequest("PUT", fmt.Sprintf("/cluster/node/%s/meta", tc.args.node), strings.NewReader(tc.data))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/cluster/node/{node}/meta", tc.handler)

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
			assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")

			if tc.expectedCode == http.StatusOK {
				n, err := envs.Get().GetStorage().Node().Get(ctx, tc.args.node)
				assert.NoError(t, err)
				assert.Equal(t, *uo.Meta.Provider, n.Meta.Provider, "provider not equal")
			}

		})
	}
}

func TestNodeConnectH(t *testing.T) {
	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	var (
		err error
		ctx = context.Background()

		n1 = getNodeAsset("test1", "", true)
		n2 = getNodeAsset("test2", "", true)
		uo = v1.Request().Node().NodeConnectOptions()
	)

	uo.Info.Hostname = "test2"
	uo.Info.Architecture = "mac"

	type args struct {
		ctx  context.Context
		node string
	}

	tests := []struct {
		name         string
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "checking create node successful",
			args:         args{ctx, n2.Meta.Name},
			handler:      node.NodeConnectH,
			data:         uo.ToJson(),
			expectedBody: "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "checking update node successful",
			args:         args{ctx, n1.Meta.Name},
			handler:      node.NodeConnectH,
			data:         uo.ToJson(),
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
			req, err := http.NewRequest("PUT", fmt.Sprintf("/cluster/node/%s", tc.args.node), strings.NewReader(tc.data))
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
			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")

			if tc.expectedCode == http.StatusOK {
				n, err := envs.Get().GetStorage().Node().Get(ctx, tc.args.node)
				if assert.NoError(t, err) {
					assert.Equal(t, uo.Info.Hostname, n.Info.Hostname, "hostname not equal")
					assert.Equal(t, uo.Info.Architecture, n.Info.Architecture, "architecture not equal")
				}

			}

		})
	}
}

func TestNodeSetStatusH(t *testing.T) {

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	var (
		err error
		ctx = context.Background()

		n1 = getNodeAsset("test1", "", true)
		n2 = getNodeAsset("test2", "", true)
		uo = v1.Request().Node().NodeStatusOptions()
	)

	uo.Resources.Capacity.Pods = 20
	uo.Resources.Allocated.Containers = 10

	type args struct {
		ctx  context.Context
		node string
	}

	tests := []struct {
		name         string
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "checking update node failed: not found",
			args:         args{ctx, n2.Meta.Name},
			handler:      node.NodeSetStatusH,
			data:         uo.ToJson(),
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Node not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update node successfully",
			args:         args{ctx, n1.Meta.Name},
			handler:      node.NodeSetStatusH,
			data:         uo.ToJson(),
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
			req, err := http.NewRequest("PUT", fmt.Sprintf("/cluster/node/%s/status", tc.args.node), strings.NewReader(tc.data))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/cluster/node/{node}/status", tc.handler)

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
			assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")

			if tc.expectedCode == http.StatusOK {
				n, err := envs.Get().GetStorage().Node().Get(ctx, tc.args.node)
				assert.NoError(t, err)
				assert.Equal(t, uo.Resources.Capacity.Pods, n.Status.Capacity.Pods, "pods not equal")
				assert.Equal(t, uo.Resources.Allocated.Containers, n.Status.Allocated.Containers, "containers not equal")
			}

		})
	}
}

func TestNodeSetPodStatusH(t *testing.T) {

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	var (
		ns  = "ns"
		svc = "svc"
		dp  = "dp"

		err error
		ctx = context.Background()

		n1 = getNodeAsset("test1", "", true)
		n2 = getNodeAsset("test2", "", true)

		p1 = getPodAsset(ns, svc, dp, "test1", "")
		p2 = getPodAsset(ns, svc, dp, "test2", "")

		uo = v1.Request().Node().NodePodStatusOptions()
	)

	uo.State = types.StateError
	uo.Message = "error message"
	uo.Containers = make(map[string]*types.PodContainer)
	uo.Containers["test"] = &types.PodContainer{
		ID: "container-id",
		Image: types.PodContainerImage{
			ID:   "image-id",
			Name: "image-name",
		},
	}

	type args struct {
		ctx  context.Context
		node string
		pod  string
	}

	tests := []struct {
		name         string
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "checking node set pod state failed: node not found",
			args:         args{ctx, n2.Meta.Name, p1.SelfLink()},
			handler:      node.NodeSetPodStatusH,
			data:         uo.ToJson(),
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Node not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking node set pod state failed: pod not found",
			args:         args{ctx, n1.Meta.Name, p2.SelfLink()},
			handler:      node.NodeSetPodStatusH,
			data:         uo.ToJson(),
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Pod not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update node successfully",
			args:         args{ctx, n1.Meta.Name, p1.SelfLink()},
			handler:      node.NodeSetPodStatusH,
			data:         uo.ToJson(),
			expectedBody: "",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		err = envs.Get().GetStorage().Node().Clear(ctx)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Node().Insert(ctx, &n1)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Pod().Insert(ctx, &p1)
		assert.NoError(t, err)

		t.Run(tc.name, func(t *testing.T) {

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("PUT", fmt.Sprintf("/cluster/node/%s/status/pod/%s", tc.args.node, tc.args.pod), strings.NewReader(tc.data))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/cluster/node/{node}/status/pod/{pod}", tc.handler)

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
			assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")

			if tc.expectedCode == http.StatusOK {
				p, err := envs.Get().GetStorage().Pod().Get(ctx, p1.Meta.Namespace, p1.Meta.Service, p1.Meta.Deployment, p1.Meta.Name)
				assert.NoError(t, err)

				assert.Equal(t, uo.State, p.Status.Stage, "pods state not equal")
				assert.Equal(t, uo.Message, p.Status.Message, "pods message not equal")

				uo.Containers = make(map[string]*types.PodContainer)
				uo.Containers["test"] = &types.PodContainer{
					ID: "container-id",
					Image: types.PodContainerImage{
						ID:   "image-id",
						Name: "image-name",
					},
				}

				b := assert.NotNil(t, p.Status.Containers["test"], "container 'test' not exists")
				if !b {
					return
				}

				assert.Equal(t, uo.Containers["test"].ID, p.Status.Containers["test"].ID, "container id not equal")
				assert.Equal(t, uo.Containers["test"].Image.ID, p.Status.Containers["test"].Image.ID, "container image id not equal")
				assert.Equal(t, uo.Containers["test"].Image.Name, p.Status.Containers["test"].Image.Name, "container image name not equal")
			}

		})
	}
}

func TestNodeSetVolumeStatusH(t *testing.T) {

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	var (
		ns = "ns"

		err error
		ctx = context.Background()

		n1 = getNodeAsset("test1", "", true)
		n2 = getNodeAsset("test2", "", true)

		vl1 = getVolumeAsset(ns, "test1", "")
		vl2 = getVolumeAsset(ns, "test2", "")

		uo = v1.Request().Node().NodeVolumeStatusOptions()
	)

	uo.State = types.StateError
	uo.Message = "error message"

	type args struct {
		ctx    context.Context
		node   string
		volume string
	}

	tests := []struct {
		name         string
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "checking ndoe set volume state failed: node not found",
			args:         args{ctx, n2.Meta.Name, vl1.SelfLink()},
			handler:      node.NodeSetVolumeStatusH,
			data:         uo.ToJson(),
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Node not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking ndoe set volume state failed: volume not found",
			args:         args{ctx, n1.Meta.Name, vl2.SelfLink()},
			handler:      node.NodeSetVolumeStatusH,
			data:         uo.ToJson(),
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Volume not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update node successfully",
			args:         args{ctx, n1.Meta.Name, vl1.SelfLink()},
			handler:      node.NodeSetVolumeStatusH,
			data:         uo.ToJson(),
			expectedBody: "",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		err = envs.Get().GetStorage().Node().Clear(ctx)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Node().Insert(ctx, &n1)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Volume().Insert(ctx, &vl1)
		assert.NoError(t, err)

		t.Run(tc.name, func(t *testing.T) {

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("PUT", fmt.Sprintf("/cluster/node/%s/status/volume/%s", tc.args.node, tc.args.volume), strings.NewReader(tc.data))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/cluster/node/{node}/status/volume/{volume}", tc.handler)

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
			assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")

			if tc.expectedCode == http.StatusOK {
				p, err := envs.Get().GetStorage().Volume().Get(ctx, vl1.Meta.Namespace, vl1.Meta.Name)
				assert.NoError(t, err)

				assert.Equal(t, uo.State, p.Status.State, "pods state not equal")
				assert.Equal(t, uo.Message, p.Status.Message, "pods message not equal")
			}

		})
	}
}

func TestNodeSetRouteStatusH(t *testing.T) {

	stg, _ := storage.GetMock()
	envs.Get().SetStorage(stg)
	viper.Set("verbose", 0)

	var (
		ns = "ns"

		err error
		ctx = context.Background()

		n1 = getNodeAsset("test1", "", true)
		n2 = getNodeAsset("test2", "", true)

		r1 = getRouteAsset(ns, "test1", "")
		r2 = getRouteAsset(ns, "test2", "")

		uo = v1.Request().Node().NodeRouteStatusOptions()
	)

	uo.State = types.StateError
	uo.Message = "error message"

	type args struct {
		ctx   context.Context
		node  string
		route string
	}

	tests := []struct {
		name         string
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "checking ndoe set route state failed: node not found",
			args:         args{ctx, n2.Meta.Name, r1.SelfLink()},
			handler:      node.NodeSetRouteStatusH,
			data:         uo.ToJson(),
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Node not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking ndoe set route state failed: route not found",
			args:         args{ctx, n1.Meta.Name, r2.SelfLink()},
			handler:      node.NodeSetRouteStatusH,
			data:         uo.ToJson(),
			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Route not found\"}",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update node successfully",
			args:         args{ctx, n1.Meta.Name, r1.SelfLink()},
			handler:      node.NodeSetRouteStatusH,
			data:         uo.ToJson(),
			expectedBody: "",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {

		err = envs.Get().GetStorage().Node().Clear(ctx)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Node().Insert(ctx, &n1)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Route().Insert(ctx, &r1)
		assert.NoError(t, err)

		t.Run(tc.name, func(t *testing.T) {

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("PUT", fmt.Sprintf("/cluster/node/%s/status/route/%s", tc.args.node, tc.args.route), strings.NewReader(tc.data))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/cluster/node/{node}/status/route/{route}", tc.handler)

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
			assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")

			if tc.expectedCode == http.StatusOK {
				p, err := envs.Get().GetStorage().Route().Get(ctx, r1.Meta.Namespace, r1.Meta.Name)
				assert.NoError(t, err)

				assert.Equal(t, uo.State, p.Status.State, "pods state not equal")
				assert.Equal(t, uo.Message, p.Status.Message, "pods message not equal")
			}

		})
	}
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
		Status: types.NodeStatus{
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
		Network: types.NetworkSpec{
			Type:  types.NetworkTypeVxLAN,
			Range: "10.0.0.1",
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

func getVolumeAsset(namespace, name, desc string) types.Volume {

	var n = types.Volume{}

	n.Meta.Name = name
	n.Meta.Namespace = namespace
	n.Meta.Description = desc

	return n
}

func getRouteAsset(namespace, name, desc string) types.Route {

	var n = types.Route{}

	n.Meta.Name = name
	n.Meta.Namespace = namespace
	n.Meta.Description = desc

	return n
}
