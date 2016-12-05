package project_test

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/libs/db"
	h "github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSwitch(t *testing.T) {
	var (
		err error
		ctx = context.Mock()
	)

	const token = "mocktoken"
	const mock_name = "mock_name"

	var mock_proj = model.Project{Name: "mock_name",
		ID:          "mock_id",
		Created:     time.Now(),
		Updated:     time.Now(),
		User:        "you",
		Description: "sample description"}

	ctx.Storage, err = db.Init()
	if err != nil {
		panic(err)

	}
	defer ctx.Storage.Close()

	ctx.Token = token
	//------------------------------------------------------------------------------------------
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tk := r.Header.Get("Authorization")
		assert.NotEmpty(t, tk, "token should be not empty")
		assert.Equal(t, tk, "Bearer "+token, "they should be equal")
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			t.Error(err)
			return
		}

		assert.Empty(t, body, "body should be empty")
		buff, err := json.Marshal(mock_proj)

		if err != nil {
			t.Error(err)
			return
		}

		w.WriteHeader(200)
		_, err = w.Write(buff)
		if err != nil {
			t.Error(err)
			return
		}
	}))
	defer server.Close()
	//------------------------------------------------------------------------------------------

	ctx.HTTP = h.New(server.URL)
	err = project.Switch(mock_name)

	if err != nil {
		t.Error(err)
		return
	}

	return

}
