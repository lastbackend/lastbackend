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
	"encoding/base64"
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
	Namespace *string `json:"namespace" yaml:"namespace"`
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
	Data  string `json:"data,omitempty"`
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

	cfg.Spec.Type = v.Spec.Type
	cfg.Spec.Data = make([]*types.ConfigSpecData, 0)

	for _, data := range v.Spec.Data {
		cData := new(types.ConfigSpecData)
		cData.Key = data.Key
		cData.File = data.File
		switch cfg.Spec.Type {
		case types.KindConfigText:
			cData.Data = []byte(base64.StdEncoding.EncodeToString([]byte(data.Value)))
			break
		case types.KindConfigFile:
			cData.Data = []byte(base64.StdEncoding.EncodeToString([]byte(data.Data)))
			break
		}

		cfg.Spec.Data = append(cfg.Spec.Data, cData)
	}
}


func (v *ConfigManifest) GetManifest() *types.ConfigManifest {
	cfg := new(types.ConfigManifest)
	cfg.Kind = v.Spec.Type
	cfg.Data = make(map[string][]byte, 0)
	for _, data := range v.Spec.Data {
		cfg.Data[data.Key] = []byte(base64.StdEncoding.EncodeToString([]byte(data.Data)))
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
			f.Data = base64.StdEncoding.EncodeToString(c)
		}
	}
	return nil
}

type ConfigRemoveOptions struct {
	Force bool
}
