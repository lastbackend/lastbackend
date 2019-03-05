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
	logJobPrefix = "state:observer:job"
)

// jobObserve manage handlers based on job state
func jobObserve(ss *JobState, s *types.Job) error {

	log.V(logLevel).Debugf("%s:> observe start: %s > %s", logJobPrefix, s.SelfLink(), s.Status.State)

	switch s.Status.State {

	// Check job created state triggers
	case types.StateCreated:
		if err := handleJobStateCreated(ss, s); err != nil {
			log.V(logLevel).Debugf("%s:observe:jobStateCreated err:> %s", logPrefix, err.Error())
			return err
		}
		break

	// Check job provision state triggers
	case types.StateRunning:
		if err := handleJobStateRunning(ss, s); err != nil {
			log.V(logLevel).Debugf("%s:observe:jobStateProvision err:> %s", logPrefix, err.Error())
			return err
		}
		break

	// Check job ready state triggers
	case types.StatePaused:
		if err := handleJobStatePaused(ss, s); err != nil {
			log.V(logLevel).Debugf("%s:observe:jobStateReady err:> %s", logPrefix, err.Error())
			return err
		}
		break

	// Check job error state triggers
	case types.StateError:
		if err := handleJobStateError(ss, s); err != nil {
			log.V(logLevel).Debugf("%s:observe:jobStateError err:> %s", logPrefix, err.Error())
			return err
		}
		break

	// Run job destroy process
	case types.StateDestroy:
		if err := handleJobStateDestroy(ss, s); err != nil {
			log.V(logLevel).Debugf("%s:observe:jobStateDestroy err:> %s", logPrefix, err.Error())
			return err
		}
		break

	// Remove job from storage if it is already destroyed
	case types.StateDestroyed:
		if err := handleJobStateDestroyed(ss, s); err != nil {
			log.V(logLevel).Debugf("%s:observe:jobStateDestroyed err:> %s", logPrefix, err.Error())
			return err
		}
		break
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
		limit = 1
		jm    = distribution.NewJobModel(context.Background(), envs.Get().GetStorage())
	)

	if js.job.Spec.Concurrency.Limit > 0 {
		limit = js.job.Spec.Concurrency.Limit
	}

	if len(js.task.active) >= limit {
		log.Debugf("%s:> limit exceeded: %d >= %d", logJobPrefix, len(js.task.active), limit)
		return nil
	}

	if len(js.task.queue) <= 0 {

		log.Debugf("%s:> there are no jobs in queue: %d", logJobPrefix, len(js.task.queue))
		if js.job.Status.State != types.StateWaiting {
			js.job.Status.State = types.StateWaiting
			if err := jm.Set(js.job); err != nil {
				log.Errorf("%s:> set job to waiting state err: %s", logJobPrefix, err.Error())
			}
		}
		return nil
	}

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

	if t == nil {
		return nil
	}

	t.Status.State = types.StateProvision
	t.Status.Status = types.StateProvision

	tm := distribution.NewTaskModel(context.Background(), envs.Get().GetStorage())
	if err := tm.Set(t); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	js.task.active[t.SelfLink().String()] = t
	delete(js.task.queue, t.SelfLink().String())

	if js.job.Status.State != types.StateRunning {
		js.job.Status.State = types.StateRunning
		if err := jm.Set(js.job); err != nil {
			log.Errorf("%s:> set job to waiting state err: %s", logJobPrefix, err.Error())
		}
	}

	return nil
}
