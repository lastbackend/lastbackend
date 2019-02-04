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

package views

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/util/resource"
)

type JobView struct{}

func (jw *JobView) New(obj *types.Job, tasks *types.TaskList, pods *types.PodList) *Job {
	j := Job{}

	j.SetMeta(obj.Meta)
	j.SetStatus(obj.Status)
	j.SetSpec(obj.Spec)

	if tasks != nil {
		j.Tasks = make(TaskList, 0)
		j.JoinTasks(tasks, pods)
	}

	return &j
}

func (j *Job) SetMeta(obj types.JobMeta) {

	jm := JobMeta{}

	jm.Namespace = obj.Namespace
	jm.Name = obj.Name

	jm.SelfLink = obj.SelfLink
	jm.Description = obj.Description

	jm.Labels = obj.Labels
	jm.Created = obj.Created
	jm.Updated = obj.Updated

	j.Meta = jm
}

func (j *Job) SetStatus(obj types.JobStatus) {
	js := JobStatus{
		State:   obj.State,
		Message: obj.Message,
		Stats: JobStatusStats{
			Total:        obj.Stats.Total,
			Active:       obj.Stats.Active,
			Failed:       obj.Stats.Failed,
			Successful:   obj.Stats.Successful,
			LastSchedule: obj.Stats.LastSchedule,
		},
		Resources: JobStatusResources{
			Allocated: JobResource{
				RAM:     resource.EncodeMemoryResource(obj.Resources.Allocated.RAM),
				CPU:     resource.EncodeCpuResource(obj.Resources.Allocated.CPU),
				Storage: resource.EncodeMemoryResource(obj.Resources.Allocated.Storage),
			},
		},
		Updated: obj.Updated,
	}

	j.Status = js
}

func (j *Job) SetSpec(obj types.JobSpec) {
	mv := new(ManifestView)
	js := JobSpec{
		Enabled:  obj.Enabled,
		Schedule: obj.Schedule,
		Concurrency: JobSpecConcurrency{
			Limit:    obj.Concurrency.Limit,
			Strategy: obj.Concurrency.Strategy,
		},
		Remote: JobSpecRemote{
			Request: JobSpecRemoteRequest{
				Endpoint: obj.Remote.Request.Endpoint,
				Method:   obj.Remote.Request.Method,
				Headers:  obj.Remote.Request.Headers,
			},
			Response: JobSpecRemoteRequest{
				Endpoint: obj.Remote.Response.Endpoint,
				Method:   obj.Remote.Response.Method,
				Headers:  obj.Remote.Response.Headers,
			},
		},
		Resources: JobResources{
			Request: JobResource{
				RAM:     resource.EncodeMemoryResource(obj.Resources.Request.RAM),
				CPU:     resource.EncodeCpuResource(obj.Resources.Request.CPU),
				Storage: resource.EncodeMemoryResource(obj.Resources.Request.Storage),
			},
			Limits: JobResource{
				RAM:     resource.EncodeMemoryResource(obj.Resources.Limits.RAM),
				CPU:     resource.EncodeCpuResource(obj.Resources.Limits.CPU),
				Storage: resource.EncodeMemoryResource(obj.Resources.Limits.Storage),
			},
		},
		Task: JobSpecTask{
			Selector: mv.NewManifestSpecSelector(obj.Task.Selector),
			Runtime:  mv.NewManifestSpecRuntime(obj.Task.Runtime),
			Template: mv.NewManifestSpecTemplate(obj.Task.Template),
		},
	}
	j.Spec = js
}

func (j *Job) JoinTasks(tasks *types.TaskList, pods *types.PodList) {

	for _, t := range tasks.Items {

		if t.Meta.Namespace != j.Meta.Namespace {
			continue
		}

		if t.Meta.Job != j.Meta.Name {
			continue
		}

		j.Tasks = append(j.Tasks, new(TaskView).New(t, pods))
	}
}

func (jw *JobView) NewList(obj *types.JobList, tasks *types.TaskList, pods *types.PodList) *JobList {

	if obj == nil {
		return nil
	}

	jl := make(JobList, 0)
	for _, v := range obj.Items {
		jl = append(jl, jw.New(v, tasks, pods))
	}

	return &jl
}
