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

package node

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/minion/server/middleware"
	h "github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/tools/logger"
)

const (
	logPrefix = "api:handler:node"
)

// Handler represent the http handler for node
type Handler struct {
}

// NewNodeHandler will initialize the node resources endpoint
func NewNodeHandler(r *mux.Router, mw middleware.Middleware) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:> init node routes", logPrefix)

	handler := &Handler{
	}

	r.Handle("/", h.Handle(mw.Authenticate(handler.NodeGetH))).Methods(http.MethodGet)
}

func (handler Handler) NodeGetH(w http.ResponseWriter, r *http.Request) {
	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	log.Debug("Handler: Node: list node")
}
