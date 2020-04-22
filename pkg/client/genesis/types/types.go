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

package types

import (
	"context"

	"github.com/lastbackend/lastbackend/pkg/client/genesis/http/v1/request"
	"github.com/lastbackend/lastbackend/pkg/client/genesis/http/v1/views"
)

type ClientV1 interface {
	Account() AccountClientV1
	Cluster() ClusterClientV1
	Registry() RegistryClientV1
}

type AccountClientV1 interface {
	Get(ctx context.Context) error
	Login(ctx context.Context, opts *request.AccountLoginOptions) (*views.Session, error)
}

type ClusterClientV1 interface {
	Get(ctx context.Context, name string) (*views.ClusterView, error)
	List(ctx context.Context) (*views.ClusterList, error)
}

type RegistryClientV1 interface {
	Get(ctx context.Context) error
}
