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

	if _, ok := js.pod.list[task.SelfLink().String()]; !ok {
		js.pod.list[task.SelfLink().String()] = make(map[string]*types.Pod)
	}

	switch task.Status.State {
	case types.StateCreated:
		if err := handleTaskStateCreated(js, task); err != nil {
			log.Errorf("%s:> handle task state create err: %s", logTaskPrefix, err.Error())
			return err
		}
		break
	case types.StateQueued:
		if err := handleTaskStateQueued(js, task); err != nil {
			log.Errorf("%s:> handle task state queued err: %s", logTaskPrefix, err.Error())
			return err
		}
		break
	case types.StateProvision:
		if err := handleTaskStateProvision(js, task); err != nil {
			log.Errorf("%s:> handle task state provision err: %s", logTaskPrefix, err.Error())
			return err
		}
		break
	case types.StateRunning:
		if err := handleTaskStateRunning(js, task); err != nil {
			log.Errorf("%s:> handle task state ready err: %s", logTaskPrefix, err.Error())
			return err
		}
		break
	case types.StateError:
		if err := handleTaskStateError(js, task); err != nil {
			log.Errorf("%s:> handle task state error err: %s", logTaskPrefix, err.Error())
			return err
		}
		break
	case types.StateExited:
		if err := handleTaskStateExited(js, task); err != nil {
			log.Errorf("%s:> handle task state degradation err: %s", logTaskPrefix, err.Error())
			return err
		}
		break
	case types.StateDestroy:
		if err := handleTaskStateDestroy(js, task); err != nil {
			log.Errorf("%s:> handle task state destroy err: %s", logTaskPrefix, err.Error())
			return err
		}
		break
	case types.StateDestroyed:
		if err := handleTaskStateDestroyed(js, task); err != nil {
			log.Errorf("%s:> handle task state destroyed err: %s", logTaskPrefix, err.Error())
			return err
		}
		break
	}

	if task.Status.State == types.StateDestroyed {
		delete(js.task.list, task.SelfLink().String())
	} else {
		js.task.list[task.SelfLink().String()] = task
	}

	if err := jobTaskProvision(js); err != nil {
		log.Errorf("%s:> job task queue pop err: %s", err.Error())
		return err
	}

	log.V(logLevel).Debugf("%s:> observe finish: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	return nil
}

func handleTaskStateCreated(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:> handleTaskStateCreated: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	check, err := taskCheckDependencies(js, task)
	if err != nil {
		log.Errorf("%s:> handle task check deps: %s, err: %s", logTaskPrefix, task.SelfLink(), err.Error())
		return err
	}

	if !check {
		task.Status.State = types.StateWaiting
		tm := distribution.NewTaskModel(context.Background(), envs.Get().GetStorage())
		if err := tm.Set(task); err != nil {
			log.Errorf("%s:> handle task create, deps update: %s, err: %s", logTaskPrefix, task.SelfLink(), err.Error())
			return err
		}
		return nil
	}

	if err := taskCheckSelectors(js, task); err != nil {
		task.Status.State = types.StateError
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

func handleTaskStateRunning(js *JobState, task *types.Task) error {

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

	if task.Status.State == types.StateDestroyed {
		return handleTaskStateDestroyed(js, task)
	}

	return nil
}

func handleTaskStateDestroyed(js *JobState, task *types.Task) error {

	log.V(logLevel).Debugf("%s:> handleTaskStateDestroyed: %s > %s", logTaskPrefix, task.SelfLink(), task.Status.State)

	link := task.SelfLink().String()

	if _, ok := js.pod.list[link]; ok && len(js.pod.list[link]) > 0 {

		if err := taskDestroy(js, task); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}

		task.Status.State = types.StateDestroy
		tm := distribution.NewTaskModel(context.Background(), envs.Get().GetStorage())
		return tm.Set(task)
	}

	if err := taskRemove(task); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	js.DelTask(task)
	return nil
}

// jobCheckDependencies function - check if job can provisioned or should wait for dependencies
func taskCheckDependencies(ss *JobState, d *types.Task) (bool, error) {

	var (
		ctx  = context.Background()
		stg  = envs.Get().GetStorage()
		vm   = distribution.NewVolumeModel(ctx, stg)
		sm   = distribution.NewSecretModel(ctx, stg)
		cm   = distribution.NewConfigModel(ctx, stg)
		deps = types.StatusDependencies{
			Volumes: make(map[string]types.StatusDependency, 0),
			Secrets: make(map[string]types.StatusDependency, 0),
			Configs: make(map[string]types.StatusDependency, 0),
		}
	)

	volumesRequiredList := make(map[string]bool, 0)
	secretsRequiredList := make(map[string]bool, 0)
	configsRequiredList := make(map[string]bool, 0)
	for _, v := range d.Spec.Template.Volumes {
		if v.Volume.Name != types.EmptyString {
			volumesRequiredList[v.Volume.Name] = true
		}
		if v.Secret.Name != types.EmptyString {
			secretsRequiredList[v.Secret.Name] = true
		}
		if v.Config.Name != types.EmptyString {
			configsRequiredList[v.Config.Name] = true
		}
	}

	for _, c := range d.Spec.Template.Containers {
		for _, e := range c.EnvVars {
			if e.Secret.Name != types.EmptyString {
				secretsRequiredList[e.Secret.Name] = true
			}

			if e.Config.Name != types.EmptyString {
				configsRequiredList[e.Config.Name] = true
			}
		}
	}

	if len(volumesRequiredList) != 0 {

		vl, err := vm.ListByNamespace(d.Meta.Namespace)
		if err != nil {
			log.Errorf("%s:> job check deps err: %s", logJobPrefix, err.Error())
			return false, err
		}

		for vr := range volumesRequiredList {
			var f = false

			for _, v := range vl.Items {
				if vr == v.Meta.Name {
					f = true
					deps.Volumes[vr] = types.StatusDependency{
						Name:   vr,
						Type:   types.KindVolume,
						Status: v.Status.State,
					}
				}
			}

			if !f {
				deps.Volumes[vr] = types.StatusDependency{
					Name:   vr,
					Type:   types.KindVolume,
					Status: types.StateNotReady,
				}
			}
		}
	}

	if len(secretsRequiredList) != 0 {

		sl, err := sm.List(d.Meta.Namespace)
		if err != nil {
			log.Errorf("%s:> job check deps err: %s", logJobPrefix, err.Error())
			return false, err
		}

		for sr := range secretsRequiredList {
			var f = false

			for _, s := range sl.Items {
				if sr == s.Meta.Name {
					f = true
					deps.Secrets[sr] = types.StatusDependency{
						Name:   sr,
						Type:   types.KindSecret,
						Status: types.StateReady,
					}
				}
			}

			if !f {
				deps.Secrets[sr] = types.StatusDependency{
					Name:   sr,
					Type:   types.KindSecret,
					Status: types.StateNotReady,
				}
			}
		}
	}

	if len(configsRequiredList) != 0 {

		cl, err := cm.List(d.Meta.Namespace)
		if err != nil {
			log.Errorf("%s:> job check deps err: %s", logJobPrefix, err.Error())
			return false, err
		}

		for cr := range configsRequiredList {
			var f = false

			for _, c := range cl.Items {
				if cr == c.Meta.Name {
					f = true
					deps.Configs[cr] = types.StatusDependency{
						Name:   cr,
						Type:   types.KindConfig,
						Status: types.StateReady,
					}
				}
			}

			if !f {
				deps.Configs[cr] = types.StatusDependency{
					Name:   cr,
					Type:   types.KindConfig,
					Status: types.StateNotReady,
				}
			}
		}
	}

	d.Status.Dependencies = deps
	if !d.Status.CheckDeps() {
		d.Status.State = types.StateWaiting
		return false, nil
	}

	return true, nil
}

// taskCheckSelectors function - handles provided selectors to match nodes
func taskCheckSelectors(ss *JobState, d *types.Task) (err error) {

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
func taskCreate(job *types.Job) (*types.Task, error) {

	tm := distribution.NewTaskModel(context.Background(), envs.Get().GetStorage())

	task := new(types.Task)

	task.Meta.Namespace = job.Meta.Namespace
	task.Meta.Job = job.SelfLink().String()

	name := strings.Split(generator.GetUUIDV4(), "-")[4][5:]
	task.Meta.Name = name
	task.Meta.SelfLink = *types.NewTaskSelfLink(job.Meta.Name, job.Meta.Name, name)

	task.Spec.Runtime = job.Spec.Task.Runtime
	task.Spec.Selector = job.Spec.Task.Selector
	task.Spec.Template = job.Spec.Task.Template

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
		tm := distribution.NewTaskModel(context.Background(), envs.Get().GetStorage())
		if err := tm.Set(task); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
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

	pods, ok := js.pod.list[task.SelfLink().String()]
	if !ok {
		pods = make(map[string]*types.Pod, 0)
	}

	var (
		total int
		state = make(map[string][]*types.Pod)
	)

	for _, p := range pods {

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

			total++
		}

		if _, ok := state[p.Status.State]; !ok {
			state[p.Status.State] = make([]*types.Pod, 0)
		}

		state[p.Status.State] = append(state[p.Status.State], p)
	}

	if total < 1 {
		p, err := podCreate(task)
		if err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
		pods[p.SelfLink().String()] = p
		provision = true
	}

	if provision {
		if task.Status.State != types.StateProvision {
			task.Status.State = types.StateProvision
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

	pl, ok := js.pod.list[task.SelfLink().String()]
	if !ok {
		task.Status.State = types.StateDestroyed
		task.Meta.Updated = time.Now()
		return nil
	}

	for _, p := range pl {

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

	if len(pl) == 0 {
		task.Status.State = types.StateDestroyed
		task.Meta.Updated = time.Now()
		return nil
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
		task.Meta.Updated = time.Now()
	}

	pl, ok := js.pod.list[task.SelfLink().String()]
	if !ok {
		task.Status.State = types.StateDestroyed
		task.Meta.Updated = time.Now()
		return nil
	}

	for _, p := range pl {
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

func taskStatusState(d *types.Task, pl map[string]*types.Pod) (err error) {

	log.V(logLevel).Infof("%s:> taskStatusState: start: %s > %s", logTaskPrefix, d.SelfLink(), d.Status.State)

	t := d.Meta.Updated
	defer func() {
		log.V(logLevel).Infof("%s:> taskStatusState: finish: %s > %s", logTaskPrefix, d.SelfLink(), d.Status.State)
		if err == nil {
			err = taskUpdate(d, t)
		}
	}()

	var (
		state   = make(map[string]int)
		message string
		total   int
	)

	for _, p := range pl {
		total++

		if p.Status.Status == types.StateError {
			message = p.Status.Message
		}

		state[p.Status.Status]++
	}

	if state[types.StateRunning] > 0 {
		if d.Status.State != types.StateRunning {
			d.Status.State = types.StateRunning
			d.Status.Message = types.EmptyString
			d.Meta.Updated = time.Now()
		}
		return nil
	}

	if state[types.StateError] > 0 {
		if d.Status.State != types.StateExited {
			d.Status.State = types.StateExited
			d.Status.Message = message
			d.Meta.Updated = time.Now()
		}
		return nil
	}

	if state[types.StateExited] > 0 && state[types.StateExited] == total {
		if d.Status.State != types.StateExited {
			d.Status.State = types.StateExited
			d.Status.Message = types.EmptyString
			d.Meta.Updated = time.Now()
		}
		return nil
	}

	return nil
}
