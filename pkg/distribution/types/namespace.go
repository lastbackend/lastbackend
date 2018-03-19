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
	"encoding/json"
	"fmt"
)

type NamespaceList []*Namespace

type Namespace struct {
	Meta      NamespaceMeta      `json:"meta"`
	Env       NamespaceEnvs      `json:"env"`
	Resources NamespaceResources `json:"resources"`
	Quotas    NamespaceQuotas    `json:"quotas,omitempty"`
	Labels    map[string]string  `json:"labels"`
}

type NamespaceEnvs []NamespaceEnv

type NamespaceEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type NamespaceMeta struct {
	Meta     `yaml:",inline"`
	Endpoint string `json:"endpoint"`
	Type     string `json:"type"`
}

type NamespaceResources struct {
	RAM    int64 `json:"ram"`
	Routes int   `json:"routes"`
}

type NamespaceQuotas struct {
	RAM      int64 `json:"ram"`
	Routes   int   `json:"routes"`
	Disabled bool  `json:"disabled"`
}

func (n *Namespace) SelfLink() string {
	if n.Meta.SelfLink == "" {
		n.Meta.SelfLink = fmt.Sprintf("%s", n.Meta.Name)
	}
	return n.Meta.SelfLink
}

func (n *Namespace) ToJson() ([]byte, error) {
	buf, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (n *NamespaceList) ToJson() ([]byte, error) {
	if n == nil {
		return []byte("[]"), nil
	}
	buf, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
