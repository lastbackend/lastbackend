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

package views

import (
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/util/resource"
)

type JobView struct{}

func (jw *JobView) New(obj *models.Job) *Job {
	j := Job{}

	j.SetMeta(obj.Meta)
	j.SetStatus(obj.Status)
	j.SetSpec(obj.Spec)

	return &j
}

func (j *Job) SetMeta(obj models.JobMeta) {

	jm := JobMeta{}

	jm.Namespace = obj.Namespace
	jm.Name = obj.Name

	jm.SelfLink = obj.SelfLink.String()
	jm.Description = obj.Description

	jm.Labels = obj.Labels
	jm.Created = obj.Created
	jm.Updated = obj.Updated

	j.Meta = jm
}

func (j *Job) SetStatus(obj models.JobStatus) {
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

func (j *Job) SetSpec(obj models.JobSpec) {
	mv := new(ManifestView)
	js := JobSpec{
		Enabled: obj.Enabled,
		Concurrency: JobSpecConcurrency{
			Limit:    obj.Concurrency.Limit,
			Strategy: obj.Concurrency.Strategy,
		},
		Provider: JobSpecProvider{
			Timeout: obj.Provider.Timeout,
		},
		Hook: JobSpecHook{},
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

	if obj.Provider.Http != nil {
		js.Provider.Http = &JobSpecProviderHTTP{
			Endpoint: obj.Provider.Http.Endpoint,
			Method:   obj.Provider.Http.Method,
			Headers:  obj.Provider.Http.Headers,
		}
	}

	if obj.Hook.Http != nil {
		js.Hook.Http = &JobSpecHookHTTP{
			Endpoint: obj.Hook.Http.Endpoint,
			Method:   obj.Hook.Http.Method,
			Headers:  obj.Hook.Http.Headers,
		}
	}

	j.Spec = js
}

func (jw *JobView) NewList(obj *models.JobList) *JobList {

	if obj == nil {
		return nil
	}

	jl := make(JobList, 0)
	for _, v := range obj.Items {
		jl = append(jl, jw.New(v))
	}

	return &jl
}
