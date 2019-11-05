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
	"encoding/base64"
	"encoding/json"
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"gopkg.in/yaml.v2"
)

type SecretManifest struct {
	Meta SecretManifestMeta `json:"meta,omitempty" yaml:"meta,omitempty"`
	Spec SecretManifestSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type SecretManifestMeta struct {
	RuntimeMeta `yaml:",inline"`
	Namespace   *string `json:"namespace" yaml:"namespace"`
}

type SecretManifestSpec struct {
	// Template volume types
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// Tempate volume selector
	Data map[string]string `json:"data,omitempty" yaml:"data,omitempty"`
}

func (v *SecretManifest) FromJson(data []byte) error {
	return json.Unmarshal(data, v)
}

func (v *SecretManifest) ToJson() ([]byte, error) {
	return json.Marshal(v)
}

func (v *SecretManifest) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, v)
}

func (v *SecretManifest) ToYaml() ([]byte, error) {
	return yaml.Marshal(v)
}

func (v *SecretManifest) SetSecretMeta(cfg *types.Secret) {

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

func (v *SecretManifest) SetAuthData(username, password string) {
	v.Spec.Data[types.SecretUsernameKey] = username
	v.Spec.Data[types.SecretPasswordKey] = password
}

// SetSecretSpec - set config spec from manifest
// TODO: check if config spec is updated => update Meta.Updated or skip
func (v *SecretManifest) SetSecretSpec(s *types.Secret) {

	s.Spec.Type = v.Spec.Type
	s.Spec.Data = make(map[string][]byte, 0)

	for key, value := range v.Spec.Data {
		s.Spec.Data[key] = []byte(base64.StdEncoding.EncodeToString([]byte(value)))
	}
}

func (v *SecretManifest) GetManifest() *types.SecretManifest {
	cfg := new(types.SecretManifest)
	cfg.Type = v.Spec.Type
	return cfg
}

// swagger:ignore
type SecretRemoveOptions struct {
	Force bool `json:"force"`
}
