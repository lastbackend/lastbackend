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
	"github.com/lastbackend/lastbackend/internal/server/server/middleware"
	h "github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/tools/logger"
)

const (
	logPrefix = "api:handler:pod"
)

// Handler represent the http handler for pod
type Handler struct {
}

// NewPodHandler will initialize the pod resources endpoint
func NewPodHandler(r *mux.Router, mw middleware.Middleware) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:> init pod routes", logPrefix)

	handler := &Handler{}

	r.Handle("/namespace/{namespace}/service/{service}/deployment/{deployment}/pod", h.Handle(mw.Authenticate(handler.PodListH))).Methods(http.MethodGet)
}

func (handler Handler) PodListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/service/{service}/deployment/{pod} pod podList
	//
	// Shows a list of pods
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: namespace
	//     in: path
	//     description: name of the namespace
	//     required: true
	//     type: string
	//   - name: service
	//     in: path
	//     description: name of the service
	//     required: true
	//     type: string
	//   - name: deployment
	//     in: path
	//     description: name of the deployment
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Deployment list response
	//     schema:
	//       "$ref": "#/definitions/views_pod_list"
	//   '404':
	//     description: Namespace not found / Service not found / Deployment not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)
	//
	//sid := util.Vars(r)["service"]
	//nid := util.Vars(r)["namespace"]
	//did := util.Vars(r)["deployment"]
	//
	//log.Debugf("%s:list:> get pod list for `%s/%s`", logPrefix, sid, nid)
	//
	//var (
	//	sm  = model.NewServiceModel(r.Context(), envs.Get().GetStorage())
	//	nsm = model.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
	//	dm  = model.NewDeploymentModel(r.Context(), envs.Get().GetStorage())
	//	pm  = model.NewPodModel(r.Context(), envs.Get().GetStorage())
	//)
	//
	//ns, err := nsm.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:list:> get namespace", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if ns == nil {
	//	err := errors.New("namespace not found")
	//	log.Errorf("%s:list:> get namespace", logPrefix, err.Error())
	//	errors.New("namespace").NotFound().Http(w)
	//	return
	//}
	//
	//srv, err := sm.Get(ns.Meta.Name, sid)
	//if err != nil {
	//	log.Errorf("%s:list:> get service by name `%s` in namespace `%s` err: %s", logPrefix, sid, ns.Meta.Name, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if srv == nil {
	//	log.Warnf("%s:list:> service `%s` in namespace `%s` not found", logPrefix, sid, ns.Meta.Name)
	//	errors.New("service").NotFound().Http(w)
	//	return
	//}
	//
	//dep, err := dm.Get(srv.Meta.Namespace, srv.Meta.Name, did)
	//if err != nil {
	//	log.Errorf("%s:list:> get deployment by deployment id `%s` err: %s", logPrefix, did, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//if dep == nil {
	//	log.Warnf("%s:list:> deployment `%s` in namespace `%s` not found", logPrefix, sid, did)
	//	errors.New("deployment").NotFound().Http(w)
	//	return
	//}
	//
	//pl, err := pm.ListByDeployment(ns.Meta.Name, srv.Meta.Name, dep.Meta.Name)
	//if err != nil {
	//	log.Errorf("%s:list:> get pod list by deployment name `%s` err: %s", logPrefix, did, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//response, err := v1.View().Pod().NewList(pl).ToJson()
	//if err != nil {
	//	log.Errorf("%s:list:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:list:> write response err: %s", logPrefix, err.Error())
		return
	}
}
