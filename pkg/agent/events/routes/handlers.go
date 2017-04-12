package routes

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/agent/runtime"
	"github.com/lastbackend/lastbackend/pkg/daemon/node/views/v1"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"io/ioutil"
	"net/http"
)

func SetPods(w http.ResponseWriter, r *http.Request) {

	log := context.Get().GetLogger()
	log.Debug("Set pods to agent")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	data := v1.Spec{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error(err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	patch := v1.FromNodeSpec(data)
	runtime.Get().Sync(patch.Pods)
	w.WriteHeader(http.StatusOK)
}
