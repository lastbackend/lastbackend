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
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/util/compare"
	"github.com/lastbackend/lastbackend/pkg/util/resource"
	"strings"
	"time"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type ManifestSpecSelector struct {
	Node   string            `json:"node,omitempty" yaml:"node,omitempty"`
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

type ManifestSpecNetwork struct {
	IP    *string  `json:"ip,omitempty" yaml:"ip,omitempty"`
	Ports []string `json:"ports,omitempty" yaml:"ports,omitempty"`
}

type ManifestSpecStrategy struct {
	Type *string `json:"type,omitempty" yaml:"type,omitempty"`
}

type ManifestSpecRuntime struct {
	Services []string                  `json:"services,omitempty"`
	Tasks    []ManifestSpecRuntimeTask `json:"tasks,omitempty"`
}

type ManifestSpecRuntimeTask struct {
	Name      string                             `json:"name"`
	Container string                             `json:"container" yaml:"container"`
	Env       []ManifestSpecTemplateContainerEnv `json:"env,omitempty" yaml:"env,omitempty"`
	Commands  []ManifestSpecRuntimeTaskCommand   `json:"commands,omitempty" yaml:"commands,omitempty"`
}

type ManifestSpecRuntimeTaskCommand struct {
	Command    string   `json:"command,omitempty" yaml:"command,omitempty"`
	Workdir    string   `json:"workdir,omitempty" yaml:"workdir,omitempty"`
	Entrypoint string   `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
	Args       []string `json:"args,omitempty" yaml:"args,omitempty"`
}

type ManifestSpecTemplate struct {
	Containers []ManifestSpecTemplateContainer `json:"containers,omitempty" yaml:"containers,omitempty"`
	Volumes    []ManifestSpecTemplateVolume    `json:"volumes,omitempty" yaml:"volumes,omitempty"`
}

type ManifestSpecTemplateContainer struct {
	Name          string                                  `json:"name,omitempty" yaml:"name,omitempty"`
	Command       string                                  `json:"command,omitempty" yaml:"command,omitempty"`
	Workdir       string                                  `json:"workdir,omitempty" yaml:"workdir,omitempty"`
	Entrypoint    string                                  `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
	Args          []string                                `json:"args,omitempty" yaml:"args,omitempty"`
	Ports         []string                                `json:"ports,omitempty" yaml:"ports,omitempty"`
	Env           []ManifestSpecTemplateContainerEnv      `json:"env,omitempty" yaml:"env,omitempty"`
	Volumes       []ManifestSpecTemplateContainerVolume   `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Image         *ManifestSpecTemplateContainerImage     `json:"image,omitempty" yaml:"image,omitempty"`
	Resources     *ManifestSpecTemplateContainerResources `json:"resources,omitempty" yaml:"resources,omitempty"`
	RestartPolicy *ManifestSpecTemplateRestartPolicy      `json:"restart,omitempty" yaml:"restart,omitempty"`
	Security      *ManifestSpecSecurity                   `json:"security,omitempty" yaml:"security,omitempty"`
}

type ManifestSpecTemplateContainerEnv struct {
	Name   string                                  `json:"name,omitempty" yaml:"name,omitempty"`
	Value  string                                  `json:"value,omitempty" yaml:"value,omitempty"`
	Secret *ManifestSpecTemplateContainerEnvSecret `json:"secret,omitempty" yaml:"secret,omitempty"`
	Config *ManifestSpecTemplateContainerEnvConfig `json:"config,omitempty" yaml:"config,omitempty"`
}

type ManifestSpecTemplateContainerEnvSecret struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Key  string `json:"key,omitempty" yaml:"key,omitempty"`
}

type ManifestSpecTemplateContainerEnvConfig struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Key  string `json:"key,omitempty" yaml:"key,omitempty"`
}

type ManifestSpecTemplateContainerImage struct {
	Name   string                                   `json:"name,omitempty" yaml:"name,omitempty"`
	Secret ManifestSpecTemplateContainerImageSecret `json:"secret,omitempty" yaml:"secret,omitempty"`
}

type ManifestSpecTemplateContainerImageSecret struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Key  string `json:"key,omitempty" yaml:"key,omitempty"`
}

type ManifestSpecTemplateContainerResources struct {
	// Limit resources
	Limits *ManifestSpecTemplateContainerResource `json:"limits,omitempty" yaml:"limits,omitempty"`
	// Request resources
	Request *ManifestSpecTemplateContainerResource `json:"quota,omitempty" yaml:"quota,omitempty"`
}

type ManifestSpecTemplateContainerVolume struct {
	// Volume name
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Volume mount mode
	Mode string `json:"mode,omitempty" yaml:"mode,omitempty"`
	// Volume mount path
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
}

type ManifestSpecTemplateContainerResource struct {
	// CPU resource option
	CPU string `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	// RAM resource option
	RAM string `json:"ram,omitempty" yaml:"ram,omitempty"`
}

type ManifestSpecTemplateVolume struct {
	// Template volume name
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Template volume types
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// Template volume from persistent volume
	Volume *ManifestSpecTemplateVolumeClaim `json:"volume,omitempty" yaml:"volume,omitempty"`
	// Template volume from secret type
	Secret *ManifestSpecTemplateSecretVolume `json:"secret,omitempty" yaml:"secret,omitempty"`
	// Template volume from config type
	Config *ManifestSpecTemplateConfigVolume `json:"config,omitempty" yaml:"config,omitempty"`
}

type ManifestSpecTemplateVolumeClaim struct {
	// Persistent volume name to mount
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Persistent volume subpath
	Subpath string `json:"subpath,omitempty" yaml:"subpath,omitempty"`
}

type ManifestSpecTemplateSecretVolume struct {
	// Secret name to mount
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Secret file key
	Binds []ManifestSpecTemplateSecretVolumeBind `json:"binds,omitempty" yaml:"bindsk,omitempty"`
}

type ManifestSpecTemplateSecretVolumeBind struct {
	Key  string `json:"key,omitempty" yaml:"key,omitempty"`
	File string `json:"file,omitempty" yaml:"file,omitempty"`
}

type ManifestSpecTemplateConfigVolume struct {
	// Secret name to mount
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Secret file key
	Binds []ManifestSpecTemplateConfigVolumeBind `json:"binds,omitempty" yaml:"binds,omitempty"`
}

type ManifestSpecTemplateConfigVolumeBind struct {
	Key  string `json:"key,omitempty" yaml:"key"`
	File string `json:"file,omitempty" yaml:"file"`
}

type ManifestSpecTemplateRestartPolicy struct {
	Policy  string `json:"policy,omitempty" yaml:"policy"`
	Attempt int    `json:"attempt,omitempty" yaml:"attempt"`
}

type ManifestSpecSecurity struct {
	Privileged bool `json:"privileged"`
}

func (m ManifestSpecSelector) GetSpec() types.SpecSelector {
	s := types.SpecSelector{}

	s.Node = m.Node
	s.Labels = m.Labels

	return s
}

func (m ManifestSpecTemplate) GetSpec() types.SpecTemplate {
	var s = types.SpecTemplate{}

	for _, t := range m.Containers {
		sp := t.GetSpec()
		s.Containers = append(s.Containers, &sp)
	}

	for _, t := range m.Volumes {
		sp := t.GetSpec()
		s.Volumes = append(s.Volumes, &sp)
	}

	return s
}

func (m ManifestSpecRuntime) GetSpec() types.SpecRuntime {
	var s = types.SpecRuntime{}

	s.Services = m.Services

	for _, t := range m.Tasks {
		ts := t.GetSpec()
		s.Tasks = append(s.Tasks, ts)
	}

	return s
}

func (m ManifestSpecRuntimeTask) GetSpec() types.SpecRuntimeTask {
	s := types.SpecRuntimeTask{}

	s.Name = m.Name
	s.Container = m.Container

	for _, e := range m.Env {
		env := types.SpecTemplateContainerEnv{
			Name:  e.Name,
			Value: e.Value,
		}

		if e.Secret != nil {
			env.Secret = types.SpecTemplateContainerEnvSecret{
				Name: e.Secret.Name,
				Key:  e.Secret.Key,
			}
		}

		if e.Config != nil {
			env.Config = types.SpecTemplateContainerEnvConfig{
				Name: e.Config.Name,
				Key:  e.Config.Key,
			}
		}

		s.EnvVars = append(s.EnvVars, &env)
	}

	for _, c := range m.Commands {
		cmd := types.SpecTemplateContainerExec{}

		if c.Command != types.EmptyString {
			cmd.Command = strings.Split(c.Command, " ")
		}
		cmd.Args = c.Args
		cmd.Workdir = c.Workdir

		if c.Entrypoint != types.EmptyString {
			cmd.Entrypoint = strings.Split(c.Entrypoint, " ")
		}

		s.Commands = append(s.Commands, cmd)
	}

	return s
}

func (m ManifestSpecTemplateVolume) GetSpec() types.SpecTemplateVolume {

	s := types.SpecTemplateVolume{
		Name: m.Name,
		Type: m.Type,
	}

	if m.Volume != nil {
		s.Volume = types.SpecTemplateVolumeClaim{
			Name:    m.Volume.Name,
			Subpath: m.Volume.Subpath,
		}
	}

	if m.Secret != nil {
		s.Secret = types.SpecTemplateSecretVolume{
			Name:  m.Secret.Name,
			Binds: make([]types.SpecTemplateSecretVolumeBind, 0),
		}

		for _, b := range m.Secret.Binds {
			s.Secret.Binds = append(s.Secret.Binds, types.SpecTemplateSecretVolumeBind{
				Key:  b.Key,
				File: b.File,
			})
		}
	}

	if m.Config != nil {
		s.Config = types.SpecTemplateConfigVolume{
			Name:  m.Config.Name,
			Binds: make([]types.SpecTemplateConfigVolumeBind, 0),
		}

		for _, b := range m.Config.Binds {
			s.Config.Binds = append(s.Config.Binds, types.SpecTemplateConfigVolumeBind{
				Key:  b.Key,
				File: b.File,
			})
		}
	}

	return s
}

func (m ManifestSpecTemplateContainer) GetSpec() types.SpecTemplateContainer {
	s := types.SpecTemplateContainer{}
	s.Name = m.Name

	s.RestartPolicy.Policy = "always"

	if m.RestartPolicy != nil {
		s.RestartPolicy.Policy = m.RestartPolicy.Policy
		s.RestartPolicy.Attempt = m.RestartPolicy.Attempt
	}

	if m.Command != types.EmptyString {
		s.Exec.Command = strings.Split(m.Command, " ")
	}
	s.Exec.Args = m.Args
	s.Exec.Workdir = m.Workdir

	if m.Entrypoint != types.EmptyString {
		s.Exec.Entrypoint = strings.Split(m.Entrypoint, " ")
	}

	for _, p := range m.Ports {
		port := new(types.SpecTemplateContainerPort)
		port.Parse(p)
		s.Ports = append(s.Ports, port)
	}

	for _, e := range m.Env {
		env := &types.SpecTemplateContainerEnv{
			Name:  e.Name,
			Value: e.Value,
		}

		if e.Secret != nil {
			env.Secret = types.SpecTemplateContainerEnvSecret{
				Name: e.Secret.Name,
				Key:  e.Secret.Key,
			}
		}

		if e.Config != nil {
			env.Config = types.SpecTemplateContainerEnvConfig{
				Name: e.Config.Name,
				Key:  e.Config.Key,
			}
		}

		s.EnvVars = append(s.EnvVars, env)
	}

	s.Image.Name = m.Image.Name
	s.Image.Secret.Name = m.Image.Secret.Name
	s.Image.Secret.Key = m.Image.Secret.Key

	if m.Security != nil {
		s.Security.Privileged = m.Security.Privileged
	}

	if m.Resources != nil {
		if m.Resources.Request != nil {
			if m.Resources.Request.RAM != types.EmptyString {
				s.Resources.Request.RAM, _ = resource.DecodeMemoryResource(m.Resources.Request.RAM)
			}

			if m.Resources.Request.CPU != types.EmptyString {
				s.Resources.Request.CPU, _ = resource.DecodeCpuResource(m.Resources.Request.CPU)
			}
		}

		if m.Resources.Limits != nil {
			if m.Resources.Limits.RAM != types.EmptyString {
				s.Resources.Limits.RAM, _ = resource.DecodeMemoryResource(m.Resources.Limits.RAM)
			}

			if m.Resources.Limits.CPU != types.EmptyString {
				s.Resources.Limits.CPU, _ = resource.DecodeCpuResource(m.Resources.Limits.CPU)
			}
		}
	}

	for _, v := range m.Volumes {

		s.Volumes = append(s.Volumes, &types.SpecTemplateContainerVolume{
			Name: v.Name,
			Mode: v.Mode,
			Path: v.Path,
		})
	}

	return s
}

func (m ManifestSpecRuntime) SetSpecRuntime(sr *types.SpecRuntime) {

	// check services in runtime spec
	if !compare.SliceOfString(m.Services, sr.Services) {
		sr.Services = m.Services
		sr.Updated = time.Now()
	}

	var te = true

	if len(m.Tasks) != len(sr.Tasks) {
		te = false
	}

	for _, mt := range m.Tasks {
		var f = false

		for _, st := range sr.Tasks {

			// check task name
			if mt.Name != st.Name {
				continue
			}

			// check container name
			if mt.Container != st.Container {
				continue
			}

			// check envs commands
			if len(mt.Env) != len(st.EnvVars) {
				continue
			}

			var ee = true
			for _, ce := range mt.Env {

				var f = false

				for _, se := range st.EnvVars {

					if ce.Name != se.Name {
						continue
					}

					if ce.Value != se.Value {
						continue
					}

					if se.Secret.Name != ce.Secret.Name || se.Secret.Key != ce.Secret.Key {
						continue
					}

					if se.Config.Name != ce.Config.Name || se.Secret.Key != ce.Config.Key {
						continue
					}

					f = true
				}

				if !f {
					ee = false
					break
				}
			}

			if !ee {
				continue
			}

			// check container commands
			if len(mt.Commands) != len(st.Commands) {
				continue
			}

			var ce = true
			for _, mc := range mt.Commands {
				var f = false

				for _, sc := range st.Commands {

					if mc.Workdir != sc.Workdir {
						continue
					}

					if !compare.SliceOfString(strings.Split(mc.Command, " "), sc.Command) {
						continue
					}

					if !compare.SliceOfString(strings.Split(mc.Entrypoint, " "), sc.Entrypoint) {
						continue
					}

					if !compare.SliceOfString(mc.Args, sc.Args) {
						continue
					}

					f = true
				}

				if !f {
					ce = false
				}
			}

			if !ce {
				continue
			}

			f = true
		}

		if !f {
			te = false
			break
		}
	}

	// apply new task manifest if not equal
	if !te {

		sr.Tasks = make([]types.SpecRuntimeTask, 0)

		for _, t := range m.Tasks {
			task := types.SpecRuntimeTask{
				Name:      t.Name,
				Container: t.Container,
				EnvVars:   make(types.SpecTemplateContainerEnvs, 0),
				Commands:  make([]types.SpecTemplateContainerExec, 0),
			}

			for _, e := range t.Env {
				env := types.SpecTemplateContainerEnv{
					Name:  e.Name,
					Value: e.Value,
					Config: types.SpecTemplateContainerEnvConfig{
						Name: e.Config.Name,
						Key:  e.Config.Key,
					},
					Secret: types.SpecTemplateContainerEnvSecret{
						Name: e.Secret.Name,
						Key:  e.Secret.Key,
					},
				}

				task.EnvVars = append(task.EnvVars, &env)
			}

			for _, c := range t.Commands {
				cmd := types.SpecTemplateContainerExec{
					Command:    strings.Split(c.Command, " "),
					Workdir:    c.Workdir,
					Args:       c.Args,
					Entrypoint: strings.Split(c.Entrypoint, " "),
				}
				task.Commands = append(task.Commands, cmd)
			}

			sr.Tasks = append(sr.Tasks, task)
		}

		sr.Updated = time.Now()
	}

}

func (m ManifestSpecSelector) SetSpecSelector(ss *types.SpecSelector) {

	if ss.Node != m.Node {
		ss.Node = m.Node
		ss.Updated = time.Now()
	}

	if len(ss.Labels) != len(m.Labels) {
		ss.Updated = time.Now()
	}

	var eq = true
	for k, v := range m.Labels {
		if _, ok := ss.Labels[k]; !ok {
			eq = false
			break
		}

		if ss.Labels[k] != v {
			eq = false
			break
		}
	}

	if !eq {
		ss.Labels = m.Labels
		ss.Updated = time.Now()
	}
}

func (m ManifestSpecTemplate) SetSpecTemplate(st *types.SpecTemplate) error {

	for _, c := range m.Containers {

		var (
			f    = false
			err  error
			spec *types.SpecTemplateContainer
		)

		for _, sc := range st.Containers {
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
			st.Updated = time.Now()
		}

		if spec.Image.Name != c.Image.Name {
			spec.Image.Name = c.Image.Name
			st.Updated = time.Now()
		}

		if spec.Image.Secret.Name != c.Image.Secret.Name || spec.Image.Secret.Key != c.Image.Secret.Key {
			spec.Image.Secret.Name = c.Image.Secret.Name
			spec.Image.Secret.Key = c.Image.Secret.Key
			st.Updated = time.Now()
		}

		if strings.Join(spec.Exec.Command, " ") != c.Command {
			spec.Exec.Command = strings.Split(c.Command, " ")
			st.Updated = time.Now()
		}

		if strings.Join(spec.Exec.Args, "") != strings.Join(c.Args, "") {
			spec.Exec.Args = c.Args
			st.Updated = time.Now()
		}

		if strings.Join(spec.Exec.Entrypoint, " ") != c.Entrypoint {
			spec.Exec.Entrypoint = strings.Split(c.Entrypoint, " ")
			st.Updated = time.Now()
		}

		if spec.Exec.Workdir != c.Workdir {
			spec.Exec.Workdir = c.Workdir
			st.Updated = time.Now()
		}

		if c.Security != nil && spec.Security.Privileged != c.Security.Privileged {
			spec.Security.Privileged = c.Security.Privileged
			st.Updated = time.Now()
		}

		// Environments check
		for _, ce := range c.Env {
			var f = false

			for _, se := range spec.EnvVars {
				if ce.Name == se.Name {
					f = true

					if se.Value != ce.Value {
						se.Value = ce.Value
						st.Updated = time.Now()
					}

					if ce.Secret != nil && (se.Secret.Name != ce.Secret.Name || se.Secret.Key != ce.Secret.Key) {
						se.Secret.Name = ce.Secret.Name
						se.Secret.Key = ce.Secret.Key
						st.Updated = time.Now()
					}

					if ce.Config != nil && (se.Config.Name != ce.Config.Name || se.Config.Key != ce.Config.Key) {
						se.Config.Name = ce.Config.Name
						se.Config.Key = ce.Config.Key
						st.Updated = time.Now()
					}
				}
			}

			if !f {
				item := &types.SpecTemplateContainerEnv{
					Name:  ce.Name,
					Value: ce.Value,
				}

				if ce.Secret != nil {
					item.Secret = types.SpecTemplateContainerEnvSecret{
						Name: ce.Secret.Name,
						Key:  ce.Secret.Key,
					}
				}

				if ce.Config != nil {
					item.Config = types.SpecTemplateContainerEnvConfig{
						Name: ce.Config.Name,
						Key:  ce.Config.Key,
					}
				}

				spec.EnvVars = append(spec.EnvVars, item)
				st.Updated = time.Now()
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
			st.Updated = time.Now()
		}
		spec.EnvVars = envs

		var (
			resourcesRequestRam int64
			resourcesRequestCPU int64

			resourcesLimitsRam int64
			resourcesLimitsCPU int64
		)

		// Resources check
		if c.Resources != nil && c.Resources.Request != nil {
			if c.Resources.Request.RAM != types.EmptyString {
				resourcesRequestRam, err = resource.DecodeMemoryResource(c.Resources.Request.RAM)
				if err != nil {
					return handleErr("request.ram", err)
				}
			}
			if c.Resources.Request.CPU != types.EmptyString {
				resourcesRequestCPU, err = resource.DecodeCpuResource(c.Resources.Request.CPU)
				if err != nil {
					return handleErr("request.cpu", err)
				}
			}
		}

		if c.Resources != nil && c.Resources.Limits != nil {
			if c.Resources.Limits.RAM != types.EmptyString {
				resourcesLimitsRam, err = resource.DecodeMemoryResource(c.Resources.Limits.RAM)
				if err != nil {
					return handleErr("limit.ram", err)
				}
			}
			if c.Resources.Limits.CPU != types.EmptyString {
				resourcesLimitsCPU, err = resource.DecodeCpuResource(c.Resources.Limits.CPU)
				if err != nil {
					return handleErr("limit.cpu", err)
				}
			}
		}

		if resourcesRequestRam != spec.Resources.Request.RAM ||
			resourcesRequestCPU != spec.Resources.Request.CPU {
			spec.Resources.Request.RAM = resourcesRequestRam
			spec.Resources.Request.CPU = resourcesRequestCPU
			st.Updated = time.Now()
		}

		if resourcesLimitsRam != spec.Resources.Limits.RAM ||
			resourcesLimitsCPU != spec.Resources.Limits.CPU {
			spec.Resources.Limits.RAM = resourcesLimitsRam
			spec.Resources.Limits.CPU = resourcesLimitsCPU
			st.Updated = time.Now()
		}

		// Volumes check
		for _, v := range c.Volumes {

			var f = false
			for _, sv := range spec.Volumes {

				if v.Name == sv.Name {
					f = true
					if sv.Mode != v.Mode || sv.Path != v.Path {
						sv.Mode = v.Mode
						sv.Path = v.Path
						st.Updated = time.Now()
					}

				}
			}
			if !f {
				spec.Volumes = append(spec.Volumes, &types.SpecTemplateContainerVolume{
					Name: v.Name,
					Mode: v.Mode,
					Path: v.Path,
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
			st.Updated = time.Now()
		}

		spec.Volumes = vlms

		// Ports check
		spec.Ports = make(types.SpecTemplateContainerPorts, 0)
		for _, cp := range c.Ports {
			port := new(types.SpecTemplateContainerPort)
			port.Parse(cp)
			spec.Ports = append(spec.Ports, port)
		}

		if !f {
			st.Containers = append(st.Containers, spec)
		}

	}

	var spcs = make([]*types.SpecTemplateContainer, 0)
	for _, ss := range st.Containers {
		for _, cs := range m.Containers {
			if ss.Name == cs.Name {
				spcs = append(spcs, ss)
			}
		}
	}

	if len(spcs) != len(st.Containers) {
		st.Updated = time.Now()
	}

	st.Containers = spcs

	for _, v := range m.Volumes {

		var (
			f    = false
			spec *types.SpecTemplateVolume
		)

		for _, sv := range st.Volumes {
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
			st.Updated = time.Now()
		}

		if v.Type != spec.Type || v.Volume != nil && (v.Volume.Name != spec.Volume.Name || v.Volume.Subpath != spec.Volume.Subpath) {
			spec.Type = v.Type
			spec.Volume.Name = v.Volume.Name
			spec.Volume.Subpath = v.Volume.Subpath
			st.Updated = time.Now()
		}

		if v.Type != spec.Type || (v.Secret != nil && v.Secret.Name != spec.Secret.Name) {
			spec.Type = v.Type
			spec.Secret.Name = v.Secret.Name
			st.Updated = time.Now()
		}

		var e = true
		if v.Secret != nil {
			for _, vf := range v.Secret.Binds {

				var f = false
				for _, sf := range spec.Secret.Binds {
					if vf.Key == sf.Key && vf.File == sf.File {
						f = true
						break
					}
				}

				if !f {
					e = false
					break
				}
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
			st.Updated = time.Now()
		}

		if v.Type != spec.Type || (v.Config != nil && v.Config.Name != spec.Config.Name) {
			spec.Type = v.Type
			spec.Config.Name = v.Config.Name
			st.Updated = time.Now()
		}

		var ce = true
		if v.Config != nil {
			for _, vf := range v.Config.Binds {

				var f = false
				for _, sf := range spec.Config.Binds {
					if vf.Key == sf.Key && vf.File == sf.File {
						f = true
						break
					}
				}

				if !f {
					ce = false
					break
				}
			}
		}

		if !ce {
			spec.Config.Binds = make([]types.SpecTemplateConfigVolumeBind, 0)
			for _, v := range v.Config.Binds {
				spec.Config.Binds = append(spec.Config.Binds, types.SpecTemplateConfigVolumeBind{
					Key:  v.Key,
					File: v.File,
				})
			}
			st.Updated = time.Now()
		}

		if !f {
			st.Volumes = append(st.Volumes, spec)
		}

	}

	var vlms = make([]*types.SpecTemplateVolume, 0)
	for _, ss := range st.Volumes {
		for _, cs := range m.Volumes {
			if ss.Name == cs.Name {
				vlms = append(vlms, ss)
			}
		}
	}

	if len(vlms) != len(st.Volumes) {
		st.Updated = time.Now()
	}

	st.Volumes = vlms

	return nil
}

func (m ManifestSpecRuntime) SetManifestSpecRuntime(sr *types.ManifestSpecRuntime) {

	// check services in runtime spec
	if !compare.SliceOfString(m.Services, sr.Services) {
		sr.Services = m.Services
	}

	var te = true

	if len(m.Tasks) != len(sr.Tasks) {
		te = false
	}

	for _, mt := range m.Tasks {
		var f = false

		for _, st := range sr.Tasks {

			// check task name
			if mt.Name != st.Name {
				continue
			}

			// check container name
			if mt.Container != st.Container {
				continue
			}

			// check envs commands
			if len(mt.Env) != len(st.Env) {
				continue
			}

			var ee = true
			for _, ce := range mt.Env {

				var f = false

				for _, se := range st.Env {

					if ce.Name != se.Name {
						continue
					}

					if ce.Value != se.Value {
						continue
					}

					if se.Secret.Name != ce.Secret.Name || se.Secret.Key != ce.Secret.Key {
						continue
					}

					if se.Config.Name != ce.Config.Name || se.Secret.Key != ce.Config.Key {
						continue
					}

					f = true
				}

				if !f {
					ee = false
					break
				}
			}

			if !ee {
				continue
			}

			// check container commands
			if len(mt.Commands) != len(st.Commands) {
				continue
			}

			var ce = true
			for _, mc := range mt.Commands {
				var f = false

				for _, sc := range st.Commands {

					if mc.Workdir != sc.Workdir {
						continue
					}

					if mc.Command != sc.Command {
						continue
					}

					if mc.Entrypoint != sc.Entrypoint {
						continue
					}

					if !compare.SliceOfString(mc.Args, sc.Args) {
						continue
					}

					f = true
				}

				if !f {
					ce = false
				}
			}

			if !ce {
				continue
			}

			f = true
		}

		if !f {
			te = false
			break
		}
	}

	// apply new task manifest if not equal
	if !te {

		sr.Tasks = make([]types.ManifestSpecRuntimeTask, 0)

		for _, t := range m.Tasks {
			task := types.ManifestSpecRuntimeTask{
				Name:      t.Name,
				Container: t.Container,
				Env:       make([]types.ManifestSpecTemplateContainerEnv, 0),
				Commands:  make([]types.ManifestSpecRuntimeTaskCommand, 0),
			}

			for _, e := range t.Env {
				env := types.ManifestSpecTemplateContainerEnv{
					Name:  e.Name,
					Value: e.Value,
				}

				if e.Secret != nil {
					env.Secret.Name = e.Secret.Name
					env.Secret.Key = e.Secret.Key
				}

				if e.Config != nil {
					env.Config.Name = e.Config.Name
					env.Config.Key = e.Config.Key
				}

				task.Env = append(task.Env, env)
			}

			for _, c := range t.Commands {
				cmd := types.ManifestSpecRuntimeTaskCommand{
					Command:    c.Command,
					Workdir:    c.Workdir,
					Args:       c.Args,
					Entrypoint: c.Entrypoint,
				}
				task.Commands = append(task.Commands, cmd)
			}

			sr.Tasks = append(sr.Tasks, task)
		}

	}

}

func (m ManifestSpecSelector) SetManifestSpecSelector(ss *types.ManifestSpecSelector) {

	if ss.Node != m.Node {
		ss.Node = m.Node
	}

	var eq = true
	for k, v := range m.Labels {
		if _, ok := ss.Labels[k]; !ok {
			eq = false
			break
		}

		if ss.Labels[k] != v {
			eq = false
			break
		}
	}

	if !eq {
		ss.Labels = m.Labels
	}
}

func (m ManifestSpecTemplate) SetManifestSpecTemplate(st *types.ManifestSpecTemplate) error {

	for _, c := range m.Containers {

		var (
			f    = false
			spec *types.ManifestSpecTemplateContainer
		)

		for _, sc := range st.Containers {
			if c.Name == sc.Name {
				f = true
				spec = &sc
			}
		}

		if spec == nil {
			spec = new(types.ManifestSpecTemplateContainer)
		}

		if spec.Name == types.EmptyString {
			spec.Name = c.Name
		}

		if spec.Image.Name != c.Image.Name {
			spec.Image.Name = c.Image.Name
		}

		if spec.Image.Secret.Name != c.Image.Secret.Name || spec.Image.Secret.Key != c.Image.Secret.Key {
			spec.Image.Secret.Name = c.Image.Secret.Name
			spec.Image.Secret.Key = c.Image.Secret.Key
		}

		if spec.Command != c.Command {
			spec.Command = c.Command
		}

		if strings.Join(spec.Args, "") != strings.Join(c.Args, "") {
			spec.Args = c.Args
		}

		if spec.Entrypoint != c.Entrypoint {
			spec.Entrypoint = c.Entrypoint
		}

		if spec.Workdir != c.Workdir {
			spec.Workdir = c.Workdir
		}

		if c.Security != nil && spec.Security.Privileged != c.Security.Privileged {
			spec.Security.Privileged = c.Security.Privileged
		}

		// Environments check
		for _, ce := range c.Env {
			var f = false

			for _, se := range spec.Env {
				if ce.Name == se.Name {
					f = true

					if se.Value != ce.Value {
						se.Value = ce.Value
					}

					if se.Secret.Name != ce.Secret.Name || se.Secret.Key != ce.Secret.Key {
						se.Secret.Name = ce.Secret.Name
						se.Secret.Key = ce.Secret.Key
					}

					if se.Config.Name != ce.Config.Name || se.Secret.Key != ce.Config.Key {
						se.Config.Name = ce.Config.Name
						se.Config.Key = ce.Config.Key
					}
				}
			}

			if !f {
				env := types.ManifestSpecTemplateContainerEnv{
					Name:  ce.Name,
					Value: ce.Value,
				}

				if ce.Secret != nil {
					env.Secret = types.ManifestSpecTemplateContainerEnvSecret{
						Name: ce.Secret.Name,
						Key:  ce.Secret.Key,
					}
				}

				if ce.Config != nil {
					env.Config = types.ManifestSpecTemplateContainerEnvConfig{
						Name: ce.Config.Name,
						Key:  ce.Config.Key,
					}
				}

				spec.Env = append(spec.Env, env)
			}
		}

		var envs = make([]types.ManifestSpecTemplateContainerEnv, 0)
		for _, se := range spec.Env {
			for _, ce := range c.Env {
				if ce.Name == se.Name {
					envs = append(envs, se)
					break
				}
			}
		}

		spec.Env = envs
		if c.Resources != nil {
			if c.Resources.Request != nil {
				spec.Resources.Request.RAM = c.Resources.Request.RAM
				spec.Resources.Request.CPU = c.Resources.Request.CPU
			}
			if c.Resources.Limits != nil {
				spec.Resources.Limits.RAM = c.Resources.Limits.RAM
				spec.Resources.Limits.CPU = c.Resources.Limits.CPU
			}
		}

		// Volumes check
		for _, v := range c.Volumes {

			var f = false
			for _, sv := range spec.Volumes {

				if v.Name == sv.Name {
					f = true
					if sv.Mode != v.Mode || sv.Path != v.Path {
						sv.Mode = v.Mode
						sv.Path = v.Path
					}

				}
			}
			if !f {
				spec.Volumes = append(spec.Volumes, types.ManifestSpecTemplateContainerVolume{
					Name: v.Name,
					Mode: v.Mode,
					Path: v.Path,
				})
			}
		}

		vlms := make([]types.ManifestSpecTemplateContainerVolume, 0)
		for _, sv := range spec.Volumes {
			for _, cv := range c.Volumes {
				if sv.Name == cv.Name {
					vlms = append(vlms, sv)
					break
				}
			}
		}

		spec.Volumes = vlms
		spec.Ports = c.Ports

		if !f {
			st.Containers = append(st.Containers, *spec)
		}

	}

	var spcs = make([]types.ManifestSpecTemplateContainer, 0)
	for _, ss := range st.Containers {
		for _, cs := range m.Containers {
			if ss.Name == cs.Name {
				spcs = append(spcs, ss)
			}
		}
	}

	st.Containers = spcs

	for _, v := range m.Volumes {

		var (
			f    = false
			spec *types.ManifestSpecTemplateVolume
		)

		for _, sv := range st.Volumes {
			if v.Name == sv.Name {
				f = true
				spec = &sv
			}
		}

		if spec == nil {
			spec = new(types.ManifestSpecTemplateVolume)
		}

		if spec.Name == types.EmptyString {
			spec.Name = v.Name
		}

		if v.Type != spec.Type || v.Volume != nil && (v.Volume.Name != spec.Volume.Name || v.Volume.Subpath != spec.Volume.Subpath) {
			spec.Type = v.Type
			spec.Volume.Name = v.Volume.Name
			spec.Volume.Subpath = v.Volume.Subpath
		}

		if v.Type != spec.Type || v.Secret != nil && (v.Secret.Name != spec.Secret.Name) {
			spec.Type = v.Type
			spec.Secret.Name = v.Secret.Name
		}

		var e = true

		if v.Secret != nil {
			for _, vf := range v.Secret.Binds {

				var f = false
				for _, sf := range spec.Secret.Binds {
					if vf.Key == sf.Key && vf.File == sf.File {
						f = true
						break
					}
				}

				if !f {
					e = false
					break
				}

			}
		}

		if !e {
			spec.Secret.Binds = make([]types.ManifestSpecTemplateSecretVolumeBind, 0)
			for _, v := range v.Secret.Binds {
				spec.Secret.Binds = append(spec.Secret.Binds, types.ManifestSpecTemplateSecretVolumeBind{
					Key:  v.Key,
					File: v.File,
				})
			}
		}

		if v.Type != spec.Type || v.Config != nil && (v.Config.Name != spec.Config.Name) {
			spec.Type = v.Type
			spec.Config.Name = v.Config.Name
		}

		var ce = true

		if v.Config != nil {
			for _, vf := range v.Config.Binds {

				var f = false
				for _, sf := range spec.Config.Binds {
					if vf.Key == sf.Key && vf.File == sf.File {
						f = true
						break
					}
				}

				if !f {
					ce = false
					break
				}

			}
		}

		if !ce {
			spec.Config.Binds = make([]types.ManifestSpecTemplateConfigVolumeBind, 0)
			for _, v := range v.Config.Binds {
				spec.Config.Binds = append(spec.Config.Binds, types.ManifestSpecTemplateConfigVolumeBind{
					Key:  v.Key,
					File: v.File,
				})
			}
		}

		if !f {
			st.Volumes = append(st.Volumes, *spec)
		}

	}

	var vlms = make([]types.ManifestSpecTemplateVolume, 0)
	for _, ss := range st.Volumes {
		for _, cs := range m.Volumes {
			if ss.Name == cs.Name {
				vlms = append(vlms, ss)
			}
		}
	}

	st.Volumes = vlms

	return nil
}

func handleErr(msg string, e error) error {
	log.Errorf("decode resource %s error: %s", msg, e.Error())
	return e
}
