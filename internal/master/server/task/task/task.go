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

package task

import (
	"context"
	"github.com/lastbackend/lastbackend/internal/api/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/model"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/internal/util/generator"
	"github.com/lastbackend/lastbackend/internal/util/resource"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/tools/log"
	"strings"
)

const (
	logPrefix = "api:handler:task"
	logLevel  = 3
)

func Fetch(ctx context.Context, namespace, job, name string) (*types.Task, *errors.Err) {

	tm := model.NewTaskModel(ctx, envs.Get().GetStorage())
	task, err := tm.Get(types.NewTaskSelfLink(namespace, job, name).String())

	if err != nil {
		log.Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
		return nil, errors.New("task").InternalServerError(err)
	}

	if task == nil {
		err := errors.New("task not found")
		log.Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
		return nil, errors.New("task").NotFound()
	}

	return task, nil
}

func Create(ctx context.Context, ns *types.Namespace, job *types.Job, mf *request.TaskManifest) (*types.Task, *errors.Err) {

	jm := model.NewJobModel(ctx, envs.Get().GetStorage())
	tm := model.NewTaskModel(ctx, envs.Get().GetStorage())

	if mf.Meta.Name != nil {

		task, err := tm.Get(types.NewTaskSelfLink(ns.Meta.Name, job.Meta.Name, *mf.Meta.Name).String())
		if err != nil {
			log.Errorf("%s:create:> get task by name `%s` in namespace `%s` err: %s", logPrefix, mf.Meta.Name, ns.Meta.Name, err.Error())
			return nil, errors.New("task").InternalServerError()

		}

		if task != nil {
			log.Warnf("%s:create:> task name `%s` in namespace `%s` not unique", logPrefix, mf.Meta.Name, ns.Meta.Name)
			return nil, errors.New("job").NotUnique("name")

		}
	}

	task := new(types.Task)
	task.Meta.SetDefault()
	task.Meta.Namespace = ns.Meta.Name
	task.Meta.Job = job.Meta.Name

	if mf.Meta.Name != nil {
		task.Meta.SelfLink = *types.NewTaskSelfLink(ns.Meta.Name, job.Meta.Name, *mf.Meta.Name)
		mf.SetTaskMeta(task)
	} else {
		name := strings.Split(generator.GetUUIDV4(), "-")[4][5:]
		task.Meta.Name = name
		task.Meta.SelfLink = *types.NewTaskSelfLink(ns.Meta.Name, job.Meta.Name, name)
	}

	task.Status.State = types.StateCreated

	task.Spec.Runtime = job.Spec.Task.Runtime
	task.Spec.Selector = job.Spec.Task.Selector
	task.Spec.Template = job.Spec.Task.Template

	if err := mf.SetTaskSpec(task); err != nil {
		log.Errorf("%s:create:> set task spec err: %s", logPrefix, err.Error())
		return nil, errors.New("task").BadParameter("spec")
	}

	if job.Spec.Resources.Limits.RAM != 0 || job.Spec.Resources.Limits.CPU != 0 {
		for _, c := range task.Spec.Template.Containers {
			if c.Resources.Limits.RAM == 0 {
				c.Resources.Limits.RAM, _ = resource.DecodeMemoryResource(types.DEFAULT_RESOURCE_LIMITS_RAM)
			}
			if c.Resources.Limits.CPU == 0 {
				c.Resources.Limits.CPU, _ = resource.DecodeCpuResource(types.DEFAULT_RESOURCE_LIMITS_CPU)
			}
		}
	}

	if err := job.AllocateResources(task.Spec.GetResourceRequest()); err != nil {
		log.Errorf("%s:create:> %s", logPrefix, err.Error())
		return nil, errors.New("job").BadRequest(err.Error())
	} else {
		if err := jm.Set(job); err != nil {
			log.Errorf("%s:update:> update namespace err: %s", logPrefix, err.Error())
			return nil, errors.New("job").InternalServerError()
		}
	}

	if _, err := tm.Create(task); err != nil {
		log.Errorf("%s:create:> create task err: %s", logPrefix, err.Error())
		return nil, errors.New("task").InternalServerError()
	}

	return task, nil
}
