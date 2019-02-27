//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/util/resource"
	"strings"
)

type ManifestView struct{}

func (mv *ManifestView) NewManifestSpecTemplate(obj types.SpecTemplate) ManifestSpecTemplate {

	mst := ManifestSpecTemplate{
		Containers: make([]ManifestSpecTemplateContainer, 0),
		Volumes:    make([]ManifestSpecTemplateVolume, 0),
	}

	for _, s := range obj.Containers {

		c := ManifestSpecTemplateContainer{
			Name:       s.Name,
			Command:    strings.Join(s.Exec.Command, " "),
			Workdir:    s.Exec.Workdir,
			Args:       s.Exec.Args,
			Entrypoint: strings.Join(s.Exec.Entrypoint, " "),
		}

		for _, p := range s.Ports {
			c.Ports = append(c.Ports, fmt.Sprintf("%d/%s", p.ContainerPort, p.Protocol))
		}

		for _, env := range s.EnvVars {
			c.Env = append(c.Env, ManifestSpecTemplateContainerEnv{
				Name:  env.Name,
				Value: env.Value,
				Secret: &ManifestSpecTemplateContainerEnvSecret{
					Name: env.Secret.Name,
					Key:  env.Secret.Key,
				},
				Config: &ManifestSpecTemplateContainerEnvConfig{
					Name: env.Config.Name,
					Key:  env.Config.Key,
				},
			})
		}

		c.Image = new(ManifestSpecTemplateContainerImage)
		c.Image.Name = s.Image.Name
		c.Image.Secret.Name = s.Image.Secret.Name
		c.Image.Secret.Key = s.Image.Secret.Key

		for _, volume := range s.Volumes {
			c.Volumes = append(c.Volumes, ManifestSpecTemplateContainerVolume{
				Name: volume.Name,
				Mode: volume.Mode,
				Path: volume.Path,
			})
		}

		c.Resources = new(ManifestSpecTemplateContainerResources)
		c.Resources.Limits = new(ManifestSpecTemplateContainerResource)
		c.Resources.Limits.RAM = resource.EncodeMemoryResource(s.Resources.Limits.RAM)
		c.Resources.Limits.CPU = resource.EncodeCpuResource(s.Resources.Limits.CPU)
		c.Resources.Request = new(ManifestSpecTemplateContainerResource)
		c.Resources.Request.RAM = resource.EncodeMemoryResource(s.Resources.Request.RAM)
		c.Resources.Request.CPU = resource.EncodeCpuResource(s.Resources.Request.CPU)

		mst.Containers = append(mst.Containers, c)
	}

	for _, s := range obj.Volumes {
		v := ManifestSpecTemplateVolume{
			Name: s.Name,
			Type: s.Type,
			Volume: &ManifestSpecTemplateVolumeClaim{
				Name:    s.Volume.Name,
				Subpath: s.Volume.Subpath,
			},
			Secret: &ManifestSpecTemplateSecretVolume{
				Name:  s.Secret.Name,
				Binds: make([]ManifestSpecTemplateSecretVolumeBind, 0),
			},
			Config: &ManifestSpecTemplateConfigVolume{
				Name:  s.Config.Name,
				Binds: make([]ManifestSpecTemplateConfigVolumeBind, 0),
			},
		}

		for _, b := range s.Secret.Binds {
			v.Secret.Binds = append(v.Secret.Binds, ManifestSpecTemplateSecretVolumeBind{
				Key:  b.Key,
				File: b.File,
			})
		}

		for _, b := range s.Config.Binds {
			v.Config.Binds = append(v.Config.Binds, ManifestSpecTemplateConfigVolumeBind{
				Key:  b.Key,
				File: b.File,
			})
		}

		mst.Volumes = append(mst.Volumes, v)
	}

	return mst
}

func (mv *ManifestView) NewManifestSpecSelector(obj types.SpecSelector) ManifestSpecSelector {
	return ManifestSpecSelector{
		Node:   obj.Node,
		Labels: obj.Labels,
	}
}

func (mv *ManifestView) NewManifestSpecRuntime(obj types.SpecRuntime) ManifestSpecRuntime {

	mfr := ManifestSpecRuntime{}
	mfr.Services = obj.Services
	mfr.Tasks = make([]ManifestSpecRuntimeTask, 0)
	for _, task := range obj.Tasks {
		mft := ManifestSpecRuntimeTask{
			Name:      task.Name,
			Container: task.Container,
			Env:       make([]ManifestSpecTemplateContainerEnv, 0),
			Commands:  make([]ManifestSpecRuntimeTaskCommand, 0),
		}

		for _, e := range task.EnvVars {
			env := ManifestSpecTemplateContainerEnv{
				Name:  e.Name,
				Value: e.Value,
				Secret: &ManifestSpecTemplateContainerEnvSecret{
					Name: e.Secret.Name,
					Key:  e.Secret.Key,
				},
				Config: &ManifestSpecTemplateContainerEnvConfig{
					Name: e.Config.Name,
					Key:  e.Config.Key,
				},
			}
			mft.Env = append(mft.Env, env)
		}

		for _, c := range task.Commands {
			cmd := ManifestSpecRuntimeTaskCommand{
				Command:    strings.Join(c.Command, " "),
				Entrypoint: strings.Join(c.Entrypoint, " "),
				Workdir:    c.Workdir,
				Args:       c.Args,
			}
			mft.Commands = append(mft.Commands, cmd)
		}

		mfr.Tasks = append(mfr.Tasks, mft)
	}

	return mfr
}
