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

package cri

import (
	"github.com/lastbackend/lastbackend/pkg/cache"
	"github.com/lastbackend/lastbackend/pkg/common/context"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"io"
	"time"
)

type CRI interface {
	ContainerCreate(ctx context.IContext, spec *types.ContainerSpec) (string, error)
	ContainerStart(ctx context.IContext, ID string) error
	ContainerRestart(ctx context.IContext, ID string, timeout *time.Duration) error
	ContainerStop(ctx context.IContext, ID string, timeout *time.Duration) error
	ContainerPause(ctx context.IContext, ID string) error
	ContainerResume(ctx context.IContext, ID string) error
	ContainerRemove(ctx context.IContext, ID string, clean bool, force bool) error
	ContainerInspect(ctx context.IContext, ID string) (*types.Container, error)
	ContainerLogs(ctx context.IContext, ID string, stdout, stderr, follow bool) (io.ReadCloser, error)

	PodList(ctx context.IContext) ([]*types.Pod, error)

	ImagePull(ctx context.IContext, spec *types.ImageSpec) (io.ReadCloser, error)
	ImageRemove(ctx context.IContext, image string) error

	ImagePush(ctx context.IContext)
	ImageBuild(ctx context.IContext)
	ImageList(ctx context.IContext)

	Subscribe(ctx context.IContext, stg *cache.PodCache) chan types.ContainerEvent
}
