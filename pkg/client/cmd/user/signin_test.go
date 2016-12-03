package user_test

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/libs/db"
	h "github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/user"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignIn_Success(t *testing.T) {

	const (
		login    string = "mock"
		password string = "mock123456"
	)

	var (
		err error
		ctx = context.Mock()
	)

	ctx.Storage, err = db.Init()
	if err != nil {
		panic(err)
	}
	defer ctx.Storage.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			return
		}

		var d = struct {
			Login    string `json:"login,omitempty"`
			Password string `json:"password,omitempty"`
		}{}

		err = json.Unmarshal(body, &d)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, d.Login, login, "they should be equal")
		assert.Equal(t, d.Password, password, "they should be equal")

		w.WriteHeader(200)
		_, err = w.Write([]byte(`{"token":"mocktoken"}`))
		if err != nil {
			t.Error(err)
			return
		}
	}))
	defer server.Close()

	ctx.HTTP = h.New(server.URL)

	err = user.SignIn(login, password)
	if err != nil {
		t.Error(err.Error())
		return
	}
}
