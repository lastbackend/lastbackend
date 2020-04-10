//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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
	"io"
	"strconv"
	"strings"
	"time"

	docker "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/tools/log"
)

func (r *Runtime) List(ctx context.Context, all bool) ([]*models.Container, error) {
	var cl = make([]*models.Container, 0)

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

func (r *Runtime) Create(ctx context.Context, manifest *models.ContainerManifest) (string, error) {

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

func (r *Runtime) Inspect(ctx context.Context, ID string) (*models.Container, error) {

	log.V(logLevel).Debug("Docker: Container Inspect")

	info, err := r.client.ContainerInspect(ctx, ID)
	if err != nil {
		log.Errorf("Docker: Container Inspect error: %s", err)
		return nil, err
	}

	c := &models.Container{
		ID:       info.ID,
		Name:     strings.Replace(info.Name, "/", "", 1),
		Image:    info.Config.Image,
		Status:   info.State.Status,
		Error:    info.State.Error,
		ExitCode: info.State.ExitCode,
		Labels:   info.Config.Labels,
		Envs:     info.Config.Env,
		Binds:    info.HostConfig.Binds,
	}

	c.Exec.Command = info.Config.Cmd
	c.Exec.Entrypoint = info.Config.Entrypoint
	c.Exec.Workdir = info.Config.WorkingDir

	c.Restart.Policy = info.HostConfig.RestartPolicy.Name
	c.Restart.Retry = info.HostConfig.RestartPolicy.MaximumRetryCount

	c.Network.Gateway = info.NetworkSettings.Gateway
	c.Network.IPAddress = info.NetworkSettings.IPAddress

	c.Network.Ports = make([]*models.SpecTemplateContainerPort, 0)
	for key, val := range info.HostConfig.PortBindings {

		for _, port := range val {

			p := &models.SpecTemplateContainerPort{
				HostIP:        port.HostIP,
				ContainerPort: uint16(key.Int()),
				Protocol:      key.Proto(),
			}

			var base = 10
			var size = 16
			pt, err := strconv.ParseUint(port.HostPort, base, size)
			if err != nil {
				continue
			}

			p.HostPort = uint16(pt)
			c.Network.Ports = append(c.Network.Ports, p)
		}
	}

	switch info.State.Status {
	case models.StateCreated:
		c.State = models.StateCreated
	case models.StateStarted:
		c.State = models.StateStarted
	case models.StatusRunning:
		c.State = models.StateStarted
	case models.StatusStopped:
		c.State = models.StatusStopped
	case models.StateExited:
		c.State = models.StatusStopped
	case models.StateError:
		c.State = models.StateError
	}

	c.Created, _ = time.Parse(time.RFC3339Nano, info.Created)
	c.Started, _ = time.Parse(time.RFC3339Nano, info.State.StartedAt)

	meta, ok := info.Config.Labels[models.ContainerTypeLBC]
	if ok {
		c.Pod = meta
	}
	c.Labels = info.Config.Labels

	return c, nil
}

func (r *Runtime) Wait(ctx context.Context, ID string) error {
	ok, err := r.client.ContainerWait(ctx, ID, container.WaitConditionNotRunning)
	select {
	case <-ok:
		return nil
	case e := <-err:
		return e
	}
}

// Copy - https://docs.docker.com/engine/api/v1.29/#operation/PutContainerArchive
func (r *Runtime) Copy(ctx context.Context, ID, path string, content io.Reader) error {
	return r.client.CopyToContainer(ctx, ID, path, content, docker.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
}
