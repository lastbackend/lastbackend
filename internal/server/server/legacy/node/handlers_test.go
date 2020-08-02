//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
//
//import (
//	"context"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//
//	"github.com/gorilla/mux"
//	"github.com/lastbackend/lastbackend/internal/api/envs"
//	"github.com/lastbackend/lastbackend/internal/master/cache"
//	"github.com/lastbackend/lastbackend/internal/master/http/node"
//	"github.com/lastbackend/lastbackend/internal/pkg/models"
//	"github.com/lastbackend/lastbackend/internal/pkg/storage"
//	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
//	"github.com/spf13/viper"
//	"github.com/stretchr/testify/assert"
//)
//
//// Testing NodeList handler
//func TestNodeListH(t *testing.T) {
//
//	v := viper.New()
//	v.SetDefault("storage.driver", "mock")
//
//	stg, _ := storage.Get(v)
//	envs.Get().SetStorage(stg)
//	v.Set("verbose", 0)
//
//	var (
//		n1 = getNodeAsset("test1", "", true)
//		n2 = getNodeAsset("test2", "", false)
//		nl = models.NewNodeList()
//	)
//
//	nl.Items = append(nl.Items, &n1)
//	nl.Items = append(nl.Items, &n2)
//
//	view, err := v1.View().Node().NewList(nl).ToJson()
//	assert.NoError(t, err)
//
//	tests := []struct {
//		name         string
//		headers      map[string]string
//		handler      func(http.ResponseWriter, *http.Request)
//		expectedBody string
//		expectedCode int
//	}{
//		{
//			name:         "checking get node list successfully",
//			handler:      node.NodeListH,
//			expectedBody: string(view),
//			expectedCode: http.StatusOK,
//		},
//	}
//
//	for _, tc := range tests {
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Node().Info(), models.EmptyString)
//		assert.NoError(t, err)
//
//		for _, n := range nl.Items {
//			err = stg.Put(context.Background(), stg.Collection().Node().Info(), n.SelfLink().String(), &n, nil)
//			assert.NoError(t, err)
//		}
//
//		t.Run(tc.name, func(t *testing.T) {
//
//			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
//			// pass 'nil' as the third parameter.
//			req, err := http.NewRequest("GET", "/cluster/node", nil)
//			assert.NoError(t, err)
//
//			if tc.headers != nil {
//				for key, val := range tc.headers {
//					req.Header.Set(key, val)
//				}
//			}
//
//			r := mux.NewRouter()
//			r.HandleFunc("/cluster/node", tc.handler)
//
//			setRequestVars(r, req)
//
//			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
//			res := httptest.NewRecorder()
//
//			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
//			// directly and pass in our Request and ResponseRecorder.
//			r.ServeHTTP(res, req)
//
//			// Check the status code is what we expect.
//			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")
//
//			body, err := ioutil.ReadAll(res.Body)
//			assert.NoError(t, err)
//
//			if res.Code == http.StatusOK {
//				assert.Equal(t, tc.expectedBody, string(view), "status code not error")
//			} else {
//				assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")
//			}
//		})
//	}
//
//}
//
//func TestNodeGetH(t *testing.T) {
//
//	v := viper.New()
//	v.SetDefault("storage.driver", "mock")
//
//	stg, _ := storage.Get(v)
//	envs.Get().SetStorage(stg)
//	v.Set("verbose", 0)
//
//	var (
//		n1 = getNodeAsset("test1", "", true)
//		n2 = getNodeAsset("test2", "", true)
//	)
//
//	view, err := v1.View().Node().New(&n1).ToJson()
//	assert.NoError(t, err)
//
//	tests := []struct {
//		name         string
//		url          string
//		headers      map[string]string
//		handler      func(http.ResponseWriter, *http.Request)
//		expectedBody string
//		expectedCode int
//	}{
//		{
//			name:         "checking get node failed: not found",
//			url:          fmt.Sprintf("/cluster/node/%s", n2.Meta.Name),
//			handler:      node.NodeInfoH,
//			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Node not found\"}",
//			expectedCode: http.StatusNotFound,
//		},
//		{
//			name:         "checking get node successfully",
//			url:          fmt.Sprintf("/cluster/node/%s", n1.Meta.Name),
//			handler:      node.NodeInfoH,
//			expectedBody: string(view),
//			expectedCode: http.StatusOK,
//		},
//	}
//
//	for _, tc := range tests {
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Node().Info(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = stg.Put(context.Background(), stg.Collection().Node().Info(), n1.SelfLink().String(), &n1, nil)
//		assert.NoError(t, err)
//
//		t.Run(tc.name, func(t *testing.T) {
//
//			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
//			// pass 'nil' as the third parameter.
//			req, err := http.NewRequest("GET", tc.url, nil)
//			assert.NoError(t, err)
//
//			if tc.headers != nil {
//				for key, val := range tc.headers {
//					req.Header.Set(key, val)
//				}
//			}
//
//			r := mux.NewRouter()
//			r.HandleFunc("/cluster/node/{node}", tc.handler)
//
//			setRequestVars(r, req)
//
//			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
//			res := httptest.NewRecorder()
//
//			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
//			// directly and pass in our Request and ResponseRecorder.
//			r.ServeHTTP(res, req)
//
//			// Check the status code is what we expect.
//			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")
//
//			body, err := ioutil.ReadAll(res.Body)
//			assert.NoError(t, err)
//			assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")
//
//		})
//	}
//}
//
//func TestNodeGetManifestH(t *testing.T) {
//	v := viper.New()
//	v.SetDefault("storage.driver", "mock")
//
//	stg, _ := storage.Get(v)
//	cg := cache.NewCache()
//
//	envs.Get().SetStorage(stg)
//	envs.Get().SetCache(cg)
//
//	v.Set("verbose", 0)
//
//	var (
//		n1 = getNodeAsset("test1", "", true)
//		n2 = getNodeAsset("test2", "", true)
//		p1 = "test1"
//		p2 = "test2"
//
//		nm = new(models.NodeManifest)
//	)
//
//	nm.Meta.Initial = true
//	nm.Exporter = new(models.ExporterManifest)
//	nm.Resolvers = make(map[string]*models.ResolverManifest, 0)
//	nm.Network = make(map[string]*models.SubnetManifest, 0)
//	nm.Pods = make(map[string]*models.PodManifest, 0)
//	nm.Pods[p1] = getPodManifest()
//	nm.Pods[p2] = getPodManifest()
//	nm.Volumes = make(map[string]*models.VolumeManifest, 0)
//	nm.Endpoints = make(map[string]*models.EndpointManifest, 0)
//
//	view, err := v1.View().Node().NewManifest(nm).ToJson()
//	assert.NoError(t, err)
//
//	tests := []struct {
//		name         string
//		url          string
//		headers      map[string]string
//		handler      func(http.ResponseWriter, *http.Request)
//		expectedBody string
//		expectedCode int
//	}{
//		{
//			name:         "node spec failed not found",
//			url:          fmt.Sprintf("/cluster/node/%s/spec", n2.Meta.Name),
//			handler:      node.NodeGetSpecH,
//			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Node not found\"}",
//			expectedCode: http.StatusNotFound,
//		},
//		{
//			name:         "node spec successfully",
//			url:          fmt.Sprintf("/cluster/node/%s/spec", n1.Meta.Name),
//			handler:      node.NodeGetSpecH,
//			expectedBody: string(view),
//			expectedCode: http.StatusOK,
//		},
//	}
//
//	for _, tc := range tests {
//
//		t.Run(tc.name, func(t *testing.T) {
//
//			err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Node().Info(), models.EmptyString)
//			assert.NoError(t, err)
//
//			err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Manifest().Pod(n1.Meta.Name), models.EmptyString)
//			assert.NoError(t, err)
//
//			err = stg.Put(context.Background(), stg.Collection().Node().Info(), n1.SelfLink().String(), &n1, nil)
//			assert.NoError(t, err)
//
//			err = stg.Put(context.Background(), stg.Collection().Manifest().Pod(n1.Meta.Name), p1, getPodManifest(), nil)
//			assert.NoError(t, err)
//
//			err = stg.Put(context.Background(), stg.Collection().Manifest().Pod(n1.Meta.Name), p2, getPodManifest(), nil)
//			assert.NoError(t, err)
//
//			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
//			// pass 'nil' as the third parameter.
//			req, err := http.NewRequest("GET", tc.url, nil)
//			assert.NoError(t, err)
//
//			if tc.headers != nil {
//				for key, val := range tc.headers {
//					req.Header.Set(key, val)
//				}
//			}
//
//			r := mux.NewRouter()
//			r.HandleFunc("/cluster/node/{node}/spec", tc.handler)
//
//			setRequestVars(r, req)
//
//			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
//			res := httptest.NewRecorder()
//
//			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
//			// directly and pass in our Request and ResponseRecorder.
//			r.ServeHTTP(res, req)
//
//			// Check the status code is what we expect.
//			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
//				return
//			}
//
//			body, err := ioutil.ReadAll(res.Body)
//			if !assert.NoError(t, err) {
//				return
//			}
//
//			assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")
//		})
//	}
//}
//
//func TestNodeRemoveH(t *testing.T) {
//	v := viper.New()
//	v.SetDefault("storage.driver", "mock")
//
//	stg, _ := storage.Get(v)
//	envs.Get().SetStorage(stg)
//	v.Set("verbose", 0)
//
//	var (
//		err error
//		n1  = getNodeAsset("test1", "", true)
//		n2  = getNodeAsset("test2", "", true)
//	)
//
//	tests := []struct {
//		name         string
//		url          string
//		headers      map[string]string
//		handler      func(http.ResponseWriter, *http.Request)
//		expectedBody string
//		expectedCode int
//	}{
//		{
//			name:         "checking remove failed: not found",
//			url:          fmt.Sprintf("/cluster/node/%s", n2.Meta.Name),
//			handler:      node.NodeRemoveH,
//			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Node not found\"}",
//			expectedCode: http.StatusNotFound,
//		},
//		{
//			name:         "checking remove node successfully",
//			url:          fmt.Sprintf("/cluster/node/%s", n1.Meta.Name),
//			handler:      node.NodeRemoveH,
//			expectedBody: "",
//			expectedCode: http.StatusOK,
//		},
//	}
//
//	for _, tc := range tests {
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Node().Info(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = stg.Put(context.Background(), stg.Collection().Node().Info(), n1.SelfLink().String(), &n1, nil)
//		assert.NoError(t, err)
//
//		t.Run(tc.name, func(t *testing.T) {
//
//			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
//			// pass 'nil' as the third parameter.
//			req, err := http.NewRequest("GET", tc.url, nil)
//			assert.NoError(t, err)
//
//			if tc.headers != nil {
//				for key, val := range tc.headers {
//					req.Header.Set(key, val)
//				}
//			}
//
//			r := mux.NewRouter()
//			r.HandleFunc("/cluster/node/{node}", tc.handler)
//
//			setRequestVars(r, req)
//
//			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
//			res := httptest.NewRecorder()
//
//			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
//			// directly and pass in our Request and ResponseRecorder.
//			r.ServeHTTP(res, req)
//
//			// Check the status code is what we expect.
//			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")
//
//			body, err := ioutil.ReadAll(res.Body)
//			assert.NoError(t, err)
//
//			assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")
//		})
//	}
//}
//
//func TestNodeSetMetaH(t *testing.T) {
//	v := viper.New()
//	v.SetDefault("storage.driver", "mock")
//
//	stg, _ := storage.Get(v)
//	envs.Get().SetStorage(stg)
//	v.Set("verbose", 0)
//	strPointer := func(s string) *string { return &s }
//
//	var (
//		ctx = context.Background()
//
//		n1 = getNodeAsset("test1", "", true)
//		n2 = getNodeAsset("test2", "", true)
//		uo = v1.Request().Node().UpdateOptions()
//	)
//
//	n1.Meta.Architecture = "test"
//
//	uo.Meta = &models.NodeUpdateMetaOptions{}
//	uo.Meta.Architecture = strPointer("test")
//
//	view, err := v1.View().Node().New(&n1).ToJson()
//	assert.NoError(t, err)
//
//	type args struct {
//		ctx  context.Context
//		node string
//	}
//
//	tests := []struct {
//		name         string
//		args         args
//		headers      map[string]string
//		handler      func(http.ResponseWriter, *http.Request)
//		data         string
//		expectedBody string
//		expectedCode int
//	}{
//		{
//			name:         "checking update node failed: not found",
//			args:         args{ctx, n2.Meta.Name},
//			handler:      node.NodeSetMetaH,
//			data:         uo.ToJson(),
//			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Node not found\"}",
//			expectedCode: http.StatusNotFound,
//		},
//		{
//			name:         "checking update node successfully",
//			args:         args{ctx, n1.Meta.Name},
//			handler:      node.NodeSetMetaH,
//			data:         uo.ToJson(),
//			expectedBody: string(view),
//			expectedCode: http.StatusOK,
//		},
//	}
//
//	for _, tc := range tests {
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Node().Info(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = stg.Put(context.Background(), stg.Collection().Node().Info(), n1.SelfLink().String(), &n1, nil)
//		assert.NoError(t, err)
//
//		t.Run(tc.name, func(t *testing.T) {
//
//			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
//			// pass 'nil' as the third parameter.
//			req, err := http.NewRequest("PUT", fmt.Sprintf("/cluster/node/%s/meta", tc.args.node), strings.NewReader(tc.data))
//			assert.NoError(t, err)
//
//			if tc.headers != nil {
//				for key, val := range tc.headers {
//					req.Header.Set(key, val)
//				}
//			}
//
//			r := mux.NewRouter()
//			r.HandleFunc("/cluster/node/{node}/meta", tc.handler)
//
//			setRequestVars(r, req)
//
//			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
//			res := httptest.NewRecorder()
//
//			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
//			// directly and pass in our Request and ResponseRecorder.
//			r.ServeHTTP(res, req)
//
//			// Check the status code is what we expect.
//			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")
//
//			body, err := ioutil.ReadAll(res.Body)
//			assert.NoError(t, err)
//			assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")
//
//			if tc.expectedCode == http.StatusOK {
//				got := new(models.Node)
//				err = envs.Get().GetStorage().Get(context.Background(), stg.Collection().Node().Info(), tc.args.node, got, nil)
//				assert.NoError(t, err)
//				if !assert.NotNil(t, got, "node should not be empty") {
//					return
//				}
//				assert.Equal(t, *uo.Meta.Architecture, got.Meta.Architecture, "Architecture not equal")
//			}
//
//		})
//	}
//}
//
//func TestNodeConnectH(t *testing.T) {
//	v := viper.New()
//	v.SetDefault("storage.driver", "mock")
//
//	stg, _ := storage.Get(v)
//	envs.Get().SetStorage(stg)
//	v.Set("verbose", 0)
//
//	var (
//		err error
//		ctx = context.Background()
//
//		n1 = getNodeAsset("test1", "", true)
//		n2 = getNodeAsset("test2", "", true)
//		uo = v1.Request().Node().NodeConnectOptions()
//	)
//
//	uo.Info.Hostname = "test2"
//	uo.Info.Architecture = "mac"
//
//	type args struct {
//		ctx  context.Context
//		node string
//	}
//
//	tests := []struct {
//		name         string
//		args         args
//		headers      map[string]string
//		handler      func(http.ResponseWriter, *http.Request)
//		data         string
//		expectedBody string
//		expectedCode int
//	}{
//		{
//			name:         "checking create node successful",
//			args:         args{ctx, n2.Meta.Name},
//			handler:      node.NodeConnectH,
//			data:         uo.ToJson(),
//			expectedBody: "",
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "checking update node successful",
//			args:         args{ctx, n1.Meta.Name},
//			handler:      node.NodeConnectH,
//			data:         uo.ToJson(),
//			expectedBody: "",
//			expectedCode: http.StatusOK,
//		},
//	}
//
//	for _, tc := range tests {
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Node().Info(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = stg.Put(context.Background(), stg.Collection().Node().Info(), n1.SelfLink().String(), &n1, nil)
//		assert.NoError(t, err)
//
//		t.Run(tc.name, func(t *testing.T) {
//
//			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
//			// pass 'nil' as the third parameter.
//			req, err := http.NewRequest("PUT", fmt.Sprintf("/cluster/node/%s", tc.args.node), strings.NewReader(tc.data))
//			assert.NoError(t, err)
//
//			if tc.headers != nil {
//				for key, val := range tc.headers {
//					req.Header.Set(key, val)
//				}
//			}
//
//			r := mux.NewRouter()
//			r.HandleFunc("/cluster/node/{node}", tc.handler)
//
//			setRequestVars(r, req)
//
//			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
//			res := httptest.NewRecorder()
//
//			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
//			// directly and pass in our Request and ResponseRecorder.
//			r.ServeHTTP(res, req)
//
//			// Check the status code is what we expect.
//			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")
//
//			body, err := ioutil.ReadAll(res.Body)
//			assert.NoError(t, err)
//			assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")
//
//			if tc.expectedCode == http.StatusOK {
//				got := new(models.Node)
//				err = envs.Get().GetStorage().Get(context.Background(), stg.Collection().Node().Info(), tc.args.node, got, nil)
//				if assert.NoError(t, err) {
//					assert.Equal(t, uo.Info.Hostname, got.Meta.Hostname, "hostname not equal")
//					assert.Equal(t, uo.Info.Architecture, got.Meta.Architecture, "architecture not equal")
//				}
//
//			}
//
//		})
//	}
//}
//
//func TestNodeSetStatusH(t *testing.T) {
//
//	v := viper.New()
//	v.SetDefault("storage.driver", "mock")
//
//	stg, _ := storage.Get(v)
//	cg := cache.NewCache()
//
//	envs.Get().SetStorage(stg)
//	envs.Get().SetCache(cg)
//
//	v.Set("verbose", 0)
//
//	var (
//		err error
//		ctx = context.Background()
//
//		n1 = getNodeAsset("test1", "", true)
//		n2 = getNodeAsset("test2", "", true)
//		uo = v1.Request().Node().NodeStatusOptions()
//		nm = new(models.NodeManifest)
//	)
//
//	nm.Meta.Initial = true
//	nm.Exporter = new(models.ExporterManifest)
//	nm.Resolvers = make(map[string]*models.ResolverManifest, 0)
//	uo.Resources.Capacity.Pods = 20
//
//	type args struct {
//		ctx  context.Context
//		node string
//	}
//
//	view, err := v1.View().Node().NewManifest(nm).ToJson()
//	assert.NoError(t, err)
//	tests := []struct {
//		name         string
//		args         args
//		headers      map[string]string
//		handler      func(http.ResponseWriter, *http.Request)
//		data         string
//		expectedBody string
//		expectedCode int
//	}{
//		{
//			name:         "checking update node failed: not found",
//			args:         args{ctx, n2.Meta.Name},
//			handler:      node.NodeSetStatusH,
//			data:         uo.ToJson(),
//			expectedBody: "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Node not found\"}",
//			expectedCode: http.StatusNotFound,
//		},
//		{
//			name:         "checking update node successfully",
//			args:         args{ctx, n1.Meta.Name},
//			handler:      node.NodeSetStatusH,
//			data:         uo.ToJson(),
//			expectedBody: string(view),
//			expectedCode: http.StatusOK,
//		},
//	}
//
//	for _, tc := range tests {
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Node().Info(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = stg.Put(context.Background(), stg.Collection().Node().Info(), n1.SelfLink().String(), &n1, nil)
//		assert.NoError(t, err)
//
//		t.Run(tc.name, func(t *testing.T) {
//
//			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
//			// pass 'nil' as the third parameter.
//			req, err := http.NewRequest("PUT", fmt.Sprintf("/cluster/node/%s/status", tc.args.node), strings.NewReader(tc.data))
//			assert.NoError(t, err)
//
//			if tc.headers != nil {
//				for key, val := range tc.headers {
//					req.Header.Set(key, val)
//				}
//			}
//
//			r := mux.NewRouter()
//			r.HandleFunc("/cluster/node/{node}/status", tc.handler)
//
//			setRequestVars(r, req)
//
//			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
//			res := httptest.NewRecorder()
//
//			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
//			// directly and pass in our Request and ResponseRecorder.
//			r.ServeHTTP(res, req)
//
//			// Check the status code is what we expect.
//			assert.Equal(t, tc.expectedCode, res.Code, "status code not equal")
//
//			body, err := ioutil.ReadAll(res.Body)
//			assert.NoError(t, err)
//			assert.Equal(t, tc.expectedBody, string(body), "incorrect status code")
//
//			if tc.expectedCode == http.StatusOK {
//				got := new(models.Node)
//				err = envs.Get().GetStorage().Get(context.Background(), stg.Collection().Node().Info(), tc.args.node, got, nil)
//				assert.NoError(t, err)
//				assert.Equal(t, uo.Resources.Capacity.Pods, got.Status.Capacity.Pods, "pods not equal")
//			}
//
//		})
//	}
//}
//
//func setRequestVars(r *mux.Router, req *http.Request) {
//	var match mux.RouteMatch
//	// Take the request and match it
//	r.Match(req, &match)
//	// Push the variable onto the context
//	req = mux.SetURLVars(req, match.Vars)
//}
//
//func getNodeAsset(name, desc string, online bool) models.Node {
//	var n = models.Node{
//		Meta: models.NodeMeta{},
//		Status: models.NodeStatus{
//			Online: true,
//			Capacity: models.NodeResources{
//				Containers: 2,
//				Pods:       2,
//				RAM:        1024,
//				CPU:        2,
//				Storage:    512,
//			},
//			Allocated: models.NodeResources{
//				Containers: 1,
//				Pods:       1,
//				RAM:        512,
//				CPU:        1,
//				Storage:    256,
//			},
//		},
//		Spec: models.NodeSpec{},
//	}
//
//	n.Meta.Name = name
//	n.Meta.Description = desc
//	n.Meta.Hostname = name
//	n.Meta.SelfLink = *models.NewNodeSelfLink(n.Meta.Hostname)
//
//	return n
//}
//
//func getPodManifest() *models.PodManifest {
//	p := models.PodManifest{}
//	return &p
//}
