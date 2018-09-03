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

// swagger:ignore
type Ingress struct {
	Runtime
	Meta   IngressMeta   `json:"meta"`
	Status IngressStatus `json:"status"`
	Spec   IngressSpec   `json:"spec"`
}

type IngressList struct {
	Runtime
	Items []*Ingress
}

type IngressMap struct {
	Runtime
	Items map[string]*Ingress
}

// swagger:ignore
type IngressMeta struct {
	Meta
	Node string `json:"node"`
}

// swagger:model types_ingress_status
type IngressStatus struct {
	Ready bool `json:"ready"`
}

// swagger:ignore
type IngressSpec struct {
}

func (n *Ingress) SelfLink() string {
	if n.Meta.SelfLink == "" {
		n.Meta.SelfLink = fmt.Sprintf("%s", n.Meta.Node)
	}
	return n.Meta.SelfLink
}

func NewIngressList() *IngressList {
	dm := new(IngressList)
	dm.Items = make([]*Ingress, 0)
	return dm
}

func NewIngressMap() *IngressMap {
	dm := new(IngressMap)
	dm.Items = make(map[string]*Ingress)
	return dm
}
