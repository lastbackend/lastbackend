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

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/service"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	s "github.com/lastbackend/lastbackend/pkg/client/storage"
	n "github.com/lastbackend/lastbackend/pkg/api/namespace/views/v1"
	sm "github.com/lastbackend/lastbackend/pkg/api/service/views/v1"
	h "github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreate(t *testing.T) {

	const (
		nName = "nspace name"

		sName = "service name"
		sDesc = "service desc"

		storageName = "test"
	)

	var (
		err error
		ctx = context.Mock()

		data = n.Namespace{
			Meta: n.NamespaceMeta{
				Name: nName,
			},
		}
	)

	storage, err := s.Init()
	assert.NoError(t, err)
	ctx.SetStorage(storage)
	defer func() {
		storage.Clear()
	}()

	//------------------------------------------------------------------------------------------
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)

		var d = struct {
			Name string `json:"name,omitempty"`
		}{}

		err = json.Unmarshal(body, &d)
		assert.NoError(t, err)

		assert.Equal(t, sName, d.Name)

		nspaceJSON, err := json.Marshal(sm.Service{
			Meta: sm.ServiceMeta{
				Name:        sName,
				Description: sDesc,
			},
		})
		assert.NoError(t, err)

		w.WriteHeader(200)
		_, err = w.Write(nspaceJSON)
		assert.NoError(t, err)
	}))
	defer server.Close()
	//------------------------------------------------------------------------------------------

	storage.Set(storageName, data)

	ctx.SetHttpClient(h.New(server.URL[7:]))

	err = service.Create(sName, "redis", "", "", nil)
	assert.NoError(t, err)
}
