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

package docker

import (
	docker "github.com/docker/docker/api/types"

	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"io"
	"strings"
	"time"
)

func (r *Runtime) ContainerList(ctx context.Context, all bool) ([]*types.Container, error) {
	var cl []*types.Container

	items, err := r.client.ContainerList(ctx, docker.ContainerListOptions{
		All: all,
	})
	if err != nil {
		return cl, err
	}

	for _, item := range items {

		c, err := r.ContainerInspect(ctx, item.ID)
		if err != nil {
			log.Errorf("Can-not inspect container", err.Error())
			continue
		}

		if c == nil {
			continue
		}

		cl = append(cl, c)
	}

	return cl, nil
}

func (r *Runtime) ContainerCreate(ctx context.Context, spec *types.SpecTemplateContainer) (string, error) {

	c, err := r.client.ContainerCreate(
		ctx,
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

func (r *Runtime) ContainerStart(ctx context.Context, ID string) error {
	return r.client.ContainerStart(ctx, ID, docker.ContainerStartOptions{})
}

func (r *Runtime) ContainerRestart(ctx context.Context, ID string, timeout *time.Duration) error {
	return r.client.ContainerRestart(ctx, ID, timeout)
}

func (r *Runtime) ContainerStop(ctx context.Context, ID string, timeout *time.Duration) error {
	return r.client.ContainerStop(ctx, ID, timeout)
}

func (r *Runtime) ContainerPause(ctx context.Context, ID string) error {
	return r.client.ContainerPause(ctx, ID)
}

func (r *Runtime) ContainerResume(ctx context.Context, ID string) error {
	return r.client.ContainerUnpause(ctx, ID)
}

func (r *Runtime) ContainerRemove(ctx context.Context, ID string, clean bool, force bool) error {
	return r.client.ContainerRemove(ctx, ID, docker.ContainerRemoveOptions{
		RemoveVolumes: clean,
		Force:         force,
	})
}

func (r *Runtime) ContainerLogs(ctx context.Context, ID string, stdout, stderr, follow bool) (io.ReadCloser, error) {
	return r.client.ContainerLogs(ctx, ID, docker.ContainerLogsOptions{
		ShowStdout: stdout,
		ShowStderr: stderr,
		Follow:     follow,
		Timestamps: true,
		Details:    true,
	})
}

func (r *Runtime) ContainerInspect(ctx context.Context, ID string) (*types.Container, error) {

	log.Debug("Docker: Container Inspect")

	info, err := r.client.ContainerInspect(ctx, ID)
	if err != nil {
		log.Errorf("Docker: Container Inspect error: %s", err)
		return nil, err
	}

	container := &types.Container{
		ID:       info.ID,
		Name:     info.Name,
		Image:    info.Config.Image,
		Status:   info.State.Status,
		ExitCode: info.State.ExitCode,
	}

	container.Network.Gateway = info.NetworkSettings.Gateway
	container.Network.IPAddress = info.NetworkSettings.IPAddress

	container.Network.Ports = make(map[string][]*types.Port, 0)
	for key, val := range info.NetworkSettings.Ports {
		item := string(key)
		container.Network.Ports[item] = make([]*types.Port, 0)
		for _, port := range val {
			container.Network.Ports[item] = append(container.Network.Ports[item], &types.Port{
				HostIP:   port.HostIP,
				HostPort: port.HostPort,
			})
		}
	}

	switch info.State.Status {
	case types.StateCreated:
		container.State = types.StateCreated
	case types.StateStarted:
		container.State = types.StateStarted
	case types.StateRunning:
		container.State = types.StateStarted
	case types.StateStopped:
		container.State = types.StateStopped
	case types.StateExited:
		container.State = types.StateStopped
	case types.StateError:
		container.State = types.StateError
	}

	container.Created, _ = time.Parse(time.RFC3339Nano, info.Created)
	container.Started, _ = time.Parse(time.RFC3339Nano, info.State.StartedAt)

	meta, ok := info.Config.Labels["LB"]
	if !ok {
		log.Debug("Docker: Container Meta not found")
		return nil, nil
	}

	if len(strings.Split(meta, ":")) < 3 {
		return container, nil
	}
	container.Pod = meta

	return container, nil
}

// ToContainerCopy - https://docs.docker.com/engine/api/v1.29/#operation/PutContainerArchive
func (r *Runtime) ToContainerCopy(ctx context.Context, ID, path string, content io.Reader) error {
	return r.client.CopyToContainer(ctx, ID, path, content, docker.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
}
