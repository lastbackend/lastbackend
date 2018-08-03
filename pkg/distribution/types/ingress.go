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

import "fmt"

type IngressList struct {
	Runtime
	Items []*Ingress
}

// swagger:ignore
type Ingress struct {
	Runtime
	Meta   IngressMeta   `json:"meta"`
	Status IngressStatus `json:"status"`
	Spec   IngressSpec   `json:"spec"`
}

// swagger:ignore
type IngressMeta struct {
	Meta
	Cluster string `json:"cluster"`
}

// swagger:model types_ingress_status
type IngressStatus struct {
	Ready bool `json:"ready"`
}

// swagger:ignore
type IngressSpec struct {
	Routes map[string]RouteSpec `json:"routes"`
}

// swagger:ignore
type IngressCreateMetaOptions struct {
	MetaCreateOptions
}

// swagger:model types_ingress_meta_update
type IngressUpdateMetaOptions struct {
	MetaUpdateOptions
}

// swagger:ignore
type IngressCreateOptions struct {
	Meta    IngressCreateMetaOptions `json:"meta",yaml:"meta"`
	Status  IngressStatus            `json:"status",yaml:"status"`
	Network NetworkSpec              `json:"network"`
}

func (m *IngressMeta) Set(meta *IngressUpdateMetaOptions) {
	if meta.Description != nil {
		m.Description = *meta.Description
	}

	if meta.Labels != nil {
		m.Labels = meta.Labels
	}

}

func (n *Ingress) SelfLink() string {
	if n.Meta.SelfLink == "" {
		n.Meta.SelfLink = fmt.Sprintf("%s:%s", n.Meta.Cluster, n.Meta.Name)
	}
	return n.Meta.SelfLink
}
