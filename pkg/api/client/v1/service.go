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

type ServiceClient struct {
	interfaces.Service
	client       http.Interface
	namespace string
	name      string
}

func (s *ServiceClient) Deployment(name string) *DeploymentClient {
	return newDeploymentClient(s.client, s.namespace, name)
}

func (s *ServiceClient) Trigger(name string) *TriggerClient {
	return newTriggerClient(s.client, s.namespace, name)
}

func (s *ServiceClient) Create(ctx context.Context, opts *rv1.ServiceCreateOptions) (*vv1.ServiceList, error) {
	return nil, nil
}

func (s *ServiceClient) List(ctx context.Context) (*vv1.ServiceList, error) {
	return nil, nil
}

func (s *ServiceClient) Get(ctx context.Context) (*vv1.Service, error) {
	return nil, nil
}

func (s *ServiceClient) Update(ctx context.Context, opts *rv1.ServiceUpdateOptions) (*vv1.NamespaceList, error) {
	return nil, nil
}

func (s *ServiceClient) Remove(ctx context.Context, opts *rv1.ServiceRemoveOptions) error {
	return nil
}

func newServiceClient(client http.Interface, namespace, name string) *ServiceClient {
	return &ServiceClient{client: client, namespace: namespace, name: name}
}
