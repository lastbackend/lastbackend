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
	"encoding/json"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type SecretView struct{}

func (rv *SecretView) New(obj *types.Secret) *Secret {
	r := Secret{}
	r.Meta = r.ToMeta(obj.Meta)
	r.Spec = r.ToSpec(obj.Spec)
	r.State = r.ToState(obj.State)
	return &r
}

func (p *Secret) ToJson() ([]byte, error) {
	return json.Marshal(p)
}

func (r *Secret) ToMeta(obj types.SecretMeta) SecretMeta {
	meta := SecretMeta{}
	meta.Name = obj.Name
	meta.Namespace = obj.Namespace
	meta.Updated = obj.Updated
	meta.Created = obj.Created

	return meta
}

func (r *Secret) ToSpec(obj types.SecretSpec) SecretSpec {
	spec := SecretSpec{}
	return spec
}

func (r *Secret) ToState(obj types.SecretState) SecretState {
	state := SecretState{}
	return state
}

func (rv SecretView) NewList(obj map[string]*types.Secret) *SecretList {
	if obj == nil {
		return nil
	}

	n := make(SecretList, 0)
	for _, v := range obj {
		n[v.Meta.Name] = rv.New(v)
	}
	return &n
}

func (n *SecretList) ToJson() ([]byte, error) {
	if n == nil {
		n = &SecretList{}
	}
	return json.Marshal(n)
}
