package docker

import (
	"fmt"
	"github.com/deployithq/deployit/drivers/interfaces"
	"github.com/fsouza/go-dockerclient"
	"strconv"
	"strings"
)

func convertImage(i *docker.Image) (interfaces.Image, error) {
	image := interfaces.Image{}

	return image, nil
}

func ConvertContainer(info *docker.Container) (interfaces.Container, error) {

	cn := interfaces.Container{
		Ports: []interfaces.Port{},
	}

	cn.CID = info.ID
	cn.Name = info.Name

	cn.State.Running = info.State.Running
	cn.State.Paused = info.State.Paused
	cn.State.Restarting = info.State.Restarting
	cn.State.OOMKilled = info.State.OOMKilled
	cn.State.Pid = info.State.Pid
	cn.State.ExitCode = info.State.ExitCode
	cn.State.Error = info.State.Error
	cn.State.Started = info.State.StartedAt
	cn.State.Finished = info.State.FinishedAt

	for port := range info.NetworkSettings.Ports {

		cPort, _ := strconv.ParseInt(port.Port(), 0, 64)
		var cHost int64
		for index := range info.NetworkSettings.Ports[port] {

			cHost, _ = strconv.ParseInt(info.NetworkSettings.Ports[port][index].HostPort, 0, 64)

			cn.Ports = append(cn.Ports, interfaces.Port{
				Container: cPort,
				Host:      cHost,
				Protocol:  port.Proto(),
			})
		}
	}

	return cn, nil
}

func CreateConfig(c interfaces.Config) docker.Config {

	config := docker.Config{}
	config.Cmd = c.Cmd

	config.Memory = c.Memory * 1024 * 1024

	config.Env = c.Env
	config.Entrypoint = c.Entrypoint

	config.Image = c.Image

	config.ExposedPorts = make(map[docker.Port]struct{})
	config.Volumes = make(map[string]struct{})

	for _, port := range c.Ports {
		item := strings.Split(port, ":")
		containerPort, _ := strconv.ParseInt(item[1], 10, 64)

		key := docker.Port(fmt.Sprintf("%d/tcp", containerPort))
		config.ExposedPorts[key] = struct{}{}
	}

	for _, volume := range c.Volumes {
		item := strings.Split(volume, ":")
		config.Volumes[item[1]] = struct{}{}
	}

	return config
}

func CreateHostConfig(c interfaces.HostConfig) docker.HostConfig {
	host := docker.HostConfig{}

	host.Privileged = c.Privileged

	host.RestartPolicy.Name = c.RestartPolicy.Name
	host.RestartPolicy.MaximumRetryCount = c.RestartPolicy.Attempt
	host.Memory = c.Memory * 1024 * 1024
	host.Binds = c.Binds

	host.PortBindings = make(map[docker.Port][]docker.PortBinding)

	for _, port := range c.Ports {
		item := strings.Split(port, ":")
		hostPort, _ := strconv.ParseInt(item[0], 10, 64)
		containerPort, _ := strconv.ParseInt(item[1], 10, 64)
		key := docker.Port(fmt.Sprintf("%d/tcp", containerPort))

		host.PortBindings[key] = append(host.PortBindings[key], docker.PortBinding{HostPort: fmt.Sprintf("%d", hostPort)})
	}

	return host
}
