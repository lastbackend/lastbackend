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

package namespace

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/views/v1"
	"github.com/lastbackend/lastbackend/pkg/storage/mock"
	"github.com/lastbackend/lastbackend/pkg/util/http/middleware"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	token                  = "demotoken"
	namespaceExistsName    = "demo"
	namespaceNotExistsName = "notexistsname"
)

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}

// Testing NamespaceInfoH handler of a successful request (status 200)
func TestNamespaceGetWithoutMiddleware(t *testing.T) {

	a := assert.New(t)

	strg, _ := mock.New()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", fmt.Sprintf("/namespace/%s", namespaceExistsName), nil)
	a.NoError(err)

	r := mux.NewRouter()
	r.HandleFunc("/namespace/{namespace}", NamespaceInfoH)

	setRequestVars(r, req)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	res := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	r.ServeHTTP(res, req)

	// Check the status code is what we expect.
	a.Equal(http.StatusNotFound, res.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", res.Code, http.StatusOK))

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("read body err: %s", err)
		return
	}

	ns := new(v1.Namespace)
	err = json.Unmarshal(body, ns)
	a.NoError(err)
	a.Equal(t, ns.Meta.Name, namespaceExistsName, "they should be equal")
}

// Testing NamespaceInfoH handler of a successful request (status 200)
func TestNamespaceGetWithAuthenticateMiddleware(t *testing.T) {

	a := assert.New(t)

	strg, _ := mock.New()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)
	viper.Set("security.token", token)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", fmt.Sprintf("/namespace/%s", namespaceExistsName), nil)
	a.NoError(err)

	// Our handler might also expect an API access token.
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	r := mux.NewRouter()
	r.HandleFunc("/namespace/{namespace}", middleware.Authenticate(NamespaceInfoH))

	setRequestVars(r, req)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	res := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	r.ServeHTTP(res, req)

	// Check the status code is what we expect.
	a.Equal(http.StatusNotFound, res.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", res.Code, http.StatusOK))

	body, err := ioutil.ReadAll(res.Body)
	a.NoError(err)

	ns := new(v1.Namespace)
	err = json.Unmarshal(body, ns)
	a.NoError(err)
	a.Equal(t, ns.Meta.Name, namespaceExistsName, "they should be equal")
}

// Testing NamespaceInfoH handler of a status 404
func TestNamespaceGetCheckStatusNotFound(t *testing.T) {

	a := assert.New(t)

	strg, _ := mock.New()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", fmt.Sprintf("/namespace/%s", namespaceNotExistsName), nil)
	a.NoError(err)

	r := mux.NewRouter()
	r.HandleFunc("/namespace/{namespace}", NamespaceInfoH)

	setRequestVars(r, req)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	res := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	r.ServeHTTP(res, req)

	// Check the status code is what we expect.
	a.Equal(http.StatusNotFound, res.Code, fmt.Sprintf("handler returned wrong status code: got %v want %v", res.Code, http.StatusOK))
}
