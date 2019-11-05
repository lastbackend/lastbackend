//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package deployment

import (
	"github.com/lastbackend/lastbackend/internal/api/envs"
	"github.com/lastbackend/lastbackend/internal/api/http/deployment/deployment"
	"github.com/lastbackend/lastbackend/internal/api/http/namespace/namespace"
	"github.com/lastbackend/lastbackend/internal/api/http/service/service"
	"github.com/lastbackend/lastbackend/internal/api/types/v1"
	"github.com/lastbackend/lastbackend/internal/api/types/v1/request"
	"github.com/lastbackend/lastbackend/internal/pkg/model"
	"net/http"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/util/http/utils"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logLevel  = 2
	logPrefix = "api:handler:deployment"
)

func DeploymentListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/service/{service}/deployment deployment deploymentList
	//
	// Shows a list of deployments
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
	// responses:
	//   '200':
	//     description: Deployment list response
	//     schema:
	//       "$ref": "#/definitions/views_deployment_list"
	//   '404':
	//     description: Namespace not found / Service not found
	//   '500':
	//     description: Internal server error

	sid := utils.Vars(r)["service"]
	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("%s:list:> get deployments list for `%s/%s`", logPrefix, sid, nid)

	var (
		sm  = model.NewServiceModel(r.Context(), envs.Get().GetStorage())
		nsm = model.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		dm  = model.NewDeploymentModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:list:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	srv, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get service by name `%s` in namespace `%s` err: %s", logPrefix, sid, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if srv == nil {
		log.V(logLevel).Warnf("%s:list:> service `%s` in namespace `%s` not found", logPrefix, sid, ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	dl, err := dm.ListByService(srv.Meta.Namespace, srv.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get deployment list by service id `%s` err: %s", logPrefix, srv.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Deployment().NewList(dl).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:list:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func DeploymentInfoH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/service/{service}/deployment/{deployment} deployment deploymentInfo
	//
	// Shows a deployment info
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
	//     description: Deployment response
	//     schema:
	//       "$ref": "#/definitions/views_deployment"
	//   '404':
	//     description: Namespace not found / Service not found
	//   '500':
	//     description: Internal server error

	sid := utils.Vars(r)["service"]
	nid := utils.Vars(r)["namespace"]
	did := utils.Vars(r)["deployment"]

	log.V(logLevel).Debugf("%s:info:> get deployments list for `%s/%s`", logPrefix, sid, nid)

	var (
		sm  = model.NewServiceModel(r.Context(), envs.Get().GetStorage())
		nsm = model.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		dm  = model.NewDeploymentModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get namespace %s err: %s", logPrefix, nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err = errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:info:> namespace %s not found", logPrefix, nid)
		errors.New("namespace").NotFound().Http(w)
		return
	}

	srv, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get service `%s` err: %s", logPrefix, srv.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if srv == nil {
		log.V(logLevel).Warnf("%s:info:> service `%s` not found", logPrefix, sid)
		errors.New("service").NotFound().Http(w)
		return
	}

	d, err := dm.Get(srv.Meta.Namespace, srv.Meta.Name, did)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get deployment by name `%s` err: %s", logPrefix, did, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if d == nil {
		log.V(logLevel).Warnf("%s:info:> deployment `%s` not found", logPrefix, did)
		errors.New("deployment").NotFound().Http(w)
		return
	}

	response, err := v1.View().Deployment().New(d).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:info:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func DeploymentCreateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation POST /namespace/{namespace}/service/{service}/deployment/ deployment deploymentCreate
	//
	// Updates deployment parameters
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
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_deployment_update"
	// responses:
	//   '200':
	//     description: Deployment was successfully updated
	//     schema:
	//       "$ref": "#/definitions/views_deployment"
	//   '404':
	//     description: Namespace not found / Service not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["service"]

	log.V(logLevel).Debugf("%s:create:> create deployment `%s` in service `%s`", logPrefix, sid, nid)

	var (
		err error
	)

	// request body struct
	opts := v1.Request().Deployment().Manifest()
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:create:> validation incoming data err: %s", logPrefix, err.Err())
		err.Http(w)
		return
	}

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	svc, e := service.Fetch(r.Context(), ns.Meta.Name, sid)
	if e != nil {
		e.Http(w)
		return
	}

	dp, e := deployment.Create(r.Context(), ns, svc, opts)
	if e != nil {
		e.Http(w)
		return
	}

	response, err := v1.View().Deployment().New(dp).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:create:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func DeploymentUpdateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /namespace/{namespace}/service/{service}/deployment/{deployment} deployment deploymentUpdate
	//
	// Updates deployment parameters
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
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_deployment_update"
	// responses:
	//   '200':
	//     description: Deployment was successfully updated
	//     schema:
	//       "$ref": "#/definitions/views_deployment"
	//   '404':
	//     description: Namespace not found / Service not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["service"]
	did := utils.Vars(r)["deployment"]

	log.V(logLevel).Debugf("%s:update:> update deployment `%s` in service `%s`", logPrefix, sid, nid)

	var (
		err error
	)

	// request body struct
	opts := v1.Request().Deployment().Manifest()
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:update:> validation incoming data err: %s", logPrefix, err.Err())
		err.Http(w)
		return
	}

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	svc, e := service.Fetch(r.Context(), ns.Meta.Name, sid)
	if e != nil {
		e.Http(w)
		return
	}

	dp, e := deployment.Fetch(r.Context(), ns.Meta.Name, svc.Meta.Name, did)
	if e != nil {
		e.Http(w)
		return
	}

	dp, e = deployment.Update(r.Context(), ns, svc, dp, opts, &request.DeploymentUpdateOptions{})
	if e != nil {
		e.Http(w)
		return
	}

	response, err := v1.View().Deployment().New(dp).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:update:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func DeploymentRemoveH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /namespace/{namespace}/service/{service}/deployment/{deployment} deployment deploymentUpdate
	//
	// Updates deployment parameters
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
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_deployment_update"
	// responses:
	//   '200':
	//     description: Deployment was successfully updated
	//     schema:
	//       "$ref": "#/definitions/views_deployment"
	//   '404':
	//     description: Namespace not found / Service not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["service"]
	did := utils.Vars(r)["deployment"]

	log.V(logLevel).Debugf("%s:update:> remove deployment `%s` in service `%s`", logPrefix, did, sid)

	var (
		err error
		dm  = model.NewDeploymentModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	svc, e := service.Fetch(r.Context(), ns.Meta.Name, sid)
	if e != nil {
		e.Http(w)
		return
	}

	dp, e := deployment.Fetch(r.Context(), ns.Meta.Name, svc.Meta.Name, did)
	if e != nil {
		e.Http(w)
		return
	}
	if dp == nil {
		log.V(logLevel).Warnf("%s:remove:> deployment name `%s` in namespace `%s` not found", logPrefix, did, ns.Meta.Name)
		errors.New("deployment").NotFound().Http(w)
		return
	}

	if err := dm.Destroy(dp); err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove deployment err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Deployment().New(dp).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:update:> write response err: %s", logPrefix, err.Error())
		return
	}
}
