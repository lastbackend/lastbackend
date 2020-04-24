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
	"time"

	"github.com/lastbackend/lastbackend/internal/pkg/models"
)

func (r *Runtime) List(ctx context.Context, all bool) ([]*models.Container, error) {
	var cl = make([]*models.Container, 0)
	return cl, nil
}

func (r *Runtime) Create(ctx context.Context, manifest *models.ContainerManifest) (string, error) {
	return "", nil
}

func (r *Runtime) Start(ctx context.Context, ID string) error {
	return nil
}

func (r *Runtime) Restart(ctx context.Context, ID string, timeout *time.Duration) error {
	return nil
}

func (r *Runtime) Stop(ctx context.Context, ID string, timeout *time.Duration) error {
	return nil
}

func (r *Runtime) Pause(ctx context.Context, ID string) error {
	return nil
}

func (r *Runtime) Resume(ctx context.Context, ID string) error {
	return nil
}

func (r *Runtime) Remove(ctx context.Context, ID string, clean bool, force bool) error {
	return nil
}

func (r *Runtime) Logs(ctx context.Context, ID string, stdout, stderr, follow bool) (io.ReadCloser, error) {
	return nil, nil
}

func (r *Runtime) Inspect(ctx context.Context, ID string) (*models.Container, error) {
	return nil, nil
}

func (r *Runtime) Wait(ctx context.Context, ID string) error {
	return nil
}

func (r *Runtime) Copy(ctx context.Context, ID, path string, content io.Reader) error {
	return nil
}

func (r *Runtime) Close() error {
	return r.client.Close()
}
