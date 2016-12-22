package service_test

import (
	"github.com/lastbackend/lastbackend/libs/db"
	h "github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/service"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestList(t *testing.T) {

	const (
		name        string = "service"
		description string = "service describe"
		token       string = "mocktoken"
	)

	var (
		err  error
		ctx  = context.Mock()
		data = model.Project{
			Name:        "mock_name",
			ID:          "mock_id",
			Created:     time.Now(),
			Updated:     time.Now(),
			User:        "mock_demo",
			Description: "sample description",
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
	if err != nil {
		t.Error(err)
		return
	}
	defer ctx.Storage.Close()

	ctx.Token = token

	//------------------------------------------------------------------------------------------
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tk := r.Header.Get("Authorization")

		assert.NotEmpty(t, tk, "token should be not empty")
		assert.Equal(t, tk, "Bearer "+token, "they should be equal")

		w.WriteHeader(200)
		_, err := w.Write([]byte(`[{"id":"mock", "name":"` + name + `", "description":"` + description + `"}]`))
		if err != nil {
			t.Error(err)
			return
		}
	}))
	defer server.Close()
	//------------------------------------------------------------------------------------------

	err = ctx.Storage.Set("project", data)
	if err != nil {
		t.Error(err)
		return
	}

	ctx.HTTP = h.New(server.URL)

	_, _, err = service.List()
	if err != nil {
		t.Error(err)
	}
}
