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
	"github.com/lastbackend/lastbackend/internal/master/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/model"
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logJobPrefix               = "state:observer:job"
	defaultConcurrentTaskLimit = 1
)

// jobObserve manage handlers based on job state
func jobObserve(js *JobState, job *types.Job) (err error) {

	log.V(logLevel).Debugf("%s:> observe start: %s > %s", logJobPrefix, job.SelfLink(), job.Status.State)

	switch job.Status.State {
	// Check job created state triggers
	case types.StateCreated:
		err = handleJobStateCreated(js, job)
	// Check job provision state triggers
	case types.StateRunning:
		err = handleJobStateRunning(js, job)
	// Check job ready state triggers
	case types.StatePaused:
		err = handleJobStatePaused(js, job)
	// Check job error state triggers
	case types.StateError:
		err = handleJobStateError(js, job)
	// Run job destroy process
	case types.StateDestroy:
		err = handleJobStateDestroy(js, job)
	// Remove job from storage if it is already destroyed
	case types.StateDestroyed:
		err = handleJobStateDestroyed(js, job)
	}
	if err != nil {
		log.V(logLevel).Debugf("%s:observe:jobStateCreated:> handle job with state %s err:> %s", logPrefix, job.Status.State, err.Error())
		return err
	}

	if js.job == nil {
		return nil
	}

	log.V(logLevel).Debugf("%s:> observe finish: %s > %s", logJobPrefix, job.SelfLink(), job.Status.State)

	return nil
}

// handleJobStateCreated handles job created state
func handleJobStateCreated(js *JobState, job *types.Job) error {
	log.V(logLevel).Debugf("%s:> handleJobStateCreated: %s > %s", logJobPrefix, job.SelfLink(), job.Status.State)

	if js.provider != nil {
		go js.Provider()
	}

	if err := jobTaskProvision(js); err != nil {
		log.Errorf("%s:> job task provision err: %s", logPrefix, err.Error())
		return err
	}

	return nil
}

// handleJobStateRunning handles job provision state
func handleJobStateRunning(js *JobState, job *types.Job) error {
	log.V(logLevel).Debugf("%s:> handleJobStateRunning: %s > %s", logJobPrefix, job.SelfLink(), job.Status.State)
	return nil
}

// handleJobStatePaused handles job ready state
func handleJobStatePaused(js *JobState, job *types.Job) error {
	log.V(logLevel).Debugf("%s:> handleJobStatePaused: %s > %s", logJobPrefix, job.SelfLink(), job.Status.State)
	return nil
}

// handleJobStateError handles job error state
func handleJobStateError(js *JobState, job *types.Job) error {
	log.V(logLevel).Debugf("%s:> handleJobStateError: %s > %s", logJobPrefix, job.SelfLink(), job.Status.State)
	return nil
}

// handleJobStateDestroy handles job destroy state
func handleJobStateDestroy(js *JobState, job *types.Job) (err error) {

	log.V(logLevel).Debugf("%s:> handleJobStateDestroy: %s > %s", logJobPrefix, job.SelfLink(), job.Status.State)

	tm := model.NewTaskModel(context.Background(), envs.Get().GetStorage())

	if len(js.task.list) == 0 {

		jm := model.NewJobModel(context.Background(), envs.Get().GetStorage())

		job.Status.State = types.StateDestroyed
		job.Meta.Updated = time.Now()

		if err := jm.Set(job); err != nil {
			return err
		}

		return nil
	}

	for _, task := range js.task.list {

		if task.Status.State == types.StateDestroyed || task.Status.State == types.StateDestroy {
			continue
		}

		if err := tm.Destroy(task); err != nil {
			return err
		}

	}

	return nil
}

// handleJobStateDestroyed handles job destroyed state
func handleJobStateDestroyed(js *JobState, job *types.Job) (err error) {

	log.V(logLevel).Debugf("%s:> handleJobStateDestroyed: %s > %s", logJobPrefix, job.SelfLink(), job.Status.State)

	if len(js.task.list) > 0 {
		tm := model.NewTaskModel(context.Background(), envs.Get().GetStorage())
		for _, task := range js.task.list {

			if task.Status.State != types.StateDestroy {
				if err = tm.Destroy(task); err != nil {
					return err
				}
			}

			if task.Status.State == types.StateDestroyed {
				if err = tm.Remove(task); err != nil {
					return err
				}
			}

		}

		job.Status.State = types.StateDestroy
		job.Meta.Updated = time.Now()

		return nil
	}

	job.Status.State = types.StateDestroyed
	job.Meta.Updated = time.Now()

	jm := model.NewJobModel(context.Background(), envs.Get().GetStorage())
	nm := model.NewNamespaceModel(context.Background(), envs.Get().GetStorage())

	ns, err := nm.Get(job.Meta.Namespace)
	if err != nil {
		log.Errorf("%s:> namespace fetch err: %s", logJobPrefix, err.Error())
	}
	if ns != nil {
		ns.ReleaseResources(job.Spec.GetResourceRequest())

		if err := nm.Update(ns); err != nil {
			log.Errorf("%s:> namespace update err: %s", logJobPrefix, err.Error())
		}
	}

	if err = jm.Remove(job); err != nil {
		log.Errorf("%s:> job remove err: %s", logJobPrefix, err.Error())
		return err
	}

	js.job = nil
	return nil
}

// jobTaskProvision function handles all cases when task needs to be created or updated
func jobTaskProvision(js *JobState) error {

	// run task if no one task are currently running and there is at least one in queue
	var (
		limit = defaultConcurrentTaskLimit
		jm    = model.NewJobModel(context.Background(), envs.Get().GetStorage())
	)

	if len(js.task.queue) == 0 {
		log.Debugf("%s:jobTaskProvision:> there are no jobs in queue: %d", logJobPrefix, len(js.task.queue))
		if js.job.Status.State != types.StateWaiting {
			js.job.Status.State = types.StateWaiting
			if err := jm.Set(js.job); err != nil {
				log.Errorf("%s:jobTaskProvision:> set job to waiting state err: %s", logJobPrefix, err.Error())
				return err
			}
		}
		return nil
	}

	if js.job.Spec.Concurrency.Limit > 0 {
		limit = js.job.Spec.Concurrency.Limit
	}

	if len(js.task.active) >= limit {
		log.Debugf("%s:jobTaskProvision:> limit exceeded: %d >= %d", logJobPrefix, len(js.task.active), limit)
		return nil
	}

	// choose the older task task
	var t *types.Task
	for _, task := range js.task.queue {
		if t == nil {
			t = task
			continue
		}

		if task.Meta.Created.Before(t.Meta.Created) {
			t = task
		}
	}

	t.Status.State = types.StateProvision

	tm := model.NewTaskModel(context.Background(), envs.Get().GetStorage())
	if err := tm.Set(t); err != nil {
		log.Errorf("%s:jobTaskProvision:> set task to provision state err: %s", logJobPrefix, err.Error())
		return err
	}

	if js.job.Status.State != types.StateRunning {
		js.job.Status.State = types.StateRunning
		if err := jm.Set(js.job); err != nil {
			log.Errorf("%s:jobTaskProvision:> set job to running state err: %s", logJobPrefix, err.Error())
		}
	}

	return nil
}
