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

package views

import (
	"github.com/lastbackend/lastbackend/internal/pkg/types"
	"time"
)

// swagger:model views_secret
type Secret struct {
	Meta SecretMeta `json:"meta"`
	Spec SecretSpec `json:"spec"`
}

type SecretSpec struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

// swagger:model views_secret_meta
type SecretMeta struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	SelfLink  string    `json:"self_link"`
	Updated   time.Time `json:"updated"`
	Created   time.Time `json:"created"`
}

// swagger:ignore
type SecretMap map[string]*Secret

// swagger:model views_secret_list
type SecretList []*Secret

func (s *Secret) Decode() *types.Secret {

	o := new(types.Secret)
	o.Meta.Name = s.Meta.Name

	o.Meta.SelfLink = *types.NewSecretSelfLink(s.Meta.Namespace, s.Meta.Name)
	o.Meta.Updated = s.Meta.Updated
	o.Meta.Created = s.Meta.Created

	o.Spec.Type = s.Spec.Type

	o.Spec.Data = make(map[string][]byte, 0)
	for k, v := range s.Spec.Data {
		o.Spec.Data[k] = []byte(v)
	}

	return o
}
