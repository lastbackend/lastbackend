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
	Meta   NamespaceMeta   `json:"meta"`
	Status NamespaceStatus `json:"status"`
	Spec   NamespaceSpec   `json:"spec"`
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
	Domain    NamespaceDomain    `json:"domain"`
	Resources NamespaceResources `json:"resources"`
}

type NamespaceStatus struct {
	Resources NamespaceStatusResources `json:"resources"`
}

type NamespaceStatusResources struct {
	Allocated NamespaceResource `json:"allocated"`
}

// swagger:model views_namespace_envs
type NamespaceEnvs []string

type NamespaceResources struct {
	Request NamespaceResource `json:"request"`
	Limits  NamespaceResource `json:"limits"`
}

// swagger:model views_namespace_resource
type NamespaceResource struct {
	RAM     string `json:"ram"`
	CPU     string `json:"cpu"`
	Storage string `json:"storage"`
}

type NamespaceDomain struct {
	Internal string `json:"internal"`
	External string `json:"external"`
}

// swagger:model views_namespace_list
type NamespaceList []*Namespace

type NamespaceApplyStatus struct {
	Configs  map[string]bool `json:"configs,omitempty"`
	Secrets  map[string]bool `json:"secrets,omitempty"`
	Volumes  map[string]bool `json:"volumes,omitempty"`
	Services map[string]bool `json:"services,omitempty"`
	Routes   map[string]bool `json:"routes,omitempty"`
	Jobs     map[string]bool `json:"jobs,omitempty"`
}
