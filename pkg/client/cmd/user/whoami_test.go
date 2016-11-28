package user_test

import (
	"fmt"
	h "github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/user"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWhoami_Success(t *testing.T) {

	const (
		username string = "mock"
		email    string = "mock@lastbackend.com"
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

	token := struct {
		Token string `json:"token"`
	}{"mocktoken"}

	err = ctx.Storage.Set("session", token)
	if err != nil {
		fmt.Println(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tk := r.Header.Get("Authorization")
		assert.NotEmpty(t, tk, "token should be not empty")
		assert.Equal(t, tk, "Bearer "+token.Token, "they should be equal")

		w.WriteHeader(200)
		_, err := w.Write([]byte(`{"id":"mock", "username":"` + username + `", "email":"` + email + `"}`))
		if err != nil {
			t.Error(err)
			return
		}
	}))
	defer server.Close()

	ctx.HTTP = h.New(server.URL)

	err = user.Whoami()
	if err != nil {
		t.Error(err)
	}
}
