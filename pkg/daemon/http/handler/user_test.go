package handler_test

import (
	"bytes"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	h "github.com/lastbackend/lastbackend/pkg/daemon/http"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserCreateH(t *testing.T) {

	_ = context.Mock()

	r := h.NewRouter()

	var json = `{"username":"mocked", "email":"mocked@mocked.com", "password":"mockedpassword"}`

	req, err := http.NewRequest("POST", "/user", bytes.NewBufferString(json))
	if err != nil {
		t.Fatal("Creating 'POST /user' request failed!")
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbSI6Im1vY2tlZEBtb2NrZWQuY29tIiwiZXhwIjoxNDkzMjI5NTEwLCJqdGkiOjE0OTMyMjk1MTAsIm9pZCI6IiIsInVpZCI6ImM0ZmU5NTFjLTNmMmEtNDJlNS1hYjAwLWM5NDM1ZDJmOTUwZiIsInVzZXIiOiJtb2NrZWQifQ.-xiAaTqbdwz50LQliQoqNsYYhtEIc77PKneLzDutAD4")
	req.Header.Add("Content-Type", "application/json")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK && res.Code != http.StatusBadRequest {
		t.Error("Server returned: ", res.Code, " instead of ", http.StatusBadRequest)
	}
}
