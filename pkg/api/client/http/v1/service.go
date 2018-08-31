//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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
	"io"
	"strconv"

	"github.com/lastbackend/lastbackend/pkg/api/client/types"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/util/http/request"
)

type ServiceClient struct {
	client *request.RESTClient

	namespace string
	name      string
}

func (sc *ServiceClient) Deployment(args ...string) types.DeploymentClientV1 {
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
	return newDeploymentClient(sc.client, sc.namespace, sc.name, name)
}

func (sc *ServiceClient) Trigger(args ...string) types.TriggerClientV1 {
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
	return newTriggerClient(sc.client, sc.namespace, sc.name, name)
}

func (sc *ServiceClient) Create(ctx context.Context, opts *rv1.ServiceManifest) (*vv1.Service, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Service
	var e *errors.Http

	err = sc.client.Post(fmt.Sprintf("/namespace/%s/service", sc.namespace)).
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

func (sc *ServiceClient) List(ctx context.Context) (*vv1.ServiceList, error) {

	var s *vv1.ServiceList
	var e *errors.Http

	err := sc.client.Get(fmt.Sprintf("/namespace/%s/service", sc.namespace)).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(vv1.ServiceList, 0)
		s = &list
	}

	return s, nil
}

func (sc *ServiceClient) Get(ctx context.Context) (*vv1.Service, error) {

	var s *vv1.Service
	var e *errors.Http

	err := sc.client.Get(fmt.Sprintf("/namespace/%s/service/%s", sc.namespace, sc.name)).
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

func (sc *ServiceClient) Update(ctx context.Context, opts *rv1.ServiceManifest) (*vv1.Service, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Service
	var e *errors.Http

	err = sc.client.Put(fmt.Sprintf("/namespace/%s/service/%s", sc.namespace, sc.name)).
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

func (sc *ServiceClient) Remove(ctx context.Context, opts *rv1.ServiceRemoveOptions) error {

	req := sc.client.Delete(fmt.Sprintf("/namespace/%s/service/%s", sc.namespace, sc.name)).
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

func (sc *ServiceClient) Logs(ctx context.Context, opts *rv1.ServiceLogsOptions) (io.ReadCloser, error) {

	res := sc.client.Get(fmt.Sprintf("/namespace/%s/service/%s/logs", sc.namespace, sc.name))

	if opts != nil {
		res.Param("deployment", opts.Deployment)
		res.Param("pod", opts.Pod)
		res.Param("container", opts.Container)

		if opts.Follow {
			res.Param("follow", strconv.FormatBool(opts.Follow))
		}
	}

	return res.Stream()
}

func newServiceClient(client *request.RESTClient, namespace, name string) *ServiceClient {
	return &ServiceClient{client: client, namespace: namespace, name: name}
}
