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

package project_test

import (
	"github.com/lastbackend/lastbackend/libs/db"
	h "github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRemove(t *testing.T) {

	const (
		name  string = "project"
		token string = "mocktoken"
	)

	var (
		err          error
		ctx          = context.Mock()
		projectmodel = new(model.Project)
		switchData   = model.Project{
			Name:        "project",
			ID:          "mock_id",
			User:        "mock_user",
			Description: "sample description",
			Created:     time.Now(),
			Updated:     time.Now(),
		}
	)

	ctx.Storage, err = db.Init()
	if err != nil {
		t.Error(err)
		return
	}
	defer (func() {
		err = ctx.Storage.Clear()
		if err != nil {
			t.Error(err)
			return
		}
		err = ctx.Storage.Close()
		if err != nil {
			t.Error(err)
			return
		}
	})()

	ctx.Token = token

	//------------------------------------------------------------------------------------------
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tk := r.Header.Get("Authorization")
		assert.NotEmpty(t, tk, "token should be not empty")
		assert.Equal(t, tk, "Bearer "+token, "they should be equal")

		w.WriteHeader(200)
		w.Write([]byte{})
	}))
	defer server.Close()
	//------------------------------------------------------------------------------------------

	err = ctx.Storage.Set("project", switchData)
	if err != nil {
		t.Error(err)
		return
	}

	ctx.HTTP = h.New(server.URL)
	err = project.Remove(name)
	if err != nil {
		t.Error(err)
	}

	err = ctx.Storage.Get("project", projectmodel)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, projectmodel.ID, "")
}
