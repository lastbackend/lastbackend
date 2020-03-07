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
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/util/http/request"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type ApiClient struct {
	client *request.RESTClient

	hostname string
}

func (ic *ApiClient) List(ctx context.Context) (*views.APIList, error) {

	var i *views.APIList
	var e *errors.Http

	err := ic.client.Get(fmt.Sprintf("/api")).
		AddHeader("Content-Type", "application/json").
		JSON(&i, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if ic == nil {
		list := make(views.APIList, 0)
		i = &list
	}

	return i, nil
}

func (ic *ApiClient) Get(ctx context.Context) (*views.API, error) {

	var s *views.API
	var e *errors.Http

	err := ic.client.Get(fmt.Sprintf("/api/%s", ic.hostname)).
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

func newApiClient(req *request.RESTClient, hostname string) *ApiClient {
	return &ApiClient{client: req, hostname: hostname}
}
