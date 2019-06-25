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
	"gopkg.in/yaml.v2"
)

type ConfigManifest struct {
	Meta ConfigManifestMeta `json:"meta,omitempty" yaml:"meta,omitempty"`
	Spec ConfigManifestSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type ConfigManifestMeta struct {
	RuntimeMeta `yaml:",inline"`
	Namespace   *string `json:"namespace" yaml:"namespace"`
}

type ConfigManifestSpec struct {
	// Config data
	Data map[string]string `json:"data,omitempty" yaml:"data,omitempty"`
}

func (v *ConfigManifest) FromJson(data []byte) error {
	return json.Unmarshal(data, v)
}

func (v *ConfigManifest) ToJson() ([]byte, error) {
	return json.Marshal(v)
}

func (v *ConfigManifest) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, v)
}

func (v *ConfigManifest) ToYaml() ([]byte, error) {
	return yaml.Marshal(v)
}

func (v *ConfigManifest) SetConfigMeta(cfg *types.Config) {

	if cfg.Meta.Name == types.EmptyString {
		cfg.Meta.Name = *v.Meta.Name
	}

	if v.Meta.Description != nil {
		cfg.Meta.Description = *v.Meta.Description
	}

	if v.Meta.Labels != nil {
		cfg.Meta.Labels = v.Meta.Labels
	}

}

// SetConfigSpec - set config spec from manifest
// TODO: check if config spec is updated => update Meta.Updated or skip
func (v *ConfigManifest) SetConfigSpec(cfg *types.Config) {

	cfg.Spec.Data = make(map[string]string, 0)

	for key, value := range v.Spec.Data {
		cfg.Spec.Data[key] = value
	}
}

func (v *ConfigManifest) GetManifest() *types.ConfigManifest {
	cfg := new(types.ConfigManifest)
	cfg.Data = make(map[string]string, 0)
	for key, value := range v.Spec.Data {
		cfg.Data[key] = value
	}
	return cfg
}

type ConfigRemoveOptions struct {
	Force bool
}
