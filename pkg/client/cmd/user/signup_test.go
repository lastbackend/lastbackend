package user_test

import (
	"encoding/json"
	h "github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/user"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUp_Success(t *testing.T) {

	const (
		username string = "mock"
		email    string = "mock@lastbackend.com"
		password string = "mock123456"
	)

	var (
		err error
		ctx = context.Mock()
	)

	err = ctx.Storage.Init()
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
			Email    string `json:"email,omitempty"`
			Username string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
		}{}

		err = json.Unmarshal(body, &d)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, d.Username, username, "they should be equal")
		assert.Equal(t, d.Email, email, "they should be equal")
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

	err = user.SignUp(username, email, password)
	if err != nil {
		t.Error(err.Error())
		return
	}
}
