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

package task

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/http/job/job"
	"github.com/lastbackend/lastbackend/pkg/api/http/namespace/namespace"
	"github.com/lastbackend/lastbackend/pkg/api/http/task/task"
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
	logPrefix   = "api:handler:job"
	BUFFER_SIZE = 512
)

func TaskListH(w http.ResponseWriter, r *http.Request) {

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

	nid := utils.Vars(r)["namespace"]
	jsl := utils.Vars(r)["job"]

	pflag := utils.QueryBool(r, "pods")

	log.V(logLevel).Debugf("%s:list:> list tasks in %s", logPrefix, nid)

	var (
		stg = envs.Get().GetStorage()
		tm  = distribution.NewTaskModel(r.Context(), stg)
	)

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	jb, e := job.Fetch(r.Context(), ns.Meta.Name, jsl)
	if e != nil {
		e.Http(w)
		return
	}

	tasks, err := tm.ListByJob(ns.Meta.Name, jb.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get task list by job id `%s` err: %s", logPrefix, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	var pods *types.PodList
	if pflag {

		if pods, err = distribution.NewPodModel(r.Context(), stg).ListByJob(ns.Meta.Name, jb.Meta.Name); err != nil {
			log.V(logLevel).Errorf("%s:list:> get pod list by job id `%s` err: %s", logPrefix, ns.Meta.Name, err.Error())
			errors.HTTP.InternalServerError(w)
			return
		}

	}

	response, err := v1.View().Task().NewList(tasks, pods).ToJson()
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

func TaskInfoH(w http.ResponseWriter, r *http.Request) {

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

	jib := utils.Vars(r)["job"]
	nid := utils.Vars(r)["namespace"]
	tsl := utils.Vars(r)["task"]

	pflag := utils.QueryBool(r, "pods")

	log.V(logLevel).Debugf("%s:info:> get task `%s` in namespace `%s`", logPrefix, jib, nid)

	var (
		err error
	)

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	jb, e := job.Fetch(r.Context(), ns.Meta.Name, jib)
	if e != nil {
		e.Http(w)
		return
	}

	tk, e := task.Fetch(r.Context(), ns.Meta.Name, jb.Meta.Name, tsl)
	if e != nil {
		e.Http(w)
		return
	}

	var pods *types.PodList
	if pflag {

		if pods, err = distribution.NewPodModel(r.Context(), envs.Get().GetStorage()).ListByTask(ns.Meta.Name, jb.Meta.Name, tk.Meta.Name); err != nil {
			log.V(logLevel).Errorf("%s:list:> get pod list by job id `%s` err: %s", logPrefix, ns.Meta.Name, err.Error())
			errors.HTTP.InternalServerError(w)
			return
		}

	}

	response, err := v1.View().Task().New(tk, pods).ToJson()
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

func TaskCreateH(w http.ResponseWriter, r *http.Request) {

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

	nid := utils.Vars(r)["namespace"]
	jid := utils.Vars(r)["job"]

	log.V(logLevel).Debugf("%s:create:> create job in namespace `%s`", logPrefix, jid)

	var (
		opts = v1.Request().Task().Manifest()
		stg  = envs.Get().GetStorage()
		tm   = distribution.NewTaskModel(r.Context(), stg)
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

	jb, e := job.Fetch(r.Context(), ns.Meta.Name, jid)
	if e != nil {
		e.Http(w)
		return
	}

	tk := new(types.Task)
	tk.Spec.Runtime = jb.Spec.Task.Runtime
	tk.Spec.Selector = jb.Spec.Task.Selector
	tk.Spec.Template = jb.Spec.Task.Template

	if opts != nil {
		if err := opts.SetTaskSpec(tk); err != nil {
			log.V(logLevel).Errorf("%s:list:> get pod list by job id `%s` err: %s", logPrefix, ns.Meta.Name, err.Error())
			errors.HTTP.BadParameter(w, "spec")
			return
		}
	}

	if _, err := tm.Create(tk); err != nil {
		log.V(logLevel).Errorf("%s:list:> get pod list by job id `%s` err: %s", logPrefix, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Task().New(tk, nil).ToJson()
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

func TaskCancelH(w http.ResponseWriter, r *http.Request) {

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

	nid := utils.Vars(r)["namespace"]
	jid := utils.Vars(r)["job"]
	tsl := utils.Vars(r)["task"]

	log.V(logLevel).Debugf("%s:remove:> remove job `%s` from app `%s`", logPrefix, jid, nid)

	var (
		stg = envs.Get().GetStorage()
		tm  = distribution.NewTaskModel(r.Context(), stg)
	)

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	jb, e := job.Fetch(r.Context(), ns.Meta.Name, jid)
	if e != nil {
		e.Http(w)
		return
	}

	tk, e := task.Fetch(r.Context(), ns.Meta.Name, jb.Meta.Name, tsl)
	if e != nil {
		e.Http(w)
		return
	}

	if err := tm.Cancel(tk); err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove job err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Task().New(tk, nil).ToJson()
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

func TaskLogsH(w http.ResponseWriter, r *http.Request) {

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

	nid := utils.Vars(r)["namespace"]
	jid := utils.Vars(r)["job"]
	tid := utils.Vars(r)["task"]
	pid := r.URL.Query().Get("pod")
	cid := r.URL.Query().Get("container")

	log.V(logLevel).Debugf("%s:logs:> get logs service `%s` in namespace `%s`", logPrefix, jid, nid)

	var (
		pm = distribution.NewPodModel(r.Context(), envs.Get().GetStorage())
		nm = distribution.NewNodeModel(r.Context(), envs.Get().GetStorage())
	)

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	jb, e := job.Fetch(r.Context(), ns.Meta.Name, jid)
	if e != nil {
		e.Http(w)
		return
	}

	tk, e := task.Fetch(r.Context(), ns.Meta.Name, jb.Meta.Name, tid)
	if e != nil {
		e.Http(w)
		return
	}

	pod, err := pm.Get(ns.Meta.Name, jb.Meta.Name, tk.Meta.Name, pid)
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

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d/pod/%s/%s/logs", node.Meta.ExternalIP, 2969, pod.Meta.SelfLink, cid), nil)
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
