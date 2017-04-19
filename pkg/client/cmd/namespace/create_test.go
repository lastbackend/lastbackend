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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/namespace"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	s "github.com/lastbackend/lastbackend/pkg/client/storage"
	n "github.com/lastbackend/lastbackend/pkg/daemon/namespace/views/v1"
	h "github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreate(t *testing.T) {

	const (
		name        string = "test name"
		description string = "test description"
	)

	var (
		err error
		ctx = context.Mock()

		prct = new(types.Namespace)
	)

	storage, err := s.Init()
	assert.NoError(t, err)
	ctx.SetStorage(storage)

	// ----------
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)

		var d = struct {
			Name        string `json:"name,omitempty"`
			Description string `json:"description,omitempty"`
		}{}

		err = json.Unmarshal(body, &d)
		assert.NoError(t, err)

		assert.Equal(t, name, d.Name)
		assert.Equal(t, description, d.Description)

		aaaa, err := json.Marshal(n.Namespace{
			Meta: n.NamespaceMeta{
				Name:        name,
				Description: description,
			},
		})
		assert.NoError(t, err)

		w.WriteHeader(200)
		_, err = w.Write(aaaa)
		assert.NoError(t, err)
	}))
	defer server.Close()
	// ----------

	ctx.SetHttpClient(h.New(server.URL[7:]))

	err = namespace.Create(name, description)
	assert.NoError(t, err)

	err = storage.Get("namespace", prct)
	assert.NoError(t, err)

	assert.Equal(t, name, prct.Meta.Name)
	assert.Equal(t, description, prct.Meta.Description)
}
