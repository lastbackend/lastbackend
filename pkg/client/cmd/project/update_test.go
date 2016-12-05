package project_test

/*
	res["mock template 1"] = []string{"first ver.", "last ver."}
	res["mock template 2"] = []string{"first ver.","ver 0.0", "last ver."}
*/

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/libs/db"
	e "github.com/lastbackend/lastbackend/libs/errors"
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

func TestList(t *testing.T) {
	var (
		err       error
		ctx       = context.Mock()
		er        = new(e.Http)
		templates = new(model.TemplateList)
	)

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

		assert.NotEmpty(t, body, "body should not be empty")

		db_project := new(model.Project)
		reqw_project := new(model.Project)

		json.Unmarshal(body, &db_project)
		data, err := db.Init()
		if err != nil {
			t.Error(err)
			return
		}
		defer data.Close()
		err = data.Get("project", &db_project)

		if err != nil {
			t.Error(err)
			return
		}

		if reqw_project.Name != "" {
			db_project.Name = reqw_project.Name
		}

		db_project.Description = reqw_project.Description

		db_project.Updated = time.Now()

		err = data.Set("project", db_project)

		if err != nil {
			t.Error(err)
			return
		}

		w.WriteHeader(200)
		_, err = w.Write(nil)
		if err != nil {
			t.Error(err)
			return
		}
	}))
	defer server.Close()
	ctx.HTTP = h.New(server.URL)
	//------------------------------------------------------------------------------------------

	err = project.Update("mock_name", "mock desc")

	if err != nil {
		t.Error(err)
		return
	}

	if er.Code == 401 {
		t.Error("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
		return
	}

	if er.Code != 0 {
		t.Error(er.Code)
		return
	}

	templates.DrawTable()

	return

}
