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
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package cri

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"io"
	"time"
)

// CRI - Container Runtime Interface
type CRI interface {
	List(ctx context.Context, all bool) ([]*types.Container, error)
	Create(ctx context.Context, spec *types.ContainerManifest) (string, error)
	Start(ctx context.Context, ID string) error
	Restart(ctx context.Context, ID string, timeout *time.Duration) error
	Stop(ctx context.Context, ID string, timeout *time.Duration) error
	Pause(ctx context.Context, ID string) error
	Resume(ctx context.Context, ID string) error
	Remove(ctx context.Context, ID string, clean bool, force bool) error
	Inspect(ctx context.Context, ID string) (*types.Container, error)
	Logs(ctx context.Context, ID string, stdout, stderr, follow bool) (io.ReadCloser, error)
	Copy(ctx context.Context, ID, path string, content io.Reader) error
	Wait(ctx context.Context, ID string) error
	Subscribe(ctx context.Context, container chan *types.Container) error
}
