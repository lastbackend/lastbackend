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

type IngressClient struct {
	interfaces.Ingress
	client   http.Interface
	hostname string
}

func (s *IngressClient) List(ctx context.Context) (*vv1.IngressList, error) {

	res := s.client.Get(fmt.Sprintf("/cluster/ingress")).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := res.Raw()
	if err != nil {
		return nil, err
	}

	var nl *vv1.IngressList

	if err := json.Unmarshal(buf, &nl); err != nil {
		return nil, err
	}

	return nl, nil
}

func (s *IngressClient) Get(ctx context.Context) (*vv1.Ingress, error) {

	res := s.client.Get(fmt.Sprintf("/cluster/ingress/%s", s.hostname)).
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

	var ns *vv1.Ingress

	if err := json.Unmarshal(buf, ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *IngressClient) GetSpec(ctx context.Context) (*vv1.IngressSpec, error) {

	res := s.client.Get(fmt.Sprintf("/cluster/ingress/%s/spec", s.hostname)).
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

	var ns *vv1.IngressSpec

	if err := json.Unmarshal(buf, &ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *IngressClient) SetMeta(ctx context.Context, opts *rv1.IngressMetaOptions) (*vv1.Ingress, error) {

	body := opts.ToJson()
	res := s.client.Put(fmt.Sprintf("/cluster/ingress/%s/Meta", s.hostname)).
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

	var ns *vv1.Ingress

	if err := json.Unmarshal(buf, &ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *IngressClient) Connect(ctx context.Context, opts *rv1.IngressConnectOptions) error {

	body := opts.ToJson()
	res := s.client.Put(fmt.Sprintf("/cluster/ingress/%s", s.hostname)).
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

func (s *IngressClient) SetStatus(ctx context.Context, opts *rv1.IngressStatusOptions) error {

	body := opts.ToJson()
	res := s.client.Put(fmt.Sprintf("/cluster/ingress/%s/status", s.hostname)).
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

func (s *IngressClient) SetRouteStatus(ctx context.Context, route string, opts *rv1.IngressRouteStatusOptions) error {

	body := opts.ToJson()
	res := s.client.Put(fmt.Sprintf("/cluster/ingress/%s/status/route/%s", s.hostname, route)).
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

func (s *IngressClient) Remove(ctx context.Context, opts *rv1.IngressRemoveOptions) error {

	req := s.client.Delete(fmt.Sprintf("/cluster/ingress/%s", s.hostname)).
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

func (s *IngressClient) Logs(ctx context.Context, pod, container string, opts *rv1.IngressLogsOptions) (io.ReadCloser, error) {
	req := s.client.Get(fmt.Sprintf("/pod/%s/%s/logs", pod, container))
	if opts != nil {
		if opts.Follow {
			req.Param("force", strconv.FormatBool(opts.Follow))
		}
	}

	return req.Stream()
}

func newIngressClient(req http.Interface, hostname string) *IngressClient {
	return &IngressClient{client: req, hostname: hostname}
}
