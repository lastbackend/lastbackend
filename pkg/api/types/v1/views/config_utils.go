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
	"encoding/json"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
)

type ConfigView struct{}

func (sv *ConfigView) New(obj *types.Config) *Config {
	s := Config{}
	s.Meta = s.ToMeta(obj.Meta)
	s.Spec = s.ToSpec(obj.Spec)
	return &s
}

func (s *Config) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Config) ToMeta(obj types.ConfigMeta) ConfigMeta {
	meta := ConfigMeta{}
	meta.Name = obj.Name
	meta.SelfLink = obj.SelfLink.String()
	meta.Namespace = obj.Namespace
	meta.Kind = obj.Kind
	meta.Updated = obj.Updated
	meta.Created = obj.Created
	return meta
}

func (s *Config) ToSpec(obj types.ConfigSpec) ConfigSpec {

	spec := ConfigSpec{}
	spec.Data = make(map[string]string, 0)

	for key, val := range obj.Data {
		spec.Data[key] = val
	}

	return spec
}

func (sv ConfigView) NewList(obj *types.ConfigList) *ConfigList {
	if obj == nil {
		return nil
	}

	sl := make(ConfigList, 0)
	for _, v := range obj.Items {
		sl = append(sl, sv.New(v))
	}
	return &sl
}

func (sl *ConfigList) ToJson() ([]byte, error) {
	if sl == nil {
		sl = &ConfigList{}
	}
	return json.Marshal(sl)
}
