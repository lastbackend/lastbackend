//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package routes_test

//import (
//	"github.com/lastbackend/lastbackend/pkg/daemon/context"
//	h "github.com/lastbackend/lastbackend/pkg/daemon/http"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestDockerRepositorySearchH(t *testing.T) {
//
//	_ = context.Mock()
//
//	r := h.NewRouter()
//
//	req, err := http.NewRequest("GET", "/docker/repo/search?name=redis", nil)
//	if err != nil {
//		t.Fatal("Creating 'GET /docker/repo/search?name=redis' request failed!")
//	}
//
//	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbSI6Im1vY2tlZEBtb2NrZWQuY29tIiwiZXhwIjoxNDkzMjI5NTEwLCJqdGkiOjE0OTMyMjk1MTAsIm9pZCI6IiIsInVpZCI6ImM0ZmU5NTFjLTNmMmEtNDJlNS1hYjAwLWM5NDM1ZDJmOTUwZiIsInVzZXIiOiJtb2NrZWQifQ.-xiAaTqbdwz50LQliQoqNsYYhtEIc77PKneLzDutAD4")
//	req.Header.Add("Content-Type", "application/json")
//
//	res := httptest.NewRecorder()
//	r.ServeHTTP(res, req)
//
//	if res.Code != http.StatusOK {
//		t.Error("Server returned: ", res.Code, " instead of ", http.StatusBadRequest)
//	}
//}
//
//func TestDockerRepositoryTagListH(t *testing.T) {
//
//	_ = context.Mock()
//
//	r := h.NewRouter()
//
//	req, err := http.NewRequest("GET", "/docker/repo/tags?owner=library&name=redis", nil)
//	if err != nil {
//		t.Fatal("Creating 'GET /docker/repo/tags?owner=library&name=redis' request failed!")
//	}
//
//	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbSI6Im1vY2tlZEBtb2NrZWQuY29tIiwiZXhwIjoxNDkzMjI5NTEwLCJqdGkiOjE0OTMyMjk1MTAsIm9pZCI6IiIsInVpZCI6ImM0ZmU5NTFjLTNmMmEtNDJlNS1hYjAwLWM5NDM1ZDJmOTUwZiIsInVzZXIiOiJtb2NrZWQifQ.-xiAaTqbdwz50LQliQoqNsYYhtEIc77PKneLzDutAD4")
//	req.Header.Add("Content-Type", "application/json")
//
//	res := httptest.NewRecorder()
//	r.ServeHTTP(res, req)
//
//	if res.Code != http.StatusOK {
//		t.Error("Server returned: ", res.Code, " instead of ", http.StatusBadRequest)
//	}
//}
