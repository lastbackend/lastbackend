//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package v1

import (
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/pod/views/v1"
	"strings"
	"time"
)

type Service struct {
	Meta ServiceMeta  `json:"meta"`
	Pods []v1.PodInfo `json:"pods,omitempty"`
	Spec []SpecInfo   `json:"spec,omitempty"`
}

type ServiceMeta struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Region      string    `json:"region"`
	Replicas    int       `json:"replicas,omitempty"`
	Namespace   string    `json:"namespace"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

type SpecInfo struct {
	Meta    SpecMeta `json:"meta,omitempty"`
	Memory  int64    `json:"memory,omitempty"`
	Command string   `json:"command,omitempty"`
	Image   string   `json:"image,omitempty"`
	Region  string   `json:"region,omitempty"`
}

type SpecMeta struct {
	// Meta id
	ID string `json:"id"`
	// Meta labels
	Labels map[string]string `json:"labels"`
	// Meta created time
	Created time.Time `json:"created"`
	// Meta updated time
	Updated time.Time `json:"updated"`
}

type ServiceList []*Service

func ToSpecInfo(spec *types.ServiceSpec) SpecInfo {
	info := SpecInfo{
		Meta:    ToSpecMeta(spec.Meta),
		Memory:  spec.Memory,
		Command: strings.Join(spec.Command, " "),
	}

	return info
}

func ToSpecMeta(meta types.SpecMeta) SpecMeta {
	m := SpecMeta{
		ID:      meta.ID,
		Labels:  meta.Labels,
		Created: meta.Created,
		Updated: meta.Updated,
	}

	if len(m.Labels) == 0 {
		m.Labels = make(map[string]string)
	}

	return m
}
