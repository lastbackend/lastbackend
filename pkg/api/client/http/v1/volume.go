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
// patents in process, and are protected by trade volume or copyright law.
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

type VolumeClient struct {
	client *request.RESTClient

	namespace string
	name      string
}

func (vc *VolumeClient) Create(ctx context.Context, opts *rv1.VolumeManifest) (*vv1.Volume, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Volume
	var e *errors.Http

	err = vc.client.Post(fmt.Sprintf("/namespace/%s/volume", vc.namespace)).
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

func (vc *VolumeClient) List(ctx context.Context) (*vv1.VolumeList, error) {

	var s *vv1.VolumeList
	var e *errors.Http

	err := vc.client.Get(fmt.Sprintf("/namespace/%s/volume", vc.namespace)).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(vv1.VolumeList, 0)
		s = &list
	}

	return s, nil
}

func (vc *VolumeClient) Get(ctx context.Context) (*vv1.Volume, error) {

	var s *vv1.Volume
	var e *errors.Http

	err := vc.client.Get(fmt.Sprintf("/namespace/%s/volume/%s", vc.namespace, vc.name)).
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

func (vc *VolumeClient) Update(ctx context.Context, opts *rv1.VolumeManifest) (*vv1.Volume, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Volume
	var e *errors.Http

	err = vc.client.Put(fmt.Sprintf("/namespace/%s/volume/%s", vc.namespace, vc.name)).
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

func (vc *VolumeClient) Remove(ctx context.Context, opts *rv1.VolumeRemoveOptions) error {

	req := vc.client.Delete(fmt.Sprintf("/namespace/%s/volume/%s", vc.namespace, vc.name)).
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

func newVolumeClient(client *request.RESTClient, namespace, name string) *VolumeClient {
	return &VolumeClient{client: client, namespace: namespace, name: name}
}
