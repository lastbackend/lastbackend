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

package types

import (
	"github.com/lastbackend/lastbackend/pkg/util/compare"
	"github.com/lastbackend/lastbackend/pkg/util/resource"
	"strings"
	"time"
)

type NodeManifest struct {
	Meta      NodeManifestMeta             `json:"meta"`
	Resolvers map[string]*ResolverManifest `json:"resolvers"`
	Exporter  *ExporterManifest            `json:"exporter"`
	Secrets   map[string]*SecretManifest   `json:"secrets"`
	Configs   map[string]*ConfigManifest   `json:"configs"`
	Endpoints map[string]*EndpointManifest `json:"endpoint"`
	Network   map[string]*SubnetManifest   `json:"network"`
	Pods      map[string]*PodManifest      `json:"pods"`
	Volumes   map[string]*VolumeManifest   `json:"volumes"`
}

type NodeManifestMeta struct {
	Initial bool `json:"initial"`
}

type ResolverManifest struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

type ExporterManifest struct {
	Endpoint string `json:"endpoint,omitempty"`
}

type IngressManifest struct {
	Meta      IngressManifestMeta          `json:"meta"`
	Resolvers map[string]*ResolverManifest `json:"resolvers"`
	Routes    map[string]*RouteManifest    `json:"routes"`
	Endpoints map[string]*EndpointManifest `json:"endpoints"`
	Network   map[string]*SubnetManifest   `json:"network"`
}

type IngressManifestMeta struct {
	Initial bool `json:"initial"`
}

type DiscoveryManifest struct {
	Meta    DiscoveryManifestMeta      `json:"meta"`
	Network map[string]*SubnetManifest `json:"network"`
}

type DiscoveryManifestMeta struct {
	Initial bool `json:"initial"`
}

type PodManifest PodSpec

type PodManifestList struct {
	System
	Items []*PodManifest
}

type PodManifestMap struct {
	System
	Items map[string]*PodManifest
}

type VolumeManifest VolumeSpec

type VolumeManifestList struct {
	System
	Items []*VolumeManifest
}

type VolumeManifestMap struct {
	System
	Items map[string]*VolumeManifest
}

type SubnetManifest struct {
	System
	SubnetSpec
}

type SubnetManifestList struct {
	System
	Items []*SubnetManifest
}

type SubnetManifestMap struct {
	System
	Items map[string]*SubnetManifest
}

type EndpointManifest struct {
	System
	EndpointSpec `json:",inline"`
	Upstreams    []string `json:"upstreams"`
}

type EndpointManifestList struct {
	System
	Items []*EndpointManifest
}

type EndpointManifestMap struct {
	System
	Items map[string]*EndpointManifest
}

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
	Services []string                  `json:"services"`
	Tasks    []ManifestSpecRuntimeTask `json:"tasks"`
}

type ManifestSpecRuntimeTask struct {
	Name      string                             `json:"name"`
	Container string                             `json:"container" yaml:"container"`
	Env       []ManifestSpecTemplateContainerEnv `json:"env,omitempty" yaml:"env,omitempty"`
	Commands  []string                           `json:"commands" yaml:"commands"`
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
	Name          string                                 `json:"name,omitempty" yaml:"name,omitempty"`
	Command       string                                 `json:"command,omitempty" yaml:"command,omitempty"`
	Workdir       string                                 `json:"workdir,omitempty" yaml:"workdir,omitempty"`
	Entrypoint    string                                 `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
	Args          []string                               `json:"args,omitempty" yaml:"args,omitempty"`
	Ports         []string                               `json:"ports,omitempty" yaml:"ports,omitempty"`
	Env           []ManifestSpecTemplateContainerEnv     `json:"env,omitempty" yaml:"env,omitempty"`
	Volumes       []ManifestSpecTemplateContainerVolume  `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Image         ManifestSpecTemplateContainerImage     `json:"image,omitempty" yaml:"image,omitempty"`
	Resources     ManifestSpecTemplateContainerResources `json:"resources,omitempty" yaml:"resources,omitempty"`
	RestartPolicy ManifestSpecTemplateRestartPolicy      `json:"restart,omitempty" yaml:"restart,omitempty"`
	Security      ManifestSpecSecurity                   `json:"security,omitempty" yaml:"security,omitempty"`
}

type ManifestSpecTemplateContainerEnv struct {
	Name   string                                 `json:"name,omitempty" yaml:"name,omitempty"`
	Value  string                                 `json:"value,omitempty" yaml:"value,omitempty"`
	Secret ManifestSpecTemplateContainerEnvSecret `json:"secret,omitempty" yaml:"secret,omitempty"`
	Config ManifestSpecTemplateContainerEnvConfig `json:"config,omitempty" yaml:"config,omitempty"`
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
	Limits ManifestSpecTemplateContainerResource `json:"limits,omitempty" yaml:"limits,omitempty"`
	// Request resources
	Request ManifestSpecTemplateContainerResource `json:"quota,omitempty" yaml:"quota,omitempty"`
}

type ManifestSpecTemplateContainerVolume struct {
	// Volume name
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Volume mount mode
	Mode string `json:"mode,omitempty" yaml:"mode,omitempty"`
	// Volume mount path
	MountPath string `json:"path,omitempty" yaml:"path,omitempty"`
	// Volume mount sub path
	SubPath string `json:"sub_path,omitempty" yaml:"sub_path,omitempty"`
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
	Volume ManifestSpecTemplateVolumeClaim `json:"volume,omitempty" yaml:"volume,omitempty"`
	// Template volume from secret type
	Secret ManifestSpecTemplateSecretVolume `json:"secret,omitempty" yaml:"secret,omitempty"`
	// Template volume from config type
	Config ManifestSpecTemplateConfigVolume `json:"config,omitempty" yaml:"config,omitempty"`
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

func (m ManifestSpecSelector) GetSpec() SpecSelector {
	s := SpecSelector{}

	s.Node = m.Node
	s.Labels = m.Labels

	return s
}

func (m ManifestSpecTemplate) GetSpec() SpecTemplate {
	var s = SpecTemplate{}

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

func (m ManifestSpecRuntime) GetSpec() SpecRuntime {
	var s = SpecRuntime{}

	s.Services = m.Services

	for _, t := range m.Tasks {
		ts := t.GetSpec()
		s.Tasks = append(s.Tasks, ts)
	}

	return s
}

func (m ManifestSpecRuntimeTask) GetSpec() SpecRuntimeTask {
	s := SpecRuntimeTask{}

	s.Name = m.Name
	s.Container = m.Container

	for _, e := range m.Env {
		s.EnvVars = append(s.EnvVars, &SpecTemplateContainerEnv{
			Name:  e.Name,
			Value: e.Value,
			Secret: SpecTemplateContainerEnvSecret{
				Name: e.Secret.Name,
				Key:  e.Secret.Key,
			},
			Config: SpecTemplateContainerEnvConfig{
				Name: e.Config.Name,
				Key:  e.Config.Key,
			},
		})
	}

	s.Commands = make([]string, 0)
	if m.Commands != nil {
		s.Commands = m.Commands
	}

	return s
}

func (m ManifestSpecTemplateVolume) GetSpec() SpecTemplateVolume {
	s := SpecTemplateVolume{
		Name: m.Name,
		Type: m.Type,
		Volume: SpecTemplateVolumeClaim{
			Name:    m.Volume.Name,
			Subpath: m.Volume.Subpath,
		},
		Secret: SpecTemplateSecretVolume{
			Name:  m.Secret.Name,
			Binds: make([]SpecTemplateSecretVolumeBind, 0),
		},
		Config: SpecTemplateConfigVolume{
			Name:  m.Config.Name,
			Binds: make([]SpecTemplateConfigVolumeBind, 0),
		},
	}

	for _, b := range m.Secret.Binds {
		s.Secret.Binds = append(s.Secret.Binds, SpecTemplateSecretVolumeBind{
			Key:  b.Key,
			File: b.File,
		})
	}

	for _, b := range m.Config.Binds {
		s.Config.Binds = append(s.Config.Binds, SpecTemplateConfigVolumeBind{
			Key:  b.Key,
			File: b.File,
		})
	}

	return s
}

func (m ManifestSpecTemplateContainer) GetSpec() SpecTemplateContainer {
	s := SpecTemplateContainer{}
	s.Name = m.Name

	s.RestartPolicy.Policy = m.RestartPolicy.Policy
	s.RestartPolicy.Attempt = m.RestartPolicy.Attempt

	if m.Command != EmptyString {
		s.Exec.Command = strings.Split(m.Command, " ")
	}
	s.Exec.Args = m.Args
	s.Exec.Workdir = m.Workdir

	if m.Entrypoint != EmptyString {
		s.Exec.Entrypoint = strings.Split(m.Entrypoint, " ")
	}

	for _, p := range m.Ports {
		port := new(SpecTemplateContainerPort)
		port.Parse(p)
		s.Ports = append(s.Ports, port)
	}

	for _, e := range m.Env {
		s.EnvVars = append(s.EnvVars, &SpecTemplateContainerEnv{
			Name:  e.Name,
			Value: e.Value,
			Secret: SpecTemplateContainerEnvSecret{
				Name: e.Secret.Name,
				Key:  e.Secret.Key,
			},
			Config: SpecTemplateContainerEnvConfig{
				Name: e.Config.Name,
				Key:  e.Config.Key,
			},
		})
	}

	s.Image.Name = m.Image.Name
	s.Image.Secret.Name = m.Image.Secret.Name
	s.Image.Secret.Key = m.Image.Secret.Key

	s.Security.Privileged = m.Security.Privileged

	if m.Resources.Request.RAM != EmptyString {
		s.Resources.Request.RAM, _ = resource.DecodeMemoryResource(m.Resources.Request.RAM)
	}

	if m.Resources.Request.CPU != EmptyString {
		s.Resources.Request.CPU, _ = resource.DecodeCpuResource(m.Resources.Request.CPU)
	}

	if m.Resources.Limits.RAM != EmptyString {
		s.Resources.Limits.RAM, _ = resource.DecodeMemoryResource(m.Resources.Limits.RAM)
	}

	if m.Resources.Limits.CPU != EmptyString {
		s.Resources.Limits.CPU, _ = resource.DecodeCpuResource(m.Resources.Limits.CPU)
	}

	for _, v := range m.Volumes {

		s.Volumes = append(s.Volumes, &SpecTemplateContainerVolume{
			Name:      v.Name,
			Mode:      v.Mode,
			MountPath: v.MountPath,
			SubPath:   v.SubPath,
		})
	}

	return s
}

func (m ManifestSpecRuntime) SetSpecRuntime(sr *SpecRuntime) {

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

					if mc != sc {
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

		sr.Tasks = make([]SpecRuntimeTask, 0)

		for _, t := range m.Tasks {
			task := SpecRuntimeTask{
				Name:      t.Name,
				Container: t.Container,
				EnvVars:   make(SpecTemplateContainerEnvs, 0),
				Commands:  make([]string, 0),
			}

			for _, e := range t.Env {
				env := SpecTemplateContainerEnv{
					Name:  e.Name,
					Value: e.Value,
					Config: SpecTemplateContainerEnvConfig{
						Name: e.Config.Name,
						Key:  e.Config.Key,
					},
					Secret: SpecTemplateContainerEnvSecret{
						Name: e.Secret.Name,
						Key:  e.Secret.Key,
					},
				}

				task.EnvVars = append(task.EnvVars, &env)
			}

			task.Commands = make([]string, 0)
			if t.Commands != nil {
				task.Commands = t.Commands
			}

			sr.Tasks = append(sr.Tasks, task)
		}

		sr.Updated = time.Now()
	}

}

func (m ManifestSpecSelector) SetSpecSelector(ss *SpecSelector) {

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

func (m ManifestSpecTemplate) SetSpecTemplate(st *SpecTemplate) error {

	for _, c := range m.Containers {

		var (
			f    = false
			err  error
			spec *SpecTemplateContainer
		)

		for _, sc := range st.Containers {
			if c.Name == sc.Name {
				f = true
				spec = sc
			}
		}

		if spec == nil {
			spec = new(SpecTemplateContainer)
		}

		if spec.Name == EmptyString {
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

		if spec.Security.Privileged != c.Security.Privileged {
			spec.Security.Privileged = c.Security.Privileged
			st.Updated = time.Now()
		}

		if spec.RestartPolicy.Policy != c.RestartPolicy.Policy || spec.RestartPolicy.Attempt != c.RestartPolicy.Attempt {
			spec.RestartPolicy.Policy = c.RestartPolicy.Policy
			spec.RestartPolicy.Attempt = c.RestartPolicy.Attempt
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

					if se.Secret.Name != ce.Secret.Name || se.Secret.Key != ce.Secret.Key {
						se.Secret.Name = ce.Secret.Name
						se.Secret.Key = ce.Secret.Key
						st.Updated = time.Now()
					}

					if se.Config.Name != ce.Config.Name || se.Secret.Key != ce.Config.Key {
						se.Config.Name = ce.Config.Name
						se.Config.Key = ce.Config.Key
						st.Updated = time.Now()
					}
				}
			}

			if !f {
				spec.EnvVars = append(spec.EnvVars, &SpecTemplateContainerEnv{
					Name:  ce.Name,
					Value: ce.Value,
					Secret: SpecTemplateContainerEnvSecret{
						Name: ce.Secret.Name,
						Key:  ce.Secret.Key,
					},
					Config: SpecTemplateContainerEnvConfig{
						Name: ce.Config.Name,
						Key:  ce.Config.Key,
					},
				})
				st.Updated = time.Now()
			}
		}

		var envs = make([]*SpecTemplateContainerEnv, 0)
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
		if c.Resources.Request.RAM != EmptyString {
			resourcesRequestRam, err = resource.DecodeMemoryResource(c.Resources.Request.RAM)
			if err != nil {
				return handleErr("request.ram", err)
			}
		}
		if c.Resources.Request.CPU != EmptyString {
			resourcesRequestCPU, err = resource.DecodeCpuResource(c.Resources.Request.CPU)
			if err != nil {
				return handleErr("request.cpu", err)
			}
		}

		if c.Resources.Limits.RAM != EmptyString {
			resourcesLimitsRam, err = resource.DecodeMemoryResource(c.Resources.Limits.RAM)
			if err != nil {
				return handleErr("limit.ram", err)
			}
		}
		if c.Resources.Limits.CPU != EmptyString {
			resourcesLimitsCPU, err = resource.DecodeCpuResource(c.Resources.Limits.CPU)
			if err != nil {
				return handleErr("limit.cpu", err)
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
					if sv.Mode != v.Mode || sv.MountPath != v.MountPath || sv.SubPath != v.SubPath {
						sv.Mode = v.Mode
						sv.MountPath = v.MountPath
						sv.SubPath = v.SubPath
						st.Updated = time.Now()
					}

				}
			}
			if !f {
				spec.Volumes = append(spec.Volumes, &SpecTemplateContainerVolume{
					Name:      v.Name,
					Mode:      v.Mode,
					MountPath: v.MountPath,
					SubPath:   v.SubPath,
				})
			}
		}

		vlms := make([]*SpecTemplateContainerVolume, 0)
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
		spec.Ports = make(SpecTemplateContainerPorts, 0)
		for _, cp := range c.Ports {
			port := new(SpecTemplateContainerPort)
			port.Parse(cp)
			spec.Ports = append(spec.Ports, port)
		}

		if !f {
			st.Containers = append(st.Containers, spec)
		}

	}

	var spcs = make([]*SpecTemplateContainer, 0)
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
			spec *SpecTemplateVolume
		)

		for _, sv := range st.Volumes {
			if v.Name == sv.Name {
				f = true
				spec = sv
			}
		}

		if spec == nil {
			spec = new(SpecTemplateVolume)
		}

		if spec.Name == EmptyString {
			spec.Name = v.Name
			st.Updated = time.Now()
		}

		if v.Type != spec.Type || v.Volume.Name != spec.Volume.Name || v.Volume.Subpath != spec.Volume.Subpath {
			spec.Type = v.Type
			spec.Volume.Name = v.Volume.Name
			spec.Volume.Subpath = v.Volume.Subpath
			st.Updated = time.Now()
		}

		if v.Type != spec.Type || v.Secret.Name != spec.Secret.Name {
			spec.Type = v.Type
			spec.Secret.Name = v.Secret.Name
			st.Updated = time.Now()
		}

		var e = true
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

		if !e {
			spec.Secret.Binds = make([]SpecTemplateSecretVolumeBind, 0)
			for _, v := range v.Secret.Binds {
				spec.Secret.Binds = append(spec.Secret.Binds, SpecTemplateSecretVolumeBind{
					Key:  v.Key,
					File: v.File,
				})
			}
			st.Updated = time.Now()
		}

		if v.Type != spec.Type || v.Config.Name != spec.Config.Name {
			spec.Type = v.Type
			spec.Config.Name = v.Config.Name
			st.Updated = time.Now()
		}

		var ce = true
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

		if !ce {
			spec.Config.Binds = make([]SpecTemplateConfigVolumeBind, 0)
			for _, v := range v.Config.Binds {
				spec.Config.Binds = append(spec.Config.Binds, SpecTemplateConfigVolumeBind{
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

	var vlms = make([]*SpecTemplateVolume, 0)
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

func NewPodManifestMap() *PodManifestMap {
	dm := new(PodManifestMap)
	dm.Items = make(map[string]*PodManifest)
	return dm
}

func NewVolumeManifestMap() *VolumeManifestMap {
	dm := new(VolumeManifestMap)
	dm.Items = make(map[string]*VolumeManifest)
	return dm
}

func NewSubnetManifestMap() *SubnetManifestMap {
	dm := new(SubnetManifestMap)
	dm.Items = make(map[string]*SubnetManifest)
	return dm
}

func NewEndpointManifestMap() *EndpointManifestMap {
	dm := new(EndpointManifestMap)
	dm.Items = make(map[string]*EndpointManifest)
	return dm
}

func handleErr(msg string, e error) error {
	return e
}
