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
	"github.com/lastbackend/lastbackend/pkg/util/http/request"
	"io"
	"net/http"
	"strconv"
)

type JobClient struct {
	client *request.RESTClient

	namespace string
	name      string
}

func (sc *JobClient) Tasks(args ...string) types.TaskClientV1 {
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
	return newTaskClient(sc.client, sc.namespace, sc.name, name)
}

func (sc *JobClient) Create(ctx context.Context, opts *rv1.JobManifest) (*vv1.Job, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Job
	var e *errors.Http

	err = sc.client.Post(fmt.Sprintf("/namespace/%s/job", sc.namespace)).
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

func (sc *JobClient) Run(ctx context.Context, opts *rv1.TaskManifest) (*vv1.Task, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Task
	var e *errors.Http

	err = sc.client.Post(fmt.Sprintf("/namespace/%s/job/%s/task", sc.namespace, sc.name)).
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

func (sc *JobClient) List(ctx context.Context) (*vv1.JobList, error) {

	var s *vv1.JobList
	var e *errors.Http

	err := sc.client.Get(fmt.Sprintf("/namespace/%s/job", sc.namespace)).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(vv1.JobList, 0)
		s = &list
	}

	return s, nil
}

func (sc *JobClient) Get(ctx context.Context) (*vv1.Job, error) {

	var s *vv1.Job
	var e *errors.Http

	err := sc.client.Get(fmt.Sprintf("/namespace/%s/job/%s", sc.namespace, sc.name)).
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

func (sc *JobClient) Update(ctx context.Context, opts *rv1.JobManifest) (*vv1.Job, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Job
	var e *errors.Http

	err = sc.client.Put(fmt.Sprintf("/namespace/%s/job/%s", sc.namespace, sc.name)).
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

func (sc *JobClient) Remove(ctx context.Context, opts *rv1.JobRemoveOptions) error {

	req := sc.client.Delete(fmt.Sprintf("/namespace/%s/job/%s", sc.namespace, sc.name)).
		AddHeader("Content-Type", "application/json")

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

func (sc *JobClient) Logs(ctx context.Context, opts *rv1.JobLogsOptions) (io.ReadCloser, *http.Response, error) {

	res := sc.client.Get(fmt.Sprintf("/namespace/%s/job/%s/logs", sc.namespace, sc.name))

	if opts != nil {
		res.Param("task", opts.Task)
		res.Param("pod", opts.Pod)
		res.Param("container", opts.Container)

		res.Param("tail", fmt.Sprintf("%d", opts.Tail))

		if opts.Follow {
			res.Param("follow", strconv.FormatBool(opts.Follow))
		}
	}

	return res.Stream()
}

func newJobClient(client *request.RESTClient, namespace, name string) *JobClient {
	return &JobClient{client: client, namespace: namespace, name: name}
}
