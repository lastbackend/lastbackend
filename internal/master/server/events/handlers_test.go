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

package events_test

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/lastbackend/lastbackend/internal/api/envs"
	"github.com/lastbackend/lastbackend/internal/master/http/events"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/internal/util/http/middleware"
	"github.com/lastbackend/lastbackend/internal/util/resource"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Testing NamespaceInfoH handler
func TestEventsSubscribe(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")
	v.SetDefault("token", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	go func() {

	}()

	var (
		ns0 = getNamespaceAsset("ns0", "demo")
		ns1 = getNamespaceAsset("ns1", "demo")

		sc0 = getServiceAsset("ns", "svc0", "desc")
		sc1 = getServiceAsset("ns", "svc1", "desc")

		dp0 = getDeploymentAsset("ns", "svc", "dp0")
		dp1 = getDeploymentAsset("ns", "svc", "dp1")

		pd0 = getPodAsset("ns", "svc", "dp", "pod0", "")
		pd1 = getPodAsset("ns", "svc", "dp", "pod1", "")

		rt0 = getRouteAsset("ns", "route0")
		rt1 = getRouteAsset("ns", "route1")

		sk0 = getSecretAsset("ns", "secret0")
		sk1 = getSecretAsset("ns", "secret1")

		cf0 = getConfigAsset("ns", "config0")
		cf1 = getConfigAsset("ns", "config1")

		vl0 = getVolumeAsset("ns", "volume0")
		vl1 = getVolumeAsset("ns", "volume1")

		nd0 = getNodeAsset("initial", "desc", true)
		nd1 = getNodeAsset("demo", "desc", true)

		dc0 = getDiscoveryAsset("initial", "desc", true)
		dc1 = getDiscoveryAsset("demo", "desc", true)

		ig0 = getIngressAsset("initial", "desc", true)
		ig1 = getIngressAsset("demo", "desc", true)
	)

	_, vns1 := getEventAsset(types.EventActionCreate, types.KindNamespace, ns1)
	_, vns2 := getEventAsset(types.EventActionUpdate, types.KindNamespace, ns0)
	_, vns3 := getEventAsset(types.EventActionDelete, types.KindNamespace, ns0)

	_, vsc1 := getEventAsset(types.EventActionCreate, types.KindService, sc1)
	_, vsc2 := getEventAsset(types.EventActionUpdate, types.KindService, sc0)
	_, vsc3 := getEventAsset(types.EventActionDelete, types.KindService, sc0)

	_, vdp1 := getEventAsset(types.EventActionCreate, types.KindDeployment, dp1)
	_, vdp2 := getEventAsset(types.EventActionUpdate, types.KindDeployment, dp0)
	_, vdp3 := getEventAsset(types.EventActionDelete, types.KindDeployment, dp0)

	_, vpd1 := getEventAsset(types.EventActionCreate, types.KindPod, pd1)
	_, vpd2 := getEventAsset(types.EventActionUpdate, types.KindPod, pd0)
	_, vpd3 := getEventAsset(types.EventActionDelete, types.KindPod, pd0)

	_, vrt1 := getEventAsset(types.EventActionCreate, types.KindRoute, rt1)
	_, vrt2 := getEventAsset(types.EventActionUpdate, types.KindRoute, rt0)
	_, vrt3 := getEventAsset(types.EventActionDelete, types.KindRoute, rt0)

	_, vsk1 := getEventAsset(types.EventActionCreate, types.KindSecret, sk1)
	_, vsk2 := getEventAsset(types.EventActionUpdate, types.KindSecret, sk0)
	_, vsk3 := getEventAsset(types.EventActionDelete, types.KindSecret, sk0)

	_, vcf1 := getEventAsset(types.EventActionCreate, types.KindConfig, cf1)
	_, vcf2 := getEventAsset(types.EventActionUpdate, types.KindConfig, cf0)
	_, vcf3 := getEventAsset(types.EventActionDelete, types.KindConfig, cf0)

	_, vvl1 := getEventAsset(types.EventActionCreate, types.KindVolume, vl1)
	_, vvl2 := getEventAsset(types.EventActionUpdate, types.KindVolume, vl0)
	_, vvl3 := getEventAsset(types.EventActionDelete, types.KindVolume, vl0)

	_, vnd1 := getEventAsset(types.EventActionCreate, types.KindNode, nd1)
	_, vnd2 := getEventAsset(types.EventActionUpdate, types.KindNode, nd0)
	_, vnd3 := getEventAsset(types.EventActionDelete, types.KindNode, nd0)

	_, vdc1 := getEventAsset(types.EventActionCreate, types.KindDiscovery, dc1)
	_, vdc2 := getEventAsset(types.EventActionUpdate, types.KindDiscovery, dc0)
	_, vdc3 := getEventAsset(types.EventActionDelete, types.KindDiscovery, dc0)

	_, vig1 := getEventAsset(types.EventActionCreate, types.KindIngress, ig1)
	_, vig2 := getEventAsset(types.EventActionUpdate, types.KindIngress, ig0)
	_, vig3 := getEventAsset(types.EventActionDelete, types.KindIngress, ig0)

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx    context.Context
		token  string
		kind   string
		action string
		obj    interface{}
	}

	var token = v.GetString("token")

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
			name:         "access check",
			args:         args{ctx, "", types.KindNamespace, types.EventActionCreate, ns1},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			err:          "{\"code\":401,\"status\":\"Not Authorized\",\"message\":\"Access denied\"}",
			want:         string(vns1),
			wantErr:      true,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "namespace create",
			args:         args{ctx, token, types.KindNamespace, types.EventActionCreate, ns1},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vns1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "namespace update",
			args:         args{ctx, token, types.KindNamespace, types.EventActionUpdate, ns0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vns2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "namespace remove",
			args:         args{ctx, token, types.KindNamespace, types.EventActionDelete, ns0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vns3),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},

		{
			name:         "service create",
			args:         args{ctx, token, types.KindService, types.EventActionCreate, sc1},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vsc1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "service update",
			args:         args{ctx, token, types.KindService, types.EventActionUpdate, sc0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vsc2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "service remove",
			args:         args{ctx, token, types.KindService, types.EventActionDelete, sc0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vsc3),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},

		{
			name:         "deployment create",
			args:         args{ctx, token, types.KindDeployment, types.EventActionCreate, dp1},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vdp1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "deployment update",
			args:         args{ctx, token, types.KindDeployment, types.EventActionUpdate, dp0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vdp2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "deployment remove",
			args:         args{ctx, token, types.KindDeployment, types.EventActionDelete, dp0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vdp3),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},

		{
			name:         "pod create",
			args:         args{ctx, token, types.KindPod, types.EventActionCreate, pd1},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vpd1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "pod update",
			args:         args{ctx, token, types.KindPod, types.EventActionUpdate, pd0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vpd2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "pod remove",
			args:         args{ctx, token, types.KindPod, types.EventActionDelete, pd0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vpd3),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},

		{
			name:         "route create",
			args:         args{ctx, token, types.KindRoute, types.EventActionCreate, rt1},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vrt1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "route update",
			args:         args{ctx, token, types.KindRoute, types.EventActionUpdate, rt0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vrt2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "route remove",
			args:         args{ctx, token, types.KindRoute, types.EventActionDelete, rt0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vrt3),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},

		{
			name:         "secret create",
			args:         args{ctx, token, types.KindSecret, types.EventActionCreate, sk1},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vsk1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "secret update",
			args:         args{ctx, token, types.KindSecret, types.EventActionUpdate, sk0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vsk2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "secret remove",
			args:         args{ctx, token, types.KindSecret, types.EventActionDelete, sk0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vsk3),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},

		{
			name:         "config create",
			args:         args{ctx, token, types.KindConfig, types.EventActionCreate, cf1},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vcf1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "config update",
			args:         args{ctx, token, types.KindConfig, types.EventActionUpdate, cf0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vcf2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "config remove",
			args:         args{ctx, token, types.KindConfig, types.EventActionDelete, cf0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vcf3),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},

		{
			name:         "volume create",
			args:         args{ctx, token, types.KindVolume, types.EventActionCreate, vl1},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vvl1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "volume update",
			args:         args{ctx, token, types.KindVolume, types.EventActionUpdate, vl0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vvl2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "volume remove",
			args:         args{ctx, token, types.KindVolume, types.EventActionDelete, vl0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vvl3),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},

		{
			name:         "node create",
			args:         args{ctx, token, types.KindNode, types.EventActionCreate, nd1},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vnd1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "node update",
			args:         args{ctx, token, types.KindNode, types.EventActionUpdate, nd0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vnd2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "node remove",
			args:         args{ctx, token, types.KindNode, types.EventActionDelete, nd0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vnd3),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},

		{
			name:         "discovery create",
			args:         args{ctx, token, types.KindDiscovery, types.EventActionCreate, dc1},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vdc1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "discovery update",
			args:         args{ctx, token, types.KindDiscovery, types.EventActionUpdate, dc0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vdc2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "discovery remove",
			args:         args{ctx, token, types.KindDiscovery, types.EventActionDelete, dc0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vdc3),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},

		{
			name:         "ingress create",
			args:         args{ctx, token, types.KindIngress, types.EventActionCreate, ig1},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vig1),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "ingress update",
			args:         args{ctx, token, types.KindIngress, types.EventActionUpdate, ig0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vig2),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "ingress remove",
			args:         args{ctx, token, types.KindIngress, types.EventActionDelete, ig0},
			fields:       fields{stg},
			handler:      events.EventSubscribeH,
			want:         string(vig3),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Service(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Deployment(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Pod(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Route(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Secret(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Config(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Volume(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Node().Info(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Node().Status(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Discovery().Info(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Discovery().Status(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Ingress().Info(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Ingress().Status(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			var done = make(chan bool)

			clear()
			defer clear()

			switch tc.args.kind {
			case types.KindNamespace:
				err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns0.SelfLink().String(), ns0, nil)
				if !assert.NoError(t, err, "initial namespace insert error") {
					return
				}
			case types.KindService:
				err := tc.fields.stg.Put(context.Background(), stg.Collection().Service(), sc0.SelfLink().String(), sc0, nil)
				if !assert.NoError(t, err, "initial service insert error") {
					return
				}
			case types.KindDeployment:
				err := tc.fields.stg.Put(context.Background(), stg.Collection().Deployment(), dp0.SelfLink().String(), dp0, nil)
				if !assert.NoError(t, err, "initial service insert error") {
					return
				}

			case types.KindPod:
				err := tc.fields.stg.Put(context.Background(), stg.Collection().Pod(), pd0.SelfLink().String(), pd0, nil)
				if !assert.NoError(t, err, "initial service insert error") {
					return
				}

			case types.KindRoute:
				err := tc.fields.stg.Put(context.Background(), stg.Collection().Route(), rt0.SelfLink().String(), rt0, nil)
				if !assert.NoError(t, err, "initial service insert error") {
					return
				}

			case types.KindSecret:
				err := tc.fields.stg.Put(context.Background(), stg.Collection().Secret(), sk0.SelfLink().String(), sk0, nil)
				if !assert.NoError(t, err, "initial service insert error") {
					return
				}

			case types.KindConfig:
				err := tc.fields.stg.Put(context.Background(), stg.Collection().Config(), cf0.SelfLink().String(), cf0, nil)
				if !assert.NoError(t, err, "initial service insert error") {
					return
				}

			case types.KindVolume:
				err := tc.fields.stg.Put(context.Background(), stg.Collection().Volume(), vl0.SelfLink().String(), vl0, nil)
				if !assert.NoError(t, err, "initial service insert error") {
					return
				}

			case types.KindNode:
				err := tc.fields.stg.Put(context.Background(), stg.Collection().Node().Status(), nd0.SelfLink().String(), nd0, nil)
				if !assert.NoError(t, err, "initial node status insert error") {
					return
				}

				err = tc.fields.stg.Put(context.Background(), stg.Collection().Node().Info(), nd0.SelfLink().String(), nd0, nil)
				if !assert.NoError(t, err, "initial node info insert error") {
					return
				}

			case types.KindDiscovery:
				err := tc.fields.stg.Put(context.Background(), stg.Collection().Discovery().Status(), dc0.SelfLink().String(), dc0, nil)
				if !assert.NoError(t, err, "initial node status insert error") {
					return
				}

				err = tc.fields.stg.Put(context.Background(), stg.Collection().Discovery().Info(), dc0.SelfLink().String(), dc0, nil)
				if !assert.NoError(t, err, "initial node info insert error") {
					return
				}

			case types.KindIngress:
				err := tc.fields.stg.Put(context.Background(), stg.Collection().Ingress().Status(), ig0.SelfLink().String(), ig0, nil)
				if !assert.NoError(t, err, "initial node status insert error") {
					return
				}

				err = tc.fields.stg.Put(context.Background(), stg.Collection().Ingress().Info(), ig0.SelfLink().String(), ig0, nil)
				if !assert.NoError(t, err, "initial node info insert error") {
					return
				}

			}

			<-time.NewTimer(50 * time.Millisecond).C

			// Create test server with the echo handler.
			s := httptest.NewServer(middleware.Authenticate(context.Background(), events.EventSubscribeH))
			defer s.Close()

			// Convert http://127.0.0.1 to ws://127.0.0.
			u := "ws" + strings.TrimPrefix(s.URL, "http") + "?x-lastbackend-token=" + tc.args.token

			// Connect to the server
			ws, _, err := websocket.DefaultDialer.Dial(u, nil)
			if err != nil {
				if tc.wantErr && tc.expectedCode == http.StatusUnauthorized {
					return
				}
				t.Fatalf("%v", err)
			}
			defer func() {
				_ = ws.Close()
			}()

			timer := time.NewTimer(50 * time.Millisecond)

			var cl = false

			go func() {
				for {
					select {
					case <-timer.C:
						t.Error("payload not received")
						cl = true

						err = ws.Close()
						assert.NoError(t, err, "websocket close err")

						done <- true
						return
					}
				}
			}()

			go func() {
				_, p, err := ws.ReadMessage()

				if err != nil {
					if !cl {
						if !assert.NoError(t, err, "websocket read err") {
							err = ws.Close()
							assert.NoError(t, err, "websocket close err")
						}
					}
					done <- true
					return
				}

				err = ws.Close()
				assert.NoError(t, err, "websocket close err")

				assert.Equal(t, tc.want, string(p), "event payload differs")
				done <- true
			}()

			<-time.NewTimer(10 * time.Millisecond).C
			switch tc.args.kind {
			case types.KindNamespace:

				item := tc.args.obj.(*types.Namespace)
				switch tc.args.action {
				case types.EventActionCreate:
					err = tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), item.SelfLink().String(), item, nil)
				case types.EventActionUpdate:
					err = tc.fields.stg.Set(context.Background(), stg.Collection().Namespace(), item.SelfLink().String(), item, nil)
				case types.EventActionDelete:
					err = tc.fields.stg.Del(context.Background(), stg.Collection().Namespace(), item.SelfLink().String())
				}

			case types.KindService:
				item := tc.args.obj.(*types.Service)
				switch tc.args.action {
				case types.EventActionCreate:
					err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), item.SelfLink().String(), item, nil)
				case types.EventActionUpdate:
					err = tc.fields.stg.Set(context.Background(), stg.Collection().Service(), item.SelfLink().String(), item, nil)
				case types.EventActionDelete:
					err = tc.fields.stg.Del(context.Background(), stg.Collection().Service(), item.SelfLink().String())
				}

			case types.KindDeployment:
				item := tc.args.obj.(*types.Deployment)
				switch tc.args.action {
				case types.EventActionCreate:
					err = tc.fields.stg.Put(context.Background(), stg.Collection().Deployment(), item.SelfLink().String(), item, nil)
				case types.EventActionUpdate:
					err = tc.fields.stg.Set(context.Background(), stg.Collection().Deployment(), item.SelfLink().String(), item, nil)
				case types.EventActionDelete:
					err = tc.fields.stg.Del(context.Background(), stg.Collection().Deployment(), item.SelfLink().String())
				}

			case types.KindPod:
				item := tc.args.obj.(*types.Pod)
				switch tc.args.action {
				case types.EventActionCreate:
					err = tc.fields.stg.Put(context.Background(), stg.Collection().Pod(), item.SelfLink().String(), item, nil)
				case types.EventActionUpdate:
					err = tc.fields.stg.Set(context.Background(), stg.Collection().Pod(), item.SelfLink().String(), item, nil)
				case types.EventActionDelete:
					err = tc.fields.stg.Del(context.Background(), stg.Collection().Pod(), item.SelfLink().String())
				}

			case types.KindVolume:
				item := tc.args.obj.(*types.Volume)
				switch tc.args.action {
				case types.EventActionCreate:
					err = tc.fields.stg.Put(context.Background(), stg.Collection().Volume(), item.SelfLink().String(), item, nil)
				case types.EventActionUpdate:
					err = tc.fields.stg.Set(context.Background(), stg.Collection().Volume(), item.SelfLink().String(), item, nil)
				case types.EventActionDelete:
					err = tc.fields.stg.Del(context.Background(), stg.Collection().Volume(), item.SelfLink().String())
				}

			case types.KindSecret:
				item := tc.args.obj.(*types.Secret)
				switch tc.args.action {
				case types.EventActionCreate:
					err = tc.fields.stg.Put(context.Background(), stg.Collection().Secret(), item.SelfLink().String(), item, nil)
				case types.EventActionUpdate:
					err = tc.fields.stg.Set(context.Background(), stg.Collection().Secret(), item.SelfLink().String(), item, nil)
				case types.EventActionDelete:
					err = tc.fields.stg.Del(context.Background(), stg.Collection().Secret(), item.SelfLink().String())
				}

			case types.KindConfig:
				item := tc.args.obj.(*types.Config)
				switch tc.args.action {
				case types.EventActionCreate:
					err = tc.fields.stg.Put(context.Background(), stg.Collection().Config(), item.SelfLink().String(), item, nil)
				case types.EventActionUpdate:
					err = tc.fields.stg.Set(context.Background(), stg.Collection().Config(), item.SelfLink().String(), item, nil)
				case types.EventActionDelete:
					err = tc.fields.stg.Del(context.Background(), stg.Collection().Config(), item.SelfLink().String())
				}

			case types.KindRoute:
				item := tc.args.obj.(*types.Route)
				switch tc.args.action {
				case types.EventActionCreate:
					err = tc.fields.stg.Put(context.Background(), stg.Collection().Route(), item.SelfLink().String(), item, nil)
				case types.EventActionUpdate:
					err = tc.fields.stg.Set(context.Background(), stg.Collection().Route(), item.SelfLink().String(), item, nil)
				case types.EventActionDelete:
					err = tc.fields.stg.Del(context.Background(), stg.Collection().Route(), item.SelfLink().String())
				}

			case types.KindNode:
				item := tc.args.obj.(*types.Node)
				switch tc.args.action {
				case types.EventActionCreate:
					err = tc.fields.stg.Put(context.Background(), stg.Collection().Node().Status(), item.SelfLink().String(), item.Status, nil)
					<-time.NewTimer(10 * time.Millisecond).C
					err = tc.fields.stg.Put(context.Background(), stg.Collection().Node().Info(), item.SelfLink().String(), item, nil)
				case types.EventActionUpdate:
					err = tc.fields.stg.Set(context.Background(), stg.Collection().Node().Status(), item.SelfLink().String(), item.Status, nil)
					<-time.NewTimer(10 * time.Millisecond).C
					err = tc.fields.stg.Set(context.Background(), stg.Collection().Node().Info(), item.SelfLink().String(), item, nil)
				case types.EventActionDelete:
					err = tc.fields.stg.Del(context.Background(), stg.Collection().Node().Status(), item.SelfLink().String())
					<-time.NewTimer(10 * time.Millisecond).C
					err = tc.fields.stg.Del(context.Background(), stg.Collection().Node().Info(), item.SelfLink().String())
				}

			case types.KindDiscovery:
				item := tc.args.obj.(*types.Discovery)
				switch tc.args.action {
				case types.EventActionCreate:
					err = tc.fields.stg.Put(context.Background(), stg.Collection().Discovery().Status(), item.SelfLink().String(), item.Status, nil)
					<-time.NewTimer(10 * time.Millisecond).C
					err = tc.fields.stg.Put(context.Background(), stg.Collection().Discovery().Info(), item.SelfLink().String(), item, nil)
				case types.EventActionUpdate:
					err = tc.fields.stg.Set(context.Background(), stg.Collection().Discovery().Status(), item.SelfLink().String(), item.Status, nil)
					<-time.NewTimer(10 * time.Millisecond).C
					err = tc.fields.stg.Set(context.Background(), stg.Collection().Discovery().Info(), item.SelfLink().String(), item, nil)
				case types.EventActionDelete:
					err = tc.fields.stg.Del(context.Background(), stg.Collection().Discovery().Status(), item.SelfLink().String())
					<-time.NewTimer(10 * time.Millisecond).C
					err = tc.fields.stg.Del(context.Background(), stg.Collection().Discovery().Info(), item.SelfLink().String())
				}

			case types.KindIngress:
				item := tc.args.obj.(*types.Ingress)
				switch tc.args.action {
				case types.EventActionCreate:
					err = tc.fields.stg.Put(context.Background(), stg.Collection().Ingress().Status(), item.SelfLink().String(), item.Status, nil)
					<-time.NewTimer(10 * time.Millisecond).C
					err = tc.fields.stg.Put(context.Background(), stg.Collection().Ingress().Info(), item.SelfLink().String(), item, nil)
				case types.EventActionUpdate:
					err = tc.fields.stg.Set(context.Background(), stg.Collection().Ingress().Status(), item.SelfLink().String(), item.Status, nil)
					<-time.NewTimer(10 * time.Millisecond).C
					err = tc.fields.stg.Set(context.Background(), stg.Collection().Ingress().Info(), item.SelfLink().String(), item, nil)
				case types.EventActionDelete:
					err = tc.fields.stg.Del(context.Background(), stg.Collection().Ingress().Status(), item.SelfLink().String())
					<-time.NewTimer(10 * time.Millisecond).C
					err = tc.fields.stg.Del(context.Background(), stg.Collection().Ingress().Info(), item.SelfLink().String())
				}

			}

			assert.NoError(t, err, "insert obj error")

			<-done
			timer.Stop()
		})
	}

}

func getNamespaceAsset(name, desc string) *types.Namespace {
	var n = types.Namespace{}

	n.Meta.Name = name
	n.Meta.Description = desc
	n.Meta.SelfLink = *types.NewNamespaceSelfLink(name)
	return &n
}

func getServiceAsset(namespace, name, desc string) *types.Service {
	var s = types.Service{}
	s.Meta.SetDefault()
	s.Meta.Namespace = namespace
	s.Meta.Name = name
	s.Meta.Description = desc
	s.Spec.Replicas = 1
	s.Spec.Template.Containers = make(types.SpecTemplateContainers, 0)
	s.Spec.Template.Containers = append(s.Spec.Template.Containers, &types.SpecTemplateContainer{
		Name: "demo",
	})
	s.Meta.SelfLink = *types.NewServiceSelfLink(namespace, name)
	return &s
}

func getDeploymentAsset(namespace, service, name string) *types.Deployment {
	var d = types.Deployment{}
	d.Meta.SetDefault()
	d.Meta.Namespace = namespace
	d.Meta.Service = service
	d.Meta.Name = name
	d.Meta.SelfLink = *types.NewDeploymentSelfLink(namespace, service, name)
	return &d
}

func getPodAsset(namespace, service, deployment, name, desc string) *types.Pod {
	p := types.Pod{}

	p.Meta.Name = name
	p.Meta.Description = desc
	p.Meta.Namespace = namespace
	psl, _ := types.NewPodSelfLink(types.KindDeployment, types.NewDeploymentSelfLink(namespace, service, deployment).String(), name)
	p.Meta.SelfLink = *psl

	return &p
}

func getSecretAsset(namespace, name string) *types.Secret {
	var s = types.Secret{}
	s.Meta.SetDefault()
	s.Meta.Name = name
	s.Meta.Namespace = namespace
	s.SelfLink()

	s.Spec.Type = types.KindSecretOpaque
	s.Spec.Data = make(map[string][]byte, 0)
	s.Meta.SelfLink = *types.NewSecretSelfLink(namespace, name)
	return &s
}

func getConfigAsset(namespace, name string) *types.Config {
	var c = types.Config{}
	c.Meta.SetDefault()
	c.Meta.Name = name
	c.Meta.Namespace = namespace
	c.Spec.Data = make(map[string]string, 0)
	c.Meta.SelfLink = *types.NewConfigSelfLink(namespace, name)
	return &c
}

func getVolumeAsset(namespace, name string) *types.Volume {
	var r = types.Volume{}
	r.Meta.SetDefault()
	r.Meta.Namespace = namespace
	r.Meta.Name = name
	r.Spec.Selector.Node = ""
	r.Spec.HostPath = "/"
	r.Spec.Capacity.Storage, _ = resource.DecodeMemoryResource("128MB")
	r.Meta.SelfLink = *types.NewVolumeSelfLink(namespace, name)
	return &r
}

func getRouteAsset(namespace, name string) *types.Route {
	var r = types.Route{}
	r.Meta.SetDefault()
	r.Meta.Namespace = namespace
	r.Meta.Name = name
	r.Spec.Endpoint = fmt.Sprintf("%s.test-domain.com", name)
	r.Spec.Rules = make([]types.RouteRule, 0)
	r.Meta.SelfLink = *types.NewRouteSelfLink(namespace, name)
	return &r
}

func getNodeAsset(name, desc string, online bool) *types.Node {
	var n = types.Node{
		Meta: types.NodeMeta{},
		Status: types.NodeStatus{
			Online: online,
			Capacity: types.NodeResources{
				Containers: 2,
				Pods:       2,
				RAM:        1024,
				CPU:        2,
				Storage:    512,
			},
			Allocated: types.NodeResources{
				Containers: 1,
				Pods:       1,
				RAM:        512,
				CPU:        1,
				Storage:    256,
			},
		},
		Spec: types.NodeSpec{},
	}

	n.Meta.Name = name
	n.Meta.Description = desc
	n.Meta.Hostname = name
	n.Meta.SetDefault()
	n.Meta.SelfLink = *types.NewNodeSelfLink(n.Meta.Hostname)

	return &n
}

func getDiscoveryAsset(name, desc string, online bool) *types.Discovery {
	var n = types.Discovery{
		Meta: types.DiscoveryMeta{},
		Status: types.DiscoveryStatus{
			Online: online,
		},
		Spec: types.DiscoverySpec{},
	}

	n.Meta.Name = name
	n.Meta.Description = desc
	n.Meta.SetDefault()
	n.Meta.SelfLink = *types.NewDiscoverySelfLink(n.Meta.Name)

	return &n
}

func getIngressAsset(name, desc string, online bool) *types.Ingress {
	var n = types.Ingress{
		Meta: types.IngressMeta{},
		Status: types.IngressStatus{
			Online: online,
		},
		Spec: types.IngressSpec{},
	}

	n.Meta.Name = name
	n.Meta.Description = desc
	n.Meta.SetDefault()
	n.Meta.SelfLink = *types.NewIngressSelfLink(n.Meta.Name)

	return &n
}

func getEventAsset(action, kind string, obj interface{}) (*types.Event, []byte) {

	e1 := types.Event{Kind: kind}
	switch kind {
	case types.KindNamespace:
		e1.Data = obj.(*types.Namespace)

	case types.KindService:
		e1.Data = obj.(*types.Service)

	case types.KindJob:
		e1.Data = obj.(*types.Task)

	case types.KindDeployment:
		e1.Data = obj.(*types.Deployment)

	case types.KindPod:
		e1.Data = obj.(*types.Pod)

	case types.KindVolume:
		e1.Data = obj.(*types.Volume)

	case types.KindSecret:
		e1.Data = obj.(*types.Secret)

	case types.KindConfig:
		e1.Data = obj.(*types.Config)

	case types.KindRoute:
		e1.Data = obj.(*types.Route)

	case types.KindNode:
		e1.Data = obj.(*types.Node)

	case types.KindIngress:
		e1.Data = obj.(*types.Ingress)

	case types.KindDiscovery:
		e1.Data = obj.(*types.Discovery)

	case types.KindCluster:
		e1.Data = obj.(*types.Cluster)
	}

	e1.Action = action

	vns1 := v1.View().Event().New(&e1)
	bns1, _ := vns1.ToJson()
	return &e1, bns1
}
