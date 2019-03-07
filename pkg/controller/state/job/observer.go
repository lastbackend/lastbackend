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
	"github.com/lastbackend/lastbackend/pkg/controller/state/job/hook"
	jh "github.com/lastbackend/lastbackend/pkg/controller/state/job/hook/hook"
	"github.com/lastbackend/lastbackend/pkg/controller/state/job/provider"
	jp "github.com/lastbackend/lastbackend/pkg/controller/state/job/provider/provider"
	"sync"
	"time"

	"github.com/lastbackend/lastbackend/pkg/controller/envs"
	"github.com/lastbackend/lastbackend/pkg/controller/state/cluster"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const (
	logLevel  = 3
	logPrefix = "state:job"
)

type JobState struct {
	lock sync.Mutex

	cluster *cluster.ClusterState
	job     *types.Job

	task struct {
		active   map[string]*types.Task
		queue    map[string]*types.Task
		list     map[string]*types.Task
		finished []*types.Task
	}

	pod struct {
		list map[string]*types.Pod
	}

	observers struct {
		job  chan *types.Job
		task chan *types.Task
		pod  chan *types.Pod
	}

	provider provider.JobProvider
	hook     hook.Hook
}

func (js *JobState) Namespace() string {
	return js.job.Meta.Namespace
}

func (js *JobState) Restore() error {

	log.V(logLevel).Debugf("%s:restore state for job: %s", logPrefix, js.job.SelfLink())

	var (
		err error
		stg = envs.Get().GetStorage()
	)

	// Get all pods
	pm := distribution.NewPodModel(context.Background(), stg)
	pl, err := pm.ListByJob(js.job.Meta.Namespace, js.job.Meta.Name)
	if err != nil {
		log.Errorf("%s:restore:> get pod map error: %v", logPrefix, err)
		return err
	}

	for _, p := range pl.Items {
		log.Infof("%s: restore: restore pod: %s", logPrefix, p.SelfLink())

		// Check if task map for pod exists
		_, sl := p.SelfLink().Parent()
		// put pod into map by task name and pod name
		js.pod.list[sl.String()] = p
	}

	// Get all tasks

	tm := distribution.NewTaskModel(context.Background(), stg)
	tl, err := tm.ListByJob(js.job.Meta.Namespace, js.job.Meta.Name)
	if err != nil {
		log.Errorf("%s:restore:> get task map error: %v", logPrefix, err)
		return err
	}

	for _, t := range tl.Items {
		log.Infof("%s: restore task: %s", logPrefix, t.SelfLink())
		js.task.list[t.SelfLink().String()] = t
	}

	// Range over pods to sync pod status
	for _, p := range js.pod.list {
		js.observers.pod <- p
	}

	// Range over tasks to sync tasks status
	for _, d := range js.task.list {
		js.observers.task <- d
	}

	// Sync job state if updated
	js.observers.job <- js.job

	if js.provider != nil {
		go js.Provider()
	}

	if err := jobTaskProvision(js); err != nil {
		log.Errorf("%s:> job task provision err: %s", logPrefix, err.Error())
		return err
	}

	return nil
}

func (js *JobState) Observe() {

	for {
		select {

		case pod := <-js.observers.pod:
			log.V(logLevel).Debugf("%s:observe:pod:> %s", logPrefix, pod.SelfLink())
			if err := PodObserve(js, pod); err != nil {
				log.Errorf("%s:observe:pod:> err: %s", logPrefix, err.Error())
			}
			break

		case task := <-js.observers.task:

			log.V(logLevel).Debugf("%s:observe:task:> %s (%s)", logPrefix, task.SelfLink(), task.Status.State)

			status := task.Status

			if err := taskObserve(js, task); err != nil {
				log.Errorf("%s:observe:task err:> %s", logPrefix, err.Error())

				attempt := 0
				for {
					attempt++

					log.V(logLevel).Debugf("%s:observe:task:> repeat task %s attempt %в", logPrefix, task.Meta.Name, attempt)
					if err := taskObserve(js, task); err == nil {
						log.Errorf("%s:observe:task err:> %s", logPrefix, err.Error())
						break
					}

					delay := time.Duration(attempt * 1)
					if delay > 10 {
						delay = 10
					}

					<-time.After(delay * time.Second)
				}
			}

			log.V(logLevel).Debugf("%s:observe:task:> check task %s status (%s) <-> (%s)", logPrefix, task.SelfLink(), status.Status, task.Status.State)

			if task.Status.State != status.State || task.Status.Status != status.Status {
				if err := js.Hook(task); err != nil {
					log.Errorf("%s:observe:task> send state err: %s", logPrefix, err.Error())
				}
			}

			break

		case job := <-js.observers.job:
			log.V(logLevel).Debugf("%s:observe:job:> %s", logPrefix, job.SelfLink())

			js.job = job

			if err := jobObserve(js, job); err != nil {
				log.Errorf("%s:observe:job:> err: %s", logPrefix, err.Error())
			}

			js.provider, _ = jp.New(job.Spec.Provider)
			js.hook, _ = jh.New(job.Spec.Hook)

			break
		}

	}
}

func (js *JobState) SetJob(job *types.Job) {
	js.observers.job <- job
}

func (js *JobState) SetTask(task *types.Task) {
	js.observers.task <- task
}

func (js *JobState) DelTask(d *types.Task) {

	if _, ok := js.pod.list[d.SelfLink().String()]; !ok {
		return
	}
	delete(js.pod.list, d.SelfLink().String())

	if _, ok := js.task.list[d.SelfLink().String()]; !ok {
		return
	}

	delete(js.task.list, d.SelfLink().String())
	delete(js.task.queue, d.SelfLink().String())
	delete(js.task.active, d.SelfLink().String())

}

func (js *JobState) SetPod(pod *types.Pod) {
	js.observers.pod <- pod
}

func (js *JobState) DelPod(pod *types.Pod) {

	_, sl := pod.SelfLink().Parent()
	if _, ok := js.pod.list[sl.String()]; !ok {
		return
	}

	delete(js.pod.list, sl.String())
}

func (js *JobState) CheckJobDeps(dep types.StatusDependency) {

	log.Debugf("%s:> check job dependency: %s", logPrefix, dep.Name)
}

func (js *JobState) CheckTaskDeps(task *types.Task, dep types.StatusDependency) {

	log.Debugf("%s:> check dependency: %s", logPrefix, dep.Name)

	if task == nil {
		log.Debugf("%s:> check dependency: %s: provision task not found", logPrefix, dep.Name)
		return
	}

	if task.Status.State == types.StateWaiting {

		switch dep.Type {
		case types.KindVolume:
			if _, ok := task.Status.Dependencies.Volumes[dep.Name]; !ok {
				return
			}

			task.Status.Dependencies.Volumes[dep.Name] = dep
			if task.Status.CheckDeps() {
				task.Status.State = types.StateCreated
				js.observers.task <- task
			}
		case types.KindSecret:
			if _, ok := task.Status.Dependencies.Secrets[dep.Name]; !ok {
				return
			}

			task.Status.Dependencies.Secrets[dep.Name] = dep
			if task.Status.CheckDeps() {
				task.Status.State = types.StateCreated
				js.observers.task <- task
			}

		case types.KindConfig:
			if _, ok := task.Status.Dependencies.Configs[dep.Name]; !ok {
				return
			}

			task.Status.Dependencies.Configs[dep.Name] = dep
			if task.Status.CheckDeps() {
				task.Status.State = types.StateCreated
				js.observers.task <- task
			}
		}

	}
}

func (js *JobState) Provider() {

	if js.provider == nil {
		return
	}

	var (
		fetch = make(chan bool)
		limit = js.job.Spec.Concurrency.Limit
	)

	if limit == 0 {
		limit = 1
	}

	go func() {
		for {
			select {
			case _ = <-fetch:

				if js.provider == nil {
					return
				}

				manifest, err := js.provider.Fetch()

				if err != nil {
					log.Errorf("%s:> provider fetch err: %s", logPrefix, err.Error())
				}

				if manifest != nil && manifest.Spec.Template == nil && manifest.Spec.Runtime == nil {
					continue
				}

				if _, err := taskCreate(js.job, manifest); err != nil {
					log.Error(err.Error())
				}

			}
		}
	}()

	for {

		if js.provider == nil || js.job == nil {
			return
		}

		if len(js.task.active) < limit {
			fetch <- true
		}

		if js.job.Spec.Provider.Timeout == types.EmptyString {
			js.job.Spec.Provider.Timeout = "5s"
		}

		t, _ := time.ParseDuration(js.job.Spec.Provider.Timeout)

		if t < 1000000 {
			t = 1000000
		}

		<-time.NewTimer(t).C

		log.Debugf("%s:> provider timeout", logPrefix)
	}

}

func (js *JobState) Hook(task *types.Task) error {

	if js.hook != nil {
		if err := js.hook.Execute(task); err != nil {
			log.Errorf("%s:hook> execute err: %s", logPrefix, err.Error())
			return err
		}
	}

	return nil
}

func NewJobState(cs *cluster.ClusterState, job *types.Job) *JobState {

	var js = new(JobState)

	js.job = job
	js.cluster = cs

	js.observers.job = make(chan *types.Job)
	js.observers.task = make(chan *types.Task)
	js.observers.pod = make(chan *types.Pod)

	js.task.list = make(map[string]*types.Task)
	js.task.queue = make(map[string]*types.Task)
	js.task.active = make(map[string]*types.Task)
	js.task.finished = make([]*types.Task, 0)

	js.pod.list = make(map[string]*types.Pod)
	js.provider, _ = jp.New(job.Spec.Provider)
	js.hook, _ = jh.New(job.Spec.Hook)

	go js.Observe()

	return js
}
