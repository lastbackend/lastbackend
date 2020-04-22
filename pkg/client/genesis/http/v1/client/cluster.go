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

package client

import (
	"context"
	"fmt"

	"github.com/lastbackend/lastbackend/pkg/client/genesis/http/v1/views"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/util/http/request"
)

type ClusterClient struct {
	client *request.RESTClient
}

func (cc *ClusterClient) Get(ctx context.Context, name string) (*views.ClusterView, error) {

	var s *views.ClusterView
	var e *errors.Http

	err := cc.client.Get(fmt.Sprintf("/cluster/%s", name)).
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

func (cc *ClusterClient) List(ctx context.Context) (*views.ClusterList, error) {

	var s *views.ClusterList
	var e *errors.Http

	err := cc.client.Get("/cluster").
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(views.ClusterList, 0)
		s = &list
	}

	return s, nil
}

func newClusterClient(req *request.RESTClient) *ClusterClient {
	return &ClusterClient{client: req}
}
