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
	"io"
	"strconv"
)

type NodeClient struct {
	interfaces.Node
	client   http.Interface
	hostname string
}

func (s *NodeClient) List(ctx context.Context) (*vv1.NodeList, error) {

	res := s.client.Get(fmt.Sprintf("/cluster/node")).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := res.Raw()
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

	res := s.client.Get(fmt.Sprintf("/cluster/node/%s", s.hostname)).
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

	var ns *vv1.Node

	if err := json.Unmarshal(buf, &ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *NodeClient) GetSpec(ctx context.Context) (*vv1.NodeSpec, error) {

	res := s.client.Get(fmt.Sprintf("/cluster/node/%s/spec", s.hostname)).
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

	var ns *vv1.NodeSpec

	if err := json.Unmarshal(buf, &ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *NodeClient) SetMeta(ctx context.Context, opts *rv1.NodeMetaOptions) (*vv1.Node, error) {

	body := opts.ToJson()
	res := s.client.Put(fmt.Sprintf("/cluster/node/%s/Meta", s.hostname)).
		AddHeader("Content-Type", "application/json").
		Body([]byte(body)).
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

	var ns *vv1.Node

	if err := json.Unmarshal(buf, &ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *NodeClient) Connect(ctx context.Context, opts *rv1.NodeConnectOptions) error {

	body := opts.ToJson()
	res := s.client.Put(fmt.Sprintf("/cluster/node/%s", s.hostname)).
		AddHeader("Content-Type", "application/json").
		Body([]byte(body)).
		Do()

	buf, err := res.Raw()
	if err != nil {
		return err
	}

	if code := res.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return err
		}
		return errors.New(e.Message)
	}

	return nil
}

func (s *NodeClient) SetStatus(ctx context.Context, opts *rv1.NodeStatusOptions) error {

	body := opts.ToJson()
	res := s.client.Put(fmt.Sprintf("/cluster/node/%s/status", s.hostname)).
		AddHeader("Content-Type", "application/json").
		Body([]byte(body)).
		Do()

	buf, err := res.Raw()
	if err != nil {
		return err
	}

	if code := res.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return err
		}
		return errors.New(e.Message)
	}

	return nil
}

func (s *NodeClient) SetPodStatus(ctx context.Context, pod string, opts *rv1.NodePodStatusOptions) error {

	body := opts.ToJson()
	res := s.client.Put(fmt.Sprintf("/cluster/node/%s/status/pod/%s", s.hostname, pod)).
		AddHeader("Content-Type", "application/json").
		Body([]byte(body)).
		Do()

	buf, err := res.Raw()
	if err != nil {
		return err
	}

	if code := res.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return err
		}
		return errors.New(e.Message)
	}

	return nil
}

func (s *NodeClient) SetVolumeStatus(ctx context.Context, volume string, opts *rv1.NodeVolumeStatusOptions) error {

	body := opts.ToJson()
	res := s.client.Put(fmt.Sprintf("/cluster/node/%s/status/volume/%s", s.hostname, volume)).
		AddHeader("Content-Type", "application/json").
		Body([]byte(body)).
		Do()

	buf, err := res.Raw()
	if err != nil {
		return err
	}

	if code := res.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return err
		}
		return errors.New(e.Message)
	}

	return nil
}

func (s *NodeClient) SetRouteStatus(ctx context.Context, route string, opts *rv1.NodeRouteStatusOptions) error {

	body := opts.ToJson()
	res := s.client.Put(fmt.Sprintf("/cluster/node/%s/status/route/%s", s.hostname, route)).
		AddHeader("Content-Type", "application/json").
		Body([]byte(body)).
		Do()

	buf, err := res.Raw()
	if err != nil {
		return err
	}

	if code := res.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return err
		}
		return errors.New(e.Message)
	}

	return nil
}

func (s *NodeClient) Remove(ctx context.Context, opts *rv1.NodeRemoveOptions) error {

	req := s.client.Delete(fmt.Sprintf("/cluster/node/%s", s.hostname)).
		AddHeader("Content-Type", "application/json")

	if opts != nil {
		if opts.Force {
			req.Param("force", strconv.FormatBool(opts.Force))
		}
	}

	res := req.Do()

	buf, err := res.Raw()
	if err != nil {
		return err
	}

	if code := res.StatusCode(); 200 > code || code > 299 {
		var e *errors.Http
		if err := json.Unmarshal(buf, &e); err != nil {
			return err
		}
		return errors.New(e.Message)
	}

	return nil
}

func (s *NodeClient) Logs(ctx context.Context, pod, container string, opts *rv1.NodeLogsOptions) (io.ReadCloser, error) {

	req := s.client.Get(fmt.Sprintf("/pod/%s/%s/logs", pod, container)).
		AddHeader("Content-Type", "application/json")

	if opts != nil {
		if opts.Follow {
			req.Param("force", strconv.FormatBool(opts.Follow))
		}
	}

	return req.Stream()
}

func newNodeClient(req http.Interface, hostname string) *NodeClient {
	return &NodeClient{client: req, hostname: hostname}
}
