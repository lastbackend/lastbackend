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

package cri

import (
	docker "github.com/docker/docker/api/types"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/node/state"

	"context"
	"io"
	"time"
)

// CRI - Container Runtime Interface
type CRI interface {
	ContainerRuntime
	ImageRuntime
	Subscribe(ctx context.Context, state *state.PodState, p chan string)
}

type ContainerRuntime interface {
	ContainerList(ctx context.Context, all bool) ([]*types.Container, error)
	ContainerCreate(ctx context.Context, spec *types.SpecTemplateContainer) (string, error)
	ContainerStart(ctx context.Context, ID string) error
	ContainerRestart(ctx context.Context, ID string, timeout *time.Duration) error
	ContainerStop(ctx context.Context, ID string, timeout *time.Duration) error
	ContainerPause(ctx context.Context, ID string) error
	ContainerResume(ctx context.Context, ID string) error
	ContainerRemove(ctx context.Context, ID string, clean bool, force bool) error
	ContainerInspect(ctx context.Context, ID string) (*types.Container, error)
	ContainerLogs(ctx context.Context, ID string, stdout, stderr, follow bool) (io.ReadCloser, error)
	ToContainerCopy(ctx context.Context, ID, path string, content io.Reader) error
}

type ImageRuntime interface {
	ImagePull(ctx context.Context, spec *types.SpecTemplateContainerImage, secret *types.Secret) (io.ReadCloser, error)
	ImageRemove(ctx context.Context, image string) error
	ImagePush(ctx context.Context, spec *types.SpecTemplateContainerImage, secret *types.Secret) (io.ReadCloser, error)
	ImageBuild(ctx context.Context, stream io.Reader, spec *types.SpecBuildImage) (io.ReadCloser, error)
	ImageList(ctx context.Context) ([]docker.ImageSummary, error)
	ImageInspect(ctx context.Context, id string) (*types.ImageInfo, []byte, error)
}
