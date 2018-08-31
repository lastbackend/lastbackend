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
	"strconv"

	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/util/http/request"
)

type NodeClient struct {
	client *request.RESTClient

	hostname string
}

func (nc NodeClient) List(ctx context.Context) (*vv1.NodeList, error) {

	var s *vv1.NodeList
	var e *errors.Http

	err := nc.client.Get(fmt.Sprintf("/cluster/node")).
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

func (nc NodeClient) Connect(ctx context.Context, opts *rv1.NodeConnectOptions) error {

	body := opts.ToJson()

	var e *errors.Http

	err := nc.client.Put(fmt.Sprintf("/cluster/node/%s", nc.hostname)).
		AddHeader("Content-Type", "application/json").
		Body([]byte(body)).
		JSON(nil, &e)

	if err != nil {
		return err
	}
	if e != nil {
		return errors.New(e.Message)
	}

	return nil
}

func (nc NodeClient) Get(ctx context.Context) (*vv1.Node, error) {

	var s *vv1.Node
	var e *errors.Http

	err := nc.client.Get(fmt.Sprintf("/cluster/node/%s", nc.hostname)).
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

func (nc NodeClient) GetSpec(ctx context.Context) (*vv1.NodeManifest, error) {

	var s *vv1.NodeManifest
	var e *errors.Http

	err := nc.client.Get(fmt.Sprintf("/cluster/node/%s/spec", nc.hostname)).
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

func (nc NodeClient) SetMeta(ctx context.Context, opts *rv1.NodeMetaOptions) (*vv1.Node, error) {

	body := opts.ToJson()

	var s *vv1.Node
	var e *errors.Http

	err := nc.client.Put(fmt.Sprintf("/cluster/node/%s/Meta", nc.hostname)).
		AddHeader("Content-Type", "application/json").
		Body([]byte(body)).
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	return s, nil
}

func (nc NodeClient) SetStatus(ctx context.Context, opts *rv1.NodeStatusOptions) error {

	body := opts.ToJson()

	var e *errors.Http

	err := nc.client.Put(fmt.Sprintf("/cluster/node/%s/status", nc.hostname)).
		AddHeader("Content-Type", "application/json").
		Body([]byte(body)).
		JSON(nil, &e)

	if err != nil {
		return err
	}
	if e != nil {
		return errors.New(e.Message)
	}

	return nil
}

func (nc NodeClient) SetPodStatus(ctx context.Context, pod string, opts *rv1.NodePodStatusOptions) error {

	body := opts.ToJson()

	var e *errors.Http

	err := nc.client.Put(fmt.Sprintf("/cluster/node/%s/status/pod/%s", nc.hostname, pod)).
		AddHeader("Content-Type", "application/json").
		Body([]byte(body)).
		JSON(nil, &e)

	if err != nil {
		return err
	}
	if e != nil {
		return errors.New(e.Message)
	}

	return nil
}

func (nc NodeClient) SetVolumeStatus(ctx context.Context, volume string, opts *rv1.NodeVolumeStatusOptions) error {

	body := opts.ToJson()

	var e *errors.Http

	err := nc.client.Put(fmt.Sprintf("/cluster/node/%s/status/volume/%s", nc.hostname, volume)).
		AddHeader("Content-Type", "application/json").
		Body([]byte(body)).
		JSON(nil, &e)

	if err != nil {
		return err
	}
	if e != nil {
		return errors.New(e.Message)
	}

	return nil
}

func (nc NodeClient) SetRouteStatus(ctx context.Context, route string, opts *rv1.NodeRouteStatusOptions) error {

	body := opts.ToJson()

	var e *errors.Http

	err := nc.client.Put(fmt.Sprintf("/cluster/node/%s/status/route/%s", nc.hostname, route)).
		AddHeader("Content-Type", "application/json").
		Body([]byte(body)).
		JSON(nil, &e)

	if err != nil {
		return err
	}
	if e != nil {
		return errors.New(e.Message)
	}

	return nil
}

func (nc NodeClient) Remove(ctx context.Context, opts *rv1.NodeRemoveOptions) error {

	req := nc.client.Delete(fmt.Sprintf("/cluster/node/%s", nc.hostname)).
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

func newNodeClient(req *request.RESTClient, hostname string) *NodeClient {
	return &NodeClient{client: req, hostname: hostname}
}
