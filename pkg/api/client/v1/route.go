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
)

type RouteClient struct {
	interfaces.Route
	client    http.Interface
	namespace string
	name      string
}

func (s *RouteClient) Create(ctx context.Context, opts *rv1.RouteCreateOptions) (*vv1.Route, error) {
	return nil, nil
}

func (s *RouteClient) List(ctx context.Context) (*vv1.RouteList, error) {
	return nil, nil
}

func (s *RouteClient) Get(ctx context.Context) (*vv1.Route, error) {
	return nil, nil
}

func (s *RouteClient) Update(ctx context.Context, opts *rv1.RouteUpdateOptions) (*vv1.Route, error) {
	return nil, nil
}

func (s *RouteClient) Remove(ctx context.Context, opts *rv1.RouteRemoveOptions) error {
	return nil
}

func newRouteClient(client http.Interface, namespace, name string) *RouteClient {
	s := new(RouteClient)
	s.client = client
	s.namespace = namespace
	s.name = name
	return s
}
