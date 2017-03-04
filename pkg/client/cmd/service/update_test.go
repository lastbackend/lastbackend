package service_test

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/libs/db"
	h "github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/service"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {

	var (
		name        string = "service"
		description string = "service describe"
		scale       int32  = 10
		token       string = "mocktoken"
	)

	var (
		err     error
		ctx     = context.Mock()
		project = model.Project{
			Name:        "mock_name",
			ID:          "mock_id",
			Created:     time.Now(),
			Updated:     time.Now(),
			User:        "mock_demo",
			Description: "sample description",
		}
		updateData = model.ServiceUpdateConfig{}
	)

	updateData.Name = &name
	updateData.Description = &description
	updateData.Replicas = &scale

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

		var d = model.ServiceUpdateConfig{}

		err = json.Unmarshal(body, &d)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, d.Name, &name, "they should be equal")
		assert.Equal(t, d.Description, &description, "they should be equal")
		assert.Equal(t, d.Replicas, &scale, "they should be equal")

		w.WriteHeader(200)
		_, err = w.Write([]byte{})
		if err != nil {
			t.Error(err)
			return
		}
	}))
	defer server.Close()
	//------------------------------------------------------------------------------------------

	err = ctx.Storage.Set("project", project)
	if err != nil {
		t.Error(err)
		return
	}

	ctx.HTTP = h.New(server.URL)
	err = service.Update(name, updateData)
	if err != nil {
		t.Error(err)
		return
	}
}
