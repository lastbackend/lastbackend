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

package cluster

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/server/server/middleware"
	h "github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/tools/logger"
)

const (
	logPrefix = "api:handler:cluster"
)

// Handler represent the http handler for cluster
type Handler struct {
}

// NewClusterHandler will initialize the cluster resources endpoint
func NewClusterHandler(r *mux.Router, mw middleware.Middleware) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:> init cluster routes", logPrefix)

	handler := &Handler{}

	r.Handle("/cluster", h.Handle(mw.Authenticate(handler.ClusterInfoH))).Methods(http.MethodGet)
}

func (handler Handler) ClusterInfoH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster cluster clusterInfo
	//
	// Shows an info about current cluster
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: Cluster response
	//     schema:
	//       "$ref": "#/definitions/views_cluster"
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	log.Debugf("%s:info:> get cluster", logPrefix)

	//var clm = model.NewClusterModel(r.Context(), envs.Get().GetStorage())
	//
	//cl, err := clm.Get()
	//if err != nil {
	//	log.Errorf("%s:info:> get cluster err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//if cl == nil {
	//	cl = new(types.Cluster)
	//	cl.Meta.SetDefault()
	//	cl.Meta.Name, cl.Meta.Description = envs.Get().GetClusterInfo()
	//	cl.Meta.SelfLink = *types.NewClusterSelfLink(cl.Meta.Name)
	//}
	//
	//response, err := v1.View().Cluster().New(cl).ToJson()
	//if err != nil {
	//	log.Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:info:> write response err: %s", logPrefix, err.Error())
		return
	}
}
