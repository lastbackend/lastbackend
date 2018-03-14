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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
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
	token         = "demotoken"
	namespaceName = "demo"
)

func TestNamespaceGet(t *testing.T) {

	strg, _ := mock.New()
	envs.Get().SetStorage(strg)
	viper.Set("verbose", 0)
	viper.Set("security.token", token)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/namespace/demo/sandbox", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Our handler might also expect an API access token.
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Set url vars for mux
	req = mux.SetURLVars(req, map[string]string{
		"namespace": namespaceName,
	})

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	res := httptest.NewRecorder()
	handler := middleware.Authenticate(NamespaceInfoH)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(res, req)

	// Check the status code is what we expect.
	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("read body err: %s", err)
		return
	}

	ns := new(types.Namespace)
	err = json.Unmarshal(body, ns)
	if err != nil {
		t.Errorf("convert struct from json err: %s", err)
		return
	}

	assert.Equal(t, ns.Meta.Name, namespaceName, "they should be equal")
}
