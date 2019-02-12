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
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/util/resource"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type NamespaceManifest struct {
	Meta NamespaceManifestMeta `json:"meta,omitempty" yaml:"meta,omitempty"`
	Spec NamespaceManifestSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type NamespaceManifestMeta struct {
	RuntimeMeta `yaml:",inline"`
}

type NamespaceManifestSpec struct {
	Domain    *string                    `json:"domain"`
	Resources *NamespaceResourcesOptions `json:"resources"`
}

func (s *NamespaceManifest) FromJson(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *NamespaceManifest) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *NamespaceManifest) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, s)
}

func (s *NamespaceManifest) ToYaml() ([]byte, error) {
	return yaml.Marshal(s)
}

func (s *NamespaceManifest) SetNamespaceMeta(ns *types.Namespace) {

	if ns.Meta.Name == types.EmptyString {
		ns.Meta.Name = *s.Meta.Name
	}

	if s.Meta.Description != nil {
		ns.Meta.Description = *s.Meta.Description
	}

	if s.Meta.Labels != nil {
		ns.Meta.Labels = s.Meta.Labels
	}

}

func (s *NamespaceManifest) SetNamespaceSpec(ns *types.Namespace) error {

	ns.Spec.Domain.Internal = viper.GetString("domain.internal")

	if s.Spec.Domain != nil {
		if len(*s.Spec.Domain) == 0 {
			ns.Spec.Domain.External = viper.GetString("domain.external")
		} else {
			ns.Spec.Domain.External = *s.Spec.Domain
		}
	}

	if s.Spec.Resources != nil {

		if s.Spec.Resources.Request != nil {

			if s.Spec.Resources.Request.RAM != nil {
				ram, err := resource.DecodeMemoryResource(*s.Spec.Resources.Request.RAM)
				if err != nil {
					return err
				}

				ns.Spec.Resources.Request.RAM = ram
			}

			if s.Spec.Resources.Request.CPU != nil {

				cpu, err := resource.DecodeCpuResource(*s.Spec.Resources.Request.CPU)
				if err != nil {
					return err
				}

				ns.Spec.Resources.Request.CPU = cpu
			}

			if s.Spec.Resources.Request.Storage != nil {
				storage, err := resource.DecodeMemoryResource(*s.Spec.Resources.Request.Storage)
				if err != nil {
					return err
				}

				ns.Spec.Resources.Request.Storage = storage
			}
		}

		if s.Spec.Resources.Limits != nil {

			if s.Spec.Resources.Limits.RAM != nil {

				ram, err := resource.DecodeMemoryResource(*s.Spec.Resources.Limits.RAM)
				if err != nil {
					return err
				}

				ns.Spec.Resources.Limits.RAM = ram
			}

			if s.Spec.Resources.Limits.CPU != nil {

				cpu, err := resource.DecodeCpuResource(*s.Spec.Resources.Limits.CPU)
				if err != nil {
					return err
				}

				ns.Spec.Resources.Limits.CPU = cpu
			}

			if s.Spec.Resources.Limits.Storage != nil {
				storage, err := resource.DecodeMemoryResource(*s.Spec.Resources.Limits.Storage)
				if err != nil {
					return err
				}

				ns.Spec.Resources.Limits.Storage = storage
			}
		}

	}

	return nil

}

type NamespaceApplyManifest struct {
	Services map[string]*ServiceManifest `json:"services"`
	Configs  map[string]*ConfigManifest  `json:"configs"`
	Secrets  map[string]*SecretManifest  `json:"secrets"`
	Routes   map[string]*RouteManifest   `json:"routes"`
	Volumes  map[string]*VolumeManifest  `json:"volumes"`
	Jobs     map[string]*JobManifest     `json:"jobs"`
}

// swagger:model request_namespace_remove
type NamespaceRemoveOptions struct {
	Force bool `json:"force"`
}

type NamespaceResourcesOptions struct {
	Request *NamespaceResourceOptions `json:"request"`
	Limits  *NamespaceResourceOptions `json:"limits"`
}

// swagger:model request_namespace_quotas
type NamespaceResourceOptions struct {
	RAM     *string `json:"ram"`
	CPU     *string `json:"cpu"`
	Storage *string `json:"storage"`
}
