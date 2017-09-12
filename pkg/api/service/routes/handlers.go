//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package routes

import (
	"github.com/lastbackend/lastbackend/pkg/api/image"
	"github.com/lastbackend/lastbackend/pkg/api/app"
	"github.com/lastbackend/lastbackend/pkg/api/pod"
	"github.com/lastbackend/lastbackend/pkg/api/service"
	"github.com/lastbackend/lastbackend/pkg/api/service/routes/request"
	"github.com/lastbackend/lastbackend/pkg/api/service/views/v1"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logLevel = 2

func ServiceListH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		id  = utils.Vars(r)["app"]
	)

	log.V(logLevel).Debugf("Handler: Service: list services in app `%s`", id)

	ns := app.New(r.Context())
	item, err := ns.Get(id)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get app by name `%s` err: %s", id, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.V(logLevel).Warnf("Handler: Service: app `%s` not found", id)
		errors.New("app").NotFound().Http(w)
		return
	}

	s := service.New(r.Context(), item.Meta)
	items, err := s.List()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get service list in app `%s` err: %s", item.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewServiceList(items).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Service: write response err: %s", err.Error())
		return
	}
}

func ServiceInfoH(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		nid = utils.Vars(r)["app"]
		sid = utils.Vars(r)["service"]
	)

	log.V(logLevel).Debugf("Handler: Service: get service `%s` in app `%s`", sid, nid)

	ns := app.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get app by name `%s` err: %s", nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.V(logLevel).Warnf("Handler: Service: app `%s` not found", nid)
		errors.New("app").NotFound().Http(w)
		return
	}

	s := service.New(r.Context(), item.Meta)
	svc, err := s.Get(sid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get service by name `%s` in app `%s` err: %s", sid, item.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("Handler: Service: service `%s` in app `%s` not found", sid, item.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	response, err := v1.NewService(svc).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Service: write response err: %s", err.Error())
		return
	}
}

func ServiceCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		nid = utils.Vars(r)["app"]
	)

	log.V(logLevel).Debug("Handler: Service: create service")

	// request body struct
	rq := new(request.RequestServiceCreateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Service: validation incoming data err: %s", err.Err().Error())
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}
	ns := app.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get app by name `%s` err: %s", nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.V(logLevel).Warnf("Handler: Service: app `%s` not found", nid)
		errors.New("app").NotFound().Http(w)
		return
	}

	s := service.New(r.Context(), item.Meta)
	svc, err := s.Get(rq.Name)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get service by name `%s` in app `%s` err: %s", rq.Name, item.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc != nil {
		log.V(logLevel).Warnf("Handler: Service: service name `%s` in app `%s` not unique", rq.Name, item.Meta.Name)
		errors.New("service").NotUnique("name").Http(w)
		return
	}
	// Load template from registry
	if rq.Template != "" {
		// TODO: Send request for get template config from registry
		// TODO: Set service source with types.SourceTemplateType type field
		// TODO: Patch template config if need
		// TODO: Template provision
	}

	// If you are not using a template, then create a standard configuration template
	//if tpl == nil {
	// TODO: Generate default template for service
	//return
	//}

	// Patch config if exists custom configurations
	//if len(rq.Spec). != 0 {
	// TODO: If have custom config, then need patch this config
	//}

	if rq.Source.Hub != "" {

		img, err := image.Get(r.Context(), rq.Registry, rq.Source)
		if err != nil && err.Error() != store.ErrKeyNotFound {
			log.V(logLevel).Errorf("Handler: Service: get image err: %s", err.Error())
			return
		}

		if err != nil && err.Error() == store.ErrKeyNotFound {
			img, err = image.Create(r.Context(), rq.Registry, rq.Source)
			if err != nil {
				log.V(logLevel).Errorf("Handler: Service: create image err: %s", err.Error())
				errors.HTTP.InternalServerError(w)
				return
			}
		}

		rq.Spec.Image = &img.Meta.Name
	} else {
		rq.Spec.Image = &rq.Image
	}

	svc, err = s.Create(rq)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: create service err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewService(svc).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Service: write response err: %s", err.Error())
		return
	}
}

func ServiceUpdateH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		nid = utils.Vars(r)["app"]
		sid = utils.Vars(r)["service"]
	)

	log.V(logLevel).Debugf("Handler: Service: update service `%s` in app `%s`", sid, nid)

	// request body struct
	rq := new(request.RequestServiceUpdateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Service: validation incoming data err: %s", err.Err().Error())
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	ns := app.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get app by name `%s` err: %s", nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.V(logLevel).Warnf("Handler: Service: app `%s` not found", nid)
		errors.New("app").NotFound().Http(w)
		return
	}

	s := service.New(r.Context(), item.Meta)
	svc, err := s.Get(sid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get service by name `%s` in app `%s` err: %s", rq.Name, item.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("Handler: Service: service name `%s` in app `%s` not found", sid, item.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	if err = s.Update(svc, rq); err != nil {
		log.V(logLevel).Errorf("Handler: Service: update service err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewService(svc).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Service: write response err: %s", err.Error())
		return
	}
}

func ServiceRemoveH(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		nid = utils.Vars(r)["app"]
		sid = utils.Vars(r)["service"]
	)

	log.V(logLevel).Debugf("Handler: Service: remove service `%s` from app `%s`", sid, nid)

	ns := app.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get app by name `%s` err: %s", nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.V(logLevel).Warnf("Handler: Service: app `%s` not found", nid)
		errors.New("app").NotFound().Http(w)
		return
	}

	s := service.New(r.Context(), item.Meta)
	svc, err := s.Get(sid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get service by name `%s` in app `%s` err: %s", sid, item.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("Handler: Service: service name `%s` in app `%s` not found", sid, item.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	// TODO: remove all activity by service name

	if err := s.Remove(svc); err != nil {
		log.V(logLevel).Errorf("Handler: Service: remove service err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.V(logLevel).Errorf("Handler: Service: write response err: %s", err.Error())
		return
	}
}

func ServiceSpecCreateH(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		nid = utils.Vars(r)["app"]
		sid = utils.Vars(r)["service"]
	)

	log.V(logLevel).Debug("Handler: Service: create spec for service `%s` in app `%s`", sid, nid)

	// request body struct
	rq := new(request.RequestServiceSpecS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Service: validation incoming data err: %s", err.Err().Error())
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	ns := app.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get app by name `%s` err: %s", nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.V(logLevel).Warnf("Handler: Service: app `%s` not found", nid)
		errors.New("app").NotFound().Http(w)
		return
	}

	s := service.New(r.Context(), item.Meta)
	svc, err := s.Get(sid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get service `%s` in app by name `%s` err: %s", sid, nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("Handler: Service: service name `%s` in app `%s` not found", sid, item.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	if err = s.AddSpec(svc, rq); err != nil {
		log.V(logLevel).Errorf("Handler: Service: create spec for service `%s` in app by name `%s` err: %s", sid, nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewService(svc).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Service: write response err: %s", err.Error())
		return
	}
}

func ServiceSpecUpdateH(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		nid  = utils.Vars(r)["app"]
		sid  = utils.Vars(r)["service"]
		spid = utils.Vars(r)["spec"]
	)

	log.V(logLevel).Debug("Handler: Service: update spec for service `%s` in app `%s`", sid, nid)

	// request body struct
	rq := new(request.RequestServiceSpecS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Service: validation incoming data err: %s", err.Err().Error())
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	ns := app.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get app by name `%s` err: %s", nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.V(logLevel).Warnf("Handler: Service: app `%s` not found", nid)
		errors.New("app").NotFound().Http(w)
		return
	}

	s := service.New(r.Context(), item.Meta)
	svc, err := s.Get(sid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get service `%s` in app by name `%s` err: %s", sid, nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("Handler: Service: service name `%s` in app `%s` not found", sid, item.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	if err = s.SetSpec(svc, spid, rq); err != nil {
		log.V(logLevel).Errorf("Handler: Service: create spec for service `%s` in app by name `%s` err: %s", sid, nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewService(svc).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Service: write response err: %s", err.Error())
		return
	}
}

func ServiceSpecRemoveH(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		nid  = utils.Vars(r)["app"]
		sid  = utils.Vars(r)["service"]
		spid = utils.Vars(r)["spec"]
	)

	log.V(logLevel).Debug("Handler: Service: remove spec from service `%s` in app `%s`", sid, nid)

	// request body struct
	rq := new(request.RequestServiceUpdateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.V(logLevel).Errorf("Handler: Service: validation incoming data err: %s", err.Err().Error())
		errors.New("Invalid incoming data").Unknown().Http(w)
		return
	}

	ns := app.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get app by name `%s` err: %s", nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.V(logLevel).Warnf("Handler: Service: app `%s` not found", nid)
		errors.New("app").NotFound().Http(w)
		return
	}

	s := service.New(r.Context(), item.Meta)
	svc, err := s.Get(sid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get service `%s` in app by name `%s` err: %s", sid, nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("Handler: Service: service name `%s` in app `%s` not found", sid, item.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	if err = s.DelSpec(svc, spid); err != nil {
		log.V(logLevel).Errorf("Handler: Service: remove spec from service `%s` in app by name `%s` err: %s", sid, nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewService(svc).ToJson()
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: convert struct to json err: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.V(logLevel).Errorf("Handler: Service: write response err: %s", err.Error())
		return
	}
}

func ServiceLogsH(w http.ResponseWriter, r *http.Request) {
	var (
		nid      = utils.Vars(r)["app"]
		sid      = utils.Vars(r)["service"]
		pid      = r.URL.Query().Get("pod")
		cid      = r.URL.Query().Get("container")
		notify   = w.(http.CloseNotifier).CloseNotify()
		doneChan = make(chan bool, 1)
	)

	log.V(logLevel).Debug("Handler: Service: get logs for service `%s` in app `%s`", sid, nid)

	ns := app.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		log.V(logLevel).Errorf("Handler: Service: get app by name `%s` err: %s", nid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		log.V(logLevel).Warnf("Handler: Service: app `%s` not found", nid)
		errors.New("app").NotFound().Http(w)
		return
	}

	s := service.New(r.Context(), item.Meta)
	svc, err := s.Get(sid)
	if err != nil {
		log.Error("Error: find service by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc == nil {
		log.V(logLevel).Warnf("Handler: Service: service name `%s` in app `%s` not found", sid, item.Meta.Name)
		errors.New("service").NotFound().Http(w)
		return
	}

	go func() {
		<-notify
		log.Debug("HTTP connection just closed.")
		doneChan <- true
	}()

	if err := pod.Logs(r.Context(), item.Meta.Name, pid, cid, w, doneChan); err != nil {
		log.V(logLevel).Warnf("Handler: Service: get logs for service`%s` in app `%s` err: ", sid, item.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
}

func ServiceActivityListH(w http.ResponseWriter, r *http.Request) {

	var (
		nid = utils.Vars(r)["app"]
		sid = utils.Vars(r)["service"]
	)

	log.V(logLevel).Debug("Handler: Service: get activities for service `%s` in app `%s`", sid, nid)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`[]`)); err != nil {
		log.V(logLevel).Errorf("Handler: Service: write response err: %s", err.Error())
		return
	}
}
