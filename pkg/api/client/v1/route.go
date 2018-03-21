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
)

type RouteClient struct {
	interfaces.Route
	client    http.Interface
	namespace string
	name      string
}

func (s *RouteClient) Create(ctx context.Context, opts *rv1.RouteCreateOptions) (*vv1.Route, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	req := s.client.Post(fmt.Sprintf("/namespace/%s/route", s.namespace)).
		AddHeader("Content-Type", "application/json").
		Body(body).
		Do()

	if err := req.Error(); err != nil {
		return nil, err
	}

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

	var rs = new(vv1.Route)

	if err := json.Unmarshal(buf, &rs); err != nil {
		return nil, err
	}

	return rs, nil
}

func (s *RouteClient) List(ctx context.Context) (*vv1.RouteList, error) {

	req := s.client.Get(fmt.Sprintf("/namespace/%s/route", s.namespace)).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := req.Raw()
	if err != nil {
		return nil, err
	}

	var rl *vv1.RouteList

	if err := json.Unmarshal(buf, &rl); err != nil {
		return nil, err
	}

	return rl, nil
}

func (s *RouteClient) Get(ctx context.Context) (*vv1.Route, error) {

	req := s.client.Get(fmt.Sprintf("/namespace/%s/route/%s", s.namespace, s.name)).
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

	var rs *vv1.Route

	if err := json.Unmarshal(buf, &rs); err != nil {
		return nil, err
	}

	return rs, nil
}

func (s *RouteClient) Update(ctx context.Context, opts *rv1.RouteUpdateOptions) (*vv1.Route, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	req := s.client.Put(fmt.Sprintf("/namespace/%s/route/%s", s.namespace, s.name)).
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

	var rs *vv1.Route

	if err := json.Unmarshal(buf, &rs); err != nil {
		return nil, err
	}

	return rs, nil
}

func (s *RouteClient) Remove(ctx context.Context, opts *rv1.RouteRemoveOptions) error {

	req := s.client.Delete(fmt.Sprintf("/namespace/%s/route/%s", s.namespace, s.name)).
		AddHeader("Content-Type", "application/json").
		Do()

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

func newRouteClient(client http.Interface, namespace, name string) *RouteClient {
	s := new(RouteClient)
	s.client = client
	s.namespace = namespace
	s.name = name
	return s
}
