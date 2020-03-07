//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade config or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package v1

import (
	"context"
	"fmt"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"strconv"

	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/util/http/request"
)

type ConfigClient struct {
	client    *request.RESTClient
	namespace string
	name      string
}

func (sc *ConfigClient) Create(ctx context.Context, opts *rv1.ConfigManifest) (*vv1.Config, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Config
	var e *errors.Http

	err = sc.client.Post(fmt.Sprintf("/namespace/%s/config", sc.namespace)).
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

func (sc *ConfigClient) Get(ctx context.Context) (*vv1.Config, error) {

	var s *vv1.Config
	var e *errors.Http

	err := sc.client.Get(fmt.Sprintf("/namespace/%s/config/%s", sc.namespace, sc.name)).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		s = new(vv1.Config)
	}

	return s, nil
}

func (sc *ConfigClient) List(ctx context.Context) (*vv1.ConfigList, error) {

	var s *vv1.ConfigList
	var e *errors.Http

	err := sc.client.Get(fmt.Sprintf("/namespace/%s/config", sc.namespace)).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(vv1.ConfigList, 0)
		s = &list
	}

	return s, nil
}

func (sc *ConfigClient) Update(ctx context.Context, opts *rv1.ConfigManifest) (*vv1.Config, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Config
	var e *errors.Http

	err = sc.client.Put(fmt.Sprintf("/namespace/%s/config/%s", sc.namespace, sc.name)).
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

func (sc *ConfigClient) Remove(ctx context.Context, opts *rv1.ConfigRemoveOptions) error {

	req := sc.client.Delete(fmt.Sprintf("/namespace/%s/config/%s", sc.namespace, sc.name)).
		AddHeader("Content-Type", "application/json")

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

func newConfigClient(client *request.RESTClient, namespace, name string) *ConfigClient {
	return &ConfigClient{client: client, namespace: namespace, name: name}
}
