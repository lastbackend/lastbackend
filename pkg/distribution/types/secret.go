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

package types

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

const (
	KindSecretOpaque = "opaque"
	KindSecretAuth   = "auth"

	SecretUsernameKey = "username"
	SecretPasswordKey = "password"
)

// swagger:ignore
// swagger:model types_secret
type Secret struct {
	System
	Meta SecretMeta `json:"meta" yaml:"meta"`
	Spec SecretSpec `json:"spec" yaml:"yaml"`
}

// swagger:ignore
type SecretList struct {
	System
	Items []*Secret
}

// swagger:ignore
type SecretMap struct {
	System
	Items map[string]*Secret
}

// swagger:ignore
// swagger:model types_secret_meta
type SecretMeta struct {
	Meta      `yaml:",inline"`
	SelfLink  SecretSelfLink `json:"self_link"`
	Namespace string         `json:"namespace" yaml:"namespace"`
}

type SecretSpec struct {
	Type string            `json:"type"`
	Data map[string][]byte `json:"data" yaml:"data"`
}

type SecretManifest struct {
	System
	State   string    `json:"state"`
	Type    string    `json:"type"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type SecretManifestList struct {
	System
	Items []*SecretManifest
}

type SecretManifestMap struct {
	System
	Items map[string]*SecretManifest
}

func NewSecretManifestList() *SecretManifestList {
	dm := new(SecretManifestList)
	dm.Items = make([]*SecretManifest, 0)
	return dm
}

func NewSecretManifestMap() *SecretManifestMap {
	dm := new(SecretManifestMap)
	dm.Items = make(map[string]*SecretManifest)
	return dm
}

func (s *Secret) EncodeSecretAuthData(d SecretAuthData) {
	s.Spec.Data = make(map[string][]byte)
	s.Spec.Data["username"] = []byte(base64.StdEncoding.EncodeToString([]byte(d.Username)))
	s.Spec.Data["password"] = []byte(base64.StdEncoding.EncodeToString([]byte(d.Password)))
}

func (s *Secret) DecodeSecretAuthData() (*SecretAuthData, error) {

	if s.Spec.Type != KindSecretAuth {
		return nil, errors.New("invalid secret type")
	}

	data := new(SecretAuthData)

	u, err := base64.StdEncoding.DecodeString(string(s.Spec.Data["username"]))
	if err != nil {
		return nil, err
	}
	data.Username = string(u)

	p, err := base64.StdEncoding.DecodeString(string(s.Spec.Data["password"]))
	if err != nil {
		return nil, err
	}
	data.Password = string(p)

	return data, nil
}

func (s *Secret) DecodeSecretTextData(key string) (string, error) {

	if s.Spec.Type != KindSecretOpaque {
		return EmptyString, errors.New("invalid secret type")
	}

	if _, ok := s.Spec.Data[key]; !ok {
		return EmptyString, errors.New("secret key not found")
	}

	fmt.Println("key:", key)

	d, err := base64.StdEncoding.DecodeString(string(s.Spec.Data[key]))
	if err != nil {
		return EmptyString, err
	}

	fmt.Println("data", string(d))

	return string(d), nil

}

type SecretAuthData struct {
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

func (s *Secret) SelfLink() *SecretSelfLink {
	return &s.Meta.SelfLink
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
