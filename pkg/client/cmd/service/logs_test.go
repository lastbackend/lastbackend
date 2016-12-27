package service_test

import (
	h "github.com/lastbackend/lastbackend/libs/http"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/service"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogs(t *testing.T) {

	const (
		projectName   string = "project"
		serviceName   string = "service"
		podName       string = "pod"
		containerName string = "container"
		token         string = "mocktoken"
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

		var (
			query          = r.URL.Query()
			podQuery       = query.Get("pod")
			containerQuery = query.Get("container")
		)

		assert.Equal(t, podQuery, podName, "they should be equal")
		assert.Equal(t, containerQuery, containerName, "they should be equal")

		w.WriteHeader(200)
	}))
	defer server.Close()
	//------------------------------------------------------------------------------------------

	ctx.HTTP = h.New(server.URL)

	_, err = service.Logs(projectName, serviceName, podName, containerName)
	if err != nil {
		t.Error(err)
	}
}
