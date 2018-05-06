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

package cpi

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"context"
)

type CPI interface {
	Info(ctx context.Context) (map[string]*types.EndpointStatus, error)
	Create(ctx context.Context, endpoint *types.EndpointSpec) (*types.EndpointStatus, error)
	Destroy(ctx context.Context, endpoint *types.EndpointSpec) (*types.EndpointStatus, error)
	Replace(ctx context.Context, spec *types.EndpointSpec) (*types.EndpointStatus, error)
}
