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
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/util/http/request"
)

type ExporterClient struct {
	client *request.RESTClient

	hostname string
}

func (ic *ExporterClient) List(ctx context.Context) (*vv1.ExporterList, error) {

	var i *vv1.ExporterList
	var e *errors.Http

	err := ic.client.Get("/exporter").
		AddHeader("Content-Type", "application/json").
		JSON(&i, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if ic == nil {
		list := make(vv1.ExporterList, 0)
		i = &list
	}

	return i, nil
}

func (ic *ExporterClient) Get(ctx context.Context) (*vv1.Exporter, error) {

	var s *vv1.Exporter
	var e *errors.Http

	err := ic.client.Get(fmt.Sprintf("/exporter/%s", ic.hostname)).
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

func (ic *ExporterClient) Connect(ctx context.Context, opts *rv1.ExporterConnectOptions) error {

	body := opts.ToJson()

	var e *errors.Http

	err := ic.client.Put(fmt.Sprintf("/exporter/%s", ic.hostname)).
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

func (ic *ExporterClient) SetStatus(ctx context.Context, opts *rv1.ExporterStatusOptions) (*vv1.ExporterManifest, error) {

	body := opts.ToJson()

	var s *vv1.ExporterManifest
	var e *errors.Http

	err := ic.client.Put(fmt.Sprintf("/exporter/%s/status", ic.hostname)).
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

func newExporterClient(req *request.RESTClient, hostname string) *ExporterClient {
	return &ExporterClient{client: req, hostname: hostname}
}
