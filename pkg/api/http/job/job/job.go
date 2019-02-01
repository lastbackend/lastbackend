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

package job

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/api/envs"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/distribution"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"net/http"
)

const (
	logPrefix = "api:handler:job"
	logLevel  = 3
)

func Fetch(ctx context.Context, namespace, name string) (*types.Job, *errors.Err) {
	jm := distribution.NewJobModel(ctx, envs.Get().GetStorage())
	job, err := jm.Get(distribution.JobSelfLink(namespace, name))

	if err != nil {
		log.V(logLevel).Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
		return nil, errors.New("service").InternalServerError(err)
	}

	if job == nil {
		err := errors.New("job not found")
		log.V(logLevel).Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
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

	return Update(ctx, job, mf)
}

func Create(ctx context.Context, ns *types.Namespace, mf *request.JobManifest) (*types.Job, *errors.Err) {

	jm := distribution.NewJobModel(ctx, envs.Get().GetStorage())

	if mf.Meta.Name != nil {

		job, err := jm.Get(new(types.Job).CreateSelfLink(ns.Meta.Name, *mf.Meta.Name))
		if err != nil {
			log.V(logLevel).Errorf("%s:create:> get service by name `%s` in namespace `%s` err: %s", logPrefix, mf.Meta.Name, ns.Meta.Name, err.Error())
			return nil, errors.New("service").InternalServerError()

		}

		if job != nil {
			log.V(logLevel).Warnf("%s:create:> service name `%s` in namespace `%s` not unique", logPrefix, mf.Meta.Name, ns.Meta.Name)
			return nil, errors.New("service").NotUnique("name")

		}
	}

	job := new(types.Job)
	mf.SetJobMeta(job)
	job.Meta.Namespace = ns.Meta.Name

	if err := mf.SetJobSpec(job); err != nil {
		return nil, errors.New("service").BadRequest(err.Error())
	}

	job, err := jm.Create(job)
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> create service err: %s", logPrefix, err.Error())
		return nil, errors.New("service").InternalServerError()
	}

	return job, nil
}

func Update(ctx context.Context, job *types.Job, mf *request.JobManifest) (*types.Job, *errors.Err) {

	sm := distribution.NewJobModel(ctx, envs.Get().GetStorage())
	mf.SetJobMeta(job)
	if err := sm.Update(job); err != nil {
		log.V(logLevel).Errorf("%s:update:> update service err: %s", logPrefix, err.Error())
		return nil, errors.New("service").InternalServerError()
	}

	return job, nil
}
