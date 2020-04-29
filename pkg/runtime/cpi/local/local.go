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
	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

type Proxy struct {
}

func New() (*Proxy, error) {
	return &Proxy{}, nil
}

func (p *Proxy) Info(ctx context.Context) (map[string]*models.EndpointState, error) {
	es := make(map[string]*models.EndpointState)
	return es, nil
}

func (p *Proxy) Create(ctx context.Context, endpoint *models.EndpointManifest) (*models.EndpointState, error) {
	return new(models.EndpointState), nil
}

func (p *Proxy) Destroy(ctx context.Context, endpoint *models.EndpointState) error {
	return nil
}

func (p *Proxy) Update(ctx context.Context, endpoint *models.EndpointState, spec *models.EndpointManifest) (*models.EndpointState, error) {
	return new(models.EndpointState), nil
}
