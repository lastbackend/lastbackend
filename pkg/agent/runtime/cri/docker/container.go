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

func (r *Runtime) ContainerInspect(ID string) (*types.Container, error) {
	log := context.Get().GetLogger()
	log.Debug("Docker: Container Inspect")

	var container *types.Container
	var pod, spc string

	info, err := r.client.ContainerInspect(context.Background(), ID)
	if err != nil {
		log.Errorf("Docker: Container Inspect error: %s", err.Error())
		return container, err
	}

	meta, ok := info.Config.Labels["LB_META"]
	if !ok {
		log.Debug("Docker: Container Meta not found")
		return container, nil
	}

	match := strings.Split(meta, "/")

	if len(match) < 3 {
		return nil, nil
	}

	pod = match[0]
	spc = match[2]

	container = &types.Container{
		ID:    info.ID,
		Pod:   pod,
		Spec:  spc,
		Image: info.Config.Image,
		State: info.State.Status,
	}

	switch info.State.Status {
	case types.StateCreated: container.State = types.StateCreated
	case types.StateStarted: container.State = types.StateStarted
	case types.StateRunning: container.State = types.StateStarted
	case types.StateStopped: container.State = types.StateStopped
	case types.StateExited: container.State = types.StateStopped
	case types.StateError: container.State = types.StateError
	}

	container.Created, _ = time.Parse(time.RFC3339Nano, info.Created)
	container.Started, _ = time.Parse(time.RFC3339Nano, info.State.StartedAt)

	return container, nil
}
