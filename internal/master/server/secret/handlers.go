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

package secret

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/master/server/middleware"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	h "github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/tools/logger"
)

const (
	logPrefix = "api:handler:secret"
)

// Handler represent the http handler for secret
type Handler struct {
	Vault *models.Vault
}

type Config struct {
	Vault *models.Vault
}

// NewSecretHandler will initialize the secret resources endpoint
func NewSecretHandler(r *mux.Router, mw middleware.Middleware, cfg Config) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:> init secret routes", logPrefix)

	handler := &Handler{
		Vault: cfg.Vault,
	}

	r.Handle("/namespace/{namespace}/secret", h.Handle(mw.Authenticate(handler.SecretCreateH))).Methods(http.MethodPost)
	r.Handle("/namespace/{namespace}/secret", h.Handle(mw.Authenticate(handler.SecretListH))).Methods(http.MethodGet)
	r.Handle("/namespace/{namespace}/secret/{secret}", h.Handle(mw.Authenticate(handler.SecretGetH))).Methods(http.MethodGet)
	r.Handle("/namespace/{namespace}/secret/{secret}", h.Handle(mw.Authenticate(handler.SecretUpdateH))).Methods(http.MethodPut)
	r.Handle("/namespace/{namespace}/secret/{secret}", h.Handle(mw.Authenticate(handler.SecretRemoveH))).Methods(http.MethodDelete)
}

func (handler Handler) SecretGetH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/secret secret secretList
	//
	// Shows a list of secrets
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
	//     description: Secret list response
	//     schema:
	//       "$ref": "#/definitions/views_secret_list"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:get:> get secret", logPrefix)
	//
	//var (
	//	sid  = util.Vars(r)["secret"]
	//	nid  = util.Vars(r)["namespace"]
	//	rm   = model.NewSecretModel(r.Context(), envs.Get().GetStorage())
	//	item *types.Secret
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//parts := strings.SplitN(sid, ":", 2)
	//
	//switch len(parts) {
	//case 1:
	//	var err error
	//	item, err = rm.Get(ns.Meta.Name, sid)
	//	if err != nil {
	//		log.Errorf("%s:get:> find secret list err: %s", logPrefix, err.Error())
	//		errors.HTTP.InternalServerError(w)
	//		return
	//	}
	//case 2:
	//
	//	if parts[0] != "vault" {
	//		log.Errorf("%s:get:> invalid secret name: %s", logPrefix, sid)
	//		errors.HTTP.InternalServerError(w)
	//		return
	//	}
	//
	//	cx, cancel := context.WithCancel(context.Background())
	//
	//	vault := envs.Get().GetVault()
	//	if vault == nil {
	//		log.Warnf("%s:get:> vault not found", logPrefix)
	//		errors.New("vault").NotFound().Http(w)
	//		return
	//	}
	//
	//	url := fmt.Sprintf("%s/vault?secret=%s&namespace=%s", vault.Endpoint, parts[1], ns.SelfLink().String())
	//	req, err := http.NewRequest(http.MethodGet, url, nil)
	//	if err != nil {
	//		log.Errorf("%s:secret:> create http client err: %s", logPrefix, err.Error())
	//		errors.HTTP.InternalServerError(w)
	//		return
	//	}
	//	req.Header.Set("x-lastbackend-token", vault.Token)
	//
	//	req.WithContext(cx)
	//	res, err := http.DefaultClient.Do(req)
	//	if err != nil {
	//		log.Errorf("%s:secret:> get secret err: %s", logPrefix, err.Error())
	//		errors.HTTP.InternalServerError(w)
	//		return
	//	}
	//
	//	body, err := ioutil.ReadAll(res.Body)
	//	if err != nil {
	//		log.Errorf("%s:secret:> read secret err: %s", logPrefix, err.Error())
	//		errors.HTTP.InternalServerError(w)
	//		return
	//	}
	//
	//	sv := views.SecretView{}
	//	item, err = sv.Parse(body)
	//	if err != nil {
	//		log.Errorf("%s:secret:> parse secret err: %s", logPrefix, err.Error())
	//		errors.HTTP.InternalServerError(w)
	//		return
	//	}
	//
	//	defer cancel()
	//}
	//
	//if item == nil {
	//	log.Warnf("%s:get:> secret `%s` not found", logPrefix, sid)
	//	errors.New("secret").NotFound().Http(w)
	//	return
	//}
	//
	//response, err := v1.View().Secret().New(item).ToJson()
	//if err != nil {
	//	log.Errorf("%s:get:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:get:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) SecretListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/secret secret secretList
	//
	// Shows a list of secrets
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
	//     description: Secret list response
	//     schema:
	//       "$ref": "#/definitions/views_secret_list"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:list:> get secrets list", logPrefix)
	//
	//var (
	//	nid = util.Vars(r)["namespace"]
	//	rm  = model.NewSecretModel(r.Context(), envs.Get().GetStorage())
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//items, err := rm.List(ns.Meta.Name)
	//if err != nil {
	//	log.Errorf("%s:list:> find secret list err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//response, err := v1.View().Secret().NewList(items).ToJson()
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

func (handler Handler) SecretCreateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation POST /namespace/{namespace}/secret secret secretCreate
	//
	// Create secret
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
	//       "$ref": "#/definitions/request_secret_create"
	// responses:
	//   '200':
	//     description: Secret was successfully created
	//     schema:
	//       "$ref": "#/definitions/views_secret"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//log.Debugf("%s:create:> create secret", logPrefix)
	//
	//var (
	//	nid  = util.Vars(r)["namespace"]
	//	opts = v1.Request().Secret().Manifest()
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//// request body struct
	//e = opts.DecodeAndValidate(r.Body)
	//if e != nil {
	//	log.Errorf("%s:create:> validation incoming data err: %s", logPrefix, e.Err())
	//	e.Http(w)
	//	return
	//}
	//
	//sct, e := secret.Create(r.Context(), ns, opts)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//response, err := v1.View().Secret().New(sct).ToJson()
	//if err != nil {
	//	log.Errorf("%s:create:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:create:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) SecretUpdateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /namespace/{namespace}/secret/{secret} secret secretUpdate
	//
	// Create secret
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
	//   - name: secret
	//     in: path
	//     description: secret id
	//     required: true
	//     type: string
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_secret_update"
	// responses:
	//   '200':
	//     description: Secret was successfully updated
	//     schema:
	//       "$ref": "#/definitions/views_secret"
	//   '404':
	//     description: Namespace not found / Secret not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//sid := util.Vars(r)["secret"]
	//
	//log.Debugf("%s:update:> update secret `%s`", logPrefix, sid)
	//
	//var (
	//	nid  = util.Vars(r)["namespace"]
	//	opts = v1.Request().Secret().Manifest()
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//// request body struct
	//e = opts.DecodeAndValidate(r.Body)
	//if e != nil {
	//	log.Errorf("%s:update:> validation incoming data err: %s", logPrefix, e.Err())
	//	e.Http(w)
	//	return
	//}
	//
	//sct, e := secret.Fetch(r.Context(), ns.Meta.Name, sid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//sct, e = secret.Update(r.Context(), ns, sct, opts)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//response, err := v1.View().Secret().New(sct).ToJson()
	//if err != nil {
	//	log.Errorf("%s:update:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:update:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) SecretRemoveH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation DELETE /namespace/{namespace}/secret/{secret} secret secretRemove
	//
	// Remove secret
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
	//   - name: secret
	//     in: path
	//     description: secret id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Secret was successfully removed
	//   '404':
	//     description: Namespace not found / Secret not found
	//   '500':
	//     description: Internal server error

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	//sid := util.Vars(r)["secret"]
	//
	//log.Debugf("%s:remove:> remove secret %s", logPrefix, sid)
	//
	//var (
	//	nid = util.Vars(r)["namespace"]
	//	sm  = model.NewSecretModel(r.Context(), envs.Get().GetStorage())
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//ss, err := sm.Get(ns.Meta.Name, sid)
	//if err != nil {
	//	log.Errorf("%s:remove:> get secret by id `%s` err: %s", logPrefix, sid, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if ss == nil {
	//	log.Warnf("%s:remove:> secret `%s` not found", logPrefix, sid)
	//	errors.New("secret").NotFound().Http(w)
	//	return
	//}
	//
	//err = sm.Remove(ss)
	//if err != nil {
	//	log.Errorf("%s:remove:> remove secret `%s` err: %s", logPrefix, sid, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.Errorf("%s:remove:> write response err: %s", logPrefix, err.Error())
		return
	}
}
