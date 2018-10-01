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

	"github.com/lastbackend/lastbackend/pkg/api/client/types"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/util/http/request"
)

type NamespaceClient struct {
	client *request.RESTClient

	name string
}

func (nc *NamespaceClient) Secret(args ...string) types.SecretClientV1 {
	name := ""
	// Get any parameters passed to us out of the args variable into "real"
	// variables we created for them.
	for i := range args {
		switch i {
		case 0: // hostname
			name = args[0]
		default:
			panic("Wrong parameter count: (is allowed from 0 to 1)")
		}
	}
	return newSecretClient(nc.client, nc.name, name)
}

func (nc *NamespaceClient) Config(args ...string) types.ConfigClientV1 {
	name := ""
	// Get any parameters passed to us out of the args variable into "real"
	// variables we created for them.
	for i := range args {
		switch i {
		case 0: // hostname
			name = args[0]
		default:
			panic("Wrong parameter count: (is allowed from 0 to 1)")
		}
	}
	return newConfigClient(nc.client, nc.name, name)
}

func (nc *NamespaceClient) Service(args ...string) types.ServiceClientV1 {
	name := ""
	// Get any parameters passed to us out of the args variable into "real"
	// variables we created for them.
	for i := range args {
		switch i {
		case 0: // hostname
			name = args[0]
		default:
			panic("Wrong parameter count: (is allowed from 0 to 1)")
		}
	}
	return newServiceClient(nc.client, nc.name, name)
}

func (nc *NamespaceClient) Route(args ...string) types.RouteClientV1 {
	name := ""
	// Get any parameters passed to us out of the args variable into "real"
	// variables we created for them.
	for i := range args {
		switch i {
		case 0: // hostname
			name = args[0]
		default:
			panic("Wrong parameter count: (is allowed from 0 to 1)")
		}
	}
	return newRouteClient(nc.client, nc.name, name)
}

func (nc *NamespaceClient) Volume(args ...string) types.VolumeClientV1 {
	name := ""
	// Get any parameters passed to us out of the args variable into "real"
	// variables we created for them.
	for i := range args {
		switch i {
		case 0: // hostname
			name = args[0]
		default:
			panic("Wrong parameter count: (is allowed from 0 to 1)")
		}
	}
	return newVolumeClient(nc.client, nc.name, name)
}

func (nc *NamespaceClient) List(ctx context.Context) (*vv1.NamespaceList, error) {

	var s *vv1.NamespaceList
	var e *errors.Http

	err := nc.client.Get(fmt.Sprintf("/namespace")).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(vv1.NamespaceList, 0)
		s = &list
	}

	return s, nil
}

func (nc *NamespaceClient) Create(ctx context.Context, opts *rv1.NamespaceCreateOptions) (*vv1.Namespace, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Namespace
	var e *errors.Http

	err = nc.client.Post("/namespace").
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

func (nc *NamespaceClient) Get(ctx context.Context) (*vv1.Namespace, error) {

	var s *vv1.Namespace
	var e *errors.Http

	err := nc.client.Get(fmt.Sprintf("/namespace/%s", nc.name)).
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

func (nc *NamespaceClient) Update(ctx context.Context, opts *rv1.NamespaceUpdateOptions) (*vv1.Namespace, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Namespace
	var e *errors.Http

	err = nc.client.Put(fmt.Sprintf("/namespace/%s", nc.name)).
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

func (nc *NamespaceClient) Remove(ctx context.Context, opts *rv1.NamespaceRemoveOptions) error {

	req := nc.client.Delete(fmt.Sprintf("/namespace/%s", nc.name)).
		AddHeader("Content-Entity", "application/json")

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

func newNamespaceClient(client *request.RESTClient, name string) types.NamespaceClientV1 {
	return &NamespaceClient{client: client, name: name}
}
