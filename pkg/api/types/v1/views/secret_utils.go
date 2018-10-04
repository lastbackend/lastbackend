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

func (sv *SecretView) New(obj *types.Secret) *Secret {
	s := Secret{}
	s.Meta = s.ToMeta(obj.Meta)
	s.Spec = s.ToSpec(obj.Spec)
	return &s
}

func (s *Secret) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Secret) ToMeta(obj types.SecretMeta) SecretMeta {
	meta := SecretMeta{}
	meta.Name = obj.Name
	meta.SelfLink = obj.SelfLink
	meta.Namespace = obj.Namespace
	meta.Updated = obj.Updated
	meta.Created = obj.Created
	return meta
}

func (s *Secret) ToSpec(obj types.SecretSpec) SecretSpec {
	spec := SecretSpec{}
	spec.Type = obj.Type
	spec.Data = make(map[string][]byte, 0)
	for key, value := range obj.Data {
		spec.Data[key]=value
	}
	return spec
}

func (sv SecretView) NewList(obj *types.SecretList) *SecretList {
	if obj == nil {
		return nil
	}

	sl := make(SecretList, 0)
	for _, v := range obj.Items {
		sl = append(sl, sv.New(v))
	}
	return &sl
}

func (sl *SecretList) ToJson() ([]byte, error) {
	if sl == nil {
		sl = &SecretList{}
	}
	return json.Marshal(sl)
}
