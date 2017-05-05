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
//	"github.com/lastbackend/lastbackend/pkg/client/cmd/service"
//	"github.com/lastbackend/lastbackend/pkg/client/context"
//	s "github.com/lastbackend/lastbackend/pkg/client/storage"
//	n "github.com/lastbackend/lastbackend/pkg/api/namespace/views/v1"
//	sm "github.com/lastbackend/lastbackend/pkg/api/service/views/v1"
//	h "github.com/lastbackend/lastbackend/pkg/util/http"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestGet(t *testing.T) {
//
//	const (
//		nName = "nspace name"
//
//		sName = "service name"
//		sDesc = "service desc"
//
//		storageName = "test"
//	)
//
//	var (
//		err error
//		ctx = context.Mock()
//
//		data = n.Namespace{
//			Meta: n.NamespaceMeta{
//				Name: nName,
//			},
//		}
//	)
//
//	storage, err := s.Init()
//	assert.NoError(t, err)
//	ctx.SetStorage(storage)
//	defer func() {
//		storage.Clear()
//	}()
//
//	//------------------------------------------------------------------------------------------
//	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//		nspaceJSON, err := json.Marshal(sm.Service{
//			Meta: sm.ServiceMeta{
//				Name:        sName,
//				Description: sDesc,
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
//	err = storage.Set(storageName, data)
//	assert.NoError(t, err)
//
//	ctx.SetHttpClient(h.New(server.URL[7:]))
//
//	res, _, err := service.Inspect(nName)
//	assert.NoError(t, err)
//	assert.Equal(t, sName, res.Meta.Name)
//	assert.Equal(t, sDesc, res.Meta.Description)
//}
