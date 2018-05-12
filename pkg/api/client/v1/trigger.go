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
	"strconv"
)

type TriggerClient struct {
	interfaces.Trigger
	client    http.Interface
	namespace string
	service   string
	name      string
}

func (tc *TriggerClient) Create(ctx context.Context, opts *rv1.TriggerCreateOptions) (*vv1.Trigger, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Trigger
	var e *errors.Http

	err = tc.client.Post(fmt.Sprintf("/namespace/%s/service/%s/trigger", tc.namespace, tc.service)).
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

func (tc *TriggerClient) List(ctx context.Context) (*vv1.TriggerList, error) {

	var s *vv1.TriggerList
	var e *errors.Http

	err := tc.client.Get(fmt.Sprintf("/namespace/%s/service/%s/trigger", tc.namespace, tc.service)).
		AddHeader("Content-Entity", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(vv1.TriggerList, 0)
		s = &list
	}

	return s, nil
}

func (tc *TriggerClient) Get(ctx context.Context) (*vv1.Trigger, error) {

	var s *vv1.Trigger
	var e *errors.Http

	err := tc.client.Get(fmt.Sprintf("/namespace/%s/service/%s/trigger/%s", tc.namespace, tc.service, tc.name)).
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

func (tc *TriggerClient) Update(ctx context.Context, opts *rv1.TriggerUpdateOptions) (*vv1.Trigger, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Trigger
	var e *errors.Http

	err = tc.client.Put(fmt.Sprintf("/namespace/%s/service/%s/trigger/%s", tc.namespace, tc.service, tc.name)).
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

func (tc *TriggerClient) Remove(ctx context.Context, opts *rv1.TriggerRemoveOptions) error {

	req := tc.client.Delete(fmt.Sprintf("/namespace/%s/service/%s/trigger/%s", tc.namespace, tc.service, tc.name)).
		AddHeader("Content-Entity", "application/json")

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

func newTriggerClient(req http.Interface, namespace, service, name string) *TriggerClient {
	return &TriggerClient{client: req, namespace: namespace, service: service, name: name}
}
