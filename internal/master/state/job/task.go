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
	"github.com/lastbackend/lastbackend/internal/master/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/model"
	"strings"
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/internal/util/generator"
	"github.com/lastbackend/lastbackend/tools/log"
)

const logTaskPrefix = "state:observer:task"

func taskObserve(js *JobState, task *types.Task) (err error) {

	log.V(logLevel).Debugf("%s:> observe start: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	switch task.Status.State {
	case types.StateCreated:
		err = handleTaskStateCreated(js, task)
	case types.StateQueued:
		err = handleTaskStateQueued(js, task)
	case types.StateProvision:
		err = handleTaskStateProvision(js, task)
	case types.StateRunning:
		err = handleTaskStateRunning(js, task)
	case types.StateError:
		err = handleTaskStateError(js, task)
	case types.StateCanceled:
		err = handleTaskStateCanceled(js, task)
	case types.StateExited:
		err = handleTaskStateExited(js, task)
	case types.StateDestroy:
		err = handleTaskStateDestroy(js, task)
	case types.StateDestroyed:
		err = handleTaskStateDestroyed(js, task)
	}
	if err != nil {
		task.Status.State = types.StateError
		task.Status.Error = true
		task.Status.Message = err.Error()
		if err := handleTaskStateError(js, task); err != nil {
			log.Errorf("%s:> handle task state %s error err: %s", logTaskPrefix, task.Status.State, err.Error())
			return err
		}
		return nil
	}

	log.V(logLevel).Debugf("%s:> observe finish: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	return nil
}

func handleTaskStateCreated(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:handleTaskStateCreated:> try to handle task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	if err := taskCheckSelectors(js, task); err != nil {
		task.Status.State = types.StateError
		task.Status.Error = true
		task.Status.Message = err.Error()
		tm := model.NewTaskModel(context.Background(), envs.Get().GetStorage())
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

func handleTaskStateQueued(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:handleTaskStateQueued:> task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	if err := taskQueue(js, task); err != nil {
		log.Errorf("%s:handleTaskStateProvision:> task queued err: %s", logTaskPrefix, err.Error())
		return err
	}

	return nil
}

func handleTaskStateProvision(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:handleTaskStateProvision:> task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	// check pods are created and state is normal state
	if err := taskProvision(js, task); err != nil {
		log.Errorf("%s:handleTaskStateProvision:> task provision err: %s", logTaskPrefix, err.Error())
		return err
	}

	return nil
}

func handleTaskStateRunning(_ *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:handleTaskStateRunning:> task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)
	// there nothing need to be done

	return nil
}

func handleTaskStateError(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:handleTaskStateError:> task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	// finish task and destroy it
	if err := taskFinish(js, task); err != nil {
		log.Errorf("%s:handleTaskStateError:> task finish err: %s", logTaskPrefix, err.Error())
		return err
	}

	return nil
}

func handleTaskStateCanceled(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:handleTaskStateCanceled:> task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	// finish task and destroy it
	if err := taskFinish(js, task); err != nil {
		log.Errorf("%s:handleTaskStateCanceled:> task finish err: %s", logTaskPrefix, err.Error())
		return err
	}

	return nil
}

func handleTaskStateExited(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:handleTaskStateExited:>: task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)
	// finish task and destroy it
	if err := taskFinish(js, task); err != nil {
		log.Errorf("%s:handleTaskStateExited:> task finish err: %s", logTaskPrefix, err.Error())
		return err
	}

	return nil
}

func handleTaskStateDestroy(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:handleTaskStateDestroy:> task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	if err := taskDestroy(js, task); err != nil {
		log.Errorf("%s:handleTaskStateDestroy:> task destroy err: %s", logTaskPrefix, err.Error())
		return err
	}

	return nil
}

func handleTaskStateDestroyed(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:handleTaskStateDestroyed:> task %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	if _, ok := js.pod.list[task.SelfLink().String()]; ok {
		if err := taskDestroy(js, task); err != nil {
			log.Errorf("%s:handleTaskStateDestroyed:> task destroy err: %s", logTaskPrefix, err.Error())
			return err
		}
		return nil
	}
	if err := taskRemove(task); err != nil {
		log.Errorf("%s:handleTaskStateDestroyed:> remove task err: %s", logTaskPrefix, err.Error())
		return err
	}
	js.DelTask(task)
	return nil
}

// taskCheckSelectors function - handles provided selectors to match nodes
func taskCheckSelectors(_ *JobState, task *types.Task) (err error) {

	var (
		ctx = context.Background()
		stg = envs.Get().GetStorage()
		vm  = model.NewVolumeModel(ctx, stg)
		vc  = make(map[string]string, 0)
	)

	for _, v := range task.Spec.Template.Volumes {
		if v.Volume.Name != types.EmptyString {
			vc[v.Volume.Name] = v.Name
		}
	}

	if len(vc) > 0 {

		var node string

		vl, err := vm.ListByNamespace(task.Meta.Namespace)
		if err != nil {
			log.V(logLevel).Errorf("%s:check_selectors:> create task, volume list err: %s", logPrefix, err.Error())
			return err
		}

		for name := range vc {

			var f = false

			for _, v := range vl.Items {

				if v.Meta.Name != name {
					continue
				}

				f = true

				if v.Status.State != types.StateReady {
					log.V(logLevel).Errorf("%s:check_selectors:> create task err: volume is not ready yet: %s", logPrefix, v.Meta.Name)
					return errors.New(v.Meta.Name).Volume().NotReady(v.Meta.Name)
				}

				if v.Meta.Node == types.EmptyString {
					log.V(logLevel).Errorf("%s:check_selectors:> create task err: volume is not provisioned yet: %s", logPrefix, v.Meta.Name)
					return errors.New(v.Meta.Name).Volume().NotProvisioned(v.Meta.Name)
				}

				if node == types.EmptyString {
					node = v.Meta.Node
				} else {
					if node != v.Meta.Node {
						return errors.New(v.Meta.Name).Volume().DifferentNodes()
					}
				}
			}

			if !f {
				log.V(logLevel).Errorf("%s:check_selectors:> create task err: volume is not found: %s", logPrefix, name)
				return errors.New(name).Volume().NotFound(name)
			}
		}

		if node != types.EmptyString {

			if task.Spec.Selector.Node != types.EmptyString {
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
func taskCreate(job *types.Job, mf *types.TaskManifest) (*types.Task, error) {

	tm := model.NewTaskModel(context.Background(), envs.Get().GetStorage())

	task := new(types.Task)
	task.Meta.SetDefault()
	task.Meta.Namespace = job.Meta.Namespace
	task.Meta.Job = job.SelfLink().String()

	if mf != nil {
		mf.SetTaskMeta(task)
	}

	if task.Meta.Name == types.EmptyString {
		name := strings.Split(generator.GetUUIDV4(), "-")[4][5:]
		task.Meta.Name = name
	}

	task.Meta.SelfLink = *types.NewTaskSelfLink(job.Meta.Namespace, job.Meta.Name, task.Meta.Name)

	task.Spec.Runtime = job.Spec.Task.Runtime
	task.Spec.Selector = job.Spec.Task.Selector
	task.Spec.Template = job.Spec.Task.Template

	if mf != nil {
		if err := mf.SetTaskSpec(task); err != nil {
			log.V(logLevel).Errorf("%s:taskCreate:> set task spec err: %v", logTaskPrefix, err)
			return nil, err
		}
	}

	d, err := tm.Create(task)
	if err != nil {
		log.V(logLevel).Errorf("%s:taskCreate:> create task err: %v", logTaskPrefix, err)
		return nil, err
	}

	return d, nil
}

func taskQueue(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:taskQueue:> move task %s to queue", logTaskPrefix, task.Meta.Name)

	// set a task as queued if it is not right now
	// ==============================================

	if task.Status.State != types.StateQueued {

		task.Status.State = types.StateQueued

		tm := model.NewTaskModel(context.Background(), envs.Get().GetStorage())
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
func taskProvision(js *JobState, task *types.Task) (err error) {

	log.V(logLevel).Debugf("%s:taskProvision:> set task %s as provision", logTaskPrefix, task.Meta.Name)

	var (
		t  = task.Meta.Updated
		pm = model.NewPodModel(context.Background(), envs.Get().GetStorage())
	)

	// set a task as provision if it is not right now
	// ==============================================

	if task.Status.State != types.StateProvision {
		task.Status.State = types.StateProvision
		task.Meta.Updated = time.Now()

		if err := taskUpdate(task, t); err != nil {
			log.Errorf("%s:taskProvision:> set task err: %s", logTaskPrefix, err.Error())
			return err
		}

		return nil
	}

	if _, ok := js.task.queue[task.SelfLink().String()]; !ok {
		task.Status.State = types.StateQueued
		task.Meta.Updated = time.Now()

		if err := taskUpdate(task, t); err != nil {
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

		if p.Status.State == types.StateDestroy || p.Status.State == types.StateDestroyed {
			return nil
		}

		if p.Meta.Node == types.EmptyString {
			err := errors.New("node not attached")
			return errors.New(fmt.Sprintf("pod %s can not be manage: %v", p.Meta.SelfLink.String(), err))
		}

		m, err := pm.ManifestGet(p.Meta.Node, p.SelfLink().String())
		if err != nil {
			return err
		}

		if m == nil {
			if err := podManifestPut(p); err != nil {
				return err
			}
		}

	} else {

		// if the task is set to the provision state,
		// then create and move task from queue to active
		// ==============================================

		pod, err := podCreate(task)
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

func taskDestroy(js *JobState, task *types.Task) (err error) {

	t := task.Meta.Updated
	defer func() {
		if err == nil {
			err = taskUpdate(task, t)
		}
	}()

	if task.Status.State != types.StateDestroy {
		task.Status.State = types.StateDestroy
		task.Meta.Updated = time.Now()
		return nil
	}

	p, ok := js.pod.list[task.SelfLink().String()]
	if !ok {
		task.Status.State = types.StateDestroyed
		task.Meta.Updated = time.Now()
		return nil
	}

	if p.Status.State != types.StateDestroy {
		if err := podDestroy(js, p); err != nil {
			return err
		}
	}

	if p.Status.State == types.StateDestroyed {
		if err := podRemove(js, p); err != nil {
			return err
		}
	}

	return nil
}

func taskUpdate(task *types.Task, timestamp time.Time) error {

	if timestamp.Before(task.Meta.Updated) {
		tm := model.NewTaskModel(context.Background(), envs.Get().GetStorage())
		if err := tm.Set(task); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	return nil
}

func taskRemove(task *types.Task) error {
	tm := model.NewTaskModel(context.Background(), envs.Get().GetStorage())
	if err := tm.Remove(task); err != nil {
		return err
	}
	return nil
}

func taskFinish(js *JobState, task *types.Task) (err error) {

	t := task.Meta.Updated
	defer func() {
		if err == nil {
			err = taskUpdate(task, t)
		}
	}()

	if task.Status.State != types.StateExited {
		task.Status.Error = task.Status.State == types.StateError
		task.Status.Canceled = task.Status.State == types.StateCanceled
		task.Status.Done = !task.Status.Error && !task.Status.Canceled
		task.Status.State = types.StateExited
		task.Meta.Updated = time.Now()
	}

	p, ok := js.pod.list[task.SelfLink().String()]
	if ok {
		if p.Status.State != types.StateDestroy {
			if err := podDestroy(js, p); err != nil {
				return err
			}
		}
		if p.Status.State == types.StateDestroyed {
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
			var t *types.Task
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

	if err = taskUpdate(task, t); err != nil {
		return err
	}

	return nil
}

func taskStatusState(js *JobState, t *types.Task, p *types.Pod) (err error) {

	log.V(logLevel).Infof("%s:task_status_state:> start: %s > %s", logTaskPrefix, t.SelfLink(), t.Status.State)

	u := t.Meta.Updated
	status := t.Status

	defer func() {
		if status.State == t.Status.State {
			return
		}

		if err == nil {
			if err := taskUpdate(t, u); err != nil {
				log.V(logLevel).Infof("%s:task_status_state:> update task %s err: %s", logTaskPrefix, t.Meta.Name, err.Error())
			}
		}

		log.V(logLevel).Debugf("%s:task_status_state:> check task %s status (%s) > (%s)", logPrefix, t.SelfLink(), status.State, t.Status.State)

		if t.Status.State != status.State || t.Status.State == types.StateRunning {
			if err := js.Hook(t); err != nil {
				log.Errorf("%s:task_status_state:task> send state err: %s", logPrefix, err.Error())
			}
		}
	}()

	t.Status.Pod = types.TaskStatusPod{
		SelfLink: p.SelfLink().String(),
		State:    p.Status.State,
		Status:   p.Status.Status,
		Runtime:  p.Status.Runtime,
	}

	switch true {
	case p.Status.State == types.StateProvision:
		if t.Status.State == types.StateProvision {
			return nil
		}
		t.Status.State = types.StateProvision
		t.Status.Message = types.EmptyString
		t.Meta.Updated = time.Now()
	case p.Status.State == types.StateError:
		t.Status.State = types.StateError
		t.Status.Error = true
		t.Status.Message = p.Status.Message
		t.Meta.Updated = time.Now()
	case p.Status.State == types.StateExited && p.Status.Status == types.StateError:
		if t.Status.State != types.StateExited {
			t.Status.State = types.StateExited
			t.Status.Done = true
			t.Status.Message = p.Status.Message
			t.Meta.Updated = time.Now()
		}
		return nil
	case p.Status.State == types.StateDestroyed:
		fallthrough
	case p.Status.State == types.StateDestroy:
		fallthrough
	case p.Status.State == types.StateExited:
		t.Status.State = types.StateExited
		t.Status.Message = types.EmptyString
		t.Status.Done = !t.Status.Error && !t.Status.Canceled
		t.Meta.Updated = time.Now()
	}

	return nil
}
