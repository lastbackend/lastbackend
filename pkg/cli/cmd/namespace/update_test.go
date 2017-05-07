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

package namespace_test

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/cli/cmd/namespace"
	"github.com/lastbackend/lastbackend/pkg/cli/context"
	storage "github.com/lastbackend/lastbackend/pkg/cli/storage/mock"
	n "github.com/lastbackend/lastbackend/pkg/api/namespace/views/v1"
	h "github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdate(t *testing.T) {

	const (
		tName = "test name"
		tDesc = "test description"

		nName = "new name"
		nDesc = "new description"
	)

	var (
		err error
		ctx = context.Mock()

		ns = &n.Namespace{
			Meta: n.NamespaceMeta{
				Name:        tName,
				Description: tDesc,
			},
		}
	)

	strg, err := storage.Get()
	assert.NoError(t, err)
	ctx.SetStorage(strg)
	defer strg.Namespace().Remove()

	//------------------------------------------------------------------------------------------
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)

		var rData struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		err = json.Unmarshal(body, &rData)
		assert.NoError(t, err)

		assert.Equal(t, nName, rData.Name)
		assert.Equal(t, nDesc, rData.Description)

		wData, err := json.Marshal(rData)
		assert.NoError(t, err)

		w.WriteHeader(200)
		_, err = w.Write(wData)
		assert.NoError(t, err)
	}))
	defer server.Close()
	//------------------------------------------------------------------------------------------

	strg.Namespace().Save(ns)
	assert.NoError(t, err)

	client, err := h.New(server.URL, &h.ReqOpts{})
	assert.NoError(t, err)
	ctx.SetHttpClient(client)

	err = namespace.Update(tName, nName, nDesc)
	assert.NoError(t, err)

	ns, err = strg.Namespace().Load()
	assert.NoError(t, err)
	assert.Equal(t, ns.Meta.Name, nName)
	assert.Equal(t, ns.Meta.Description, nDesc)
}
