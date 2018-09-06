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

// swagger:ignore
type NamespaceMap struct {
	Runtime
	Items map[string]*Namespace
}

// swagger:ignore
type NamespaceList struct {
	Runtime
	Items []*Namespace
}

// swagger:ignore
type Namespace struct {
	Meta NamespaceMeta `json:"meta"`
	Spec NamespaceSpec `json:"spec"`
}

// swagger:ignore
type NamespaceEnvs []NamespaceEnv

// swagger:ignore
type NamespaceEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// swagger:ignore
// swagger:model types_namespace_meta
type NamespaceMeta struct {
	Meta     `yaml:",inline"`
	Endpoint string `json:"endpoint"`
	Type     string `json:"type"`
}

// swagger:ignore
type NamespaceSpec struct {
	Quotas    NamespaceQuotas    `json:"quotas"`
	Resources NamespaceResources `json:"resources"`
	Env       NamespaceEnvs      `json:"env"`
	Domain    NamespaceDomain    `json:"domain"`
}

type NamespaceDomain struct {
	Internal string `json:"internal"`
	External string `json:"external"`
}

// swagger:ignore
type NamespaceQuotas struct {
	RAM      int64 `json:"ram"`
	Routes   int   `json:"routes"`
	Disabled bool  `json:"disabled"`
}

// swagger:ignore
type NamespaceResources struct {
	RAM    int64 `json:"ram"`
	Routes int   `json:"routes"`
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

// swagger:ignore
type NamespaceCreateOptions struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Domain      *string                 `json:"domain"`
	Quotas      *NamespaceQuotasOptions `json:"quotas"`
}

// swagger:ignore
type NamespaceUpdateOptions struct {
	Description *string                 `json:"description"`
	Domain      *string                 `json:"domain"`
	Quotas      *NamespaceQuotasOptions `json:"quotas"`
}

// swagger:ignore
type NamespaceRemoveOptions struct {
	Force bool `json:"force"`
}

// swagger:ignore
type NamespaceQuotasOptions struct {
	Disabled bool  `json:"disabled"`
	RAM      int64 `json:"ram"`
	Routes   int   `json:"routes"`
}

func NewNamespaceList() *NamespaceList {
	dm := new(NamespaceList)
	dm.Items = make([]*Namespace, 0)
	return dm
}

func NewNamespaceMap() *NamespaceMap {
	dm := new(NamespaceMap)
	dm.Items = make(map[string]*Namespace)
	return dm
}
