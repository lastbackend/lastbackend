package docker

import (
	docker "github.com/docker/docker/api/types"
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/satori/go.uuid"
	"time"
)

func GetPodMetaFromContainer(c docker.Container) types.PodMeta {
	log := context.Get().GetLogger()
	log.Debug("Docker: get pod meta from contanier")

	var meta types.PodMeta
	var err error

	meta = types.PodMeta{}

	ID, ok := c.Labels["LB_POD_ID"]
	if !ok {
		return meta
	}

	meta.ID, err = uuid.FromString(ID)
	if err != nil {
		return meta
	}

	meta.Owner, ok = c.Labels["LB_POD_OWNER"]
	if !ok {
		return meta
	}

	meta.Project, ok = c.Labels["LB_POD_PROJECT"]
	if !ok {
		return meta
	}

	meta.Service, ok = c.Labels["LB_POD_SERVICE"]
	if !ok {
		return meta
	}

	return meta

}

func GetContainerSpecFromContainer(c docker.ContainerJSON) *types.ContainerSpec {
	log := context.Get().GetLogger()
	log.Debug("Docker: get pod spec from contanier")
	var spec *types.ContainerSpec

	log.Debugf("Container Image: %s", c.Image)
	log.Debugf("Container Image name: %s", c.Config.Image)
	spec = &types.ContainerSpec{

		Network: types.ContainerNetworkSpec{
			Hostname: c.Config.Hostname,
			Domain:   c.Config.Domainname,
			//			Network: c.NetworkSettings.Networks["lo"].NetworkID,
			Mode: string(c.HostConfig.NetworkMode),
		},

		Ports: []types.ContainerPortSpec{},

		Command: []string(c.Config.Cmd),

		Entrypoint: []string(c.Config.Entrypoint),

		Envs: c.Config.Env,

		Labels: c.Config.Labels,

		DNS: types.ContainerDNSSpec{
			Server:  c.HostConfig.DNS,
			Options: c.HostConfig.DNSOptions,
			Search:  c.HostConfig.DNSSearch,
		},

		Quota: types.ContainerQuotaSpec{
			Memory:    c.HostConfig.Memory,
			CPUShares: c.HostConfig.CPUShares,
		},

		RestartPolicy: types.ContainerRestartPolicySpec{
			Name:    c.HostConfig.RestartPolicy.Name,
			Attempt: c.HostConfig.RestartPolicy.MaximumRetryCount,
		},
	}

	for port := range c.HostConfig.PortBindings {
		portSpec := types.ContainerPortSpec{
			ContainerPort: port.Int(),
			Protocol:      port.Proto(),
		}

		spec.Ports = append(spec.Ports, portSpec)
	}

	return spec
}

func GetContainer(dc docker.Container, info docker.ContainerJSON) *types.Container {
	log := context.Get().GetLogger()
	log.Debug("Docker: convert container format")

	var c *types.Container

	c = &types.Container{
		ID:      dc.ID,
		State:   dc.State,
		Status:  dc.Status,
		Created: time.Unix(dc.Created, 0),
	}

	t, _ := time.Parse(time.RFC3339Nano, info.State.StartedAt)
	c.Started = t
	c.Image.ID = dc.ImageID
	c.Image.Name = dc.Image
	return c

}
