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

type NodeClient struct {
	interfaces.Node
	client   http.Interface
	hostname string
}

func (s *NodeClient) List(ctx context.Context) (*vv1.NodeList, error) {

	req := s.client.Get(fmt.Sprintf("/cluster/node")).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := req.Raw()
	if err != nil {
		return nil, err
	}

	var nl *vv1.NodeList

	if err := json.Unmarshal(buf, &nl); err != nil {
		return nil, err
	}

	return nl, nil
}

func (s *NodeClient) Get(ctx context.Context) (*vv1.Node, error) {

	req := s.client.Get("/cluster/node").
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

	var ns *vv1.Node

	if err := json.Unmarshal(buf, &ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *NodeClient) Update(ctx context.Context, opts *rv1.NodeUpdateOptions) (*vv1.Node, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	req := s.client.Put(fmt.Sprintf("/cluster/namespace/%s", s.hostname)).
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

	var ns *vv1.Node

	if err := json.Unmarshal(buf, &ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *NodeClient) Remove(ctx context.Context, opts *rv1.NodeRemoveOptions) error {

	req := s.client.Delete(fmt.Sprintf("/cluster/node/%s", s.hostname)).
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

func newNodeClient(req http.Interface, hostname string) *NodeClient {
	return &NodeClient{client: req, hostname: hostname}
}
