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
	"fmt"
	"time"
)

const (
	KindConfigText = "text"
)

// swagger:ignore
// swagger:model types_config
type Config struct {
	System
	Meta ConfigMeta `json:"meta" yaml:"meta"`
	Spec ConfigSpec `json:"spec" yaml:"spec"`
}

// swagger:ignore
type ConfigList struct {
	System
	Items []*Config
}

// swagger:ignore
type ConfigMap struct {
	System
	Items map[string]*Config
}

// swagger:ignore
// swagger:model types_config_meta
type ConfigMeta struct {
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Meta      `yaml:",inline"`
}

type ConfigSpec struct {
	Type string            `json:"type" yaml:"type"`
	Data map[string]string `json:"data" yaml:"data"`
}

type ConfigManifest struct {
	System
	State   string            `json:"state"`
	Type    string            `json:"kind"`
	Data    map[string]string `json:"data" yaml:"data"`
	Created time.Time         `json:"created"`
	Updated time.Time         `json:"updated"`
}

type ConfigText struct {
	Text string `json:"text"`
}

type ConfigFile struct {
	Files map[string][]byte `json:"text"`
}

func (c *ConfigManifest) Set(cfg *Config) {
	c.Type = cfg.Spec.Type
	c.Data = cfg.Spec.Data
}

func (s *Config) GetHash() string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s", s.Meta.Name)))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (s *Config) SelfLink() string {
	if s.Meta.SelfLink == "" {
		s.Meta.SelfLink = s.CreateSelfLink(s.Meta.Namespace, s.Meta.Name)
	}
	return s.Meta.SelfLink
}

func (s *Config) CreateSelfLink(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
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
