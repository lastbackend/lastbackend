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

type ManifestSpecSelector struct {
	Node   *string           `json:"node,omitempty" yaml:"node,omitempty"`
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

type ManifestSpecNetwork struct {
	IP    *string           `json:"ip,omitempty" yaml:"ip,omitempty"`
	Ports map[uint16]string `json:"ports,omitempty" yaml:"ports,omitempty"`
}

type ManifestSpecStrategy struct {
	Type *string `json:"type,omitempty" yaml:"type,omitempty"`
}

type ManifestSpecTemplate struct {
	Containers []ManifestSpecTemplateContainer `json:"containers,omitempty" yaml:"containers"`
	Volumes    []ManifestSpecTemplateVolume    `json:"volumes,omitempty" yaml:"volumes"`
}

type ManifestSpecTemplateContainer struct {
	Name       string                                 `json:"name,omitempty" yaml:"name,omitempty"`
	Command    string                                 `json:"command,omitempty" yaml:"command,omitempty"`
	Workdir    string                                 `json:"workdir,omitempty" yaml:"workdir,omitempty"`
	Entrypoint string                                 `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
	Args       []string                               `json:"args,omitempty" yaml:"args,omitempty"`
	Env        []ManifestSpecTemplateContainerEnv     `json:"env,omitempty" yaml:"env,omitempty"`
	Image      ManifestSpecTemplateContainerImage     `json:"image,omitempty" yaml:"image,omitempty"`
	Resources  ManifestSpecTemplateContainerResources `json:"resources,omitempty" yaml:"resources,omitempty"`
	Volumes    []ManifestSpecTemplateContainerVolume  `json:"volumes,omitempty" yaml:"volumes,omitempty"`
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
	From  ManifestSpecTemplateSecretVolume `json:"from,omitempty" yaml:"name,omitempty"`
}

type ManifestSpecTemplateSecretVolume struct {
	// Secret name to mount
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Secret file key
	Files  []string `json:"files,omitempty" yaml:"files,omitempty"`
}
