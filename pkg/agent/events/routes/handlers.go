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

package routes

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/agent/runtime"
	"github.com/lastbackend/lastbackend/pkg/api/node/views/v1"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"io/ioutil"
	"net/http"
	"github.com/lastbackend/lastbackend/pkg/log"
)

func SetPods(w http.ResponseWriter, r *http.Request) {

	log.Debug("Set pods to agent")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	data := v1.Spec{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error(err)
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	patch := v1.FromNodeSpec(data)
	runtime.Get().Sync(patch.Pods)
	w.WriteHeader(http.StatusOK)
}
