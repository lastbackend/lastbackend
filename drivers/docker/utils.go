package docker

import (
	"encoding/json"
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

	for key := range info.Config.Env {
		result := strings.Split(info.Config.Env[key], "=")
		if result[0] == `LB` {
			if err := json.Unmarshal([]byte(result[1]), &cn.LB); err != nil {
				return cn, err
			}
			break
		}
	}

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

	config.Memory = c.Memory.Total * 1024 * 1024
	config.MemorySwap = c.Memory.Swap * 1024 * 1024
	config.MemoryReservation = c.Memory.Reservation * 1024 * 1024
	config.KernelMemory = c.Memory.Kernel * 1024 * 1024

	config.CPUShares = c.CPU.Shares
	config.CPUSet = c.CPU.Set

	config.Env = c.Env
	config.Entrypoint = c.Entrypoint

	config.DNS = c.DNS.Server
	config.Image = c.Image

	config.ExposedPorts = make(map[docker.Port]struct{})
	config.Volumes = make(map[string]struct{})

	for _, port := range c.Ports {
		key := docker.Port(fmt.Sprintf("%d/%s", port.Container, port.Protocol))
		config.ExposedPorts[key] = struct{}{}
	}

	for _, volume := range c.Volumes {
		config.Volumes[volume.Container] = struct{}{}
	}

	return config
}

func CreateHostconfig(c interfaces.HostConfig) docker.HostConfig {
	host := docker.HostConfig{}

	host.Privileged = c.Privileged

	host.DNS = c.DNS.Server
	host.DNSOptions = c.DNS.Options
	host.DNSSearch = c.DNS.Search

	host.RestartPolicy.Name = c.RestartPolicy.Name
	host.RestartPolicy.MaximumRetryCount = c.RestartPolicy.Attempt
	host.Memory = c.Memory.Total * 1024 * 1024
	host.MemorySwap = c.Memory.Swap

	host.CPUShares = c.CPU.Shares
	host.CPUSet = c.CPU.Set
	host.CPUSetCPUs = c.CPU.CPUs
	host.CPUSetMEMs = c.CPU.MEMs
	host.CPUQuota = c.CPU.Quota
	host.CPUPeriod = c.CPU.Period
	host.Binds = c.Binds

	host.PortBindings = make(map[docker.Port][]docker.PortBinding)

	for _, port := range c.Ports {
		key := docker.Port(fmt.Sprintf("%d/%s", port.Container, port.Protocol))
		host.PortBindings[key] = append(host.PortBindings[key], docker.PortBinding{HostPort: fmt.Sprintf("%d", port.Host)})
	}

	return host
}
