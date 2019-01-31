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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/log"

	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/storage"
)

const (
	logJobPrefix = "distribution:task"
)

type Task struct {
	context context.Context
	storage storage.Storage
}

func (t *Task) Runtime() (*types.System, error) {

	log.V(logLevel).Debugf("%s:get:> get task runtime info", logJobPrefix)
	runtime, err := t.storage.Info(t.context, t.storage.Collection().Task(), "")
	if err != nil {
		log.V(logLevel).Errorf("%s:get:> get runtime info error: %s", logJobPrefix, err)
		return &runtime.System, err
	}
	return &runtime.System, nil
}

func (t *Task) Get(selflink string) (*types.Task, error) {

	log.V(logLevel).Debugf("%s:get:> get by selflink %s", logJobPrefix, selflink)

	task := new(types.Task)
	err := t.storage.Get(t.context, t.storage.Collection().Task(), selflink, task, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> get task by selflink %s not found", logJobPrefix, selflink)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> get task by selflink %s error: %v", logJobPrefix, selflink, err)
		return nil, err
	}

	return task, nil
}

func (t *Task) ListByNamespace(namespace string) (*types.TaskList, error) {
	log.V(logLevel).Debugf("%s:list:> by namespace %s", logJobPrefix, namespace)
	tasks := types.NewTaskList()

	q := t.storage.Filter().Task().ByNamespace(namespace)
	err := t.storage.List(t.context, t.storage.Collection().Task(), q, tasks, nil)
	if err != nil {
		log.V(logLevel).Error("%s:list:> by namespace %s err: %v", logJobPrefix, namespace, err)
		return nil, err
	}

	log.V(logLevel).Debugf("%s:list:> by namespace %s result: %d", logJobPrefix, namespace, len(tasks.Items))

	return tasks, nil
}

func (t *Task) ListByJob(namespace, job string) (*types.TaskList, error) {
	log.V(logLevel).Debugf("%s:list:> by namespace %s", logJobPrefix, namespace)
	tasks := types.NewTaskList()

	q := t.storage.Filter().Task().ByJob(namespace, job)
	err := t.storage.List(t.context, t.storage.Collection().Task(), q, tasks, nil)
	if err != nil {
		log.V(logLevel).Error("%s:list:> by namespace %s err: %v", logJobPrefix, namespace, err)
		return nil, err
	}

	log.V(logLevel).Debugf("%s:list:> by namespace %s result: %d", logJobPrefix, namespace, len(tasks.Items))

	return tasks, nil
}

func (t *Task) Create(task *types.Task) (*types.Task, error) {

	if err := t.storage.Put(t.context, t.storage.Collection().Task(),
		task.SelfLink(), task, nil); err != nil {
		log.Errorf("%s:create:> task %s create err: %v", logJobPrefix, task.Meta.SelfLink, err)
		return nil, err
	}

	return task, nil
}

// Update task
func (t *Task) Cancel(task *types.Task) error {

	log.V(logLevel).Debugf("%s:cancel:> cancel task %s", logJobPrefix, task.Meta.Name)

	// mark task for destroy
	task.Spec.State.Cancel = true
	// mark task for cancel
	task.Status.SetCancel()

	if err := t.storage.Set(t.context, t.storage.Collection().Task(),
		task.SelfLink(), task, nil); err != nil {
		log.V(logLevel).Debugf("%s:destroy: destroy task %s err: %v", logJobPrefix, task.Meta.Name, err)
		return err
	}

	return nil
}

// Destroy task
func (t *Task) Destroy(task *types.Task) error {

	log.V(logLevel).Debugf("%s:destroy:> destroy task %s", logJobPrefix, task.Meta.Name)

	// mark task for destroy
	task.Spec.State.Destroy = true
	// mark task for destroy
	task.Status.SetDestroy()

	if err := t.storage.Set(t.context, t.storage.Collection().Task(),
		task.SelfLink(), task, nil); err != nil {
		log.V(logLevel).Debugf("%s:destroy:> destroy task %s err: %v", logJobPrefix, task.Meta.Name, err)
		return err
	}

	return nil
}

// Remove task
func (t *Task) Remove(task *types.Task) error {

	log.V(logLevel).Debugf("%s:remove:> remove task %s", logJobPrefix, task.Meta.Name)
	if err := t.storage.Del(t.context, t.storage.Collection().Task(),
		task.SelfLink()); err != nil {
		log.V(logLevel).Debugf("%s:remove:> remove task %s err: %v", logJobPrefix, task.Meta.Name, err)
		return err
	}

	return nil
}

// Watch task changes
func (t *Task) Watch(dt chan types.JobEvent, rev *int64) error {

	done := make(chan bool)
	watcher := storage.NewWatcher()

	log.V(logLevel).Debugf("%s:watch:> watch tasks", logJobPrefix)

	go func() {
		for {
			select {
			case <-t.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := types.JobEvent{}
				res.Action = e.Action
				res.Name = e.Name

				task := new(types.Task)

				if err := json.Unmarshal(e.Data.([]byte), task); err != nil {
					log.Errorf("%s:> parse data err: %v", logJobPrefix, err)
					continue
				}

				res.Data = task

				dt <- res
			}
		}
	}()

	opts := storage.GetOpts()
	opts.Rev = rev
	if err := t.storage.Watch(t.context, t.storage.Collection().Task(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewTaskModel(ctx context.Context, stg storage.Storage) *Task {
	return &Task{ctx, stg}
}

func TaskSelfLink(namespace, job, name string) string {
	return fmt.Sprintf("%s:%s:%s", namespace, job, name)
}
