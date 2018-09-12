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
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ConfigManifest struct {
	Meta ConfigManifestMeta `json:"meta,omitempty" yaml:"meta,omitempty"`
	Spec ConfigManifestSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type ConfigManifestMeta struct {
	RuntimeMeta `yaml:",inline"`
}

type ConfigManifestSpec struct {
	// Template volume types
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// Tempate volume selector
	Data []*ConfigManifestData `json:"data,omitempty" yaml:"data,omitempty"`
}

type ConfigManifestData struct {
	Key   string `json:"key,omitempty" yaml:"key,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
	File  string `json:"file,omitempty" yaml:"file,omitempty"`
	Data  []byte `json:"data,omitempty"`
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

func (v *ConfigManifest) SetConfigMeta(vol *types.Config) {

	if vol.Meta.Name == types.EmptyString {
		vol.Meta.Name = *v.Meta.Name
	}

	if v.Meta.Description != nil {
		vol.Meta.Description = *v.Meta.Description
	}

	if v.Meta.Labels != nil {
		vol.Meta.Labels = v.Meta.Labels
	}

}

func (v *ConfigManifest) SetConfigSpec(vol *types.Config) {

}

func (v *ConfigManifest) GetManifest() *types.ConfigManifest {
	cfg := new(types.ConfigManifest)
	cfg.Kind = v.Spec.Type
	cfg.Data = make(map[string][]byte, 0)
	for _, data := range v.Spec.Data {
		cfg.Data[data.Key] = data.Data
	}
	return cfg
}

func (v *ConfigManifest) ReadData() error {
	for _, f := range v.Spec.Data {
		if f.File != types.EmptyString {
			c, err := ioutil.ReadFile(f.File)
			if err != nil {
				_ = fmt.Errorf("failed read data from file: %s", f)
				return err
			}
			f.Data = c
		}
	}
	return nil
}

type ConfigRemoveOptions struct {
	Force bool
}
