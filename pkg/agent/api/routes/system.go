package routes

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/agent/runtime"
	"github.com/lastbackend/lastbackend/pkg/apis/views/v1/node"
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

	data := node.Spec{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error(err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	patch := node.FromNodeSpec(data)
	runtime.Get().Sync(patch)

	w.WriteHeader(http.StatusOK)
}
