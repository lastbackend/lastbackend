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
	n "github.com/lastbackend/lastbackend/pkg/api/namespace/views/v1"
	"github.com/lastbackend/lastbackend/pkg/cli/cmd/namespace"
	"github.com/lastbackend/lastbackend/pkg/cli/context"
	h "github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestList(t *testing.T) {

	const (
		tName1 = "test name1"
		tDesc1 = "test description1"

		tName2 = "test name2"
		tDesc2 = "test description2"
	)

	var (
		err error
		ctx = context.Mock()
	)

	//------------------------------------------------------------------------------------------
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		nspaceJSON, err := json.Marshal(n.NamespaceList{
			&n.Namespace{
				Meta: n.NamespaceMeta{
					Name:        tName1,
					Description: tDesc1,
				},
			},
			&n.Namespace{
				Meta: n.NamespaceMeta{
					Name:        tName2,
					Description: tDesc2,
				},
			},
		})
		assert.NoError(t, err)

		w.WriteHeader(200)
		_, err = w.Write(nspaceJSON)
		assert.NoError(t, err)
	}))
	defer server.Close()
	//------------------------------------------------------------------------------------------

	client, err := h.New(server.URL, &h.ReqOpts{})
	assert.NoError(t, err)
	ctx.SetHttpClient(client)

	ns, err := namespace.List()
	assert.NoError(t, err)
	assert.Equal(t, tName1, (*ns)[0].Meta.Name)
	assert.Equal(t, tDesc1, (*ns)[0].Meta.Description)
	assert.Equal(t, tName2, (*ns)[1].Meta.Name)
	assert.Equal(t, tDesc2, (*ns)[1].Meta.Description)
}
