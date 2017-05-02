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
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/context"
	"io"
	"time"
	"github.com/lastbackend/lastbackend/pkg/pod"
)

type CRI interface {
	ContainerCreate(ctx context.Context, spec *types.ContainerSpec) (string, error)
	ContainerStart(ctx context.Context, ID string) error
	ContainerRestart(ctx context.Context, ID string, timeout *time.Duration) error
	ContainerStop(ctx context.Context, ID string, timeout *time.Duration) error
	ContainerPause(ctx context.Context, ID string) error
	ContainerResume(ctx context.Context, ID string) error
	ContainerRemove(ctx context.Context, ID string, clean bool, force bool) error
	ContainerInspect(ctx context.Context, ID string) (*types.Container, error)

	PodList(ctx context.Context) ([]*types.Pod, error)

	ImagePull(ctx context.Context, spec *types.ImageSpec) (io.ReadCloser, error)
	ImageRemove(ctx context.Context, image string) error

	ImagePush(ctx context.Context)
	ImageBuild(ctx context.Context)
	ImageList(ctx context.Context)

	Subscribe(ctx context.Context, stg *pod.PodStorage) chan types.ContainerEvent
}
