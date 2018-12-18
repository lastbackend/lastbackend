//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/util/resource"
	"gopkg.in/yaml.v2"
	"strconv"
	"strings"
	"time"
)

type ServiceManifest struct {
	Meta ServiceManifestMeta `json:"meta,omitempty" yaml:"meta,omitempty"`
	Spec ServiceManifestSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type ServiceManifestMeta struct {
	RuntimeMeta `yaml:",inline"`
}

type ServiceManifestSpec struct {
	Replicas *int                  `json:"replicas,omitempty" yaml:"replicas,omitempty"`
	Selector *ManifestSpecSelector `json:"selector,omitempty" yaml:"selector,omitempty"`
	Network  *ManifestSpecNetwork  `json:"network,omitempty" yaml:"network,omitempty"`
	Strategy *ManifestSpecStrategy `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	Template *ManifestSpecTemplate `json:"template,omitempty" yaml:"template,omitempty"`
}

func (s *ServiceManifest) FromJson(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *ServiceManifest) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *ServiceManifest) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, s)
}

func (s *ServiceManifest) ToYaml() ([]byte, error) {
	return yaml.Marshal(s)
}

func (s *ServiceManifest) SetServiceMeta(svc *types.Service) {

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

func (s *ServiceManifest) SetServiceSpec(svc *types.Service) {

	tn := svc.Spec.Network.Updated
	tc := svc.Spec.Template.Updated

	defer func() {
		if s.Spec.Replicas != nil {
			svc.Status.State = types.StateProvision
			return
		}

		if tn.Before(svc.Spec.Network.Updated) || tc.Before(svc.Spec.Template.Updated) {
			svc.Status.State = types.StateProvision
			return
		}
	}()

	if s.Spec.Replicas != nil {
		svc.Spec.Replicas = *s.Spec.Replicas
	}

	if s.Spec.Network != nil {

		if s.Spec.Network.IP != nil {
			svc.Spec.Network.IP = *s.Spec.Network.IP
		}

		if len(s.Spec.Network.Ports) > 0 {

			svc.Spec.Network.Ports = make(map[uint16]string, 0)

			for _, p := range s.Spec.Network.Ports {
				mp := strings.Split(p, ":")
				var base = 10
				var size = 16
				port, err := strconv.ParseUint(mp[0], base, size)
				if err != nil {
					continue
				}
				if len(mp) == 1 {
					svc.Spec.Network.Ports[uint16(port)] = mp[0]
				}

				if len(mp) == 2 {
					svc.Spec.Network.Ports[uint16(port)] = mp[1]
				}

			}
		}

		svc.Spec.Network.Updated = time.Now()
	}

	if s.Spec.Selector != nil {

		if svc.Spec.Selector.Node != s.Spec.Selector.Node {
			svc.Spec.Selector.Node = s.Spec.Selector.Node
			svc.Spec.Selector.Updated = time.Now()
		}

		if len(svc.Spec.Selector.Labels) != len(s.Spec.Selector.Labels) {
			svc.Spec.Selector.Updated = time.Now()
		}

		var eq = true
		for k, v := range s.Spec.Selector.Labels {
			if _, ok := svc.Spec.Selector.Labels[k]; !ok {
				eq = false
				break
			}

			if svc.Spec.Selector.Labels[k] != v {
				eq = false
				break
			}
		}

		if !eq {
			svc.Spec.Selector.Labels = s.Spec.Selector.Labels
			svc.Spec.Selector.Updated = time.Now()
		}

	} else {

		if svc.Spec.Selector.Node != types.EmptyString {
			svc.Spec.Selector.Node = types.EmptyString
			svc.Spec.Selector.Updated = time.Now()
		}

		if len(svc.Spec.Selector.Labels) > 0 {
			svc.Spec.Selector.Labels = make(map[string]string)
			svc.Spec.Selector.Updated = time.Now()
		}
	}

	if s.Spec.Strategy != nil {
		if s.Spec.Strategy.Type != nil {
			svc.Spec.Strategy.Type = *s.Spec.Strategy.Type
		}
	}

	if s.Spec.Template != nil {

		for _, c := range s.Spec.Template.Containers {

			var (
				f    = false
				spec *types.SpecTemplateContainer
			)

			for _, sc := range svc.Spec.Template.Containers {
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
				svc.Spec.Template.Updated = time.Now()
			}

			if spec.Image.Name != c.Image.Name {
				spec.Image.Name = c.Image.Name
				svc.Spec.Template.Updated = time.Now()
			}

			if spec.Image.Secret != c.Image.Secret {
				spec.Image.Secret = c.Image.Secret
				svc.Spec.Template.Updated = time.Now()
			}

			if strings.Join(spec.Exec.Command, " ") != c.Command {
				spec.Exec.Command = strings.Split(c.Command, " ")
				svc.Spec.Template.Updated = time.Now()
			}

			if strings.Join(spec.Exec.Args, "") != strings.Join(c.Args, "") {
				spec.Exec.Args = c.Args
				svc.Spec.Template.Updated = time.Now()
			}

			if strings.Join(spec.Exec.Entrypoint, " ") != c.Entrypoint {
				spec.Exec.Entrypoint = strings.Split(c.Entrypoint, " ")
				svc.Spec.Template.Updated = time.Now()
			}

			if spec.Exec.Workdir != c.Workdir {
				spec.Exec.Workdir = c.Workdir
				svc.Spec.Template.Updated = time.Now()
			}


			// Environments check
			for _, ce := range c.Env {
				var f = false

				for _, se := range spec.EnvVars {
					if ce.Name == se.Name {
						f = true

						if se.Value != ce.Value {
							se.Value = ce.Value
							svc.Spec.Template.Updated = time.Now()
						}

						if se.Secret.Name != ce.Secret.Name || se.Secret.Key != ce.Secret.Key {
							se.Secret.Name = ce.Secret.Name
							se.Secret.Key = ce.Secret.Key
							svc.Spec.Template.Updated = time.Now()
						}

						if se.Config.Name != ce.Config.Name || se.Secret.Key != ce.Config.Key {
							se.Config.Name = ce.Config.Name
							se.Config.Key = ce.Config.Key
							svc.Spec.Template.Updated = time.Now()
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
					svc.Spec.Template.Updated = time.Now()
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
				svc.Spec.Template.Updated = time.Now()
			}
			spec.EnvVars = envs

			// Resources check
			resourcesRequestRam, _ := resource.DecodeMemoryResource(c.Resources.Request.RAM)
			resourcesRequestCPU, _ := resource.DecodeCpuResource(c.Resources.Request.CPU)

			resourcesLimitsRam, _ := resource.DecodeMemoryResource(c.Resources.Limits.RAM)
			resourcesLimitsCPU, _ := resource.DecodeCpuResource(c.Resources.Limits.CPU)


			if resourcesRequestRam != spec.Resources.Request.RAM ||
				resourcesRequestCPU != spec.Resources.Request.CPU {
				spec.Resources.Request.RAM = resourcesRequestRam
				spec.Resources.Request.CPU = resourcesRequestCPU
				svc.Spec.Template.Updated = time.Now()
			}

			if resourcesLimitsRam != spec.Resources.Limits.RAM ||
				resourcesLimitsCPU != spec.Resources.Limits.CPU {
				spec.Resources.Limits.RAM = resourcesLimitsRam
				spec.Resources.Limits.CPU = resourcesLimitsCPU
				svc.Spec.Template.Updated = time.Now()
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
							svc.Spec.Template.Updated = time.Now()
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
				svc.Spec.Template.Updated = time.Now()
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
				svc.Spec.Template.Containers = append(svc.Spec.Template.Containers, spec)
			}

		}

		var spcs = make([]*types.SpecTemplateContainer, 0)
		for _, ss := range svc.Spec.Template.Containers {
			for _, cs := range s.Spec.Template.Containers {
				if ss.Name == cs.Name {
					spcs = append(spcs, ss)
				}
			}
		}

		if len(spcs) != len(svc.Spec.Template.Containers) {
			svc.Spec.Template.Updated = time.Now()
		}

		svc.Spec.Template.Containers = spcs

		for _, v := range s.Spec.Template.Volumes {

			var (
				f    = false
				spec *types.SpecTemplateVolume
			)

			for _, sv := range svc.Spec.Template.Volumes {
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
				svc.Spec.Template.Updated = time.Now()
			}

			if v.Type != spec.Type || v.Volume.Name != spec.Volume.Name || v.Volume.Subpath != spec.Volume.Subpath {
				spec.Type = v.Type
				spec.Volume.Name = v.Volume.Name
				spec.Volume.Subpath = v.Volume.Subpath
				svc.Spec.Template.Updated = time.Now()
			}

			if v.Type != spec.Type || v.Secret.Name != spec.Secret.Name {
				spec.Type = v.Type
				spec.Secret.Name = v.Secret.Name
				svc.Spec.Template.Updated = time.Now()
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
				spec.Secret.Binds = make([]types.SpecTemplateSecretVolumeBind, 0)
				for _, v := range v.Secret.Binds {
					spec.Secret.Binds = append(spec.Secret.Binds, types.SpecTemplateSecretVolumeBind{
						Key:  v.Key,
						File: v.File,
					})
				}
				svc.Spec.Template.Updated = time.Now()
			}

			if v.Type != spec.Type || v.Config.Name != spec.Config.Name {
				spec.Type = v.Type
				spec.Config.Name = v.Config.Name
				svc.Spec.Template.Updated = time.Now()
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
				spec.Config.Binds = make([]types.SpecTemplateConfigVolumeBind, 0)
				for _, v := range v.Config.Binds {
					spec.Config.Binds = append(spec.Config.Binds, types.SpecTemplateConfigVolumeBind{
						Key:  v.Key,
						File: v.File,
					})
				}
				svc.Spec.Template.Updated = time.Now()
			}

			if !f {
				svc.Spec.Template.Volumes = append(svc.Spec.Template.Volumes, spec)
			}

		}

		var vlms = make([]*types.SpecTemplateVolume, 0)
		for _, ss := range svc.Spec.Template.Volumes {
			for _, cs := range s.Spec.Template.Volumes {
				if ss.Name == cs.Name {
					vlms = append(vlms, ss)
				}
			}
		}

		if len(vlms) != len(svc.Spec.Template.Volumes) {
			svc.Spec.Template.Updated = time.Now()
		}

		svc.Spec.Template.Volumes = vlms

	}

}

func (s *ServiceManifest) GetManifest() *types.ServiceManifest {
	sm := new(types.ServiceManifest)
	return sm
}

// swagger:ignore
// swagger:model request_service_remove
type ServiceRemoveOptions struct {
	Force bool `json:"force"`
}

// swagger:ignore
// swagger:model request_service_logs
type ServiceLogsOptions struct {
	Deployment string `json:"deployment"`
	Pod        string `json:"pod"`
	Container  string `json:"container"`
	Follow     bool   `json:"follow"`
}
