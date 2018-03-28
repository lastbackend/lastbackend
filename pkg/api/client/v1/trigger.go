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

	"encoding/json"
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

func (s *TriggerClient) Create(ctx context.Context, opts *rv1.TriggerCreateOptions) (*vv1.Trigger, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	res := s.client.Post(fmt.Sprintf("/namespace/%s/service/%s/trigger", s.namespace, s.service)).
		AddHeader("Content-Type", "application/json").
		Body(body).
		Do()

	if err := res.Error(); err != nil {
		return nil, err
	}

	buf, err := res.Raw()
	if err != nil {
		return nil, err
	}

	if code := res.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	}

	var ts = new(vv1.Trigger)

	if err := json.Unmarshal(buf, &ts); err != nil {
		return nil, err
	}

	return ts, nil
}

func (s *TriggerClient) List(ctx context.Context) (*vv1.TriggerList, error) {

	res := s.client.Get(fmt.Sprintf("/namespace/%s/service/%s/trigger", s.namespace, s.service)).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := res.Raw()
	if err != nil {
		return nil, err
	}

	var tl *vv1.TriggerList

	if err := json.Unmarshal(buf, &tl); err != nil {
		return nil, err
	}

	return tl, nil
}

func (s *TriggerClient) Get(ctx context.Context) (*vv1.Trigger, error) {

	res := s.client.Get(fmt.Sprintf("/namespace/%s/service/%s/trigger/%s", s.namespace, s.service, s.name)).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := res.Raw()
	if err != nil {
		return nil, err
	}

	if code := res.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	}

	var ts *vv1.Trigger

	if err := json.Unmarshal(buf, &ts); err != nil {
		return nil, err
	}

	return ts, nil
}

func (s *TriggerClient) Update(ctx context.Context, opts *rv1.TriggerUpdateOptions) (*vv1.Trigger, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	res := s.client.Put(fmt.Sprintf("/namespace/%s/service/%s/trigger/%s", s.namespace, s.service, s.name)).
		AddHeader("Content-Type", "application/json").
		Body(body).
		Do()

	buf, err := res.Raw()
	if err != nil {
		return nil, err
	}

	if code := res.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	}

	var ts *vv1.Trigger

	if err := json.Unmarshal(buf, &ts); err != nil {
		return nil, err
	}

	return ts, nil
}

func (s *TriggerClient) Remove(ctx context.Context, opts *rv1.TriggerRemoveOptions) error {

	res := s.client.Delete(fmt.Sprintf("/namespace/%s/service/%s/trigger/%s", s.namespace, s.service, s.name)).
		AddHeader("Content-Type", "application/json")

	if opts != nil {
		if opts.Force {
			res.Param("force", strconv.FormatBool(opts.Force))
		}
	}

	req := res.Do()

	buf, err := req.Raw()
	if err != nil {
		return err
	}

	if code := req.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return err
		}
		return errors.New(e.Message)
	}

	return nil
}

func newTriggerClient(req http.Interface, namespace, service, name string) *TriggerClient {
	return &TriggerClient{client: req, namespace: namespace, service: service, name: name}
}
