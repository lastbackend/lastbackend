package handler_test

import (
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	h "github.com/lastbackend/lastbackend/pkg/daemon/http"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDockerRepositorySearchH(t *testing.T) {

	_ = context.Mock()

	r := h.NewRouter()

	req, err := http.NewRequest("GET", "/docker/repo/search?name=redis", nil)
	if err != nil {
		t.Fatal("Creating 'GET /docker/repo/search?name=redis' request failed!")
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbSI6Im1vY2tlZEBtb2NrZWQuY29tIiwiZXhwIjoxNDkzMjI5NTEwLCJqdGkiOjE0OTMyMjk1MTAsIm9pZCI6IiIsInVpZCI6ImM0ZmU5NTFjLTNmMmEtNDJlNS1hYjAwLWM5NDM1ZDJmOTUwZiIsInVzZXIiOiJtb2NrZWQifQ.-xiAaTqbdwz50LQliQoqNsYYhtEIc77PKneLzDutAD4")
	req.Header.Add("Content-Type", "application/json")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Error("Server returned: ", res.Code, " instead of ", http.StatusBadRequest)
	}
}

func TestDockerRepositoryTagListH(t *testing.T) {

	_ = context.Mock()

	r := h.NewRouter()

	req, err := http.NewRequest("GET", "/docker/repo/tags?owner=library&name=redis", nil)
	if err != nil {
		t.Fatal("Creating 'GET /docker/repo/tags?owner=library&name=redis' request failed!")
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbSI6Im1vY2tlZEBtb2NrZWQuY29tIiwiZXhwIjoxNDkzMjI5NTEwLCJqdGkiOjE0OTMyMjk1MTAsIm9pZCI6IiIsInVpZCI6ImM0ZmU5NTFjLTNmMmEtNDJlNS1hYjAwLWM5NDM1ZDJmOTUwZiIsInVzZXIiOiJtb2NrZWQifQ.-xiAaTqbdwz50LQliQoqNsYYhtEIc77PKneLzDutAD4")
	req.Header.Add("Content-Type", "application/json")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Error("Server returned: ", res.Code, " instead of ", http.StatusBadRequest)
	}
}
