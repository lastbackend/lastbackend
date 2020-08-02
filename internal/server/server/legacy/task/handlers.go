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

package task

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lastbackend/lastbackend/internal/server/server/legacy/middleware"
	h "github.com/lastbackend/lastbackend/internal/util/http"
	"github.com/lastbackend/lastbackend/tools/logger"
)

const (
	logPrefix  = "api:handler:job"
	BufferSize = 512
)

// Handler represent the http handler for task
type Handler struct {
}

// NewTaskHandler will initialize the task resources endpoint
func NewTaskHandler(r *mux.Router, mw middleware.Middleware) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	log.Infof("%s:> init task routes", logPrefix)

	handler := &Handler{
	}

	r.Handle("/namespace/{namespace}/job/{job}/task", h.Handle(mw.Authenticate(handler.TaskCreateH))).Methods(http.MethodPost)
	r.Handle("/namespace/{namespace}/job/{job}/task", h.Handle(mw.Authenticate(handler.TaskListH))).Methods(http.MethodGet)
	r.Handle("/namespace/{namespace}/job/{job}/task/{task}", h.Handle(mw.Authenticate(handler.TaskInfoH))).Methods(http.MethodGet)
	r.Handle("/namespace/{namespace}/job/{job}/task/{task}", h.Handle(mw.Authenticate(handler.TaskCancelH))).Methods(http.MethodPut)
	r.Handle("/namespace/{namespace}/job/{job}/task/{task}", h.Handle(mw.Authenticate(handler.TaskRemoveH))).Methods(http.MethodDelete)
	r.Handle("/namespace/{namespace}/job/{job}/task/{task}/logs", h.Handle(mw.Authenticate(handler.TaskLogsH))).Methods(http.MethodGet)
}

func (handler Handler) TaskListH(w http.ResponseWriter, r *http.Request) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	// swagger:operation GET /namespace/{namespace}/job job jobList
	//
	// Shows a list of tasks
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
	//     description: Task list response
	//     schema:
	//       "$ref": "#/definitions/views_job_list"
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	//nid := util.Vars(r)["namespace"]
	//jsl := util.Vars(r)["job"]
	//
	//log.Debugf("%s:list:> list tasks in %s", logPrefix, nid)
	//
	//var (
	//	stg = envs.Get().GetStorage()
	//	tm  = model.NewTaskModel(r.Context(), stg)
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//jb, e := job.Fetch(r.Context(), ns.Meta.Name, jsl)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//tasks, err := tm.ListByJob(ns.Meta.Name, jb.Meta.Name)
	//if err != nil {
	//	log.Errorf("%s:list:> get task list by job id `%s` err: %s", logPrefix, ns.Meta.Name, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//response, err := v1.View().Task().NewList(tasks).ToJson()
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

func (handler Handler) TaskInfoH(w http.ResponseWriter, r *http.Request) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	// swagger:operation GET /namespace/{namespace}/job/{job} job jobInfo
	//
	// Shows an info about job
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
	//   - name: job
	//     in: path
	//     description: job id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Job list response
	//     schema:
	//       "$ref": "#/definitions/views_job"
	//   '404':
	//     description: Namespace not found / Job not found
	//   '500':
	//     description: Internal server error

	//jib := util.Vars(r)["job"]
	//nid := util.Vars(r)["namespace"]
	//tsl := util.Vars(r)["task"]
	//
	//log.Debugf("%s:info:> get task `%s` in namespace `%s`", logPrefix, jib, nid)
	//
	//var (
	//	err error
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//jb, e := job.Fetch(r.Context(), ns.Meta.Name, jib)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//tk, e := task.Fetch(r.Context(), ns.Meta.Name, jb.Meta.Name, tsl)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//response, err := v1.View().Task().New(tk).ToJson()
	//if err != nil {
	//	log.Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:get write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) TaskCreateH(w http.ResponseWriter, r *http.Request) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	// swagger:operation POST /namespace/{namespace}/job job jobCreate
	//
	// Create new job
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
	//       "$ref": "#/definitions/request_job_create"
	// responses:
	//   '200':
	//     description: Job was successfully created
	//     schema:
	//       "$ref": "#/definitions/views_job"
	//   '400':
	//     description: Name is already in use
	//   '404':
	//     description: Namespace not found
	//   '500':
	//     description: Internal server error

	//nid := util.Vars(r)["namespace"]
	//jid := util.Vars(r)["job"]
	//
	//log.Debugf("%s:create:> create job in namespace `%s`", logPrefix, jid)
	//
	//var (
	//	opts = v1.Request().Task().Manifest()
	//)
	//
	//// request body struct
	//if err := opts.DecodeAndValidate(r.Body); err != nil {
	//	log.Errorf("%s:create:> validation incoming data err: %s", logPrefix, err.Err())
	//	err.Http(w)
	//	return
	//}
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//jb, e := job.Fetch(r.Context(), ns.Meta.Name, jid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//tk, e := task.Create(r.Context(), ns, jb, opts)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//response, err := v1.View().Task().New(tk).ToJson()
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

func (handler Handler) TaskCancelH(w http.ResponseWriter, r *http.Request) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	// swagger:operation DELETE /namespace/{namespace}/job/{job} job jobRemove
	//
	// Remove job
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
	//   - name: job
	//     in: path
	//     description: job id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Job was successfully removed
	//   '404':
	//     description: Namespace not found / Job not found
	//   '500':
	//     description: Internal server error

	//nid := util.Vars(r)["namespace"]
	//jid := util.Vars(r)["job"]
	//tsl := util.Vars(r)["task"]
	//
	//log.Debugf("%s:remove:> remove job `%s` from app `%s`", logPrefix, jid, nid)
	//
	//var (
	//	stg = envs.Get().GetStorage()
	//	tm  = model.NewTaskModel(r.Context(), stg)
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//jb, e := job.Fetch(r.Context(), ns.Meta.Name, jid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//tk, e := task.Fetch(r.Context(), ns.Meta.Name, jb.Meta.Name, tsl)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//if err := tm.Cancel(tk); err != nil {
	//	log.Errorf("%s:remove:> remove job err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//response, err := v1.View().Task().New(tk).ToJson()
	//if err != nil {
	//	log.Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:get write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) TaskRemoveH(w http.ResponseWriter, r *http.Request) {

	ctx := logger.NewContext(context.Background(), nil)
	log := logger.WithContext(ctx)

	// swagger:operation DELETE /namespace/{namespace}/job/{job} job jobRemove
	//
	// Remove job
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
	//   - name: job
	//     in: path
	//     description: job id
	//     required: true
	//     type: string
	// responses:
	//   '200':
	//     description: Job was successfully removed
	//   '404':
	//     description: Namespace not found / Job not found
	//   '500':
	//     description: Internal server error

	//nid := util.Vars(r)["namespace"]
	//jid := util.Vars(r)["job"]
	//tsl := util.Vars(r)["task"]
	//
	//log.Debugf("%s:remove:> remove job `%s` from app `%s`", logPrefix, jid, nid)
	//
	//var (
	//	stg = envs.Get().GetStorage()
	//	tm  = model.NewTaskModel(r.Context(), stg)
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//jb, e := job.Fetch(r.Context(), ns.Meta.Name, jid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//tk, e := task.Fetch(r.Context(), ns.Meta.Name, jb.Meta.Name, tsl)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//tk.Status.State = types.StateDestroy
	//tk.Spec.State.Destroy = true
	//
	//if err := tm.Set(tk); err != nil {
	//	log.Errorf("%s:remove:> remove job err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//response, err := v1.View().Task().New(tk).ToJson()
	//if err != nil {
	//	log.Errorf("%s:info:> convert struct to json err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}

	response := []byte{}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Errorf("%s:get write response err: %s", logPrefix, err.Error())
		return
	}
}

func (handler Handler) TaskLogsH(w http.ResponseWriter, r *http.Request) {

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
	//     description: Task logs received
	//   '404':
	//     description: Namespace not found / Task not found
	//   '500':
	//     description: Internal server error

	//nid := util.Vars(r)["namespace"]
	//jid := util.Vars(r)["job"]
	//tid := util.Vars(r)["task"]
	//pid := r.URL.Query().Get("pod")
	//cid := r.URL.Query().Get("container")
	//
	//log.Debugf("%s:logs:> get logs service `%s` in namespace `%s`", logPrefix, jid, nid)
	//
	//var (
	//	pm = model.NewPodModel(r.Context(), envs.Get().GetStorage())
	//	nm = model.NewNodeModel(r.Context(), envs.Get().GetStorage())
	//)
	//
	//ns, e := namespace.FetchFromRequest(r.Context(), nid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//jb, e := job.Fetch(r.Context(), ns.Meta.Name, jid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//tk, e := task.Fetch(r.Context(), ns.Meta.Name, jb.Meta.Name, tid)
	//if e != nil {
	//	e.Http(w)
	//	return
	//}
	//
	//sl, err := types.NewPodSelfLink(types.KindDeployment, tk.SelfLink().String(), pid)
	//if err != nil {
	//	log.Errorf("%s:logs:> pod selflink create err: %s", logPrefix, err.Error())
	//	errors.HTTP.BadRequest(w, "params")
	//	return
	//}
	//
	//pod, err := pm.Get(sl.String())
	//if err != nil {
	//	log.Errorf("%s:logs:> get pod by name` err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if pod == nil {
	//	log.Warnf("%s:logs:> pod `%s` not found", logPrefix, pid)
	//	errors.New("service").NotFound().Http(w)
	//	return
	//}
	//
	//node, err := nm.Get(pod.Meta.Node)
	//if err != nil {
	//	log.Errorf("%s:logs:> get node by name err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//if node == nil {
	//	log.Warnf("%s:logs:> node %s not found", logPrefix, pod.Meta.Node)
	//	errors.New("service").NotFound().Http(w)
	//	return
	//}
	//
	//req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d/pod/%s/%s/logs", node.Meta.ExternalIP, 2969, pod.SelfLink().String(), cid), nil)
	//if err != nil {
	//	log.Errorf("%s:logs:> create http client err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//res, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	log.Errorf("%s:logs:> get pod logs err: %s", logPrefix, err.Error())
	//	errors.HTTP.InternalServerError(w)
	//	return
	//}
	//
	//notify := w.(http.CloseNotifier).CloseNotify()
	//done := make(chan bool, 1)
	//
	//go func() {
	//	<-notify
	//	log.Debugf("%s:logs:> HTTP connection just closed.", logPrefix)
	//	done <- true
	//}()
	//
	//var buffer = make([]byte, BufferSize)
	//
	//for {
	//	select {
	//	case <-done:
	//		res.Body.Close()
	//		return
	//	default:
	//
	//		n, err := res.Body.Read(buffer)
	//		if err != nil {
	//
	//			if err == context.Canceled {
	//				log.Debug("Stream is canceled")
	//				return
	//			}
	//
	//			log.Errorf("Error read bytes from stream %s", err)
	//			return
	//		}
	//
	//		_, err = func(p []byte) (n int, err error) {
	//
	//			n, err = w.Write(p)
	//			if err != nil {
	//				log.Errorf("Error write bytes to stream %s", err)
	//				return n, err
	//			}
	//
	//			if f, ok := w.(http.Flusher); ok {
	//				f.Flush()
	//			}
	//
	//			return n, nil
	//		}(buffer[0:n])
	//
	//		if err != nil {
	//			log.Errorf("Error written to stream %s", err)
	//			return
	//		}
	//
	//		for i := 0; i < n; i++ {
	//			buffer[i] = 0
	//		}
	//	}
	//}

}
