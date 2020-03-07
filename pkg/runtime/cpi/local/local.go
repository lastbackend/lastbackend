//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package local

import (
	"context"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/pkg/runtime/cpi"
)

type Proxy struct {
	cpi.CPI
}

func (p *Proxy) Info(ctx context.Context) (map[string]*types.EndpointState, error) {
	es := make(map[string]*types.EndpointState)
	return es, nil
}

func (p *Proxy) Create(ctx context.Context, endpoint *types.EndpointManifest) (*types.EndpointState, error) {
	return new(types.EndpointState), nil
}

func (p *Proxy) Destroy(ctx context.Context, endpoint *types.EndpointState) error {
	return nil
}

func (p *Proxy) Update(ctx context.Context, endpoint *types.EndpointState, spec *types.EndpointManifest) (*types.EndpointState, error) {
	return new(types.EndpointState), nil
}

func New() (*Proxy, error) {
	return &Proxy{}, nil
}
