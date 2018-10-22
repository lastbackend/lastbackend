//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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
	"net/http"
	"strings"

	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
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
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		dm  = distribution.NewDeploymentModel(r.Context(), envs.Get().GetStorage())
		pdm = distribution.NewPodModel(r.Context(), envs.Get().GetStorage())
	)

	ns, err := nsm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:info:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	srv, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get service by name `%s` in namespace `%s` err: %s", logPrefix, sid, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if srv == nil {
		log.V(logLevel).Warnf("%s:info:> service `%s` in namespace `%s` not found", logPrefix, sid, ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	dl, err := dm.ListByService(srv.Meta.Namespace, srv.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get pod list by service id `%s` err: %s", logPrefix, srv.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	pods, err := pdm.ListByService(srv.Meta.Namespace, srv.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:info:> get pod list by service id `%s` err: %s", logPrefix, srv.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Service().NewWithDeployment(srv, dl, pods).ToJson()
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
		nm   = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm   = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		opts = v1.Request().Service().Manifest()
	)

	// request body struct
	if err := opts.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("%s:create:> validation incoming data err: %s", logPrefix, err.Err())
		err.Http(w)
		return
	}

	ns, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:create:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	if opts.Meta.Name != nil {

		srv, err := sm.Get(ns.Meta.Name, *opts.Meta.Name)
		if err != nil {
			log.V(logLevel).Errorf("%s:create:> get service by name `%s` in namespace `%s` err: %s", logPrefix, opts.Meta.Name, ns.Meta.Name, err.Error())
			errors.HTTP.InternalServerError(w)
			return
		}

		if srv != nil {
			log.V(logLevel).Warnf("%s:create:> service name `%s` in namespace `%s` not unique", logPrefix, opts.Meta.Name, ns.Meta.Name)
			errors.New("service").NotUnique("name").Http(w)
			return
		}
	}

	svc := new(types.Service)
	opts.SetServiceMeta(svc)
	svc.Meta.Namespace = ns.Meta.Name
	svc.Meta.Endpoint = fmt.Sprintf("%s.%s", strings.ToLower(svc.Meta.Name), ns.Meta.Endpoint)

	opts.SetServiceSpec(svc)

	srv, err := sm.Create(ns, svc)
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> create service err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Service().NewWithDeployment(srv, nil, nil).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> convert struct to json err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("%s:create:> write response err: %s", logPrefix, err.Error())
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

	var (
		nm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
	)

	// request body struct
	opts := v1.Request().Service().Manifest()
	if e := opts.DecodeAndValidate(r.Body); e != nil {
		log.V(logLevel).Errorf("%s:update:> validation incoming data err: %s", logPrefix, e.Err())
		e.Http(w)
		return
	}

	ns, err := nm.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> get namespace", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if ns == nil {
		err := errors.New("namespace not found")
		log.V(logLevel).Errorf("%s:update:> get namespace", logPrefix, err.Error())
		errors.New("namespace").NotFound().Http(w)
		return
	}

	svc, err := sm.Get(ns.Meta.Name, sid)
	if err != nil {
		log.V(logLevel).Errorf("%s: get service by name` err: %s", logPrefix, sid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("%s:update:> service name `%s` in namespace `%s` not found", logPrefix, sid, ns.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	opts.SetServiceMeta(svc)
	svc.Meta.Endpoint = fmt.Sprintf("%s.%s", strings.ToLower(svc.Meta.Name), ns.Meta.Endpoint)
	opts.SetServiceSpec(svc)

	srv, err := sm.Update(svc)
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> update service err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Service().NewWithDeployment(srv, nil, nil).ToJson()
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
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
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
	did := r.URL.Query().Get("deployment")
	pid := r.URL.Query().Get("pod")
	cid := r.URL.Query().Get("container")

	log.V(logLevel).Debugf("%s:logs:> get logs service `%s` in namespace `%s`", logPrefix, sid, nid)

	var (
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		sm  = distribution.NewServiceModel(r.Context(), envs.Get().GetStorage())
		pm  = distribution.NewPodModel(r.Context(), envs.Get().GetStorage())
		dm  = distribution.NewDeploymentModel(r.Context(), envs.Get().GetStorage())
		nm  = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
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

	deployment, err := dm.Get(ns.Meta.Name, svc.Meta.Name, did)
	if err != nil {
		log.V(logLevel).Errorf("%s:logs:> get deployment by name err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if deployment == nil {
		log.V(logLevel).Warnf("%s:logs:> deployment `%s` not found", logPrefix, pid)
		errors.New("service").NotFound().Http(w)
		return
	}

	pod, err := pm.Get(ns.Meta.Name, svc.Meta.Name, did, pid)
	if err != nil {
		log.V(logLevel).Errorf("%s:logs:> get pod by name` err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if pod == nil {
		log.V(logLevel).Warnf("%s:logs:> pod `%s` not found", logPrefix, pid)
		errors.New("service").NotFound().Http(w)
		return
	}

	node, err := nm.Get(pod.Meta.Node)
	if err != nil {
		log.V(logLevel).Errorf("%s:logs:> get node by name err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if node == nil {
		log.V(logLevel).Warnf("%s:logs:> node %s not found", logPrefix, pod.Meta.Node)
		errors.New("service").NotFound().Http(w)
		return
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d/pod/%s/%s/logs", node.Meta.InternalIP, 2969, pod.Meta.SelfLink, cid), nil)
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

	notify := w.(http.CloseNotifier).CloseNotify()
	done := make(chan bool, 1)

	go func() {
		<-notify
		log.V(logLevel).Debugf("%s:logs:> HTTP connection just closed.", logPrefix)
		done <- true
	}()

	var buffer = make([]byte, BUFFER_SIZE)

	for {
		select {
		case <-done:
			res.Body.Close()
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
