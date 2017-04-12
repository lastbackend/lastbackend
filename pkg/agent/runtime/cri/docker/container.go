package docker

import (
	docker "github.com/docker/docker/api/types"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"strings"
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
		RemoveVolumes: clean,
		Force:         force,
	})
}

func (r *Runtime) ContainerInspect(ID string) (*types.Container, string, error) {
	log := context.Get().GetLogger()
	log.Debug("Docker: Container Inspect")

	var container *types.Container
	var pod string

	info, err := r.client.ContainerInspect(context.Background(), ID)
	if err != nil {
		log.Errorf("Docker: Container Inspect error: %s", err.Error())
		return container, pod, err
	}

	meta, ok := info.Config.Labels["LB_META"]
	if !ok {
		log.Debug("Docker: Container Meta not found")
		return container, pod, nil
	}

	pod = strings.Split(meta, "/")[0]
	container = &types.Container{
		ID:    info.ID,
		Image: info.Config.Image,
		State: info.State.Status,
	}

	container.Created, _ = time.Parse(time.RFC3339Nano, info.Created)
	container.Started, _ = time.Parse(time.RFC3339Nano, info.State.StartedAt)

	return container, pod, nil
}
