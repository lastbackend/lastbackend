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

package service

import (
	"context"
	"encoding/json"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/pkg/storage"
	"github.com/lastbackend/lastbackend/tools/log"
)

const (
	logTaskPrefix = "distribution:task"
)

type Task struct {
	context context.Context
	storage storage.IStorage
}

func (t *Task) Runtime() (*models.System, error) {

	log.Debugf("%s:get:> get task runtime info", logTaskPrefix)
	runtime, err := t.storage.Info(t.context, t.storage.Collection().Task(), "")
	if err != nil {
		log.Errorf("%s:get:> get runtime info error: %s", logTaskPrefix, err)
		return &runtime.System, err
	}
	return &runtime.System, nil
}

func (t *Task) Get(selflink string) (*models.Task, error) {

	log.Debugf("%s:get:> get by selflink %s", logTaskPrefix, selflink)

	task := new(models.Task)
	err := t.storage.Get(t.context, t.storage.Collection().Task(), selflink, task, nil)
	if err != nil {

		if errors.Storage().IsErrEntityNotFound(err) {
			log.Warnf("%s:get:> get task by selflink %s not found", logTaskPrefix, selflink)
			return nil, nil
		}

		log.Errorf("%s:get:> get task by selflink %s error: %v", logTaskPrefix, selflink, err)
		return nil, err
	}

	return task, nil
}

func (t *Task) ListByNamespace(namespace string) (*models.TaskList, error) {
	log.Debugf("%s:list:> by namespace %s", logTaskPrefix, namespace)
	tasks := models.NewTaskList()

	q := t.storage.Filter().Task().ByNamespace(namespace)
	err := t.storage.List(t.context, t.storage.Collection().Task(), q, tasks, nil)
	if err != nil {
		log.Error("%s:list:> by namespace %s err: %v", logTaskPrefix, namespace, err)
		return nil, err
	}

	log.Debugf("%s:list:> by namespace %s result: %d", logTaskPrefix, namespace, len(tasks.Items))

	return tasks, nil
}

func (t *Task) ListByJob(namespace, job string) (*models.TaskList, error) {
	log.Debugf("%s:list:> by namespace %s", logTaskPrefix, namespace)
	tasks := models.NewTaskList()

	q := t.storage.Filter().Task().ByJob(namespace, job)
	err := t.storage.List(t.context, t.storage.Collection().Task(), q, tasks, nil)
	if err != nil {
		log.Error("%s:list:> by namespace %s err: %v", logTaskPrefix, namespace, err)
		return nil, err
	}

	log.Debugf("%s:list:> by namespace %s result: %d", logTaskPrefix, namespace, len(tasks.Items))

	return tasks, nil
}

func (t *Task) Create(task *models.Task) (*models.Task, error) {

	task.Status.State = models.StateCreated

	if err := t.storage.Put(t.context, t.storage.Collection().Task(),
		task.SelfLink().String(), task, nil); err != nil {
		log.Errorf("%s:create:> task %s create err: %v", logTaskPrefix, task.Meta.SelfLink.String(), err)
		return nil, err
	}

	return task, nil
}

// Cancel task
func (t *Task) Cancel(task *models.Task) error {

	log.Debugf("%s:cancel:> cancel task %s", logTaskPrefix, task.Meta.Name)

	// mark task for destroy
	task.Spec.State.Cancel = true
	// mark task for cancel
	task.Status.State = models.StateCanceled

	if err := t.storage.Set(t.context, t.storage.Collection().Task(),
		task.SelfLink().String(), task, nil); err != nil {
		log.Debugf("%s:destroy: destroy task %s err: %v", logTaskPrefix, task.Meta.Name, err)
		return err
	}

	return nil
}

// Update task
func (t *Task) Set(task *models.Task) error {

	log.Debugf("%s:set:> set task %s", logTaskPrefix, task.Meta.Name)

	if err := t.storage.Set(t.context, t.storage.Collection().Task(),
		task.SelfLink().String(), task, nil); err != nil {
		log.Debugf("%s:set: set task %s err: %v", logTaskPrefix, task.Meta.Name, err)
		return err
	}

	return nil
}

// Destroy task
func (t *Task) Destroy(task *models.Task) error {

	log.Debugf("%s:destroy:> destroy task %s", logTaskPrefix, task.Meta.Name)

	// mark task for destroy
	task.Spec.State.Destroy = true
	// mark task for destroy
	task.Status.State = models.StateDestroyed

	if err := t.storage.Set(t.context, t.storage.Collection().Task(),
		task.SelfLink().String(), task, nil); err != nil {
		log.Debugf("%s:destroy:> destroy task %s err: %v", logTaskPrefix, task.Meta.Name, err)
		return err
	}

	return nil
}

// Remove task
func (t *Task) Remove(task *models.Task) error {

	log.Debugf("%s:remove:> remove task %s", logTaskPrefix, task.Meta.Name)
	if err := t.storage.Del(t.context, t.storage.Collection().Task(),
		task.SelfLink().String()); err != nil {
		log.Debugf("%s:remove:> remove task %s err: %v", logTaskPrefix, task.Meta.Name, err)
		return err
	}

	return nil
}

// Watch task changes
func (t *Task) Watch(dt chan models.TaskEvent, rev *int64) error {

	done := make(chan bool)
	watcher := storage.NewWatcher()

	log.Debugf("%s:watch:> watch tasks", logTaskPrefix)

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

				res := models.TaskEvent{}
				res.Action = e.Action
				res.Name = e.Name

				task := new(models.Task)

				if err := json.Unmarshal(e.Data.([]byte), task); err != nil {
					log.Errorf("%s:> parse data err: %v", logTaskPrefix, err)
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

func NewTaskModel(ctx context.Context, stg storage.IStorage) *Task {
	return &Task{ctx, stg}
}
