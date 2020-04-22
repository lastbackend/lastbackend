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

package v1

import (
	"context"
	"fmt"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/client/cluster/types"
	"strconv"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	t "github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/internal/util/http/request"
)

type TaskClient struct {
	client *request.RESTClient

	namespace t.NamespaceSelfLink
	job       t.JobSelfLink
	selflink  t.TaskSelfLink
}

func (tc *TaskClient) Pod(args ...string) types.PodClientV1 {
	name := ""
	// Get any parameters passed to us out of the args variable into "real"
	// variables we created for them.
	for i := range args {
		switch i {
		case 0: // hostname
			name = args[0]
		default:
			panic("Wrong parameter count: (is allowed from 0 to 1)")
		}
	}
	return newPodClient(tc.client, tc.namespace.String(), t.KindTask, tc.selflink.String(), name)
}

func (tc *TaskClient) List(ctx context.Context) (*views.TaskList, error) {

	var s *views.TaskList
	var e *errors.Http

	err := tc.client.Get(fmt.Sprintf("/namespace/%s/job/%s/task", tc.namespace.String(), tc.job.Name())).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(views.TaskList, 0)
		s = &list
	}

	return s, nil
}

func (tc *TaskClient) Create(ctx context.Context, opts *rv1.TaskManifest) (*views.Task, error) {
	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *views.Task
	var e *errors.Http

	err = tc.client.Post(fmt.Sprintf("/namespace/%s/job/%s/task", tc.namespace.String(), tc.job.Name())).
		AddHeader("Content-Type", "application/json").
		Body(body).
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	return s, nil
}

func (tc *TaskClient) Get(ctx context.Context) (*views.Task, error) {

	var s *views.Task
	var e *errors.Http

	err := tc.client.Get(fmt.Sprintf("/namespace/%s/job/%s/task/%s", tc.namespace.String(), tc.job.Name(), tc.selflink.Name())).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	return s, nil
}

func (tc *TaskClient) Cancel(ctx context.Context, opts *rv1.TaskCancelOptions) (*views.Task, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *views.Task
	var e *errors.Http

	err = tc.client.Post(fmt.Sprintf("/namespace/%s/job/%s/task/%s", tc.namespace.String(), tc.job.Name(), tc.selflink.Name())).
		AddHeader("Content-Type", "application/json").
		Body(body).
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	return s, nil
}

func (tc *TaskClient) Remove(ctx context.Context, opts *rv1.TaskRemoveOptions) error {

	req := tc.client.Delete(fmt.Sprintf("/namespace/%s/job/%s/task/%s", tc.namespace.String(), tc.job.Name(), tc.selflink.Name())).
		AddHeader("Content-Entity", "application/json")

	if opts != nil {
		if opts.Force {
			req.Param("force", strconv.FormatBool(opts.Force))
		}
	}

	var e *errors.Http

	if err := req.JSON(nil, &e); err != nil {
		return err
	}
	if e != nil {
		return errors.New(e.Message)
	}

	return nil
}

func newTaskClient(client *request.RESTClient, namespace, job, name string) *TaskClient {
	return &TaskClient{
		client:    client,
		namespace: *t.NewNamespaceSelfLink(namespace),
		job:       *t.NewJobSelfLink(namespace, job),
		selflink:  *t.NewTaskSelfLink(namespace, job, name),
	}
}
