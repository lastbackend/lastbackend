package docker

import (
	"context"
	docker "github.com/docker/docker/api/types"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"time"
)

func (r *Runtime) ContainerCreate(spec *types.ContainerSpec) (string, error) {

	c, err := r.client.ContainerCreate(
		context.Background(),
		GetConfig(spec),
		GetHostConfig(spec),
		GetNetworkConfig(spec),
		"",
	)

	if err != nil {
		return "", err
	}

	return c.ID, err
}

func (r *Runtime) ContainerStart(ID string) error {
	return r.client.ContainerStart(context.Background(), ID, docker.ContainerStartOptions{})
}

func (r *Runtime) ContainerRestart(ID string, timeout *time.Duration) error {
	return r.client.ContainerRestart(context.Background(), ID, timeout)
}

func (r *Runtime) ContainerStop(ID string, timeout *time.Duration) error {
	return r.client.ContainerStop(context.Background(), ID, timeout)
}

func (r *Runtime) ContainerPause(ID string) error {
	return r.client.ContainerPause(context.Background(), ID)
}

func (r *Runtime) ContainerResume(ID string) error {
	return r.client.ContainerUnpause(context.Background(), ID)
}

func (r *Runtime) ContainerRemove(ID string, clean bool, force bool) error {
	return r.client.ContainerRemove(context.Background(), ID, docker.ContainerRemoveOptions{
		RemoveLinks:   clean,
		RemoveVolumes: clean,
		Force:         force,
	})
}

func (r *Runtime) ContainerInspect() {

}
