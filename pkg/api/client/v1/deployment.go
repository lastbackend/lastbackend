//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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
)

type DeploymentClient struct {
	interfaces.Deployment
	req       http.Interface
	namespace string
	service   string
}

func (s *DeploymentClient) Pod(deployment string) *PodClient {
	return newPodClient(s.req, s.namespace, s.service, deployment)
}

func (s *DeploymentClient) List(ctx context.Context) (*vv1.DeploymentList, error) {
	return nil, nil
}

func (s *DeploymentClient) Get(ctx context.Context, na string) (*vv1.Deployment, error) {
	return nil, nil
}

func (s *DeploymentClient) Update(ctx context.Context, namespace, service, deployment string, opts *rv1.DeploymentUpdateOptions) (*vv1.Deployment, error) {
	return nil, nil
}

func newDeploymentClient(req http.Interface, namespace, service string) *DeploymentClient {
	return &DeploymentClient{req: req, namespace: namespace, service: service}
}
