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

	"github.com/lastbackend/lastbackend/pkg/api/client/http"
	"github.com/lastbackend/lastbackend/pkg/api/client/interfaces"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"fmt"
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
)

type DeploymentClient struct {
	interfaces.Deployment
	client    http.Interface
	namespace string
	service   string
	name      string
}

func (s *DeploymentClient) Pod(deployment string) *PodClient {
	return newPodClient(s.client, s.namespace, s.service, deployment)
}

func (s *DeploymentClient) List(ctx context.Context) (*vv1.DeploymentList, error) {

	req := s.client.Get(fmt.Sprintf("/namespace/%s/service/%s/deploymet", s.namespace, s.service)).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := req.Raw()
	if err != nil {
		return nil, err
	}

	var dl *vv1.DeploymentList

	if err := json.Unmarshal(buf, &dl); err != nil {
		return nil, err
	}

	return dl, nil
}

func (s *DeploymentClient) Get(ctx context.Context) (*vv1.Deployment, error) {

	req := s.client.Get(fmt.Sprintf("/namespace/%s/service/%s/deploymet/%s", s.namespace, s.service, s.name)).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := req.Raw()
	if err != nil {
		return nil, err
	}

	if code := req.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	}

	var ds *vv1.Deployment

	if err := json.Unmarshal(buf, &ds); err != nil {
		return nil, err
	}

	return ds, nil
}

func (s *DeploymentClient) Update(ctx context.Context, opts *rv1.DeploymentUpdateOptions) (*vv1.Deployment, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	req := s.client.Put(fmt.Sprintf("/namespace/%s/service/%s/deployment/%s", s.namespace, s.service, s.name)).
		AddHeader("Content-Type", "application/json").
		Body(body).
		Do()

	buf, err := req.Raw()
	if err != nil {
		return nil, err
	}

	if code := req.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	}

	var ds *vv1.Deployment

	if err := json.Unmarshal(buf, &ds); err != nil {
		return nil, err
	}

	return ds, nil
}

func newDeploymentClient(client http.Interface, namespace, service, name string) *DeploymentClient {
	return &DeploymentClient{client: client, namespace: namespace, service: service, name: name}
}
