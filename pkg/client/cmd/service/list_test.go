package service_test

import (
	h "github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/service"
)

func TestList(t *testing.T) {

	const (
		name        string = "service"
		description string = "service describe"
		token       string = "mocktoken"
	)

	var (
		err error
		ctx = context.Mock()
	)

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

	ctx.HTTP = h.New(server.URL)

	_, err = service.List()
	if err != nil {
		t.Error(err)
	}
}
