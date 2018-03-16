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

package app_test

//
//import (
//	"encoding/json"
//	n "github.com/lastbackend/lastbackend/pkg/api/app/views/v1"
//	"github.com/lastbackend/lastbackend/pkg/cli/cmd/app"
//	"github.com/lastbackend/lastbackend/pkg/cli/context"
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
//		tName = "test name"
//		tDesc = "test description"
//	)
//
//	var (
//		err error
//		ctx = context.Mock()
//	)
//
//	//------------------------------------------------------------------------------------------
//	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		nspaceJSON, err := json.Marshal(n.App{
//			Meta: n.AppMeta{
//				Name:        tName,
//				Description: tDesc,
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
//	client, err := h.New(server.URL, &h.ReqOpts{})
//	assert.NoError(t, err)
//	ctx.SetHttpClient(client)
//
//	ns, err := app.Get(tName)
//	assert.NoError(t, err)
//	assert.Equal(t, tName, ns.Meta.Name)
//	assert.Equal(t, tDesc, ns.Meta.Description)
//}
