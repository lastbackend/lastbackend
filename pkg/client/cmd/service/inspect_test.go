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
//	"github.com/lastbackend/lastbackend/pkg/apis/types"
//	"github.com/lastbackend/lastbackend/pkg/client/cmd/service"
//	"github.com/lastbackend/lastbackend/pkg/client/context"
//	"github.com/lastbackend/lastbackend/pkg/client/storage"
//	h "github.com/lastbackend/lastbackend/pkg/util/http"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//	"time"
//)
//
//func TestGet(t *testing.T) {
//
//	const (
//		name  string = "project"
//		token string = "mocktoken"
//	)
//
//	var (
//		err  error
//		ctx  = context.Mock()
//		data = types.Namespace{
//			Name:        "mock_name",
//			Created:     time.Now(),
//			Updated:     time.Now(),
//			User:        "mock_demo",
//			Description: "sample description",
//		}
//	)
//
//	ctx.Storage, err = storage.Init()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	defer (func() {
//		err = ctx.Storage.Clear()
//		if err != nil {
//			t.Error(err)
//			return
//		}
//		err = ctx.Storage.Close()
//		if err != nil {
//			t.Error(err)
//			return
//		}
//	})()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	defer ctx.Storage.Close()
//
//	ctx.Token = token
//
//	//------------------------------------------------------------------------------------------
//	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//		tk := r.Header.Get("Authorization")
//		assert.NotEmpty(t, tk, "token should be not empty")
//		assert.Equal(t, tk, "Bearer "+token, "they should be equal")
//
//		w.WriteHeader(200)
//		_, err := w.Write([]byte(`{"id":"mock", "name":"` + name + `"}`))
//		if err != nil {
//			t.Error(err)
//			return
//		}
//	}))
//	defer server.Close()
//	//------------------------------------------------------------------------------------------
//
//	err = ctx.Storage.Set("project", data)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	ctx.HTTP = h.New(server.URL)
//
//	_, _, err = service.Inspect(name)
//	if err != nil {
//		t.Error(err)
//	}
//}
