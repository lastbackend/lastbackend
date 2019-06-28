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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"io"
	"net/http"
	"strconv"

	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/util/http/request"
)

type PodClient struct {
	client *request.RESTClient

	parent struct {
		kind     string
		selflink string
	}

	namespace string
	name      string
}

func (pc *PodClient) List(ctx context.Context) (*vv1.PodList, error) {

	var s *vv1.PodList
	var e *errors.Http

	var url string

	switch pc.parent.kind {
	case types.KindDeployment:
		dsl := types.DeploymentSelfLink{}
		if err := dsl.Parse(pc.parent.selflink); err != nil {
			return nil, err
		}
		_, svc := dsl.Parent()
		url = fmt.Sprintf("/namespace/%s/service/%s/deploymet/%s/pod", pc.namespace, svc.Name(), dsl.Name())
	case types.KindTask:
		tsl := types.TaskSelfLink{}
		if err := tsl.Parse(pc.parent.selflink); err != nil {
			return nil, err
		}
		_, job := tsl.Parent()
		url = fmt.Sprintf("/namespace/%s/job/%s/task/%s/pod", pc.namespace, job.Name(), tsl.Name())
	}

	err := pc.client.Get(url).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(vv1.PodList, 0)
		s = &list
	}

	return s, nil
}

func (pc *PodClient) Get(ctx context.Context) (*vv1.Pod, error) {

	var s *vv1.Pod
	var e *errors.Http

	var url string

	switch pc.parent.kind {
	case types.KindDeployment:
		dsl := types.DeploymentSelfLink{}
		if err := dsl.Parse(pc.parent.selflink); err != nil {
			return nil, err
		}
		_, svc := dsl.Parent()
		url = fmt.Sprintf("/namespace/%s/service/%s/deploymet/%s/pod/%s", pc.namespace, svc.Name(), dsl.Name(), pc.name)
	case types.KindTask:
		tsl := types.TaskSelfLink{}
		if err := tsl.Parse(pc.parent.selflink); err != nil {
			return nil, err
		}
		_, job := tsl.Parent()
		url = fmt.Sprintf("/namespace/%s/job/%s/task/%s/pod/%s", pc.namespace, job.Name(), tsl.Name(), pc.name)
	}

	err := pc.client.Get(url).
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

func (pc *PodClient) Logs(ctx context.Context, opts *rv1.PodLogsOptions) (io.ReadCloser, *http.Response, error) {

	var url, parent string

	switch pc.parent.kind {
	case types.KindDeployment:
		dsl := types.DeploymentSelfLink{}
		if err := dsl.Parse(pc.parent.selflink); err != nil {
			return nil, nil, err
		}
		parent = dsl.Name()
		_, svc := dsl.Parent()
		url = fmt.Sprintf("/namespace/%s/service/%s/logs", pc.namespace, svc.Name())
	case types.KindTask:
		tsl := types.TaskSelfLink{}
		if err := tsl.Parse(pc.parent.selflink); err != nil {
			return nil, nil, err
		}
		parent = tsl.Name()
		_, job := tsl.Parent()
		url = fmt.Sprintf("/namespace/%s/job/%s/logs", pc.namespace, job.Name())
	}

	res := pc.client.Get(url)

	if opts != nil {

		switch pc.parent.kind {
		case types.KindDeployment:
			res.Param("deployment", parent)
		case types.KindTask:
			res.Param("task", parent)
		}

		res.Param("pod", pc.name)
		res.Param("container", opts.Container)

		if opts.Follow {
			res.Param("follow", strconv.FormatBool(opts.Follow))
		}
	}

	return res.Stream()
}

func newPodClient(client *request.RESTClient, namespace, kind, parent, name string) *PodClient {
	pc := PodClient{client: client, namespace: namespace, name: name}
	pc.parent.kind = kind
	pc.parent.selflink = parent
	return &pc
}
