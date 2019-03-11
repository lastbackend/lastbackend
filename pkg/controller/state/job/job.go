//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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

	"time"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const (
	logJobPrefix               = "state:observer:job"
	defaultConcurrentTaskLimit = 1
)

// jobObserve manage handlers based on job state
func jobObserve(ss *JobState, s *types.Job) (err error) {

	log.V(logLevel).Debugf("%s:> observe start: %s > %s", logJobPrefix, s.SelfLink(), s.Status.State)

	switch s.Status.State {
	// Check job created state triggers
	case types.StateCreated:
		err = handleJobStateCreated(ss, s)
	// Check job provision state triggers
	case types.StateRunning:
		err = handleJobStateRunning(ss, s)
	// Check job ready state triggers
	case types.StatePaused:
		err = handleJobStatePaused(ss, s)
	// Check job error state triggers
	case types.StateError:
		err = handleJobStateError(ss, s)
	// Run job destroy process
	case types.StateDestroy:
		err = handleJobStateDestroy(ss, s)
	// Remove job from storage if it is already destroyed
	case types.StateDestroyed:
		err = handleJobStateDestroyed(ss, s)
	}
	if err != nil {
		log.V(logLevel).Debugf("%s:observe:jobStateCreated:> handle job with state %s err:> %s", logPrefix, s.Status.State, err.Error())
		return err
	}

	if ss.job == nil {
		return nil
	}

	log.V(logLevel).Debugf("%s:> observe finish: %s > %s", logJobPrefix, s.SelfLink(), s.Status.State)

	return nil
}

// handleJobStateCreated handles job created state
func handleJobStateCreated(js *JobState, job *types.Job) error {
	log.V(logLevel).Debugf("%s:> handleJobStateCreated: %s > %s", logJobPrefix, job.SelfLink(), job.Status.State)
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

	dm := distribution.NewTaskModel(context.Background(), envs.Get().GetStorage())
	if len(js.task.list) == 0 {
		sm := distribution.NewJobModel(context.Background(), envs.Get().GetStorage())
		if err = sm.Remove(job); err != nil {
			log.Errorf("%s:> job remove err: %s", logJobPrefix, err.Error())
			return err
		}

		js.job = nil
		return nil
	}

	for _, d := range js.task.list {

		if d.Status.State == types.StateDestroyed {
			continue
		}

		if d.Status.State != types.StateDestroy {
			if err := dm.Destroy(d); err != nil {
				return err
			}
		}
	}

	if len(js.task.list) == 0 {
		job.Status.State = types.StateDestroyed
		job.Meta.Updated = time.Now()
	}

	return nil
}

// handleJobStateDestroyed handles job destroyed state
func handleJobStateDestroyed(js *JobState, job *types.Job) (err error) {

	log.V(logLevel).Debugf("%s:> handleJobStateDestroyed: %s > %s", logJobPrefix, job.SelfLink(), job.Status.State)

	job.Status.State = types.StateDestroy
	job.Meta.Updated = time.Now()

	if len(js.task.list) > 0 {
		dm := distribution.NewTaskModel(context.Background(), envs.Get().GetStorage())
		for _, d := range js.task.list {

			if d.Status.State == types.StateDestroyed {
				if err = dm.Remove(d); err != nil {
					return err
				}
			}

			if d.Status.State != types.StateDestroy {
				if err = dm.Destroy(d); err != nil {
					return err
				}
			}

		}

		job.Status.State = types.StateDestroy
		job.Meta.Updated = time.Now()

		return nil
	}

	sm := distribution.NewJobModel(context.Background(), envs.Get().GetStorage())
	nm := distribution.NewNamespaceModel(context.Background(), envs.Get().GetStorage())

	ns, err := nm.Get(job.Meta.Namespace)
	if err != nil {
		log.Errorf("%s:> namespece fetch err: %s", logJobPrefix, err.Error())
	}

	if ns != nil {
		ns.ReleaseResources(job.Spec.GetResourceRequest())

		if err := nm.Update(ns); err != nil {
			log.Errorf("%s:> namespece update err: %s", logJobPrefix, err.Error())
		}
	}

	if err = sm.Remove(job); err != nil {
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
		jm    = distribution.NewJobModel(context.Background(), envs.Get().GetStorage())
	)

	if len(js.task.queue) == 0 {
		log.Debugf("%s:jobTaskProvision:> there are no jobs in queue: %d", logJobPrefix, len(js.task.queue))
		if js.job.Status.State != types.StateWaiting {
			js.job.Status.State = types.StateWaiting
			if err := jm.Set(js.job); err != nil {
				log.Errorf("%s:jobTaskProvision:> set job to waiting state err: %s", logJobPrefix, err.Error())
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
	t.Status.Status = types.StateProvision

	tm := distribution.NewTaskModel(context.Background(), envs.Get().GetStorage())
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
