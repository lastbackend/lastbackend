//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/internal/pkg/errors"
	"github.com/lastbackend/lastbackend/internal/util/http/request"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type IngressClient struct {
	client *request.RESTClient

	hostname string
}

func (ic *IngressClient) List(ctx context.Context) (*views.IngressList, error) {

	var i *views.IngressList
	var e *errors.Http

	err := ic.client.Get(fmt.Sprintf("/ingress")).
		AddHeader("Content-Type", "application/json").
		JSON(&i, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if ic == nil {
		list := make(views.IngressList, 0)
		i = &list
	}

	return i, nil
}

func (ic *IngressClient) Get(ctx context.Context) (*views.Ingress, error) {

	var s *views.Ingress
	var e *errors.Http

	err := ic.client.Get(fmt.Sprintf("/ingress/%s", ic.hostname)).
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

func (ic *IngressClient) Connect(ctx context.Context, opts *rv1.IngressConnectOptions) error {

	body := opts.ToJson()

	var e *errors.Http

	err := ic.client.Put(fmt.Sprintf("/ingress/%s", ic.hostname)).
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

func (ic *IngressClient) SetStatus(ctx context.Context, opts *rv1.IngressStatusOptions) (*views.IngressManifest, error) {

	body := opts.ToJson()

	var s *views.IngressManifest
	var e *errors.Http

	err := ic.client.Put(fmt.Sprintf("/ingress/%s/status", ic.hostname)).
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

func newIngressClient(req *request.RESTClient, hostname string) *IngressClient {
	return &IngressClient{client: req, hostname: hostname}
}
