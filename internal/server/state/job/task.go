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

package job

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/internal/util/generator"
	"github.com/lastbackend/lastbackend/tools/log"
)

const logTaskPrefix = "state:observer:task"

func taskObserve(js *JobState, task *models.Task) (err error) {

	log.Debugf("%s:> observe start: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	switch task.Status.State {
	case models.StateCreated:
		err = handleTaskStateCreated(js, task)
	case models.StateQueued:
		err = handleTaskStateQueued(js, task)
	case models.StateProvision:
		err = handleTaskStateProvision(js, task)
	case models.StateRunning:
		err = handleTaskStateRunning(js, task)
	case models.StateError:
		err = handleTaskStateError(js, task)
	case models.StateCanceled:
		err = handleTaskStateCanceled(js, task)
	case models.StateExited:
		err = handleTaskStateExited(js, task)
	case models.StateDestroy:
		err = handleTaskStateDestroy(js, task)
	case models.StateDestroyed:
		err = handleTaskStateDestroyed(js, task)
	}
	if err != nil {
		task.Status.State = models.StateError
		task.Status.Error = true
		task.Status.Message = err.Error()
		if err := handleTaskStateError(js, task); err != nil {
			log.Errorf("%s:> handle task state %s error err: %s", logTaskPrefix, task.Status.State, err.Error())
			return err
		}
		return nil
	}

	log.Debugf("%s:> observe finish: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	return nil
}

func handleTaskStateCreated(js *JobState, task *models.Task) error {

	log.Debugf("%s:handleTaskStateCreated:> try to handle task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	if err := taskCheckSelectors(js, task); err != nil {
		task.Status.State = models.StateError
		task.Status.Error = true
		task.Status.Message = err.Error()
		tm := service.NewTaskModel(context.Background(), js.storage)
		if err := tm.Set(task); err != nil {
			log.Errorf("%s:handleTaskStateCreated:> handle task create, deps update: %s, err: %s", logTaskPrefix, task.SelfLink(), err.Error())
			return err
		}
		return nil
	}

	if err := taskQueue(js, task); err != nil {
		log.Errorf("%s:handleTaskStateCreated:> move task %s to queue err: %s", logTaskPrefix, task.Meta.Name, err.Error())
		return err
	}

	return nil
}

func handleTaskStateQueued(js *JobState, task *models.Task) error {

	log.Debugf("%s:handleTaskStateQueued:> task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	if err := taskQueue(js, task); err != nil {
		log.Errorf("%s:handleTaskStateProvision:> task queued err: %s", logTaskPrefix, err.Error())
		return err
	}

	return nil
}

func handleTaskStateProvision(js *JobState, task *models.Task) error {

	log.Debugf("%s:handleTaskStateProvision:> task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	// check pods are created and state is normal state
	if err := taskProvision(js, task); err != nil {
		log.Errorf("%s:handleTaskStateProvision:> task provision err: %s", logTaskPrefix, err.Error())
		return err
	}

	return nil
}

func handleTaskStateRunning(_ *JobState, task *models.Task) error {

	log.Debugf("%s:handleTaskStateRunning:> task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)
	// there nothing need to be done

	return nil
}

func handleTaskStateError(js *JobState, task *models.Task) error {

	log.Debugf("%s:handleTaskStateError:> task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	// finish task and destroy it
	if err := taskFinish(js, task); err != nil {
		log.Errorf("%s:handleTaskStateError:> task finish err: %s", logTaskPrefix, err.Error())
		return err
	}

	return nil
}

func handleTaskStateCanceled(js *JobState, task *models.Task) error {

	log.Debugf("%s:handleTaskStateCanceled:> task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	// finish task and destroy it
	if err := taskFinish(js, task); err != nil {
		log.Errorf("%s:handleTaskStateCanceled:> task finish err: %s", logTaskPrefix, err.Error())
		return err
	}

	return nil
}

func handleTaskStateExited(js *JobState, task *models.Task) error {

	log.Debugf("%s:handleTaskStateExited:>: task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)
	// finish task and destroy it
	if err := taskFinish(js, task); err != nil {
		log.Errorf("%s:handleTaskStateExited:> task finish err: %s", logTaskPrefix, err.Error())
		return err
	}

	return nil
}

func handleTaskStateDestroy(js *JobState, task *models.Task) error {

	log.Debugf("%s:handleTaskStateDestroy:> task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	if err := taskDestroy(js, task); err != nil {
		log.Errorf("%s:handleTaskStateDestroy:> task destroy err: %s", logTaskPrefix, err.Error())
		return err
	}

	return nil
}

func handleTaskStateDestroyed(js *JobState, task *models.Task) error {

	log.Debugf("%s:handleTaskStateDestroyed:> task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	if _, ok := js.pod.list[task.SelfLink().String()]; ok {
		if err := taskDestroy(js, task); err != nil {
			log.Errorf("%s:handleTaskStateDestroyed:> task destroy err: %s", logTaskPrefix, err.Error())
			return err
		}
		return nil
	}
	if err := taskRemove(js.storage, task); err != nil {
		log.Errorf("%s:handleTaskStateDestroyed:> remove task err: %s", logTaskPrefix, err.Error())
		return err
	}
	js.DelTask(task)
	return nil
}

// taskCheckSelectors function - handles provided selectors to match nodes
func taskCheckSelectors(js *JobState, task *models.Task) (err error) {

	var (
		ctx = context.Background()
		vm  = service.NewVolumeModel(ctx, js.storage)
		vc  = make(map[string]string, 0)
	)

	for _, v := range task.Spec.Template.Volumes {
		if v.Volume.Name != models.EmptyString {
			vc[v.Volume.Name] = v.Name
		}
	}

	if len(vc) > 0 {

		var node string

		vl, err := vm.ListByNamespace(task.Meta.Namespace)
		if err != nil {
			log.Errorf("%s:check_selectors:> create task, volume list err: %s", logPrefix, err.Error())
			return err
		}

		for name := range vc {

			var f = false

			for _, v := range vl.Items {

				if v.Meta.Name != name {
					continue
				}

				f = true

				if v.Status.State != models.StateReady {
					log.Errorf("%s:check_selectors:> create task err: volume is not ready yet: %s", logPrefix, v.Meta.Name)
					return errors.New(v.Meta.Name).Volume().NotReady(v.Meta.Name)
				}

				if v.Meta.Node == models.EmptyString {
					log.Errorf("%s:check_selectors:> create task err: volume is not provisioned yet: %s", logPrefix, v.Meta.Name)
					return errors.New(v.Meta.Name).Volume().NotProvisioned(v.Meta.Name)
				}

				if node == models.EmptyString {
					node = v.Meta.Node
				} else {
					if node != v.Meta.Node {
						return errors.New(v.Meta.Name).Volume().DifferentNodes()
					}
				}
			}

			if !f {
				log.Errorf("%s:check_selectors:> create task err: volume is not found: %s", logPrefix, name)
				return errors.New(name).Volume().NotFound(name)
			}
		}

		if node != models.EmptyString {

			if task.Spec.Selector.Node != models.EmptyString {
				if task.Spec.Selector.Node != node {
					return errors.New("spec.selector.node not matched with attached volumes")
				}

				return nil
			}

			task.Spec.Selector.Node = node
		}

	}

	return nil
}

// taskCreate - create a new task from current job
// usually used by cron or other time repeatable jobs
func taskCreate(stg storage.IStorage, job *models.Job, mf *models.TaskManifest) (*models.Task, error) {

	tm := service.NewTaskModel(context.Background(), stg)

	task := new(models.Task)
	task.Meta.SetDefault()
	task.Meta.Namespace = job.Meta.Namespace
	task.Meta.Job = job.SelfLink().String()

	if mf != nil {
		mf.SetTaskMeta(task)
	}

	if task.Meta.Name == models.EmptyString {
		name := strings.Split(generator.GetUUIDV4(), "-")[4][5:]
		task.Meta.Name = name
	}

	task.Meta.SelfLink = *models.NewTaskSelfLink(job.Meta.Namespace, job.Meta.Name, task.Meta.Name)

	task.Spec.Runtime = job.Spec.Task.Runtime
	task.Spec.Selector = job.Spec.Task.Selector
	task.Spec.Template = job.Spec.Task.Template

	if mf != nil {
		if err := mf.SetTaskSpec(task); err != nil {
			log.Errorf("%s:taskCreate:> set task spec err: %v", logTaskPrefix, err)
			return nil, err
		}
	}

	d, err := tm.Create(task)
	if err != nil {
		log.Errorf("%s:taskCreate:> create task err: %v", logTaskPrefix, err)
		return nil, err
	}

	return d, nil
}

func taskQueue(js *JobState, task *models.Task) error {

	log.Debugf("%s:taskQueue:> move task %s to queue", logTaskPrefix, task.Meta.Name)

	// set a task as queued if it is not right now
	// ==============================================

	if task.Status.State != models.StateQueued {

		task.Status.State = models.StateQueued

		tm := service.NewTaskModel(context.Background(), js.storage)
		if err := tm.Set(task); err != nil {
			log.Errorf("%s:taskQueue:> set task err: %s", logTaskPrefix, err.Error())
			return err
		}

		return nil
	}

	// if the task is set to the queued state,
	// then add the task to queue and set to
	// provision state if available by limits
	// ==============================================

	js.task.queue[task.SelfLink().String()] = task

	if err := jobTaskProvision(js); err != nil {
		log.Errorf("%s:taskQueue:> job task queue pop err: %s", logTaskPrefix, err.Error())
		return err
	}

	return nil
}

// taskProvision - handles task provision logic
// based on current task state and current pod list of provided task
func taskProvision(js *JobState, task *models.Task) (err error) {

	log.Debugf("%s:taskProvision:> set task %s as provision", logTaskPrefix, task.Meta.Name)

	var (
		t  = task.Meta.Updated
		pm = service.NewPodModel(context.Background(), js.storage)
	)

	// set a task as provision if it is not right now
	// ==============================================

	if task.Status.State != models.StateProvision {
		task.Status.State = models.StateProvision
		task.Meta.Updated = time.Now()

		if err := taskUpdate(js.storage, task, t); err != nil {
			log.Errorf("%s:taskProvision:> set task err: %s", logTaskPrefix, err.Error())
			return err
		}

		return nil
	}

	if _, ok := js.task.queue[task.SelfLink().String()]; !ok {
		task.Status.State = models.StateQueued
		task.Meta.Updated = time.Now()

		if err := taskUpdate(js.storage, task, t); err != nil {
			log.Errorf("%s:taskProvision:> set task err: %s", logTaskPrefix, err.Error())
			return err
		}
		return nil
	}

	p, ok := js.pod.list[task.SelfLink().String()]
	if ok {
		// we look for pod for task.
		// if exists, then check the node binding
		// and the presence of the manifest.
		// ==============================================

		if p.Status.State == models.StateDestroy || p.Status.State == models.StateDestroyed {
			return nil
		}

		if p.Meta.Node == models.EmptyString {
			err := errors.New("node not attached")
			return fmt.Errorf("pod %s can not be manage: %v", p.Meta.SelfLink.String(), err)
		}

		m, err := pm.ManifestGet(p.Meta.Node, p.SelfLink().String())
		if err != nil {
			return err
		}

		if m == nil {
			if err := podManifestPut(js.storage, p); err != nil {
				return err
			}
		}

	} else {

		// if the task is set to the provision state,
		// then create and move task from queue to active
		// ==============================================

		pod, err := podCreate(js.storage, task)
		if err != nil {
			log.Errorf("%s:taskProvision:> creates new pod based on task spec err:", err.Error())
			return err
		}

		js.pod.list[task.SelfLink().String()] = pod
	}

	js.task.active[task.SelfLink().String()] = task
	delete(js.task.queue, task.SelfLink().String())

	return nil
}

func taskDestroy(js *JobState, task *models.Task) (err error) {

	t := task.Meta.Updated
	defer func() {
		if err == nil {
			err = taskUpdate(js.storage, task, t)
		}
	}()

	if task.Status.State != models.StateDestroy {
		task.Status.State = models.StateDestroy
		task.Meta.Updated = time.Now()
		return nil
	}

	p, ok := js.pod.list[task.SelfLink().String()]
	if !ok {
		task.Status.State = models.StateDestroyed
		task.Meta.Updated = time.Now()
		return nil
	}

	if p.Status.State != models.StateDestroy {
		if err := podDestroy(js, p); err != nil {
			return err
		}
	}

	if p.Status.State == models.StateDestroyed {
		if err := podRemove(js, p); err != nil {
			return err
		}
	}

	return nil
}

func taskUpdate(stg storage.IStorage, task *models.Task, timestamp time.Time) error {

	if timestamp.Before(task.Meta.Updated) {
		tm := service.NewTaskModel(context.Background(), stg)
		if err := tm.Set(task); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	return nil
}

func taskRemove(stg storage.IStorage, task *models.Task) error {
	tm := service.NewTaskModel(context.Background(), stg)
	if err := tm.Remove(task); err != nil {
		return err
	}
	return nil
}

func taskFinish(js *JobState, task *models.Task) (err error) {

	t := task.Meta.Updated
	defer func() {
		if err == nil {
			err = taskUpdate(js.storage, task, t)
		}
	}()

	if task.Status.State != models.StateExited {
		task.Status.Error = task.Status.State == models.StateError
		task.Status.Canceled = task.Status.State == models.StateCanceled
		task.Status.Done = !task.Status.Error && !task.Status.Canceled
		task.Status.State = models.StateExited
		task.Meta.Updated = time.Now()
	}

	p, ok := js.pod.list[task.SelfLink().String()]
	if ok {
		if p.Status.State != models.StateDestroy {
			if err := podDestroy(js, p); err != nil {
				return err
			}
		}
		if p.Status.State == models.StateDestroyed {
			if err := podRemove(js, p); err != nil {
				return err
			}
		}
	}

	js.task.finished = append(js.task.finished, task)
	delete(js.task.queue, task.SelfLink().String())
	delete(js.task.active, task.SelfLink().String())

	for {
		if len(js.task.finished) > 5 {
			var t *models.Task
			t, js.task.finished = js.task.finished[0], js.task.finished[1:]
			if t != nil {
				if err := taskDestroy(js, t); err != nil {
					log.Errorf("%s:> clean up task from finished list err: %s", logTaskPrefix, err.Error())
					break
				}
			}
			continue
		}
		break
	}

	if err = taskUpdate(js.storage, task, t); err != nil {
		return err
	}

	return nil
}

func taskStatusState(js *JobState, t *models.Task, p *models.Pod) (err error) {

	log.Infof("%s:task_status_state:> start: %s > %s", logTaskPrefix, t.SelfLink(), t.Status.State)

	u := t.Meta.Updated
	status := t.Status

	defer func() {
		if status.State == t.Status.State {
			return
		}

		if err == nil {
			if err := taskUpdate(js.storage, t, u); err != nil {
				log.Infof("%s:task_status_state:> update task %s err: %s", logTaskPrefix, t.Meta.Name, err.Error())
			}
		}

		log.Debugf("%s:task_status_state:> check task %s status (%s) > (%s)", logPrefix, t.SelfLink(), status.State, t.Status.State)

		if t.Status.State != status.State || t.Status.State == models.StateRunning {
			if err := js.Hook(t); err != nil {
				log.Errorf("%s:task_status_state:task> send state err: %s", logPrefix, err.Error())
			}
		}
	}()

	t.Status.Pod = models.TaskStatusPod{
		SelfLink: p.SelfLink().String(),
		State:    p.Status.State,
		Status:   p.Status.Status,
		Runtime:  p.Status.Runtime,
	}

	switch true {
	case p.Status.State == models.StateProvision:
		if t.Status.State == models.StateProvision {
			return nil
		}
		t.Status.State = models.StateProvision
		t.Status.Message = models.EmptyString
		t.Meta.Updated = time.Now()
	case p.Status.State == models.StateError:
		t.Status.State = models.StateError
		t.Status.Error = true
		t.Status.Message = p.Status.Message
		t.Meta.Updated = time.Now()
	case p.Status.State == models.StateExited && p.Status.Status == models.StateError:
		if t.Status.State != models.StateExited {
			t.Status.State = models.StateExited
			t.Status.Done = true
			t.Status.Message = p.Status.Message
			t.Meta.Updated = time.Now()
		}
		return nil
	case p.Status.State == models.StateDestroyed:
		fallthrough
	case p.Status.State == models.StateDestroy:
		fallthrough
	case p.Status.State == models.StateExited:
		t.Status.State = models.StateExited
		t.Status.Message = models.EmptyString
		t.Status.Done = !t.Status.Error && !t.Status.Canceled
		t.Meta.Updated = time.Now()
	}

	return nil
}
