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
//	"github.com/lastbackend/lastbackend/pkg/client/cmd/service"
//	"github.com/lastbackend/lastbackend/pkg/client/context"
//	s "github.com/lastbackend/lastbackend/pkg/client/storage"
//	n "github.com/lastbackend/lastbackend/pkg/api/namespace/views/v1"
//	h "github.com/lastbackend/lastbackend/pkg/util/http"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestRemove(t *testing.T) {
//
//	const (
//		sName = "service name"
//
//		nName = "namespace name"
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
//		w.WriteHeader(200)
//		w.Write([]byte{})
//	}))
//	defer server.Close()
//	//------------------------------------------------------------------------------------------
//
//	ctx.SetHttpClient(h.New(server.URL[7:]))
//
//	storage.Set(storageName, data)
//
//	err = service.Remove(sName)
//	assert.NoError(t, err)
//}
