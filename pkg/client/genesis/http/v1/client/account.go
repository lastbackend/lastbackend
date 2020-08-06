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

	rv1 "github.com/lastbackend/lastbackend/pkg/client/genesis/http/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/client/genesis/http/v1/views"
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/util/http/request"
)

type AccountClient struct {
	client *request.RESTClient
}

func (ac *AccountClient) Get(ctx context.Context) error {
	return nil
}

func (ac *AccountClient) Login(ctx context.Context, opts *rv1.AccountLoginOptions) (*vv1.Session, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Session
	var e *errors.Http

	err = ac.client.Post("/session").
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

func newAccountClient(req *request.RESTClient) *AccountClient {
	return &AccountClient{client: req}
}
