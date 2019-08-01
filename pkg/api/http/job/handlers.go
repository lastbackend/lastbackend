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

package job

import (
	"context"
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/http/job/job"
	"github.com/lastbackend/lastbackend/pkg/api/http/namespace/namespace"
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

func JobListH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/job job jobList
	//
	// Shows a list of jobs
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

	log.V(logLevel).Debugf("%s:list:> list jobs in %s", logPrefix, nid)

	var (
		stg = envs.Get().GetStorage()
		jm  = distribution.NewJobModel(r.Context(), stg)
	)

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	jobs, err := jm.ListByNamespace(ns.Meta.Name)
	if err != nil {
		log.V(logLevel).Errorf("%s:list:> get job list in namespace `%s` err: %s", logPrefix, ns.Meta.Name, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.View().Job().NewList(jobs).ToJson()
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

func JobInfoH(w http.ResponseWriter, r *http.Request) {

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

	sid := utils.Vars(r)["job"]
	nid := utils.Vars(r)["namespace"]

	log.V(logLevel).Debugf("%s:info:> get job `%s` in namespace `%s`", logPrefix, sid, nid)

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	jb, e := job.Fetch(r.Context(), ns.Meta.Name, sid)
	if e != nil {
		e.Http(w)
		return
	}

	response, err := v1.View().Job().New(jb).ToJson()
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

func JobCreateH(w http.ResponseWriter, r *http.Request) {

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

	log.V(logLevel).Debugf("%s:create:> create job in namespace `%s`", logPrefix, nid)

	var (
		opts = v1.Request().Job().Manifest()
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

	jb, e := job.Create(r.Context(), ns, opts)
	if e != nil {
		e.Http(w)
		return
	}

	response, err := v1.View().Job().New(jb).ToJson()
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

func JobUpdateH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation PUT /namespace/{namespace}/job/{job} job jobUpdate
	//
	// Update job
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
	//   - name: body
	//     in: body
	//     required: true
	//     schema:
	//       "$ref": "#/definitions/request_job_update"
	// responses:
	//   '200':
	//     description: Job was successfully updated
	//     schema:
	//       "$ref": "#/definitions/views_job"
	//   '404':
	//     description: Namespace not found / Job not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]
	sid := utils.Vars(r)["job"]

	log.V(logLevel).Debugf("%s:update:> update job `%s` in namespace `%s`", logPrefix, sid, nid)

	// request body struct
	opts := v1.Request().Job().Manifest()
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

	jb, e := job.Fetch(r.Context(), ns.Meta.Name, sid)
	if e != nil {
		e.Http(w)
		return
	}

	jb, e = job.Update(r.Context(), ns, jb, opts)
	if e != nil {
		e.Http(w)
		return
	}

	response, err := v1.View().Job().New(jb).ToJson()
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

func JobRemoveH(w http.ResponseWriter, r *http.Request) {

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
	sid := utils.Vars(r)["job"]

	log.V(logLevel).Debugf("%s:remove:> remove job `%s` from app `%s`", logPrefix, sid, nid)

	var (
		stg = envs.Get().GetStorage()
		jm  = distribution.NewJobModel(r.Context(), stg)
	)

	ns, e := namespace.FetchFromRequest(r.Context(), nid)
	if e != nil {
		e.Http(w)
		return
	}

	jb, e := job.Fetch(r.Context(), ns.Meta.Name, sid)
	if e != nil {
		e.Http(w)
		return
	}

	jb, err := jm.Destroy(jb)
	if err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove job err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte{}); err != nil {
		log.V(logLevel).Errorf("%s:remove:> write response err: %s", logPrefix, err.Error())
		return
	}
}

func JobLogsH(w http.ResponseWriter, r *http.Request) {

	// swagger:operation GET /namespace/{namespace}/job/{job}/logs job jobLogs
	//
	// Shows logs of the job
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
	//     description: Applications logs received
	//   '404':
	//     description: Namespace not found / Service not found
	//   '500':
	//     description: Internal server error

	nid := utils.Vars(r)["namespace"]
	jid := utils.Vars(r)["job"]
	tid := utils.QueryString(r, "task")

	tail := utils.QueryInt(r, "tail")
	flw := utils.QueryBool(r, "follow")

	log.V(logLevel).Debugf("%s:logs:> get logs for job `%s` in namespace `%s`", logPrefix, jid, nid)

	var (
		nsm = distribution.NewNamespaceModel(r.Context(), envs.Get().GetStorage())
		jm  = distribution.NewJobModel(r.Context(), envs.Get().GetStorage())
		em  = distribution.NewExporterModel(r.Context(), envs.Get().GetStorage())
		tm  = distribution.NewTaskModel(r.Context(), envs.Get().GetStorage())
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

	jsl := types.NewJobSelfLink(ns.Meta.Name, jid)
	job, err := jm.Get(jsl.String())
	if err != nil {
		log.V(logLevel).Errorf("%s:logs:> get job by name `%s` err: %s", logPrefix, jid, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if job == nil {
		log.V(logLevel).Warnf("%s:logs:> job name `%s` in namespace `%s` not found", logPrefix, jid, ns.Meta.Name)
		errors.New("job").NotFound().Http(w)
		return
	}

	var task *types.Task
	if tid == types.EmptyString {
		tl, err := tm.ListByNamespace(ns.SelfLink().String())
		if err != nil {
			log.V(logLevel).Errorf("%s:logs:> get task list `%s` err: %s", logPrefix, err.Error())
			errors.HTTP.InternalServerError(w)
			return
		}

		for _, t := range tl.Items {
			if t.Status.State == types.StateRunning || t.Status.State == types.StateProvision {
				if task == nil {
					task = t
					continue
				}

				if task.Meta.Created.Before(t.Meta.Created) {
					task = t
				}
			}
		}

		if task == nil {
			for _, t := range tl.Items {

				if t.Status.State == types.StateWaiting {
					if task == nil {
						task = t
						continue
					}

					if task.Meta.Created.Before(t.Meta.Created) {
						task = t
					}
				}
			}
		}

		if task == nil {
			for _, t := range tl.Items {

				if t.Status.State == types.StateExited {
					if task == nil {
						task = t
						continue
					}

					if task.Meta.Created.Before(t.Meta.Created) {
						task = t
					}
				}
			}
		}

		if task == nil {
			errors.New("task").NotFound().Http(w)
			return
		}

	} else {
		tsl := types.NewTaskSelfLink(ns.Meta.Name, job.Meta.Name, tid)
		task, err = tm.Get(tsl.String())
		if err != nil {
			log.V(logLevel).Errorf("%s:logs:> get task by name `%s` err: %s", logPrefix, tsl.String(), err.Error())
			errors.HTTP.InternalServerError(w)
			return
		}
		if task == nil {
			log.V(logLevel).Warnf("%s:logs:> task name `%s` in namespace `%s` not found", logPrefix, tsl.String(), ns.Meta.Name)
			errors.New("task").NotFound().Http(w)
			return
		}

	}

	el, err := em.List()
	if err != nil {
		log.V(logLevel).Errorf("%s:logs:> get exporters", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}
	if len(el.Items) == 0 {
		log.V(logLevel).Errorf("%s:logs:>exporters not found", logPrefix)
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

	follow := "false"
	if flw && task.Status.State != types.StateExited {
		follow = "true"
	}

	cx, cancel := context.WithCancel(context.Background())
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%d/logs?kind=%s&selflink=%s&lines=%d&follow=%s",
		exp.Status.Http.IP, exp.Status.Http.Port, types.KindTask, task.SelfLink().String(), tail, follow), nil)
	if err != nil {
		log.V(logLevel).Errorf("%s:logs:> create http client err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	req.WithContext(cx)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", envs.Get().GetAccessToken()))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.V(logLevel).Errorf("%s:logs:> get pod logs err: %s", logPrefix, err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	defer cancel()

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
