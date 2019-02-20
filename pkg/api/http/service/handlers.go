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

package service

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/http/namespace/namespace"
	"github.com/lastbackend/lastbackend/pkg/api/http/service/service"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

const (
	logLevel    = 2
	logPrefix   = "api:handler:service"
	BUFFER_SIZE = 512
)

func ServiceListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/service service serviceList
	//
	// Shows a list of services
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: namespace
	//     in: path
	//     description: namespace id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Service list response
	//     schema:
	//       "$ref": "#/definitions/views_service_list"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("%s:list:> list services in %s", logPrefix, nid)

	var (
		stg = envs.Get().GetStorage()
		sm  = distribution.NewServiceModel(r.Context(), stg)
		nsm = distribution.NewNamespaceModel(r.Context(), stg)
		dm  = distribution.NewDeploymentModel(r.Context(), stg)
		pm  = distribution.NewPodModel(r.Context(), stg)
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

	items, err := sm.List(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get service list in namespace `%s` err: %s", logPrefix, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	dl, err := dm.ListByNamespace(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get pod list by service id `%s` err: %s", logPrefix, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	pl, err := pm.ListByNamespace(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get pod list by service id `%s` err: %s", logPrefix, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Service().NewList(items, dl, pl).ToJson()
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

func ServiceInfoH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/service/{service} service serviceInfo
	//
	// Shows an info about service
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: namespace
	//     in: path
	//     description: namespace id
	//     required: true
	//     type: string
	//   - name: service
	//     in: path
	//     description: service id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Service list response
	//     schema:
	//       "$ref": "#/definitions/views_service"
	//   '404':
	//     description: Namespace not found / Service not found
	//   '500':
	//     description: Internal server error

	sid := utils.Vars(r)["service"]
	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("%s:info:> get service `%s` in namespace `%s`", logPrefix, sid, nid)

	var (
		dm  = distribution.NewDeploymentModel(r.Context(), envs.Get().GetStorage())
		pdm = distribution.NewPodModel(r.Context(), envs.Get().GetStorage())
	)

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

	dl, err := dm.ListByService(svc.Meta.Namespace, svc.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get pod list by service id `%s` err: %s", logPrefix, svc.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	pods, err := pdm.ListByService(svc.Meta.Namespace, svc.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get pod list by service id `%s` err: %s", logPrefix, svc.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Service().NewWithDeployment(svc, dl, pods).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:get write response err: %s", logPrefix, err.Error())
		return
	}
}

func ServiceCreateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation POST /namespace/{namespace}/service service serviceCreate
	//
	// Create new service
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: namespace
	//     in: path
	//     description: namespace id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_service_create"
	// responses:
	//   '200':
	//     description: Service was successfully created
	//     schema:
	//       "$ref": "#/definitions/views_service"
	//   '400':
	//     description: Name is already in use
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("%s:create:> create service in namespace `%s`", logPrefix, nid)

	var (
		opts = v1.Request().Service().Manifest()
	)

	// request body struct
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

	svc, e := service.Create(r.Context(), ns, opts)
	if e != nil {
		e.Http(w)
		return
	}

	response, err := v1.View().Service().NewWithDeployment(svc, nil, nil).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:update:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func ServiceUpdateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /namespace/{namespace}/service/{service} service serviceUpdate
	//
	// Update service
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: namespace
	//     in: path
	//     description: namespace id
	//     required: true
	//     type: string
	//   - name: service
	//     in: path
	//     description: service id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_service_update"
	// responses:
	//   '200':
	//     description: Service was successfully updated
	//     schema:
	//       "$ref": "#/definitions/views_service"
	//   '404':
	//     description: Namespace not found / Service not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["service"]

	log.V(logLevel).Debugf("%s:update:> update service `%s` in namespace `%s`", logPrefix, sid, nid)

	// request body struct
	opts := v1.Request().Service().Manifest()
	if e := opts.DecodeAndValidate(r.Body); e != nil {
		log.V(logLevel).Errorf("%s:update:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
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

	svc, e = service.Update(r.Context(), ns, svc, opts)
	if e != nil {
		e.Http(w)
		return
	}

	response, err := v1.View().Service().NewWithDeployment(svc, nil, nil).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:update:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func ServiceRemoveH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation DELETE /namespace/{namespace}/service/{service} service serviceRemove
	//
	// Remove service
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: namespace
	//     in: path
	//     description: namespace id
	//     required: true
	//     type: string
	//   - name: service
	//     in: path
	//     description: service id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Service was successfully removed
	//   '404':
	//     description: Namespace not found / Service not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["service"]

	log.V(logLevel).Debugf("%s:remove:> remove service `%s` from app `%s`", logPrefix, sid, nid)

	var (
		stg = envs.Get().GetStorage()
		nsm = distribution.NewNamespaceModel(r.Context(), stg)
		sm  = distribution.NewServiceModel(r.Context(), stg)
		rm  = distribution.NewRouteModel(r.Context(), stg)
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:remove:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	svc, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:> get service by name `%s` in namespace `%s` err: %s", logPrefix, sid, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("%s:remove:> service name `%s` in namespace `%s` not found", logPrefix, sid, ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	rl, err := rm.ListByNamespace(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:> get routes list in namespace `%s` err: %s", logPrefix, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	// check routes attached to routes
	for _, r := range rl.Items {
		for _, rule := range r.Spec.Rules {
			if rule.Service == svc.Meta.Name {
				log.V(logLevel).Errorf("%s:remove:> service used in route `%s` err: %s", logPrefix, r.Meta.Name, err.Error())
				errors.HTTP.BadRequest(w, errors.New(r.Meta.Name).Service().RouteBinded(r.Meta.Name).Error())
				return
			}
		}
	}

	if _, err := sm.Destroy(svc); err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove service err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.V(logLevel).Errorf("%s:remove:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func ServiceLogsH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/service/{service}/logs service serviceLogs
	//
	// Shows logs of the service
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: namespace
	//     in: path
	//     description: namespace id
	//     required: true
	//     type: string
	//   - name: service
	//     in: path
	//     description: service id
	//     required: true
	//     type: string
	//   - name: deployment
	//     in: query
	//     description: deployment id
	//     required: true
	//     type: string
	//   - name: pod
	//     in: query
	//     description: pod id
	//     required: true
	//     type: string
	//   - name: container
	//     in: query
	//     description: container id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Service logs received
	//   '404':
	//     description: Namespace not found / Service not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["service"]

	//did := r.URL.Query().Get("deployment")
	//pid := r.URL.Query().Get("pod")
	//cid := r.URL.Query().Get("container")

	log.V(logLevel).Debugf("%s:logs:> get logs for service `%s` in namespace `%s`", logPrefix, sid, nid)

	var (
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		em  = distribution.NewExporterModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:logs:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:logs:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	svc, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("%s:logs:> get service by name `%s` err: %s", logPrefix, sid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("%s:logs:> service name `%s` in namespace `%s` not found", logPrefix, sid, ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	el, err := em.List()
	if err != nil {
		log.V(logLevel).Errorf("%s:logs:> get exporters", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if len(el.Items) == 0 {
		log.V(logLevel).Errorf("%s:logs:>exporters not found", logPrefix, err.Error())
		errors.HTTP.NotFound(w)
		return
	}

	exp := new(types.Exporter)

	for _, e := range el.Items {
		if e.Status.Ready {
			exp = e
			break
		}
	}

	if exp == nil {
		log.V(logLevel).Errorf("%s:logs:> active exporters not found", logPrefix, err.Error())
		errors.HTTP.NotFound(w)
		return
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d/logs?kind=%s&selflink=%s", exp.Status.Http.IP, exp.Status.Http.Port, types.KindService, svc.SelfLink().String()), nil)
	if err != nil {
		log.V(logLevel).Errorf("%s:logs:> create http client err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.V(logLevel).Errorf("%s:logs:> get pod logs err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	var buffer = make([]byte, BUFFER_SIZE)

	for {

		select {
		case <-r.Context().Done():
			return
		default:

			n, err := res.Body.Read(buffer)
			if err != nil {

				if err == context.Canceled {
					log.V(logLevel).Debug("Stream is canceled")
					return
				}

				log.Errorf("Error read bytes from stream %s", err)
				return
			}

			_, err = func(p []byte) (n int, err error) {

				n, err = w.Write(p)
				if err != nil {
					log.Errorf("Error write bytes to stream %s", err)
					return n, err
				}

				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}

				return n, nil
			}(buffer[0:n])

			if err != nil {
				log.Errorf("Error written to stream %s", err)
				return
			}

			for i := 0; i < n; i++ {
				buffer[i] = 0
			}
		}
	}

}
