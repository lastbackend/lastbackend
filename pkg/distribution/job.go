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

package distribution

import (
	"context"
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/log"
	"time"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

const (
	logJobPrefix = "distribution:job"
)

// Job structure describe
type Job struct {
	context context.Context
	storage storage.Storage
}

// System - get job runtime
func (j *Job) Runtime() (*types.System, error) {

	log.V(logLevel).Debugf("%s:get:> get job runtime info", logJobPrefix)
	runtime, err := j.storage.Info(j.context, j.storage.Collection().Job(), "")
	if err != nil {
		log.V(logLevel).Errorf("%s:get:> get runtime info error: %s", logJobPrefix, err)
		return &runtime.System, err
	}
	return &runtime.System, nil
}

// Get job by selflink
func (j *Job) Get(selflink string) (*types.Job, error) {

	log.V(logLevel).Debugf("%s:get:> get by selflink %s", logJobPrefix, selflink)

	job := new(types.Job)
	err := j.storage.Get(j.context, j.storage.Collection().Job(), selflink, job, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> get job by selflink %s not found", logJobPrefix, selflink)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> get job by selflink %s error: %v", logJobPrefix, selflink, err)
		return nil, err
	}

	return job, nil
}

// ListByNamespace jobs
func (j *Job) ListByNamespace(namespace string) (*types.JobList, error) {
	log.V(logLevel).Debugf("%s:list:> by namespace %s", logJobPrefix, namespace)
	jobs := types.NewJobList()

	q := j.storage.Filter().Job().ByNamespace(namespace)
	err := j.storage.List(j.context, j.storage.Collection().Job(), q, jobs, nil)
	if err != nil {
		log.V(logLevel).Error("%s:list:> by namespace %s err: %v", logJobPrefix, namespace, err)
		return nil, err
	}

	log.V(logLevel).Debugf("%s:list:> by namespace %s result: %d", logJobPrefix, namespace, len(jobs.Items))

	return jobs, nil
}

// Create job
func (j *Job) Create(job *types.Job) (*types.Job, error) {

	job.Meta.Created = time.Now()
	if err := j.storage.Put(j.context, j.storage.Collection().Job(),
		job.SelfLink().String(), job, nil); err != nil {
		log.Errorf("%s:create:> job %s create err: %v", logJobPrefix, job.Meta.SelfLink, err)
		return nil, err
	}

	return job, nil
}

// Update job
func (j *Job) Set(job *types.Job) error {

	job.Meta.Updated = time.Now()
	log.V(logLevel).Debugf("%s:update:> update job %s", logJobPrefix, job.Meta.Name)

	if err := j.storage.Set(j.context, j.storage.Collection().Job(),
		job.SelfLink().String(), job, nil); err != nil {
		log.Errorf("%s:update:> update for job %s err: %v", logJobPrefix, job.Meta.Name, err)
		return err
	}

	return nil
}

// Pause job
func (j *Job) Pause(job *types.Job) error {

	log.V(logLevel).Debugf("%s:pause:> pause job %s", logJobPrefix, job.Meta.Name)

	// mark job for destroy
	job.Spec.Enabled = false
	// mark job for cancel
	job.Status.SetPaused()

	if err := j.storage.Set(j.context, j.storage.Collection().Job(),
		job.SelfLink().String(), job, nil); err != nil {
		log.V(logLevel).Debugf("%s:pause: pause job %s err: %v", logJobPrefix, job.Meta.Name, err)
		return err
	}

	return nil
}

// Start job
func (j *Job) Start(job *types.Job) error {

	log.V(logLevel).Debugf("%s:start:> start job %s", logJobPrefix, job.Meta.Name)

	// mark job for destroy
	job.Spec.Enabled = true
	// mark job for cancel
	job.Status.SetRunning()

	if err := j.storage.Set(j.context, j.storage.Collection().Job(),
		job.SelfLink().String(), job, nil); err != nil {
		log.V(logLevel).Debugf("%s:destroy:> destroy job %s err: %v", logJobPrefix, job.Meta.Name, err)
		return err
	}

	return nil
}

// Destroy job
func (j *Job) Destroy(job *types.Job) (*types.Job, error) {

	log.V(logLevel).Debugf("%s:destroy:> destroy job %s", logServicePrefix, job.SelfLink().String())

	job.Status.State = types.StateDestroy
	job.Spec.State.Destroy = true

	if err := j.storage.Set(j.context, j.storage.Collection().Job(),
		job.SelfLink().String(), job, nil); err != nil {
		log.V(logLevel).Errorf("%s:destroy:> destroy job err: %v", logServicePrefix, err)
		return nil, err
	}
	return job, nil
}

// Remove job
func (j *Job) Remove(job *types.Job) error {

	log.V(logLevel).Debugf("%s:remove:> remove job %s", logJobPrefix, job.Meta.Name)
	if err := j.storage.Del(j.context, j.storage.Collection().Job(),
		job.SelfLink().String()); err != nil {
		log.V(logLevel).Debugf("%s:remove:> remove job %s err: %v", logJobPrefix, job.Meta.Name, err)
		return err
	}

	return nil
}

// Watch job changes
func (j *Job) Watch(dt chan types.JobEvent, rev *int64) error {

	done := make(chan bool)
	watcher := storage.NewWatcher()

	log.V(logLevel).Debugf("%s:watch:> watch jobs", logJobPrefix)

	go func() {
		for {
			select {
			case <-j.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := types.JobEvent{}
				res.Action = e.Action
				res.Name = e.Name

				job := new(types.Job)

				if err := json.Unmarshal(e.Data.([]byte), job); err != nil {
					log.Errorf("%s:> parse data err: %v", logJobPrefix, err)
					continue
				}

				res.Data = job

				dt <- res
			}
		}
	}()

	opts := storage.GetOpts()
	opts.Rev = rev
	if err := j.storage.Watch(j.context, j.storage.Collection().Job(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewJobModel(ctx context.Context, stg storage.Storage) *Job {
	return &Job{ctx, stg}
}
