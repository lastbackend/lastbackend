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
	"github.com/lastbackend/lastbackend/pkg/api/context"
	"github.com/lastbackend/lastbackend/pkg/api/image"
	"github.com/lastbackend/lastbackend/pkg/api/namespace"
	"github.com/lastbackend/lastbackend/pkg/api/service"
	"github.com/lastbackend/lastbackend/pkg/api/service/routes/request"
	"github.com/lastbackend/lastbackend/pkg/api/service/views/v1"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"net/http"
)

func ServiceListH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		log = context.Get().GetLogger()
		id  = utils.Vars(r)["namespace"]
	)

	log.Debug("List service handler")
	ns := namespace.New(r.Context())
	item, err := ns.Get(id)
	if err != nil {
		if item == nil {
			errors.New("namespace").NotFound().Http(w)
			return
		}

		log.Error("Error: find namespace by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	s := service.New(r.Context(), item.Meta)
	items, err := s.List()
	if err != nil {
		log.Error("Error: find service list by user", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewServiceList(items).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceInfoH(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		log = context.Get().GetLogger()
		nid = utils.Vars(r)["namespace"]
		sid = utils.Vars(r)["service"]
	)

	log.Debug("Get service handler")

	ns := namespace.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		if item == nil {
			errors.New("namespace").NotFound().Http(w)
			return
		}

		log.Errorf("Error: find namespace by name: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	s := service.New(r.Context(), item.Meta)
	svc, err := s.Get(sid)
	if err != nil {
		if svc == nil {
			errors.New("service").NotFound().Http(w)
			return
		}

		log.Errorf("Error: find service by name: %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewService(svc).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceWatchH(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		log = context.Get().GetLogger()
		nid = utils.Vars(r)["namespace"]
	)

	log.Debug("Get service handler")

	ns := namespace.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		if item == nil {
			errors.New("namespace").NotFound().Http(w)
			return
		}

		log.Error("Error: find namespace by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte("")); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceCreateH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		log = context.Get().GetLogger()
		nid = utils.Vars(r)["namespace"]
		sid = utils.Vars(r)["service"]
	)

	log.Debug("Create service handler")

	// request body struct
	rq := new(request.RequestServiceCreateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	ns := namespace.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		if item == nil {
			errors.New("namespace").NotFound().Http(w)
			return
		}

		log.Error("Error: find namespace by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	s := service.New(r.Context(), item.Meta)
	svc, err := s.Get(sid)
	if err != nil && err.Error() != store.ErrKeyNotFound {
		log.Error("Error: find service by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if svc != nil {
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
			return
		}

		if err != nil && err.Error() == store.ErrKeyNotFound {
			img, err = image.Create(r.Context(), rq.Registry, rq.Source)
			if err != nil {
				log.Error("Error: insert service to db", err)
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
		log.Error("Error: insert service to db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewService(svc).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceUpdateH(w http.ResponseWriter, r *http.Request) {

	var (
		err error
		log = context.Get().GetLogger()
		nid = utils.Vars(r)["namespace"]
		sid = utils.Vars(r)["service"]
	)

	log.Debug("Update service handler")

	// request body struct
	rq := new(request.RequestServiceUpdateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	ns := namespace.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		if item == nil {
			errors.New("namespace").NotFound().Http(w)
			return
		}

		log.Error("Error: find namespace by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	s := service.New(r.Context(), item.Meta)
	svc, err := s.Get(sid)
	if err != nil {
		if svc == nil {
			errors.New("service").NotFound().Http(w)
			return
		}

		log.Error("Error: find service by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	if err = s.Update(svc, rq); err != nil {
		log.Error("Error: update service error", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	// TODO: spec generate
	response, err := v1.NewService(svc).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceRemoveH(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		log = context.Get().GetLogger()

		nid = utils.Vars(r)["namespace"]
		sid = utils.Vars(r)["service"]
	)

	log.Info("Remove service")

	ns := namespace.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		log.Error("Error: find namespace by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		errors.New("namespace").NotFound().Http(w)
		return
	}

	s := service.New(r.Context(), item.Meta)
	svc, err := s.Get(sid)
	if err != nil {
		if svc == nil {
			errors.New("service").NotFound().Http(w)
			return
		}

		log.Error("Error: find service by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	// Todo: remove all activity by service name

	if err := s.Remove(svc); err != nil {
		log.Error("Error: remove service from db", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceActivityListH(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`[]`)); err != nil {
		context.Get().GetLogger().Error("Error: write response", err.Error())
		return
	}
}

func ServiceSpecCreateH(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		log = context.Get().GetLogger()
		nid = utils.Vars(r)["namespace"]
		sid = utils.Vars(r)["service"]
	)

	// request body struct
	rq := new(request.RequestServiceSpecCreateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	ns := namespace.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		log.Error("Error: find namespace by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		errors.New("namespace").NotFound().Http(w)
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
		errors.New("service").NotFound().Http(w)
		return
	}

	if err = s.AddSpec(svc, rq); err != nil {
		log.Error("Error: add spec to service error", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte{}); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceSpecUpdateH(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		log  = context.Get().GetLogger()
		nid  = utils.Vars(r)["namespace"]
		sid  = utils.Vars(r)["service"]
		spid = utils.Vars(r)["spec"]
	)

	// request body struct
	rq := new(request.RequestServiceSpecUpdateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	ns := namespace.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		log.Error("Error: find namespace by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		errors.New("namespace").NotFound().Http(w)
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
		errors.New("service").NotFound().Http(w)
		return
	}

	if err = s.SetSpec(svc, spid, rq); err != nil {
		log.Error("Error: update service spec error", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte{}); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceSpecRemoveH(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		log  = context.Get().GetLogger()
		nid  = utils.Vars(r)["namespace"]
		sid  = utils.Vars(r)["service"]
		spid = utils.Vars(r)["spec"]
	)

	// request body struct
	rq := new(request.RequestServiceUpdateS)
	if err := rq.DecodeAndValidate(r.Body); err != nil {
		log.Error("Error: validation incomming data", err)
		errors.New("Invalid incomming data").Unknown().Http(w)
		return
	}

	ns := namespace.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		log.Error("Error: find namespace by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		errors.New("namespace").NotFound().Http(w)
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
		errors.New("service").NotFound().Http(w)
		return
	}

	if err = s.DelSpec(svc, spid); err != nil {
		log.Error("Error: remove spec from service error", err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte{}); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func ServiceLogsH(w http.ResponseWriter, r *http.Request) {
	var (
		log      = context.Get().GetLogger()
		nid      = utils.Vars(r)["namespace"]
		sid      = utils.Vars(r)["service"]
		pid      = r.URL.Query().Get("pod")
		cid      = r.URL.Query().Get("container")
		notify   = w.(http.CloseNotifier).CloseNotify()
		doneChan = make(chan bool, 1)
	)

	log.Debug("Get service logs")

	ns := namespace.New(r.Context())
	item, err := ns.Get(nid)
	if err != nil {
		log.Error("Error: find namespace by name", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if item == nil {
		errors.New("namespace").NotFound().Http(w)
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
		errors.New("service").NotFound().Http(w)
		return
	}

	go func() {
		<-notify
		log.Debug("HTTP connection just closed.")
		doneChan <- true
	}()

	if err := service.Logs(r.Context(), item.Meta.Name, svc.Meta.Name, pid, cid, w, doneChan); err != nil {
		log.Errorf("Error: get service logs err %s", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
}
