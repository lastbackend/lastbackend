package template_test

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/libs/db"
	e "github.com/lastbackend/lastbackend/libs/errors"
	h "github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const token = "mocktoken"

func TestList(t *testing.T) {
	var (
		err       error
		ctx       = context.Mock()
		er        = new(e.Http)
		templates = new(model.TemplateList)
	)

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

		var temp = make(map[string][]string)
		temp["test_temp_1"] = []string{"ver. 1.1", "ver 2.2"}
		temp["test_temp_2"] = []string{"ver. 1.1", "ver 2.2", "ver. 3.3"}
		byte, _ := json.Marshal(temp)

		w.WriteHeader(200)
		_, err = w.Write(byte)
		if err != nil {
			t.Error(err)
			return
		}
	}))
	defer server.Close()
	ctx.HTTP = h.New(server.URL)
	//------------------------------------------------------------------------------------------

	err = List()

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
