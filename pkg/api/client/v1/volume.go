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
// patents in process, and are protected by trade volume or copyright law.
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

type VolumeClient struct {
	interfaces.Volume
	client    http.Interface
	namespace string
	name      string
}

func (s *VolumeClient) Create(ctx context.Context, opts *rv1.VolumeCreateOptions) (*vv1.Volume, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	req := s.client.Post(fmt.Sprintf("/namespace/%s/volume", s.namespace)).
		AddHeader("Content-Type", "application/json").
		Body(body).
		Do()

	if err := req.Error(); err != nil {
		return nil, err
	}

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

	var vs = new(vv1.Volume)

	if err := json.Unmarshal(buf, &vs); err != nil {
		return nil, err
	}

	return vs, nil
}

func (s *VolumeClient) List(ctx context.Context) (*vv1.VolumeList, error) {

	req := s.client.Get(fmt.Sprintf("/namespace/%s/volume", s.namespace)).
		AddHeader("Content-Type", "application/json").
		Do()

	buf, err := req.Raw()
	if err != nil {
		return nil, err
	}

	var vl *vv1.VolumeList

	if err := json.Unmarshal(buf, &vl); err != nil {
		return nil, err
	}

	return vl, nil
}

func (s *VolumeClient) Get(ctx context.Context) (*vv1.Volume, error) {

	req := s.client.Get(fmt.Sprintf("/namespace/%s/volume/%s", s.namespace, s.name)).
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

	var vs *vv1.Volume

	if err := json.Unmarshal(buf, &vs); err != nil {
		return nil, err
	}

	return vs, nil
}

func (s *VolumeClient) Update(ctx context.Context, opts *rv1.VolumeUpdateOptions) (*vv1.Volume, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	req := s.client.Put(fmt.Sprintf("/namespace/%s/volume/%s", s.namespace, s.name)).
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

	var vs *vv1.Volume

	if err := json.Unmarshal(buf, &vs); err != nil {
		return nil, err
	}

	return vs, nil
}

func (s *VolumeClient) Remove(ctx context.Context, opts *rv1.VolumeRemoveOptions) error {

	req := s.client.Delete(fmt.Sprintf("/namespace/%s/volume/%s", s.namespace, s.name)).
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

func newVolumeClient(client http.Interface, namespace, name string) *VolumeClient {
	s := new(VolumeClient)
	s.client = client
	s.namespace = namespace
	s.name = name
	return s
}
