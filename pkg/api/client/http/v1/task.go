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

package v1

import (
	"context"
	"fmt"

	"github.com/lastbackend/lastbackend/pkg/api/client/types"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	t "github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/util/http/request"
)

type TaskClient struct {
	client *request.RESTClient

	namespace t.NamespaceSelfLink
	job       t.JobSelfLink
	selflink  t.TaskSelfLink
}

func (dc *TaskClient) Pod(args ...string) types.PodClientV1 {
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
	return newPodClient(dc.client, dc.namespace.String(), t.KindTask, dc.selflink.String(), name)
}

func (dc *TaskClient) List(ctx context.Context) (*vv1.TaskList, error) {

	var s *vv1.TaskList
	var e *errors.Http

	err := dc.client.Get(fmt.Sprintf("/namespace/%s/job/%s/task", dc.namespace.String(), dc.job.Name())).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(vv1.TaskList, 0)
		s = &list
	}

	return s, nil
}

func (dc *TaskClient) Get(ctx context.Context) (*vv1.Task, error) {

	var s *vv1.Task
	var e *errors.Http

	err := dc.client.Get(fmt.Sprintf("/namespace/%s/job/%s/task/%s", dc.namespace.String(), dc.job.Name(), dc.selflink.Name())).
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

func (dc *TaskClient) Cancel(ctx context.Context, opts *rv1.TaskCancelOptions) (*vv1.Task, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Task
	var e *errors.Http

	err = dc.client.Delete(fmt.Sprintf("/namespace/%s/job/%s/deployment/%s", dc.namespace.String(), dc.job.Name(), dc.selflink.Name())).
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

func newTaskClient(client *request.RESTClient, namespace, job, name string) *TaskClient {
	return &TaskClient{
		client:    client,
		namespace: *t.NewNamespaceSelfLink(namespace),
		job:       *t.NewJobSelfLink(namespace, job),
		selflink:  *t.NewTaskSelfLink(namespace, job, name),
	}
}
