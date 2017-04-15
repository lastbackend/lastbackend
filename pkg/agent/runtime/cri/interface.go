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
	"io"
	"time"
)

type CRI interface {
	ContainerCreate(spec *types.ContainerSpec) (string, error)
	ContainerStart(ID string) error
	ContainerRestart(ID string, timeout *time.Duration) error
	ContainerStop(ID string, timeout *time.Duration) error
	ContainerPause(ID string) error
	ContainerResume(ID string) error
	ContainerRemove(ID string, clean bool, force bool) error
	ContainerInspect(ID string) (*types.Container, error)

	PodList() ([]*types.Pod, error)

	ImagePull(spec *types.ImageSpec) (io.ReadCloser, error)
	ImageRemove(image string) error

	ImagePush()
	ImageBuild()
	ImageList()

	Subscribe() chan types.ContainerEvent
}
