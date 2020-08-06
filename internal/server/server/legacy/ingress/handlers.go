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

package ingress

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/server/server/middleware"
	h "github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/tools/logger"
)

const (
	logPrefix = "api:handler:ingress"
)

// Handler represent the http handler for ingress
type Handler struct {
}

// NewIngressHandler will initialize the ingress resources endpoint
func NewIngressHandler(r *mux.Router, mw middleware.Middleware) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:> init ingress routes", logPrefix)

	handler := &Handler{}

	r.Handle("/ingress", h.Handle(mw.Authenticate(handler.IngressListH))).Methods(http.MethodGet)
	r.Handle("/ingress/{ingress}", h.Handle(mw.Authenticate(handler.IngressInfoH))).Methods(http.MethodGet)
	r.Handle("/ingress/{ingress}", h.Handle(mw.Authenticate(handler.IngressConnectH))).Methods(http.MethodPut)
	r.Handle("/ingress/{ingress}", h.Handle(mw.Authenticate(handler.IngressRemoveH))).Methods(http.MethodDelete)
	r.Handle("/ingress/{ingress}/status", h.Handle(mw.Authenticate(handler.IngressSetStatusH))).Methods(http.MethodPut)
}

func (handler Handler) IngressInfoH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/ingress/{ingress} ingress ingressInfo
	//
	// Shows an ingress info
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: ingress
	//     in: path
	//     description: ingress id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Ingress response
	//     schema:
	//       "$ref": "#/definitions/views_ingress_list"
	//   '404':
	//     description: Ingress not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:info:> get ingress", logPrefix)
	//
	//var (
	//	im  = model.NewIngressModel(r.Context(), envs.Get().GetStorage())
	//	nid = util.Vars(r)["ingress"]
	//)
	//
	//n, err := im.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:info:> get ingress err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if n == nil {
	//	log.Warnf("%s:info:> ingress `%s` not found", logPrefix, nid)
	//	errors.New("ingress").NotFound().Http(w)
	//	return
	//}
	//
	//response, err := v1.View().Ingress().New(n).ToJson()
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

func (handler Handler) IngressListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/ingress ingress ingressList
	//
	// Shows an ingress list
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: Ingress list response
	//     schema:
	//       "$ref": "#/definitions/views_ingress_list"
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:list:> get ingresss list", logPrefix)
	//
	//var (
	//	im = model.NewIngressModel(r.Context(), envs.Get().GetStorage())
	//)
	//
	//ingresss, err := im.List()
	//if err != nil {
	//	log.Errorf("%s:list:> get ingresss list err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//response, err := v1.View().Ingress().NewList(ingresss).ToJson()
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

func (handler Handler) IngressConnectH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /cluster/ingress/{ingress} ingress ingressInfo
	//
	// Shows an ingress info
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: ingress
	//     in: path
	//     description: ingress id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Ingress response
	//     schema:
	//       "$ref": "#/definitions/views_ingress_list"
	//   '404':
	//     description: Ingress not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:info:> get ingress", logPrefix)
	//
	//var (
	//	im    = model.NewIngressModel(r.Context(), envs.Get().GetStorage())
	//	sn    = model.NewNetworkModel(r.Context(), envs.Get().GetStorage())
	//	nid   = util.Vars(r)["ingress"]
	//	cache = envs.Get().GetCache().Ingress()
	//)
	//
	//// request body struct
	//opts := new(request.IngressConnectOptions)
	//if err := opts.DecodeAndValidate(r.Body); err != nil {
	//	log.Errorf("%s:setstatus:> validation incoming data", logPrefix, err.Err())
	//	err.Http(w)
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
	//ing, err := im.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:info:> get ingress err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if ing == nil {
	//	log.Warnf("%s:info:> ingress `%s` not found", logPrefix, nid)
	//
	//	ingress := new(types.Ingress)
	//	ingress.Meta.SetDefault()
	//	ingress.Meta.Name = opts.Info.Hostname
	//	ingress.Meta.SelfLink = *types.NewIngressSelfLink(opts.Info.Hostname)
	//	ingress.Status.Ready = opts.Status.Ready
	//
	//	im.Put(ingress)
	//
	//	if snet == nil {
	//		if _, err := sn.SubnetPut(opts.Info.Hostname, opts.Network.SubnetSpec); err != nil {
	//			log.Errorf("%s:connect:> snet put error: %s", logPrefix, err.Error())
	//			errors.HTTP.InternalServerError(w)
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
	//ing.Status.Ready = opts.Status.Ready
	//if err := im.Set(ing); err != nil {
	//	log.Errorf("%s:connect:> get ingress set err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//if snet == nil {
	//	if _, err := sn.SubnetPut(opts.Info.Hostname, opts.Network.SubnetSpec); err != nil {
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
	//cache.Clear(nid)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:connect:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) IngressSetStatusH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /ingress/{ingress}/status ingress ingressSetStatus
	//
	// Set ingress status
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: ingress
	//     in: path
	//     description: ingress id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_ingress_status"
	// responses:
	//   '200':
	//     description: Successfully set ingress status
	//   '400':
	//     description: Bad request
	//   '404':
	//     description: Node not found / Pod not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:setstatus:> ingress set state", logPrefix)
	//
	//var (
	//	im  = model.NewIngressModel(r.Context(), envs.Get().GetStorage())
	//	rm  = model.NewRouteModel(r.Context(), envs.Get().GetStorage())
	//	nid = util.Vars(r)["ingress"]
	//)
	//
	//// request body struct
	//opts := new(request.IngressStatusOptions)
	//if err := opts.DecodeAndValidate(r.Body); err != nil {
	//	log.Errorf("%s:setstatus:> validation incoming data", logPrefix, err.Err())
	//	err.Http(w)
	//	return
	//}
	//
	//ingress, err := im.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:setstatus:> get ingresss list err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if ingress == nil {
	//	log.Warnf("%s:setstatus:> update ingress `%s` not found", logPrefix, nid)
	//	errors.New("ingress").NotFound().Http(w)
	//	return
	//}
	//
	//ingress.Status.Ready = opts.Status.Ready
	//ingress.Status.Online = true
	//
	//if err := im.Set(ingress); err != nil {
	//	log.Errorf("%s:setstatus:> set status err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//for r, s := range opts.Routes {
	//
	//	log.Debugf("set route status: %s> %s", r, s.State)
	//
	//	keys := strings.Split(r, ":")
	//	if len(keys) != 2 {
	//		log.Errorf("%s:setroutestatus:> invalid route selflink err: %s", logPrefix, r)
	//		errors.HTTP.BadRequest(w)
	//		return
	//	}
	//
	//	route, err := rm.Get(keys[0], keys[1])
	//	if err != nil {
	//		log.Errorf("%s:setroutestatus:> route found err: %s", logPrefix, r)
	//		errors.HTTP.InternalServerError(w)
	//		return
	//	}
	//	if route == nil {
	//		log.Warnf("%s:setroutestatus:> route not found `%s` not found", logPrefix, r)
	//		if err := rm.ManifestDel(nid, r); err != nil {
	//			if !errors.Storage().IsErrEntityNotFound(err) {
	//				log.Warnf("%s:setroutestatus:> route manifest del err `%s` ", logPrefix, err.Error())
	//				continue
	//			}
	//		}
	//		continue
	//	}
	//
	//	route.Status.State = s.State
	//	route.Status.Message = s.Message
	//
	//	if _, err := rm.Set(route); err != nil {
	//		log.Errorf("%s:setroutestatus:> update route err: %s", logPrefix, err.Error())
	//		errors.HTTP.InternalServerError(w)
	//		return
	//	}
	//}
	//
	//spec, err := getIngressManifest(r.Context(), ingress)
	//if err != nil {
	//	errors.HTTP.InternalServerError(w)
	//}
	//
	//response, err := v1.View().Ingress().NewManifest(spec).ToJson()
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

func (handler Handler) IngressRemoveH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation DELETE /cluster/ingress/{ingress} ingress ingressRemove
	//
	// Remove ingress
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: ingress
	//     in: path
	//     description: ingress id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Ingress removed
	//   '404':
	//     description: Ingress not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(r.Context(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:remove:>_ create ingress", logPrefix)
	//
	//var (
	//	nm  = model.NewIngressModel(r.Context(), envs.Get().GetStorage())
	//	nid = util.Vars(r)["ingress"]
	//)
	//
	//n, err := nm.Get(nid)
	//if err != nil {
	//	log.Errorf("%s:remove:>_ remove ingress err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//if n == nil {
	//	log.Warnf("%s:remove:>_ remove ingress `%s` not found", logPrefix, nid)
	//	errors.New("ingress").NotFound().Http(w)
	//	return
	//}
	//
	//if err := nm.Remove(n); err != nil {
	//	log.Errorf("%s:remove:>_ remove ingress err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:remove:> write response err: %s", logPrefix, err.Error())
		return
	}
}

//func getIngressManifest(ctx context.Context, ing *types.Ingress) (*types.IngressManifest, error) {
//
//	var (
//		cache = envs.Get().GetCache().Ingress()
//		spec  = cache.Get(ing.SelfLink().String())
//		stg   = envs.Get().GetStorage()
//		ns    = model.NewNetworkModel(ctx, stg)
//		em    = model.NewEndpointModel(ctx, stg)
//		rm    = model.NewRouteModel(ctx, stg)
//	)
//
//	if spec == nil {
//		spec = new(types.IngressManifest)
//		spec.Meta.Initial = true
//
//		spec.Resolvers = cache.GetResolvers()
//		spec.Routes = make(map[string]*types.RouteManifest, 0)
//		routes, err := rm.ManifestMap(ing.SelfLink().String())
//		if err != nil {
//			log.Errorf("%s:getmanifest:> get endpoint manifests for node err: %s", logPrefix, err.Error())
//			return spec, err
//		}
//
//		if routes != nil {
//			for s, r := range routes.Items {
//				spec.Routes[s] = r
//			}
//		}
//
//		endpoints, err := em.ManifestMap()
//		if err != nil {
//			log.Errorf("%s:getmanifest:> get endpoint manifests for node err: %s", logPrefix, err.Error())
//			return spec, err
//		}
//		spec.Endpoints = endpoints.Items
//
//		subnets, err := ns.SubnetManifestMap()
//		if err != nil {
//			log.Errorf("%s:getmanifest:> get endpoint manifests for ingress err: %s", logPrefix, err.Error())
//			return spec, err
//		}
//		spec.Network = subnets.Items
//	}
//
//	cache.Flush(ing.SelfLink().String())
//	return spec, nil
//
//}
