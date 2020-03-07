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
	"github.com/lastbackend/lastbackend/internal/api/envs"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/pkg/model"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/internal/util/resource"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/tools/log"
	"net/http"
)

const (
	logPrefix = "api:handler:job"
	logLevel  = 3
)

func Fetch(ctx context.Context, namespace, name string) (*types.Job, *errors.Err) {
	jm := model.NewJobModel(ctx, envs.Get().GetStorage())
	job, err := jm.Get(types.NewJobSelfLink(namespace, name).String())

	if err != nil {
		log.Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
		return nil, errors.New("job").InternalServerError(err)
	}

	if job == nil {
		err := errors.New("job not found")
		log.Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
		return nil, errors.New("job").NotFound()
	}

	return job, nil
}

func Apply(ctx context.Context, ns *types.Namespace, mf *request.JobManifest) (*types.Job, *errors.Err) {

	if mf.Meta.Name == nil {
		return nil, errors.New("job").BadParameter("meta.name")
	}

	job, err := Fetch(ctx, ns.Meta.Name, *mf.Meta.Name)
	if err != nil {
		if err.Code != http.StatusText(http.StatusNotFound) {
			return nil, errors.New("job").InternalServerError()
		}
	}

	if job == nil {
		return Create(ctx, ns, mf)
	}

	return Update(ctx, ns, job, mf)
}

func Create(ctx context.Context, ns *types.Namespace, mf *request.JobManifest) (*types.Job, *errors.Err) {

	jm := model.NewJobModel(ctx, envs.Get().GetStorage())
	nm := model.NewNamespaceModel(ctx, envs.Get().GetStorage())

	if mf.Meta.Name != nil {

		job, err := jm.Get(types.NewJobSelfLink(ns.Meta.Name, *mf.Meta.Name).String())
		if err != nil {
			log.Errorf("%s:create:> get job by name `%s` in namespace `%s` err: %s", logPrefix, mf.Meta.Name, ns.Meta.Name, err.Error())
			return nil, errors.New("job").InternalServerError()

		}

		if job != nil {
			log.Warnf("%s:create:> job name `%s` in namespace `%s` not unique", logPrefix, mf.Meta.Name, ns.Meta.Name)
			return nil, errors.New("job").NotUnique("name")

		}
	}

	job := new(types.Job)
	mf.SetJobMeta(job)
	job.Meta.SelfLink = *types.NewJobSelfLink(ns.Meta.Name, *mf.Meta.Name)
	job.Meta.Namespace = ns.Meta.Name
	job.Status.State = types.StateCreated

	if err := mf.SetJobSpec(job); err != nil {
		return nil, errors.New("job").BadRequest(err.Error())
	}

	if ns.Spec.Resources.Limits.RAM != 0 || ns.Spec.Resources.Limits.CPU != 0 {
		for _, c := range job.Spec.Task.Template.Containers {
			if c.Resources.Limits.RAM == 0 {
				c.Resources.Limits.RAM, _ = resource.DecodeMemoryResource(types.DEFAULT_RESOURCE_LIMITS_RAM)
			}
			if c.Resources.Limits.CPU == 0 {
				c.Resources.Limits.CPU, _ = resource.DecodeCpuResource(types.DEFAULT_RESOURCE_LIMITS_CPU)
			}
		}
	}

	if err := ns.AllocateResources(job.Spec.GetResourceRequest()); err != nil {
		log.Errorf("%s:create:> %s", logPrefix, err.Error())
		return nil, errors.New("job").BadRequest(err.Error())

	} else {
		if err := nm.Update(ns); err != nil {
			log.Errorf("%s:create:> update namespace err: %s", logPrefix, err.Error())
			return nil, errors.New("job").InternalServerError()
		}
	}

	job, err := jm.Create(job)
	if err != nil {
		log.Errorf("%s:create:> create job err: %s", logPrefix, err.Error())
		return nil, errors.New("job").InternalServerError()
	}

	return job, nil
}

func Update(ctx context.Context, ns *types.Namespace, job *types.Job, mf *request.JobManifest) (*types.Job, *errors.Err) {

	jm := model.NewJobModel(ctx, envs.Get().GetStorage())
	nm := model.NewNamespaceModel(ctx, envs.Get().GetStorage())

	resources := job.Spec.GetResourceRequest()

	mf.SetJobMeta(job)
	if err := mf.SetJobSpec(job); err != nil {
		return nil, errors.New("job").BadRequest(err.Error())
	}

	requestedResources := job.Spec.GetResourceRequest()
	if !resources.Equal(requestedResources) {
		allocatedResources := ns.Status.Resources.Allocated
		ns.ReleaseResources(resources)

		if err := ns.AllocateResources(job.Spec.GetResourceRequest()); err != nil {
			ns.Status.Resources.Allocated = allocatedResources
			log.Errorf("%s:update:> %s", logPrefix, err.Error())
			return nil, errors.New("job").BadRequest(err.Error())
		} else {
			if err := nm.Update(ns); err != nil {
				log.Errorf("%s:update:> update namespace err: %s", logPrefix, err.Error())
				return nil, errors.New("job").InternalServerError()
			}

		}
	}

	if err := jm.Set(job); err != nil {
		log.Errorf("%s:update:> update job err: %s", logPrefix, err.Error())
		return nil, errors.New("job").InternalServerError()
	}

	return job, nil
}
