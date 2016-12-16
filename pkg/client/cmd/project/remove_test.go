package project_test

import (
	h "github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/lastbackend/lastbackend/libs/model"
	"time"
	"github.com/lastbackend/lastbackend/libs/db"
)

func TestRemove(t *testing.T) {

	const (
		name  string = "project"
		token string = "mocktoken"
	)

	var (
		err error
		ctx = context.Mock()
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
