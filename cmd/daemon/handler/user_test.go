package handler_test

import (
	"bytes"
	"github.com/lastbackend/lastbackend/cmd/daemon"
	"github.com/lastbackend/lastbackend/cmd/daemon/context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserCreateH(t *testing.T) {

	_ = context.Mock()

	r := daemon.NewRouter()

	var json = `{"username":"mocked", "email":"mocked@mocked.com", "password":"mockedpassword"}`

	req, err := http.NewRequest("POST", "/user", bytes.NewBufferString(json))
	if err != nil {
		t.Fatal("Creating 'POST /user' request failed!")
	}

	req.Header.Add("Content-Type", "application/json")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK && res.Code != http.StatusBadRequest {
		t.Error("Server returned: ", res.Code, " instead of ", http.StatusBadRequest)
	}
}

func TestUserGetH(t *testing.T) {

	_ = context.Mock()

	r := daemon.NewRouter()

	req, err := http.NewRequest("GET", "/user", nil)
	if err != nil {
		t.Fatal("Creating 'GET /user' request failed!")
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbSI6Im1vY2tlZEBtb2NrZWQuY29tIiwiZXhwIjoxNDg3NDE3ODk1LCJqdGkiOjE0ODc0MTc4OTUsIm9pZCI6IiIsInVpZCI6IjU2MmYwY2EwLTI2ZWEtNGFiNC1hZDBmLTU1N2NmYjJmYjgwNyIsInVzZXIiOiJtb2NrZWQifQ.VjHgKRqJCwf7TDphHPHhMl6njwL7agE1dzPVeGy5HFI")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Error("Server returned: ", res.Code, " instead of ", http.StatusBadRequest)
	}
}
