//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package pod

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/agent/server/middleware"
	"github.com/lastbackend/lastbackend/internal/agent/state"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	h "github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/tools/logger"
)

const (
	logPrefix = "api:handler:pod"
)

// Handler represent the http handler for pod
type Handler struct {
	state *state.State
}

// NewPodHandler will initialize the pod resources endpoint
func NewPodHandler(r *mux.Router, mw middleware.Middleware, state *state.State) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:> init pod routes", logPrefix)

	handler := &Handler{
		state: state,
	}

	r.Handle("/pod/{pod}", h.Handle(mw.Authenticate(handler.PodGetH))).Methods(http.MethodGet)
	r.Handle("/pod/{pod}/{container}/logs", h.Handle(mw.Authenticate(handler.PodLogsH))).Methods(http.MethodGet)
}

func (handler Handler) PodGetH(w http.ResponseWriter, r *http.Request) {
	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)
	log.Debug("node:http:pod:get:> get pod info")
}

// PodLogsH handler streams pod logs into response writer
func (handler Handler) PodLogsH(w http.ResponseWriter, r *http.Request) {
	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)
	log.Debug("node:http:pod:get:> get pod logs")

	var (
		c      = mux.Vars(r)["container"]
		p      = handler.state.Pods().GetPod(mux.Vars(r)["pod"])
		notify = w.(http.CloseNotifier).CloseNotify()
		done   = make(chan bool, 1)
	)

	go func() {
		<-notify
		log.Debug("HTTP connection just closed.")
		done <- true
	}()

	if p == nil {
		log.Errorf("node:http:pod:get:> pod not found")
		errors.New("pod").NotFound().Http(w)
		return
	}

	if _, ok := p.Runtime.Services[c]; !ok {
		log.Errorf("node:http:pod:get:> container not found")
		errors.New("pod").NotFound().Http(w)
		return
	}

	// TODO: Get stream for pod logs
	//if err := runtime.PodLogs(r.Context(), c, true, w, done); err != nil {
	//	log.Errorf("node:http:pod:get:> get pod logs err: %s", err.Error())
	//}

	return
}
