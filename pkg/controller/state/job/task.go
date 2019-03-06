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
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/util/generator"
	"strings"

	"time"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logTaskPrefix = "state:observer:task"

func taskObserve(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:> observe start: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	switch task.Status.State {
	case types.StateCreated:
		if err := handleTaskStateCreated(js, task); err != nil {
			log.Errorf("%s:> handle task state create err: %s", logTaskPrefix, err.Error())
			return err
		}
	case types.StateQueued:
		if err := handleTaskStateQueued(js, task); err != nil {
			log.Errorf("%s:> handle task state queued err: %s", logTaskPrefix, err.Error())
			return err
		}
	case types.StateProvision:
		if err := handleTaskStateProvision(js, task); err != nil {
			log.Errorf("%s:> handle task state provision err: %s", logTaskPrefix, err.Error())
			return err
		}
	case types.StateRunning:
		if err := handleTaskStateRunning(js, task); err != nil {
			log.Errorf("%s:> handle task state ready err: %s", logTaskPrefix, err.Error())
			return err
		}
	case types.StateError:
		if err := handleTaskStateError(js, task); err != nil {
			log.Errorf("%s:> handle task state error err: %s", logTaskPrefix, err.Error())
			return err
		}
	case types.StateExited:
		if err := handleTaskStateExited(js, task); err != nil {
			log.Errorf("%s:> handle task state degradation err: %s", logTaskPrefix, err.Error())
			return err
		}
	case types.StateDestroy:
		if err := handleTaskStateDestroy(js, task); err != nil {
			log.Errorf("%s:> handle task state destroy err: %s", logTaskPrefix, err.Error())
			return err
		}
	case types.StateDestroyed:
		if err := handleTaskStateDestroyed(js, task); err != nil {
			log.Errorf("%s:> handle task state destroyed err: %s", logTaskPrefix, err.Error())
			return err
		}
	}

	if task.Status.State == types.StateDestroyed {
		delete(js.task.list, task.SelfLink().String())
	} else {
		js.task.list[task.SelfLink().String()] = task
	}

	log.V(logLevel).Debugf("%s:> observe finish: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	return nil
}

func handleTaskStateCreated(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:> handleTaskStateCreated: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	if err := taskCheckSelectors(js, task); err != nil {
		task.Status.State = types.StateError
		task.Status.Status = types.StateError
		task.Status.Message = err.Error()
		tm := distribution.NewTaskModel(context.Background(), envs.Get().GetStorage())
		if err := tm.Set(task); err != nil {
			log.Errorf("%s:> handle task create, deps update: %s, err: %s", logTaskPrefix, task.SelfLink(), err.Error())
			return err
		}
		return nil
	}

	if err := taskQueue(js, task); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handleTaskStateQueued(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:> handleTaskStateQueued: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	if err := taskQueue(js, task); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handleTaskStateProvision(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:> handleTaskStateProvision: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	// check pods are created and state is normal state
	if err := taskProvision(js, task); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handleTaskStateRunning(_ *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:> handleTaskStateRunning: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)
	// there nothing need to be done

	return nil
}

func handleTaskStateError(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:> handleTaskStateError: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)
	// finish task and destroy it

	if err := taskFinish(js, task); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handleTaskStateExited(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:> handleTaskStateExited: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)
	// finish task and destroy it

	if err := taskFinish(js, task); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handleTaskStateDestroy(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:> handleTaskStateDestroy: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	if err := taskDestroy(js, task); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handleTaskStateDestroyed(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:> handleTaskStateDestroyed: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	link := task.SelfLink().String()

	if _, ok := js.pod.list[link]; ok {

		if err := taskDestroy(js, task); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
		return nil
	}

	if err := taskRemove(task); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	// TODO: check of nil
	js.DelTask(task)
	return nil
}

// taskCheckSelectors function - handles provided selectors to match nodes
func taskCheckSelectors(_ *JobState, d *types.Task) (err error) {

	var (
		ctx = context.Background()
		stg = envs.Get().GetStorage()
		vm  = distribution.NewVolumeModel(ctx, stg)
		vc  = make(map[string]string, 0)
	)

	for _, v := range d.Spec.Template.Volumes {
		if v.Volume.Name != types.EmptyString {
			vc[v.Volume.Name] = v.Name
		}
	}

	if len(vc) > 0 {

		var node string

		vl, err := vm.ListByNamespace(d.Meta.Namespace)
		if err != nil {
			log.V(logLevel).Errorf("%s:create:> create task, volume list err: %s", logPrefix, err.Error())
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
					log.V(logLevel).Errorf("%s:create:> create task err: volume is not ready yet: %s", logPrefix, v.Meta.Name)
					return errors.New(v.Meta.Name).Volume().NotReady(v.Meta.Name)
				}

				if v.Meta.Node == types.EmptyString {
					log.V(logLevel).Errorf("%s:create:> create task err: volume is not provisioned yet: %s", logPrefix, v.Meta.Name)
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
				log.V(logLevel).Errorf("%s:create:> create task err: volume is not found: %s", logPrefix, name)
				return errors.New(name).Volume().NotFound(name)
			}
		}

		if node != types.EmptyString {

			if d.Spec.Selector.Node != types.EmptyString {
				if d.Spec.Selector.Node != node {
					return errors.New("spec.selector.node not matched with attached volumes")
				}

				return nil
			}

			d.Spec.Selector.Node = node
		}

	}

	return nil
}

// taskCreate - create a new task from current job
// usually used by cron or other time repeatable jobs
func taskCreate(job *types.Job, mf *types.TaskManifest) (*types.Task, error) {

	tm := distribution.NewTaskModel(context.Background(), envs.Get().GetStorage())

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
			return nil, err
		}
	}

	d, err := tm.Create(task)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func taskQueue(js *JobState, task *types.Task) error {

	js.task.queue[task.SelfLink().String()] = task

	if task.Status.State != types.StateQueued {
		task.Status.State = types.StateQueued
		task.Status.Status = types.StateQueued
		tm := distribution.NewTaskModel(context.Background(), envs.Get().GetStorage())
		if err := tm.Set(task); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
		return nil
	}

	if err := jobTaskProvision(js); err != nil {
		log.Errorf("%s:> job task queue pop err: %s", logTaskPrefix, err.Error())
		return err
	}

	return nil
}

// taskProvision - handles task provision logic
// based on current task state and current pod list of provided task
func taskProvision(js *JobState, task *types.Task) (err error) {

	t := task.Meta.Updated

	var (
		provision = false
	)

	defer func() {
		if err == nil {
			err = taskUpdate(task, t)
		}
	}()

	var (
		pm = distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	)

	p, ok := js.pod.list[task.SelfLink().String()]
	if ok {
		if p.Status.State != types.StateDestroy && p.Status.State != types.StateDestroyed {

			if p.Meta.Node != types.EmptyString {

				m, e := pm.ManifestGet(p.Meta.Node, p.SelfLink().String())
				if err != nil {
					err = e
					return e
				}

				if m == nil {
					if err = podManifestPut(p); err != nil {
						return err
					}
				}

			}

		}
	} else {
		_, err := podCreate(task)
		if err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
		provision = true
	}

	if provision {
		if task.Status.State != types.StateProvision {
			task.Status.State = types.StateProvision
			task.Status.Status = types.StateProvision
			task.Meta.Updated = time.Now()
		}

		js.task.active[task.SelfLink().String()] = task
	}

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
	}

	p, ok := js.pod.list[task.SelfLink().String()]
	if !ok {
		task.Status.State = types.StateDestroyed
		task.Meta.Updated = time.Now()
		return nil
	} else {
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

	return nil
}

func taskUpdate(task *types.Task, timestamp time.Time) error {
	if timestamp.Before(task.Meta.Updated) {
		tm := distribution.NewTaskModel(context.Background(), envs.Get().GetStorage())
		if err := tm.Set(task); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	return nil
}

func taskRemove(task *types.Task) error {
	dm := distribution.NewTaskModel(context.Background(), envs.Get().GetStorage())
	if err := dm.Remove(task); err != nil {
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
		task.Status.State = types.StateExited
		task.Status.Status = types.StateReady
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

	js.task.finished = append(js.task.finished, task)
	return nil
}

func taskStatusState(t *types.Task, p *types.Pod) (err error) {

	log.V(logLevel).Infof("%s:task_status_state:> start: %s > %s", logTaskPrefix, t.SelfLink(), t.Status.State)

	u := t.Meta.Updated
	defer func() {
		log.V(logLevel).Infof("%s:task_status_state:> finish: %s > %s", logTaskPrefix, t.SelfLink(), t.Status.State)
		if err == nil {
			err = taskUpdate(t, u)
		}
	}()

	t.Status.Pod = types.TaskStatusPod{
		SelfLink: p.SelfLink().String(),
		State:    p.Status.State,
		Status:   p.Status.Status,
		Runtime:  p.Status.Runtime,
	}

	switch true {
	case p.Status.State == types.StateError:
		if t.Status.State != types.StateExited {
			t.Status.State = types.StateExited
			t.Status.Status = types.StateError
			t.Status.Message = p.Status.Message
			t.Meta.Updated = time.Now()
		}
		return nil
	case p.Status.Status == types.StateError:
		if t.Status.State != types.StateExited {
			t.Status.State = types.StateExited
			t.Status.Status = types.StateReady
			t.Status.Message = p.Status.Message
			t.Meta.Updated = time.Now()
		}
		return nil
	case p.Status.Status == types.StateRunning:
		if t.Status.State != types.StateRunning {
			t.Status.State = types.StateRunning
			t.Status.Status = types.StateRunning
			t.Status.Message = types.EmptyString
			t.Meta.Updated = time.Now()
		}
		return nil
	case p.Status.Status == types.StateExited:
		if t.Status.State != types.StateExited {
			t.Status.State = types.StateExited
			t.Status.Status = types.StateReady
			t.Status.Message = types.EmptyString
			t.Meta.Updated = time.Now()
		}
		return nil
	}

	return nil
}
