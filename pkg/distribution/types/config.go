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
// patents in process, and are protected by trade config or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package types

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

const (
	KindConfigText = "text"
	KindConfigFile = "file"
)

// swagger:ignore
// swagger:model types_config
type Config struct {
	Runtime
	Meta ConfigMeta        `json:"meta" yaml:"meta"`
	Data map[string][]byte `json:"data" yaml:"data"`
}

// swagger:ignore
type ConfigList struct {
	Runtime
	Items []*Config
}

// swagger:ignore
type ConfigMap struct {
	Runtime
	Items map[string]*Config
}

// swagger:ignore
// swagger:model types_config_meta
type ConfigMeta struct {
	Kind string `json:"kind"`
	Meta `yaml:",inline"`
}

type ConfigManifest struct {
	Runtime
	State   string            `json:"state"`
	Kind    string            `json:"kind"`
	Data    map[string][]byte `json:"data" yaml:"data"`
	Created time.Time         `json:"created"`
	Updated time.Time         `json:"updated"`
}

type ConfigManifestList struct {
	Runtime
	Items []*ConfigManifest
}

type ConfigManifestMap struct {
	Runtime
	Items map[string]*ConfigManifest
}

func NewConfigManifestList() *ConfigManifestList {
	dm := new(ConfigManifestList)
	dm.Items = make([]*ConfigManifest, 0)
	return dm
}

func NewConfigManifestMap() *ConfigManifestMap {
	dm := new(ConfigManifestMap)
	dm.Items = make(map[string]*ConfigManifest)
	return dm
}

func (s *Config) DecodeConfigTextData(key string) (string, error) {

	if s.Meta.Kind != KindConfigText {
		return EmptyString, errors.New("invalid config type")
	}

	if _, ok := s.Data[key]; !ok {
		return EmptyString, errors.New("config key not found")
	}

	d, err := base64.StdEncoding.DecodeString(string(s.Data[key]))
	if err != nil {
		return EmptyString, err
	}

	return string(d), nil

}

type ConfigText struct {
	Text string `json:"text"`
}

type ConfigFile struct {
	Files map[string][]byte `json:"text"`
}

func (s *Config) GetHash() string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s", s.Meta.Name)))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (s *Config) SelfLink() string {
	if s.Meta.SelfLink == "" {
		s.Meta.SelfLink = s.CreateSelfLink(s.Meta.Name)
	}
	return s.Meta.SelfLink
}

func (s *Config) CreateSelfLink(name string) string {
	return fmt.Sprintf("%s", name)
}

func (s *Config) DecodeRegistry() {

}

// swagger:ignore
type ConfigCreateOptions struct {
	Name string
	Kind string
	Data map[string][]byte
}

// swagger:ignore
type ConfigUpdateOptions struct {
	Kind string
	Data map[string][]byte
}

// swagger:ignore
type ConfigRemoveOptions struct {
	Force bool `json:"force"`
}

func NewConfigList() *ConfigList {
	dm := new(ConfigList)
	dm.Items = make([]*Config, 0)
	return dm
}

func NewConfigMap() *ConfigMap {
	dm := new(ConfigMap)
	dm.Items = make(map[string]*Config)
	return dm
}
