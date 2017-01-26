package handler_test

import (
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	h "github.com/lastbackend/lastbackend/pkg/daemon/http"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceHookListH(t *testing.T) {

	_ = c.Mock()

	r := h.NewRouter()

	req, err := http.NewRequest("GET", "/project/mocked/service/mocked/hook", nil)
	if err != nil {
		t.Fatal("Getting 'GET /project/mocked/service/mocked/hook' request failed!")
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbSI6Im1vY2tlZEBtb2NrZWQuY29tIiwiZXhwIjoxNDkzMjI5NTEwLCJqdGkiOjE0OTMyMjk1MTAsIm9pZCI6IiIsInVpZCI6ImM0ZmU5NTFjLTNmMmEtNDJlNS1hYjAwLWM5NDM1ZDJmOTUwZiIsInVzZXIiOiJtb2NrZWQifQ.-xiAaTqbdwz50LQliQoqNsYYhtEIc77PKneLzDutAD4")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Error("Server returned: ", res.Code, " instead of ", http.StatusBadRequest)
	}
}

func TestHookInsertH(t *testing.T) {

	_ = c.Mock()

	r := h.NewRouter()

	req, err := http.NewRequest("POST", "/project/mocked/service/mocked/hook", nil)
	if err != nil {
		t.Fatal("Creating 'POST /project/mocked/service/mocked/hook' request failed!")
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbSI6Im1vY2tlZEBtb2NrZWQuY29tIiwiZXhwIjoxNDkzMjI5NTEwLCJqdGkiOjE0OTMyMjk1MTAsIm9pZCI6IiIsInVpZCI6ImM0ZmU5NTFjLTNmMmEtNDJlNS1hYjAwLWM5NDM1ZDJmOTUwZiIsInVzZXIiOiJtb2NrZWQifQ.-xiAaTqbdwz50LQliQoqNsYYhtEIc77PKneLzDutAD4")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Error("Server returned: ", res.Code, " instead of ", http.StatusBadRequest)
	}
}

func TestHookRemoveH(t *testing.T) {

	_ = c.Mock()

	r := h.NewRouter()

	req, err := http.NewRequest("DELETE", "/project/mocked/service/mocked/hook/mocked", nil)
	if err != nil {
		t.Fatal("Deleting 'DELETE /project/mocked/service/mocked/hook/mocked' request failed!")
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbSI6Im1vY2tlZEBtb2NrZWQuY29tIiwiZXhwIjoxNDkzMjI5NTEwLCJqdGkiOjE0OTMyMjk1MTAsIm9pZCI6IiIsInVpZCI6ImM0ZmU5NTFjLTNmMmEtNDJlNS1hYjAwLWM5NDM1ZDJmOTUwZiIsInVzZXIiOiJtb2NrZWQifQ.-xiAaTqbdwz50LQliQoqNsYYhtEIc77PKneLzDutAD4")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Error("Server returned: ", res.Code, " instead of ", http.StatusBadRequest)
	}
}
