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

package views

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"time"
)

// swagger:model views_secret
type Secret struct {
	Meta SecretMeta        `json:"meta"`
	Data map[string][]byte `json:"data"`
}

// swagger:model views_secret_meta
type SecretMeta struct {
	Name     string    `json:"name"`
	Kind     string    `json:"kind"`
	SelfLink string    `json:"self_link"`
	Updated  time.Time `json:"updated"`
	Created  time.Time `json:"created"`
}

// swagger:ignore
type SecretMap map[string]*Secret

// swagger:model views_secret_list
type SecretList []*Secret

func (s *Secret) Decode() *types.Secret {

	o := new(types.Secret)
	o.Meta.Name = s.Meta.Name
	o.Meta.Kind = s.Meta.Kind
	o.Meta.SelfLink = s.Meta.SelfLink
	o.Meta.Updated = s.Meta.Updated
	o.Meta.Created = s.Meta.Created

	o.Data = s.Data

	return o
}
