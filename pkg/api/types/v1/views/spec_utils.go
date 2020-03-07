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

package views

import (
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/internal/util/resource"
)

type SpecView struct{}

func (sv *SpecView) NewSpecTemplateContainers(cl types.SpecTemplateContainers) SpecTemplateContainers {

	cs := SpecTemplateContainers{}
	for _, c := range cl {
		cs = append(cs, sv.NewContainer(c))
	}

	return cs
}

func (sv *SpecView) NewSpecTemplateVolumes(vl types.SpecTemplateVolumeList) SpecTemplateVolumeList {
	vs := SpecTemplateVolumeList{}
	for _, v := range vl {
		vs = append(vs, sv.NewVolume(v))
	}

	return vs
}

func (sv *SpecView) NewContainer(c *types.SpecTemplateContainer) *SpecTemplateContainer {

	s := new(SpecTemplateContainer)
	s.ID = c.ID
	s.Name = c.Name
	s.Role = c.Role
	s.AutoRemove = c.AutoRemove
	s.Labels = c.Labels
	s.Image = SpecTemplateContainerImage{
		Name: c.Image.Name,
		Sha:  c.Image.Sha,
		Secret: SpecTemplateContainerImageSecret{
			Name: c.Image.Secret.Name,
			Key:  c.Image.Secret.Key,
		},
		Policy: c.Image.Policy,
	}

	s.Ports = SpecTemplateContainerPorts{}
	for _, p := range c.Ports {
		s.Ports = append(s.Ports, &SpecTemplateContainerPort{
			ContainerPort: p.ContainerPort,
			HostPort:      p.HostPort,
			HostIP:        p.HostIP,
			Protocol:      p.Protocol,
		})
	}

	s.EnvVars = SpecTemplateContainerEnvs{}
	for _, e := range c.EnvVars {
		s.EnvVars = append(s.EnvVars, &SpecTemplateContainerEnv{
			Name:  e.Name,
			Value: e.Value,
			Secret: SpecTemplateContainerEnvSecret{
				Name: e.Secret.Name,
				Key:  e.Secret.Key,
			},
			Config: SpecTemplateContainerEnvConfig{
				Name: e.Secret.Name,
				Key:  e.Secret.Key,
			},
		})
	}

	s.Resources = SpecTemplateContainerResources{
		Limits: SpecTemplateContainerResource{
			CPU: resource.EncodeCpuResource(c.Resources.Limits.CPU),
			RAM: resource.EncodeMemoryResource(c.Resources.Limits.RAM),
		},
		Request: SpecTemplateContainerResource{
			CPU: resource.EncodeCpuResource(c.Resources.Request.CPU),
			RAM: resource.EncodeMemoryResource(c.Resources.Request.RAM),
		},
	}

	s.Exec = SpecTemplateContainerExec{
		Command:    c.Exec.Command,
		Entrypoint: c.Exec.Entrypoint,
		Workdir:    c.Exec.Workdir,
		Args:       c.Exec.Args,
	}

	s.Volumes = SpecTemplateContainerVolumes{}
	for _, v := range c.Volumes {
		s.Volumes = append(s.Volumes, &SpecTemplateContainerVolume{
			Name:      v.Name,
			MountPath: v.MountPath,
			SubPath:   v.SubPath,
			Mode:      v.Mode,
		})
	}

	s.Probes.LiveProbe = SpecTemplateContainerProbe{
		Exec:                c.Probes.LiveProbe.Exec,
		Socket:              c.Probes.LiveProbe.Socket,
		InitialDelaySeconds: c.Probes.LiveProbe.InitialDelaySeconds,
		TimeoutSeconds:      c.Probes.LiveProbe.TimeoutSeconds,
		PeriodSeconds:       c.Probes.LiveProbe.PeriodSeconds,
		ThresholdSuccess:    c.Probes.LiveProbe.ThresholdSuccess,
		ThresholdFailure:    c.Probes.LiveProbe.ThresholdFailure,
	}

	s.Probes.ReadProbe = SpecTemplateContainerProbe{
		Exec:                c.Probes.ReadProbe.Exec,
		Socket:              c.Probes.ReadProbe.Socket,
		InitialDelaySeconds: c.Probes.ReadProbe.InitialDelaySeconds,
		TimeoutSeconds:      c.Probes.ReadProbe.TimeoutSeconds,
		PeriodSeconds:       c.Probes.ReadProbe.PeriodSeconds,
		ThresholdSuccess:    c.Probes.ReadProbe.ThresholdSuccess,
		ThresholdFailure:    c.Probes.ReadProbe.ThresholdFailure,
	}

	s.Security = SpecTemplateContainerSecurity{
		Privileged: c.Security.Privileged,
		LinuxOptions: SpecTemplateContainerSecurityLinuxOptions{
			Level: c.Security.LinuxOptions.Level,
		},
		User: c.Security.User,
	}

	s.Network = SpecTemplateContainerNetwork{
		Hostname: c.Network.Hostname,
		Network:  c.Network.Network,
		Domain:   c.Network.Domain,
		Mode:     c.Network.Mode,
	}

	s.DNS = SpecTemplateContainerDNS{
		Server:  c.DNS.Server,
		Options: c.DNS.Options,
		Search:  c.DNS.Search,
	}

	s.ExtraHosts = c.ExtraHosts
	s.PublishAllPorts = c.PublishAllPorts
	s.Links = make([]SpecTemplateContainerLink, 0)

	for _, l := range c.Links {
		s.Links = append(s.Links, SpecTemplateContainerLink{
			Link:  l.Link,
			Alias: l.Alias,
		})
	}

	s.RestartPolicy = SpecTemplateRestartPolicy{
		Policy:  c.RestartPolicy.Policy,
		Attempt: c.RestartPolicy.Attempt,
	}

	return s
}

func (sv *SpecView) NewVolume(v *types.SpecTemplateVolume) *SpecTemplateVolume {
	s := new(SpecTemplateVolume)
	s.Name = v.Name
	s.Type = v.Type

	s.Config = SpecTemplateConfigVolume{
		Name: v.Config.Name,
	}

	for _, b := range v.Config.Binds {
		s.Config.Binds = append(s.Config.Binds, SpecTemplateConfigVolumeBind{
			Key:  b.Key,
			File: b.File,
		})
	}

	s.Secret = SpecTemplateSecretVolume{
		Name: v.Secret.Name,
	}

	for _, b := range v.Secret.Binds {
		s.Secret.Binds = append(s.Secret.Binds, SpecTemplateSecretVolumeBind{
			Key:  b.Key,
			File: b.File,
		})
	}

	s.Volume = SpecTemplateVolumeClaim{
		Name:    v.Volume.Name,
		Subpath: v.Volume.Subpath,
	}

	return s
}
