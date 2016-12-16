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

func TestUpdate(t *testing.T) {

	const (
		name           string = "project"
		newProjectName string = "newname"
		description    string = "new description"
		token                 = "mocktoken"
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

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			return
		}

		var d = struct {
			Name        string `json:"name,omitempty"`
			Description string `json:"description,omitempty"`
		}{}

		err = json.Unmarshal(body, &d)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, d.Name, newProjectName, "they should be equal")
		assert.Equal(t, d.Description, description, "they should be equal")

		w.WriteHeader(200)
		_, err = w.Write([]byte(`{"id":"mock", "name":"` + name + `", "description":"` + description + `"}`))
		if err != nil {
			t.Error(err)
			return
		}
	}))
	defer server.Close()
	//------------------------------------------------------------------------------------------

	err = ctx.Storage.Set("project", switchData)
	if err != nil {
		t.Error(err)
		return
	}

	ctx.HTTP = h.New(server.URL)
	err = project.Update(name, newProjectName, description)
	if err != nil {
		t.Error(err)
		return
	}

	err = ctx.Storage.Get("project", projectmodel)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, projectmodel.Name, newProjectName)
	assert.Equal(t, projectmodel.Description, description)
}
