package cri

import (
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/satori/go.uuid"
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

	PodList() (map[uuid.UUID]*types.Pod, error)

	ImagePull(spec *types.ImageSpec) (io.ReadCloser, error)
	ImageRemove(image string) error

	ImagePush()
	ImageBuild()
	ImageList()
}
