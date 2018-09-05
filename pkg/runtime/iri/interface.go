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

package iri

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"io"
)

// IMI - Image Runtime Interface
type IRI interface {
	Auth(ctx context.Context, secret *types.SecretAuthData) (string, error)
	Pull(ctx context.Context, spec *types.ImageManifest, out io.Writer) (*types.Image, error)
	Remove(ctx context.Context, image string) error
	Push(ctx context.Context, spec *types.ImageManifest, out io.Writer) (*types.Image, error)
	Build(ctx context.Context, stream io.Reader, spec *types.SpecBuildImage, out io.Writer) (*types.Image, error)
	List(ctx context.Context) ([]*types.Image, error)
	Inspect(ctx context.Context, id string) (*types.Image, error)
	Subscribe(ctx context.Context) (chan *types.Image, error)
}
