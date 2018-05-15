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
	"io"
	"strconv"
)

type IngressClient struct {
	interfaces.Ingress
	client   http.Interface
	hostname string
}

func (ic *IngressClient) List(ctx context.Context) (*vv1.IngressList, error) {

	var i *vv1.IngressList
	var e *errors.Http

	err := ic.client.Get(fmt.Sprintf("/cluster/ingress")).
		AddHeader("Content-Type", "application/json").
		JSON(&i, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if ic == nil {
		list := make(vv1.IngressList, 0)
		i = &list
	}

	return i, nil
}

func (ic *IngressClient) Get(ctx context.Context) (*vv1.Ingress, error) {

	var s *vv1.Ingress
	var e *errors.Http

	err := ic.client.Get(fmt.Sprintf("/cluster/ingress/%s", ic.hostname)).
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

func (ic *IngressClient) GetSpec(ctx context.Context) (*vv1.IngressSpec, error) {

	var s *vv1.IngressSpec
	var e *errors.Http

	err := ic.client.Get(fmt.Sprintf("/cluster/ingress/%s/spec", ic.hostname)).
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

func (ic *IngressClient) SetMeta(ctx context.Context, opts *rv1.IngressMetaOptions) (*vv1.Ingress, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Ingress
	var e *errors.Http

	err = ic.client.Get(fmt.Sprintf("/cluster/ingress/%s/Meta", ic.hostname)).
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

func (ic *IngressClient) Connect(ctx context.Context, opts *rv1.IngressConnectOptions) error {

	body, err := opts.ToJson()
	if err != nil {
		return err
	}

	var e *errors.Http

	err = ic.client.Put(fmt.Sprintf("/cluster/ingress/%s", ic.hostname)).
		AddHeader("Content-Type", "application/json").
		Body(body).
		JSON(nil, &e)

	if err != nil {
		return err
	}
	if e != nil {
		return errors.New(e.Message)
	}

	return nil
}

func (ic *IngressClient) SetStatus(ctx context.Context, opts *rv1.IngressStatusOptions) error {

	body, err := opts.ToJson()
	if err != nil {
		return err
	}

	var e *errors.Http

	err = ic.client.Put(fmt.Sprintf("/cluster/ingress/%s/status", ic.hostname)).
		AddHeader("Content-Type", "application/json").
		Body(body).
		JSON(nil, &e)

	if err != nil {
		return err
	}
	if e != nil {
		return errors.New(e.Message)
	}

	return nil
}

func (ic *IngressClient) SetRouteStatus(ctx context.Context, route string, opts *rv1.IngressRouteStatusOptions) error {

	body, err := opts.ToJson()
	if err != nil {
		return err
	}

	var e *errors.Http

	err = ic.client.Put(fmt.Sprintf("/cluster/ingress/%s/status/route/%s", ic.hostname, route)).
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

func (ic *IngressClient) Remove(ctx context.Context, opts *rv1.IngressRemoveOptions) error {

	req := ic.client.Delete(fmt.Sprintf("/cluster/ingress/%s", ic.hostname)).
		AddHeader("Content-Type", "application/json")

	if opts != nil {
		if opts.Force {
			req.Param("force", strconv.FormatBool(opts.Force))
		}
	}

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

func (ic *IngressClient) Logs(ctx context.Context, pod, container string, opts *rv1.IngressLogsOptions) (io.ReadCloser, error) {

	req := ic.client.Get(fmt.Sprintf("/pod/%s/%s/logs", pod, container))

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
