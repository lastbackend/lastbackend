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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"strings"
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
}

type ManifestSpecTemplateContainerEnv struct {
	Name  string                                 `json:"name,omitempty" yaml:"name,omitempty"`
	Value string                                 `json:"value,omitempty" yaml:"value,omitempty"`
	From  ManifestSpecTemplateContainerEnvSecret `json:"from,omitempty" yaml:"from,omitempty"`
}

type ManifestSpecTemplateContainerEnvSecret struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Key  string `json:"key,omitempty" yaml:"key,omitempty"`
}

type ManifestSpecTemplateContainerImage struct {
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Secret string `json:"secret,omitempty" yaml:"secret,omitempty"`
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
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
}

type ManifestSpecTemplateContainerResource struct {
	// CPU resource option
	CPU int64 `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	// RAM resource option
	RAM int64 `json:"ram,omitempty" yaml:"ram,omitempty"`
}

type ManifestSpecTemplateVolume struct {
	// Template volume name
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Template volume types
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// Template volume from secret type
	Secret ManifestSpecTemplateSecretVolume `json:"secret,omitempty" yaml:"secret,omitempty"`
	// Template volume from config type
	Config ManifestSpecTemplateConfigVolume `json:"config,omitempty" yaml:"config,omitempty"`
}

type ManifestSpecTemplateSecretVolume struct {
	// Secret name to mount
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Secret file key
	Files []string `json:"files,omitempty" yaml:"files,omitempty"`
}

type ManifestSpecTemplateConfigVolume struct {
	// Secret name to mount
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Secret file key
	Files []string `json:"files,omitempty" yaml:"files,omitempty"`
}

type ManifestSpecTemplateRestartPolicy struct {
	Policy  string `json:"policy" yaml:"policy"`
	Attempt int    `json:"attempt" yaml:"attempt"`
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

func (m ManifestSpecTemplateVolume) GetSpec() types.SpecTemplateVolume {
	s := types.SpecTemplateVolume{
		Name: m.Name,
		Type: m.Type,
		Secret: types.SpecTemplateSecretVolume{
			Name:  m.Secret.Name,
			Files: m.Secret.Files,
		},
		Config: types.SpecTemplateConfigVolume{
			Name:  m.Config.Name,
			Files: m.Config.Files,
		},
	}
	return s
}

func (m ManifestSpecTemplateContainer) GetSpec() types.SpecTemplateContainer {
	s := types.SpecTemplateContainer{}
	s.Name = m.Name

	s.RestartPolicy.Policy = m.RestartPolicy.Policy
	s.RestartPolicy.Attempt = m.RestartPolicy.Attempt

	s.Exec.Command = strings.Split(m.Command, " ")
	s.Exec.Args = m.Args
	s.Exec.Workdir = m.Workdir
	s.Exec.Entrypoint = strings.Split(m.Entrypoint, " ")

	for _, p := range m.Ports {
		port := new(types.SpecTemplateContainerPort)
		port.Parse(p)
		s.Ports = append(s.Ports, port)
	}

	for _, e := range m.Env {
		s.EnvVars = append(s.EnvVars, &types.SpecTemplateContainerEnv{
			Name:  e.Name,
			Value: e.Value,
			Secret: types.SpecTemplateContainerEnvSecret{
				Name: e.From.Name,
				Key:  e.From.Key,
			},
			Config: types.SpecTemplateContainerEnvConfig{
				Name: e.From.Name,
				Key:  e.From.Key,
			},
		})
	}

	s.Image.Name = m.Image.Name
	s.Image.Secret = m.Image.Secret

	s.Resources.Request.RAM = m.Resources.Request.RAM
	s.Resources.Request.CPU = m.Resources.Request.CPU

	s.Resources.Limits.RAM = m.Resources.Limits.RAM
	s.Resources.Limits.CPU = m.Resources.Limits.CPU

	for _, v := range m.Volumes {

		s.Volumes = append(s.Volumes, &types.SpecTemplateContainerVolume{
			Name: v.Name,
			Mode: v.Mode,
			Path: v.Path,
		})
	}

	return s
}
