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
	"sync"
	"time"

	"github.com/lastbackend/lastbackend/internal/master/state/cluster"
	h "github.com/lastbackend/lastbackend/internal/master/state/job/hook"
	"github.com/lastbackend/lastbackend/internal/master/state/job/hook/hook"
	p "github.com/lastbackend/lastbackend/internal/master/state/job/provider"
	"github.com/lastbackend/lastbackend/internal/master/state/job/provider/provider"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logLevel  = 3
	logPrefix = "state:job"
)

type JobState struct {
	lock sync.Mutex

	storage storage.IStorage

	cluster *cluster.ClusterState
	job     *models.Job

	task struct {
		active   map[string]*models.Task
		queue    map[string]*models.Task
		list     map[string]*models.Task
		finished []*models.Task
	}

	pod struct {
		list map[string]*models.Pod
	}

	observers struct {
		job  chan *models.Job
		task chan *models.Task
		pod  chan *models.Pod
	}

	provider p.JobProvider
	hook     h.Hook
}

type JobTaskState struct {
	Active   int
	Queue    int
	List     int
	Finished int
}

type JobPodState struct {
	List int
}

func (js *JobState) Namespace() string {
	return js.job.Meta.Namespace
}

func (js *JobState) Restore() error {

	log.V(logLevel).Debugf("%s:restore state for job: %s", logPrefix, js.job.SelfLink())

	var (
		err error
	)

	// Get all pods
	pm := service.NewPodModel(context.Background(), js.storage)
	pl, err := pm.ListByJob(js.job.Meta.Namespace, js.job.Meta.Name)
	if err != nil {
		log.Errorf("%s:restore:> get pod map error: %v", logPrefix, err)
		return err
	}

	js.lock.Lock()
	for _, pod := range pl.Items {
		log.Infof("%s: restore: restore pod: %s", logPrefix, pod.SelfLink())

		// Check if task map for pod exists
		_, sl := pod.SelfLink().Parent()
		// put pod into map by task name and pod name
		js.pod.list[sl.String()] = pod
	}
	js.lock.Unlock()

	// Get all tasks
	tm := service.NewTaskModel(context.Background(), js.storage)
	tl, err := tm.ListByJob(js.job.Meta.Namespace, js.job.Meta.Name)
	if err != nil {
		log.Errorf("%s:restore:> get task map error: %v", logPrefix, err)
		return err
	}

	js.lock.Lock()
	for _, task := range tl.Items {
		log.Infof("%s: restore task: %s", logPrefix, task.SelfLink())
		js.task.list[task.SelfLink().String()] = task
	}
	js.lock.Unlock()

	// Range over pods to sync pod status
	for _, pod := range js.pod.list {
		js.observers.pod <- pod
	}

	// Range over tasks to sync tasks status
	for _, task := range js.task.list {
		js.observers.task <- task
	}

	// Sync job state if updated
	js.SetJob(js.job)

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
				log.V(logLevel).Errorf("%s:observe:pod:> err: %s", logPrefix, err.Error())
				break
			}
		case task := <-js.observers.task:

			log.V(logLevel).Debugf("%s:observe:task:> %s (%s)", logPrefix, task.SelfLink(), task.Status.State)

			if err := taskObserve(js, task); err != nil {
				log.V(logLevel).Errorf("%s:observe:task err:> %s", logPrefix, err.Error())
				break
			}
		case job := <-js.observers.job:
			log.V(logLevel).Debugf("%s:observe:job:> %s", logPrefix, job.SelfLink())

			js.job = job

			if err := jobObserve(js, job); err != nil {
				log.V(logLevel).Errorf("%s:observe:job:> err: %s", logPrefix, err.Error())
				break
			}

			js.provider, _ = provider.New(job.Spec.Provider)
			js.hook, _ = hook.New(job.Spec.Hook)
		}

	}
}

func (js *JobState) SetJob(job *models.Job) {
	js.observers.job <- job
}

func (js *JobState) SetTask(task *models.Task) {
	js.observers.task <- task
}

func (js *JobState) DelTask(t *models.Task) {
	js.lock.Lock()
	delete(js.task.list, t.SelfLink().String())
	delete(js.task.queue, t.SelfLink().String())
	delete(js.task.active, t.SelfLink().String())
	delete(js.pod.list, t.SelfLink().String())
	js.lock.Unlock()
}

func (js *JobState) SetPod(pod *models.Pod) {
	js.observers.pod <- pod
}

func (js *JobState) DelPod(pod *models.Pod) {

	_, sl := pod.SelfLink().Parent()
	if _, ok := js.pod.list[sl.String()]; !ok {
		return
	}

	js.lock.Lock()
	delete(js.pod.list, sl.String())
	js.lock.Unlock()
}

func (js *JobState) CheckJobDeps(dep models.StatusDependency) {
	log.Debugf("%s:> check job dependency: %s", logPrefix, dep.Name)
}

func (js *JobState) CheckTaskDeps(task *models.Task, dep models.StatusDependency) {

	log.Debugf("%s:> check dependency: %s", logPrefix, dep.Name)

	if task == nil {
		log.Debugf("%s:> check dependency: %s: provision task not found", logPrefix, dep.Name)
		return
	}

	if task.Status.State == models.StateWaiting {

		switch dep.Type {
		case models.KindVolume:
			if _, ok := task.Status.Dependencies.Volumes[dep.Name]; !ok {
				return
			}

			task.Status.Dependencies.Volumes[dep.Name] = dep
			if task.Status.CheckDeps() {
				task.Status.State = models.StateCreated
				js.observers.task <- task
			}
		case models.KindSecret:
			if _, ok := task.Status.Dependencies.Secrets[dep.Name]; !ok {
				return
			}

			task.Status.Dependencies.Secrets[dep.Name] = dep
			if task.Status.CheckDeps() {
				task.Status.State = models.StateCreated
				js.observers.task <- task
			}

		case models.KindConfig:
			if _, ok := task.Status.Dependencies.Configs[dep.Name]; !ok {
				return
			}

			task.Status.Dependencies.Configs[dep.Name] = dep
			if task.Status.CheckDeps() {
				task.Status.State = models.StateCreated
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
			case <-fetch:

				if js.provider == nil {
					return
				}

				manifest, err := js.provider.Fetch()
				if err != nil {
					log.Errorf("%s:> provider fetch err: %v", logPrefix, err)
					continue
				}

				if manifest != nil && manifest.Spec.Template == nil && manifest.Spec.Runtime == nil {
					continue
				}

				task, err := taskCreate(js.storage, js.job, manifest)
				if err != nil {
					log.Errorf("%s:> create task err: %v", logPrefix, err)
					continue
				}

				js.task.list[task.SelfLink().String()] = task

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

		if js.job.Spec.Provider.Timeout == models.EmptyString {
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

func (js *JobState) Hook(task *models.Task) error {

	if js.hook != nil {
		if err := js.hook.Execute(task); err != nil {
			log.Errorf("%s:hook> execute err: %s", logPrefix, err.Error())
			return err
		}
	}

	return nil
}

func NewJobState(stg storage.IStorage, cs *cluster.ClusterState, job *models.Job) *JobState {

	var js = new(JobState)

	js.storage = stg
	js.job = job
	js.cluster = cs

	js.observers.job = make(chan *models.Job)
	js.observers.task = make(chan *models.Task)
	js.observers.pod = make(chan *models.Pod)

	js.task.list = make(map[string]*models.Task, 0)
	js.task.queue = make(map[string]*models.Task, 0)
	js.task.active = make(map[string]*models.Task, 0)
	js.task.finished = make([]*models.Task, 0)

	js.pod.list = make(map[string]*models.Pod, 0)
	js.provider, _ = provider.New(job.Spec.Provider)
	js.hook, _ = hook.New(job.Spec.Hook)

	go js.Observe()

	return js
}
