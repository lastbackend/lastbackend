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
	"errors"
)

const (
	KindSecretText     = "text"
	KindSecretRegistry = "registry"
	KindSecretFiles    = "files"
)

// swagger:ignore
// swagger:model types_secret
type Secret struct {
	Runtime
	Meta SecretMeta        `json:"meta" yaml:"meta"`
	Data map[string][]byte `json:"data" yaml:"data"`
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
	Kind string `json:"kind"`
	Meta `yaml:",inline"`
}

func (s *Secret) EncodeSecretRegistryData(d SecretRegistryData) {
	s.Data["registry"] = []byte(base64.StdEncoding.EncodeToString([]byte(d.Registry)))
	s.Data["username"] = []byte(base64.StdEncoding.EncodeToString([]byte(d.Username)))
	s.Data["password"] = []byte(base64.StdEncoding.EncodeToString([]byte(d.Password)))
}

func (s *Secret) DecodeSecretRegistryData() (*SecretRegistryData, error) {

	if s.Meta.Kind != KindSecretRegistry {
		return nil, errors.New("Invalid secret type")
	}

	data := new(SecretRegistryData)

	r, err := base64.StdEncoding.DecodeString(string(s.Data["registry"]))
	if err != nil {
		return nil, err
	}
	data.Registry = string(r)

	u, err := base64.StdEncoding.DecodeString(string(s.Data["username"]))
	if err != nil {
		return nil, err
	}
	data.Username = string(u)

	p, err := base64.StdEncoding.DecodeString(string(s.Data["password"]))
	if err != nil {
		return nil, err
	}
	data.Password = string(p)

	return data, nil
}

type SecretRegistryData struct {
	Registry string `json:"registry"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SecretText struct {
	Text string `json:"text"`
}

type SecretFile struct {
	Files map[string][]byte `json:"text"`
}



func (s *Secret) GetHash() string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s", s.Meta.Name)))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (s *Secret) SelfLink() string {
	if s.Meta.SelfLink == "" {
		s.Meta.SelfLink = s.CreateSelfLink(s.Meta.Name)
	}
	return s.Meta.SelfLink
}

func (s *Secret) CreateSelfLink(name string) string {
	return fmt.Sprintf("%s", name)
}

func (s *Secret) DecodeRegistry() {

}


// swagger:ignore
type SecretCreateOptions struct {
	Name string
	Kind string
	Data map[string][]byte
}

// swagger:ignore
type SecretUpdateOptions struct {
	Kind string
	Data map[string][]byte
}

// swagger:ignore
type SecretRemoveOptions struct {
	Force bool `json:"force"`
}

func NewSecretList() *SecretList {
	dm := new(SecretList)
	dm.Items = make([]*Secret, 0)
	return dm
}

func NewSecretMap() *SecretMap {
	dm := new(SecretMap)
	dm.Items = make(map[string]*Secret)
	return dm
}
