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

package http

import (
	"context"

	"github.com/lastbackend/lastbackend/pkg/api/client/interfaces"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type RouteClient struct {
	interfaces.Route
}

func (s *RouteClient) Create(ctx context.Context, namespace string, opts rv1.RouteCreateOptions) (*vv1.Route, error) {
	return nil, nil
}

func (s *RouteClient) List(ctx context.Context, namespace string) (*vv1.RouteList, error) {
	return nil, nil
}

func (s *RouteClient) Get(ctx context.Context, namespace, name string) (*vv1.Route, error) {
	return nil, nil
}

func (s *RouteClient) Update(ctx context.Context, namespace, name string, opts rv1.RouteUpdateOptions) (*vv1.Route, error) {
	return nil, nil
}

func (s *RouteClient) Remove(ctx context.Context, namespace, name string, opts rv1.RouteRemoveOptions) error {
	return nil
}

func newRouteClient() *RouteClient {
	s := new(RouteClient)
	return s
}