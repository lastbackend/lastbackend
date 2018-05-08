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
)

type DeploymentClient struct {
	interfaces.Deployment
	client    http.Interface
	namespace string
	service   string
	name      string
}

func (dc *DeploymentClient) Pod(name string) *PodClient {
	return newPodClient(dc.client, dc.namespace, dc.service, dc.name, name)
}

func (dc *DeploymentClient) List(ctx context.Context) (*vv1.DeploymentList, error) {

	var s *vv1.DeploymentList
	var e *errors.Http

	err := dc.client.Get(fmt.Sprintf("/namespace/%s/service/%s/deploymet", dc.namespace, dc.service)).
		AddHeader("Content-Entity", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(vv1.DeploymentList, 0)
		s = &list
	}

	return s, nil
}

func (dc *DeploymentClient) Get(ctx context.Context) (*vv1.Deployment, error) {

	var s *vv1.Deployment
	var e *errors.Http

	err := dc.client.Get(fmt.Sprintf("/namespace/%s/service/%s/deploymet/%s", dc.namespace, dc.service, dc.name)).
		AddHeader("Content-Entity", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	return s, nil
}

func (dc *DeploymentClient) Update(ctx context.Context, opts *rv1.DeploymentUpdateOptions) (*vv1.Deployment, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Deployment
	var e *errors.Http

	err = dc.client.Put(fmt.Sprintf("/namespace/%s/service/%s/deployment/%s", dc.namespace, dc.service, dc.name)).
		AddHeader("Content-Entity", "application/json").
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

func newDeploymentClient(client http.Interface, namespace, service, name string) *DeploymentClient {
	return &DeploymentClient{client: client, namespace: namespace, service: service, name: name}
}
