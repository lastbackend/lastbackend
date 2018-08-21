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

package types

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

// swagger:ignore
// swagger:model types_secret
type Secret struct {
	Runtime
	Meta SecretMeta `json:"meta" yaml:"meta"`
	Data string     `json:"data" yaml:"data"`
}

// swagger:ignore
type SecretList struct {
	Runtime
	Items []*Secret
}

// swagger:ignore
type SecretMap struct {
	Runtime
	Items map[string]*Secret
}

// swagger:ignore
// swagger:model types_secret_meta
type SecretMeta struct {
	Meta      `yaml:",inline"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

func (s *Secret) GetHash() string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s:%s", s.Meta.Namespace, s.Meta.Name)))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (s *Secret) SelfLink() string {
	if s.Meta.SelfLink == "" {
		s.Meta.SelfLink = s.CreateSelfLink(s.Meta.Namespace, s.Meta.Name)
	}
	return s.Meta.SelfLink
}

func (s *Secret) CreateSelfLink(namespace, name string) string {
	return fmt.Sprintf("%s:%s", namespace, name)
}

// swagger:ignore
type SecretCreateOptions struct {
	Data *string
}

// swagger:ignore
type SecretUpdateOptions struct {
	Data *string
}

// swagger:ignore
type SecretRemoveOptions struct {
	Force bool `json:"force"`
}

func NewSecretList () *SecretList {
	dm := new(SecretList)
	dm.Items = make([]*Secret, 0)
	return dm
}

func NewSecretMap () *SecretMap {
	dm := new(SecretMap)
	dm.Items = make(map[string]*Secret)
	return dm
}