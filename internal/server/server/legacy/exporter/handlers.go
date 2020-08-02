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

package exporter

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/server/server/legacy/middleware"
	h "github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/tools/logger"
)

const (
	logPrefix = "api:handler:exporter"
)

// Handler represent the http handler for exporter
type Handler struct {
}

// NewExporterHandler will initialize the exporter resources endpoint
func NewExporterHandler(r *mux.Router, mw middleware.Middleware) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:> init exporter routes", logPrefix)

	handler := &Handler{
	}

	r.Handle("/exporter", h.Handle(mw.Authenticate(handler.ExporterListH))).Methods(http.MethodGet)
	r.Handle("/exporter/{exporter}", h.Handle(mw.Authenticate(handler.ExporterInfoH))).Methods(http.MethodGet)
	r.Handle("/exporter/{exporter}", h.Handle(mw.Authenticate(handler.ExporterConnectH))).Methods(http.MethodPut)
	r.Handle("/exporter/{exporter}", h.Handle(mw.Authenticate(handler.ExporterRemoveH))).Methods(http.MethodDelete)
	r.Handle("/exporter/{exporter}/status", h.Handle(mw.Authenticate(handler.ExporterSetStatusH))).Methods(http.MethodPut)
}

func (handler Handler) ExporterInfoH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/exporter/{exporter} exporter exporterInfo
	//
	// Shows an exporter info
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: exporter
	//     in: path
	//     description: exporter id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Exporter response
	//     schema:
	//       "$ref": "#/definitions/views_exporter_list"
	//   '404':
	//     description: Exporter not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:info:> get exporter", logPrefix)
	//
	//var (
	//	im  = model.NewExporterModel(r.Context(), envs.Get().GetStorage())
	//	nid = util.Vars(r)["exporter"]
	//)
	//
	//n, err := im.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:info:> get exporter err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if n == nil {
	//	log.Warnf("%s:info:> exporter `%s` not found", logPrefix, nid)
	//	errors.New("exporter").NotFound().Http(w)
	//	return
	//}
	//
	//response, err := v1.View().Exporter().New(n).ToJson()
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

func (handler Handler) ExporterListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/exporter exporter exporterList
	//
	// Shows an exporter list
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: Exporter list response
	//     schema:
	//       "$ref": "#/definitions/views_exporter_list"
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:list:> get exporters list", logPrefix)
	//
	//var (
	//	im = model.NewExporterModel(r.Context(), envs.Get().GetStorage())
	//)
	//
	//exporters, err := im.List()
	//if err != nil {
	//	log.Errorf("%s:list:> get exporters list err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//response, err := v1.View().Exporter().NewList(exporters).ToJson()
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

func (handler Handler) ExporterConnectH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/exporter/{exporter} exporter exporterInfo
	//
	// Shows an exporter info
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: exporter
	//     in: path
	//     description: exporter id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Exporter response
	//     schema:
	//       "$ref": "#/definitions/views_exporter_list"
	//   '404':
	//     description: Exporter not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:info:> exporter connect", logPrefix)
	//
	//var (
	//	stg   = envs.Get().GetStorage()
	//	dm    = model.NewExporterModel(r.Context(), stg)
	//	sn    = model.NewNetworkModel(r.Context(), stg)
	//	nid   = util.Vars(r)["exporter"]
	//	cache = envs.Get().GetCache().Exporter()
	//)
	//
	//// request body struct
	//opts := new(request.ExporterConnectOptions)
	//if err := opts.DecodeAndValidate(r.Body); err != nil {
	//	log.Errorf("%s:setstatus:> validation incoming data", logPrefix, err.Err())
	//	err.Http(w)
	//	return
	//}
	//
	//dvc, err := dm.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:info:> get exporter err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//snet, err := sn.SubnetGet(opts.Network.CIDR)
	//if err != nil {
	//	log.Errorf("%s:connect:> get subnet err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//if dvc == nil {
	//	log.Debugf("%s:info:> create new exporter `%s`", logPrefix, nid)
	//
	//	exporter := new(types.Exporter)
	//	exporter.Meta.SetDefault()
	//	exporter.Meta.Name = opts.Info.Hostname
	//
	//	exporter.Meta.SelfLink = *types.NewExporterSelfLink(opts.Info.Hostname)
	//
	//	exporter.Status.Listener.Port = opts.Status.Listener.Port
	//	exporter.Status.Listener.IP = opts.Status.Listener.IP
	//	exporter.Status.Http.Port = opts.Status.Http.Port
	//	exporter.Status.Http.IP = opts.Status.Http.IP
	//	exporter.Status.Ready = opts.Status.Ready
	//
	//	if err := dm.Put(exporter); err != nil {
	//		log.Errorf("can not add exporter: %s", err.Error())
	//		errors.HTTP.InternalServerError(w)
	//		return
	//	}
	//
	//	if snet == nil {
	//		if _, err := sn.SubnetPut(exporter.SelfLink().String(), opts.Network.SubnetSpec); err != nil {
	//			log.Errorf("%s:connect:> snet put error: %s", logPrefix, err.Error())
	//			errors.HTTP.InternalServerError(w)
	//			return
	//		}
	//	}
	//
	//	w.WriteHeader(http.StatusOK)
	//	if _, err := w.Write([]byte{}); err != nil {
	//		log.Errorf("%s:connect:> write response err: %s", logPrefix, err.Error())
	//		return
	//	}
	//
	//	return
	//}
	//
	//dvc.Status.Ready = opts.Status.Ready
	//dvc.Status.Listener.Port = opts.Status.Listener.Port
	//dvc.Status.Listener.IP = opts.Status.Listener.IP
	//dvc.Status.Http.Port = opts.Status.Http.Port
	//dvc.Status.Http.IP = opts.Status.Http.IP
	//
	//if err := dm.Set(dvc); err != nil {
	//	log.Errorf("%s:connect:> get exporter set err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//if snet == nil {
	//	if _, err := sn.SubnetPut(dvc.SelfLink().String(), opts.Network.SubnetSpec); err != nil {
	//		log.Errorf("%s:connect:> snet put error: %s", logPrefix, err.Error())
	//		errors.HTTP.InternalServerError(w)
	//	}
	//} else {
	//	if !sn.SubnetEqual(snet, opts.Network.SubnetSpec) {
	//		snet.Spec = opts.Network.SubnetSpec
	//		if err := sn.SubnetSet(snet); err != nil {
	//			log.Errorf("%s:connect:> get subnet set err: %s", logPrefix, err.Error())
	//			errors.HTTP.InternalServerError(w)
	//			return
	//		}
	//	}
	//}
	//
	//cache.Clear(dvc.SelfLink().String())

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:connect:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) ExporterSetStatusH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /exporter/{exporter}/status exporter exporterSetStatus
	//
	// Set exporter status
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: exporter
	//     in: path
	//     description: exporter id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_exporter_status"
	// responses:
	//   '200':
	//     description: Successfully set exporter status
	//   '400':
	//     description: Bad request
	//   '404':
	//     description: Node not found / Pod not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:setstatus:> exporter set state", logPrefix)
	//
	//var (
	//	dm  = model.NewExporterModel(r.Context(), envs.Get().GetStorage())
	//	nid = util.Vars(r)["exporter"]
	//)
	//
	//// request body struct
	//opts := new(request.ExporterStatusOptions)
	//if err := opts.DecodeAndValidate(r.Body); err != nil {
	//	log.Errorf("%s:setstatus:> validation incoming data", logPrefix, err.Err())
	//	err.Http(w)
	//	return
	//}
	//
	//exporter, err := dm.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:setstatus:> get exporters list err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if exporter == nil {
	//	log.Warnf("%s:setstatus:> update exporter `%s` not found", logPrefix, nid)
	//	errors.New("exporter").NotFound().Http(w)
	//	return
	//}
	//
	//exporter.Status.Ready = opts.Ready
	//exporter.Status.Listener.Port = opts.Listener.Port
	//exporter.Status.Listener.IP = opts.Listener.IP
	//exporter.Status.Http.Port = opts.Http.Port
	//exporter.Status.Http.IP = opts.Http.IP
	//
	//exporter.Status.Online = true
	//
	//if err := dm.Set(exporter); err != nil {
	//	log.Errorf("%s:setstatus:> set status err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//spec, err := getCollectorManifest(r.Context(), exporter)
	//if err != nil {
	//	errors.HTTP.InternalServerError(w)
	//}
	//
	//response, err := v1.View().Exporter().NewManifest(spec).ToJson()
	//if err != nil {
	//	log.Errorf("%s:getspec:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:setstatus:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) ExporterRemoveH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation DELETE /cluster/exporter/{exporter} exporter exporterRemove
	//
	// Remove exporter
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: exporter
	//     in: path
	//     description: exporter id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Exporter removed
	//   '404':
	//     description: Exporter not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:remove:>_ create exporter", logPrefix)
	//
	//var (
	//	nm  = model.NewExporterModel(r.Context(), envs.Get().GetStorage())
	//	nid = util.Vars(r)["exporter"]
	//)
	//
	//n, err := nm.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:remove:>_ remove exporter err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//if n == nil {
	//	log.Warnf("%s:remove:>_ remove exporter `%s` not found", logPrefix, nid)
	//	errors.New("exporter").NotFound().Http(w)
	//	return
	//}
	//
	//if err := nm.Remove(n); err != nil {
	//	log.Errorf("%s:remove:>_ remove exporter err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:remove:>_ write response err: %s", logPrefix, err.Error())
		return
	}
}

//func getCollectorManifest(ctx context.Context, dns *types.Exporter) (*types.ExporterManifest, error) {
//
//	var (
//		cache = envs.Get().GetCache().Exporter()
//		spec  = cache.Get(dns.SelfLink().String())
//	)
//
//	cache.Flush(dns.SelfLink().String())
//	return spec, nil
//
//}
