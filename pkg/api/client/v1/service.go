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
	"github.com/lastbackend/lastbackend/pkg/api/client/http"
	"github.com/lastbackend/lastbackend/pkg/api/client/interfaces"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"io"
	"strconv"
	"github.com/lastbackend/lastbackend/pkg/api/client/watcher"
)

type ServiceClient struct {
	interfaces.Service
	client    http.Interface
	namespace string
	name      string
}

func (sc *ServiceClient) Deployment(name string) *DeploymentClient {
	return newDeploymentClient(sc.client, sc.namespace, sc.name, name)
}

func (sc *ServiceClient) Trigger(name string) *TriggerClient {
	return newTriggerClient(sc.client, sc.namespace, sc.name, name)
}

func (sc *ServiceClient) Create(ctx context.Context, opts *rv1.ServiceCreateOptions) (*vv1.Service, error) {

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

func (sc *ServiceClient) Update(ctx context.Context, opts *rv1.ServiceUpdateOptions) (*vv1.Service, error) {

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

func (sc *ServiceClient) Watch(ctx context.Context) (watcher.IWatcher, error) {
	return sc.client.Get(fmt.Sprintf("/namespace/%s/service/%s/watch", sc.namespace, sc.name)).Watch()
}

func newServiceClient(client http.Interface, namespace, name string) *ServiceClient {
	return &ServiceClient{client: client, namespace: namespace, name: name}
}
