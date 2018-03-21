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

type Namespace struct {
	Meta      NamespaceMeta      `json:"meta"`
	Env       NamespaceEnvs      `json:"env"`
	Resources NamespaceResources `json:"resources"`
	Quotas    NamespaceQuotas    `json:"quotas"`
}

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

type NamespaceEnvs []string

type NamespaceResources struct {
	RAM    int64 `json:"ram"`
	Routes int   `json:"routes"`
}

type NamespaceQuotas struct {
	Disabled bool  `json:"disabled"`
	RAM      int64 `json:"ram"`
	Routes   int   `json:"routes"`
}

type NamespaceResource struct {
}

type NamespaceList []*Namespace
