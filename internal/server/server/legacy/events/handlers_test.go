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
//
//import (
//	"context"
//	"fmt"
//	"github.com/gorilla/websocket"
//	"github.com/lastbackend/lastbackend/internal/api/envs"
//	"github.com/lastbackend/lastbackend/internal/master/http/events"
//	"github.com/lastbackend/lastbackend/internal/pkg/storage"
//	"github.com/lastbackend/lastbackend/internal/pkg/models"
//	"github.com/lastbackend/lastbackend/internal/util/http/middleware"
//	"github.com/lastbackend/lastbackend/internal/util/resource"
//	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
//
//	"github.com/spf13/viper"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//	"time"
//)
//
//// Testing NamespaceInfoH handler
//func TestEventsSubscribe(t *testing.T) {
//
//	var ctx = context.Background()
//
//	v := viper.New()
//	v.SetDefault("storage.driver", "mock")
//	v.SetDefault("token", "mock")
//
//	stg, _ := storage.Get(v)
//	envs.Get().SetStorage(stg)
//
//	go func() {
//
//	}()
//
//	var (
//		ns0 = getNamespaceAsset("ns0", "demo")
//		ns1 = getNamespaceAsset("ns1", "demo")
//
//		sc0 = getServiceAsset("ns", "svc0", "desc")
//		sc1 = getServiceAsset("ns", "svc1", "desc")
//
//		dp0 = getDeploymentAsset("ns", "svc", "dp0")
//		dp1 = getDeploymentAsset("ns", "svc", "dp1")
//
//		pd0 = getPodAsset("ns", "svc", "dp", "pod0", "")
//		pd1 = getPodAsset("ns", "svc", "dp", "pod1", "")
//
//		rt0 = getRouteAsset("ns", "route0")
//		rt1 = getRouteAsset("ns", "route1")
//
//		sk0 = getSecretAsset("ns", "secret0")
//		sk1 = getSecretAsset("ns", "secret1")
//
//		cf0 = getConfigAsset("ns", "config0")
//		cf1 = getConfigAsset("ns", "config1")
//
//		vl0 = getVolumeAsset("ns", "volume0")
//		vl1 = getVolumeAsset("ns", "volume1")
//
//		nd0 = getNodeAsset("initial", "desc", true)
//		nd1 = getNodeAsset("demo", "desc", true)
//
//		dc0 = getDiscoveryAsset("initial", "desc", true)
//		dc1 = getDiscoveryAsset("demo", "desc", true)
//
//		ig0 = getIngressAsset("initial", "desc", true)
//		ig1 = getIngressAsset("demo", "desc", true)
//	)
//
//	_, vns1 := getEventAsset(models.EventActionCreate, models.KindNamespace, ns1)
//	_, vns2 := getEventAsset(models.EventActionUpdate, models.KindNamespace, ns0)
//	_, vns3 := getEventAsset(models.EventActionDelete, models.KindNamespace, ns0)
//
//	_, vsc1 := getEventAsset(models.EventActionCreate, models.KindService, sc1)
//	_, vsc2 := getEventAsset(models.EventActionUpdate, models.KindService, sc0)
//	_, vsc3 := getEventAsset(models.EventActionDelete, models.KindService, sc0)
//
//	_, vdp1 := getEventAsset(models.EventActionCreate, models.KindDeployment, dp1)
//	_, vdp2 := getEventAsset(models.EventActionUpdate, models.KindDeployment, dp0)
//	_, vdp3 := getEventAsset(models.EventActionDelete, models.KindDeployment, dp0)
//
//	_, vpd1 := getEventAsset(models.EventActionCreate, models.KindPod, pd1)
//	_, vpd2 := getEventAsset(models.EventActionUpdate, models.KindPod, pd0)
//	_, vpd3 := getEventAsset(models.EventActionDelete, models.KindPod, pd0)
//
//	_, vrt1 := getEventAsset(models.EventActionCreate, models.KindRoute, rt1)
//	_, vrt2 := getEventAsset(models.EventActionUpdate, models.KindRoute, rt0)
//	_, vrt3 := getEventAsset(models.EventActionDelete, models.KindRoute, rt0)
//
//	_, vsk1 := getEventAsset(models.EventActionCreate, models.KindSecret, sk1)
//	_, vsk2 := getEventAsset(models.EventActionUpdate, models.KindSecret, sk0)
//	_, vsk3 := getEventAsset(models.EventActionDelete, models.KindSecret, sk0)
//
//	_, vcf1 := getEventAsset(models.EventActionCreate, models.KindConfig, cf1)
//	_, vcf2 := getEventAsset(models.EventActionUpdate, models.KindConfig, cf0)
//	_, vcf3 := getEventAsset(models.EventActionDelete, models.KindConfig, cf0)
//
//	_, vvl1 := getEventAsset(models.EventActionCreate, models.KindVolume, vl1)
//	_, vvl2 := getEventAsset(models.EventActionUpdate, models.KindVolume, vl0)
//	_, vvl3 := getEventAsset(models.EventActionDelete, models.KindVolume, vl0)
//
//	_, vnd1 := getEventAsset(models.EventActionCreate, models.KindNode, nd1)
//	_, vnd2 := getEventAsset(models.EventActionUpdate, models.KindNode, nd0)
//	_, vnd3 := getEventAsset(models.EventActionDelete, models.KindNode, nd0)
//
//	_, vdc1 := getEventAsset(models.EventActionCreate, models.KindDiscovery, dc1)
//	_, vdc2 := getEventAsset(models.EventActionUpdate, models.KindDiscovery, dc0)
//	_, vdc3 := getEventAsset(models.EventActionDelete, models.KindDiscovery, dc0)
//
//	_, vig1 := getEventAsset(models.EventActionCreate, models.KindIngress, ig1)
//	_, vig2 := getEventAsset(models.EventActionUpdate, models.KindIngress, ig0)
//	_, vig3 := getEventAsset(models.EventActionDelete, models.KindIngress, ig0)
//
//	type fields struct {
//		stg storage.IStorage
//	}
//
//	type args struct {
//		ctx    context.Context
//		token  string
//		kind   string
//		action string
//		obj    interface{}
//	}
//
//	var token = v.GetString("token")
//
//	tests := []struct {
//		name         string
//		fields       fields
//		args         args
//		headers      map[string]string
//		handler      func(http.ResponseWriter, *http.Request)
//		err          string
//		want         string
//		wantErr      bool
//		expectedCode int
//	}{
//		{
//			name:         "access check",
//			args:         args{ctx, "", models.KindNamespace, models.EventActionCreate, ns1},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			err:          "{\"code\":401,\"status\":\"Not Authorized\",\"message\":\"Access denied\"}",
//			want:         string(vns1),
//			wantErr:      true,
//			expectedCode: http.StatusUnauthorized,
//		},
//		{
//			name:         "namespace create",
//			args:         args{ctx, token, models.KindNamespace, models.EventActionCreate, ns1},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vns1),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "namespace update",
//			args:         args{ctx, token, models.KindNamespace, models.EventActionUpdate, ns0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vns2),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "namespace remove",
//			args:         args{ctx, token, models.KindNamespace, models.EventActionDelete, ns0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vns3),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//
//		{
//			name:         "service create",
//			args:         args{ctx, token, models.KindService, models.EventActionCreate, sc1},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vsc1),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "service update",
//			args:         args{ctx, token, models.KindService, models.EventActionUpdate, sc0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vsc2),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "service remove",
//			args:         args{ctx, token, models.KindService, models.EventActionDelete, sc0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vsc3),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//
//		{
//			name:         "deployment create",
//			args:         args{ctx, token, models.KindDeployment, models.EventActionCreate, dp1},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vdp1),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "deployment update",
//			args:         args{ctx, token, models.KindDeployment, models.EventActionUpdate, dp0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vdp2),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "deployment remove",
//			args:         args{ctx, token, models.KindDeployment, models.EventActionDelete, dp0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vdp3),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//
//		{
//			name:         "pod create",
//			args:         args{ctx, token, models.KindPod, models.EventActionCreate, pd1},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vpd1),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "pod update",
//			args:         args{ctx, token, models.KindPod, models.EventActionUpdate, pd0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vpd2),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "pod remove",
//			args:         args{ctx, token, models.KindPod, models.EventActionDelete, pd0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vpd3),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//
//		{
//			name:         "route create",
//			args:         args{ctx, token, models.KindRoute, models.EventActionCreate, rt1},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vrt1),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "route update",
//			args:         args{ctx, token, models.KindRoute, models.EventActionUpdate, rt0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vrt2),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "route remove",
//			args:         args{ctx, token, models.KindRoute, models.EventActionDelete, rt0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vrt3),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//
//		{
//			name:         "secret create",
//			args:         args{ctx, token, models.KindSecret, models.EventActionCreate, sk1},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vsk1),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "secret update",
//			args:         args{ctx, token, models.KindSecret, models.EventActionUpdate, sk0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vsk2),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "secret remove",
//			args:         args{ctx, token, models.KindSecret, models.EventActionDelete, sk0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vsk3),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//
//		{
//			name:         "config create",
//			args:         args{ctx, token, models.KindConfig, models.EventActionCreate, cf1},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vcf1),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "config update",
//			args:         args{ctx, token, models.KindConfig, models.EventActionUpdate, cf0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vcf2),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "config remove",
//			args:         args{ctx, token, models.KindConfig, models.EventActionDelete, cf0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vcf3),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//
//		{
//			name:         "volume create",
//			args:         args{ctx, token, models.KindVolume, models.EventActionCreate, vl1},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vvl1),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "volume update",
//			args:         args{ctx, token, models.KindVolume, models.EventActionUpdate, vl0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vvl2),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "volume remove",
//			args:         args{ctx, token, models.KindVolume, models.EventActionDelete, vl0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vvl3),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//
//		{
//			name:         "node create",
//			args:         args{ctx, token, models.KindNode, models.EventActionCreate, nd1},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vnd1),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "node update",
//			args:         args{ctx, token, models.KindNode, models.EventActionUpdate, nd0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vnd2),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "node remove",
//			args:         args{ctx, token, models.KindNode, models.EventActionDelete, nd0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vnd3),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//
//		{
//			name:         "discovery create",
//			args:         args{ctx, token, models.KindDiscovery, models.EventActionCreate, dc1},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vdc1),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "discovery update",
//			args:         args{ctx, token, models.KindDiscovery, models.EventActionUpdate, dc0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vdc2),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "discovery remove",
//			args:         args{ctx, token, models.KindDiscovery, models.EventActionDelete, dc0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vdc3),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//
//		{
//			name:         "ingress create",
//			args:         args{ctx, token, models.KindIngress, models.EventActionCreate, ig1},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vig1),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "ingress update",
//			args:         args{ctx, token, models.KindIngress, models.EventActionUpdate, ig0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vig2),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "ingress remove",
//			args:         args{ctx, token, models.KindIngress, models.EventActionDelete, ig0},
//			fields:       fields{stg},
//			handler:      events.EventSubscribeH,
//			want:         string(vig3),
//			wantErr:      false,
//			expectedCode: http.StatusOK,
//		},
//	}
//
//	clear := func() {
//		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Service(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Deployment(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Pod(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Route(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Secret(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Config(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Volume(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Node().Info(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Node().Status(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Discovery().Info(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Discovery().Status(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Ingress().Info(), models.EmptyString)
//		assert.NoError(t, err)
//
//		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Ingress().Status(), models.EmptyString)
//		assert.NoError(t, err)
//	}
//
//	for _, tc := range tests {
//
//		t.Run(tc.name, func(t *testing.T) {
//
//			var done = make(chan bool)
//
//			clear()
//			defer clear()
//
//			switch tc.args.kind {
//			case models.KindNamespace:
//				err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns0.SelfLink().String(), ns0, nil)
//				if !assert.NoError(t, err, "initial namespace insert error") {
//					return
//				}
//			case models.KindService:
//				err := tc.fields.stg.Put(context.Background(), stg.Collection().Service(), sc0.SelfLink().String(), sc0, nil)
//				if !assert.NoError(t, err, "initial service insert error") {
//					return
//				}
//			case models.KindDeployment:
//				err := tc.fields.stg.Put(context.Background(), stg.Collection().Deployment(), dp0.SelfLink().String(), dp0, nil)
//				if !assert.NoError(t, err, "initial service insert error") {
//					return
//				}
//
//			case models.KindPod:
//				err := tc.fields.stg.Put(context.Background(), stg.Collection().Pod(), pd0.SelfLink().String(), pd0, nil)
//				if !assert.NoError(t, err, "initial service insert error") {
//					return
//				}
//
//			case models.KindRoute:
//				err := tc.fields.stg.Put(context.Background(), stg.Collection().Route(), rt0.SelfLink().String(), rt0, nil)
//				if !assert.NoError(t, err, "initial service insert error") {
//					return
//				}
//
//			case models.KindSecret:
//				err := tc.fields.stg.Put(context.Background(), stg.Collection().Secret(), sk0.SelfLink().String(), sk0, nil)
//				if !assert.NoError(t, err, "initial service insert error") {
//					return
//				}
//
//			case models.KindConfig:
//				err := tc.fields.stg.Put(context.Background(), stg.Collection().Config(), cf0.SelfLink().String(), cf0, nil)
//				if !assert.NoError(t, err, "initial service insert error") {
//					return
//				}
//
//			case models.KindVolume:
//				err := tc.fields.stg.Put(context.Background(), stg.Collection().Volume(), vl0.SelfLink().String(), vl0, nil)
//				if !assert.NoError(t, err, "initial service insert error") {
//					return
//				}
//
//			case models.KindNode:
//				err := tc.fields.stg.Put(context.Background(), stg.Collection().Node().Status(), nd0.SelfLink().String(), nd0, nil)
//				if !assert.NoError(t, err, "initial node status insert error") {
//					return
//				}
//
//				err = tc.fields.stg.Put(context.Background(), stg.Collection().Node().Info(), nd0.SelfLink().String(), nd0, nil)
//				if !assert.NoError(t, err, "initial node info insert error") {
//					return
//				}
//
//			case models.KindDiscovery:
//				err := tc.fields.stg.Put(context.Background(), stg.Collection().Discovery().Status(), dc0.SelfLink().String(), dc0, nil)
//				if !assert.NoError(t, err, "initial node status insert error") {
//					return
//				}
//
//				err = tc.fields.stg.Put(context.Background(), stg.Collection().Discovery().Info(), dc0.SelfLink().String(), dc0, nil)
//				if !assert.NoError(t, err, "initial node info insert error") {
//					return
//				}
//
//			case models.KindIngress:
//				err := tc.fields.stg.Put(context.Background(), stg.Collection().Ingress().Status(), ig0.SelfLink().String(), ig0, nil)
//				if !assert.NoError(t, err, "initial node status insert error") {
//					return
//				}
//
//				err = tc.fields.stg.Put(context.Background(), stg.Collection().Ingress().Info(), ig0.SelfLink().String(), ig0, nil)
//				if !assert.NoError(t, err, "initial node info insert error") {
//					return
//				}
//
//			}
//
//			<-time.NewTimer(50 * time.Millisecond).C
//
//			// Create test server with the echo handler.
//			s := httptest.NewServer(middleware.Authenticate(context.Background(), events.EventSubscribeH))
//			defer s.Close()
//
//			// Convert http://127.0.0.1 to ws://127.0.0.
//			u := "ws" + strings.TrimPrefix(s.URL, "http") + "?x-lastbackend-token=" + tc.args.token
//
//			// Connect to the server
//			ws, _, err := websocket.DefaultDialer.Dial(u, nil)
//			if err != nil {
//				if tc.wantErr && tc.expectedCode == http.StatusUnauthorized {
//					return
//				}
//				t.Fatalf("%v", err)
//			}
//			defer func() {
//				_ = ws.Close()
//			}()
//
//			timer := time.NewTimer(50 * time.Millisecond)
//
//			var cl = false
//
//			go func() {
//				for {
//					select {
//					case <-timer.C:
//						t.Error("payload not received")
//						cl = true
//
//						err = ws.Close()
//						assert.NoError(t, err, "websocket close err")
//
//						done <- true
//						return
//					}
//				}
//			}()
//
//			go func() {
//				_, p, err := ws.ReadMessage()
//
//				if err != nil {
//					if !cl {
//						if !assert.NoError(t, err, "websocket read err") {
//							err = ws.Close()
//							assert.NoError(t, err, "websocket close err")
//						}
//					}
//					done <- true
//					return
//				}
//
//				err = ws.Close()
//				assert.NoError(t, err, "websocket close err")
//
//				assert.Equal(t, tc.want, string(p), "event payload differs")
//				done <- true
//			}()
//
//			<-time.NewTimer(10 * time.Millisecond).C
//			switch tc.args.kind {
//			case models.KindNamespace:
//
//				item := tc.args.obj.(*models.Namespace)
//				switch tc.args.action {
//				case models.EventActionCreate:
//					err = tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), item.SelfLink().String(), item, nil)
//				case models.EventActionUpdate:
//					err = tc.fields.stg.Set(context.Background(), stg.Collection().Namespace(), item.SelfLink().String(), item, nil)
//				case models.EventActionDelete:
//					err = tc.fields.stg.Del(context.Background(), stg.Collection().Namespace(), item.SelfLink().String())
//				}
//
//			case models.KindService:
//				item := tc.args.obj.(*models.Service)
//				switch tc.args.action {
//				case models.EventActionCreate:
//					err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), item.SelfLink().String(), item, nil)
//				case models.EventActionUpdate:
//					err = tc.fields.stg.Set(context.Background(), stg.Collection().Service(), item.SelfLink().String(), item, nil)
//				case models.EventActionDelete:
//					err = tc.fields.stg.Del(context.Background(), stg.Collection().Service(), item.SelfLink().String())
//				}
//
//			case models.KindDeployment:
//				item := tc.args.obj.(*models.Deployment)
//				switch tc.args.action {
//				case models.EventActionCreate:
//					err = tc.fields.stg.Put(context.Background(), stg.Collection().Deployment(), item.SelfLink().String(), item, nil)
//				case models.EventActionUpdate:
//					err = tc.fields.stg.Set(context.Background(), stg.Collection().Deployment(), item.SelfLink().String(), item, nil)
//				case models.EventActionDelete:
//					err = tc.fields.stg.Del(context.Background(), stg.Collection().Deployment(), item.SelfLink().String())
//				}
//
//			case models.KindPod:
//				item := tc.args.obj.(*models.Pod)
//				switch tc.args.action {
//				case models.EventActionCreate:
//					err = tc.fields.stg.Put(context.Background(), stg.Collection().Pod(), item.SelfLink().String(), item, nil)
//				case models.EventActionUpdate:
//					err = tc.fields.stg.Set(context.Background(), stg.Collection().Pod(), item.SelfLink().String(), item, nil)
//				case models.EventActionDelete:
//					err = tc.fields.stg.Del(context.Background(), stg.Collection().Pod(), item.SelfLink().String())
//				}
//
//			case models.KindVolume:
//				item := tc.args.obj.(*models.Volume)
//				switch tc.args.action {
//				case models.EventActionCreate:
//					err = tc.fields.stg.Put(context.Background(), stg.Collection().Volume(), item.SelfLink().String(), item, nil)
//				case models.EventActionUpdate:
//					err = tc.fields.stg.Set(context.Background(), stg.Collection().Volume(), item.SelfLink().String(), item, nil)
//				case models.EventActionDelete:
//					err = tc.fields.stg.Del(context.Background(), stg.Collection().Volume(), item.SelfLink().String())
//				}
//
//			case models.KindSecret:
//				item := tc.args.obj.(*models.Secret)
//				switch tc.args.action {
//				case models.EventActionCreate:
//					err = tc.fields.stg.Put(context.Background(), stg.Collection().Secret(), item.SelfLink().String(), item, nil)
//				case models.EventActionUpdate:
//					err = tc.fields.stg.Set(context.Background(), stg.Collection().Secret(), item.SelfLink().String(), item, nil)
//				case models.EventActionDelete:
//					err = tc.fields.stg.Del(context.Background(), stg.Collection().Secret(), item.SelfLink().String())
//				}
//
//			case models.KindConfig:
//				item := tc.args.obj.(*models.Config)
//				switch tc.args.action {
//				case models.EventActionCreate:
//					err = tc.fields.stg.Put(context.Background(), stg.Collection().Config(), item.SelfLink().String(), item, nil)
//				case models.EventActionUpdate:
//					err = tc.fields.stg.Set(context.Background(), stg.Collection().Config(), item.SelfLink().String(), item, nil)
//				case models.EventActionDelete:
//					err = tc.fields.stg.Del(context.Background(), stg.Collection().Config(), item.SelfLink().String())
//				}
//
//			case models.KindRoute:
//				item := tc.args.obj.(*models.Route)
//				switch tc.args.action {
//				case models.EventActionCreate:
//					err = tc.fields.stg.Put(context.Background(), stg.Collection().Route(), item.SelfLink().String(), item, nil)
//				case models.EventActionUpdate:
//					err = tc.fields.stg.Set(context.Background(), stg.Collection().Route(), item.SelfLink().String(), item, nil)
//				case models.EventActionDelete:
//					err = tc.fields.stg.Del(context.Background(), stg.Collection().Route(), item.SelfLink().String())
//				}
//
//			case models.KindNode:
//				item := tc.args.obj.(*models.Node)
//				switch tc.args.action {
//				case models.EventActionCreate:
//					err = tc.fields.stg.Put(context.Background(), stg.Collection().Node().Status(), item.SelfLink().String(), item.Status, nil)
//					<-time.NewTimer(10 * time.Millisecond).C
//					err = tc.fields.stg.Put(context.Background(), stg.Collection().Node().Info(), item.SelfLink().String(), item, nil)
//				case models.EventActionUpdate:
//					err = tc.fields.stg.Set(context.Background(), stg.Collection().Node().Status(), item.SelfLink().String(), item.Status, nil)
//					<-time.NewTimer(10 * time.Millisecond).C
//					err = tc.fields.stg.Set(context.Background(), stg.Collection().Node().Info(), item.SelfLink().String(), item, nil)
//				case models.EventActionDelete:
//					err = tc.fields.stg.Del(context.Background(), stg.Collection().Node().Status(), item.SelfLink().String())
//					<-time.NewTimer(10 * time.Millisecond).C
//					err = tc.fields.stg.Del(context.Background(), stg.Collection().Node().Info(), item.SelfLink().String())
//				}
//
//			case models.KindDiscovery:
//				item := tc.args.obj.(*models.Discovery)
//				switch tc.args.action {
//				case models.EventActionCreate:
//					err = tc.fields.stg.Put(context.Background(), stg.Collection().Discovery().Status(), item.SelfLink().String(), item.Status, nil)
//					<-time.NewTimer(10 * time.Millisecond).C
//					err = tc.fields.stg.Put(context.Background(), stg.Collection().Discovery().Info(), item.SelfLink().String(), item, nil)
//				case models.EventActionUpdate:
//					err = tc.fields.stg.Set(context.Background(), stg.Collection().Discovery().Status(), item.SelfLink().String(), item.Status, nil)
//					<-time.NewTimer(10 * time.Millisecond).C
//					err = tc.fields.stg.Set(context.Background(), stg.Collection().Discovery().Info(), item.SelfLink().String(), item, nil)
//				case models.EventActionDelete:
//					err = tc.fields.stg.Del(context.Background(), stg.Collection().Discovery().Status(), item.SelfLink().String())
//					<-time.NewTimer(10 * time.Millisecond).C
//					err = tc.fields.stg.Del(context.Background(), stg.Collection().Discovery().Info(), item.SelfLink().String())
//				}
//
//			case models.KindIngress:
//				item := tc.args.obj.(*models.Ingress)
//				switch tc.args.action {
//				case models.EventActionCreate:
//					err = tc.fields.stg.Put(context.Background(), stg.Collection().Ingress().Status(), item.SelfLink().String(), item.Status, nil)
//					<-time.NewTimer(10 * time.Millisecond).C
//					err = tc.fields.stg.Put(context.Background(), stg.Collection().Ingress().Info(), item.SelfLink().String(), item, nil)
//				case models.EventActionUpdate:
//					err = tc.fields.stg.Set(context.Background(), stg.Collection().Ingress().Status(), item.SelfLink().String(), item.Status, nil)
//					<-time.NewTimer(10 * time.Millisecond).C
//					err = tc.fields.stg.Set(context.Background(), stg.Collection().Ingress().Info(), item.SelfLink().String(), item, nil)
//				case models.EventActionDelete:
//					err = tc.fields.stg.Del(context.Background(), stg.Collection().Ingress().Status(), item.SelfLink().String())
//					<-time.NewTimer(10 * time.Millisecond).C
//					err = tc.fields.stg.Del(context.Background(), stg.Collection().Ingress().Info(), item.SelfLink().String())
//				}
//
//			}
//
//			assert.NoError(t, err, "insert obj error")
//
//			<-done
//			timer.Stop()
//		})
//	}
//
//}
//
//func getNamespaceAsset(name, desc string) *models.Namespace {
//	var n = models.Namespace{}
//
//	n.Meta.Name = name
//	n.Meta.Description = desc
//	n.Meta.SelfLink = *models.NewNamespaceSelfLink(name)
//	return &n
//}
//
//func getServiceAsset(namespace, name, desc string) *models.Service {
//	var s = models.Service{}
//	s.Meta.SetDefault()
//	s.Meta.Namespace = namespace
//	s.Meta.Name = name
//	s.Meta.Description = desc
//	s.Spec.Replicas = 1
//	s.Spec.Template.Containers = make(models.SpecTemplateContainers, 0)
//	s.Spec.Template.Containers = append(s.Spec.Template.Containers, &models.SpecTemplateContainer{
//		Name: "demo",
//	})
//	s.Meta.SelfLink = *models.NewServiceSelfLink(namespace, name)
//	return &s
//}
//
//func getDeploymentAsset(namespace, service, name string) *models.Deployment {
//	var d = models.Deployment{}
//	d.Meta.SetDefault()
//	d.Meta.Namespace = namespace
//	d.Meta.Service = service
//	d.Meta.Name = name
//	d.Meta.SelfLink = *models.NewDeploymentSelfLink(namespace, service, name)
//	return &d
//}
//
//func getPodAsset(namespace, service, deployment, name, desc string) *models.Pod {
//	p := models.Pod{}
//
//	p.Meta.Name = name
//	p.Meta.Description = desc
//	p.Meta.Namespace = namespace
//	psl, _ := models.NewPodSelfLink(models.KindDeployment, models.NewDeploymentSelfLink(namespace, service, deployment).String(), name)
//	p.Meta.SelfLink = *psl
//
//	return &p
//}
//
//func getSecretAsset(namespace, name string) *models.Secret {
//	var s = models.Secret{}
//	s.Meta.SetDefault()
//	s.Meta.Name = name
//	s.Meta.Namespace = namespace
//	s.SelfLink()
//
//	s.Spec.Type = models.KindSecretOpaque
//	s.Spec.Data = make(map[string][]byte, 0)
//	s.Meta.SelfLink = *models.NewSecretSelfLink(namespace, name)
//	return &s
//}
//
//func getConfigAsset(namespace, name string) *models.Config {
//	var c = models.Config{}
//	c.Meta.SetDefault()
//	c.Meta.Name = name
//	c.Meta.Namespace = namespace
//	c.Spec.Data = make(map[string]string, 0)
//	c.Meta.SelfLink = *models.NewConfigSelfLink(namespace, name)
//	return &c
//}
//
//func getVolumeAsset(namespace, name string) *models.Volume {
//	var r = models.Volume{}
//	r.Meta.SetDefault()
//	r.Meta.Namespace = namespace
//	r.Meta.Name = name
//	r.Spec.Selector.Node = ""
//	r.Spec.HostPath = "/"
//	r.Spec.Capacity.Storage, _ = resource.DecodeMemoryResource("128MB")
//	r.Meta.SelfLink = *models.NewVolumeSelfLink(namespace, name)
//	return &r
//}
//
//func getRouteAsset(namespace, name string) *models.Route {
//	var r = models.Route{}
//	r.Meta.SetDefault()
//	r.Meta.Namespace = namespace
//	r.Meta.Name = name
//	r.Spec.Endpoint = fmt.Sprintf("%s.test-domain.com", name)
//	r.Spec.Rules = make([]models.RouteRule, 0)
//	r.Meta.SelfLink = *models.NewRouteSelfLink(namespace, name)
//	return &r
//}
//
//func getNodeAsset(name, desc string, online bool) *models.Node {
//	var n = models.Node{
//		Meta: models.NodeMeta{},
//		Status: models.NodeStatus{
//			Online: online,
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
//	n.Meta.SetDefault()
//	n.Meta.SelfLink = *models.NewNodeSelfLink(n.Meta.Hostname)
//
//	return &n
//}
//
//func getDiscoveryAsset(name, desc string, online bool) *models.Discovery {
//	var n = models.Discovery{
//		Meta: models.DiscoveryMeta{},
//		Status: models.DiscoveryStatus{
//			Online: online,
//		},
//		Spec: models.DiscoverySpec{},
//	}
//
//	n.Meta.Name = name
//	n.Meta.Description = desc
//	n.Meta.SetDefault()
//	n.Meta.SelfLink = *models.NewDiscoverySelfLink(n.Meta.Name)
//
//	return &n
//}
//
//func getIngressAsset(name, desc string, online bool) *models.Ingress {
//	var n = models.Ingress{
//		Meta: models.IngressMeta{},
//		Status: models.IngressStatus{
//			Online: online,
//		},
//		Spec: models.IngressSpec{},
//	}
//
//	n.Meta.Name = name
//	n.Meta.Description = desc
//	n.Meta.SetDefault()
//	n.Meta.SelfLink = *models.NewIngressSelfLink(n.Meta.Name)
//
//	return &n
//}
//
//func getEventAsset(action, kind string, obj interface{}) (*models.Event, []byte) {
//
//	e1 := models.Event{Kind: kind}
//	switch kind {
//	case models.KindNamespace:
//		e1.Data = obj.(*models.Namespace)
//
//	case models.KindService:
//		e1.Data = obj.(*models.Service)
//
//	case models.KindJob:
//		e1.Data = obj.(*models.Task)
//
//	case models.KindDeployment:
//		e1.Data = obj.(*models.Deployment)
//
//	case models.KindPod:
//		e1.Data = obj.(*models.Pod)
//
//	case models.KindVolume:
//		e1.Data = obj.(*models.Volume)
//
//	case models.KindSecret:
//		e1.Data = obj.(*models.Secret)
//
//	case models.KindConfig:
//		e1.Data = obj.(*models.Config)
//
//	case models.KindRoute:
//		e1.Data = obj.(*models.Route)
//
//	case models.KindNode:
//		e1.Data = obj.(*models.Node)
//
//	case models.KindIngress:
//		e1.Data = obj.(*models.Ingress)
//
//	case models.KindDiscovery:
//		e1.Data = obj.(*models.Discovery)
//
//	case models.KindCluster:
//		e1.Data = obj.(*models.Cluster)
//	}
//
//	e1.Action = action
//
//	vns1 := v1.View().Event().New(&e1)
//	bns1, _ := vns1.ToJson()
//	return &e1, bns1
//}
