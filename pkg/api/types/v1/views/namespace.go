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
	"time"
)

// swagger:model views_namespace
type Namespace struct {
	Meta NamespaceMeta `json:"meta"`
	Spec NamespaceSpec `json:"spec"`
}

// swagger:model views_namespace_meta
type NamespaceMeta struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	SelfLink    string            `json:"self_link"`
	Endpoint    string            `json:"endpoint"`
	Suffix      string            `json:"suffix"`
	Labels      map[string]string `json:"labels"`
	Created     time.Time         `json:"created"`
	Updated     time.Time         `json:"updated"`
}

// swagger:model views_namespace_spec
type NamespaceSpec struct {
	Env       NamespaceEnvs      `json:"env"`
	Resources NamespaceResources `json:"resources"`
	Quotas    NamespaceQuotas    `json:"quotas"`
}

// swagger:model views_namespace_envs
type NamespaceEnvs []string

// swagger:model views_namespace_resource
type NamespaceResources struct {
	RAM    int64 `json:"ram"`
	Routes int   `json:"routes"`
}

// swagger:model views_namespace_quotas
type NamespaceQuotas struct {
	Disabled bool  `json:"disabled"`
	RAM      int64 `json:"ram"`
	Routes   int   `json:"routes"`
}

// swagger:ignore
// swagger:model views_namespace_resource
type NamespaceResource struct {
}

// swagger:model views_namespace_list
type NamespaceList []*Namespace
