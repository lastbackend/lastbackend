//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

	"github.com/lastbackend/lastbackend/pkg/api/client/interfaces"
	rv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
	"github.com/lastbackend/lastbackend/pkg/api/client/http"
)

type VolumeClient struct {
	interfaces.Volume
	client    http.Interface
	namespace string
	name      string
}

func (s *VolumeClient) Create(ctx context.Context, opts rv1.VolumeCreateOptions) (*vv1.Volume, error) {
	return nil, nil
}

func (s *VolumeClient) List(ctx context.Context) (*vv1.VolumeList, error) {
	return nil, nil
}

func (s *VolumeClient) Get(ctx context.Context) (*vv1.Volume, error) {
	return nil, nil
}

func (s *VolumeClient) Update(ctx context.Context, opts rv1.VolumeUpdateOptions) (*vv1.Volume, error) {
	return nil, nil
}

func (s *VolumeClient) Remove(ctx context.Context, opts rv1.VolumeRemoveOptions) error {
	return nil
}

func newVolumeClient(client http.Interface, namespace, name string) *VolumeClient {
	s := new(VolumeClient)
	s.client = client
	s.namespace = namespace
	s.name = name
	return s
}
