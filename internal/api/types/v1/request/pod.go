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

package request

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"github.com/lastbackend/lastbackend/internal/util/resource"
	"gopkg.in/yaml.v2"
	"strings"
	"time"
)

type PodManifest struct {
	Meta PodManifestMeta `json:"meta,omitempty" yaml:"meta,omitempty"`
	Spec PodManifestSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type PodManifestMeta struct {
	RuntimeMeta `yaml:",inline"`
}

type PodManifestSpec struct {
	Selector *ManifestSpecSelector `json:"selector,omitempty" yaml:"selector,omitempty"`
	Runtime  *ManifestSpecRuntime  `json:"runtime,omitempty" yaml:"runtime,omitempty"`
	Template *ManifestSpecTemplate `json:"template,omitempty" yaml:"template,omitempty"`
}

func (s *PodManifest) FromJson(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *PodManifest) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *PodManifest) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, s)
}

func (s *PodManifest) ToYaml() ([]byte, error) {
	return yaml.Marshal(s)
}

func (s *PodManifest) SetPodMeta(svc *types.Pod) {

	if svc.Meta.Name == types.EmptyString {
		svc.Meta.Name = *s.Meta.Name
	}

	if s.Meta.Description != nil {
		svc.Meta.Description = *s.Meta.Description
	}

	if s.Meta.Labels != nil {
		svc.Meta.Labels = s.Meta.Labels
	}

}

func (s *PodManifest) SetPodSpec(pod *types.Pod) {

	if s.Spec.Selector != nil {

		if s.Spec.Selector.Node != types.EmptyString && pod.Spec.Selector.Node != s.Spec.Selector.Node {
			pod.Spec.Selector.Node = s.Spec.Selector.Node
		}

		if s.Spec.Selector.Labels != nil {
			pod.Spec.Selector.Labels = s.Spec.Selector.Labels
		}

	}

	if s.Spec.Runtime != nil {

		// update services
		if len(s.Spec.Runtime.Services) != len(pod.Spec.Runtime.Services) {
			pod.Spec.Runtime.Services = s.Spec.Runtime.Services
			pod.Spec.Template.Updated = time.Now()
		} else {
			for _, svc := range s.Spec.Runtime.Services {
				var f = false

				for _, ssvc := range pod.Spec.Runtime.Services {
					if svc == ssvc {
						f = true
					}
				}

				if !f {
					pod.Spec.Runtime.Services = s.Spec.Runtime.Services
					pod.Spec.Template.Updated = time.Now()
					break
				}
			}
		}

		// update tasks
		var tne = false
		if len(s.Spec.Runtime.Tasks) != len(pod.Spec.Runtime.Tasks) {
			tne = true
		}

		if tne {
			for _, st := range s.Spec.Runtime.Tasks {

				var f = false
				for _, pt := range pod.Spec.Runtime.Tasks {

					if st.Name != pt.Name && st.Container != pt.Container {
						continue
					}

					// check envs
					if len(st.Env) != len(pt.EnvVars) {
						break
					}

					for _, ce := range st.Env {
						var f = false
						for _, se := range pt.EnvVars {

							if ce.Name != se.Name {
								continue
							}

							if se.Value != ce.Value {
								continue
							}

							if se.Secret.Name != ce.Secret.Name || se.Secret.Key != ce.Secret.Key {
								continue
							}

							if se.Config.Name != ce.Config.Name || se.Config.Key != ce.Config.Key {
								continue
							}

							f = true
						}

						if !f {
							break
						}
					}

					// check commands
					if len(st.Commands) != len(pt.Commands) {
						break
					}
					for _, cc := range st.Commands {
						var f = false

						for _, sc := range pt.Commands {

							if sc != cc {
								continue
							}

							f = true
						}

						if !f {
							break
						}
					}

					f = true
				}

				if !f {
					tne = true
					break
				}
			}
		}

		if tne {
			pod.Spec.Runtime.Tasks = make([]types.SpecRuntimeTask, 0)
			for _, t := range s.Spec.Runtime.Tasks {

				task := types.SpecRuntimeTask{
					Name:      t.Name,
					Container: t.Container,
					EnvVars:   make(types.SpecTemplateContainerEnvs, 0),
					Commands:  make([]string, 0),
				}

				for _, env := range t.Env {
					task.EnvVars = append(task.EnvVars, &types.SpecTemplateContainerEnv{
						Name:  env.Name,
						Value: env.Value,
						Secret: types.SpecTemplateContainerEnvSecret{
							Name: env.Secret.Name,
							Key:  env.Secret.Key,
						},
						Config: types.SpecTemplateContainerEnvConfig{
							Name: env.Config.Name,
							Key:  env.Config.Key,
						},
					})
				}

				task.Commands = make([]string, 0)
				if t.Commands != nil {
					task.Commands = t.Commands
				}

				pod.Spec.Runtime.Tasks = append(pod.Spec.Runtime.Tasks, task)
			}

			pod.Spec.Template.Updated = time.Now()
		}

	}

	if s.Spec.Template != nil {

		for _, c := range s.Spec.Template.Containers {

			var (
				f    = false
				spec *types.SpecTemplateContainer
			)

			for _, sc := range pod.Spec.Template.Containers {
				if c.Name == sc.Name {
					f = true
					spec = sc
				}
			}

			if spec == nil {
				spec = new(types.SpecTemplateContainer)
			}

			if spec.Name == types.EmptyString {
				spec.Name = c.Name
				pod.Spec.Template.Updated = time.Now()
			}

			if spec.Image.Name != c.Image.Name {
				spec.Image.Name = c.Image.Name
				pod.Spec.Template.Updated = time.Now()
			}

			if spec.Image.Sha != c.Image.Sha {
				spec.Image.Sha = c.Image.Sha
				pod.Spec.Template.Updated = time.Now()
			}

			if spec.Image.Secret.Name != c.Image.Secret.Name || spec.Image.Secret.Key != c.Image.Secret.Key {
				spec.Image.Secret.Name = c.Image.Secret.Name
				spec.Image.Secret.Key = c.Image.Secret.Key
				pod.Spec.Template.Updated = time.Now()
			}

			if strings.Join(spec.Exec.Command, " ") != c.Command {
				spec.Exec.Command = strings.Split(c.Command, " ")
				pod.Spec.Template.Updated = time.Now()
			}

			if strings.Join(spec.Exec.Args, "") != strings.Join(c.Args, "") {
				spec.Exec.Args = c.Args
				pod.Spec.Template.Updated = time.Now()
			}

			if strings.Join(spec.Exec.Entrypoint, " ") != c.Entrypoint {
				spec.Exec.Entrypoint = strings.Split(c.Entrypoint, " ")
				pod.Spec.Template.Updated = time.Now()
			}

			if spec.Exec.Workdir != c.Workdir {
				spec.Exec.Workdir = c.Workdir
				pod.Spec.Template.Updated = time.Now()
			}

			if spec.Security.Privileged != c.Security.Privileged {
				spec.Security.Privileged = c.Security.Privileged
				pod.Spec.Template.Updated = time.Now()
			}

			for _, ce := range c.Env {
				var f = false

				for _, se := range spec.EnvVars {
					if ce.Name == se.Name {
						f = true
						if se.Value != ce.Value {
							se.Value = ce.Value
							pod.Spec.Template.Updated = time.Now()
						}

						if se.Secret.Name != ce.Secret.Name || se.Secret.Key != ce.Secret.Key {
							se.Secret.Name = ce.Secret.Name
							se.Secret.Key = ce.Secret.Key
							pod.Spec.Template.Updated = time.Now()
						}

						if se.Config.Name != ce.Config.Name || se.Config.Key != ce.Config.Key {
							se.Config.Name = ce.Config.Name
							se.Config.Key = ce.Config.Key
							pod.Spec.Template.Updated = time.Now()
						}
					}
				}

				if !f {
					spec.EnvVars = append(spec.EnvVars, &types.SpecTemplateContainerEnv{
						Name:  ce.Name,
						Value: ce.Value,
						Secret: types.SpecTemplateContainerEnvSecret{
							Name: ce.Secret.Name,
							Key:  ce.Secret.Key,
						},
						Config: types.SpecTemplateContainerEnvConfig{
							Name: ce.Config.Name,
							Key:  ce.Config.Key,
						},
					})
				}
			}

			var envs = make([]*types.SpecTemplateContainerEnv, 0)
			for _, se := range spec.EnvVars {
				for _, ce := range c.Env {
					if ce.Name == se.Name {
						envs = append(envs, se)
						break
					}
				}
			}

			if len(spec.EnvVars) != len(envs) {
				pod.Spec.Template.Updated = time.Now()
			}
			spec.EnvVars = envs

			resourcesRequestRam, _ := resource.DecodeMemoryResource(c.Resources.Request.RAM)
			resourcesRequestCPU, _ := resource.DecodeCpuResource(c.Resources.Request.CPU)

			resourcesLimitsRam, _ := resource.DecodeMemoryResource(c.Resources.Limits.RAM)
			resourcesLimitsCPU, _ := resource.DecodeCpuResource(c.Resources.Limits.CPU)

			if resourcesRequestRam != spec.Resources.Request.RAM ||
				resourcesRequestCPU != spec.Resources.Request.CPU {
				spec.Resources.Request.RAM = resourcesRequestRam
				spec.Resources.Request.CPU = resourcesRequestCPU
				pod.Spec.Template.Updated = time.Now()
			}

			if resourcesLimitsRam != spec.Resources.Limits.RAM ||
				resourcesLimitsCPU != spec.Resources.Limits.CPU {
				spec.Resources.Limits.RAM = resourcesLimitsRam
				spec.Resources.Limits.CPU = resourcesLimitsCPU
				pod.Spec.Template.Updated = time.Now()
			}

			for _, v := range c.Volumes {

				var f = false
				for _, sv := range spec.Volumes {

					if v.Name == sv.Name {
						f = true
						if sv.Mode != v.Mode || sv.MountPath != v.MountPath || sv.SubPath != v.SubPath {
							sv.Mode = v.Mode
							sv.MountPath = v.MountPath
							sv.SubPath = v.SubPath
							pod.Spec.Template.Updated = time.Now()
						}

					}
				}
				if !f {
					spec.Volumes = append(spec.Volumes, &types.SpecTemplateContainerVolume{
						Name:      v.Name,
						Mode:      v.Mode,
						MountPath: v.MountPath,
						SubPath:   v.SubPath,
					})
				}
			}

			vlms := make([]*types.SpecTemplateContainerVolume, 0)
			for _, sv := range spec.Volumes {
				for _, cv := range c.Volumes {
					if sv.Name == cv.Name {
						vlms = append(vlms, sv)
						break
					}
				}
			}

			if len(vlms) != len(spec.Volumes) {
				pod.Spec.Template.Updated = time.Now()
			}

			spec.Volumes = vlms

			if !f {
				pod.Spec.Template.Containers = append(pod.Spec.Template.Containers, spec)
			}

		}

		var spcs = make([]*types.SpecTemplateContainer, 0)
		for _, ss := range pod.Spec.Template.Containers {
			for _, cs := range s.Spec.Template.Containers {
				if ss.Name == cs.Name {
					spcs = append(spcs, ss)
				}
			}
		}

		if len(spcs) != len(pod.Spec.Template.Containers) {
			pod.Spec.Template.Updated = time.Now()
		}

		pod.Spec.Template.Containers = spcs

		for _, v := range s.Spec.Template.Volumes {

			var (
				f    = false
				spec *types.SpecTemplateVolume
			)

			for _, sv := range pod.Spec.Template.Volumes {
				if v.Name == sv.Name {
					f = true
					spec = sv
				}
			}

			if spec == nil {
				spec = new(types.SpecTemplateVolume)
			}

			if spec.Name == types.EmptyString {
				spec.Name = v.Name
				pod.Spec.Template.Updated = time.Now()
			}

			if v.Type != spec.Type || v.Secret.Name != spec.Secret.Name || v.Config.Name != spec.Config.Name {
				spec.Type = v.Type
				spec.Secret.Name = v.Secret.Name
				spec.Config.Name = v.Config.Name
				pod.Spec.Template.Updated = time.Now()
			}

			var e = true
			for _, vf := range v.Secret.Binds {

				var f = false
				for _, sf := range spec.Secret.Binds {
					if (vf.Key == sf.Key) && (vf.File == sf.File) {
						f = true
						break
					}
				}

				if !f {
					e = false
					break
				}

			}

			if !e {
				spec.Secret.Binds = make([]types.SpecTemplateSecretVolumeBind, 0)
				for _, v := range v.Secret.Binds {
					spec.Secret.Binds = append(spec.Secret.Binds, types.SpecTemplateSecretVolumeBind{
						Key:  v.Key,
						File: v.File,
					})
				}
				pod.Spec.Template.Updated = time.Now()
			}

			var ec = true
			for _, vf := range v.Config.Binds {

				var f = false
				for _, sf := range spec.Config.Binds {
					if (vf.Key == sf.Key) && (vf.File == sf.File) {
						f = true
						break
					}
				}

				if !f {
					ec = false
					break
				}

			}

			if !ec {
				spec.Config.Binds = make([]types.SpecTemplateConfigVolumeBind, 0)
				for _, v := range v.Config.Binds {
					spec.Config.Binds = append(spec.Config.Binds, types.SpecTemplateConfigVolumeBind{
						Key:  v.Key,
						File: v.File,
					})
				}
				pod.Spec.Template.Updated = time.Now()
			}

			if !f {
				pod.Spec.Template.Volumes = append(pod.Spec.Template.Volumes, spec)
			}

		}

		var vlms = make([]*types.SpecTemplateVolume, 0)
		for _, ss := range pod.Spec.Template.Volumes {
			for _, cs := range s.Spec.Template.Volumes {
				if ss.Name == cs.Name {
					vlms = append(vlms, ss)
				}
			}
		}

		if len(vlms) != len(pod.Spec.Template.Volumes) {
			pod.Spec.Template.Updated = time.Now()
		}

		pod.Spec.Template.Volumes = vlms

	}

}

func (s *PodManifest) GetManifest() *types.PodManifest {
	sm := new(types.PodManifest)
	if s.Spec.Selector != nil {
		sm.Selector = s.Spec.Selector.GetSpec()
	}

	if s.Spec.Runtime != nil {
		sm.Runtime = s.Spec.Runtime.GetSpec()
	}

	if s.Spec.Template != nil {
		sm.Template = s.Spec.Template.GetSpec()
	}

	return sm
}

type PodLogsOptions struct {
	Container string `json:"container"`
	Follow    bool   `json:"follow"`
}
