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

package containerd

import (
	"context"
	"io"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

func (r *Runtime) Auth(ctx context.Context, secret *models.SecretAuthData) (string, error) {
	return models.EmptyString, nil
}

func (r *Runtime) Pull(ctx context.Context, spec *models.ImageManifest, out io.Writer) (*models.Image, error) {

	ctx = namespaces.WithNamespace(ctx, "example")

	image, err := r.client.Pull(ctx, spec.Name, containerd.WithPullUnpack)
	if err != nil {
		return nil, err
	}

	// TODO: set structure fields
	img := new(models.Image)
	img.Meta.Name = image.Name()
	img.Meta.Digest = image.Target().Digest.String()

	return img, nil
}

func (r *Runtime) Push(ctx context.Context, spec *models.ImageManifest, out io.Writer) (*models.Image, error) {
	//err := r.client.Push(ctx, "docker.io/library/redis:latest", image.Target())
	//if err != nil {
	//	return nil, err
	//}
	//
	//// TODO: set structure fields
	//img := new(models.Image)
	//img.Meta.Name = image.Name()
	//img.Meta.Digest = image.Target().Digest.String()

	return nil, nil
}

func (r *Runtime) Build(ctx context.Context, stream io.Reader, spec *models.SpecBuildImage, out io.Writer) (*models.Image, error) {
	return nil, nil
}

func (r *Runtime) Remove(ctx context.Context, ID string) error {
	return nil
}

func (r *Runtime) List(ctx context.Context) ([]*models.Image, error) {
	return nil, nil
}

func (r *Runtime) Inspect(ctx context.Context, id string) (*models.Image, error) {
	return nil, nil
}

func (r *Runtime) Close() error {
	return r.client.Close()
}
