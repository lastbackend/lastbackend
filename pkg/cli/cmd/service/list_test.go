//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package service_test

//
//import (
//	"encoding/json"
//	n "github.com/lastbackend/lastbackend/pkg/api/namespace/views/v1"
//	sm "github.com/lastbackend/lastbackend/pkg/api/service/views/v1"
//	"github.com/lastbackend/lastbackend/pkg/cli/cmd/service"
//	"github.com/lastbackend/lastbackend/pkg/cli/context"
//	storage "github.com/lastbackend/lastbackend/pkg/cli/storage/mock"
//	h "github.com/lastbackend/lastbackend/pkg/util/http"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestList(t *testing.T) {
//
//	const (
//		nName = "nspace name"
//
//		sName1 = "service name1"
//		sDesc1 = "service desc1"
//
//		sName2 = "service name2"
//		sDesc2 = "service desc2"
//	)
//
//	var (
//		err error
//		ctx = context.Mock()
//
//		ns = &n.Namespace{
//			Meta: n.NamespaceMeta{
//				Name: nName,
//			},
//		}
//	)
//
//	strg, err := storage.Get()
//	assert.NoError(t, err)
//	ctx.SetStorage(strg)
//	defer strg.Namespace().Remove()
//
//	//------------------------------------------------------------------------------------------
//	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//		nspaceJSON, err := json.Marshal(sm.ServiceList{
//			&sm.Service{
//				Meta: sm.ServiceMeta{
//					Name:        sName1,
//					Description: sDesc1,
//				},
//			},
//			&sm.Service{
//				Meta: sm.ServiceMeta{
//					Name:        sName2,
//					Description: sDesc2,
//				},
//			},
//		})
//		assert.NoError(t, err)
//
//		w.WriteHeader(200)
//		_, err = w.Write(nspaceJSON)
//		assert.NoError(t, err)
//	}))
//	defer server.Close()
//	//------------------------------------------------------------------------------------------
//
//	strg.Namespace().Save(ns)
//	assert.NoError(t, err)
//
//	client, err := h.New(server.URL, &h.ReqOpts{})
//	assert.NoError(t, err)
//	ctx.SetHttpClient(client)
//
//	sl, _, err := service.List()
//	assert.NoError(t, err)
//	assert.Equal(t, sName1, (*sl)[0].Meta.Name)
//	assert.Equal(t, sDesc1, (*sl)[0].Meta.Description)
//	assert.Equal(t, sName2, (*sl)[1].Meta.Name)
//	assert.Equal(t, sDesc2, (*sl)[1].Meta.Description)
//}
