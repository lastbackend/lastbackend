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
	"context"
	docker "github.com/docker/docker/api/types"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"io"
	"strconv"
	"strings"
	"time"
)

func (r *Runtime) List(ctx context.Context, all bool) ([]*types.Container, error) {
	var cl = make([]*types.Container, 0)

	items, err := r.client.ContainerList(ctx, docker.ContainerListOptions{
		All: all,
	})
	if err != nil {
		return cl, err
	}

	for _, item := range items {

		c, err := r.Inspect(ctx, item.ID)
		if err != nil {
			log.Errorf("Can-not inspect container err: %v", err)
			continue
		}

		if c == nil {
			continue
		}

		cl = append(cl, c)
	}

	return cl, nil
}

func (r *Runtime) Create(ctx context.Context, manifest *types.ContainerManifest) (string, error) {

	c, err := r.client.ContainerCreate(
		ctx,
		GetConfig(manifest),
		GetHostConfig(manifest),
		GetNetworkConfig(manifest),
		manifest.Name,
	)
	if err != nil {
		return "", err
	}

	return c.ID, err
}

func (r *Runtime) Start(ctx context.Context, ID string) error {
	return r.client.ContainerStart(ctx, ID, docker.ContainerStartOptions{})
}

func (r *Runtime) Restart(ctx context.Context, ID string, timeout *time.Duration) error {
	return r.client.ContainerRestart(ctx, ID, timeout)
}

func (r *Runtime) Stop(ctx context.Context, ID string, timeout *time.Duration) error {
	return r.client.ContainerStop(ctx, ID, timeout)
}

func (r *Runtime) Pause(ctx context.Context, ID string) error {
	return r.client.ContainerPause(ctx, ID)
}

func (r *Runtime) Resume(ctx context.Context, ID string) error {
	return r.client.ContainerUnpause(ctx, ID)
}

func (r *Runtime) Remove(ctx context.Context, ID string, clean bool, force bool) error {
	return r.client.ContainerRemove(ctx, ID, docker.ContainerRemoveOptions{
		RemoveVolumes: clean,
		Force:         force,
	})
}

func (r *Runtime) Logs(ctx context.Context, ID string, stdout, stderr, follow bool) (io.ReadCloser, error) {
	return r.client.ContainerLogs(ctx, ID, docker.ContainerLogsOptions{
		ShowStdout: stdout,
		ShowStderr: stderr,
		Follow:     follow,
		Timestamps: true,
		Details:    true,
	})
}

func (r *Runtime) Inspect(ctx context.Context, ID string) (*types.Container, error) {

	log.V(logLevel).Debug("Docker: Container Inspect")

	info, err := r.client.ContainerInspect(ctx, ID)
	if err != nil {
		log.Errorf("Docker: Container Inspect error: %s", err)
		return nil, err
	}

	container := &types.Container{
		ID:       info.ID,
		Name:     strings.Replace(info.Name,"/","",1),
		Image:    info.Config.Image,
		Status:   info.State.Status,
		ExitCode: info.State.ExitCode,
		Labels:   info.Config.Labels,
		Envs: info.Config.Env,
		Binds: info.HostConfig.Binds,
	}

	container.Exec.Command = info.Config.Cmd
	container.Exec.Entrypoint = info.Config.Entrypoint
	container.Exec.Workdir = info.Config.WorkingDir

	container.Restart.Policy = info.HostConfig.RestartPolicy.Name
	container.Restart.Retry = info.HostConfig.RestartPolicy.MaximumRetryCount

	container.Network.Gateway = info.NetworkSettings.Gateway
	container.Network.IPAddress = info.NetworkSettings.IPAddress

	container.Network.Ports = make([]*types.SpecTemplateContainerPort, 0)
	for key, val := range info.HostConfig.PortBindings {

		for _, port := range val {

			p := &types.SpecTemplateContainerPort{
				HostIP: port.HostIP,
				ContainerPort: uint16(key.Int()),
				Protocol: key.Proto(),
			}

			var base = 10
			var size = 16
			pt, err := strconv.ParseUint(port.HostPort, base, size)
			if err != nil {
				continue
			}

			p.HostPort = uint16(pt)
			container.Network.Ports = append(container.Network.Ports, p)
		}
	}

	switch info.State.Status {
	case types.StateCreated:
		container.State = types.StateCreated
	case types.StateStarted:
		container.State = types.StateStarted
	case types.StatusRunning:
		container.State = types.StateStarted
	case types.StatusStopped:
		container.State = types.StatusStopped
	case types.StateExited:
		container.State = types.StatusStopped
	case types.StateError:
		container.State = types.StateError
	}

	container.Created, _ = time.Parse(time.RFC3339Nano, info.Created)
	container.Started, _ = time.Parse(time.RFC3339Nano, info.State.StartedAt)

	meta, ok := info.Config.Labels[types.ContainerTypeLBC]
	if ok {
		container.Pod = meta
	} else if _, ok := info.Config.Labels[types.ContainerTypeLBR]; !ok {
		log.V(logLevel).Debug("Docker: Container Meta not found")
	}

	container.Labels = info.Config.Labels

	return container, nil
}

// Copy - https://docs.docker.com/engine/api/v1.29/#operation/PutContainerArchive
func (r *Runtime) Copy(ctx context.Context, ID, path string, content io.Reader) error {
	return r.client.CopyToContainer(ctx, ID, path, content, docker.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
}
