package project_test

import (
	"github.com/lastbackend/lastbackend/libs/db"
	h "github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/project"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRemove_Success(t *testing.T) {

	const (
		name  string = "project"
		token string = "mocktoken"
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

	ctx.Token = token

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tk := r.Header.Get("Authorization")
		assert.NotEmpty(t, tk, "token should be not empty")
		assert.Equal(t, tk, "Bearer "+token, "they should be equal")

		w.WriteHeader(200)
		w.Write([]byte{})
	}))
	defer server.Close()

	ctx.HTTP = h.New(server.URL)

	err = project.Remove(name)
	if err != nil {
		t.Error(err)
	}
}
